package configs

import (
	"github.com/envy-tool/envy/internal/addrs"

	"github.com/hashicorp/hcl2/hcl"
)

// SharedObject represents a single "shared" block in a configuration.
type SharedObject struct {
	Name      string
	DeclRange hcl.Range

	Attributes hcl.Attributes
}

// Addr returns the address for the shared object that was declared.
func (o *SharedObject) Addr() addrs.SharedObject {
	return addrs.SharedObject{Name: o.Name}
}

func decodeSharedObjectBlock(block *hcl.Block) (*SharedObject, hcl.Diagnostics) {
	so := &SharedObject{
		Name:      block.Labels[0],
		DeclRange: block.DefRange,
	}

	var diags hcl.Diagnostics
	so.Attributes, diags = block.Body.JustAttributes()

	if !validName(so.Name) {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid shared object name",
			Detail:   "All object names must begin with a letter and contain only letters, digits, and underscores.",
			Subject:  block.LabelRanges[0].Ptr(),
		})
	}

	return so, diags
}
