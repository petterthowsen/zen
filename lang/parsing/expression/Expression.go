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

// IdentifierExpression represents a variable reference in the AST
type IdentifierExpression struct {
	Name     string
	Location *common.SourceLocation
}

func NewIdentifierExpression(name string, location *common.SourceLocation) *IdentifierExpression {
	return &IdentifierExpression{
		Name:     name,
		Location: location,
	}
}

func (e *IdentifierExpression) Accept(visitor ast.Visitor) interface{} {
	return visitor.VisitIdentifier(e)
}

func (e *IdentifierExpression) GetLocation() *common.SourceLocation {
	return e.Location
}

func (e *IdentifierExpression) IsExpression() {}

func (e *IdentifierExpression) String(indent int) string {
	return fmt.Sprintf("%sIdentifier: %s\n", strings.Repeat("  ", indent), e.Name)
}

// CallExpression represents a function call in the AST
type CallExpression struct {
	Callee    ast.Expression
	Arguments []ast.Expression
	Location  *common.SourceLocation
}

func NewCallExpression(callee ast.Expression, arguments []ast.Expression, location *common.SourceLocation) *CallExpression {
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
