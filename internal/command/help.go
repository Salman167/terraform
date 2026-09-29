// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"strings"

	"github.com/hashicorp/terraform/internal/tfdiags"
)

// HelpCommand helps users who run "terraform help" by directing them to the
// existing global -help flag. See https://github.com/hashicorp/terraform/issues/26967
type HelpCommand struct {
	Meta
}

func (c *HelpCommand) Run(_ []string) int {
	c.showDiagnostics(tfdiags.Sourceless(
		tfdiags.Error,
		"Command \"terraform help\" does not exist",
		"Use \"terraform -help\" instead. For a specific command, run \"terraform <command> -help\".",
	))
	return 1
}

func (c *HelpCommand) Help() string {
	return strings.TrimSpace(`
Usage: terraform [global options] help

  Terraform's help is a global flag, not a subcommand.

  Run "terraform -help" to list top-level commands, or
  "terraform <command> -help" for a specific command.
`)
}

func (c *HelpCommand) Synopsis() string {
	return "Suggest the terraform -help flag"
}
