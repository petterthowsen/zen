package lexing

import (
	"testing"
	"zen/lang/lexing"
)

/*
// single-line
var config = {"title": "hello", "volume": 0.5}

// multi-line

	var config = {
	    "title": "hello"
	    "volume": 0.5
	}
*/
func TestMaps(t *testing.T) {
	expected := []TokenAssert{

		// one line
		{Type: lexing.KEYWORD, Literal: "var"},
		{Type: lexing.IDENTIFIER, Literal: "config"},
		{Type: lexing.ASSIGN, Literal: "="},
		{Type: lexing.LEFT_BRACE, Literal: "{"},

		{Type: lexing.STRING, Literal: "title"},
		{Type: lexing.COLON, Literal: ":"},
		{Type: lexing.STRING, Literal: "hello"},
		{Type: lexing.COMMA, Literal: ","},

		{Type: lexing.STRING, Literal: "volume"},
		{Type: lexing.COLON, Literal: ":"},
		{Type: lexing.FLOAT, Literal: "0.5"},

		{Type: lexing.RIGHT_BRACE, Literal: "}"},

		// multi-line without comma
		{Type: lexing.KEYWORD, Literal: "var"},
		{Type: lexing.IDENTIFIER, Literal: "config"},
		{Type: lexing.ASSIGN, Literal: "="},
		{Type: lexing.LEFT_BRACE, Literal: "{"},

		{Type: lexing.STRING, Literal: "title"},
		{Type: lexing.COLON, Literal: ":"},
		{Type: lexing.STRING, Literal: "hello"},

		{Type: lexing.STRING, Literal: "volume"},
		{Type: lexing.COLON, Literal: ":"},
		{Type: lexing.FLOAT, Literal: "0.5"},

		{Type: lexing.RIGHT_BRACE, Literal: "}"},
	}

	LoadAndAssertTokens(t, "maps.zen", expected)
}
