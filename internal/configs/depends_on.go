package configs

import (
	"github.com/hashicorp/hcl2/hcl"
)

func decodeDependsOn(expr hcl.Expression) ([]Reference, hcl.Diagnostics) {
	if expr == nil {
		return nil, nil
	}
	if v, diags := expr.Value(nil); !diags.HasErrors() && v.IsNull() {
		return nil, nil
	}

	var ret []Reference
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
