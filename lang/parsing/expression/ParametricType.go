package expression

import (
	"fmt"
	"strings"
	"zen/lang/common"
	"zen/lang/parsing/ast"
)

// Parameter represents a type parameter which can be either a type name or an integer literal
type Parameter struct {
	// If IsType is true, Value is either:
	// - a string (type name)
	// - a BasicType
	// - a ParametricType
	// If IsType is false, Value is an int64 (size/range constraint)
	Value    interface{}
	IsType   bool
	Location *common.SourceLocation
}

// ParametricType represents a generic type with parameters, e.g., Array<int, 5> or Map<string, any>
type ParametricType struct {
	// The base type name (e.g., "Array" or "Map")
	BaseType string
	// The type parameters (e.g., [int, 5] for Array<int, 5>)
	Parameters []Parameter
	Location   *common.SourceLocation
}

func NewParametricType(baseType string, parameters []Parameter, location *common.SourceLocation) *ParametricType {
	return &ParametricType{
		BaseType:   baseType,
		Parameters: parameters,
		Location:   location,
	}
}

func (p *ParametricType) Accept(visitor ast.Visitor) interface{} {
	return visitor.VisitParametricType(p)
}

func (p *ParametricType) GetLocation() *common.SourceLocation {
	return p.Location
}

func (p *ParametricType) IsExpression() {}

func (p *ParametricType) String(indent int) string {
	params := make([]string, len(p.Parameters))
	for i, param := range p.Parameters {
		if param.IsType {
			switch v := param.Value.(type) {
			case string:
				params[i] = v
			case *BasicType:
				params[i] = v.Name
			case *ParametricType:
				params[i] = v.String(0) // Don't indent nested types
			default:
				params[i] = fmt.Sprintf("<%T>", v) // For debugging
			}
		} else {
			params[i] = fmt.Sprintf("%d", param.Value.(int64))
		}
	}
	return strings.Repeat("  ", indent) + fmt.Sprintf("%s<%s>", p.BaseType, strings.Join(params, ", "))
}
