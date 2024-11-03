package statement

import (
	"strings"
	"zen/lang/common"
	"zen/lang/parsing/ast"
)

// ContinueStatement represents a continue statement in a loop
type ContinueStatement struct {
	Location *common.SourceLocation
}

func NewContinueStatement(location *common.SourceLocation) *ContinueStatement {
	return &ContinueStatement{
		Location: location,
	}
}

func (s *ContinueStatement) Accept(visitor ast.Visitor) interface{} {
	return visitor.VisitContinueStatement(s)
}

func (s *ContinueStatement) GetLocation() *common.SourceLocation {
	return s.Location
}

func (s *ContinueStatement) IsStatement() {}

func (s *ContinueStatement) String(indent int) string {
	return strings.Repeat("  ", indent) + "Continue\n"
}
