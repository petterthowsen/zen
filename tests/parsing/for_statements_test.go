package parsing

import (
	"testing"
	"zen/lang/parsing/ast"
	"zen/lang/parsing/expression"
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
		AssertInit("i", 0),            // i = 0
		AssertCondition("i", "<", 10), // i < 10
		AssertIncrement("i"),          // i++
		1,                             // number of body statements
	)

	// Test for loop with expression initialization: for x = 0; x < 5; x = x + 2 { print(x) }
	AssertForLoop(t, program.Statements[1],
		AssertInit("x", 0),           // x = 0
		AssertCondition("x", "<", 5), // x < 5
		AssertAdd("x", 2),            // x = x + 2
		1,                            // number of body statements
	)

	// Test for loop with multiple body statements
	AssertForLoop(t, program.Statements[2],
		AssertInit("i", 0),           // i = 0
		AssertCondition("i", "<", 3), // i < 3
		AssertIncrement("i"),         // i++
		2,                            // number of body statements
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

// AssertInit checks initialization statement (e.g., i = 0)
func AssertInit(name string, value int64) func(*testing.T, ast.Statement) {
	return func(t *testing.T, stmt ast.Statement) {
		exprStmt, ok := stmt.(*statement.ExpressionStatement)
		if !ok {
			t.Errorf("Expected ExpressionStatement, got %T", stmt)
			return
		}

		binary, ok := exprStmt.Expression.(*expression.BinaryExpression)
		if !ok {
			t.Errorf("Expected BinaryExpression, got %T", exprStmt.Expression)
			return
		}

		if binary.Operator != "=" {
			t.Errorf("Expected assignment operator '=', got '%s'", binary.Operator)
		}

		ident, ok := binary.Left.(*expression.IdentifierExpression)
		if !ok {
			t.Errorf("Expected IdentifierExpression, got %T", binary.Left)
			return
		}

		if ident.Name != name {
			t.Errorf("Expected identifier '%s', got '%s'", name, ident.Name)
		}

		literal, ok := binary.Right.(*expression.LiteralExpression)
		if !ok {
			t.Errorf("Expected LiteralExpression, got %T", binary.Right)
			return
		}

		if literal.Value != value {
			t.Errorf("Expected value %d, got %v", value, literal.Value)
		}
	}
}

// AssertCondition checks condition expression (e.g., i < 10)
func AssertCondition(name, op string, value int64) func(*testing.T, ast.Expression) {
	return func(t *testing.T, expr ast.Expression) {
		binary, ok := expr.(*expression.BinaryExpression)
		if !ok {
			t.Errorf("Expected BinaryExpression, got %T", expr)
			return
		}

		if binary.Operator != op {
			t.Errorf("Expected operator '%s', got '%s'", op, binary.Operator)
		}

		ident, ok := binary.Left.(*expression.IdentifierExpression)
		if !ok {
			t.Errorf("Expected IdentifierExpression, got %T", binary.Left)
			return
		}

		if ident.Name != name {
			t.Errorf("Expected identifier '%s', got '%s'", name, ident.Name)
		}

		literal, ok := binary.Right.(*expression.LiteralExpression)
		if !ok {
			t.Errorf("Expected LiteralExpression, got %T", binary.Right)
			return
		}

		if literal.Value != value {
			t.Errorf("Expected value %d, got %v", value, literal.Value)
		}
	}
}

// AssertIncrement checks increment expression (i++)
func AssertIncrement(name string) func(*testing.T, ast.Statement) {
	return func(t *testing.T, stmt ast.Statement) {
		exprStmt, ok := stmt.(*statement.ExpressionStatement)
		if !ok {
			t.Errorf("Expected ExpressionStatement, got %T", stmt)
			return
		}

		assign, ok := exprStmt.Expression.(*expression.BinaryExpression)
		if !ok {
			t.Errorf("Expected BinaryExpression, got %T", exprStmt.Expression)
			return
		}

		if assign.Operator != "=" {
			t.Errorf("Expected assignment operator '=', got '%s'", assign.Operator)
		}

		ident, ok := assign.Left.(*expression.IdentifierExpression)
		if !ok {
			t.Errorf("Expected IdentifierExpression, got %T", assign.Left)
			return
		}

		if ident.Name != name {
			t.Errorf("Expected identifier '%s', got '%s'", name, ident.Name)
		}

		add, ok := assign.Right.(*expression.BinaryExpression)
		if !ok {
			t.Errorf("Expected BinaryExpression, got %T", assign.Right)
			return
		}

		if add.Operator != "+" {
			t.Errorf("Expected addition operator '+', got '%s'", add.Operator)
		}

		identRight, ok := add.Left.(*expression.IdentifierExpression)
		if !ok {
			t.Errorf("Expected IdentifierExpression, got %T", add.Left)
			return
		}

		if identRight.Name != name {
			t.Errorf("Expected identifier '%s', got '%s'", name, identRight.Name)
		}

		literal, ok := add.Right.(*expression.LiteralExpression)
		if !ok {
			t.Errorf("Expected LiteralExpression, got %T", add.Right)
			return
		}

		if literal.Value != int64(1) {
			t.Errorf("Expected value 1, got %v", literal.Value)
		}
	}
}

// AssertAdd checks addition assignment (x = x + 2)
func AssertAdd(name string, value int64) func(*testing.T, ast.Statement) {
	return func(t *testing.T, stmt ast.Statement) {
		exprStmt, ok := stmt.(*statement.ExpressionStatement)
		if !ok {
			t.Errorf("Expected ExpressionStatement, got %T", stmt)
			return
		}

		assign, ok := exprStmt.Expression.(*expression.BinaryExpression)
		if !ok {
			t.Errorf("Expected BinaryExpression, got %T", exprStmt.Expression)
			return
		}

		if assign.Operator != "=" {
			t.Errorf("Expected assignment operator '=', got '%s'", assign.Operator)
		}

		ident, ok := assign.Left.(*expression.IdentifierExpression)
		if !ok {
			t.Errorf("Expected IdentifierExpression, got %T", assign.Left)
			return
		}

		if ident.Name != name {
			t.Errorf("Expected identifier '%s', got '%s'", name, ident.Name)
		}

		add, ok := assign.Right.(*expression.BinaryExpression)
		if !ok {
			t.Errorf("Expected BinaryExpression, got %T", assign.Right)
			return
		}

		if add.Operator != "+" {
			t.Errorf("Expected addition operator '+', got '%s'", add.Operator)
		}

		identRight, ok := add.Left.(*expression.IdentifierExpression)
		if !ok {
			t.Errorf("Expected IdentifierExpression, got %T", add.Left)
			return
		}

		if identRight.Name != name {
			t.Errorf("Expected identifier '%s', got '%s'", name, identRight.Name)
		}

		literal, ok := add.Right.(*expression.LiteralExpression)
		if !ok {
			t.Errorf("Expected LiteralExpression, got %T", add.Right)
			return
		}

		if literal.Value != value {
			t.Errorf("Expected value %d, got %v", value, literal.Value)
		}
	}
}
