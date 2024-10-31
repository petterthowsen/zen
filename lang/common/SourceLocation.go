package common

import (
	"fmt"
	"strings"
)

// SourceLocation represents a specific location within a source code.
type SourceLocation struct {
	Source SourceCode
	Line   int
	Column int
}

// String returns a formatted string representing the source location in the format "Line %d, Column %d".
func (sl *SourceLocation) String() string {
	return fmt.Sprintf("Line %d, Column %d", sl.Line, sl.Column)
}

// GetLine returns the source code text at the specified line within Source.
func (sl *SourceLocation) GetLine() string {
	return sl.Source.GetLine(sl.Line)
}

// GetLineWithMarker returns the code at the line with blank line with a caret (^) pointing to the column
func (sl *SourceLocation) GetLineWithMarker() string {
	return fmt.Sprintf("%s\n%s^", sl.GetLine(), strings.Repeat(" ", max(0, sl.Column)))
}
