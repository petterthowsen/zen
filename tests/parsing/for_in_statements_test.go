package parsing

import (
	"testing"
	"zen/lang/parsing/ast"
	"zen/lang/parsing/expression"
	"zen/lang/parsing/statement"
)

func TestForInStatements(t *testing.T) {
	program := ParseTestFile(t, "for_in_statements.zen")
	if program == nil {
		return
	}

	// We expect 4 statements (4 for-in loops)
	if len(program.Statements) != 4 {
		t.Errorf("Expected 4 statements, got %d", len(program.Statements))
		return
	}

	// Test basic for-in loop: for value in items { print(value) }
	AssertForInLoop(t, program.Statements[0],
		"",      // no key
		"value", // value variable
		"items", // container
		1,       // number of body statements
	)

	// Test for-in loop with key and value: for key, value in items { print(key); print(value) }
	AssertForInLoop(t, program.Statements[1],
		"key",   // key variable
		"value", // value variable
		"items", // container
		2,       // number of body statements
	)

	// Test for-in loop with function call: for item in getItems() { print(item) }
	AssertForInLoop(t, program.Statements[2],
		"",     // no key
		"item", // value variable
		"",     // container (will check separately)
		1,      // number of body statements
	)

	// Verify the function call container
	forIn := program.Statements[2].(*statement.ForInStatement)
	call, ok := forIn.Container.(*expression.CallExpression)
	if !ok {
		t.Errorf("Expected CallExpression for container, got %T", forIn.Container)
	} else {
		callee := call.Callee.(*expression.IdentifierExpression)
		if callee.Name != "getItems" {
			t.Errorf("Expected function name 'getItems', got '%s'", callee.Name)
		}
		if len(call.Arguments) != 0 {
			t.Errorf("Expected 0 arguments, got %d", len(call.Arguments))
		}
	}

	// Test for-in loop with multiple statements: for k, v in users { ... }
	AssertForInLoop(t, program.Statements[3],
		"k",     // key variable
		"v",     // value variable
		"users", // container
		3,       // number of body statements
	)
}

// AssertForInLoop verifies the structure of a for-in loop statement
func AssertForInLoop(t *testing.T, stmt ast.Statement, key, value, containerName string, bodyLen int) {
	forInStmt, ok := stmt.(*statement.ForInStatement)
	if !ok {
		t.Errorf("Expected ForInStatement, got %T", stmt)
		return
	}

	// Check key (if expected)
	if key != "" && forInStmt.Key != key {
		t.Errorf("Expected key '%s', got '%s'", key, forInStmt.Key)
	}

	// Check value
	if forInStmt.Value != value {
		t.Errorf("Expected value '%s', got '%s'", value, forInStmt.Value)
	}

	// Check container (if name provided)
	if containerName != "" {
		container, ok := forInStmt.Container.(*expression.IdentifierExpression)
		if !ok {
			t.Errorf("Expected IdentifierExpression for container, got %T", forInStmt.Container)
		} else if container.Name != containerName {
			t.Errorf("Expected container name '%s', got '%s'", containerName, container.Name)
		}
	}

	// Check body length
	if len(forInStmt.Body) != bodyLen {
		t.Errorf("Expected %d body statements, got %d", bodyLen, len(forInStmt.Body))
	}
}
