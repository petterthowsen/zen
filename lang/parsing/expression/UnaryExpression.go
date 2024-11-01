package expression

import (
	"fmt"
	"strings"
	"zen/lang/common"
	"zen/lang/parsing/ast"
)

// UnaryExpression represents a unary operation in the AST
type UnaryExpression struct {
	Operator   string
	Expression ast.Expression
	Location   *common.SourceLocation
}

func NewUnaryExpression(operator string, expression ast.Expression, location *common.SourceLocation) *UnaryExpression {
	return &UnaryExpression{
		Operator:   operator,
		Expression: expression,
		Location:   location,
	}
}

func (e *UnaryExpression) Accept(visitor ast.Visitor) interface{} {
	return visitor.VisitUnary(e)
}

func (e *UnaryExpression) GetLocation() *common.SourceLocation {
	return e.Location
}

func (e *UnaryExpression) IsExpression() {}

func (e *UnaryExpression) String(indent int) string {
	return fmt.Sprintf("%sUnary: %s\n%s",
		strings.Repeat("  ", indent),
		e.Operator,
		e.Expression.String(indent+1))
}
