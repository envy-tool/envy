package configs

import (
	"envy.pw/cli/internal/addrs"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
)

// Command represents a single Command block in a configuration.
type Command struct {
	Name      string
	DeclRange hcl.Range

	// Executable and CommandLine are mutually-exclusive.
	//
	// Executable sets only the prefix of the command to run, taking all of the
	// other arguments from the envy command line verbatim.
	//
	// CommandLine, on the other hand, sets the full sequence of command line
	// arguments, overriding whatever was present on the command line. However,
	// the expression may still refer to what was given on the command line,
	// allowing the resulting command to be constructed dynamically.
	Executable  hcl.Expression
	CommandLine hcl.Expression

	// Environment is a mapping of environment variables to set when launching
	// the command.
	//
	// InheritEnvironment is a boolean flag, defauting to true, that controls
	// whether the command also inherits the environment variables available
	// to the envy process.
	Environment        hcl.Expression
	InheritEnvironment hcl.Expression

	// WorkDir, if set, decides which directory will be used as the current
	// working directory when launching the command. By default the command
	// inherits the current working directory from when envy itself was
	// launched.
	WorkDir hcl.Expression

	// OnUpdate selects a behavior for when the Executable, CommandLine,
	// Environment, or InheritEnvironment expression results change while the
	// command is already running. We can't directly change any of these while
	// the program is running, so for a program to react to these it must
	// be signalled some other way.
	OnUpdate ProcessAction

	// OnError selects a behavior for when a dependency of Executable,
	// CommandLine, Environment, or InheritEnvironment enters error state,
	// which effectively puts this command also in an error state.
	OnError ProcessAction

	// Dependencies is a collection of references to other objects that
	// must exist and be active for the command to function, even though
	// they are not referenced in any of the other configuration expressions.
	Dependencies []Reference
}

// Addr returns the address for the command that was declared.
func (c *Command) Addr() addrs.Command {
	return addrs.Command{Name: c.Name}
}

func decodeCommandBlock(block *hcl.Block) (*Command, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	cmd := &Command{
		Name:      block.Labels[0],
		DeclRange: block.DefRange,
	}

	type DecodeCommand struct {
		Executable         hcl.Expression `hcl:"exec"`
		CommandLine        hcl.Expression `hcl:"cmdline"`
		Environment        hcl.Expression `hcl:"env"`
		InheritEnvironment hcl.Expression `hcl:"inherit_env"`
		WorkDir            hcl.Expression `hcl:"work_dir"`
		OnUpdate           hcl.Expression `hcl:"on_update"`
		OnError            hcl.Expression `hcl:"on_error"`
		Dependencies       hcl.Expression `hcl:"depends_on"`
	}
	var decCmd DecodeCommand
	moreDiags := gohcl.DecodeBody(block.Body, nil, &decCmd)
	diags = append(diags, moreDiags...)

	cmd.Executable = decCmd.Executable
	cmd.CommandLine = decCmd.CommandLine
	cmd.Environment = decCmd.Environment
	cmd.InheritEnvironment = decCmd.InheritEnvironment
	cmd.WorkDir = decCmd.WorkDir

	cmd.OnUpdate, moreDiags = decodeProcessAction(decCmd.OnUpdate)
	diags = append(diags, moreDiags...)

	cmd.OnError, moreDiags = decodeProcessAction(decCmd.OnError)
	diags = append(diags, moreDiags...)

	cmd.Dependencies, moreDiags = decodeDependsOn(decCmd.Dependencies)
	diags = append(diags, moreDiags...)

	if !validName(cmd.Name) {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid command name",
			Detail:   "All object names must begin with a letter and contain only letters, digits, and underscores.",
			Subject:  block.LabelRanges[0].Ptr(),
		})
	}

	return cmd, diags
}
