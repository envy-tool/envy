package configs

import (
	"strings"

	"github.com/hashicorp/hcl2/hcl/hclsyntax"
)

func validName(candidate string) bool {
	if !hclsyntax.ValidIdentifier(candidate) {
		return false
	}
	if strings.Contains(candidate, "-") { // We further constrain HCL syntax by disallowing dashes in our own names
		return false
	}
	return true
}
