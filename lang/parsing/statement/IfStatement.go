package statement

import (
	"strings"
	"zen/lang/common"
	"zen/lang/parsing/ast"
)

// IfStatement represents an if statement in the AST
type IfStatement struct {
	Location         *common.SourceLocation
	PrimaryCondition ast.Expression
	PrimaryBlock     []ast.Statement
	ElseIfBlocks     []*IfConditionBlock
	ElseBlock        []ast.Statement
}

func NewIfStatement(condition ast.Expression, body []ast.Statement, elseIfBlocks []*IfConditionBlock, elseBlock []ast.Statement, location *common.SourceLocation) *IfStatement {
	return &IfStatement{
		Location:         location,
		PrimaryCondition: condition,
		PrimaryBlock:     body,
		ElseIfBlocks:     elseIfBlocks,
		ElseBlock:        elseBlock,
	}
}

func (i *IfStatement) Accept(visitor ast.Visitor) interface{} {
	return visitor.VisitIfStatement(i)
}

func (i *IfStatement) GetLocation() *common.SourceLocation {
	return i.Location
}

func (i *IfStatement) IsStatement() {}

func (i *IfStatement) String(indent int) string {
	var builder strings.Builder
	indentStr := strings.Repeat("  ", indent)

	builder.WriteString(indentStr + "If\n")

	// Write primary condition
	builder.WriteString(indentStr + "  Primary Condition:\n")
	builder.WriteString(i.PrimaryCondition.String(indent + 2))

	// Write primary block
	builder.WriteString(indentStr + "  Primary Block:\n")
	for _, stmt := range i.PrimaryBlock {
		builder.WriteString(stmt.String(indent + 2))
	}

	// Write else if blocks
	builder.WriteString(indentStr + "  Else If Blocks:\n")
	for _, elseIfBlock := range i.ElseIfBlocks {
		builder.WriteString(elseIfBlock.String(indent + 2))
	}

	// write else block
	builder.WriteString(indentStr + "  Else Block:\n")
	for _, stmt := range i.ElseBlock {
		builder.WriteString(stmt.String(indent + 2))
	}

	return builder.String()
}
