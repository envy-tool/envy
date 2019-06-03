package runs

import (
	"envy.pw/cli/internal/configs"
	"envy.pw/cli/internal/graphs"
)

type commandExecNode struct {
	graphs.CommandNode
	Config *configs.Command
}
