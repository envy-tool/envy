package cmd

import (
	"context"
	"fmt"

	"envy.pw/cli/internal/addrs"
	"envy.pw/cli/internal/nvdiags"
	"envy.pw/cli/internal/runs"
)

// runCommand is a command for running commands.
type runCommand struct {
	Context *RunContext
	Args    []string
}

func (c *runCommand) Run() (int, nvdiags.Diagnostics) {
	var diags nvdiags.Diagnostics
	cmdName := c.Args[0]
	if !addrs.ValidName(cmdName) {
		diags = diags.Append(nvdiags.Sourceless(
			nvdiags.Error,
			"Invalid command name",
			fmt.Sprintf("The name %q is not a valid name for a command.", cmdName),
		))
		return 126, diags
	}

	cfg, moreDiags := c.Context.LoadConfig()
	diags = diags.Append(moreDiags)
	if moreDiags.HasErrors() {
		return 126, diags
	}

	runner, moreDiags := c.Context.NewRunner()
	diags = diags.Append(moreDiags)
	if moreDiags.HasErrors() {
		return 126, diags
	}

	call := &runs.CommandCall{
		Addr:    addrs.MakeCommand(cmdName),
		Args:    c.Args[1:],
		Environ: nil, // TODO
	}
	status, moreDiags := runner.RunCommand(context.Background(), call, cfg)
	diags = diags.Append(moreDiags)
	return status, diags
}
