package expression

import (
	"fmt"
	"strings"
	"zen/lang/common"
	"zen/lang/parsing/ast"
)

// LiteralExpression represents a literal value in the AST
type LiteralExpression struct {
	Value    interface{}
	Location *common.SourceLocation
}

func NewLiteralExpression(value interface{}, location *common.SourceLocation) *LiteralExpression {
	return &LiteralExpression{
		Value:    value,
		Location: location,
	}
}

func (e *LiteralExpression) Accept(visitor ast.Visitor) interface{} {
	return visitor.VisitLiteral(e)
}

func (e *LiteralExpression) GetLocation() *common.SourceLocation {
	return e.Location
}

// IsExpression implements ast.Expression interface
func (e *LiteralExpression) IsExpression() {}

func (e *LiteralExpression) String(indent int) string {
	return fmt.Sprintf("%sLiteral: %v\n", strings.Repeat("  ", indent), e.Value)
}
