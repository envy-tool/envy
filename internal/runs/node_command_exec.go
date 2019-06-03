package runs

import (
	"envy.pw/cli/internal/configs"
	"envy.pw/cli/internal/graphs"
)

type commandExecNode struct {
	graphs.CommandNode
	Config *configs.Command
}

func (n *commandExecNode) References() []configs.Reference {
	return n.Config.AllReferences()
}
