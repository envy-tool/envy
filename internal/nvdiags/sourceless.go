package nvdiags

func Sourceless(
	severity Severity,
	summary string,
	detail string,
) Diagnostic {
	return &sourceless{
		severity: severity,
		summary:  summary,
		detail:   detail,
	}
}

type sourceless struct {
	severity Severity
	summary  string
	detail   string
}

func (diag sourceless) Severity() Severity {
	return diag.severity
}

func (diag sourceless) Messages() Messages {
	return Messages{
		Summary: diag.summary,
		Detail:  diag.detail,
	}
}

func (diag sourceless) Locations() Locations {
	return Locations{}
}

func (diag sourceless) ExprContext() *ExprContext {
	return nil
}
