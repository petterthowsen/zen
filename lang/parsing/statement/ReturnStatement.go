package statement

import (
	"strings"
	"zen/lang/common"
	"zen/lang/parsing/ast"
)

type ReturnStatmenet struct {
	Location   *common.SourceLocation
	Expression ast.Expression
}

func NewReturnStatement(expression ast.Expression, location *common.SourceLocation) *ReturnStatmenet {
	return &ReturnStatmenet{
		Location:   location,
		Expression: expression,
	}
}

func (s *ReturnStatmenet) Accept(visitor ast.Visitor) interface{} {
	return visitor.VisitReturnStatement(s)
}

func (s *ReturnStatmenet) GetLocation() *common.SourceLocation {
	return s.Location
}

func (s *ReturnStatmenet) IsStatement() {
}

func (s *ReturnStatmenet) String(indent int) string {
	var sb strings.Builder
	indentStr := strings.Repeat("  ", indent)

	sb.WriteString(indentStr + "Return\n")
	if s.Expression != nil {
		sb.WriteString(s.Expression.String(indent + 1))
	}

	return sb.String()
}
