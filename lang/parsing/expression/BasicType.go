package expression

import (
	"strings"
	"zen/lang/common"
	"zen/lang/parsing/ast"
)

// BasicType represents a simple type name, which can be either:
// - A primitive type (string, int, float64, etc.) from keywords
// - A class type (Array, Map, MyClass, etc.) from identifiers
type BasicType struct {
	Name     string
	Location *common.SourceLocation
}

func NewBasicType(name string, location *common.SourceLocation) *BasicType {
	return &BasicType{
		Name:     name,
		Location: location,
	}
}

func (t *BasicType) Accept(visitor ast.Visitor) interface{} {
	return visitor.VisitBasicType(t)
}

func (t *BasicType) GetLocation() *common.SourceLocation {
	return t.Location
}

func (t *BasicType) IsExpression() {}

func (t *BasicType) String(indent int) string {
	return strings.Repeat("  ", indent) + t.Name
}
