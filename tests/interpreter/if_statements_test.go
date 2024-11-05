package interpreter

import (
	"testing"
)

func TestIfStatements(t *testing.T) {
	i := InterpretTestFile(t, "if_statements.zen")
	if i == nil {
		t.Fatal("Failed to interpret test file")
	}

	// Test basic if statement
	AssertValue(t, i, "result", "greater")

	// Test if-else statement
	AssertValue(t, i, "elseResult", "lesser")

	// Test if-elif-else chain
	AssertValue(t, i, "chainResult", "medium")

	// Test nested scopes
	AssertValue(t, i, "outer", "inner")
	AssertUndefined(t, i, "inner") // inner should not be accessible outside if block

	// Test multiple conditions
	AssertValue(t, i, "multiResult", "both true")

	// Test complex condition
	AssertValue(t, i, "mathResult", "complex true")

	// Test nested if statements
	AssertValue(t, i, "nested", "deep")

	// Test multiple elif conditions
	AssertValue(t, i, "gradeResult", "B")
}

func TestIfStatementErrors(t *testing.T) {
	// Test non-boolean condition
	AssertInterpretError(t, `
		if "string" {
			var x = 1
		}
	`)

	// Test invalid type in condition
	AssertInterpretError(t, `
		var x = {}
		if x {
			var y = 1
		}
	`)

	// Test invalid logical operation
	AssertInterpretError(t, `
		var x = "string"
		var y = 42
		if x and y {
			var z = 1
		}
	`)

	// Test invalid comparison
	AssertInterpretError(t, `
		var x = true
		var y = 42
		if x < y {
			var z = 1
		}
	`)

	// Test missing condition
	AssertInterpretError(t, `
		if {
			var x = 1
		}
	`)

	// Test missing block
	AssertInterpretError(t, `
		if true
			var x = 1
	`)
}

func TestIfStatementScoping(t *testing.T) {
	// Test variable shadowing
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

	// Test nested scopes
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

	// Test elif and else scoping
	i, err = InterpretString(`
		var x = "start"
		if false {
			var a = "a"
			x = a
		} elif true {
			var b = "b"
			x = b
		} else {
			var c = "c"
			x = c
		}
	`)
	if err != nil {
		t.Fatalf("Failed to interpret code: %v", err)
	}

	AssertValue(t, i, "x", "b")
	AssertUndefined(t, i, "a")
	AssertUndefined(t, i, "b")
	AssertUndefined(t, i, "c")
}
