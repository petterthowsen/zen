package expression

import (
	"fmt"
	"strings"
	"zen/lang/common"
	"zen/lang/parsing/ast"
)

// ArrayLiteralExpression represents an array literal in the AST
type ArrayLiteralExpression struct {
	Elements []ast.Expression
	Location *common.SourceLocation
}

func NewArrayLiteralExpression(elements []ast.Expression, location *common.SourceLocation) *ArrayLiteralExpression {
	return &ArrayLiteralExpression{
		Elements: elements,
		Location: location,
	}
}

func (e *ArrayLiteralExpression) Accept(visitor ast.Visitor) interface{} {
	return visitor.VisitArrayLiteral(e)
}

func (e *ArrayLiteralExpression) GetLocation() *common.SourceLocation {
	return e.Location
}

func (e *ArrayLiteralExpression) IsExpression() {}

func (e *ArrayLiteralExpression) String(indent int) string {
	var elements []string
	for _, elem := range e.Elements {
		elements = append(elements, elem.String(indent+1))
	}
	return fmt.Sprintf("%sArrayLiteral:\n%s", strings.Repeat("  ", indent), strings.Join(elements, ""))
}
