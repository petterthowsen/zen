package expression

import (
	"fmt"
	"strings"
	"zen/lang/common"
	"zen/lang/parsing/ast"
)

// CallExpression represents a function call in the AST
type CallExpression struct {
	Callee    ast.Expression
	Arguments []ast.Expression
	Location  *common.SourceLocation
}

func NewCallExpression(callee *IdentifierExpression, arguments []ast.Expression, location *common.SourceLocation) *CallExpression {
	return &CallExpression{
		Callee:    callee,
		Arguments: arguments,
		Location:  location,
	}
}

func (e *CallExpression) Accept(visitor ast.Visitor) interface{} {
	return visitor.VisitCall(e)
}

func (e *CallExpression) GetLocation() *common.SourceLocation {
	return e.Location
}

func (e *CallExpression) IsExpression() {}

func (e *CallExpression) String(indent int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%sCall\n", strings.Repeat("  ", indent)))
	sb.WriteString(fmt.Sprintf("%sCallee:\n%s", strings.Repeat("  ", indent+1), e.Callee.String(indent+2)))
	if len(e.Arguments) > 0 {
		sb.WriteString(fmt.Sprintf("%sArguments:\n", strings.Repeat("  ", indent+1)))
		for _, arg := range e.Arguments {
			sb.WriteString(arg.String(indent + 2))
		}
	}
	return sb.String()
}
