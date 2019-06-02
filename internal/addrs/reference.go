package addrs

import (
	"github.com/envy-tool/envy/internal/nvdiags"
)

// Reference represents a reference to a referenceable address somewhere in
// the configuration, retaining information about where that reference was
// made.
type Reference struct {
	Addr        Referenceable
	SourceRange nvdiags.SourceRange
}

// Referenceable is an interface implemented by all of the address types that
// can be referenced.
type Referenceable interface {
	isReference()
	String() string
}
