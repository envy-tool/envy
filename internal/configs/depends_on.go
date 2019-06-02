package configs

import (
	"github.com/envy-tool/envy/internal/addrs"

	"github.com/hashicorp/hcl2/hcl"
)

// reservedHelperTypeNames are names that cannot be used as helper types
// because they indicate special references.
var reservedHelperTypeNames = map[string]struct{}{
	"call":    struct{}{},
	"command": struct{}{},
	"service": struct{}{},
	"shared":  struct{}{},
	"socket":  struct{}{},
	"pipe":    struct{}{},
	"path":    struct{}{},
}

func decodeDependsOn(expr hcl.Expression) ([]addrs.Reference, hcl.Diagnostics) {
	if expr == nil {
		return nil, nil
	}
	if v, diags := expr.Value(nil); !diags.HasErrors() && v.IsNull() {
		return nil, nil
	}

	var ret []addrs.Reference
	var diags hcl.Diagnostics

	exprs, moreDiags := hcl.ExprList(expr)
	diags = append(diags, moreDiags...)
	if moreDiags.HasErrors() {
		return nil, diags
	}

	for _, expr := range exprs {
		traversal, moreDiags := hcl.AbsTraversalForExpr(expr)
		diags = append(diags, moreDiags...)
		if moreDiags.HasErrors() {
			continue
		}
		ref, remain, moreDiags := DecodeReference(traversal)
		if moreDiags.HasErrors() {
			continue
		}
		if len(remain) > 0 {
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid reference",
				Detail:   "References in depends_on must be to whole objects, not individual attributes.",
				Subject:  remain.SourceRange().Ptr(),
			})
		}
		ret = append(ret, ref)
	}

	return ret, diags
}

// DecodeReference decodes a reference address given as an HCL traversal
// and produces the equivalent reference value, if the reference is valid.
func DecodeReference(traversal hcl.Traversal) (addrs.Reference, hcl.Traversal, hcl.Diagnostics) {
	return addrs.Reference{}, nil, hcl.Diagnostics{
		{
			Severity: hcl.DiagError,
			Summary:  "Reference decoding isn't implemented yet",
			Subject:  traversal.SourceRange().Ptr(),
		},
	}
}
