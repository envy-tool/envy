package addrs

import (
	"fmt"
)

// Helper identifies a particular helper defined in the configuration.
type Helper struct {
	Type string
	Name string
}

func (h Helper) isReference() {} // marker for interface Referenceable

func (h Helper) String() string {
	return fmt.Sprintf("%s.%s", h.Type, h.Name)
}
