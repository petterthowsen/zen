package statement

import (
	"strings"
	"zen/lang/common"
	"zen/lang/parsing/ast"
)

// ForInStatement represents a for-in loop in the AST
// Syntax:
//
//	for key, value in map { body }
//	for value in map { body }
type ForInStatement struct {
	Location  *common.SourceLocation
	Key       string         // Optional key variable (can be empty)
	Value     string         // Value variable
	Container ast.Expression // Expression being iterated over (e.g., map, array)
	Body      []ast.Statement
}

func NewForInStatement(
	key string,
	value string,
	container ast.Expression,
	body []ast.Statement,
	location *common.SourceLocation,
) *ForInStatement {
	return &ForInStatement{
		Key:       key,
		Value:     value,
		Container: container,
		Body:      body,
		Location:  location,
	}
}

func (s *ForInStatement) Accept(visitor ast.Visitor) interface{} {
	return visitor.VisitForInStatement(s)
}

func (s *ForInStatement) GetLocation() *common.SourceLocation {
	return s.Location
}

func (s *ForInStatement) IsStatement() {}

func (s *ForInStatement) String(indent int) string {
	var sb strings.Builder
	indentStr := strings.Repeat("  ", indent)

	sb.WriteString(indentStr + "ForInStatement\n")

	if s.Key != "" {
		sb.WriteString(indentStr + "  Key: " + s.Key + "\n")
	}
	sb.WriteString(indentStr + "  Value: " + s.Value + "\n")

	sb.WriteString(indentStr + "  Container:\n")
	sb.WriteString(s.Container.String(indent + 2))

	sb.WriteString(indentStr + "  Body:\n")
	for _, stmt := range s.Body {
		sb.WriteString(stmt.String(indent + 2))
	}

	return sb.String()
}
