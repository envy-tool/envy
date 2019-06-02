package nvdiags

import (
	"github.com/hashicorp/hcl2/hcl"
)

// SourceRange represents a range of characters, possibly spanning across
// multiple lines, in one of the configuration source files.
type SourceRange struct {
	Filename   string
	Start, End SourcePos
}

// SourcePos represents a specific source position.
//
// It generally used only as a part of SourceRange, since otherwise it has no
// information about which file it relates to.
type SourcePos struct {
	Line, Column, Byte int
}

// SourceRangeFromHCL converts a HCL source range into a SourceRange.
func SourceRangeFromHCL(rng hcl.Range) SourceRange {
	return SourceRange{
		Filename: rng.Filename,
		Start:    SourcePosFromHCL(rng.Start),
		End:      SourcePosFromHCL(rng.End),
	}
}

// SourcePosFromHCL converts a HCL source position into a SourcePos.
func SourcePosFromHCL(rng hcl.Pos) SourcePos {
	return SourcePos{
		Line:   rng.Line,
		Column: rng.Column,
		Byte:   rng.Byte,
	}
}

// ToHCL converts the receiver into an HCL SourceRange.
func (r SourceRange) ToHCL() hcl.Range {
	return hcl.Range{
		Filename: r.Filename,
		Start:    r.Start.ToHCL(),
		End:      r.End.ToHCL(),
	}
}

// ToHCL converts the receiver into an HCL SourcePos.
func (p SourcePos) ToHCL() hcl.Pos {
	return hcl.Pos{
		Line:   p.Line,
		Column: p.Column,
		Byte:   p.Byte,
	}
}
