package statement

import (
	"strings"
	"zen/lang/common"
	"zen/lang/parsing/ast"
)

// BreakStatement represents a break statement in a loop
type BreakStatement struct {
	Location *common.SourceLocation
}

func NewBreakStatement(location *common.SourceLocation) *BreakStatement {
	return &BreakStatement{
		Location: location,
	}
}

func (s *BreakStatement) Accept(visitor ast.Visitor) interface{} {
	return visitor.VisitBreakStatement(s)
}

func (s *BreakStatement) GetLocation() *common.SourceLocation {
	return s.Location
}

func (s *BreakStatement) IsStatement() {}

func (s *BreakStatement) String(indent int) string {
	return strings.Repeat("  ", indent) + "Break\n"
}
