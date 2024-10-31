package lexing

import (
	"fmt"
	"strings"
	"unicode"
	"zen/lang/common"
)

// Lexer represents a lexical analyzer for tokenizing source code.
// SourceCode is the input string being parsed.
// Index is the current position in the SourceCode.
// Line is the current line number in the SourceCode.
// Column is the current column number in the SourceCode.
type Lexer struct {
	SourceCode common.SourceCode
	Index      int
	Line       int
	Column     int
	Errors     []common.SyntaxError
	tokens     []Token
}

func NewLexer(SourceCode common.SourceCode) *Lexer {
	return &Lexer{
		SourceCode: SourceCode,
	}
}

func (l *Lexer) getLocation() *common.SourceLocation {
	return l.SourceCode.GetLocation(l.Line, l.Column)
}

// Scan tokenizes the source code stored in the Lexer and returns a slice of Tokens and an error, if any.
func (l *Lexer) Scan() ([]Token, error) {
	l.tokens = []Token{}
	l.Index = 0
	l.Line = 1
	l.Column = 0
	l.Errors = []common.SyntaxError{}

	for l.Index <= l.SourceCode.GetLength() {
		ch := l.Peek()

		switch {
		case l.IsEOF():
			eofToken := Token{
				Type:     EOF,
				Literal:  "",
				Location: l.SourceCode.GetLocation(l.Line, l.Column),
			}
			l.tokens = append(l.tokens, eofToken)
			l.Consume()
		case string(ch) == "/" && string(l.Next()) == "/":
			l.ConsumeAllExcept('\n')
		case string(ch) == "\"":
			l.tokens = append(l.tokens, l.scanString())
		case unicode.IsLetter(ch):
			l.tokens = append(l.tokens, l.scanIdentifierOrKeyword())
		case unicode.IsDigit(ch):
			l.tokens = append(l.tokens, l.scanNumber())
		case l.isSequence("++"):
			l.scanSequence("++", INCREMENT)
		case l.isSequence("--"):
			l.scanSequence("--", DECREMENT)
		case l.isSequence("+="):
			l.scanSequence("+=", PLUS_ASSIGN)
		case l.isSequence("-="):
			l.scanSequence("-=", MINUS_ASSIGN)
		case l.isSequence("*="):
			l.scanSequence("*=", MULTIPLY_ASSIGN)
		case l.isSequence("/="):
			l.scanSequence("/=", DIVIDE_ASSIGN)
		case string(ch) == "+":
			l.ConsumeToken(PLUS)
		case string(ch) == "-":
			l.ConsumeToken(MINUS)
		case string(ch) == "*":
			l.ConsumeToken(MULTIPLY)
		case string(ch) == "/":
			l.ConsumeToken(DIVIDE)
		case l.isSequence("=="):
			l.scanSequence("==", EQUALS)
		case l.isSequence("!="):
			l.scanSequence("!=", NOT_EQUALS)
		case l.isSequence(">="):
			l.scanSequence(">=", GREATER_EQUALS)
		case l.isSequence("<="):
			l.scanSequence("<=", LESS_EQUALS)
		case l.isSequence(">"):
			l.scanSequence(">", GREATER)
		case l.isSequence("<"):
			l.scanSequence("<", LESS)
		case string(ch) == ",":
			l.ConsumeToken(COMMA)
		case string(ch) == ";":
			l.ConsumeToken(SEMICOLON)
		case string(ch) == "(":
			l.ConsumeToken(LEFT_PAREN)
		case string(ch) == ")":
			l.ConsumeToken(RIGHT_PAREN)
		case string(ch) == "{":
			l.ConsumeToken(LEFT_BRACE)
		case string(ch) == "}":
			l.ConsumeToken(RIGHT_BRACE)
		case string(ch) == "[":
			l.ConsumeToken(LEFT_BRACKET)
		case string(ch) == "]":
			l.ConsumeToken(RIGHT_BRACKET)
		case string(ch) == ":":
			l.ConsumeToken(COLON)
		case string(ch) == ".":
			l.ConsumeToken(DOT)
		case string(ch) == "=":
			l.ConsumeToken(ASSIGN)
		case string(ch) == "?":
			l.ConsumeToken(QMARK)
		case unicode.IsSpace(ch):
			l.IgnoreWhitespace()
		default:
			l.addError("Unexpected character")
			l.Consume()
		}
	}

	if len(l.Errors) > 0 {
		return l.tokens, fmt.Errorf("syntax error(s) found in source code")
	}

	return l.tokens, nil
}

func (l *Lexer) scanString() Token {
	start := l.Index

	// consume the opening quote
	l.Consume()

	for {
		ch := l.Peek()

		// Check for end of file
		if l.IsEOF() {
			l.addError("Unterminated string literal")
			break
		}

		// Check for escaped quote
		if ch == '\\' {
			l.Consume() // Consume the backslash
			nextCh := l.Peek()
			if nextCh == '"' || nextCh == '\\' || nextCh == 'n' || nextCh == 't' {
				// Consume the escaped character
				l.Consume()
				continue
			} else {
				l.addError(fmt.Sprintf("Invalid escape sequence: \\%c", nextCh))
				break
			}
		} else if ch == '"' {
			// Closing quote found
			l.Consume()
			break
		} else if ch == '\n' {
			l.addError("Unexpected newline in string literal")
			break
		}

		// Consume regular characters
		l.Consume()
	}

	unquoted := l.SourceCode.GetText()[start+1 : l.Index-1]

	// unescape the \" sequences
	unquoted = strings.Replace(unquoted, "\\\"", "", -1)

	return Token{
		Type:     STRING,
		Literal:  unquoted,
		Location: l.SourceCode.GetLocation(l.Line, l.Column),
	}
}

