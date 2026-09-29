// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: BUSL-1.1

package e2etest

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/internal/e2e"
)

func TestInitializeCommandSuggestsInit(t *testing.T) {
	fixturePath := filepath.Join("testdata", "empty")
	tf := e2e.NewBinary(t, terraformBin, fixturePath)

	cmd := tf.Cmd("initialize", "-no-color")
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
		`Command "terraform initialize" does not exist`,
		`Use "terraform init" instead.`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in stderr:\n%s", want, got)
		}
	}
}
