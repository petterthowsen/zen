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
	indentStr := strings.Repeat("  ", indent)
	var leftStr, rightStr string

	// Handle nil values gracefully
	if e.Left == nil {
		leftStr = indentStr + "  <nil>\n"
	} else {
		leftStr = e.Left.String(indent + 1)
	}

	if e.Right == nil {
		rightStr = indentStr + "  <nil>\n"
	} else {
		rightStr = e.Right.String(indent + 1)
	}

	return fmt.Sprintf("%sBinary: %s\n%s%s",
		indentStr,
		e.Operator,
		leftStr,
		rightStr)
}
