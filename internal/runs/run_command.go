package runs

import (
	"context"
	"fmt"

	"envy.pw/cli/internal/addrs"
	"envy.pw/cli/internal/configs"
	"envy.pw/cli/internal/graphs"
	"envy.pw/cli/internal/nvdiags"
)

// CommandCall represents a call of a command defined in the configuration.
type CommandCall struct {
	Addr    addrs.Command
	Args    []string
	Environ []string
}

// RunCommand creates all of the necessary context to run the given command
// as configured in the given configuration, and then runs the command.
//
// This function blocks until the command has terminated and all of its
// associated helpers are cleaned up.
func (r *Runner) RunCommand(ctx context.Context, call *CommandCall, cfg *configs.Config) (status int, diags nvdiags.Diagnostics) {
	graph, moreDiags := graphForRunCommand(call, cfg)
	diags = diags.Append(moreDiags)
	if moreDiags.HasErrors() {
		return 126, diags
	}

	fmt.Printf("TODO: execute run with graph:\n%s\n", graph.DebugString())

	return 0, diags
}

func graphForRunCommand(call *CommandCall, cfg *configs.Config) (*graphs.Graph, nvdiags.Diagnostics) {
	var diags nvdiags.Diagnostics
	g := graphs.NewGraph()

	cc, exists := cfg.Commands[call.Addr]
	if !exists {
		diags = diags.Append(nvdiags.Sourceless(
			nvdiags.Error,
			"Command not found",
			fmt.Sprintf("There is no command named %q defined in the configuration.", call.Addr.Name),
		))
		return g, diags
	}

	root := &commandExecNode{
		CommandNode: graphs.CommandNode{
			Addr: call.Addr,
		},
		Config: cc,
	}
	g.AddWithReferents(root, func(referrer addrs.Referenceable, ref configs.Reference) (graphs.Node, nvdiags.Diagnostics) {
		var diags nvdiags.Diagnostics
		switch addr := ref.Addr.(type) {

		case addrs.Helper:
			return makeHelperRunNode(addr, ref.SourceRange, cfg)

		case addrs.Path:
			return nil, nil // No node required for a path

		default:
			// This default error message is not actionable and lacks
			// explanation, so we should try to catch most error cases
			// above and return better error messages where possible.
			diags = diags.Append(nvdiags.WithSource(
				nvdiags.Error,
				"Invalid reference",
				fmt.Sprintf("Cannot refer to %s from %s.", ref.Addr, referrer),
				ref.SourceRange,
			))
			return nil, diags
		}
	})

	return g, diags
}
