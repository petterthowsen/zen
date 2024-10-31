package parsing

import (
	"zen/lang/common"
	"zen/lang/lexing"
	"zen/lang/parsing/ast"
	"zen/lang/parsing/statement"
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
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(*common.SyntaxError); ok {
				// Expected panic from error() method
				return
			}
			// Unexpected panic, re-panic
			panic(r)
		}
	}()

	statements := make([]ast.Statement, 0)

	for !p.isAtEnd() {
		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
		}

		// Error recovery: skip until we find a statement boundary
		if len(p.errors) > 0 && !p.stopAtFirstError {
			p.synchronize()
		}
	}

	// Return nil if we had errors
	if len(p.errors) > 0 {
		return nil, p.errors
	}

	return ast.NewProgramNode(statements), p.errors
}

// synchronize skips tokens until a statement boundary is found
func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		// Synchronize on statement-starting keywords
		if p.checkKeyword("var") || p.checkKeyword("const") ||
			p.checkKeyword("func") || p.checkKeyword("class") ||
			p.checkKeyword("if") || p.checkKeyword("for") ||
			p.checkKeyword("while") || p.checkKeyword("return") {
			return
		}
		p.advance()
	}
}

// isAtEnd returns true if we've reached the end of the tokens
func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens)
}

// peek returns the current token
func (p *Parser) peek() lexing.Token {
	if p.current >= len(p.tokens) {
		return lexing.Token{} // Return empty token
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
		return false
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

// match returns true if the current token type matches any of the given TokenType
func (p *Parser) match(types ...lexing.TokenType) bool {
	for _, typ := range types {
		if p.check(typ) {
			p.advance()
			return true
		}
	}
	return false
}

// matchKeyword returns true if the current token is a keyword with any of the given literals
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

	p.error(message)
	return lexing.Token{}
}

// error adds a SyntaxError to the errors array
func (p *Parser) error(message string) {
	err := common.NewSyntaxError(message, p.peek().Location)
	p.errors = append(p.errors, err)

	if p.stopAtFirstError {
		panic(err) // Will be caught in Parse()
	}
}

// parseStatement: Initial statement parsing - we'll expand this as we implement more features
func (p *Parser) parseStatement() ast.Statement {
	if p.matchKeyword("var", "const") {
		return p.parseVarDeclaration()
	}

	// TODO: Add other statement types
	p.error("Expected statement")
	return nil
}

// parseVarDeclaration: Parses a variable declaration
func (p *Parser) parseVarDeclaration() ast.Statement {
	isConstant := p.previous().Literal == "const"

	// Parse variable name
	name := p.consume(lexing.IDENTIFIER, "Expected variable name")
	if len(p.errors) > 0 {
		return nil
	}

	// Parse optional type annotation
	var varType string
	if p.match(lexing.COLON) {
		typeToken := p.consume(lexing.IDENTIFIER, "Expected type name")
		if len(p.errors) > 0 {
			return nil
		}
		varType = typeToken.Literal

		// Handle nullable type
		if p.match(lexing.QMARK) {
			varType += "?"
		}
	}

	// Parse initializer
	p.consume(lexing.ASSIGN, "Expected '=' after variable name")
	if len(p.errors) > 0 {
		return nil
	}

	// TODO: Implement expression parsing
	// For now, just consume tokens until semicolon
	for !p.isAtEnd() && !p.check(lexing.SEMICOLON) {
		p.advance()
	}

	if !p.isAtEnd() {
		p.advance() // Consume semicolon
	}

	return statement.NewVarDeclarationNode(
		name.Literal,
		varType,
		nil, // TODO: Add initializer expression
		isConstant,
		name.Location,
	)
}
