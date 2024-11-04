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

	// print tokens
	// print("tokens:\n")
	// for _, token := range tokens {
	// 	print(token.String() + "\n")
	// }

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

	// Save AST to file regardless of errors
	if program != nil {
		astStr := program.String(0)
		astPath := strings.TrimSuffix(path, ".zen") + ".ast"
		if err := os.WriteFile(astPath, []byte(astStr), 0644); err != nil {
			t.Logf("Warning: Failed to write AST file: %v", err)
		}
	}

	// Report any errors
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

	return program
}

// AssertVarDeclaration checks if a statement is a variable declaration with expected properties
func AssertVarDeclaration(t *testing.T, stmt ast.Statement, name string, isConst bool, nullable bool) *statement.VarDeclarationNode {
	varDecl, ok := stmt.(*statement.VarDeclarationNode)
	if !ok {
		t.Errorf("Expected VarDeclarationNode, got %T", stmt)
		return nil
	}

	if varDecl.Name != name {
		t.Errorf("Expected name %q, got %q", name, varDecl.Name)
	}

	if varDecl.IsConstant != isConst {
		t.Errorf("Expected const=%v, got %v", isConst, varDecl.IsConstant)
	}

	if varDecl.IsNullable != nullable {
		t.Errorf("Expected nullable=%v, got %v", nullable, varDecl.IsNullable)
	}

	return varDecl
}

// AssertVarDeclarationWithType checks if a statement is a variable declaration with expected type
func AssertVarDeclarationWithType(t *testing.T, stmt ast.Statement, name string, isConst bool, nullable bool, typeCheck func(t *testing.T, typ ast.Expression)) *statement.VarDeclarationNode {
	varDecl := AssertVarDeclaration(t, stmt, name, isConst, nullable)
	if varDecl == nil {
		return nil
	}

	if varDecl.Type == nil {
		t.Error("Expected type annotation, got nil")
		return nil
	}

	typeCheck(t, varDecl.Type)
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

// AssertBasicType checks if an expression is a basic type with expected name
func AssertBasicType(t *testing.T, expr ast.Expression, expectedName string) *expression.BasicType {
	basicType, ok := expr.(*expression.BasicType)
	if !ok {
		t.Errorf("Expected BasicType, got %T", expr)
		return nil
	}
	if basicType.Name != expectedName {
		t.Errorf("Expected type name '%s', got '%s'", expectedName, basicType.Name)
	}
	return basicType
}

// AssertParametricType checks if an expression is a parametric type with expected base type
func AssertParametricType(t *testing.T, expr ast.Expression, expectedBaseType string, expectedParamCount int) *expression.ParametricType {
	paramType, ok := expr.(*expression.ParametricType)
	if !ok {
		t.Errorf("Expected ParametricType, got %T", expr)
		return nil
	}
	if paramType.BaseType != expectedBaseType {
		t.Errorf("Expected base type '%s', got '%s'", expectedBaseType, paramType.BaseType)
	}
	if len(paramType.Parameters) != expectedParamCount {
		t.Errorf("Expected %d parameters, got %d", expectedParamCount, len(paramType.Parameters))
	}
	return paramType
}

// AssertTypeParameter checks if a parameter is a type parameter with expected name
func AssertTypeParameter(t *testing.T, param expression.Parameter, expectedName string) {
	if !param.IsType {
		t.Error("Expected type parameter, got value parameter")
		return
	}
	name, ok := param.Value.(string)
	if !ok {
		t.Errorf("Expected string value, got %T", param.Value)
		return
	}
	if name != expectedName {
		t.Errorf("Expected type name '%s', got '%s'", expectedName, name)
	}
}

// AssertValueParameter checks if a parameter is a value parameter with expected value
func AssertValueParameter(t *testing.T, param expression.Parameter, expectedValue int64) {
	if param.IsType {
		t.Error("Expected value parameter, got type parameter")
		return
	}
	value, ok := param.Value.(int64)
	if !ok {
		t.Errorf("Expected int64 value, got %T", param.Value)
		return
	}
	if value != expectedValue {
		t.Errorf("Expected value %d, got %d", expectedValue, value)
	}
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

// AssertMemberAccess checks if an expression is a member access with expected properties
func AssertMemberAccess(t *testing.T, expr ast.Expression, objectName string, propertyName string) *expression.MemberAccessExpression {
	memberAccess, ok := expr.(*expression.MemberAccessExpression)
	if !ok {
		t.Errorf("Expected MemberAccessExpression, got %T", expr)
		return nil
	}

	if memberAccess.Property != propertyName {
		t.Errorf("Expected property name '%s', got '%s'", propertyName, memberAccess.Property)
	}

	// Object can be either an identifier or another member access expression
	switch obj := memberAccess.Object.(type) {
	case *expression.IdentifierExpression:
		if obj.Name != objectName {
			t.Errorf("Expected object name '%s', got '%s'", objectName, obj.Name)
		}
	case *expression.MemberAccessExpression:
		// For chained access, we only check the property name
		// The full chain should be checked with AssertChainedMemberAccess
		if obj.Property != objectName {
			t.Errorf("Expected object property '%s', got '%s'", objectName, obj.Property)
		}
	default:
		t.Errorf("Expected IdentifierExpression or MemberAccessExpression for object, got %T", memberAccess.Object)
		return nil
	}

	return memberAccess
}

// AssertChainedMemberAccess checks if an expression is a chained member access
func AssertChainedMemberAccess(t *testing.T, expr ast.Expression, chain ...string) ast.Expression {
	if len(chain) < 2 {
		t.Error("Chain must have at least 2 elements")
		return nil
	}

	// For chained access like person.address.city, we get MemberAccess nodes
	// nested from right to left, so we traverse them in reverse
	current := expr
	for i := len(chain) - 1; i > 0; i-- {
		memberAccess, ok := current.(*expression.MemberAccessExpression)
		if !ok {
			t.Errorf("Expected MemberAccessExpression at chain index %d, got %T", i, current)
			return nil
		}

		if memberAccess.Property != chain[i] {
			t.Errorf("Expected property name '%s' at chain index %d, got '%s'", chain[i], i, memberAccess.Property)
			return nil
		}

		if i == len(chain)-1 {
			// First level should have an identifier as object
			ident, ok := memberAccess.Object.(*expression.IdentifierExpression)
			if !ok {
				t.Errorf("Expected IdentifierExpression for first object, got %T", memberAccess.Object)
				return nil
			}
			if ident.Name != chain[0] {
				t.Errorf("Expected object name '%s', got '%s'", chain[0], ident.Name)
			}
		}

		current = memberAccess.Object
	}

	return expr
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
