package nvdiags

// SourceRange represents a range of characters, possibly spanning across
// multiple lines, in one of the configuration source files.
type SourceRange struct {
	Filename string
	Start, End SourcePos
}

// SourcePos represents a specific source position.
//
// It generally used only as a part of SourceRange, since otherwise it has no
// information about which file it relates to.
type SourcePos struct {
	Line, Column, Byte int
}
