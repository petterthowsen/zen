package runtime

// EnvironmentInterface defines the interface for the environment
// that Zen uses to manage the scope chain
type EnvironmentInterface interface {

	// BeginScope creates a new scope with the current scope as its parent
	// Returns the new environment state
	BeginScope() *EnvironmentInterface

	// EndScope ends the current scope and returns to the parent scope
	// Returns error if attempting to end global scope
	EndScope() error

	// Define creates a new variable in the current scope
	Define(name string, value interface{}) error

	// DefineConst creates a new constant in the current scope
	DefineConst(name string, value interface{}) error

	// DefineNullable creates a new nullable variable in the current scope
	DefineNullable(name string, value interface{}) error

	// DefineGlobal creates a new variable in the global scope
	DefineGlobal(name string, value interface{}) error

	// Get retrieves a variable's value from the current scope chain
	Get(name string) (interface{}, error)

	// GetGlobal retrieves a variable's value from the global scope only
	GetGlobal(name string) (interface{}, error)

	// Assign updates an existing variable in the current scope chain
	Assign(name string, value interface{}) error

	// AssignGlobal updates an existing variable in the global scope only
	AssignGlobal(name string, value interface{}) error
}
