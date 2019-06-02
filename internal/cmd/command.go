// Package cmd handles the command line interface for Envy.
package cmd

import (
	"fmt"
	"os"

	"github.com/envy-tool/envy/internal/nvdiags"

	"github.com/spf13/cobra"
)

// Execute is the main entry point for this package.
//
// It does not return.
func Execute() {
	var ctx *RunContext
	var configDir string
	var workingDir string

	type Command interface {
		Run() nvdiags.Diagnostics
	}
	var command Command

	var rootCmd = &cobra.Command{
		Use:   "envy",
		Short: "Envy is a command launcher supporting dynamic credentials and data",
		Long:  `A command launcher that can generate and provide dynamic credentials and other data to programs that need them.`,
		Args:  cobra.NoArgs,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var err error
			ctx, err = newRunContext(configDir, workingDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
				os.Exit(1)
			}

			// We'll switch to the config directory as our working directory
			// so that paths in the config are config-relative.
			// Workingdir-relative paths can still be constructed from
			// ctx.WorkingDir if needed.
			err = os.Chdir(ctx.ConfigDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
				os.Exit(1)
			}
		},
		TraverseChildren: true,
	}
	rootCmd.Flags().StringVarP(&configDir, "config-dir", "c", "", "directory to search for configuration files")
	rootCmd.Flags().StringVarP(&workingDir, "working-dir", "w", "", "directory to use as the working directory when running commands")

	var runCmd = &cobra.Command{
		Use:   "run <command-name> [args...]",
		Short: "Run a configured command",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			command = &runCommand{
				Context: ctx,
				Args:    args,
			}
		},
	}
	runCmd.Flags().SetInterspersed(false) // Everything after the command name appaers in "args", including flag-like strings
	rootCmd.AddCommand(runCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	diags := command.Run()
	// TODO: Print out the diagnostics, if any
	if diags.HasErrors() {
		os.Exit(1)
	}
	os.Exit(0)
}
