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

const StopAtFirstError = false

// ParseString parses source code string and returns the AST
func ParseString(source string) (*ast.ProgramNode, []error) {
	sourceCode := common.NewInlineSourceCode(source)
	lexer := lexing.NewLexer(sourceCode)
	tokens, err := lexer.Scan()
	if err != nil {
		return nil, []error{err}
	}

	parser := parsing.NewParser(tokens, StopAtFirstError)
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

	parser := parsing.NewParser(tokens, StopAtFirstError)
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

func AssertFuncDeclaration(t *testing.T, stmt ast.Statement) *statement.FuncDeclaration {
	funcDecl, ok := stmt.(*statement.FuncDeclaration)
	if !ok {
		t.Errorf("Expected FuncDeclaration, got %T", stmt)
		return nil
	}
	return funcDecl
}

// AssertParseError checks if parsing a string produces an error
func AssertParseError(t *testing.T, input string) {
	_, errors := ParseString(input)
	if len(errors) == 0 {
		t.Errorf("Input %q: expected error, got none", input)
	}
}

// BinaryCheck holds information for checking a binary expression
type BinaryCheck struct {
	LeftName   string
	Operator   string
	RightValue interface{}
}

// Check verifies a binary expression matches expected properties
func (bc *BinaryCheck) Check(t *testing.T, expr *expression.BinaryExpression) {
	if expr.Operator != bc.Operator {
		t.Errorf("Expected operator '%s', got '%s'", bc.Operator, expr.Operator)
	}

	ident, ok := expr.Left.(*expression.IdentifierExpression)
	if !ok {
		t.Errorf("Expected IdentifierExpression, got %T", expr.Left)
		return
	}

	if ident.Name != bc.LeftName {
		t.Errorf("Expected identifier '%s', got '%s'", bc.LeftName, ident.Name)
	}

	literal, ok := expr.Right.(*expression.LiteralExpression)
	if !ok {
		t.Errorf("Expected LiteralExpression, got %T", expr.Right)
		return
	}

	if literal.Value != bc.RightValue {
		t.Errorf("Expected value %v, got %v", bc.RightValue, literal.Value)
	}
}

// AssertBinaryAssignment checks if a statement is an assignment with expected properties
func AssertBinaryAssignment(t *testing.T, stmt ast.Statement, varName string, operator string, rightValue interface{}) {
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

	if ident.Name != varName {
		t.Errorf("Expected identifier '%s', got '%s'", varName, ident.Name)
	}

	// Handle different right-hand side patterns
	switch v := rightValue.(type) {
	case int64:
		// Simple assignment: x = 5
		literal, ok := binary.Right.(*expression.LiteralExpression)
		if !ok {
			t.Errorf("Expected LiteralExpression, got %T", binary.Right)
			return
		}
		if literal.Value != v {
			t.Errorf("Expected value %d, got %v", v, literal.Value)
		}
	case *BinaryCheck:
		// Complex assignment: x = x + 2
		rightBinary, ok := binary.Right.(*expression.BinaryExpression)
		if !ok {
			t.Errorf("Expected BinaryExpression, got %T", binary.Right)
			return
		}
		v.Check(t, rightBinary)
	}
}

// AssertBinaryComparison checks if an expression is a comparison with expected properties
func AssertBinaryComparison(t *testing.T, expr ast.Expression, leftName string, operator string, rightValue interface{}) {
	binary, ok := expr.(*expression.BinaryExpression)
	if !ok {
		t.Errorf("Expected BinaryExpression, got %T", expr)
		return
	}

	if binary.Operator != operator {
		t.Errorf("Expected operator '%s', got '%s'", operator, binary.Operator)
	}

	ident, ok := binary.Left.(*expression.IdentifierExpression)
	if !ok {
		t.Errorf("Expected IdentifierExpression, got %T", binary.Left)
		return
	}

	if ident.Name != leftName {
		t.Errorf("Expected identifier '%s', got '%s'", leftName, ident.Name)
	}

	literal, ok := binary.Right.(*expression.LiteralExpression)
	if !ok {
		t.Errorf("Expected LiteralExpression, got %T", binary.Right)
		return
	}

	if literal.Value != rightValue {
		t.Errorf("Expected value %v, got %v", rightValue, literal.Value)
	}
}

// AssertIncrement checks if a statement is an increment operation (x = x + 1)
func AssertIncrement(t *testing.T, stmt ast.Statement, varName string) {
	AssertBinaryAssignment(t, stmt, varName, "=", &BinaryCheck{
		LeftName:   varName,
		Operator:   "+",
		RightValue: int64(1),
	})
}
