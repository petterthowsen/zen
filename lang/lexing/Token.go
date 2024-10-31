package lexing

import "zen/lang/common"

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF

	IDENTIFIER //myVar
	KEYWORD    // if, else, for etc.

	INT    // 123
	FLOAT  // 3.14
	STRING // "hello, world"

	DOT
	COMMA
	COLON
	SEMICOLON
	QMARK

	LEFT_PAREN
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	LEFT_BRACKET
	RIGHT_BRACKET

	LESS
	GREATER

	//OPERATORS
	PLUS
	MINUS
	MULTIPLY
	DIVIDE
	PERCENT
	EQUALS
	NOT_EQUALS
	GREATER_EQUALS
	LESS_EQUALS

	INCREMENT
	DECREMENT

	ASSIGN
	PLUS_ASSIGN
	MINUS_ASSIGN
	MULTIPLY_ASSIGN
	DIVIDE_ASSIGN
)

var tokenTypeNames = map[TokenType]string{
	ILLEGAL: "Illegal",
	EOF:     "EOF",

	IDENTIFIER: "Identifier", //main, car.drive
	KEYWORD:    "Keyword",    // if, else, for etc.

	INT:    "Int",
	FLOAT:  "Float",
	STRING: "String",

	DOT:       "Dot",
	COMMA:     "Comma",
	COLON:     "Colon",
	SEMICOLON: "Semicolon",
	QMARK:     "QuestionMark",

	LEFT_PAREN:    "LeftParen",
	RIGHT_PAREN:   "RightParen",
	LEFT_BRACE:    "LeftBrace",
	RIGHT_BRACE:   "RightBrace",
	LEFT_BRACKET:  "LeftBracket",
	RIGHT_BRACKET: "RightBracket",

	LESS:    "Less",
	GREATER: "Greater",

	//OPERATORS
	PLUS:           "Plus",
	MINUS:          "Minus",
	MULTIPLY:       "Multiply",
	DIVIDE:         "Divide",
	PERCENT:        "Percent",
	EQUALS:         "Equals",
	NOT_EQUALS:     "NotEqual",
	GREATER_EQUALS: "GreaterEqual",
	LESS_EQUALS:    "LessEqual",
	INCREMENT:      "Increment",
	DECREMENT:      "Decrement",

	ASSIGN:          "Assign",
	PLUS_ASSIGN:     "PlusAssign",
	MINUS_ASSIGN:    "MinusAssign",
	MULTIPLY_ASSIGN: "MultiplyAssign",
	DIVIDE_ASSIGN:   "DivideAssign",
}

var keywords = []string{
	"import",
	"package",

	"var",
	"const",

	"if",
	"else",

	"for",
	"in",
	"while",
	"when",
	"where",
	"break",
	"continue",

	"return",

	"func",
	"class",
	"interface",
	"implements",
	"extends",
	"new",
	"this",
	"super",
	"pub",

	"true",
	"false",
	"null",
}

type Token struct {
	Type     TokenType
	Literal  string
	Location *common.SourceLocation
}

// NewToken creates a new Token instance with the given type, literal, and source location.
func NewToken(tokenType TokenType, literal string, sourceLocation *common.SourceLocation) *Token {
	return &Token{
		Type:     tokenType,
		Literal:  literal,
		Location: sourceLocation,
	}
}

// Name returns the name of the token type
func (t *Token) Name() string {
	return tokenTypeNames[t.Type]
}

// TokenName returns the string name of the token type
func TokenName(t TokenType) string {
	return tokenTypeNames[t]
}

// String returns the string representation of the Token.
func (t *Token) String() string {
	var str string = t.Name()
	if t.Literal != "" {
		str += "(" + t.Literal + ")"
	}
	return str
}

func IsKeyword(text string) bool {
	for _, keyword := range keywords {
		if keyword == text {
			return true
		}
	}

	return false
}
