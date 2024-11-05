package interpreter

import (
	"testing"
)

func TestBasicVariables(t *testing.T) {
	i := InterpretTestFile(t, "variables.zen")
	if i == nil {
		t.Fatal("Failed to interpret test file")
	}

	// Test basic variable declarations
	AssertValue(t, i, "x", 42)
	AssertValue(t, i, "y", "hello")
	AssertValue(t, i, "z", true)

	// Test constants
	AssertValue(t, i, "PI", 3.14159)
	AssertValue(t, i, "MAX_SIZE", 100)
	AssertValue(t, i, "GREETING", "hello")

	// Test nullable variables
	AssertValue(t, i, "nullableInt", nil)
	AssertValue(t, i, "nullableString", nil)
	AssertValue(t, i, "initializedNullable", "can be null")

	// Test variable modification
	AssertValue(t, i, "counter", 2) // (0 + 1) * 2

	// Test multiple variables
	AssertValue(t, i, "name", "Alice")
	AssertValue(t, i, "age", 30)
	AssertValue(t, i, "isStudent", false)
	AssertValue(t, i, "gpa", 3.85)

	// Test scope
	AssertValue(t, i, "global", "modified")
	AssertUndefined(t, i, "local")

	// Test expressions
	AssertValue(t, i, "sum", 15)
	AssertValue(t, i, "product", 50)
	AssertValue(t, i, "comparison", true)

	// Test string operations
	AssertValue(t, i, "fullName", "John Doe")

	// Test boolean operations
	AssertValue(t, i, "andResult", false)
	AssertValue(t, i, "orResult", true)
}

func TestVariableErrors(t *testing.T) {
	// Test redeclaration of variable
	AssertInterpretError(t, `
		var x = 1
		var x = 2
	`)

	// Test redeclaration of constant
	AssertInterpretError(t, `
		const X = 1
		const X = 2
	`)

	// Test modification of constant
	AssertInterpretError(t, `
		const X = 1
		X = 2
	`)

	// Test use of undefined variable
	AssertInterpretError(t, `
		var x = y
	`)

	// Test non-nullable without initialization
	AssertInterpretError(t, `
		var x : int
	`)

	// Test invalid type operations
	AssertInterpretError(t, `
		var x = "string"
		var y = 42
		var z = x + y
	`)
}

func TestVariableScoping(t *testing.T) {
	// Test nested scope variable shadowing
	i, err := InterpretString(`
		var x = "outer"
		if true {
			var x = "inner"
		}
		var final = x  // Should be "outer"
	`)
	if err != nil {
		t.Fatalf("Failed to interpret code: %v", err)
	}

	AssertValue(t, i, "x", "outer")
	AssertValue(t, i, "final", "outer")

	// Test modification of outer variable
	i, err = InterpretString(`
		var x = "initial"
		if true {
			x = "modified"  // Modifies outer x
		}
	`)
	if err != nil {
		t.Fatalf("Failed to interpret code: %v", err)
	}

	AssertValue(t, i, "x", "modified")

	// Test multiple nested scopes
	i, err = InterpretString(`
		var x = "level1"
		if true {
			var y = "level2"
			if true {
				var z = "level3"
				x = z
			}
		}
	`)
	if err != nil {
		t.Fatalf("Failed to interpret code: %v", err)
	}

	AssertValue(t, i, "x", "level3")
	AssertUndefined(t, i, "y")
	AssertUndefined(t, i, "z")
}

func TestNullableVariables(t *testing.T) {
	// Test nullable variable initialization
	i, err := InterpretString(`
		var x : int? = null
		var y : string? = "initial"
	`)
	if err != nil {
		t.Fatalf("Failed to interpret code: %v", err)
	}

	AssertValue(t, i, "x", nil)
	AssertValue(t, i, "y", "initial")

	// Test nullable variable assignment
	i, err = InterpretString(`
		var x : string? = null
		x = "value"
		x = null
	`)
	if err != nil {
		t.Fatalf("Failed to interpret code: %v", err)
	}

	AssertValue(t, i, "x", nil)

	// Test non-nullable to nullable assignment error
	AssertInterpretError(t, `
		var x = "non-null"
		var y : string? = null
		x = y
	`)
}
