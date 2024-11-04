package parsing

import (
	"testing"
	"zen/lang/parsing/expression"
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
	paramType, ok := sayParamSomething.Type.(*expression.BasicType)
	if !ok {
		t.Errorf("Expected BasicType for parameter type, got %T", sayParamSomething.Type)
	} else if paramType.Name != "string" {
		t.Errorf("Expected parameter type 'string', got %s", paramType.Name)
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

	paramTypeA, ok := addA.Type.(*expression.BasicType)
	if !ok {
		t.Errorf("Expected BasicType for parameter type, got %T", addA.Type)
	} else if paramTypeA.Name != "int" {
		t.Errorf("Expected parameter type 'int', got %s", paramTypeA.Name)
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

	paramTypeB, ok := addB.Type.(*expression.BasicType)
	if !ok {
		t.Errorf("Expected BasicType for parameter type, got %T", addB.Type)
	} else if paramTypeB.Name != "int" {
		t.Errorf("Expected parameter type 'int', got %s", paramTypeB.Name)
	}

	// verify return type
	returnType, ok := add.ReturnType.(*expression.BasicType)
	if !ok {
		t.Errorf("Expected BasicType for return type, got %T", add.ReturnType)
	} else if returnType.Name != "int" {
		t.Errorf("Expected return type 'int', got %s", returnType.Name)
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
	messageType, ok := logMessage.Type.(*expression.BasicType)
	if !ok {
		t.Errorf("Expected BasicType for parameter type, got %T", logMessage.Type)
	} else if messageType.Name != "string" {
		t.Errorf("Expected parameter type 'string', got %s", messageType.Name)
	}

	// not nullable
	if logMessage.IsNullable != false {
		t.Errorf("Expected parameter message to be NOT nullable, got %v", logMessage.IsNullable)
	}

	// verify logSuffix is called "suffix"
	if logSuffix.Name != "suffix" {
		t.Errorf("Expected parameter 'suffix', got %s", logSuffix.Name)
	}

	suffixType, ok := logSuffix.Type.(*expression.BasicType)
	if !ok {
		t.Errorf("Expected BasicType for parameter type, got %T", logSuffix.Type)
	} else if suffixType.Name != "string" {
		t.Errorf("Expected parameter type 'string', got %s", suffixType.Name)
	}

	// verify suffix is nullable
	if logSuffix.IsNullable != true {
		t.Errorf("Expected parameter suffix to be nullable, got %v", logSuffix.IsNullable)
	}

	// verify return type
	logReturnType, ok := log.ReturnType.(*expression.BasicType)
	if !ok {
		t.Errorf("Expected BasicType for return type, got %T", log.ReturnType)
	} else if logReturnType.Name != "void" {
		t.Errorf("Expected return type 'void', got %s", logReturnType.Name)
	}
}
