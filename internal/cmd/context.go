package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"envy.pw/cli/internal/configs"
	"envy.pw/cli/internal/nvdiags"
	"envy.pw/cli/internal/runs"

	"github.com/apparentlymart/go-userdirs/userdirs"
)

// RunContext contains some general contextual information about a run of
// the "envy" CLI.
type RunContext struct {
	ConfigDir  string
	WorkingDir string
}

func newRunContext(configDir, workingDir string) (*RunContext, error) {
	if !supportedOS() {
		return nil, fmt.Errorf("envy cannot run on %s", runtime.GOOS)
	}

	dirs := userdirs.ForApp("Envy", "Envy", "pw.envy")
	if configDir == "" {
		configDir = filepath.Join(dirs.ConfigHome(), "config")
		// If we're using the default directory then we'll try to make it now.
		// If this fails then we'll just ignore that and catch it as a
		// "not found" error downstream.
		os.MkdirAll(configDir, 0700)
	}

	if workingDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("cannot determine working directory: %s", err)
		}
		workingDir = wd
	}

	return &RunContext{
		ConfigDir:  configDir,
		WorkingDir: workingDir,
	}, nil
}

// LoadConfig loads a configuration from the context's configuration directory.
func (c *RunContext) LoadConfig() (*configs.Config, nvdiags.Diagnostics) {
	var diags nvdiags.Diagnostics
	cfg, hclDiags := configs.LoadConfig(c.ConfigDir)
	diags = diags.Append(hclDiags)
	return cfg, diags
}

// NewRunner creates a runner using the settings from the context.
func (c *RunContext) NewRunner() (*runs.Runner, nvdiags.Diagnostics) {
	// TODO: Eventually this should be doing a bunch more work to compute
	// various other contextual information, such as a set of available
	// provider plugins.
	return runs.NewRunner(), nil
}

func supportedOS() bool {
	// We can only support operating systems that userdirs can run on
	return userdirs.SupportedOS()
}
