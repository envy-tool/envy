package addrs

// Service identifies a named service from the envy configuration.
type Service struct {
	Name string
}

func (s Service) isReference() {} // marker for interface Referenceable

func (s Service) String() string {
	return "service." + s.Name
}
