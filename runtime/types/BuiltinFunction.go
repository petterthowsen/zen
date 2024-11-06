package types

import "zen/runtime"

// BuiltinFunction represents a built-in function in zen that runs a go function within the zen runtime
type BuiltinFunction struct {
	Name       string
	Parameters []*FunctionParameterHint
	ReturnType interface{}
	Async      bool
	Func       func(env *runtime.EnvironmentInterface, args map[string]Value) (Value, error)
}

func (f *BuiltinFunction) String() string {

}

func (f *BuiltinFunction) Type() Type {
	return TypeFunction
}

// IsCallable implement Callable
func (f *BuiltinFunction) IsCallable() bool {
	return true
}

// Call calls the underlying function with 0 or more Value parameters
// Returns a Value and an error
func (f *BuiltinFunction) Call(env *runtime.EnvironmentInterface, params map[string]Value) (Value, error) {
	return f.Func(env, params)
}

// NewBuiltinFunction creates a new BuiltinFunction
func NewBuiltinFunction(name string, parameters []*FunctionParameterHint, returnType interface{}, async bool, funcFunc func(env *runtime.EnvironmentInterface, params map[string]Value) (Value, error)) BuiltinFunction {
	return BuiltinFunction{
		Name:       name,
		Parameters: parameters,
		ReturnType: returnType,
		Async:      async,
		Func:       funcFunc,
	}
}
