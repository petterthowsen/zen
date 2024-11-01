package parsing

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"zen/lang/common"
	"zen/lang/lexing"
	"zen/lang/parsing"
	"zen/lang/parsing/ast"
	"zen/lang/parsing/expression"
	"zen/lang/parsing/statement"
)

// ParseString parses source code string and returns the AST
func ParseString(source string) (*ast.ProgramNode, []error) {
	sourceCode := common.NewInlineSourceCode(source)
	lexer := lexing.NewLexer(sourceCode)
	tokens, err := lexer.Scan()
	if err != nil {
		return nil, []error{err}
	}

	parser := parsing.NewParser(tokens, false)
	program, syntaxErrors := parser.Parse()

	// Convert syntax errors to regular errors
	errors := make([]error, len(syntaxErrors))
	for i, err := range syntaxErrors {
		errors[i] = err
	}

	return program, errors
}

// ParseFile parses a source file and returns the AST
func ParseFile(path string, content string) (*ast.ProgramNode, []error) {
	sourceCode := common.NewFileSourceCode(path, content)
	lexer := lexing.NewLexer(sourceCode)
	tokens, err := lexer.Scan()
	if err != nil {
		return nil, []error{err}
	}

	parser := parsing.NewParser(tokens, false)
	program, syntaxErrors := parser.Parse()

	// Convert syntax errors to regular errors
	errors := make([]error, len(syntaxErrors))
	for i, err := range syntaxErrors {
		errors[i] = err
	}

	return program, errors
}

// getTestDataPath returns the absolute path to the test data file
func getTestDataPath(filename string) string {
	_, currentFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(currentFile)
	return filepath.Join(dir, filename)
}

// ParseTestFile is a helper that loads, parses and saves AST for a test file
func ParseTestFile(t *testing.T, filename string) *ast.ProgramNode {
	// Load and parse the test file
	path := getTestDataPath(filename)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", path, err)
	}

	program, errors := ParseFile(path, string(content))
	if len(errors) > 0 {
		t.Errorf("Parser errors:")
		for _, err := range errors {
			t.Errorf("  %s", err.Error())
		}
		return nil
	}

	if program == nil {
		t.Error("Expected program node, got nil")
		return nil
	}

	// Save AST to file
	astStr := program.String(0)
	astPath := strings.TrimSuffix(path, ".zen") + ".ast"
	if err := os.WriteFile(astPath, []byte(astStr), 0644); err != nil {
		t.Logf("Warning: Failed to write AST file: %v", err)
	}

	return program
}

// AssertVarDeclaration checks if a statement is a variable declaration with expected properties
func AssertVarDeclaration(t *testing.T, stmt ast.Statement, name string, typ string, isConst bool, nullable bool) *statement.VarDeclarationNode {
	varDecl, ok := stmt.(*statement.VarDeclarationNode)
	if !ok {
		t.Errorf("Expected VarDeclarationNode, got %T", stmt)
		return nil
	}

	if varDecl.Name != name {
		t.Errorf("Expected name %q, got %q", name, varDecl.Name)
	}

	if varDecl.Type != typ {
		t.Errorf("Expected type %q, got %q", typ, varDecl.Type)
	}

	if varDecl.IsConstant != isConst {
		t.Errorf("Expected const=%v, got %v", isConst, varDecl.IsConstant)
	}

	if varDecl.IsNullable != nullable {
		t.Errorf("Expected nullable=%v, got %v", nullable, varDecl.IsNullable)
	}

	return varDecl
}

// AssertIfStatement checks if a statement is an if statement and returns it
func AssertIfStatement(t *testing.T, stmt ast.Statement) *statement.IfStatement {
	ifStmt, ok := stmt.(*statement.IfStatement)
	if !ok {
		t.Errorf("Expected IfStatement, got %T", stmt)
		return nil
	}
	return ifStmt
}

// AssertLiteralExpression checks if an expression is a literal with expected value
func AssertLiteralExpression(t *testing.T, expr ast.Expression, expectedValue interface{}) *expression.LiteralExpression {
	lit, ok := expr.(*expression.LiteralExpression)
	if !ok {
		t.Errorf("Expected LiteralExpression, got %T", expr)
		return nil
	}
	if lit.Value != expectedValue {
		t.Errorf("Expected literal value %v, got %v", expectedValue, lit.Value)
	}
	return lit
}

// AssertBinaryExpression checks if an expression is a binary expression with expected operator
func AssertBinaryExpression(t *testing.T, expr ast.Expression, expectedOp string) *expression.BinaryExpression {
	bin, ok := expr.(*expression.BinaryExpression)
	if !ok {
		t.Errorf("Expected BinaryExpression, got %T", expr)
		return nil
	}
	if bin.Operator != expectedOp {
		t.Errorf("Expected operator %s, got %s", expectedOp, bin.Operator)
	}
	return bin
}

// AssertUnaryExpression checks if an expression is a unary expression with expected operator
func AssertUnaryExpression(t *testing.T, expr ast.Expression, expectedOp string) *expression.UnaryExpression {
	un, ok := expr.(*expression.UnaryExpression)
	if !ok {
		t.Errorf("Expected UnaryExpression, got %T", expr)
		return nil
	}
	if un.Operator != expectedOp {
		t.Errorf("Expected operator %s, got %s", expectedOp, un.Operator)
	}
	return un
}

// AssertIdentifierExpression checks if an expression is an identifier with expected name
func AssertIdentifierExpression(t *testing.T, expr ast.Expression, expectedName string) *expression.IdentifierExpression {
	id, ok := expr.(*expression.IdentifierExpression)
	if !ok {
		t.Errorf("Expected IdentifierExpression, got %T", expr)
		return nil
	}
	if id.Name != expectedName {
		t.Errorf("Expected identifier %s, got %s", expectedName, id.Name)
	}
	return id
}

// AssertCallExpression checks if an expression is a function call with expected number of arguments
func AssertCallExpression(t *testing.T, expr ast.Expression, expectedArgCount int) *expression.CallExpression {
	call, ok := expr.(*expression.CallExpression)
	if !ok {
		t.Errorf("Expected CallExpression, got %T", expr)
		return nil
	}
	if len(call.Arguments) != expectedArgCount {
		t.Errorf("Expected %d arguments, got %d", expectedArgCount, len(call.Arguments))
	}
	return call
}

// AssertParseError checks if parsing a string produces an error
func AssertParseError(t *testing.T, input string) {
	program, errors := ParseString(input)
	if len(errors) == 0 {
		t.Errorf("Input %q: expected error, got none", input)
	}
	if program != nil {
		t.Errorf("Input %q: expected nil program for error case", input)
	}
}
