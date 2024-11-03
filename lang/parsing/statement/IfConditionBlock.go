package statement

import (
	"strings"
	"zen/lang/common"
	"zen/lang/parsing/ast"
)

// IfConditionBlock represents: [condition] { [body] }
// see IfStatement for the statement that encapsulates a complete if statement with else if and else blocks.
type IfConditionBlock struct {
	Location  *common.SourceLocation
	Condition ast.Expression
	Body      []ast.Statement
}

func NewIfConditionBlock(condition ast.Expression, body []ast.Statement, location *common.SourceLocation) *IfConditionBlock {
	return &IfConditionBlock{
		Condition: condition,
		Body:      body,
		Location:  location,
	}
}

func (i *IfConditionBlock) Accept(visitor ast.Visitor) interface{} {
	return visitor.VisitIfStatement(i)
}

func (i *IfConditionBlock) GetLocation() *common.SourceLocation {
	return i.Location
}

func (i *IfConditionBlock) IsStatement() {}

func (i *IfConditionBlock) String(indent int) string {
	var builder strings.Builder
	indentStr := strings.Repeat("  ", indent)

	builder.WriteString(indentStr + "Condition:\n")
	builder.WriteString(i.Condition.String(indent + 2))

	builder.WriteString(indentStr + "Body:\n")
	for _, stmt := range i.Body {
		builder.WriteString(stmt.String(indent + 2))
	}

	return builder.String()
}
