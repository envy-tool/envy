package nvdiags

import (
	"github.com/hashicorp/hcl2/hcl"
)

func WithSource(
	severity Severity,
	summary string,
	detail string,
	srcRange SourceRange,
) Diagnostic {
	return withSource{
		Diagnostic: &hcl.Diagnostic{
			Severity: severity,
			Summary:  summary,
			Detail:   detail,
			Subject:  &srcRange,
		},
	}
}

func FromHCL(diag *hcl.Diagnostic) Diagnostic {
	return withSource{
		Diagnostic: diag,
	}
}

type withSource struct {
	*hcl.Diagnostic
}

func (diag withSource) Severity() Severity {
	return diag.Diagnostic.Severity
}

func (diag withSource) Messages() Messages {
	return Messages{
		Summary: diag.Diagnostic.Summary,
		Detail:  diag.Diagnostic.Detail,
	}
}

func (diag withSource) Locations() Locations {
	return Locations{
		Subject: diag.Diagnostic.Subject,
		Context: diag.Diagnostic.Context,
	}
}

func (diag withSource) ExprContext() *ExprContext {
	if diag.Diagnostic.Expression != nil {
		return &ExprContext{
			Expression:  diag.Diagnostic.Expression,
			EvalContext: diag.Diagnostic.EvalContext,
		}
	}
	return nil
}
