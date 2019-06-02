package addrs

// SharedObject represents the address of a particular shared custom object.
//
// Shared objects live inside the user's envy agent but can be accessed from
// any envy command or service.
type SharedObject struct {
	Name string
}

// MakeSharedObject returns a SharedObject address for the given name.
//
// It will panic if the given name is not a valid identifier.
func MakeSharedObject(name string) SharedObject {
	assertValidName(name)
	return SharedObject{Name: name}
}

func (o SharedObject) isReference() {} // marker for interface Referenceable

func (o SharedObject) String() string {
	return "shared." + o.Name
}
