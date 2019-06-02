package addrs

import (
	"fmt"
)

// Helper identifies a particular helper defined in the configuration.
type Helper struct {
	Type string
	Name string
}

// MakeHelper returns a Helper address for the given type and name.
//
// It will panic if the type and/or the name is not a valid identifier.
func MakeHelper(typeName, name string) Helper {
	assertValidName(typeName)
	assertValidName(name)
	return Helper{Type: typeName, Name: name}
}

func (h Helper) isReference() {} // marker for interface Referenceable

func (h Helper) String() string {
	return fmt.Sprintf("%s.%s", h.Type, h.Name)
}
