package expression

import (
	"fmt"
	"strings"
	"zen/lang/common"
	"zen/lang/parsing/ast"
)

// MapAccessExpression represents map key access in the AST (e.g., map{"key"})
type MapAccessExpression struct {
	Map      ast.Expression
	Key      ast.Expression
	Location *common.SourceLocation
}

func NewMapAccessExpression(mapExpr ast.Expression, key ast.Expression, location *common.SourceLocation) *MapAccessExpression {
	return &MapAccessExpression{
		Map:      mapExpr,
		Key:      key,
		Location: location,
	}
}

func (e *MapAccessExpression) Accept(visitor ast.Visitor) interface{} {
	return visitor.VisitMapAccess(e)
}

func (e *MapAccessExpression) GetLocation() *common.SourceLocation {
	return e.Location
}

func (e *MapAccessExpression) IsExpression() {}

func (e *MapAccessExpression) String(indent int) string {
	return fmt.Sprintf("%sMapAccess:\n%sMap:\n%s%sKey:\n%s",
		strings.Repeat("  ", indent),
		strings.Repeat("  ", indent+1),
		e.Map.String(indent+2),
		strings.Repeat("  ", indent+1),
		e.Key.String(indent+2))
}
