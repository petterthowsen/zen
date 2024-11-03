package parsing

import (
	"testing"
	"zen/lang/parsing/ast"
	"zen/lang/parsing/statement"
)

func TestBreakContinueStatements(t *testing.T) {
	program := ParseTestFile(t, "break_continue_statements.zen")
	if program == nil {
		return
	}

	// We expect 5 statements (5 loops with break/continue)
	if len(program.Statements) != 5 {
		t.Errorf("Expected 5 statements, got %d", len(program.Statements))
		return
	}

	// Test break in traditional for loop
	AssertForLoopWithBreak(t, program.Statements[0])

	// Test continue in traditional for loop
	AssertForLoopWithContinue(t, program.Statements[1])

	// Test break in for-in loop
	AssertForInLoopWithBreak(t, program.Statements[2])

	// Test continue in for-in loop
	AssertForInLoopWithContinue(t, program.Statements[3])

	// Test nested loops with break and continue
	AssertNestedLoopsWithBreakContinue(t, program.Statements[4])
}

// AssertForLoopWithBreak verifies a for loop containing a break statement
func AssertForLoopWithBreak(t *testing.T, stmt ast.Statement) {
	forStmt, ok := stmt.(*statement.ForStatement)
	if !ok {
		t.Errorf("Expected ForStatement, got %T", stmt)
		return
	}

	// Verify loop body contains if statement with break
	if len(forStmt.Body) != 2 { // if statement and print statement
		t.Errorf("Expected 2 body statements, got %d", len(forStmt.Body))
		return
	}

	ifStmt := forStmt.Body[0].(*statement.IfStatement)
	if len(ifStmt.PrimaryBlock) != 1 {
		t.Errorf("Expected 1 statement in if block, got %d", len(ifStmt.PrimaryBlock))
		return
	}

	_, ok = ifStmt.PrimaryBlock[0].(*statement.BreakStatement)
	if !ok {
		t.Errorf("Expected BreakStatement, got %T", ifStmt.PrimaryBlock[0])
	}
}

// AssertForLoopWithContinue verifies a for loop containing a continue statement
func AssertForLoopWithContinue(t *testing.T, stmt ast.Statement) {
	forStmt, ok := stmt.(*statement.ForStatement)
	if !ok {
		t.Errorf("Expected ForStatement, got %T", stmt)
		return
	}

	// Verify loop body contains if statement with continue
	if len(forStmt.Body) != 2 { // if statement and print statement
		t.Errorf("Expected 2 body statements, got %d", len(forStmt.Body))
		return
	}

	ifStmt := forStmt.Body[0].(*statement.IfStatement)
	if len(ifStmt.PrimaryBlock) != 1 {
		t.Errorf("Expected 1 statement in if block, got %d", len(ifStmt.PrimaryBlock))
		return
	}

	_, ok = ifStmt.PrimaryBlock[0].(*statement.ContinueStatement)
	if !ok {
		t.Errorf("Expected ContinueStatement, got %T", ifStmt.PrimaryBlock[0])
	}
}

// AssertForInLoopWithBreak verifies a for-in loop containing a break statement
func AssertForInLoopWithBreak(t *testing.T, stmt ast.Statement) {
	forInStmt, ok := stmt.(*statement.ForInStatement)
	if !ok {
		t.Errorf("Expected ForInStatement, got %T", stmt)
		return
	}

	// Verify loop body contains if statement with break
	if len(forInStmt.Body) != 2 { // if statement and print statement
		t.Errorf("Expected 2 body statements, got %d", len(forInStmt.Body))
		return
	}

	ifStmt := forInStmt.Body[0].(*statement.IfStatement)
	if len(ifStmt.PrimaryBlock) != 1 {
		t.Errorf("Expected 1 statement in if block, got %d", len(ifStmt.PrimaryBlock))
		return
	}

	_, ok = ifStmt.PrimaryBlock[0].(*statement.BreakStatement)
	if !ok {
		t.Errorf("Expected BreakStatement, got %T", ifStmt.PrimaryBlock[0])
	}
}

// AssertForInLoopWithContinue verifies a for-in loop containing a continue statement
func AssertForInLoopWithContinue(t *testing.T, stmt ast.Statement) {
	forInStmt, ok := stmt.(*statement.ForInStatement)
	if !ok {
		t.Errorf("Expected ForInStatement, got %T", stmt)
		return
	}

	// Verify loop body contains if statement with continue
	if len(forInStmt.Body) != 2 { // if statement and print statement
		t.Errorf("Expected 2 body statements, got %d", len(forInStmt.Body))
		return
	}

	ifStmt := forInStmt.Body[0].(*statement.IfStatement)
	if len(ifStmt.PrimaryBlock) != 1 {
		t.Errorf("Expected 1 statement in if block, got %d", len(ifStmt.PrimaryBlock))
		return
	}

	_, ok = ifStmt.PrimaryBlock[0].(*statement.ContinueStatement)
	if !ok {
		t.Errorf("Expected ContinueStatement, got %T", ifStmt.PrimaryBlock[0])
	}
}

// AssertNestedLoopsWithBreakContinue verifies nested loops with break and continue statements
func AssertNestedLoopsWithBreakContinue(t *testing.T, stmt ast.Statement) {
	forStmt, ok := stmt.(*statement.ForStatement)
	if !ok {
		t.Errorf("Expected ForStatement, got %T", stmt)
		return
	}

	// Verify outer loop body contains inner loop
	if len(forStmt.Body) != 1 {
		t.Errorf("Expected 1 body statement (inner loop), got %d", len(forStmt.Body))
		return
	}

	innerLoop, ok := forStmt.Body[0].(*statement.ForStatement)
	if !ok {
		t.Errorf("Expected ForStatement (inner loop), got %T", forStmt.Body[0])
		return
	}

	// Verify inner loop body contains two if statements and two print statements
	if len(innerLoop.Body) != 4 {
		t.Errorf("Expected 4 body statements, got %d", len(innerLoop.Body))
		return
	}

	// Check first if statement contains break
	firstIf := innerLoop.Body[0].(*statement.IfStatement)
	_, ok = firstIf.PrimaryBlock[0].(*statement.BreakStatement)
	if !ok {
		t.Errorf("Expected BreakStatement, got %T", firstIf.PrimaryBlock[0])
	}

	// Check second if statement contains continue
	secondIf := innerLoop.Body[1].(*statement.IfStatement)
	_, ok = secondIf.PrimaryBlock[0].(*statement.ContinueStatement)
	if !ok {
		t.Errorf("Expected ContinueStatement, got %T", secondIf.PrimaryBlock[0])
	}
}
