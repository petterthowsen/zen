package parsing

import (
	"testing"
	"zen/lang/parsing/statement"
)

func TestAwaitExpressions(t *testing.T) {
	print("parsing await_expressions.zen\n")
	programNode := ParseTestFile(t, "await_expressions.zen")

	if programNode == nil {
		return
	}

	// await somePromise
	stmt := programNode.Statements[0]
	exprStmt, ok := stmt.(*statement.ExpressionStatement)
	if !ok {
		t.Errorf("Expected ExpressionStatement, got %T", stmt)
		return
	}
	await := AssertAwaitExpression(t, exprStmt.Expression)
	AssertIdentifierExpression(t, await.Expression, "somePromise")

	// await getData()
	stmt = programNode.Statements[1]
	exprStmt, ok = stmt.(*statement.ExpressionStatement)
	if !ok {
		t.Errorf("Expected ExpressionStatement, got %T", stmt)
		return
	}
	await = AssertAwaitExpression(t, exprStmt.Expression)
	call := AssertCallExpression(t, await.Expression, 0)
	AssertIdentifierExpression(t, call.Callee, "getData")

	// await obj.method()
	stmt = programNode.Statements[2]
	exprStmt, ok = stmt.(*statement.ExpressionStatement)
	if !ok {
		t.Errorf("Expected ExpressionStatement, got %T", stmt)
		return
	}
	await = AssertAwaitExpression(t, exprStmt.Expression)
	call = AssertCallExpression(t, await.Expression, 0)
	AssertMemberAccess(t, call.Callee, "obj", "method")

	// await (1 + 2)
	stmt = programNode.Statements[3]
	exprStmt, ok = stmt.(*statement.ExpressionStatement)
	if !ok {
		t.Errorf("Expected ExpressionStatement, got %T", stmt)
		return
	}
	await = AssertAwaitExpression(t, exprStmt.Expression)
	binary := AssertBinaryExpression(t, await.Expression, "+")
	AssertLiteralExpression(t, binary.Left, int64(1))
	AssertLiteralExpression(t, binary.Right, int64(2))

	// async func test() { ... }
	stmt = programNode.Statements[4]
	funcDecl := AssertFuncDeclaration(t, stmt)
	if !funcDecl.Async {
		t.Error("Expected function to be async")
	}

	// Check function body
	if len(funcDecl.Body) != 2 {
		t.Errorf("Expected 2 statements in function body, got %d", len(funcDecl.Body))
		return
	}

	// var result = await fetch()
	varDecl := AssertVarDeclaration(t, funcDecl.Body[0], "result", false, false)
	await = AssertAwaitExpression(t, varDecl.Initializer)
	call = AssertCallExpression(t, await.Expression, 0)
	AssertIdentifierExpression(t, call.Callee, "fetch")

	// await result.json()
	exprStmt, ok = funcDecl.Body[1].(*statement.ExpressionStatement)
	if !ok {
		t.Errorf("Expected ExpressionStatement, got %T", funcDecl.Body[1])
		return
	}
	await = AssertAwaitExpression(t, exprStmt.Expression)
	call = AssertCallExpression(t, await.Expression, 0)
	AssertMemberAccess(t, call.Callee, "result", "json")
}
