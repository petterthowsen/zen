package lexing

import (
	"testing"
	"zen/lang/lexing"
)

/*
// single line
var numbers = [1, 2, 3]

// multiline
var strings = [

	"one"
	"two"
	"three"

]
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
	}

	LoadAndAssertTokens(t, "arrays.zen", expected)
}
