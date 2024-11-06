package environment

import (
	"fmt"
	"zen/builtins/global"
	"zen/runtime/types"
)

// Environment manages the scope chain and provides high-level operations
// for variable management and scope creation
type Environment struct {
	// Global scope is the root of all scopes
	global *Scope
	// Current scope being executed
	current *Scope
}

// NewEnvironment creates a new environment with a global scope
func NewEnvironment() *Environment {
	global := NewScope(nil)
	return &Environment{
		global:  global,
		current: global,
	}
}

// BeginScope creates a new scope with the current scope as its parent
// Returns the new environment state
func (e *Environment) BeginScope() *Environment {
	e.current = NewScope(e.current)
	return e
}

// EndScope ends the current scope and returns to the parent scope
// Returns error if attempting to end global scope
func (e *Environment) EndScope() error {
	if e.current == e.global {
		return &ScopeError{Message: "Cannot end global scope"}
	}
	e.current = e.current.parent
	return nil
}

// Define creates a new variable in the current scope
func (e *Environment) Define(name string, value interface{}) error {
	return e.current.Define(name, value)
}

// DefineConst creates a new constant in the current scope
func (e *Environment) DefineConst(name string, value interface{}) error {
	return e.current.DefineConst(name, value)
}

// DefineNullable creates a new nullable variable in the current scope
func (e *Environment) DefineNullable(name string, value interface{}) error {
	return e.current.DefineNullable(name, value)
}

// DefineGlobal creates a new variable in the global scope
func (e *Environment) DefineGlobal(name string, value interface{}) error {
	return e.global.Define(name, value)
}

func (e *Environment) RegisterBuiltInFunctions() {
	printFn := types.BuiltinFunction{
		Name: "print",
		Parameters: []*types.FunctionParameterHint{
			types.NewFunctionParameterHint("str", types.TypeString, false),
		},
		ReturnType: nil,
		Async:      false,
		Func:       global.Print,
	}

	e.global.Define("print", printFn)
}

// Get retrieves a variable's value from the current scope chain
func (e *Environment) Get(name string) (types.Type, error) {
	return e.current.Get(name)
}

// GetGlobal retrieves a variable's value from the global scope only
func (e *Environment) GetGlobal(name string) (interface{}, error) {
	return e.global.Get(name)
}

// Assign updates an existing variable in the current scope chain
func (e *Environment) Assign(name string, value interface{}) error {
	found, err := e.current.Set(name, value)
	if err != nil {
		return err
	}
	if !found {
		return &UndefinedError{Name: name}
	}
	return nil
}

// AssignGlobal updates an existing variable in the global scope only
func (e *Environment) AssignGlobal(name string, value interface{}) error {
	found, err := e.global.Set(name, value)
	if err != nil {
		return err
	}
	if !found {
		return &UndefinedError{Name: name}
	}
	return nil
}

// ScopeError represents an error related to scope operations
type ScopeError struct {
	Message string
}

func (e *ScopeError) Error() string {
	return "Scope error: " + e.Message
}

// AssignmentError represents an error related to variable assignment
type AssignmentError struct {
	Message string
}

func (e *AssignmentError) Error() string {
	return fmt.Sprintf("Assignment error: %s", e.Message)
}
