package lexing

import (
	"testing"
	"zen/lang/lexing"
)

func TestControlFlow(t *testing.T) {
	expected := []TokenAssert{
		// if x > 0 { print("positive") }
		{Type: lexing.KEYWORD, Literal: "if"},
		{Type: lexing.IDENTIFIER, Literal: "x"},
		{Type: lexing.GREATER, Literal: ">"},
		{Type: lexing.INT, Literal: "0"},
		{Type: lexing.LEFT_BRACE, Literal: "{"},
		{Type: lexing.IDENTIFIER, Literal: "print"},
		{Type: lexing.LEFT_PAREN, Literal: "("},
		{Type: lexing.STRING, Literal: "positive"},
		{Type: lexing.RIGHT_PAREN, Literal: ")"},
		{Type: lexing.RIGHT_BRACE, Literal: "}"},

		// else if x < 0 { print("negative") }
		{Type: lexing.KEYWORD, Literal: "else"},
		{Type: lexing.KEYWORD, Literal: "if"},
		{Type: lexing.IDENTIFIER, Literal: "x"},
		{Type: lexing.LESS, Literal: "<"},
		{Type: lexing.INT, Literal: "0"},
		{Type: lexing.LEFT_BRACE, Literal: "{"},
		{Type: lexing.IDENTIFIER, Literal: "print"},
		{Type: lexing.LEFT_PAREN, Literal: "("},
		{Type: lexing.STRING, Literal: "negative"},
		{Type: lexing.RIGHT_PAREN, Literal: ")"},
		{Type: lexing.RIGHT_BRACE, Literal: "}"},

		// else { print("zero") }
		{Type: lexing.KEYWORD, Literal: "else"},
		{Type: lexing.LEFT_BRACE, Literal: "{"},
		{Type: lexing.IDENTIFIER, Literal: "print"},
		{Type: lexing.LEFT_PAREN, Literal: "("},
		{Type: lexing.STRING, Literal: "zero"},
		{Type: lexing.RIGHT_PAREN, Literal: ")"},
		{Type: lexing.RIGHT_BRACE, Literal: "}"},

		// while count > 0 { count-- }
		{Type: lexing.KEYWORD, Literal: "while"},
		{Type: lexing.IDENTIFIER, Literal: "count"},
		{Type: lexing.GREATER, Literal: ">"},
		{Type: lexing.INT, Literal: "0"},
		{Type: lexing.LEFT_BRACE, Literal: "{"},
		{Type: lexing.IDENTIFIER, Literal: "count"},
		{Type: lexing.DECREMENT, Literal: "--"},
		{Type: lexing.RIGHT_BRACE, Literal: "}"},

		// for var i = 0; i < 10; i++ { print(i) }
		{Type: lexing.KEYWORD, Literal: "for"},
		{Type: lexing.KEYWORD, Literal: "var"},
		{Type: lexing.IDENTIFIER, Literal: "i"},
		{Type: lexing.ASSIGN, Literal: "="},
		{Type: lexing.INT, Literal: "0"},
		{Type: lexing.SEMICOLON, Literal: ";"},
		{Type: lexing.IDENTIFIER, Literal: "i"},
		{Type: lexing.LESS, Literal: "<"},
		{Type: lexing.INT, Literal: "10"},
		{Type: lexing.SEMICOLON, Literal: ";"},
		{Type: lexing.IDENTIFIER, Literal: "i"},
		{Type: lexing.INCREMENT, Literal: "++"},
		{Type: lexing.LEFT_BRACE, Literal: "{"},
		{Type: lexing.IDENTIFIER, Literal: "print"},
		{Type: lexing.LEFT_PAREN, Literal: "("},
		{Type: lexing.IDENTIFIER, Literal: "i"},
		{Type: lexing.RIGHT_PAREN, Literal: ")"},
		{Type: lexing.RIGHT_BRACE, Literal: "}"},

		// for number in [1, 5, 7] { print(number) }
		{Type: lexing.KEYWORD, Literal: "for"},
		{Type: lexing.IDENTIFIER, Literal: "number"},
		{Type: lexing.KEYWORD, Literal: "in"},
		{Type: lexing.LEFT_BRACKET, Literal: "["},
		{Type: lexing.INT, Literal: "1"},
		{Type: lexing.COMMA, Literal: ","},
		{Type: lexing.INT, Literal: "5"},
		{Type: lexing.COMMA, Literal: ","},
		{Type: lexing.INT, Literal: "7"},
		{Type: lexing.RIGHT_BRACKET, Literal: "]"},
		{Type: lexing.LEFT_BRACE, Literal: "{"},
		{Type: lexing.IDENTIFIER, Literal: "print"},
		{Type: lexing.LEFT_PAREN, Literal: "("},
		{Type: lexing.IDENTIFIER, Literal: "number"},
		{Type: lexing.RIGHT_PAREN, Literal: ")"},
		{Type: lexing.RIGHT_BRACE, Literal: "}"},

		// for index, value in items { ... }
		{Type: lexing.KEYWORD, Literal: "for"},
		{Type: lexing.IDENTIFIER, Literal: "index"},
		{Type: lexing.COMMA, Literal: ","},
		{Type: lexing.IDENTIFIER, Literal: "value"},
		{Type: lexing.KEYWORD, Literal: "in"},
		{Type: lexing.IDENTIFIER, Literal: "items"},
		{Type: lexing.LEFT_BRACE, Literal: "{"},
		{Type: lexing.IDENTIFIER, Literal: "print"},
		{Type: lexing.LEFT_PAREN, Literal: "("},
		{Type: lexing.IDENTIFIER, Literal: "index"},
		{Type: lexing.RIGHT_PAREN, Literal: ")"},
		{Type: lexing.IDENTIFIER, Literal: "print"},
		{Type: lexing.LEFT_PAREN, Literal: "("},
		{Type: lexing.IDENTIFIER, Literal: "value"},
		{Type: lexing.RIGHT_PAREN, Literal: ")"},
		{Type: lexing.RIGHT_BRACE, Literal: "}"},

		// for article in articles where article.views > 100 { ... }
		{Type: lexing.KEYWORD, Literal: "for"},
		{Type: lexing.IDENTIFIER, Literal: "article"},
		{Type: lexing.KEYWORD, Literal: "in"},
		{Type: lexing.IDENTIFIER, Literal: "articles"},
		{Type: lexing.KEYWORD, Literal: "where"},
		{Type: lexing.IDENTIFIER, Literal: "article"},
		{Type: lexing.DOT, Literal: "."},
		{Type: lexing.IDENTIFIER, Literal: "views"},
		{Type: lexing.GREATER, Literal: ">"},
		{Type: lexing.INT, Literal: "100"},
		{Type: lexing.LEFT_BRACE, Literal: "{"},
		{Type: lexing.IDENTIFIER, Literal: "print"},
		{Type: lexing.LEFT_PAREN, Literal: "("},
		{Type: lexing.IDENTIFIER, Literal: "article"},
		{Type: lexing.DOT, Literal: "."},
		{Type: lexing.IDENTIFIER, Literal: "title"},
		{Type: lexing.RIGHT_PAREN, Literal: ")"},
		{Type: lexing.RIGHT_BRACE, Literal: "}"},

		// when value { ... }
		{Type: lexing.KEYWORD, Literal: "when"},
		{Type: lexing.IDENTIFIER, Literal: "value"},
		{Type: lexing.LEFT_BRACE, Literal: "{"},
		{Type: lexing.STRING, Literal: "john"},
		{Type: lexing.LEFT_BRACE, Literal: "{"},
		{Type: lexing.IDENTIFIER, Literal: "print"},
		{Type: lexing.LEFT_PAREN, Literal: "("},
		{Type: lexing.STRING, Literal: "name is john!"},
		{Type: lexing.RIGHT_PAREN, Literal: ")"},
		{Type: lexing.RIGHT_BRACE, Literal: "}"},
		{Type: lexing.STRING, Literal: "jane"},
		{Type: lexing.LEFT_BRACE, Literal: "{"},
		{Type: lexing.IDENTIFIER, Literal: "print"},
		{Type: lexing.LEFT_PAREN, Literal: "("},
		{Type: lexing.STRING, Literal: "name is jane!"},
		{Type: lexing.RIGHT_PAREN, Literal: ")"},
		{Type: lexing.RIGHT_BRACE, Literal: "}"},
		{Type: lexing.KEYWORD, Literal: "else"},
		{Type: lexing.LEFT_BRACE, Literal: "{"},
		{Type: lexing.IDENTIFIER, Literal: "print"},
		{Type: lexing.LEFT_PAREN, Literal: "("},
		{Type: lexing.STRING, Literal: "unknown name"},
		{Type: lexing.RIGHT_PAREN, Literal: ")"},
		{Type: lexing.RIGHT_BRACE, Literal: "}"},
		{Type: lexing.RIGHT_BRACE, Literal: "}"},

		// when age { ... }
		{Type: lexing.KEYWORD, Literal: "when"},
		{Type: lexing.IDENTIFIER, Literal: "age"},
		{Type: lexing.LEFT_BRACE, Literal: "{"},
		{Type: lexing.LESS, Literal: "<"},
		{Type: lexing.INT, Literal: "13"},
		{Type: lexing.LEFT_BRACE, Literal: "{"},
		{Type: lexing.IDENTIFIER, Literal: "print"},
		{Type: lexing.LEFT_PAREN, Literal: "("},
		{Type: lexing.STRING, Literal: "child"},
		{Type: lexing.RIGHT_PAREN, Literal: ")"},
		{Type: lexing.RIGHT_BRACE, Literal: "}"},
		{Type: lexing.LESS, Literal: "<"},
		{Type: lexing.INT, Literal: "20"},
		{Type: lexing.LEFT_BRACE, Literal: "{"},
		{Type: lexing.IDENTIFIER, Literal: "print"},
		{Type: lexing.LEFT_PAREN, Literal: "("},
		{Type: lexing.STRING, Literal: "teenager"},
		{Type: lexing.RIGHT_PAREN, Literal: ")"},
		{Type: lexing.RIGHT_BRACE, Literal: "}"},
		{Type: lexing.KEYWORD, Literal: "else"},
		{Type: lexing.LEFT_BRACE, Literal: "{"},
		{Type: lexing.IDENTIFIER, Literal: "print"},
		{Type: lexing.LEFT_PAREN, Literal: "("},
		{Type: lexing.STRING, Literal: "adult"},
		{Type: lexing.RIGHT_PAREN, Literal: ")"},
		{Type: lexing.RIGHT_BRACE, Literal: "}"},
		{Type: lexing.RIGHT_BRACE, Literal: "}"},
	}

	LoadAndAssertTokens(t, "control_flow.zen", expected)
}
