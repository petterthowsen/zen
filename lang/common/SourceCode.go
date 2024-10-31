package common

// SourceCode represents a piece of zen code
type SourceCode interface {
	GetText() string
	GetLength() int
	GetLine(line int) string
	GetChar(index int) rune
	GetLocation(line int, column int) *SourceLocation
}
