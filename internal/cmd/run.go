package cmd

import (
	"fmt"

	"envy.pw/cli/internal/nvdiags"
)

// runCommand is a command for running commands.
type runCommand struct {
	Context *RunContext
	Args    []string
}

func (c *runCommand) Run() nvdiags.Diagnostics {
	fmt.Printf("Stubbed out 'run' command in %s with args %#v\n", c.Context.WorkingDir, c.Args)
	return nil
}
