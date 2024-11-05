package lexing

import (
	"testing"
	"zen/lang/lexing"
)

/*
var numbers = [1, 2, 3]

var strings = [

	"one"
	"two"
	"three"

]

numbers[0] = 2

strings[1] = "one edited"
*/
func TestArrays(t *testing.T) {
	expected := []TokenAssert{

		{Type: lexing.KEYWORD, Literal: "var"},
		{Type: lexing.IDENTIFIER, Literal: "numbers"},
		{Type: lexing.ASSIGN, Literal: "="},
		{Type: lexing.LEFT_BRACKET, Literal: "["},

		{Type: lexing.INT, Literal: "1"},
		{Type: lexing.COMMA, Literal: ","},
		{Type: lexing.INT, Literal: "2"},
		{Type: lexing.COMMA, Literal: ","},
		{Type: lexing.INT, Literal: "3"},

		{Type: lexing.RIGHT_BRACKET, Literal: "]"},

		{Type: lexing.KEYWORD, Literal: "var"},
		{Type: lexing.IDENTIFIER, Literal: "strings"},
		{Type: lexing.ASSIGN, Literal: "="},
		{Type: lexing.LEFT_BRACKET, Literal: "["},

		{Type: lexing.STRING, Literal: "one"},
		{Type: lexing.STRING, Literal: "two"},
		{Type: lexing.STRING, Literal: "three"},

		{Type: lexing.RIGHT_BRACKET, Literal: "]"},

		// numbers[0] = 2
		{Type: lexing.IDENTIFIER, Literal: "numbers"},
		{Type: lexing.LEFT_BRACKET, Literal: "["},
		{Type: lexing.INT, Literal: "0"},

		{Type: lexing.RIGHT_BRACKET, Literal: "]"},

		{Type: lexing.ASSIGN, Literal: "="},

		{Type: lexing.INT, Literal: "2"},

		// strings[1] = "one edited"
		{Type: lexing.IDENTIFIER, Literal: "strings"},
		{Type: lexing.LEFT_BRACKET, Literal: "["},
		{Type: lexing.INT, Literal: "1"},
		{Type: lexing.RIGHT_BRACKET, Literal: "]"},
		{Type: lexing.ASSIGN, Literal: "="},
		{Type: lexing.STRING, Literal: "one edited"},
	}

	LoadAndAssertTokens(t, "arrays.zen", expected)
}
