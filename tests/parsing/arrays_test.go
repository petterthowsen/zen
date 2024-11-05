package parsing

import (
	"testing"
	"zen/lang/parsing/expression"
	"zen/lang/parsing/statement"
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

	// print(numbers[1])
	exprStmt, ok := programNode.Statements[2].(*statement.ExpressionStatement)
	if !ok {
		t.Errorf("Expected ExpressionStatement, got %T", programNode.Statements[2])
		return
	}
	call := AssertCallExpression(t, exprStmt.Expression, 1)
	AssertIdentifierExpression(t, call.Callee, "print")
	arrayAccess, ok := call.Arguments[0].(*expression.ArrayAccessExpression)
	if !ok {
		t.Errorf("Expected ArrayAccessExpression, got %T", call.Arguments[0])
		return
	}
	AssertIdentifierExpression(t, arrayAccess.Array, "numbers")
	AssertLiteralExpression(t, arrayAccess.Index, int64(1))

	// const PI = floats[0]
	varDecl = AssertVarDeclaration(t, programNode.Statements[3], "PI", true, false)
	arrayAccess, ok = varDecl.Initializer.(*expression.ArrayAccessExpression)
	if !ok {
		t.Errorf("Expected ArrayAccessExpression, got %T", varDecl.Initializer)
		return
	}
	AssertIdentifierExpression(t, arrayAccess.Array, "floats")
	AssertLiteralExpression(t, arrayAccess.Index, int64(0))
}
