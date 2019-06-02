// This is the main package for the "Envy" CLI tool.
//
// This is a Go program, not a library. There is no public API for other
// Go programs to call.
package main

import (
	"github.com/envy-tool/envy/internal/cmd"
)

func main() {
	cmd.Execute()
}
