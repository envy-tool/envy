package addrs

// Service identifies a named service from the envy configuration.
type Service struct {
	Name string
}

// MakeService returns a Service address for the given name.
//
// It will panic if the given name is not a valid identifier.
func MakeService(name string) Service {
	assertValidName(name)
	return Service{Name: name}
}

func (s Service) isReference() {} // marker for interface Referenceable

func (s Service) String() string {
	return "service." + s.Name
}
