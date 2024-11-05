package expression

import (
	"fmt"
	"strings"
	"zen/lang/common"
	"zen/lang/parsing/ast"
)

// MapEntry represents a key-value pair in a map literal
type MapEntry struct {
	Key   ast.Expression
	Value ast.Expression
}

// MapLiteralExpression represents a map literal in the AST
type MapLiteralExpression struct {
	Entries  []MapEntry
	Location *common.SourceLocation
}

func NewMapLiteralExpression(entries []MapEntry, location *common.SourceLocation) *MapLiteralExpression {
	return &MapLiteralExpression{
		Entries:  entries,
		Location: location,
	}
}

func (e *MapLiteralExpression) Accept(visitor ast.Visitor) interface{} {
	return visitor.VisitMapLiteral(e)
}

func (e *MapLiteralExpression) GetLocation() *common.SourceLocation {
	return e.Location
}

func (e *MapLiteralExpression) IsExpression() {}

func (e *MapLiteralExpression) String(indent int) string {
	var entries []string
	for _, entry := range e.Entries {
		keyStr := entry.Key.String(indent + 1)
		valueStr := entry.Value.String(indent + 2)
		entries = append(entries, fmt.Sprintf("%sKey:\n%s%sValue:\n%s",
			strings.Repeat("  ", indent+1), keyStr,
			strings.Repeat("  ", indent+1), valueStr))
	}
	return fmt.Sprintf("%sMapLiteral:\n%s", strings.Repeat("  ", indent), strings.Join(entries, "\n"))
}
