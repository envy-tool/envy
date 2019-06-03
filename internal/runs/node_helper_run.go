package runs

import (
	"envy.pw/cli/internal/configs"
	"envy.pw/cli/internal/graphs"
)

type helperRunNode struct {
	graphs.HelperNode
	Config *configs.Helper
}
