package statement

import (
	"strings"
	"zen/lang/common"
	"zen/lang/parsing/ast"
)

// WhileStatement represents a while loop in the AST
// Syntax: while condition { body }
type WhileStatement struct {
	Location  *common.SourceLocation
	Condition ast.Expression
	Body      []ast.Statement
}

func NewWhileStatement(
	condition ast.Expression,
	body []ast.Statement,
	location *common.SourceLocation,
) *WhileStatement {
	return &WhileStatement{
		Condition: condition,
		Body:      body,
		Location:  location,
	}
}

func (s *WhileStatement) Accept(visitor ast.Visitor) interface{} {
	return visitor.VisitWhileStatement(s)
}

func (s *WhileStatement) GetLocation() *common.SourceLocation {
	return s.Location
}

func (s *WhileStatement) IsStatement() {}

func (s *WhileStatement) String(indent int) string {
	var sb strings.Builder
	indentStr := strings.Repeat("  ", indent)

	sb.WriteString(indentStr + "WhileStatement\n")

	sb.WriteString(indentStr + "  Condition:\n")
	sb.WriteString(s.Condition.String(indent + 2))

	sb.WriteString(indentStr + "  Body:\n")
	for _, stmt := range s.Body {
		sb.WriteString(stmt.String(indent + 2))
	}

	return sb.String()
}
