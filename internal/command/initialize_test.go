// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"strings"
	"testing"
)

func TestInitializeCommand(t *testing.T) {
	ui := testUiWrapped(t)
	command := &InitializeCommand{
		Meta: Meta{Ui: ui},
	}

	if code := command.Run(nil); code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}

	got := ui.ErrorWriter.String()
	for _, want := range []string{
		`Command "terraform initialize" does not exist`,
		`Use "terraform init" instead.`,
		`terraform init -help`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in diagnostic:\n%s", want, got)
		}
	}
}
