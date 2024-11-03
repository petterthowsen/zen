package parsing

import (
	"testing"
	"zen/lang/parsing/ast"
	"zen/lang/parsing/expression"
	"zen/lang/parsing/statement"
)

func TestWhileStatements(t *testing.T) {
	program := ParseTestFile(t, "while_statements.zen")
	if program == nil {
		return
	}

	// We expect 5 statements (5 while loops)
	if len(program.Statements) != 5 {
		t.Errorf("Expected 5 statements, got %d", len(program.Statements))
		return
	}

	// Test basic while loop
	AssertWhileLoop(t, program.Statements[0],
		func(t *testing.T, expr ast.Expression) {
			AssertBinaryComparison(t, expr, "x", "<", int64(10))
		},
		2, // number of body statements
	)

	// Test while loop with break
	AssertWhileLoopWithBreak(t, program.Statements[1])

	// Test while loop with continue
	AssertWhileLoopWithContinue(t, program.Statements[2])

	// Test while loop with complex condition
	AssertWhileLoop(t, program.Statements[3],
		func(t *testing.T, expr ast.Expression) {
			// Check that it's a logical AND expression
			binary, ok := expr.(*expression.BinaryExpression)
			if !ok {
				t.Errorf("Expected BinaryExpression, got %T", expr)
				return
			}
			if binary.Operator != "and" {
				t.Errorf("Expected 'and' operator, got '%s'", binary.Operator)
			}
		},
		4, // number of body statements
	)

	// Test nested while loops
	AssertNestedWhileLoops(t, program.Statements[4])
}

// AssertWhileLoop verifies a while loop's structure
func AssertWhileLoop(t *testing.T, stmt ast.Statement,
	conditionChecker func(*testing.T, ast.Expression),
	bodyLen int) {

	whileStmt, ok := stmt.(*statement.WhileStatement)
	if !ok {
		t.Errorf("Expected WhileStatement, got %T", stmt)
		return
	}

	// Check condition
	if whileStmt.Condition != nil {
		conditionChecker(t, whileStmt.Condition)
	}

	// Check body length
	if len(whileStmt.Body) != bodyLen {
		t.Errorf("Expected %d body statements, got %d", bodyLen, len(whileStmt.Body))
	}
}

// AssertWhileLoopWithBreak verifies a while loop containing a break statement
func AssertWhileLoopWithBreak(t *testing.T, stmt ast.Statement) {
	whileStmt, ok := stmt.(*statement.WhileStatement)
	if !ok {
		t.Errorf("Expected WhileStatement, got %T", stmt)
		return
	}

	// Verify condition is 'true'
	literal, ok := whileStmt.Condition.(*expression.LiteralExpression)
	if !ok {
		t.Errorf("Expected LiteralExpression, got %T", whileStmt.Condition)
		return
	}
	if literal.Value != true {
		t.Errorf("Expected true literal, got %v", literal.Value)
	}

	// Verify body contains if statement with break
	if len(whileStmt.Body) != 3 { // if statement and two other statements
		t.Errorf("Expected 3 body statements, got %d", len(whileStmt.Body))
		return
	}

	ifStmt := whileStmt.Body[0].(*statement.IfStatement)
	if len(ifStmt.PrimaryBlock) != 1 {
		t.Errorf("Expected 1 statement in if block, got %d", len(ifStmt.PrimaryBlock))
		return
	}

	_, ok = ifStmt.PrimaryBlock[0].(*statement.BreakStatement)
	if !ok {
		t.Errorf("Expected BreakStatement, got %T", ifStmt.PrimaryBlock[0])
	}
}

// AssertWhileLoopWithContinue verifies a while loop containing a continue statement
func AssertWhileLoopWithContinue(t *testing.T, stmt ast.Statement) {
	whileStmt, ok := stmt.(*statement.WhileStatement)
	if !ok {
		t.Errorf("Expected WhileStatement, got %T", stmt)
		return
	}

	// Verify body contains if statement with continue
	if len(whileStmt.Body) != 3 { // assignment, if statement, and print statement
		t.Errorf("Expected 3 body statements, got %d", len(whileStmt.Body))
		return
	}

	ifStmt := whileStmt.Body[1].(*statement.IfStatement)
	if len(ifStmt.PrimaryBlock) != 1 {
		t.Errorf("Expected 1 statement in if block, got %d", len(ifStmt.PrimaryBlock))
		return
	}

	_, ok = ifStmt.PrimaryBlock[0].(*statement.ContinueStatement)
	if !ok {
		t.Errorf("Expected ContinueStatement, got %T", ifStmt.PrimaryBlock[0])
	}
}

// AssertNestedWhileLoops verifies nested while loops
func AssertNestedWhileLoops(t *testing.T, stmt ast.Statement) {
	whileStmt, ok := stmt.(*statement.WhileStatement)
	if !ok {
		t.Errorf("Expected WhileStatement, got %T", stmt)
		return
	}

	// Verify outer loop body contains assignment, inner loop, and increment
	if len(whileStmt.Body) != 3 {
		t.Errorf("Expected 3 body statements, got %d", len(whileStmt.Body))
		return
	}

	// Check inner while loop
	innerLoop, ok := whileStmt.Body[1].(*statement.WhileStatement)
	if !ok {
		t.Errorf("Expected WhileStatement (inner loop), got %T", whileStmt.Body[1])
		return
	}

	// Verify inner loop body contains if statement with break and other statements
	if len(innerLoop.Body) != 4 { // if with break, two prints, increment
		t.Errorf("Expected 4 body statements, got %d", len(innerLoop.Body))
		return
	}

	// Check if statement contains break
	ifStmt := innerLoop.Body[0].(*statement.IfStatement)
	_, ok = ifStmt.PrimaryBlock[0].(*statement.BreakStatement)
	if !ok {
		t.Errorf("Expected BreakStatement, got %T", ifStmt.PrimaryBlock[0])
	}
}
