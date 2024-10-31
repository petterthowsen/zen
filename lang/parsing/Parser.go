package parsing

import (
	"strconv"
	"zen/lang/common"
	"zen/lang/lexing"
	"zen/lang/parsing/ast"
	"zen/lang/parsing/expression"
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
		if p.check(lexing.EOF) {
			break
		}

		stmt := p.parseStatement()
		if stmt != nil {
			statements = append(statements, stmt)
		} else if len(p.errors) > 0 && !p.stopAtFirstError {
			p.synchronize()
		}
	}

	if len(p.errors) > 0 {
		return nil, p.errors
	}

	return ast.NewProgramNode(statements), nil
}

// synchronize skips tokens until a statement boundary is found
func (p *Parser) synchronize() {
	for !p.isAtEnd() {
		if p.check(lexing.EOF) {
			return
		}

		// If we're at a statement-starting keyword, we can start parsing again
		if p.checkKeyword("var") || p.checkKeyword("const") ||
			p.checkKeyword("func") || p.checkKeyword("class") ||
			p.checkKeyword("if") || p.checkKeyword("for") ||
			p.checkKeyword("while") || p.checkKeyword("return") || p.checkKeyword("when") {
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
	var err *common.SyntaxError
	if token.Type == lexing.EOF {
		err = common.NewSyntaxError("Unexpected end of file: "+message, token.Location)
	} else {
		err = common.NewSyntaxError("Unexpected '"+token.Literal+"': "+message, token.Location)
	}
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

	token := p.peek()
	if token.Type == lexing.EOF {
		return nil // End of file is not an error
	}

	p.errorAtToken(token, "Expected statement")
	return nil
}

// parseVarDeclaration: Parses a variable declaration
func (p *Parser) parseVarDeclaration() ast.Statement {
	isConstant := p.previous().Literal == "const"
	startToken := p.previous() // Save the 'var' or 'const' token for error reporting

	// Parse variable name
	name := p.consume(lexing.IDENTIFIER, "Expected variable name")
	if len(p.errors) > 0 {
		return nil
	}

	isNullable := false

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
			isNullable = true
		}
	}

	// Parse optional initializer
	var initializer ast.Expression
	if p.match(lexing.ASSIGN) {
		initializer = p.parseExpression()
		if initializer == nil {
			return nil
		}
	}

	return statement.NewVarDeclarationNode(
		name.Literal,
		varType,
		initializer,
		isConstant,
		isNullable,
		startToken.Location,
	)
}

// parseExpression: Parses an expression
func (p *Parser) parseExpression() ast.Expression {
	return p.parsePrimary()
}

// parsePrimary: Parses primary expressions (literals for now)
func (p *Parser) parsePrimary() ast.Expression {
	token := p.peek()

	switch token.Type {
	case lexing.STRING:
		p.advance()
		return expression.NewLiteralExpression(token.Literal, token.Location)

	case lexing.INT:
		p.advance()
		// Convert string to int
		value, err := strconv.ParseInt(token.Literal, 10, 64)
		if err != nil {
			p.errorAtToken(token, "Invalid integer literal")
			return nil
		}
		return expression.NewLiteralExpression(value, token.Location)

	case lexing.KEYWORD:
		if token.Literal == "true" || token.Literal == "false" {
			p.advance()
			return expression.NewLiteralExpression(token.Literal == "true", token.Location)
		} else if token.Literal == "null" {
			p.advance()
			return expression.NewLiteralExpression(nil, token.Location)
		}
	}

	p.errorAtToken(token, "Expected expression")
	return nil
}
