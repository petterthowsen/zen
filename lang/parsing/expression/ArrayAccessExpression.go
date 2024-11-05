package expression

import (
	"fmt"
	"strings"
	"zen/lang/common"
	"zen/lang/parsing/ast"
)

// ArrayAccessExpression represents array indexing in the AST (e.g., array[index])
type ArrayAccessExpression struct {
	Array    ast.Expression
	Index    ast.Expression
	Location *common.SourceLocation
}

func NewArrayAccessExpression(array ast.Expression, index ast.Expression, location *common.SourceLocation) *ArrayAccessExpression {
	return &ArrayAccessExpression{
		Array:    array,
		Index:    index,
		Location: location,
	}
}

func (e *ArrayAccessExpression) Accept(visitor ast.Visitor) interface{} {
	return visitor.VisitArrayAccess(e)
}

func (e *ArrayAccessExpression) GetLocation() *common.SourceLocation {
	return e.Location
}

func (e *ArrayAccessExpression) IsExpression() {}

func (e *ArrayAccessExpression) String(indent int) string {
	return fmt.Sprintf("%sArrayAccess:\n%sArray:\n%s%sIndex:\n%s",
		strings.Repeat("  ", indent),
		strings.Repeat("  ", indent+1),
		e.Array.String(indent+2),
		strings.Repeat("  ", indent+1),
		e.Index.String(indent+2))
}
