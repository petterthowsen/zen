package expression

import (
	"fmt"
	"strings"
	"zen/lang/common"
	"zen/lang/parsing/ast"
)

// AwaitExpression represents an await expression in the AST
type AwaitExpression struct {
	Expression ast.Expression
	Location   *common.SourceLocation
}

func NewAwaitExpression(expr ast.Expression, location *common.SourceLocation) *AwaitExpression {
	return &AwaitExpression{
		Expression: expr,
		Location:   location,
	}
}

func (e *AwaitExpression) Accept(visitor ast.Visitor) interface{} {
	return visitor.VisitAwait(e)
}

func (e *AwaitExpression) GetLocation() *common.SourceLocation {
	return e.Location
}

func (e *AwaitExpression) IsExpression() {}

func (e *AwaitExpression) String(indent int) string {
	return fmt.Sprintf("%sAwait:\n%s", strings.Repeat("  ", indent), e.Expression.String(indent+1))
}
