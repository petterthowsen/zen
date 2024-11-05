package environment

// VarInfo holds information about a variable
type VarInfo struct {
	value      interface{}
	isConstant bool
	isNullable bool
}

// Scope represents a single scope level in the environment chain
type Scope struct {
	// Parent scope in the chain, nil for global scope
	parent *Scope
	// Variables stored in this scope
	variables map[string]VarInfo
}

// NewScope creates a new scope with an optional parent
func NewScope(parent *Scope) *Scope {
	return &Scope{
		parent:    parent,
		variables: make(map[string]VarInfo),
	}
}

// Define creates a new variable in the current scope
func (s *Scope) Define(name string, value interface{}) error {
	if _, exists := s.variables[name]; exists {
		return &RedefinitionError{Name: name}
	}
	s.variables[name] = VarInfo{
		value:      value,
		isConstant: false,
		isNullable: value == nil,
	}
	return nil
}

// DefineConst creates a new constant in the current scope
func (s *Scope) DefineConst(name string, value interface{}) error {
	if _, exists := s.variables[name]; exists {
		return &RedefinitionError{Name: name}
	}
	if value == nil {
		return &AssignmentError{Message: "Cannot define constant with null value"}
	}
	s.variables[name] = VarInfo{
		value:      value,
		isConstant: true,
		isNullable: false,
	}
	return nil
}

// DefineNullable creates a new nullable variable in the current scope
func (s *Scope) DefineNullable(name string, value interface{}) error {
	if _, exists := s.variables[name]; exists {
		return &RedefinitionError{Name: name}
	}
	s.variables[name] = VarInfo{
		value:      value,
		isConstant: false,
		isNullable: true,
	}
	return nil
}

// Get retrieves a variable's value from this scope or any parent scope
func (s *Scope) Get(name string) (interface{}, error) {
	if info, exists := s.variables[name]; exists {
		return info.value, nil
	}

	if s.parent != nil {
		return s.parent.Get(name)
	}

	return nil, &UndefinedError{Name: name}
}

// GetInfo retrieves a variable's full info from this scope or any parent scope
func (s *Scope) GetInfo(name string) (*VarInfo, error) {
	if info, exists := s.variables[name]; exists {
		return &info, nil
	}

	if s.parent != nil {
		return s.parent.GetInfo(name)
	}

	return nil, &UndefinedError{Name: name}
}

// Set updates a variable's value in this scope or any parent scope
// Returns true if the variable was found and updated, false otherwise
func (s *Scope) Set(name string, value interface{}) (bool, error) {
	// First check if variable exists in current scope
	if info, exists := s.variables[name]; exists {
		if info.isConstant {
			return false, &AssignmentError{Message: "Cannot assign to constant"}
		}
		if !info.isNullable && value == nil {
			return false, &AssignmentError{Message: "Cannot assign null to non-nullable variable"}
		}
		info.value = value
		s.variables[name] = info
		return true, nil
	}

	// If not in current scope and we have a parent, try to set in parent
	if s.parent != nil {
		found, err := s.parent.Set(name, value)
		if err != nil {
			return false, err
		}
		return found, nil
	}

	return false, nil
}

// RedefinitionError occurs when attempting to define a variable that already exists
type RedefinitionError struct {
	Name string
}

func (e *RedefinitionError) Error() string {
	return "Cannot redefine variable: " + e.Name
}

// UndefinedError occurs when attempting to access a variable that doesn't exist
type UndefinedError struct {
	Name string
}

func (e *UndefinedError) Error() string {
	return "Undefined variable: " + e.Name
}
