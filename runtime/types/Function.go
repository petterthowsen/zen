package types

import "zen/lang/parsing/ast"

// FunctionParameter represents a parameter of a function
type FunctionParameter struct {
	Name     string
	Type     interface{}
	Nullable Bool
}

// Function represents a function which have a set of parameters and a return type
type Function struct {
	Name       string
	Parameters []FunctionParameter
	ReturnType interface{}
	Body       []*ast.Statement
	Async      bool
}

// implement Callable
func (f Function) IsCallable() bool {
	return true
}

// NewFunction creates a new Function
func NewFunction(name string, parameters []FunctionParameter, returnType interface{}, async bool) Function {
	return Function{
		Name:       name,
		Parameters: parameters,
		ReturnType: returnType,
		Async:      async,
	}
}
