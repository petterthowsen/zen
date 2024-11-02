package statement

import (
	"strings"
	"zen/lang/common"
	"zen/lang/parsing/ast"
	"zen/lang/parsing/expression"
)

type FuncDeclaration struct {
	Name       string
	Parameters []expression.FuncParameterExpression
	ReturnType string
	Body       []ast.Statement
	location   *common.SourceLocation
}

func NewFuncDeclarationStatement(name string, parameters []expression.FuncParameterExpression, returnType string, body []ast.Statement, location *common.SourceLocation) *FuncDeclaration {
	return &FuncDeclaration{
		Name:       name,
		Parameters: parameters,
		ReturnType: returnType,
		Body:       body,
		location:   location,
	}
}

func (n *FuncDeclaration) Accept(visitor ast.Visitor) interface{} {
	return visitor.VisitFuncDeclaration(n)
}

func (n *FuncDeclaration) GetLocation() *common.SourceLocation {
	return n.location
}

// IsStatement implements ast.Statement interface
func (n *FuncDeclaration) IsStatement() {}

func NewFuncDeclaration(name string, parameters []expression.FuncParameterExpression, returnType string, body []ast.Statement, location *common.SourceLocation) *FuncDeclaration {
	return &FuncDeclaration{
		Name:       name,
		Parameters: parameters,
		ReturnType: returnType,
		Body:       body,
		location:   location,
	}
}

func (n *FuncDeclaration) String(indent int) string {
	var sb strings.Builder
	indentStr := strings.Repeat("  ", indent)

	// Write name
	sb.WriteString(indentStr + "FuncDeclaration " + n.Name + "\n")

	// Write parameters
	sb.WriteString(indentStr + "  Parameters:\n")
	for _, p := range n.Parameters {
		sb.WriteString(p.String(indent + 1))
	}

	// Write return type
	sb.WriteString(indentStr + "  ReturnType: " + n.ReturnType + "\n")

	// Write body
	sb.WriteString(indentStr + "  Body:\n")
	for _, stmt := range n.Body {
		sb.WriteString(stmt.String(indent+2) + "\n")
	}

	return sb.String()
}