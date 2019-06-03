package nvdiags

import (
	"github.com/hashicorp/hcl2/hcl"
)

// SourceRange represents a range of characters, possibly spanning across
// multiple lines, in one of the configuration source files.
type SourceRange = hcl.Range

// SourcePos represents a specific source position.
//
// It generally used only as a part of SourceRange, since otherwise it has no
// information about which file it relates to.
type SourcePos = hcl.Pos