func (l *Lexer) addError(message string) {
	l.Errors = append(l.Errors, common.SyntaxError{
		Message:  message,
		Location: l.SourceCode.GetLocation(l.Line, l.Column),
	})
}

func (l *Lexer) scanNumber() Token {
	start := l.Index

	for {
		if unicode.IsDigit(l.Peek()) {
			l.Consume()
		} else if l.Previous() != '.' && l.Peek() == '.' {
			l.Consume()
		} else {
			break
		}
	}

	text := l.SourceCode.GetText()[start:l.Index]

	// check if it contains a .
	if strings.Contains(text, ".") {
		return Token{
			Type:     FLOAT,
			Literal:  text,
			Location: l.SourceCode.GetLocation(l.Line, l.Column),
		}
	} else {
		return Token{
			Type:     INT,
			Literal:  text,
			Location: l.SourceCode.GetLocation(l.Line, l.Column),
		}
	}
}

// scanIdentifierOrKeyword scans a sequence starting with a letter,
// followed by letters or digits, and determines if it's an identifier or keyword.
func (l *Lexer) scanIdentifierOrKeyword() Token {
	start := l.Index

	// is letter, digit or underscore?
	for unicode.IsLetter(l.Peek()) || unicode.IsDigit(l.Peek()) || l.Peek() == '_' {
		l.Consume()
	}

	text := l.SourceCode.GetText()[start:l.Index]

	// keyword or identifier?
	if IsKeyword(text) {
		return Token{
			Type:     KEYWORD,
			Literal:  text,
			Location: l.SourceCode.GetLocation(l.Line, l.Column),
		}
	} else {
		return Token{
			Type:     IDENTIFIER,
			Literal:  text,
			Location: l.SourceCode.GetLocation(l.Line, l.Column),
		}
	}
}

func (l *Lexer) scanSequence(sequence string, tokenType TokenType) {
	for i := 0; i < len(sequence); i++ {
		if l.Peek() != rune(sequence[i]) {
			l.addError("scanSequence failed: " + sequence)
		}
		l.Consume() // Consume matching characters
	}

	l.tokens = append(l.tokens, Token{
		Type:     tokenType,
		Literal:  sequence,
		Location: l.SourceCode.GetLocation(l.Line, l.Column),
	})
}

// Peek returns the current rune at the lexer's index without advancing the position.
// If the index is at or beyond the end of the source code, it returns 0.
func (l *Lexer) Peek() rune {
	if l.Index >= l.SourceCode.GetLength() {
		return 0
	}
	return l.SourceCode.GetChar(l.Index)
}

func (l *Lexer) isSequence(sequence string) bool {
	idx := l.Index

	for _, char := range sequence {
		if l.SourceCode.GetChar(idx) != char {
			return false
		}
		idx += 1
	}

	return true
}

func (l *Lexer) consumeSequence(sequence string) {
	for i := 0; i < len(sequence); i++ {
		if l.Peek() != rune(sequence[i]) {
			return // Exit if characters don't match
		}
		l.Consume() // Consume matching characters
	}
}

func (l *Lexer) Previous() rune {
	prev := l.Index - 1
	if prev < 0 || prev >= l.SourceCode.GetLength() {
		return 0
	}

	return l.SourceCode.GetChar(prev)
}

// Next returns the next rune in the source code. Returns 0 if at the end of the source code.
func (l *Lexer) Next() rune {
	if l.Index+1 >= l.SourceCode.GetLength() {
		return 0
	}
	return l.SourceCode.GetChar(l.Index + 1)
}

// Consume reads the next rune from the source code, updates the lexer's position, and returns the read rune.
// Advances the Index, and if a newline is encountered, increments the Line and resets Column to 0, otherwise increments Column.
func (l *Lexer) Consume() rune {
	ch := l.Peek()
	l.Index++
	if ch == '\n' {
		l.Line++
		l.Column = 0
	} else {
		l.Column++
	}
	return ch
}

func (l *Lexer) ConsumeAll(character rune) {
	for l.Peek() == character {
		l.Consume()
	}
}

func (l *Lexer) ConsumeAllExcept(character rune) {
	for l.Peek() != character {
		l.Consume()
	}
}

// ConsumeToken consumes the character and creates a new Token.
func (l *Lexer) ConsumeToken(tokenType TokenType) {
	l.tokens = append(l.tokens, Token{
		Type:     tokenType,
		Literal:  string(l.Peek()),
		Location: l.SourceCode.GetLocation(l.Line, l.Column),
	})
	l.Consume()
}

// IsEOF checks if the Lexer has reached or exceeded the end of the SourceCode.
func (l *Lexer) IsEOF() bool {
	return l.Index >= l.SourceCode.GetLength()
}

// IgnoreWhitespace consumes all contiguous whitespace characters from the Lexer's current position.
func (l *Lexer) IgnoreWhitespace() {
	for unicode.IsSpace(l.Peek()) && !l.IsEOF() {
		l.Consume()
	}
}
