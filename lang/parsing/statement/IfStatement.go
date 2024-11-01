package statement

import (
	"strings"
	"zen/lang/common"
	"zen/lang/parsing/ast"
)

// IfStatement represents an if statement in the AST
type IfStatement struct {
	Location  *common.SourceLocation
	Condition ast.Expression
	Body      []ast.Statement
}

func (i *IfStatement) Accept(visitor ast.Visitor) interface{} {
	return visitor.VisitIfStatement(i)
}

func (i *IfStatement) GetLocation() *common.SourceLocation {
	return i.Location
}

func (i *IfStatement) IsStatement() {}

func (i *IfStatement) String(indent int) string {
	var builder strings.Builder
	indentStr := strings.Repeat("  ", indent)

	builder.WriteString(indentStr + "If\n")
	builder.WriteString(indentStr + "  Condition:\n")
	builder.WriteString(i.Condition.String(indent + 2))
	builder.WriteString(indentStr + "  Body:\n")
	for _, stmt := range i.Body {
		builder.WriteString(stmt.String(indent + 2))
	}

	return builder.String()
}
