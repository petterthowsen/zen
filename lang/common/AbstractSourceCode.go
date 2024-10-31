package common

import (
	"strings"
)

// AbstractSourceCode represents source code, which may come from anywhere (stdin, file etc.)
type AbstractSourceCode struct {
	text      string
	lineCache []string
}

// Initialize the line cache (called when text is set or updated)
func (src *AbstractSourceCode) initLineCache() {
	src.lineCache = strings.Split(src.text, "\n")
}

// GetText returns the entire source code text
func (src *AbstractSourceCode) GetText() string {
	return src.text
}

// GetLength returns the length of the source code text
func (src *AbstractSourceCode) GetLength() int {
	return len(src.text)
}

// GetChar returns the character at a given index
func (src *AbstractSourceCode) GetChar(index int) rune {
	if index < 0 || index >= len(src.text) {
		return rune(0)
	}
	return rune(src.text[index])
}

// GetLine returns the source code at a given line
func (src *AbstractSourceCode) GetLine(line int) string {
	// Initialize the cache if it's nil
	if src.lineCache == nil {
		src.initLineCache()
	}
	if line < 1 || line > len(src.lineCache) {
		return ""
	}
	return src.lineCache[line-1]
}

// GetLocation returns a SourceLocation at a given line and column
func (src *AbstractSourceCode) GetLocation(line int, column int) *SourceLocation {
	return &SourceLocation{
		Source: src,
		Line:   line,
		Column: column,
	}
}
