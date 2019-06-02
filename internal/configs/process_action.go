package configs

import (
	"github.com/hashicorp/hcl2/hcl"
)

// ProcessAction represents an action to be taken against a child process in
// response to some event.
type ProcessAction int

const (
	// ProcessIgnore indicates that no action should be taken on the process.
	ProcessIgnore ProcessAction = iota

	// ProcessRestart indicates that the process should be restarted.
	ProcessRestart

	// ProcessTerminate indicates that the process should be terminated.
	ProcessTerminate
)

func decodeProcessAction(expr hcl.Expression) (ProcessAction, hcl.Diagnostics) {
	if expr == nil {
		return ProcessIgnore, nil
	}
	if v, diags := expr.Value(nil); !diags.HasErrors() && v.IsNull() {
		return ProcessIgnore, nil
	}

	kw := hcl.ExprAsKeyword(expr)
	switch kw {
	case "ignore":
		return ProcessIgnore, nil
	case "restart":
		return ProcessRestart, nil
	case "terminate":
		return ProcessTerminate, nil
	default:
		return ProcessIgnore, hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "Invalid process action",
				Detail:   "Must be one of the following keywords: ignore, restart, or terminate.",
				Subject:  expr.StartRange().Ptr(),
			},
		}
	}
}
