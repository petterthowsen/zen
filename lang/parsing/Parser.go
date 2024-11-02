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
		}

		p.advance()
		return true
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

// parseStatement: Initial statement parsing - we'll expand this as we implement more features
func (p *Parser) parseStatement() ast.Statement {
	// var/const declaration
	if p.matchKeyword("var", "const") {
		return p.parseVarDeclaration()
	}

	// func declaration
	if p.matchKeyword("func") {
		return p.parseFuncDeclaration()
	}

	// If Statement
	if p.matchKeyword("if") {
		return p.parseIfStatement()
	}

	// Return statement
	if p.matchKeyword("return") {
		return p.parseReturnStatement()
	}

	// Try parsing an expression statement
	expr := p.parseExpression()
	if expr != nil {
		stmt := &statement.ExpressionStatement{
			Location:   p.previous().Location,
			Expression: expr,
		}
		return stmt
	}

	token := p.peek()
	if token.Type == lexing.EOF {
		return nil // End of file is not an error
	}

	p.errorAtToken(token, "Expected statement")
	return nil
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

// parseIfStatement parses an if statement
func (p *Parser) parseIfStatement() ast.Statement {
	startToken := p.previous() // The 'if' token

	// parse the first IfConditionBlock explicitly since it is required
	primaryCondition := p.parseExpression()
	if primaryCondition == nil {
		p.errorAtToken(startToken, "Expected expression after 'if'")
		return nil
	}

	// parse the primary block
	if !p.match(lexing.LEFT_BRACE) {
		p.error("Expected '{' after 'if'")
		return nil
	}

	primaryBlock := p.parseBlock()

	if !p.match(lexing.RIGHT_BRACE) {
		p.error("Expected '}' after if block")
		return nil
	}

	// parse optional 'elif' condition blocks

	elifBlocks := make([]*statement.IfConditionBlock, 0)

	for {
		tok := p.peek()
		if tok.Type == lexing.EOF {
			break
		}

		if p.matchKeyword("elif") {
			block := p.parseIfConditionBlock()

			if !p.match(lexing.RIGHT_BRACE) {
				p.error("Expected '}' after elif block")
				return nil
			}

			elifBlocks = append(elifBlocks, block)
		} else {
			break
		}
	}

	// parse optional 'else' block
	elseBlock := make([]ast.Statement, 0)

	if p.matchKeyword("else") {
		if !p.match(lexing.LEFT_BRACE) {
			p.error("Expected '{' after else")
			return nil
		}
		elseBlock = p.parseBlock()

		if !p.match(lexing.RIGHT_BRACE) {
			p.error("Expected '}' after else block")
			return nil
		}
	}

	return &statement.IfStatement{
		Location:         startToken.Location,
		PrimaryCondition: primaryCondition,
		PrimaryBlock:     primaryBlock,
		ElseIfBlocks:     elifBlocks,
		ElseBlock:        elseBlock,
	}
}

// parseIfConditionBlock  parses a condition followed by a block
// I.E:
// 1+1 == 2 {
// }
func (p *Parser) parseIfConditionBlock() *statement.IfConditionBlock {
	condition := p.parseExpression()
	if condition == nil {
		return nil
	}

	// Parse body
	if !p.match(lexing.LEFT_BRACE) {
		p.error("Expected '{' after if condition")
		return nil
	}

	body := p.parseBlock()

	// Consume the right brace
	if !p.match(lexing.RIGHT_BRACE) {
		p.error("Expected '}' after if condition")
		return nil
	}

	return &statement.IfConditionBlock{
		Location:  condition.GetLocation(),
		Condition: condition,
		Body:      body,
	}
}

func (p *Parser) parseReturnStatement() ast.Statement {
	startToken := p.previous() // The 'return' token

	// after 'return' could be a closing } or an expression
	var exp ast.Expression = nil

	if !p.check(lexing.RIGHT_BRACE) {
		exp = p.parseExpression()
	}

	return &statement.ReturnStatmenet{
		Location:   startToken.Location,
		Expression: exp,
	}
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

// parseFuncDeclaration: Parses a function declaration
func (p *Parser) parseFuncDeclaration() ast.Statement {
	startToken := p.previous() // Save the 'func' token for error reporting

	// Parse function name
	name := p.consume(lexing.IDENTIFIER, "Expected function name")
	if len(p.errors) > 0 {
		return nil
	}

	// Parse function parameters
	if !p.match(lexing.LEFT_PAREN) {
		p.error("Expected '(' after function name")
		return nil
	}

	// Parse function parameters with comma between each parameter
	parameters := make([]expression.FuncParameterExpression, 0)

	// Handle parameters
	if !p.check(lexing.RIGHT_PAREN) {
		for {
			// Parse parameter
			param := p.parseFuncParameter()
			if param == nil {
				// Error already reported by parseFuncParameter
				return nil
			}
			parameters = append(parameters, *param)

			// Check for comma or end of parameters
			if p.check(lexing.RIGHT_PAREN) {
				break
			}

			if !p.match(lexing.COMMA) {
				p.error("Expected ',' or ')' after function parameter")
				return nil
			}
		}
	}

	// Consume the closing parenthesis
	if !p.match(lexing.RIGHT_PAREN) {
		p.error("Expected ')' after function parameters")
		return nil
	}

	// Parse optional return type (defaults to "void" if not specified)
	returnType := "void"
	if p.match(lexing.COLON) {
		tok := p.peek()

		if tok.Type == lexing.KEYWORD || tok.Type == lexing.IDENTIFIER {
			returnType = tok.Literal
		} else {
			p.error("Invalid function return type " + tok.Literal)
		}

		p.advance()
	}

	// Parse function body
	if !p.match(lexing.LEFT_BRACE) {
		p.error("Expected '{' after function declaration")
		return nil
	}

	// Parse the function body
	body := p.parseBlock()
	if body == nil {
		// Error already reported by parseBlock
		return nil
	}

	// Consume the closing brace
	if !p.match(lexing.RIGHT_BRACE) {
		p.error("Expected '}' after function body")
		return nil
	}

	return statement.NewFuncDeclaration(
		name.Literal,
		parameters,
		returnType,
		body,
		startToken.Location,
	)
}

// parseFuncParameter: Parses a function parameter
// parseFuncParameter: Parses a function parameter
func (p *Parser) parseFuncParameter() *expression.FuncParameterExpression {
	name := p.consume(lexing.IDENTIFIER, "Expected parameter name")
	if len(p.errors) > 0 {
		return nil
	}

	// type
	if !p.match(lexing.COLON) {
		p.error("Expected ':' after parameter name")
		return nil
	}

	typeToken := p.consume(lexing.IDENTIFIER, "Expected parameter type")
	if len(p.errors) > 0 {
		return nil
	}

	// Check for nullable type marker
	isNullable := p.match(lexing.QMARK)

	// optional default value
	var defaultValue ast.Expression
	if p.match(lexing.ASSIGN) {
		defaultValue = p.parseExpression()
		if defaultValue == nil {
			p.error("Expected expression after '='")
			return nil
		}
	}

	return expression.NewFuncParameterExpression(
		name.Literal,
		typeToken.Literal,
		isNullable,
		name.Location,
		defaultValue,
	)
}

// parseExpression: Entry point for expression parsing
func (p *Parser) parseExpression() ast.Expression {
	return p.parseLogicalOr()
}

// parseLogicalOr parses logical OR expressions
func (p *Parser) parseLogicalOr() ast.Expression {
	expr := p.parseLogicalAnd()

	for p.matchKeyword("or") {
		operator := p.previous().Literal
		right := p.parseLogicalAnd()
		if right == nil {
			p.error("Expected expression after 'or'")
			return nil
		}
		expr = expression.NewBinaryExpression(expr, operator, right, p.previous().Location)
	}

	return expr
}

// parseLogicalAnd parses logical AND expressions
func (p *Parser) parseLogicalAnd() ast.Expression {
	expr := p.parseEquality()

	for p.matchKeyword("and") {
		operator := p.previous().Literal
		right := p.parseEquality()
		if right == nil {
			p.error("Expected expression after 'and'")
			return nil
		}
		expr = expression.NewBinaryExpression(expr, operator, right, p.previous().Location)
	}

	return expr
}

// parseEquality parses equality expressions
func (p *Parser) parseEquality() ast.Expression {
	expr := p.parseComparison()

	for p.match(lexing.EQUALS, lexing.NOT_EQUALS) {
		operator := p.previous().Literal
		right := p.parseComparison()
		if right == nil {
			p.error("Expected expression after comparison operator")
			return nil
		}
		expr = expression.NewBinaryExpression(expr, operator, right, p.previous().Location)
	}

	return expr
}

// parseComparison parses comparison expressions
func (p *Parser) parseComparison() ast.Expression {
	expr := p.parseAdditive()

	for p.match(lexing.LESS, lexing.LESS_EQUALS, lexing.GREATER, lexing.GREATER_EQUALS) {
		operator := p.previous().Literal
		right := p.parseAdditive()
		if right == nil {
			p.error("Expected expression after comparison operator")
			return nil
		}
		expr = expression.NewBinaryExpression(expr, operator, right, p.previous().Location)
	}

	return expr
}

// parseAdditive: Parses addition and subtraction
func (p *Parser) parseAdditive() ast.Expression {
	expr := p.parseMultiplicative()

	for p.match(lexing.PLUS, lexing.MINUS) {
		operator := p.previous().Literal
		right := p.parseMultiplicative()
		if right == nil {
			if p.isAtEnd() {
				p.errorAtToken(p.peek(), "Expected expression after operator")
			} else if p.peek().Type == lexing.PLUS {
				p.errorAtToken(p.peek(), "Expected expression between operators")
			} else {
				p.errorAtToken(p.peek(), "Expected expression after operator")
			}
			return nil
		}
		expr = expression.NewBinaryExpression(expr, operator, right, p.previous().Location)
	}

	return expr
}

// parseMultiplicative: Parses multiplication and division
func (p *Parser) parseMultiplicative() ast.Expression {
	if p.check(lexing.MULTIPLY) || p.check(lexing.DIVIDE) {
		p.errorAtToken(p.peek(), "Expected expression before operator")
		p.advance() // Skip the operator
		return nil
	}

	expr := p.parseUnary()

	for p.match(lexing.MULTIPLY, lexing.DIVIDE) {
		operator := p.previous().Literal
		right := p.parseUnary()
		if right == nil {
			p.errorAtToken(p.peek(), "Expected expression after operator")
			return nil
		}
		expr = expression.NewBinaryExpression(expr, operator, right, p.previous().Location)
	}

	return expr
}

// parseUnary: Parses unary operators
func (p *Parser) parseUnary() ast.Expression {
	if p.match(lexing.MINUS) || p.matchKeyword("not") {
		operator := p.previous().Literal
		expr := p.parseUnary()
		if expr == nil {
			p.errorAtToken(p.peek(), "Expected expression after unary operator")
			return nil
		}
		return expression.NewUnaryExpression(operator, expr, p.previous().Location)
	}

	return p.parseCall()
}

// parseCall: Parses function calls
// parseCall parses function calls and primary expressions
func (p *Parser) parseCall() ast.Expression {
	expr := p.parsePrimary()
	if expr == nil {
		return nil
	}

	for {
		if p.match(lexing.LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else {
			break
		}
	}

	return expr
}

// finishCall handles the parsing of function call arguments after '(' has been matched
func (p *Parser) finishCall(callee ast.Expression) ast.Expression {
	args := make([]ast.Expression, 0)
	if !p.check(lexing.RIGHT_PAREN) {
		for {
			arg := p.parseExpression()
			if arg == nil {
				p.error("Expected argument in function call")
				return nil
			}
			args = append(args, arg)

			if !p.match(lexing.COMMA) {
				break
			}
		}
	}

	if !p.match(lexing.RIGHT_PAREN) {
		p.error("Expected ')' after function arguments")
		return nil
	}

	return &expression.CallExpression{
		Callee:    callee,
		Arguments: args,
		Location:  p.previous().Location,
	}
}

// parsePrimary: Parses primary expressions (literals, identifiers, and parentheses)
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

	case lexing.FLOAT:
		p.advance()
		// Convert string to float
		value, err := strconv.ParseFloat(token.Literal, 64)
		if err != nil {
			p.errorAtToken(token, "Invalid float literal")
			return nil
		}
		return expression.NewLiteralExpression(value, token.Location)

	case lexing.IDENTIFIER:
		p.advance()
		return expression.NewIdentifierExpression(token.Literal, token.Location)

	case lexing.KEYWORD:
		if token.Literal == "true" || token.Literal == "false" {
			p.advance()
			return expression.NewLiteralExpression(token.Literal == "true", token.Location)
		} else if token.Literal == "null" {
			p.advance()
			return expression.NewLiteralExpression(nil, token.Location)
		}

	case lexing.LEFT_PAREN:
		p.advance()
		expr := p.parseExpression()
		if expr == nil {
			return nil
		}
		if !p.match(lexing.RIGHT_PAREN) {
			p.errorAtToken(p.peek(), "Expected closing parenthesis")
			return nil
		}
		return expr
	}

	p.errorAtToken(token, "Expected expression")
	return nil
}
