package runs

import (
	"fmt"

	"envy.pw/cli/internal/addrs"
	"envy.pw/cli/internal/configs"
	"envy.pw/cli/internal/graphs"
	"envy.pw/cli/internal/nvdiags"
)

type helperRunNode struct {
	graphs.HelperNode
	Config *configs.Helper
}

func makeHelperRunNode(addr addrs.Helper, rng nvdiags.SourceRange, cfg *configs.Config) (*helperRunNode, nvdiags.Diagnostics) {
	var diags nvdiags.Diagnostics

	hc, exists := cfg.Helpers[addr]
	if !exists {
		diags = diags.Append(nvdiags.WithSource(
			nvdiags.Error,
			"Reference to undeclared helper",
			fmt.Sprintf("No helper %q %q is declared in the configuration.", addr.Type, addr.Name),
			rng,
		))
		return nil, diags
	}

	return &helperRunNode{
		HelperNode: graphs.HelperNode{
			Addr: addr,
		},
		Config: hc,
	}, diags
}
