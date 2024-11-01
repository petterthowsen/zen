package statement

import (
	"strings"
	"zen/lang/common"
	"zen/lang/parsing/ast"
)

// ExpressionStatement represents an expression used as a statement
type ExpressionStatement struct {
	Location   *common.SourceLocation
	Expression ast.Expression
}

func (e *ExpressionStatement) Accept(visitor ast.Visitor) interface{} {
	return visitor.VisitExpressionStatement(e)
}

func (e *ExpressionStatement) GetLocation() *common.SourceLocation {
	return e.Location
}

func (e *ExpressionStatement) IsStatement() {}

func (e *ExpressionStatement) String(indent int) string {
	var builder strings.Builder
	indentStr := strings.Repeat("  ", indent)

	builder.WriteString(indentStr + "ExpressionStatement\n")
	builder.WriteString(e.Expression.String(indent + 1))

	return builder.String()
}
