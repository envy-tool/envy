package configs

import (
	"fmt"

	"envy.pw/cli/internal/addrs"

	"github.com/hashicorp/hcl2/hcl"
)

// Helper represents a single "helper" block in a configuration.
type Helper struct {
	Type      string
	Name      string
	DeclRange hcl.Range

	Body hcl.Body
}

// Addr returns the address for the helper that was declared.
func (h *Helper) Addr() addrs.Helper {
	return addrs.Helper{
		Type: h.Type,
		Name: h.Name,
	}
}

func decodeHelperBlock(block *hcl.Block) (*Helper, hcl.Diagnostics) {
	h := &Helper{
		Type:      block.Labels[0],
		Name:      block.Labels[1],
		DeclRange: block.DefRange,
	}

	var diags hcl.Diagnostics
	h.Body = block.Body

	if !validName(h.Type) {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid helper type",
			Detail:   "Helper type names begin with a letter and contain only letters, digits, and underscores.",
			Subject:  block.LabelRanges[0].Ptr(),
		})
	}
	if IsReservedHelperType(h.Type) {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid helper type",
			Detail:   fmt.Sprintf("The helper type name %q is reserved for accessing %s objects.", h.Type, h.Type),
			Subject:  block.LabelRanges[0].Ptr(),
		})
	}
	if !validName(h.Name) {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid helper name",
			Detail:   "All object names must begin with a letter and contain only letters, digits, and underscores.",
			Subject:  block.LabelRanges[1].Ptr(),
		})
	}

	return h, diags
}

// IsReservedHelperType returns true if the given helper type would create
// a helper that cannot be referenced, due to its type name instead indicating
// that a different kind of object is being referenced.
func IsReservedHelperType(proposed string) bool {
	_, exists := reservedHelperTypeNames[proposed]
	return exists
}
