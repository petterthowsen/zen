package interpreter

import (
	"testing"
)

func TestWhileLoops(t *testing.T) {
	i := InterpretTestFile(t, "while_loops.zen")
	if i == nil {
		t.Fatal("Failed to interpret test file")
	}

	// Test basic counter loop
	AssertValue(t, i, "counter", 3)

	// Test loop with false condition
	AssertValue(t, i, "shouldNotChange", 42)

	// Test nested scope
	AssertValue(t, i, "outer", "done")
	AssertUndefined(t, i, "inner") // inner should not be accessible outside loop

	// Test multiple conditions
	AssertValue(t, i, "x", 2)
	AssertValue(t, i, "y", 3) // y = 0 + 1 + 2
}

func TestWhileLoopErrors(t *testing.T) {
	// Test undefined variable in condition
	AssertInterpretError(t, `
		while undefinedVar < 5 {
			var x = 1
		}
	`)

	// Test non-boolean condition
	AssertInterpretError(t, `
		while 42 {
			var x = 1
		}
	`)

	// Test invalid operation in condition
	AssertInterpretError(t, `
		var x = "string"
		while x < 5 {
			x = x + 1
		}
	`)
}
