package addrs

// Shared represents the address of a particular shared custom object.
//
// Shared objects live inside the user's envy agent but can be accessed from
// any envy command or service.
type Shared struct {
	Name string
}

func (s Shared) isReference() {} // marker for interface Referenceable

func (s Shared) String() string {
	return "shared." + s.Name
}
