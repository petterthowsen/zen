package types

import (
	"zen/lang/parsing/ast"
)

// FunctionParameterHint represents a parameter for a function
type FunctionParameterHint struct {
	Name     string
	Type     Type
	Nullable bool
}

// NewFunctionParameterHint creates a new FunctionParameterHint
func NewFunctionParameterHint(name string, typ Type, nullable bool) *FunctionParameterHint {
	return &FunctionParameterHint{
		Name:     name,
		Type:     typ,
		Nullable: nullable,
	}
}

// UserFunction represents a function which have a set of parameters and a return type
type UserFunction struct {
	Name       string
	Parameters []FunctionParameterHint
	ReturnType interface{}
	Body       []*ast.Statement
	Async      bool
}

// IsCallable implement Callable
func (f UserFunction) IsCallable() bool {
	return true
}

// NewUserFunction creates a new UserFunction
func NewUserFunction(name string, parameters []FunctionParameterHint, returnType interface{}, async bool) UserFunction {
	return UserFunction{
		Name:       name,
		Parameters: parameters,
		ReturnType: returnType,
		Async:      async,
	}
}
