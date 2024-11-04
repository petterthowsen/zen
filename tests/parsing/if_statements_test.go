package parsing

import (
	"testing"
)

func TestIfStatements(t *testing.T) {
	program := ParseTestFile(t, "if_statements.zen")
	if program == nil {
		return
	}

	t.Log("Program parsed successfully")

	// Test basic if with boolean literal
	stmt := program.Statements[0]
	ifStmt := AssertIfStatement(t, stmt)
	if ifStmt == nil {
		t.Fatal("Failed to parse if statement")
		return
	}
	AssertLiteralExpression(t, ifStmt.PrimaryCondition, true)

	// Test if with comparison and 'or'
	stmt = program.Statements[3]
	ifStmt = AssertIfStatement(t, stmt)
	if ifStmt == nil {
		t.Fatal("Failed to parse if statement")
		return
	}
	binary := AssertBinaryExpression(t, ifStmt.PrimaryCondition, "or")
	if binary == nil {
		t.Fatal("Failed to parse binary expression")
		return
	}

	// name == "john"
	left := AssertBinaryExpression(t, binary.Left, "==")
	if left == nil {
		t.Fatal("Failed to parse left comparison")
		return
	}
	AssertIdentifierExpression(t, left.Left, "name")
	AssertLiteralExpression(t, left.Right, "john")

	// age > 20
	right := AssertBinaryExpression(t, binary.Right, ">")
	if right == nil {
		t.Fatal("Failed to parse right comparison")
		return
	}
	AssertIdentifierExpression(t, right.Left, "age")
	AssertLiteralExpression(t, right.Right, int64(20))

	// Test if with complex logical operations
	stmt = program.Statements[4]
	ifStmt = AssertIfStatement(t, stmt)
	if ifStmt == nil {
		t.Fatal("Failed to parse if statement")
		return
	}
	binary = AssertBinaryExpression(t, ifStmt.PrimaryCondition, "or")
	if binary == nil {
		t.Fatal("Failed to parse binary expression")
		return
	}

	// (name == "john" and age > 20)
	left = AssertBinaryExpression(t, binary.Left, "and")
	if left == nil {
		t.Fatal("Failed to parse left and expression")
		return
	}
	nameEq := AssertBinaryExpression(t, left.Left, "==")
	AssertIdentifierExpression(t, nameEq.Left, "name")
	AssertLiteralExpression(t, nameEq.Right, "john")
	ageGt := AssertBinaryExpression(t, left.Right, ">")
	AssertIdentifierExpression(t, ageGt.Left, "age")
	AssertLiteralExpression(t, ageGt.Right, int64(20))

	// (name == "jane" and age > 40)
	right = AssertBinaryExpression(t, binary.Right, "and")
	if right == nil {
		t.Fatal("Failed to parse right and expression")
		return
	}
	nameEq = AssertBinaryExpression(t, right.Left, "==")
	AssertIdentifierExpression(t, nameEq.Left, "name")
	AssertLiteralExpression(t, nameEq.Right, "jane")
	ageGt = AssertBinaryExpression(t, right.Right, ">")
	AssertIdentifierExpression(t, ageGt.Left, "age")
	AssertLiteralExpression(t, ageGt.Right, int64(40))

	// Test if with else
	stmt = program.Statements[5]
	ifStmt = AssertIfStatement(t, stmt)
	if ifStmt == nil {
		t.Fatal("Failed to parse if-else statement")
		return
	}
	if ifStmt.ElseBlock == nil {
		t.Fatal("Expected else block but got nil")
		return
	}
	binary = AssertBinaryExpression(t, ifStmt.PrimaryCondition, "==")
	AssertIdentifierExpression(t, binary.Left, "name")
	AssertLiteralExpression(t, binary.Right, "john")

	// Test if with else if
	stmt = program.Statements[6]
	ifStmt = AssertIfStatement(t, stmt)
	if ifStmt == nil {
		t.Fatal("Failed to parse if-else-if statement")
		return
	}
	if len(ifStmt.ElseIfBlocks) != 1 {
		t.Fatalf("Expected 1 else-if block but got %d", len(ifStmt.ElseIfBlocks))
		return
	}
	binary = AssertBinaryExpression(t, ifStmt.PrimaryCondition, "==")
	AssertIdentifierExpression(t, binary.Left, "name")
	AssertLiteralExpression(t, binary.Right, "john")

	elseIfCond := AssertBinaryExpression(t, ifStmt.ElseIfBlocks[0].Condition, "==")
	AssertIdentifierExpression(t, elseIfCond.Left, "name")
	AssertLiteralExpression(t, elseIfCond.Right, "jane")

	// Test multiple if-else-if with else
	stmt = program.Statements[7]
	ifStmt = AssertIfStatement(t, stmt)
	if ifStmt == nil {
		t.Fatal("Failed to parse multiple if-else-if statement")
		return
	}
	if len(ifStmt.ElseIfBlocks) != 1 {
		t.Fatalf("Expected 1 else-if block but got %d", len(ifStmt.ElseIfBlocks))
		return
	}
	if ifStmt.ElseBlock == nil {
		t.Fatal("Expected else block but got nil")
		return
	}
	binary = AssertBinaryExpression(t, ifStmt.PrimaryCondition, "==")
	AssertIdentifierExpression(t, binary.Left, "name")
	AssertLiteralExpression(t, binary.Right, "john")

	elseIfCond = AssertBinaryExpression(t, ifStmt.ElseIfBlocks[0].Condition, "==")
	AssertIdentifierExpression(t, elseIfCond.Left, "name")
	AssertLiteralExpression(t, elseIfCond.Right, "jane")
}
