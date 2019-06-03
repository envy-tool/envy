package nvdiags

import (
	"fmt"

	"github.com/hashicorp/hcl2/hcl"
)

// Diagnostics represents a set of diagnostic messages aimed at the end user.
type Diagnostics []Diagnostic

// Diagnostic represents a single diagnostic message.
type Diagnostic interface {
	Severity() Severity
	Messages() Messages
	Locations() Locations
	ExprContext() *ExprContext
}

// Severity represents the severity level of a diagnostic
type Severity = hcl.DiagnosticSeverity

const (
	// Error is a diagnostic severity indicating a problem that prevented
	// the successful completion of an operation.
	Error Severity = hcl.DiagError

	// Warning is a diagnostic severity indicating a problem that the user
	// should be aware of but that did not prevent the successful completion
	// of an operation.
	Warning Severity = hcl.DiagWarning
)

// Messages reprseents the natural-language messages associated with a diagnostic.
type Messages struct {
	Summary string
	Detail  string
}

// Locations represents the source locations associated with a diagnostic, if any.
type Locations struct {
	Subject *SourceRange
	Context *SourceRange
}

// ExprContext represents the expression context associated with a diagnostic, if any.
type ExprContext struct {
	Expression  hcl.Expression
	EvalContext *hcl.EvalContext
}

func (diags Diagnostics) HasErrors() bool {
	for _, diag := range diags {
		if diag.Severity() == Error {
			return true
		}
	}
	return false
}

// Append appends the given values to the receiving diagnostics and returns
// the result.
//
// Append arguments can be of the following types: Diagnostics, Diagnostic,
// hcl.Diagnostics, or *hcl.Diagnostic. If any other type is passed, this
// function will panic.
func (diags Diagnostics) Append(more ...interface{}) Diagnostics {
	for _, v := range more {
		switch tv := v.(type) {
		case Diagnostics:
			diags = append(diags, tv...)
		case Diagnostic:
			diags = append(diags, tv)
		case hcl.Diagnostics:
			for _, diag := range tv {
				diags = append(diags, FromHCL(diag))
			}
		case *hcl.Diagnostic:
			diags = append(diags, FromHCL(tv))
		default:
			panic(fmt.Sprintf("cannot append %T to diagnostics", v))
		}
	}
	return diags
}
