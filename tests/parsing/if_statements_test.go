package parsing

import (
	"testing"
	"zen/lang/common"
	"zen/lang/lexing"
	"zen/lang/parsing"
)

func TestIfStatements(t *testing.T) {
	// First test lexing
	t.Log("Testing lexer...")
	content := `if true {
    print("first")
}

var name = "john"
var age = 25

if name == "john" or age > 20 {
    print("second")
}

if (name == "john" and age > 20) or (name == "jane" and age > 40) {
    print("third")
}`

	sourceCode := common.NewInlineSourceCode(content)
	lexer := lexing.NewLexer(sourceCode)
	tokens, err := lexer.Scan()
	if err != nil {
		t.Fatalf("Lexer error: %v", err)
	}

	t.Log("Tokens:")
	for _, token := range tokens {
		t.Logf("%s", token.String())
	}

	// Now test parsing
	t.Log("\nTesting parser...")
	parser := parsing.NewParser(tokens, false)
	program, errors := parser.Parse()
	if len(errors) > 0 {
		t.Error("Parser errors:")
		for _, err := range errors {
			t.Errorf("  %v", err)
		}
		return
	}

	if program == nil {
		t.Fatal("Failed to parse program")
		return
	}

	t.Log("Program parsed successfully")
	t.Logf("Found %d statements", len(program.Statements))

	// Test basic if with boolean literal
	t.Log("Testing first if statement (if true)")
	stmt := program.Statements[0]
	ifStmt := AssertIfStatement(t, stmt)
	if ifStmt == nil {
		t.Fatal("Failed to parse if statement")
		return
	}
	AssertLiteralExpression(t, ifStmt.Condition, true)

	// Test if with comparison
	t.Log("Testing second if statement (if name == 'john' or age > 20)")
	stmt = program.Statements[3]
	ifStmt = AssertIfStatement(t, stmt)
	if ifStmt == nil {
		t.Fatal("Failed to parse if statement")
		return
	}
	binary := AssertBinaryExpression(t, ifStmt.Condition, "or")
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
	t.Log("Testing third if statement (complex logical operations)")
	stmt = program.Statements[4]
	ifStmt = AssertIfStatement(t, stmt)
	if ifStmt == nil {
		t.Fatal("Failed to parse if statement")
		return
	}
	binary = AssertBinaryExpression(t, ifStmt.Condition, "or")
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

	t.Log("All tests passed successfully")
}
