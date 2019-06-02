package addrs

// Referenceable is an interface implemented by all of the address types that
// can be referenced.
type Referenceable interface {
	isReference()
	String() string
}
