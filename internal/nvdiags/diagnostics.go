package nvdiags

// Diagnostics represents a set of diagnostic messages aimed at the end user.
type Diagnostics []Diagnostic

// Diagnostic represents a single diagnostic message.
type Diagnostic interface {
	// TODO: Define this
}

func (diags Diagnostics) HasErrors() bool {
	return false // TODO: Implement
}
