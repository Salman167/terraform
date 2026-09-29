// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"strings"

	"github.com/hashicorp/terraform/internal/tfdiags"
)

// FormatCommand helps users who expect the formatting command to be named
// "format" by directing them to Terraform's existing "fmt" command.
type FormatCommand struct {
	Meta
}

func (c *FormatCommand) Run(_ []string) int {
	c.showDiagnostics(tfdiags.Sourceless(
		tfdiags.Error,
		"Command \"terraform format\" does not exist",
		"Use \"terraform fmt\" instead. Run \"terraform fmt -help\" for more information.",
	))
	return 1
}

func (c *FormatCommand) Help() string {
	return strings.TrimSpace(`
Usage: terraform [global options] format

  Terraform's configuration formatting command is named "fmt".

  Run "terraform fmt -help" for information about available formatting
  options.
`)
}

func (c *FormatCommand) Synopsis() string {
	return "Suggest the terraform fmt command"
}
