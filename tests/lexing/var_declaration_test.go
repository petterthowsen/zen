package lexing

import (
	"testing"
	"zen/lang/lexing"
)

func TestVarDeclaration(t *testing.T) {
	expected := []TokenAssert{
		// var name = "zen"
		{Type: lexing.KEYWORD, Literal: "var"},
		{Type: lexing.IDENTIFIER, Literal: "name"},
		{Type: lexing.ASSIGN, Literal: "="},
		{Type: lexing.STRING, Literal: "zen"},

		// var age = 25
		{Type: lexing.KEYWORD, Literal: "var"},
		{Type: lexing.IDENTIFIER, Literal: "age"},
		{Type: lexing.ASSIGN, Literal: "="},
		{Type: lexing.INT, Literal: "25"},

		// var price = 19.99
		{Type: lexing.KEYWORD, Literal: "var"},
		{Type: lexing.IDENTIFIER, Literal: "price"},
		{Type: lexing.ASSIGN, Literal: "="},
		{Type: lexing.FLOAT, Literal: "19.99"},

		// var enabled = true
		{Type: lexing.KEYWORD, Literal: "var"},
		{Type: lexing.IDENTIFIER, Literal: "enabled"},
		{Type: lexing.ASSIGN, Literal: "="},
		{Type: lexing.KEYWORD, Literal: "true"},

		// var count:int = 0
		{Type: lexing.KEYWORD, Literal: "var"},
		{Type: lexing.IDENTIFIER, Literal: "count"},
		{Type: lexing.COLON, Literal: ":"},
		{Type: lexing.KEYWORD, Literal: "int"},
		{Type: lexing.ASSIGN, Literal: "="},
		{Type: lexing.INT, Literal: "0"},

		// var message:string = "hello"
		{Type: lexing.KEYWORD, Literal: "var"},
		{Type: lexing.IDENTIFIER, Literal: "message"},
		{Type: lexing.COLON, Literal: ":"},
		{Type: lexing.KEYWORD, Literal: "string"},
		{Type: lexing.ASSIGN, Literal: "="},
		{Type: lexing.STRING, Literal: "hello"},

		// var rate:float = 0.5
		{Type: lexing.KEYWORD, Literal: "var"},
		{Type: lexing.IDENTIFIER, Literal: "rate"},
		{Type: lexing.COLON, Literal: ":"},
		{Type: lexing.KEYWORD, Literal: "float"},
		{Type: lexing.ASSIGN, Literal: "="},
		{Type: lexing.FLOAT, Literal: "0.5"},

		// var valid:bool = false
		{Type: lexing.KEYWORD, Literal: "var"},
		{Type: lexing.IDENTIFIER, Literal: "valid"},
		{Type: lexing.COLON, Literal: ":"},
		{Type: lexing.KEYWORD, Literal: "bool"},
		{Type: lexing.ASSIGN, Literal: "="},
		{Type: lexing.KEYWORD, Literal: "false"},

		// var optional:string?
		{Type: lexing.KEYWORD, Literal: "var"},
		{Type: lexing.IDENTIFIER, Literal: "optional"},
		{Type: lexing.COLON, Literal: ":"},
		{Type: lexing.KEYWORD, Literal: "string"},
		{Type: lexing.QMARK, Literal: "?"},

		// var maybe:int? = 42
		{Type: lexing.KEYWORD, Literal: "var"},
		{Type: lexing.IDENTIFIER, Literal: "maybe"},
		{Type: lexing.COLON, Literal: ":"},
		{Type: lexing.KEYWORD, Literal: "int"},
		{Type: lexing.QMARK, Literal: "?"},
		{Type: lexing.ASSIGN, Literal: "="},
		{Type: lexing.INT, Literal: "42"},

		// var numbers:Array<int, 3> = [1, 2, 3]
		{Type: lexing.KEYWORD, Literal: "var"},
		{Type: lexing.IDENTIFIER, Literal: "numbers"},
		{Type: lexing.COLON, Literal: ":"},
		{Type: lexing.IDENTIFIER, Literal: "Array"},
		{Type: lexing.LESS, Literal: "<"},
		{Type: lexing.KEYWORD, Literal: "int"},
		{Type: lexing.COMMA, Literal: ","},
		{Type: lexing.INT, Literal: "3"},
		{Type: lexing.GREATER, Literal: ">"},
		{Type: lexing.ASSIGN, Literal: "="},
		{Type: lexing.LEFT_BRACKET, Literal: "["},
		{Type: lexing.INT, Literal: "1"},
		{Type: lexing.COMMA, Literal: ","},
		{Type: lexing.INT, Literal: "2"},
		{Type: lexing.COMMA, Literal: ","},
		{Type: lexing.INT, Literal: "3"},
		{Type: lexing.RIGHT_BRACKET, Literal: "]"},

		// var dynamic:Array<string, ?> = ["a", "b", "c"]
		{Type: lexing.KEYWORD, Literal: "var"},
		{Type: lexing.IDENTIFIER, Literal: "dynamic"},
		{Type: lexing.COLON, Literal: ":"},
		{Type: lexing.IDENTIFIER, Literal: "Array"},
		{Type: lexing.LESS, Literal: "<"},
		{Type: lexing.KEYWORD, Literal: "string"},
		{Type: lexing.COMMA, Literal: ","},
		{Type: lexing.QMARK, Literal: "?"},
		{Type: lexing.GREATER, Literal: ">"},
		{Type: lexing.ASSIGN, Literal: "="},
		{Type: lexing.LEFT_BRACKET, Literal: "["},
		{Type: lexing.STRING, Literal: "a"},
		{Type: lexing.COMMA, Literal: ","},
		{Type: lexing.STRING, Literal: "b"},
		{Type: lexing.COMMA, Literal: ","},
		{Type: lexing.STRING, Literal: "c"},
		{Type: lexing.RIGHT_BRACKET, Literal: "]"},

		// var config:Map<string, any> = { ... }
		{Type: lexing.KEYWORD, Literal: "var"},
		{Type: lexing.IDENTIFIER, Literal: "config"},
		{Type: lexing.COLON, Literal: ":"},
		{Type: lexing.IDENTIFIER, Literal: "Map"},
		{Type: lexing.LESS, Literal: "<"},
		{Type: lexing.KEYWORD, Literal: "string"},
		{Type: lexing.COMMA, Literal: ","},
		{Type: lexing.KEYWORD, Literal: "any"},
		{Type: lexing.GREATER, Literal: ">"},
		{Type: lexing.ASSIGN, Literal: "="},
		{Type: lexing.LEFT_BRACE, Literal: "{"},
		{Type: lexing.STRING, Literal: "debug"},
		{Type: lexing.COLON, Literal: ":"},
		{Type: lexing.KEYWORD, Literal: "true"},
		{Type: lexing.STRING, Literal: "port"},
		{Type: lexing.COLON, Literal: ":"},
		{Type: lexing.INT, Literal: "8080"},
		{Type: lexing.RIGHT_BRACE, Literal: "}"},
	}

	LoadAndAssertTokens(t, "var_declaration.zen", expected)
}
