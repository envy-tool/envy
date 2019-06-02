package configs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclparse"
)

// File represents the content of a single configuration file. Multiple files
// combine together to produce a whole configuration, represented by type
// Config.
type File struct {
	Commands      []*Command
	SharedObjects []*SharedObject
}

func newFile() *File {
	return &File{}
}

// LoadConfigFile loads a single configuration file.
//
// The suffix of the given filename must be either ".nv.hcl" or ".nv.json",
// and is used to determine whether to expect HCL native syntax or HCL JSON.
//
// If the returned diagnostics contains errors, the file may be incomplete
// but will include the subset of the file that was read successfully.
//
// It's rare to need to call this function directly. Instead, prefer to call
// LoadConfig to load a full configuration from a directory.
func LoadConfigFile(name string) (*File, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	file := newFile()

	parser := hclparse.NewParser()

	var rawFile *hcl.File
	switch {
	case strings.HasSuffix(name, ".nv.hcl"):
		var moreDiags hcl.Diagnostics
		rawFile, moreDiags = parser.ParseHCLFile(name)
		diags = append(diags, moreDiags...)
	case strings.HasSuffix(name, ".nv.json"):
		var moreDiags hcl.Diagnostics
		rawFile, moreDiags = parser.ParseJSONFile(name)
		diags = append(diags, moreDiags...)
	default:
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid configuration file name",
			Detail:   fmt.Sprintf("The file %q cannot be used as a configuration file: name must have either a .nv.hcl or .nv.json suffix.", name),
		})
		return file, diags
	}

	if rawFile == nil {
		return file, diags
	}

	content, moreDiags := rawFile.Body.Content(configFileSchema)
	diags = append(diags, moreDiags...)

	for _, block := range content.Blocks {
		switch block.Type {

		case "command":
			cmd, moreDiags := decodeCommandBlock(block)
			file.Commands = append(file.Commands, cmd)
			diags = append(diags, moreDiags...)

		case "shared":
			so, moreDiags := decodeSharedObjectBlock(block)
			file.SharedObjects = append(file.SharedObjects, so)
			diags = append(diags, moreDiags...)

		default:
			// Should never get here because Body.Content should ensure
			// everything fits our schema and the above cases should cover
			// all of the defined types.
			panic(fmt.Sprintf("unexpected root configuration block type %q at %s", block.Type, block.TypeRange))
		}
	}

	return file, diags
}

var configFileSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{Type: "command", LabelNames: []string{"name"}},
		{Type: "helper", LabelNames: []string{"type", "name"}},
		{Type: "service", LabelNames: []string{"name"}},
		{Type: "shared", LabelNames: []string{"name"}},
	},
}
