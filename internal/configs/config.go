package configs

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"envy.pw/cli/internal/addrs"

	"github.com/hashicorp/hcl2/hcl"
)

// Config represents an entire configuration, with all of the core portions
// guaranteed valid if it was returned without any error diagnostics.
//
// Conventionally, a config is built by merging the contents of all of the
// files in a particular directory. That behavior is implemented by the
// function LoadConfig, which is the primary way to create a Config object.
//
// Callers should consider a Config object to be immutable after it is created.
// If any changes are made, the guarantee of validity no longer applies.
type Config struct {
	// BaseDir is the directory that any configuration-relative paths must
	// be resolved relative to.
	BaseDir string

	Commands      map[addrs.Command]*Command
	Helpers       map[addrs.Helper]*Helper
	SharedObjects map[addrs.SharedObject]*SharedObject
}

func newConfig(baseDir string) *Config {
	return &Config{
		BaseDir: baseDir,

		Commands:      map[addrs.Command]*Command{},
		Helpers:       map[addrs.Helper]*Helper{},
		SharedObjects: map[addrs.SharedObject]*SharedObject{},
	}
}

// LoadConfig reads all of the configuration files (.nv.hcl and .nv.json
// extensions) in the given directory, merges them into a single Config,
// and returns it.
//
// If the returned diagnostics contains errors then the returned Config may
// be incomplete, but will include whatever subset of the configuration could
// be understood in spite of those errors.
func LoadConfig(dir string) (*Config, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	cfg := newConfig(dir)

	items, err := ioutil.ReadDir(dir)
	if err != nil {
		switch {
		case os.IsNotExist(err):
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Configuration directory not found",
				Detail:   fmt.Sprintf("The configuration directory %q does not exist.", dir),
			})
		default:
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Cannot open configuration directory",
				Detail:   fmt.Sprintf("Cannot open %q: %s.", dir, err),
			})
		}
		return cfg, diags
	}

	for _, info := range items {
		if info.IsDir() {
			continue
		}
		name := info.Name()
		if !IsConfigFile(info.Name()) {
			continue
		}

		file, moreDiags := LoadConfigFile(name)
		diags = append(diags, moreDiags...)

		moreDiags = cfg.mergeFile(file)
		diags = append(diags, moreDiags...)
	}

	return cfg, diags
}

// BuildConfig is similar to LoadConfig but it works with some files already
// loaded some other way (e.g. via LoadConfigFile) rather than reading the
// files from disk itself.
func BuildConfig(baseDir string, files []*File) (*Config, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	cfg := newConfig(baseDir)

	for _, file := range files {
		moreDiags := cfg.mergeFile(file)
		diags = append(diags, moreDiags...)
	}

	return cfg, diags
}

func (c *Config) mergeFile(f *File) hcl.Diagnostics {
	var diags hcl.Diagnostics

	for _, cmd := range f.Commands {
		addr := cmd.Addr()
		if existing, exists := c.Commands[addr]; exists {
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Command name conflict",
				Detail:   fmt.Sprintf("A command named %q was already declared at %s.", cmd.Name, existing.DeclRange),
				Subject:  cmd.DeclRange.Ptr(),
			})
			continue
		}
		c.Commands[addr] = cmd
	}

	for _, h := range f.Helpers {
		addr := h.Addr()
		if existing, exists := c.Helpers[addr]; exists {
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Helper name conflict",
				Detail:   fmt.Sprintf("A helper %q %q was already declared at %s.", h.Type, h.Name, existing.DeclRange),
				Subject:  h.DeclRange.Ptr(),
			})
			continue
		}
		c.Helpers[addr] = h
	}

	for _, so := range f.SharedObjects {
		addr := so.Addr()
		if existing, exists := c.SharedObjects[addr]; exists {
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Shared object name conflict",
				Detail:   fmt.Sprintf("A shared object named %q was already declared at %s.", so.Name, existing.DeclRange),
				Subject:  so.DeclRange.Ptr(),
			})
			continue
		}
		c.SharedObjects[addr] = so
	}

	return nil
}

// IsConfigFile returns true if the given filename should be recognized as
// a configuration file.
//
// It does not access the filesystem itself, so it's the caller's
// responsibility to ensure that the given name exists and refers to a file,
// and in particular not to a directory.
func IsConfigFile(name string) bool {
	if strings.HasPrefix(name, ".") || strings.HasSuffix(name, "~") || (strings.HasPrefix(name, "#") && strings.HasSuffix(name, "#")) {
		// Various "hidden" files: unix-style hidden, and emacs/vim cache/backup files
		return false
	}

	return strings.HasSuffix(name, ".nv.hcl") || strings.HasSuffix(name, ".nv.json")
}
