package addrs

import (
	"fmt"
)

// Path identifies a path type.
//
// The set of valid path types is fixed, because these path values are built
// in and can only be referenced in config, not declared.
type Path string

const (
	// PathWorking is a Path representing the current working directory.
	PathWorking Path = "cwd"

	// PathConfig is a Path representing the main configuration directory.
	PathConfig Path = "config"

	// PathTemp is a Path representing an object-specific temporary directory
	// that will be deleted when the owning object is destroyed.
	PathTemp Path = "temp"
)

// All of the path addresses declared above must also be included here so that
// ValidPathType can return a correct result.
var validPathTypes = map[Path]struct{}{
	PathWorking: struct{}{},
	PathConfig:  struct{}{},
	PathTemp:    struct{}{},
}

// MakePath returns a Path address for the given type.
//
// It will panic if the given name is not a valid path type, as defined by
// ValidPathType.
func MakePath(typeName string) Path {
	if !ValidPathType(typeName) {
		panic(fmt.Sprintf("invalid path address type %q", typeName))
	}
	return Path(typeName)
}

func (p Path) isReference() {} // marker for interface Referenceable

func (p Path) String() string {
	return "path." + string(p)
}

// ValidPathType returns true if and only if the given string is a valid type
// for a Path address.
func ValidPathType(candidate string) bool {
	_, ok := validPathTypes[Path(candidate)]
	return ok
}
