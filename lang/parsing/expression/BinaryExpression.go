package expression

import (
	"fmt"
	"strings"
	"zen/lang/common"
	"zen/lang/parsing/ast"
)

// BinaryExpression represents a binary operation in the AST
type BinaryExpression struct {
	Left     ast.Expression
	Operator string
	Right    ast.Expression
	Location *common.SourceLocation
}

func NewBinaryExpression(left ast.Expression, operator string, right ast.Expression, location *common.SourceLocation) *BinaryExpression {
	return &BinaryExpression{
		Left:     left,
		Operator: operator,
		Right:    right,
		Location: location,
	}
}

func (e *BinaryExpression) Accept(visitor ast.Visitor) interface{} {
	return visitor.VisitBinary(e)
}

func (e *BinaryExpression) GetLocation() *common.SourceLocation {
	return e.Location
}

func (e *BinaryExpression) IsExpression() {}

func (e *BinaryExpression) String(indent int) string {
	return fmt.Sprintf("%sBinary: %s\n%s%s",
		strings.Repeat("  ", indent),
		e.Operator,
		e.Left.String(indent+1),
		e.Right.String(indent+1))
}
