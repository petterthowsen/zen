package interpreter

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"zen/interpreter"
	"zen/tests/parsing"
)

// InterpretString parses and interprets source code string
func InterpretString(source string) (*interpreter.Interpreter, error) {
	program, errors := parsing.ParseString(source)
	if len(errors) > 0 {
		return nil, errors[0]
	}

	i := interpreter.NewInterpreter()
	err := i.Execute(program)
	return i, err
}

// getTestDataPath returns the absolute path to the test data file
func getTestDataPath(filename string) string {
	_, currentFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(currentFile)
	return filepath.Join(dir, filename)
}

// InterpretTestFile loads and interprets a test file
func InterpretTestFile(t *testing.T, filename string) *interpreter.Interpreter {
	t.Helper()
	// Load the test file
	path := getTestDataPath(filename)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", path, err)
	}

	// Parse and interpret
	program, errors := parsing.ParseFile(path, string(content))
	if len(errors) > 0 {
		t.Errorf("Parser errors:")
		for _, err := range errors {
			t.Errorf("  %s", err.Error())
		}
		return nil
	}

	if program == nil {
		t.Error("Expected program node, got nil")
		return nil
	}

	i := interpreter.NewInterpreter()
	if err := i.Execute(program); err != nil {
		t.Errorf("Interpreter error: %v", err)
		return nil
	}

	return i
}

// AssertValue checks if a variable has the expected value
func AssertValue(t *testing.T, i *interpreter.Interpreter, name string, expected interface{}) {
	t.Helper()
	result, err := i.GetValue(name)
	if err != nil {
		t.Errorf("Variable %s should be defined: %v", name, err)
		return
	}

	// Convert numeric types for comparison
	switch expected := expected.(type) {
	case int:
		switch val := result.(type) {
		case int32:
			if int32(expected) != val {
				t.Errorf("Variable %s: expected %v, got %v", name, expected, val)
			}
		case int64:
			if int64(expected) != val {
				t.Errorf("Variable %s: expected %v, got %v", name, expected, val)
			}
		default:
			t.Errorf("Variable %s: expected int, got %T", name, result)
		}
		return
	case float32:
		if val, ok := result.(float32); ok {
			if expected != val {
				t.Errorf("Variable %s: expected %v, got %v", name, expected, val)
			}
			return
		}
		t.Errorf("Variable %s: expected float32, got %T", name, result)
		return
	case float64:
		if val, ok := result.(float64); ok {
			if expected != val {
				t.Errorf("Variable %s: expected %v, got %v", name, expected, val)
			}
			return
		}
		t.Errorf("Variable %s: expected float64, got %T", name, result)
		return
	case string:
		if val, ok := result.(string); ok {
			if expected != val {
				t.Errorf("Variable %s: expected %q, got %q", name, expected, val)
			}
			return
		}
		t.Errorf("Variable %s: expected string, got %T", name, result)
		return
	case bool:
		if val, ok := result.(bool); ok {
			if expected != val {
				t.Errorf("Variable %s: expected %v, got %v", name, expected, val)
			}
			return
		}
		t.Errorf("Variable %s: expected bool, got %T", name, result)
		return
	case nil:
		if result != nil {
			t.Errorf("Variable %s: expected nil, got %v", name, result)
		}
		return
	default:
		t.Errorf("Variable %s: unsupported expected type %T", name, expected)
	}
}

// AssertUndefined checks if a variable is undefined
func AssertUndefined(t *testing.T, i *interpreter.Interpreter, name string) {
	t.Helper()
	_, err := i.GetValue(name)
	if err == nil {
		t.Errorf("Variable %s should not be defined", name)
	}
}

// AssertInterpretError checks if interpreting code produces an error
func AssertInterpretError(t *testing.T, source string) {
	t.Helper()
	_, err := InterpretString(source)
	if err == nil {
		t.Errorf("Expected interpretation error for:\n\n%s\n", source)
	}
}
