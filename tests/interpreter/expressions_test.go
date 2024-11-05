package interpreter

import (
	"testing"
)

func TestArithmeticExpressions(t *testing.T) {
	i := InterpretTestFile(t, "expressions.zen")
	if i == nil {
		t.Fatal("Failed to interpret test file")
	}

	// Test basic arithmetic
	AssertValue(t, i, "add", 8)
	AssertValue(t, i, "sub", 6)
	AssertValue(t, i, "mul", 42)
	AssertValue(t, i, "div", 4)
	AssertValue(t, i, "complex", 20) // (2 + 3) * 4

	// Test string operations
	AssertValue(t, i, "greeting", "Hello World")

	// Test comparisons
	AssertValue(t, i, "eq", true)
	AssertValue(t, i, "neq", true)
	AssertValue(t, i, "lt", true)
	AssertValue(t, i, "lte", true)
	AssertValue(t, i, "gt", true)
	AssertValue(t, i, "gte", true)

	// Test logical operations
	AssertValue(t, i, "and1", true)
	AssertValue(t, i, "and2", false)
	AssertValue(t, i, "or1", true)
	AssertValue(t, i, "or2", false)
	AssertValue(t, i, "not1", false)
	AssertValue(t, i, "not2", true)

	// Test complex expressions
	AssertValue(t, i, "complex1", 17) // 5 * 3 + 2
	AssertValue(t, i, "complex2", 25) // 5 * (3 + 2)
	AssertValue(t, i, "complex3", 16) // (5 + 3) * (5 - 3)

	// Test string comparisons
	AssertValue(t, i, "strEq", true)
	AssertValue(t, i, "strNeq", true)

	// Test logical expressions with comparisons
	AssertValue(t, i, "logicalComplex", true)

	// Test operator precedence
	AssertValue(t, i, "precedence1", 14)    // 2 + (3 * 4)
	AssertValue(t, i, "precedence2", 20)    // (2 + 3) * 4
	AssertValue(t, i, "precedence3", false) // not true and false
	AssertValue(t, i, "precedence4", true)  // not (true and false)

	// Test unary operations
	AssertValue(t, i, "unary1", -5)
	AssertValue(t, i, "unary2", 5)
	AssertValue(t, i, "unary3", true)

	// Test chained comparisons
	AssertValue(t, i, "chain1", true)
	AssertValue(t, i, "chain2", true)
}

func TestExpressionErrors(t *testing.T) {
	// Test invalid numeric operations
	AssertInterpretError(t, `
		var x = "hello" + 42
	`)

	// Test invalid logical operations
	AssertInterpretError(t, `
		var x = "hello" or false
	`)

	// Test invalid type operations
	AssertInterpretError(t, `
		var x = true * 3
	`)

	// Test division by zero
	AssertInterpretError(t, `
		var x = 5 / 0
	`)

	// Test invalid comparisons
	AssertInterpretError(t, `
		var x = "hello" < 42
	`)

	// Test invalid logical operations
	AssertInterpretError(t, `
		var x = true and 42
	`)
}

func TestComplexExpressions(t *testing.T) {
	// Test nested arithmetic
	i, err := InterpretString(`
		var x = 2 + 3 * 4 - 6 / 2
		var y = (2 + 3) * (4 - 6 / 2)
	`)
	if err != nil {
		t.Fatalf("Failed to interpret code: %v", err)
	}

	AssertValue(t, i, "x", 11) // 2 + 12 - 3
	AssertValue(t, i, "y", 5)  // 5 * 1

	// Test nested logical expressions
	i, err = InterpretString(`
		var x = not (true and false) or (true and not false)
		var y = not true and not false or true and not false
	`)
	if err != nil {
		t.Fatalf("Failed to interpret code: %v", err)
	}

	AssertValue(t, i, "x", true)
	AssertValue(t, i, "y", true)

	// Test mixed arithmetic and logical
	i, err = InterpretString(`
		var x = (5 + 3 > 7) and (10 * 2 == 20)
		var y = not (5 - 3 <= 1) or (15 / 3 == 5)
	`)
	if err != nil {
		t.Fatalf("Failed to interpret code: %v", err)
	}

	AssertValue(t, i, "x", true)
	AssertValue(t, i, "y", true)
}
