package parsing

import (
	"testing"
	"zen/lang/parsing/ast"
	"zen/lang/parsing/statement"
)

func TestForStatements(t *testing.T) {
	program := ParseTestFile(t, "for_statements.zen")
	if program == nil {
		return
	}

	// We expect 3 statements (3 for loops)
	if len(program.Statements) != 3 {
		t.Errorf("Expected 3 statements, got %d", len(program.Statements))
		return
	}

	// Test basic for loop: for i = 0; i < 10; i++ { print(i) }
	AssertForLoop(t, program.Statements[0],
		func(t *testing.T, stmt ast.Statement) {
			AssertBinaryAssignment(t, stmt, "i", "=", int64(0))
		},
		func(t *testing.T, expr ast.Expression) {
			AssertBinaryComparison(t, expr, "i", "<", int64(10))
		},
		func(t *testing.T, stmt ast.Statement) {
			AssertIncrement(t, stmt, "i")
		},
		1, // number of body statements
	)

	// Test for loop with expression initialization: for x = 0; x < 5; x = x + 2 { print(x) }
	AssertForLoop(t, program.Statements[1],
		func(t *testing.T, stmt ast.Statement) {
			AssertBinaryAssignment(t, stmt, "x", "=", int64(0))
		},
		func(t *testing.T, expr ast.Expression) {
			AssertBinaryComparison(t, expr, "x", "<", int64(5))
		},
		func(t *testing.T, stmt ast.Statement) {
			AssertBinaryAssignment(t, stmt, "x", "=", &BinaryCheck{
				LeftName:   "x",
				Operator:   "+",
				RightValue: int64(2),
			})
		},
		1, // number of body statements
	)

	// Test for loop with multiple body statements
	AssertForLoop(t, program.Statements[2],
		func(t *testing.T, stmt ast.Statement) {
			AssertBinaryAssignment(t, stmt, "i", "=", int64(0))
		},
		func(t *testing.T, expr ast.Expression) {
			AssertBinaryComparison(t, expr, "i", "<", int64(3))
		},
		func(t *testing.T, stmt ast.Statement) {
			AssertIncrement(t, stmt, "i")
		},
		2, // number of body statements
	)
}

// AssertForLoop verifies the structure of a for loop statement
func AssertForLoop(t *testing.T, stmt ast.Statement,
	initChecker func(*testing.T, ast.Statement),
	condChecker func(*testing.T, ast.Expression),
	updateChecker func(*testing.T, ast.Statement),
	bodyLen int) {

	forStmt, ok := stmt.(*statement.ForStatement)
	if !ok {
		t.Errorf("Expected ForStatement, got %T", stmt)
		return
	}

	// Check initialization
	if forStmt.Init != nil {
		initChecker(t, forStmt.Init)
	}

	// Check condition
	if forStmt.Condition != nil {
		condChecker(t, forStmt.Condition)
	}

	// Check increment
	if forStmt.Increment != nil {
		updateChecker(t, forStmt.Increment)
	}

	// Check body length
	if len(forStmt.Body) != bodyLen {
		t.Errorf("Expected %d body statements, got %d", bodyLen, len(forStmt.Body))
	}
}
