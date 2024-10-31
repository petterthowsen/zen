package ast

import (
	"strings"
	"zen/lang/common"
)

// Node represents a node in the AST
type Node interface {
	// Accept implements the visitor pattern
	Accept(visitor Visitor) interface{}
	// GetLocation returns the source location of this node
	GetLocation() *common.SourceLocation
	// String returns a string representation of the node with proper indentation
	String(indent int) string
}

// Expression represents an expression node in the AST
type Expression interface {
	Node
	IsExpression() // Marker method must be exported for cross-package use
}

// Statement represents a statement node in the AST
type Statement interface {
	Node
	IsStatement() // Marker method must be exported for cross-package use
}

// Visitor interface for implementing the visitor pattern
type Visitor interface {
	// Add visitor methods as we implement more node types
	VisitProgram(node *ProgramNode) interface{}
	VisitVarDecl(node Statement) interface{}
	VisitLiteral(node Expression) interface{}
}

// ProgramNode represents the root node of the AST
type ProgramNode struct {
	Statements []Statement
	Location   *common.SourceLocation
}

func NewProgramNode(statements []Statement) *ProgramNode {
	var location *common.SourceLocation
	if len(statements) > 0 {
		location = statements[0].GetLocation()
	}

	return &ProgramNode{
		Statements: statements,
		Location:   location,
	}
}

func (n *ProgramNode) Accept(visitor Visitor) interface{} {
	return visitor.VisitProgram(n)
}

func (n *ProgramNode) GetLocation() *common.SourceLocation {
	return n.Location
}

func (n *ProgramNode) String(indent int) string {
	var sb strings.Builder
	indentStr := strings.Repeat("  ", indent)

	sb.WriteString(indentStr + "Program\n")
	for _, stmt := range n.Statements {
		sb.WriteString(stmt.String(indent + 1))
	}

	return sb.String()
}

// Helper function to create indented strings
func Indent(indent int) string {
	return strings.Repeat("  ", indent)
}
