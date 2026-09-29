// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: BUSL-1.1

package e2etest

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/internal/e2e"
)

func TestFormatCommandSuggestsFmt(t *testing.T) {
	fixturePath := filepath.Join("testdata", "fmt")
	tf := e2e.NewBinary(t, terraformBin, fixturePath)

	cmd := tf.Cmd("format", "-no-color")
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	exitErr, ok := err.(*exec.ExitError)
	if !ok || exitErr.ExitCode() != 1 {
		t.Fatalf("expected exit status 1, got %v", err)
	}
	if stdout.Len() != 0 {
		t.Errorf("expected empty stdout, got:\n%s", stdout)
	}

	got := stderr.String()
	for _, want := range []string{
		`Command "terraform format" does not exist`,
		`Use "terraform fmt" instead.`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in stderr:\n%s", want, got)
		}
	}
}

// Reproduction of the scenario reported in https://github.com/hashicorp/terraform/issues/39299
func TestFmt_errorWritingToFile(t *testing.T) {
	switch runtime.GOOS {
	case "darwin", "linux":
	default:
		t.Skipf("test requires a Unix shell; `sh` unsupported on %s", runtime.GOOS)
	}

	fixturePath := filepath.Join("testdata", "fmt")
	tf := e2e.NewBinary(t, terraformBin, fixturePath)

	// Assert that main.tf has content before running fmt
	mainPath := filepath.Join(tf.WorkDir(), "main.tf")
	content, err := os.ReadFile(mainPath)
	if err != nil {
		t.Fatalf("unexpected error reading test file: %s", err)
	}
	if len(content) == 0 {
		t.Fatal("expected main.tf to contain config, but it is empty")
	}

	// Assert that there's a formatting issue present
	// With `-check`, this is confirmed by error code 3.
	preFmtCmd := tf.Cmd("fmt", "-check", "-no-color")
	err = preFmtCmd.Run()
	if err.Error() != "exit status 3" {
		t.Fatalf("expected exit status 3 error, got: %s", err)
	}

	// The fmt command we're testing ulimit with, which will attempt to write to main.tf
	fmtCmd := tf.Cmd("fmt", "-no-color")
	fmtCmd.Stdin = nil
	fmtCmd.Stdout = &bytes.Buffer{}
	fmtCmd.Stderr = &bytes.Buffer{}

	// But, we need to wrap the command above in a shell command to enforce `ulimit -f 0`.
	// This:
	//   * Causes an error in fmt when writing formatted content to the file,
	//     resulting in the file being left empty.
	//   * Only impacts this command and not the entire test process.
	cmd := exec.Command(
		"/bin/sh", "-c",
		`ulimit -f 0; exec "$@"`,
		"sh", fmtCmd.Path,
	)
	cmd.Args = append(cmd.Args, fmtCmd.Args[1:]...)
	cmd.Dir = fmtCmd.Dir
	cmd.Env = fmtCmd.Env
	cmd.Stdin, cmd.Stdout, cmd.Stderr = fmtCmd.Stdin, fmtCmd.Stdout, fmtCmd.Stderr

	err = cmd.Run()
	if err == nil {
		t.Fatal("expected error when writing to file with ulimit -f 0, but got none")
	}

	stderr := cmd.Stderr.(*bytes.Buffer).String()
	expectErr := "Error: Failed to write main.tf"
	if !strings.Contains(stderr, expectErr) {
		t.Fatalf("expected stderr to contain '%s', but got: %s", expectErr, stderr)
	}

	// Finally, confirm that the error made main.tf empty
	content, err = os.ReadFile(mainPath)
	if err != nil {
		t.Fatalf("unexpected error reading test file: %s", err)
	}
	if len(content) != 0 {
		t.Fatalf("expected main.tf to be empty after error, but got: %s", string(content))
	}
}
