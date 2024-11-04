package parsing

import (
	"testing"
	"zen/lang/parsing/ast"
)

func TestExpressionsFile(t *testing.T) {
	program := ParseTestFile(t, "expressions.zen")
	if program == nil {
		return
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
				AssertLiteralExpression(t, init, int64(42))
			},
		},
		{
			name:    "PI",
			isConst: true,
			checkInit: func(t *testing.T, init ast.Expression) {
				AssertLiteralExpression(t, init, float64(3.14))
			},
		},
		{
			name: "zen",
			checkInit: func(t *testing.T, init ast.Expression) {
				AssertLiteralExpression(t, init, true)
			},
		},
		{
			name: "nope",
			checkInit: func(t *testing.T, init ast.Expression) {
				AssertLiteralExpression(t, init, false)
			},
		},
		{
			name: "title",
			checkInit: func(t *testing.T, init ast.Expression) {
				AssertLiteralExpression(t, init, "Zen Lang!")
			},
		},
		{
			name: "plus",
			checkInit: func(t *testing.T, init ast.Expression) {
				bin := AssertBinaryExpression(t, init, "+")
				if bin == nil {
					return
				}
				AssertLiteralExpression(t, bin.Left, int64(1))
				AssertLiteralExpression(t, bin.Right, int64(2))
			},
		},
		{
			name: "minus",
			checkInit: func(t *testing.T, init ast.Expression) {
				bin := AssertBinaryExpression(t, init, "-")
				if bin == nil {
					return
				}
				AssertLiteralExpression(t, bin.Left, int64(5))
				AssertLiteralExpression(t, bin.Right, int64(6))
			},
		},
		{
			name: "mult",
			checkInit: func(t *testing.T, init ast.Expression) {
				bin := AssertBinaryExpression(t, init, "*")
				if bin == nil {
					return
				}
				AssertLiteralExpression(t, bin.Left, int64(3))
				AssertLiteralExpression(t, bin.Right, int64(4))
			},
		},
		{
			name: "divide",
			checkInit: func(t *testing.T, init ast.Expression) {
				bin := AssertBinaryExpression(t, init, "/")
				if bin == nil {
					return
				}
				AssertLiteralExpression(t, bin.Left, int64(10))
				AssertLiteralExpression(t, bin.Right, int64(2))
			},
		},
		{
			name: "complex1",
			checkInit: func(t *testing.T, init ast.Expression) {
				bin := AssertBinaryExpression(t, init, "+")
				if bin == nil {
					return
				}
				AssertLiteralExpression(t, bin.Left, int64(2))

				right := AssertBinaryExpression(t, bin.Right, "*")
				if right == nil {
					return
				}
				AssertLiteralExpression(t, right.Left, int64(3))
				AssertLiteralExpression(t, right.Right, int64(4))
			},
		},
		{
			name: "complex2",
			checkInit: func(t *testing.T, init ast.Expression) {
				bin := AssertBinaryExpression(t, init, "*")
				if bin == nil {
					return
				}

				left := AssertBinaryExpression(t, bin.Left, "+")
				if left == nil {
					return
				}
				AssertLiteralExpression(t, left.Left, int64(2))
				AssertLiteralExpression(t, left.Right, int64(3))

				AssertLiteralExpression(t, bin.Right, int64(4))
			},
		},
		{
			name: "negative42",
			checkInit: func(t *testing.T, init ast.Expression) {
				un := AssertUnaryExpression(t, init, "-")
				if un == nil {
					return
				}
				AssertLiteralExpression(t, un.Expression, int64(42))
			},
		},
		{
			name: "negated",
			checkInit: func(t *testing.T, init ast.Expression) {
				un := AssertUnaryExpression(t, init, "not")
				if un == nil {
					return
				}
				AssertLiteralExpression(t, un.Expression, true)
			},
		},
		{
			name: "result",
			checkInit: func(t *testing.T, init ast.Expression) {
				// -5 + plus * 1.5 + myFunc()
				bin1 := AssertBinaryExpression(t, init, "+")
				if bin1 == nil {
					return
				}

				bin2 := AssertBinaryExpression(t, bin1.Left, "+")
				if bin2 == nil {
					return
				}

				un := AssertUnaryExpression(t, bin2.Left, "-")
				if un == nil {
					return
				}
				AssertLiteralExpression(t, un.Expression, int64(5))

				bin3 := AssertBinaryExpression(t, bin2.Right, "*")
				if bin3 == nil {
					return
				}
				AssertIdentifierExpression(t, bin3.Left, "plus")
				AssertLiteralExpression(t, bin3.Right, float64(1.5))

				call := AssertCallExpression(t, bin1.Right, 0)
				if call == nil {
					return
				}
				AssertIdentifierExpression(t, call.Callee, "myFunc")
			},
		},
	}

	for i, exp := range expectedDecls {
		if i >= len(program.Statements) {
			t.Errorf("Missing declaration %d: %s", i, exp.name)
			continue
		}

		stmt := program.Statements[i]
		//t.Logf("Checking declaration %d: %s", i, stmt.String(0))

		varDecl := AssertVarDeclaration(t, stmt, exp.name, exp.isConst, false)
		if varDecl == nil {
			continue
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
	cases := []string{
		"var x = 1 +",
		"var x = * 2",
		"var x = (1 + 2",
		"var x = 1 + + 2",
	}

	for _, input := range cases {
		AssertParseError(t, input)
	}
}
