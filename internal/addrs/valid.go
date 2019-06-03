package addrs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl2/hcl/hclsyntax"
)

// ValidName returns true if and only if the given string is valid as a name.
func ValidName(candidate string) bool {
	if !hclsyntax.ValidIdentifier(candidate) {
		return false
	}
	if strings.Contains(candidate, "-") {
		return false
	}
	return true
}

// assertValidName is like validName except that it panics if the given name
// isn't valid.
func assertValidName(candidate string) {
	if !ValidName(candidate) {
		panic(fmt.Sprintf("invalid name %q", candidate))
	}
}
