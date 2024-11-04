package expression

import (
	"strings"
	"zen/lang/common"
	"zen/lang/parsing/ast"
)

type FuncParameterExpression struct {
	ast.Expression
	Name         string
	Type         ast.Expression
	IsNullable   bool
	Location     *common.SourceLocation
	DefaultValue ast.Expression
}

func NewFuncParameterExpression(name string, typ ast.Expression, isNullable bool, location *common.SourceLocation, defaultValue ast.Expression) *FuncParameterExpression {
	return &FuncParameterExpression{
		Name:         name,
		Type:         typ,
		IsNullable:   isNullable,
		Location:     location,
		DefaultValue: defaultValue,
	}
}

func (n *FuncParameterExpression) Accept(visitor ast.Visitor) interface{} {
	return visitor.VisitFuncParameterExpression(n)
}

func (n *FuncParameterExpression) GetLocation() *common.SourceLocation {
	return n.Location
}

func (n *FuncParameterExpression) IsExpression() {}

func (n *FuncParameterExpression) String(indent int) string {
	var sb strings.Builder
	indentStr := strings.Repeat("  ", indent)
	sb.WriteString(indentStr + "  FuncParameterExpression:\n")

	sb.WriteString(indentStr + "    Name: " + n.Name + "\n")
	sb.WriteString(indentStr + "    Type: " + n.Type.String(indent+2) + "\n")
	if n.DefaultValue != nil {
		sb.WriteString(indentStr + "    DefaultValue:")
		sb.WriteString(n.DefaultValue.String(indent+2) + "\n")
	}
	return sb.String()
}
