package addrs

// Command identifies a named command from the envy configuration.
type Command struct {
	Name string
}

// MakeCommand returns a Command address for the given name.
//
// It will panic if the given name is not a valid identifier.
func MakeCommand(name string) Command {
	assertValidName(name)
	return Command{Name: name}
}

func (c Command) isReference() {} // marker for interface Referenceable

func (c Command) String() string {
	return "command." + c.Name
}
