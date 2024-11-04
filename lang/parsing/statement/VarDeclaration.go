package statement

import (
	"fmt"
	"strings"
	"zen/lang/common"
	"zen/lang/parsing/ast"
)

// VarDeclarationNode represents a variable or constant declaration in the AST
// it can optionally be constant (const), have a type, be nullable and/or have an initializer
type VarDeclarationNode struct {
	Name        string
	Type        ast.Expression // Can be either string literal (simple type) or ParametricType
	Initializer ast.Expression
	IsConstant  bool
	IsNullable  bool
	Location    *common.SourceLocation
}

// NewVarDeclarationNode creates a new VarDeclarationNode instance.
func NewVarDeclarationNode(
	name string,
	typ ast.Expression,
	initializer ast.Expression,
	isConstant bool,
	isNullable bool,
	location *common.SourceLocation,
) *VarDeclarationNode {
	return &VarDeclarationNode{
		Name:        name,
		Type:        typ,
		Initializer: initializer,
		IsConstant:  isConstant,
		IsNullable:  isNullable,
		Location:    location,
	}
}

func (n *VarDeclarationNode) Accept(visitor ast.Visitor) interface{} {
	return visitor.VisitVarDecl(n)
}

func (n *VarDeclarationNode) GetLocation() *common.SourceLocation {
	return n.Location
}

// IsStatement implements ast.Statement interface
func (n *VarDeclarationNode) IsStatement() {}

func (n *VarDeclarationNode) String(indent int) string {
	var sb strings.Builder
	indentStr := strings.Repeat("  ", indent)

	// Write declaration type (var/const)
	if n.IsConstant {
		sb.WriteString(indentStr + "Const Declaration\n")
	} else {
		sb.WriteString(indentStr + "Var Declaration\n")
	}

	// Write name and type
	sb.WriteString(fmt.Sprintf("%s  Name: %s\n", indentStr, n.Name))
	if n.Type != nil {
		sb.WriteString(fmt.Sprintf("%s  Type:\n", indentStr))
		sb.WriteString(n.Type.String(indent+2) + "\n")
	}

	// Write initializer if present
	if n.Initializer != nil {
		sb.WriteString(fmt.Sprintf("%s  Initializer:\n", indentStr))
		sb.WriteString(n.Initializer.String(indent + 2))
	}

	return sb.String()
}
