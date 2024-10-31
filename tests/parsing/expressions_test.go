package parsing

import (
	"os"
	"strings"
	"testing"
	"zen/lang/parsing/ast"
	"zen/lang/parsing/expression"
	"zen/lang/parsing/statement"
)

func TestExpressionsFile(t *testing.T) {
	// Load and parse the test file
	path := getTestDataPath("expressions.zen")
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
		return
	}

	if program == nil {
		t.Error("Expected program node, got nil")
		return
	}

	// Save AST to file
	astStr := program.String(0)
	astPath := strings.TrimSuffix(path, ".zen") + ".ast"
	if err := os.WriteFile(astPath, []byte(astStr), 0644); err != nil {
		t.Logf("Warning: Failed to write AST file: %v", err)
	}

	// Verify expressions in variable declarations
	expectedDecls := []struct {
		name      string
		isConst   bool
		checkInit func(t *testing.T, init ast.Expression)
	}{
		{
			name: "answer",
			checkInit: func(t *testing.T, init ast.Expression) {
				lit, ok := init.(*expression.LiteralExpression)
				if !ok {
					t.Errorf("Expected LiteralExpression, got %T", init)
					return
				}
				if lit.Value != int64(42) {
					t.Errorf("Expected literal value 42, got %v", lit.Value)
				}
			},
		},
		{
			name:    "PI",
			isConst: true,
			checkInit: func(t *testing.T, init ast.Expression) {
				lit, ok := init.(*expression.LiteralExpression)
				if !ok {
					t.Errorf("Expected LiteralExpression, got %T", init)
					return
				}
				if lit.Value != float64(3.14) {
					t.Errorf("Expected literal value 3.14, got %v", lit.Value)
				}
			},
		},
		{
			name: "zen",
			checkInit: func(t *testing.T, init ast.Expression) {
				lit, ok := init.(*expression.LiteralExpression)
				if !ok {
					t.Errorf("Expected LiteralExpression, got %T", init)
					return
				}
				if lit.Value != true {
					t.Errorf("Expected literal value true, got %v", lit.Value)
				}
			},
		},
		{
			name: "nope",
			checkInit: func(t *testing.T, init ast.Expression) {
				lit, ok := init.(*expression.LiteralExpression)
				if !ok {
					t.Errorf("Expected LiteralExpression, got %T", init)
					return
				}
				if lit.Value != false {
					t.Errorf("Expected literal value false, got %v", lit.Value)
				}
			},
		},
		{
			name: "title",
			checkInit: func(t *testing.T, init ast.Expression) {
				lit, ok := init.(*expression.LiteralExpression)
				if !ok {
					t.Errorf("Expected LiteralExpression, got %T", init)
					return
				}
				if lit.Value != "Zen Lang!" {
					t.Errorf("Expected literal value \"Zen Lang!\", got %v", lit.Value)
				}
			},
		},
		{
			name: "plus",
			checkInit: func(t *testing.T, init ast.Expression) {
				bin, ok := init.(*expression.BinaryExpression)
				if !ok {
					t.Errorf("Expected BinaryExpression, got %T", init)
					return
				}
				if bin.Operator != "+" {
					t.Errorf("Expected operator +, got %s", bin.Operator)
				}
				left, ok := bin.Left.(*expression.LiteralExpression)
				if !ok || left.Value != int64(1) {
					t.Errorf("Expected left operand 1, got %v", bin.Left)
				}
				right, ok := bin.Right.(*expression.LiteralExpression)
				if !ok || right.Value != int64(2) {
					t.Errorf("Expected right operand 2, got %v", bin.Right)
				}
			},
		},
		{
			name: "minus",
			checkInit: func(t *testing.T, init ast.Expression) {
				bin, ok := init.(*expression.BinaryExpression)
				if !ok {
					t.Errorf("Expected BinaryExpression, got %T", init)
					return
				}
				if bin.Operator != "-" {
					t.Errorf("Expected operator -, got %s", bin.Operator)
				}
				left, ok := bin.Left.(*expression.LiteralExpression)
				if !ok || left.Value != int64(5) {
					t.Errorf("Expected left operand 5, got %v", bin.Left)
				}
				right, ok := bin.Right.(*expression.LiteralExpression)
				if !ok || right.Value != int64(6) {
					t.Errorf("Expected right operand 6, got %v", bin.Right)
				}
			},
		},
		{
			name: "mult",
			checkInit: func(t *testing.T, init ast.Expression) {
				bin, ok := init.(*expression.BinaryExpression)
				if !ok {
					t.Errorf("Expected BinaryExpression, got %T", init)
					return
				}
				if bin.Operator != "*" {
					t.Errorf("Expected operator *, got %s", bin.Operator)
				}
				left, ok := bin.Left.(*expression.LiteralExpression)
				if !ok || left.Value != int64(3) {
					t.Errorf("Expected left operand 3, got %v", bin.Left)
				}
				right, ok := bin.Right.(*expression.LiteralExpression)
				if !ok || right.Value != int64(4) {
					t.Errorf("Expected right operand 4, got %v", bin.Right)
				}
			},
		},
		{
			name: "divide",
			checkInit: func(t *testing.T, init ast.Expression) {
				bin, ok := init.(*expression.BinaryExpression)
				if !ok {
					t.Errorf("Expected BinaryExpression, got %T", init)
					return
				}
				if bin.Operator != "/" {
					t.Errorf("Expected operator /, got %s", bin.Operator)
				}
				left, ok := bin.Left.(*expression.LiteralExpression)
				if !ok || left.Value != int64(10) {
					t.Errorf("Expected left operand 10, got %v", bin.Left)
				}
				right, ok := bin.Right.(*expression.LiteralExpression)
				if !ok || right.Value != int64(2) {
					t.Errorf("Expected right operand 2, got %v", bin.Right)
				}
			},
		},
		{
			name: "complex1",
			checkInit: func(t *testing.T, init ast.Expression) {
				bin, ok := init.(*expression.BinaryExpression)
				if !ok {
					t.Errorf("Expected BinaryExpression, got %T", init)
					return
				}
				if bin.Operator != "+" {
					t.Errorf("Expected operator +, got %s", bin.Operator)
				}
				left, ok := bin.Left.(*expression.LiteralExpression)
				if !ok || left.Value != int64(2) {
					t.Errorf("Expected left operand 2, got %v", bin.Left)
				}
				right, ok := bin.Right.(*expression.BinaryExpression)
				if !ok {
					t.Errorf("Expected right operand to be BinaryExpression, got %T", bin.Right)
					return
				}
				if right.Operator != "*" {
					t.Errorf("Expected right operator *, got %s", right.Operator)
				}
				rightLeft, ok := right.Left.(*expression.LiteralExpression)
				if !ok || rightLeft.Value != int64(3) {
					t.Errorf("Expected right left operand 3, got %v", right.Left)
				}
				rightRight, ok := right.Right.(*expression.LiteralExpression)
				if !ok || rightRight.Value != int64(4) {
					t.Errorf("Expected right right operand 4, got %v", right.Right)
				}
			},
		},
		{
			name: "complex2",
			checkInit: func(t *testing.T, init ast.Expression) {
				bin, ok := init.(*expression.BinaryExpression)
				if !ok {
					t.Errorf("Expected BinaryExpression, got %T", init)
					return
				}
				if bin.Operator != "*" {
					t.Errorf("Expected operator *, got %s", bin.Operator)
				}
				left, ok := bin.Left.(*expression.BinaryExpression)
				if !ok {
					t.Errorf("Expected left operand to be BinaryExpression, got %T", bin.Left)
					return
				}
				if left.Operator != "+" {
					t.Errorf("Expected left operator +, got %s", left.Operator)
				}
				leftLeft, ok := left.Left.(*expression.LiteralExpression)
				if !ok || leftLeft.Value != int64(2) {
					t.Errorf("Expected left left operand 2, got %v", left.Left)
				}
				leftRight, ok := left.Right.(*expression.LiteralExpression)
				if !ok || leftRight.Value != int64(3) {
					t.Errorf("Expected left right operand 3, got %v", left.Right)
				}
				right, ok := bin.Right.(*expression.LiteralExpression)
				if !ok || right.Value != int64(4) {
					t.Errorf("Expected right operand 4, got %v", bin.Right)
				}
			},
		},
		{
			name: "negative42",
			checkInit: func(t *testing.T, init ast.Expression) {
				un, ok := init.(*expression.UnaryExpression)
				if !ok {
					t.Errorf("Expected UnaryExpression, got %T", init)
					return
				}
				if un.Operator != "-" {
					t.Errorf("Expected operator -, got %s", un.Operator)
				}
				lit, ok := un.Expression.(*expression.LiteralExpression)
				if !ok || lit.Value != int64(42) {
					t.Errorf("Expected operand 42, got %v", un.Expression)
				}
			},
		},
		{
			name: "negated",
			checkInit: func(t *testing.T, init ast.Expression) {
				un, ok := init.(*expression.UnaryExpression)
				if !ok {
					t.Errorf("Expected UnaryExpression, got %T", init)
					return
				}
				if un.Operator != "not" {
					t.Errorf("Expected operator not, got %s", un.Operator)
				}
				lit, ok := un.Expression.(*expression.LiteralExpression)
				if !ok || lit.Value != true {
					t.Errorf("Expected operand true, got %v", un.Expression)
				}
			},
		},
		{
			name: "result",
			checkInit: func(t *testing.T, init ast.Expression) {
				// -5 + plus * 1.5 + function()
				bin1, ok := init.(*expression.BinaryExpression)
				if !ok || bin1.Operator != "+" {
					t.Errorf("Expected top-level binary + expression, got %T", init)
					return
				}

				bin2, ok := bin1.Left.(*expression.BinaryExpression)
				if !ok || bin2.Operator != "+" {
					t.Errorf("Expected second-level binary + expression, got %T", bin1.Left)
					return
				}

				un, ok := bin2.Left.(*expression.UnaryExpression)
				if !ok || un.Operator != "-" {
					t.Errorf("Expected unary - expression, got %T", bin2.Left)
					return
				}

				lit1, ok := un.Expression.(*expression.LiteralExpression)
				if !ok || lit1.Value != int64(5) {
					t.Errorf("Expected literal 5, got %v", un.Expression)
					return
				}

				bin3, ok := bin2.Right.(*expression.BinaryExpression)
				if !ok || bin3.Operator != "*" {
					t.Errorf("Expected binary * expression, got %T", bin2.Right)
					return
				}

				id1, ok := bin3.Left.(*expression.IdentifierExpression)
				if !ok || id1.Name != "plus" {
					t.Errorf("Expected identifier 'plus', got %T", bin3.Left)
					return
				}

				lit2, ok := bin3.Right.(*expression.LiteralExpression)
				if !ok || lit2.Value != float64(1.5) {
					t.Errorf("Expected literal 1.5, got %v", bin3.Right)
					return
				}

				call, ok := bin1.Right.(*expression.CallExpression)
				if !ok {
					t.Errorf("Expected CallExpression, got %T", bin1.Right)
					return
				}

				id2, ok := call.Callee.(*expression.IdentifierExpression)
				if !ok || id2.Name != "myFunc" {
					t.Errorf("Expected identifier 'myFunc', got %T", call.Callee)
					return
				}

				if len(call.Arguments) != 0 {
					t.Errorf("Expected 0 arguments, got %d", len(call.Arguments))
				}
			},
		},
	}

	for i, exp := range expectedDecls {
		if i >= len(program.Statements) {
			t.Errorf("Missing declaration %d: %s", i, exp.name)
			continue
		}

		stmt := program.Statements[i]
		t.Logf("Checking declaration %d: %s", i, stmt.String(0))

		// Type assert to VarDeclarationNode
		varDecl, ok := stmt.(*statement.VarDeclarationNode)
		if !ok {
			t.Errorf("Statement %d: expected VarDeclarationNode, got %T", i, stmt)
			continue
		}

		if varDecl.Name != exp.name {
			t.Errorf("Declaration %d: expected name %q, got %q", i, exp.name, varDecl.Name)
		}

		if varDecl.IsConstant != exp.isConst {
			t.Errorf("Declaration %d: expected const=%v, got %v", i, exp.isConst, varDecl.IsConstant)
		}

		if exp.checkInit != nil {
			if varDecl.Initializer == nil {
				t.Errorf("Declaration %d: missing initializer", i)
				continue
			}
			exp.checkInit(t, varDecl.Initializer)
		}
	}
}

func TestExpressionErrors(t *testing.T) {
	cases := []struct {
		input string
	}{
		{"var x = 1 +"},
		{"var x = * 2"},
		{"var x = (1 + 2"},
		{"var x = 1 + + 2"},
	}

	for _, c := range cases {
		program, errors := ParseString(c.input)
		if len(errors) == 0 {
			t.Errorf("Input %q: expected error, got none", c.input)
			continue
		}
		if program != nil {
			t.Errorf("Input %q: expected nil program for error case", c.input)
		}
	}
}
