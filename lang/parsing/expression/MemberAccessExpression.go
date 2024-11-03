package expression

import (
	"fmt"
	"strings"
	"zen/lang/common"
	"zen/lang/parsing/ast"
)

// MemberAccessExpression represents a member access operation (e.g., obj.prop)
type MemberAccessExpression struct {
	Object   ast.Expression // The object being accessed
	Property string         // The name of the property being accessed
	Location *common.SourceLocation
}

func NewMemberAccessExpression(
	object ast.Expression,
	property string,
	location *common.SourceLocation,
) *MemberAccessExpression {
	return &MemberAccessExpression{
		Object:   object,
		Property: property,
		Location: location,
	}
}

func (e *MemberAccessExpression) Accept(visitor ast.Visitor) interface{} {
	return visitor.VisitMemberAccess(e)
}

func (e *MemberAccessExpression) GetLocation() *common.SourceLocation {
	return e.Location
}

func (e *MemberAccessExpression) IsExpression() {}

func (e *MemberAccessExpression) String(indent int) string {
	var sb strings.Builder
	indentStr := strings.Repeat("  ", indent)

	sb.WriteString(fmt.Sprintf("%sMemberAccess(%s)\n", indentStr, e.Property))
	sb.WriteString(e.Object.String(indent + 1))

	return sb.String()
}
