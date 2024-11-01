package expression

import (
	"fmt"
	"strings"
	"zen/lang/common"
	"zen/lang/parsing/ast"
)

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
