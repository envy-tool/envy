package addrs

import (
	"fmt"

	"github.com/hashicorp/hcl2/hcl/hclsyntax"
)

func validName(candidate string) bool {
	return hclsyntax.ValidIdentifier(candidate)
}

// assertValidName is like validName except that it panics if the given name
// isn't valid.
func assertValidName(candidate string) {
	if !validName(candidate) {
		panic(fmt.Sprintf("invalid name %q", candidate))
	}
}
