package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

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

func supportedOS() bool {
	// We can only support operating systems that userdirs can run on
	return userdirs.SupportedOS()
}
