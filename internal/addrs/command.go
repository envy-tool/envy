package addrs

// Command identifies a named command from the envy configuration.
type Command struct {
	Name string
}

func (c Command) isReference() {} // marker for interface Referenceable

func (c Command) String() string {
	return "command." + c.Name
}
