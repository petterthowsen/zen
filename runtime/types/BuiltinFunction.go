package types

import "zen/runtime"

// BuiltinFunction represents a built-in function in zen that runs a go function within the zen runtime
type BuiltinFunction struct {
	Name       string
	Parameters []FunctionParameter
	ReturnType interface{}
	Async      bool
	Func       func(env *runtime.EnvironmentInterface, args ...Value) (Value, error)
}

// implement Callable
func (f *BuiltinFunction) IsCallable() bool {
	return true
}

// Call calls the underlying function with 0 or more Value as arguments
// and returns the result of the call
func (f *BuiltinFunction) Call(env *runtime.EnvironmentInterface, args ...Value) (Value, error) {
	return f.Func(env, args...)
}

// NewBuiltinFunction creates a new BuiltinFunction
func NewBuiltinFunction(name string, parameters []FunctionParameter, returnType interface{}, async bool, funcFunc func(env *runtime.EnvironmentInterface, args ...Value) (Value, error)) BuiltinFunction {
	return BuiltinFunction{
		Name:       name,
		Parameters: parameters,
		ReturnType: returnType,
		Async:      async,
		Func:       funcFunc,
	}
}
