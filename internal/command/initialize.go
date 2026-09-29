// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"strings"

	"github.com/hashicorp/terraform/internal/tfdiags"
)

// InitializeCommand helps users who expect the initialization command to be
// named "initialize" by directing them to Terraform's existing "init" command.
type InitializeCommand struct {
	Meta
}

func (c *InitializeCommand) Run(_ []string) int {
	c.showDiagnostics(tfdiags.Sourceless(
		tfdiags.Error,
		"Command \"terraform initialize\" does not exist",
		"Use \"terraform init\" instead. Run \"terraform init -help\" for more information.",
	))
	return 1
}

func (c *InitializeCommand) Help() string {
	return strings.TrimSpace(`
Usage: terraform [global options] initialize

  Terraform's working-directory initialization command is named "init".

  Run "terraform init -help" for information about available initialization
  options.
`)
}

func (c *InitializeCommand) Synopsis() string {
	return "Suggest the terraform init command"
}
