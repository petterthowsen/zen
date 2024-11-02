package parsing

import (
	"testing"
	"zen/lang/parsing/statement"
)

func TestFuncDeclarations(t *testing.T) {
	print("parsing func_declarations.zen\n")
	programNode := ParseTestFile(t, "func_declarations.zen")

	if programNode == nil {
		return
	}

	// func sayHello() {
	// 		print("Hello!")
	// }
	sayHello := AssertFuncDeclaration(t, programNode.Statements[0])
	if sayHello.Name != "sayHello" {
		t.Errorf("Expected func sayHello, got func %s", sayHello.Name)
		return
	}

	// verify that sayHello.Body[0] is instance of CallExpression
	_, ok := sayHello.Body[0].(*statement.ExpressionStatement)
	if !ok {
		t.Errorf("Expected CallStmt, got %T", sayHello.Body[0])
	}

	say := AssertFuncDeclaration(t, programNode.Statements[1])
	if say.Name != "say" {
		t.Errorf("Expected func say, got func %s", say.Name)
		return
	}

	// func say(something:string) {
	// 		print(something)
	// }
	if say.Name != "say" {
		// here,
		t.Errorf("Expected func say, got func %s", say.Name)
		return
	}

	// verify that say.Parameters[0] exists
	if len(say.Parameters) != 1 {
		t.Errorf("Expected 1 parameter, got %d", len(say.Parameters))
	}

	// and that it is named "something"
	sayParamSomething := say.Parameters[0]
	if sayParamSomething.Name != "something" {
		t.Errorf("Expected parameter 'something', got %s", sayParamSomething.Name)
	}

	// and is of type 'string'
	if sayParamSomething.Type != "string" {
		t.Errorf("Expected parameter type 'string', got %s", sayParamSomething.Type)
	}

	add := AssertFuncDeclaration(t, programNode.Statements[2])
	if add.Name != "add" {
		t.Errorf("Expected func add, got func %s", add.Name)
		return
	}

	if len(add.Parameters) != 2 {
		t.Errorf("Expected 2 parameters for add function, got %d", len(add.Parameters))
		return
	}

	addA := add.Parameters[0]
	addB := add.Parameters[1]

	// verify paramater a
	if addA.Name != "a" {
		t.Errorf("Expected parameter 'a', got %s", addA.Name)
	}

	if addA.Type != "int" {
		t.Errorf("Expected parameter type 'int', got %s", addA.Type)
	}

	if addA.IsNullable != false {
		t.Errorf("Expected parameter a to be NOT nullable, got %v", addA.IsNullable)
	}

	// verify parameter b
	if addB.Name != "b" {
		t.Errorf("Expected parameter 'b', got %s", addB.Name)
	}

	if addB.IsNullable != false {
		t.Errorf("Expected parameter b to be NOT nullable, got %v", addB.IsNullable)
	}

	if addB.Type != "int" {
		t.Errorf("Expected parameter type 'int', got %s", addB.Type)
	}

	// verify return type
	if add.ReturnType != "int" {
		t.Errorf("Expected return type 'int', got %s", add.ReturnType)
	}

	// optional parameter
	log := AssertFuncDeclaration(t, programNode.Statements[3])
	if log.Name != "log" {
		t.Errorf("Expected func log, got func %s", log.Name)
		return
	}

	if len(log.Parameters) != 2 {
		t.Errorf("Expected 2 parameters for log function, got %d", len(log.Parameters))
		return
	}

	logMessage := log.Parameters[0]
	logSuffix := log.Parameters[1]

	// verify logMessage is called "message"
	if logMessage.Name != "message" {
		t.Errorf("Expected parameter 'message', got %s", logMessage.Name)
	}

	// logMessage is of type string
	if logMessage.Type != "string" {
		t.Errorf("Expected parameter type 'string', got %s", logMessage.Type)
	}

	// not nullable
	if logMessage.IsNullable != false {
		t.Errorf("Expected parameter message to be NOT nullable, got %v", logMessage.IsNullable)
	}

	// verify logSuffix is called "suffix"
	if logSuffix.Name != "suffix" {
		t.Errorf("Expected parameter 'suffix', got %s", logSuffix.Name)
	}

	if logSuffix.Type != "string" {
		t.Errorf("Expected parameter type 'string', got %s", logSuffix.Type)
	}

	// verify suffix is nullable
	if logSuffix.IsNullable != true {
		t.Errorf("Expected parameter suffix to be nullable, got %v", logSuffix.IsNullable)
	}

	// verify return type
	if log.ReturnType != "void" {
		t.Errorf("Expected return type 'void', got %s", log.ReturnType)
	}
}
