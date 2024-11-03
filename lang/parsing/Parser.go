package parsing

import (
	"zen/lang/common"
	"zen/lang/lexing"
	"zen/lang/parsing/ast"
)

type Parser struct {
	tokens           []lexing.Token
	current          int
	stopAtFirstError bool
	errors           []*common.SyntaxError
}

// NewParser creates a new Parser instance
func NewParser(tokens []lexing.Token, stopAtFirstError bool) *Parser {
	return &Parser{
		tokens:           tokens,
		current:          0,
		stopAtFirstError: stopAtFirstError,
		errors:           make([]*common.SyntaxError, 0),
	}
}

// Parse takes an array of tokens and produces an AST with a ProgramNode as the root node.
func (p *Parser) Parse() (*ast.ProgramNode, []*common.SyntaxError) {
	//defer func() {
	//	if r := recover(); r != nil {
	//		if _, ok := r.(*common.SyntaxError); ok {
	//			// Expected panic from error() method
	//			return
	//		}
	//		// Unexpected panic, re-panic
	//		panic(r)
	//	}
	//}()

	statements := make([]ast.Statement, 0)

	for !p.isAtEnd() {
		if p.check(lexing.EOF) {
			break
		}

		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
		} else if len(p.errors) > 0 && !p.stopAtFirstError {
			if !p.synchronize() {
				break
			}
		}
	}

	programNode := ast.NewProgramNode(statements)

	// return with any errors
	if len(p.errors) > 0 {
		return programNode, p.errors
	} else {
		return programNode, nil
	}
}

// synchronize skips tokens until a statement boundary is found
// returns true if successful, false if not
func (p *Parser) synchronize() bool {
	for !p.isAtEnd() {
		if p.check(lexing.EOF) {
			return false
		}

		// If we're at a statement-starting keyword, we can start parsing again
		if p.checkKeyword("var") || p.checkKeyword("const") ||
			p.checkKeyword("func") || p.checkKeyword("class") ||
			p.checkKeyword("if") || p.checkKeyword("for") ||
			p.checkKeyword("while") || p.checkKeyword("return") || p.checkKeyword("when") {
			return true // Found a synchronization point
		}

		p.advance() // Skip tokens until we find a statement boundary
	}

	return false
}

// isAtEnd returns true if we've reached the end of the tokens
func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens)
}

// peek returns the current token
func (p *Parser) peek() lexing.Token {
	if p.current >= len(p.tokens) {
		return lexing.Token{Type: lexing.EOF} // Return EOF token
	}
	return p.tokens[p.current]
}

// previous returns the previous token
func (p *Parser) previous() lexing.Token {
	return p.tokens[p.current-1]
}

// advance returns the current token and advances to the next
func (p *Parser) advance() lexing.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

// check returns true if the current token type matches the given TokenType
func (p *Parser) check(typ lexing.TokenType) bool {
	if p.isAtEnd() {
		return typ == lexing.EOF
	}
	return p.peek().Type == typ
}

// checkKeyword returns true if the current token is a keyword with the given literal
func (p *Parser) checkKeyword(keyword string) bool {
	if p.isAtEnd() {
		return false
	}
	token := p.peek()
	return token.Type == lexing.KEYWORD && token.Literal == keyword
}

// match checks and consumes (advances) the current token if it matches any of the given types
func (p *Parser) match(types ...lexing.TokenType) bool {
	for _, typ := range types {
		if p.check(typ) {
			p.advance()
			return true
		}
	}
	return false
}

// matchKeyword returns true and consumes if the current token is a keyword with any of the given literals
func (p *Parser) matchKeyword(keywords ...string) bool {
	for _, keyword := range keywords {
		if p.checkKeyword(keyword) {
			p.advance()
			return true
		}
	}
	return false
}

// consume advances if the current token matches the expected type, otherwise reports an error
func (p *Parser) consume(typ lexing.TokenType, message string) lexing.Token {
	if p.check(typ) {
		return p.advance()
	}

	token := p.peek()
	p.errorAtToken(token, message)
	return lexing.Token{}
}

// error adds a SyntaxError to the errors array
func (p *Parser) error(message string) {
	p.errorAtToken(p.peek(), message)
}

// errorAtToken adds a SyntaxError at the specified token
func (p *Parser) errorAtToken(token lexing.Token, message string) {
	err := common.NewSyntaxError(message, token.Location)
	p.errors = append(p.errors, err)

	if p.stopAtFirstError {
		panic(err) // Will be caught in Parse()
	}
}

// parseBlock parses a block of statements until it encounters a right brace or the end of the input.
func (p *Parser) parseBlock() []ast.Statement {
	body := make([]ast.Statement, 0)

	// Check for premature end of input
	if p.isAtEnd() {
		p.error("Unexpected end of input while parsing block")
		return body
	}

	// Keep parsing statements until we hit a right brace or EOF or no statements
	for !p.check(lexing.RIGHT_BRACE) {
		if p.isAtEnd() {
			p.error("Unterminated block - expected '}'")
			return body
		}

		stmt := p.parseStatement()
		if stmt == nil {
			break
		}
		body = append(body, stmt)
	}

	// Note: We don't consume the right brace here because that should be done by the calling method
	// This allows the calling method to handle any syntax after the block
	return body
}
