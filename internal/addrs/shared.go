package addrs

// SharedObject represents the address of a particular shared custom object.
//
// Shared objects live inside the user's envy agent but can be accessed from
// any envy command or service.
type SharedObject struct {
	Name string
}

func (o SharedObject) isReference() {} // marker for interface Referenceable

func (o SharedObject) String() string {
	return "shared." + o.Name
}
