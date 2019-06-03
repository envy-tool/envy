package configs

import (
	"fmt"

	"envy.pw/cli/internal/addrs"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
)

// Reference represents a reference to a referenceable address somewhere in
// the configuration, retaining information about where that reference was
// made.
type Reference struct {
	Addr        addrs.Referenceable
	SourceRange hcl.Range
}

// reservedHelperTypeNames are names that cannot be used as helper types
// because they indicate special references.
var reservedHelperTypeNames = map[string]struct{}{
	"call":    struct{}{},
	"command": struct{}{},
	"service": struct{}{},
	"shared":  struct{}{},
	"socket":  struct{}{},
	"pipe":    struct{}{},
	"path":    struct{}{},
}

// DecodeReference decodes a reference address given as an HCL absolute
// traversal and produces the equivalent reference value, if the reference is
// valid.
//
// If any additional traversal steps appear after the reference, they are
// returned as a relative traversal in the second return value.
//
// If the returned diagnostics contains errors, the reference and remaining
// traversal are not valid.
func DecodeReference(traversal hcl.Traversal) (Reference, hcl.Traversal, hcl.Diagnostics) {
	if traversal.IsRelative() {
		// Programming error: contract requires only absolute traversals
		panic("DecodeReference with relative traversal")
	}

	switch rootName := traversal.RootName(); rootName {

	case "command":
		const errSummary = "Invalid command reference"
		if len(traversal) < 2 {
			return Reference{}, nil, hcl.Diagnostics{
				{
					Severity: hcl.DiagError,
					Summary:  errSummary,
					Detail:   "The keyword \"command\" must be followed by a command name using attribute access syntax.",
					Subject:  traversal.SourceRange().Ptr(),
				},
			}
		}
		nameStep, ok := traversal[1].(hcl.TraverseAttr)
		if !ok {
			return Reference{}, nil, hcl.Diagnostics{
				{
					Severity: hcl.DiagError,
					Summary:  errSummary,
					Detail:   "The keyword \"command\" must be followed by a helper name using attribute access syntax.",
					Subject:  traversal.SourceRange().Ptr(),
				},
			}
		}
		if !validName(nameStep.Name) {
			return Reference{}, nil, hcl.Diagnostics{
				{
					Severity: hcl.DiagError,
					Summary:  errSummary,
					Detail:   fmt.Sprintf("%q is not a valid command name.", nameStep.Name),
					Subject:  nameStep.SourceRange().Ptr(),
				},
			}
		}
		return Reference{
			Addr:        addrs.MakeCommand(nameStep.Name),
			SourceRange: traversal.SourceRange(),
		}, traversal[2:], nil

	default:
		if IsReservedHelperType(rootName) {
			// Should not get here; indicates we didn't handle one of the
			// reserved names above.
			return Reference{}, nil, hcl.Diagnostics{
				{
					Severity: hcl.DiagError,
					Summary:  "Unhandled reference type",
					Detail:   fmt.Sprintf("The reference parser did not handle address type %q. This is a bug in Envy.", rootName),
					Subject:  traversal.SourceRange().Ptr(),
				},
			}
		}

		const errSummary = "Invalid helper reference"
		if !validName(rootName) {
			return Reference{}, nil, hcl.Diagnostics{
				{
					Severity: hcl.DiagError,
					Summary:  errSummary,
					Detail:   fmt.Sprintf("%q is not a valid helper type name.", rootName),
					Subject:  traversal[0].SourceRange().Ptr(),
				},
			}
		}
		if len(traversal) < 2 {
			return Reference{}, nil, hcl.Diagnostics{
				{
					Severity: hcl.DiagError,
					Summary:  errSummary,
					Detail:   fmt.Sprintf("The helper type name %q must be followed by a helper name using attribute access syntax.", rootName),
					Subject:  traversal.SourceRange().Ptr(),
				},
			}
		}
		nameStep, ok := traversal[1].(hcl.TraverseAttr)
		if !ok {
			return Reference{}, nil, hcl.Diagnostics{
				{
					Severity: hcl.DiagError,
					Summary:  errSummary,
					Detail:   fmt.Sprintf("The helper type name %q must be followed by a helper name using attribute access syntax.", rootName),
					Subject:  traversal.SourceRange().Ptr(),
				},
			}
		}
		if !validName(nameStep.Name) {
			return Reference{}, nil, hcl.Diagnostics{
				{
					Severity: hcl.DiagError,
					Summary:  errSummary,
					Detail:   fmt.Sprintf("%q is not a valid helper name.", nameStep.Name),
					Subject:  nameStep.SourceRange().Ptr(),
				},
			}
		}
		return Reference{
			Addr:        addrs.MakeHelper(rootName, nameStep.Name),
			SourceRange: traversal.SourceRange(),
		}, traversal[2:], nil
	}
}

// ParseReferenceStr is like DecodeReference but it works with a string
// representation of a reference address, rather than a traversal object.
//
// It first parses the string using the HCL native reference syntax, and then
// passes it to DecodeReference. Because this function includes a parsing
// step, the returned diagnostics may include parse errors.
//
// If the returned diagnostics contains errors, the reference and remaining
// traversal are not valid.
func ParseReferenceStr(str string) (Reference, hcl.Traversal, hcl.Diagnostics) {
	traversal, diags := hclsyntax.ParseTraversalAbs([]byte(str), "", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return Reference{}, nil, diags
	}

	return DecodeReference(traversal)
}

func exprReferences(expr hcl.Expression) []Reference {
	traversals := expr.Variables()
	refs := make([]Reference, 0, len(traversals))
	for _, traversal := range traversals {
		ref, _, diags := DecodeReference(traversal)
		if diags.HasErrors() {
			continue
		}
		refs = append(refs, ref)
	}
	return refs
}
