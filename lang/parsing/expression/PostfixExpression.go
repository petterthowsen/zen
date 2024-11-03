package expression

import (
	"fmt"
	"strings"
	"zen/lang/common"
	"zen/lang/parsing/ast"
)

// PostfixExpression represents a postfix operation (e.g., i++, i--)
type PostfixExpression struct {
	Operand  ast.Expression
	Operator string // ++ or --
	Location *common.SourceLocation
}

func NewPostfixExpression(operand ast.Expression, operator string, location *common.SourceLocation) *PostfixExpression {
	return &PostfixExpression{
		Operand:  operand,
		Operator: operator,
		Location: location,
	}
}

func (e *PostfixExpression) Accept(visitor ast.Visitor) interface{} {
	return visitor.VisitPostfix(e)
}

func (e *PostfixExpression) GetLocation() *common.SourceLocation {
	return e.Location
}

func (e *PostfixExpression) IsExpression() {}

func (e *PostfixExpression) String(indent int) string {
	var sb strings.Builder
	indentStr := strings.Repeat("  ", indent)

	sb.WriteString(fmt.Sprintf("%sPostfixExpression(%s)\n", indentStr, e.Operator))
	sb.WriteString(e.Operand.String(indent + 1))

	return sb.String()
}
