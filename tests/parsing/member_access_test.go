package parsing

import (
	"testing"
	"zen/lang/parsing/expression"
	"zen/lang/parsing/statement"
)

func TestMemberAccess(t *testing.T) {
	program := ParseTestFile(t, "member_access.zen")
	if program == nil {
		return
	}
	// The following tests are fine
	// Basic member access: person.name
	AssertVarDeclaration(t, program.Statements[0], "name", "", false, false)
	varDecl := program.Statements[0].(*statement.VarDeclarationNode)
	AssertMemberAccess(t, varDecl.Initializer, "person", "name")

	// Another basic member access: user.age
	AssertVarDeclaration(t, program.Statements[1], "age", "", false, false)
	varDecl = program.Statements[1].(*statement.VarDeclarationNode)
	AssertMemberAccess(t, varDecl.Initializer, "user", "age")

	// Chained member access: person.address.city
	AssertVarDeclaration(t, program.Statements[2], "city", "", false, false)
	varDecl = program.Statements[2].(*statement.VarDeclarationNode)
	memberAccess := varDecl.Initializer.(*expression.MemberAccessExpression)
	if memberAccess.Property != "city" {
		t.Errorf("Expected property 'city', got '%s'", memberAccess.Property)
	}
	innerAccess := memberAccess.Object.(*expression.MemberAccessExpression)
	if innerAccess.Property != "address" {
		t.Errorf("Expected property 'address', got '%s'", innerAccess.Property)
	}
	ident := innerAccess.Object.(*expression.IdentifierExpression)
	if ident.Name != "person" {
		t.Errorf("Expected identifier 'person', got '%s'", ident.Name)
	}

	// Member access with function calls: name.length()
	AssertVarDeclaration(t, program.Statements[3], "length", "", false, false)
	varDecl = program.Statements[3].(*statement.VarDeclarationNode)
	call := varDecl.Initializer.(*expression.CallExpression)
	AssertMemberAccess(t, call.Callee, "name", "length")

	// Member access with function calls: list.getItems()
	AssertVarDeclaration(t, program.Statements[4], "items", "", false, false)
	varDecl = program.Statements[4].(*statement.VarDeclarationNode)
	call = varDecl.Initializer.(*expression.CallExpression)
	AssertMemberAccess(t, call.Callee, "list", "getItems")

	// Member access in expressions: person.firstName + " " + person.lastName
	AssertVarDeclaration(t, program.Statements[5], "fullName", "", false, false)
	varDecl = program.Statements[5].(*statement.VarDeclarationNode)
	binary := varDecl.Initializer.(*expression.BinaryExpression)
	leftBinary := binary.Left.(*expression.BinaryExpression)
	AssertMemberAccess(t, leftBinary.Left, "person", "firstName")
	AssertMemberAccess(t, binary.Right, "person", "lastName")

	// Member access in assignments: person.age = 25
	exprStmt := program.Statements[6].(*statement.ExpressionStatement)
	binary = exprStmt.Expression.(*expression.BinaryExpression)
	AssertMemberAccess(t, binary.Left, "person", "age")

	// Chained member access in assignments: user.contact.email = "test@example.com"
	exprStmt = program.Statements[7].(*statement.ExpressionStatement)
	binary = exprStmt.Expression.(*expression.BinaryExpression)
	memberAccess = binary.Left.(*expression.MemberAccessExpression)
	if memberAccess.Property != "email" {
		t.Errorf("Expected property 'email', got '%s'", memberAccess.Property)
	}
	innerAccess = memberAccess.Object.(*expression.MemberAccessExpression)
	if innerAccess.Property != "contact" {
		t.Errorf("Expected property 'contact', got '%s'", innerAccess.Property)
	}
	ident = innerAccess.Object.(*expression.IdentifierExpression)
	if ident.Name != "user" {
		t.Errorf("Expected identifier 'user', got '%s'", ident.Name)
	}

	// Member access in conditions: person.age > 18
	ifStmt := program.Statements[8].(*statement.IfStatement)
	binary = ifStmt.PrimaryCondition.(*expression.BinaryExpression)
	AssertMemberAccess(t, binary.Left, "person", "age")

	// this one fails:

	// Member access in function calls: print(person.name)
	exprStmt = ifStmt.PrimaryBlock[0].(*statement.ExpressionStatement)
	call = exprStmt.Expression.(*expression.CallExpression)
	AssertMemberAccess(t, call.Arguments[0], "person", "name")

	// Member access in function calls: print(user.getFullName())
	exprStmt = program.Statements[9].(*statement.ExpressionStatement)
	call = exprStmt.Expression.(*expression.CallExpression)
	innerCall := call.Arguments[0].(*expression.CallExpression)
	AssertMemberAccess(t, innerCall.Callee, "user", "getFullName")

	// Chained member access with function calls: validate(person.address.isValid())
	exprStmt = program.Statements[10].(*statement.ExpressionStatement)
	call = exprStmt.Expression.(*expression.CallExpression)
	innerCall = call.Arguments[0].(*expression.CallExpression)
	memberAccess = innerCall.Callee.(*expression.MemberAccessExpression)
	if memberAccess.Property != "isValid" {
		t.Errorf("Expected property 'isValid', got '%s'", memberAccess.Property)
	}
	innerAccess = memberAccess.Object.(*expression.MemberAccessExpression)
	if innerAccess.Property != "address" {
		t.Errorf("Expected property 'address', got '%s'", innerAccess.Property)
	}
	ident = innerAccess.Object.(*expression.IdentifierExpression)
	if ident.Name != "person" {
		t.Errorf("Expected identifier 'person', got '%s'", ident.Name)
	}
}
