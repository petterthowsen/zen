package parsing

import (
	"testing"
	"zen/lang/parsing/expression"
)

func TestArrayLiterals(t *testing.T) {
	print("parsing arrays.zen\n")
	programNode := ParseTestFile(t, "arrays.zen")
	if programNode == nil {
		return
	}

	// var numbers = [42, 1337, 7]
	varDecl := AssertVarDeclaration(t, programNode.Statements[0], "numbers", false, false)
	arrayLit, ok := varDecl.Initializer.(*expression.ArrayLiteralExpression)
	if !ok {
		t.Errorf("Expected ArrayLiteralExpression, got %T", varDecl.Initializer)
		return
	}

	if len(arrayLit.Elements) != 3 {
		t.Errorf("Expected 3 elements, got %d", len(arrayLit.Elements))
		return
	}

	AssertLiteralExpression(t, arrayLit.Elements[0], int64(42))
	AssertLiteralExpression(t, arrayLit.Elements[1], int64(1337))
	AssertLiteralExpression(t, arrayLit.Elements[2], int64(7))

	// var floats : Array<float, 2> = [3.14, 1.62]
	varDecl = AssertVarDeclaration(t, programNode.Statements[1], "floats", false, false)

	// Check array literal
	arrayLit, ok = varDecl.Initializer.(*expression.ArrayLiteralExpression)
	if !ok {
		t.Errorf("Expected ArrayLiteralExpression, got %T", varDecl.Initializer)
		return
	}

	if len(arrayLit.Elements) != 2 {
		t.Errorf("Expected 2 elements, got %d", len(arrayLit.Elements))
		return
	}

	AssertLiteralExpression(t, arrayLit.Elements[0], float64(3.14))
	AssertLiteralExpression(t, arrayLit.Elements[1], float64(1.62))
}
