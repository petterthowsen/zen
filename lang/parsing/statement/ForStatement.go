package statement

import (
	"strings"
	"zen/lang/common"
	"zen/lang/parsing/ast"
)

// ForStatement represents a for loop in the AST
type ForStatement struct {
	Location  *common.SourceLocation
	Init      ast.Statement  // Initialization (e.g., var i = 0)
	Condition ast.Expression // Loop condition (e.g., i < 10)
	Increment ast.Statement  // Increment (e.g., i = i + 1)
	Body      []ast.Statement
}

func NewForStatement(
	init ast.Statement,
	condition ast.Expression,
	increment ast.Statement,
	body []ast.Statement,
	location *common.SourceLocation,
) *ForStatement {
	return &ForStatement{
		Init:      init,
		Condition: condition,
		Increment: increment,
		Body:      body,
		Location:  location,
	}
}

func (s *ForStatement) Accept(visitor ast.Visitor) interface{} {
	return visitor.VisitForStatement(s)
}

func (s *ForStatement) GetLocation() *common.SourceLocation {
	return s.Location
}

func (s *ForStatement) IsStatement() {}

func (s *ForStatement) String(indent int) string {
	var sb strings.Builder
	indentStr := strings.Repeat("  ", indent)

	sb.WriteString(indentStr + "ForStatement\n")

	if s.Init != nil {
		sb.WriteString(indentStr + "  Init:\n")
		sb.WriteString(s.Init.String(indent + 2))
	}

	if s.Condition != nil {
		sb.WriteString(indentStr + "  Condition:\n")
		sb.WriteString(s.Condition.String(indent + 2))
	}

	if s.Increment != nil {
		sb.WriteString(indentStr + "  Increment:\n")
		sb.WriteString(s.Increment.String(indent + 2))
	}

	sb.WriteString(indentStr + "  Body:\n")
	for _, stmt := range s.Body {
		sb.WriteString(stmt.String(indent + 2))
	}

	return sb.String()
}
