# Zen Language Parser Proposal

## Overview

This document outlines a proposal for implementing a recursive descent parser for the Zen programming language. The parser will transform a sequence of tokens from the lexer into an Abstract Syntax Tree (AST) that represents the program structure.

## Grammar

See grammar.txt

## Parser Design

### Core Components

1. **Parser Structure**
```go
// Parser maintains the state during parsing
type Parser struct {
    tokens []lexing.Token    // Input token stream
    current int              // Current token index
    errors []common.SyntaxError  // Accumulated errors
    panicMode bool           // Error recovery state
}

// NewParser creates a new parser instance
func NewParser(tokens []lexing.Token) *Parser {
    return &Parser{
        tokens: tokens,
        current: 0,
        errors: make([]common.SyntaxError, 0),
        panicMode: false,
    }
}
```

2. **Helper Methods**
```go
// Core parsing utilities
func (p *Parser) peek() lexing.Token {
    // Returns current token without consuming
    return p.tokens[p.current]
}

func (p *Parser) previous() lexing.Token {
    // Returns last consumed token
    return p.tokens[p.current-1]
}

func (p *Parser) advance() lexing.Token {
    // Moves to next token and returns current
    if !p.isAtEnd() {
        p.current++
    }
    return p.previous()
}

func (p *Parser) match(types ...lexing.TokenType) bool {
    // Checks if current token matches any of the given types
    for _, t := range types {
        if p.check(t) {
            p.advance()
            return true
        }
    }
    return false
}

func (p *Parser) consume(typ lexing.TokenType, message string) lexing.Token {
    // Consumes token of expected type or reports error
    if p.check(typ) {
        return p.advance()
    }
    p.error(p.peek(), message)
    return lexing.Token{}
}

func (p *Parser) synchronize() {
    // Recovers from error by finding next safe parsing point
    p.panicMode = false
    
    for !p.isAtEnd() {
        if p.previous().Type == lexing.SEMICOLON {
            return
        }

        switch p.peek().Type {
        case lexing.CLASS, lexing.FUNC, lexing.VAR, lexing.FOR, 
             lexing.IF, lexing.WHILE, lexing.RETURN:
            return
        }

        p.advance()
    }
}
```

### AST Node Types

1. **Base Node Interface**
```go
// Node represents any node in the AST
type Node interface {
    // String returns a debug representation
    String() string
    // Location returns source position
    Location() *common.SourceLocation
}
```

2. **Expressions**
```go
// Expression represents any value-producing construct
type Expression interface {
    Node
    expressionNode()
}

// Literal expressions
type IntegerLiteral struct {
    Token lexing.Token
    Value int64
}

type FloatLiteral struct {
    Token lexing.Token
    Value float64
}

type StringLiteral struct {
    Token lexing.Token
    Value string
}

type BooleanLiteral struct {
    Token lexing.Token
    Value bool
}

// Binary expressions (a + b, a * b, etc.)
type BinaryExpression struct {
    Left     Expression
    Operator lexing.Token
    Right    Expression
}

// Unary expressions (-a, !b, etc.)
type UnaryExpression struct {
    Operator lexing.Token
    Right    Expression
}

// Variable reference
type Identifier struct {
    Token lexing.Token
    Name  string
}

// Array expressions
type ArrayLiteral struct {
    Elements []Expression
    Size     Expression // Optional for dynamic arrays
    Type     Expression
}

// Map expressions
type MapLiteral struct {
    Pairs map[Expression]Expression
    Type  Expression
}

// Function call
type CallExpression struct {
    Callee    Expression
    Arguments []Expression
    Location  *common.SourceLocation
}

// Property access (obj.prop)
type PropertyExpression struct {
    Object   Expression
    Property *Identifier
}
```

3. **Statements**
```go
// Statement represents any executable construct
type Statement interface {
    Node
    statementNode()
}

// Variable declaration
type VarDeclaration struct {
    Name        string
    Type        Expression // Optional
    Initializer Expression
    IsConst     bool
    Location    *common.SourceLocation
}

// Function declaration
type FunctionDeclaration struct {
    Name       string
    Parameters []Parameter
    ReturnType Expression
    Body       []Statement
    Location   *common.SourceLocation
}

// Parameter represents a function parameter
type Parameter struct {
    Name     string
    Type     Expression
    Optional bool
}

// Class declaration
type ClassDeclaration struct {
    Name       string
    SuperClass string
    Interfaces []string
    Methods    []FunctionDeclaration
    Properties []VarDeclaration
    Location   *common.SourceLocation
}

// Control flow statements
type IfStatement struct {
    Condition   Expression
    ThenBranch  []Statement
    ElseBranch  []Statement
    Location    *common.SourceLocation
}

type ForStatement struct {
    // Traditional for loop
    Initializer Statement
    Condition   Expression
    Increment   Statement
    // Range-based for loop
    LoopVar     string
    Collection  Expression
    // Common
    Body        []Statement
    Location    *common.SourceLocation
}

type WhileStatement struct {
    Condition Expression
    Body      []Statement
    Location  *common.SourceLocation
}

type MatchStatement struct {
    Value    Expression
    Cases    []MatchCase
    ElseCase []Statement
    Location *common.SourceLocation
}

type MatchCase struct {
    Pattern Expression
    Body    []Statement
}

// Return statement
type ReturnStatement struct {
    Value    Expression
    Location *common.SourceLocation
}
```

## Parsing Strategy

### 1. Expression Parsing

The parser uses precedence climbing for expressions, which provides a clean and efficient way to handle operator precedence:

```go
// Operator precedence levels
var precedence = map[lexing.TokenType]int{
    lexing.OR:             1,  // or
    lexing.AND:            2,  // and
    lexing.EQUALS:         3,  // ==
    lexing.NOT_EQUALS:     3,  // !=
    lexing.LESS:           4,  // <
    lexing.LESS_EQUALS:    4,  // <=
    lexing.GREATER:        4,  // >
    lexing.GREATER_EQUALS: 4,  // >=
    lexing.PLUS:           5,  // +
    lexing.MINUS:          5,  // -
    lexing.MULTIPLY:       6,  // *
    lexing.DIVIDE:         6,  // /
    lexing.PERCENT:        6,  // %
    lexing.DOT:           7,  // .
    lexing.LEFT_PAREN:    8,  // (
    lexing.LEFT_BRACKET:  8,  // [
}
```

Expression parsing methods with examples:

```go
func (p *Parser) expression() Expression {
    return p.assignment()
}

// Handles variable assignment
// Example: x = 5
func (p *Parser) assignment() Expression {
    expr := p.logicalOr()
    
    if p.match(lexing.ASSIGN) {
        equals := p.previous()
        value := p.assignment()
        
        if id, ok := expr.(*Identifier); ok {
            return &AssignmentExpression{
                Name: id.Name,
                Value: value,
            }
        }
        p.error(equals, "Invalid assignment target")
    }
    
    return expr
}

// Handles binary operations based on precedence
func (p *Parser) binary(precedenceLevel int) Expression {
    left := p.unary()
    
    for p.peek().Type.Precedence() > precedenceLevel {
        operator := p.advance()
        right := p.binary(operator.Type.Precedence())
        left = &BinaryExpression{
            Left: left,
            Operator: operator,
            Right: right,
        }
    }
    
    return left
}

// Example AST for: 1 + 2 * 3
//        +
//       / \
//      1   *
//         / \
//        2   3
```

### 2. Statement Parsing

Each statement type has its own parsing method. Here's how they work:

```go
// Parses any statement
func (p *Parser) statement() Statement {
    switch {
    case p.match(lexing.VAR):
        return p.varDeclaration()
    case p.match(lexing.FUNC):
        return p.functionDeclaration()
    case p.match(lexing.CLASS):
        return p.classDeclaration()
    case p.match(lexing.IF):
        return p.ifStatement()
    case p.match(lexing.FOR):
        return p.forStatement()
    case p.match(lexing.WHILE):
        return p.whileStatement()
    case p.match(lexing.MATCH):
        return p.matchStatement()
    case p.match(lexing.RETURN):
        return p.returnStatement()
    default:
        return p.expressionStatement()
    }
}

// Example: Variable Declaration
// var x:int = 5
func (p *Parser) varDeclaration() Statement {
    name := p.consume(lexing.IDENTIFIER, "Expected variable name")
    
    var typeExpr Expression
    if p.match(lexing.COLON) {
        typeExpr = p.typeExpression()
    }
    
    var initializer Expression
    if p.match(lexing.ASSIGN) {
        initializer = p.expression()
    }
    
    p.consume(lexing.SEMICOLON, "Expected ';' after variable declaration")
    
    return &VarDeclaration{
        Name: name.Literal,
        Type: typeExpr,
        Initializer: initializer,
        Location: name.Location,
    }
}

// Example: If Statement
// if x > 0 { print(x) } else { print("negative") }
func (p *Parser) ifStatement() Statement {
    condition := p.expression()
    thenBranch := p.block()
    
    var elseBranch []Statement
    if p.match(lexing.ELSE) {
        if p.match(lexing.IF) {
            elseBranch = []Statement{p.ifStatement()}
        } else {
            elseBranch = p.block()
        }
    }
    
    return &IfStatement{
        Condition: condition,
        ThenBranch: thenBranch,
        ElseBranch: elseBranch,
        Location: condition.Location(),
    }
}
```

### 3. Type Parsing

Zen supports generic types and nullable types. Here's how they're parsed:

```go
func (p *Parser) typeExpression() Expression {
    base := p.identifier()
    
    // Handle generics: Array<int, 3>
    if p.match(lexing.LESS) {
        var params []Expression
        
        if !p.check(lexing.GREATER) {
            do {
                params = append(params, p.typeExpression())
            } while p.match(lexing.COMMA)
        }
        
        p.consume(lexing.GREATER, "Expected '>' after type parameters")
        
        base = &GenericType{
            Base: base,
            Parameters: params,
        }
    }
    
    // Handle nullable types: string?
    if p.match(lexing.QMARK) {
        base = &NullableType{
            BaseType: base,
        }
    }
    
    return base
}
```

### 4. Error Recovery

The parser implements panic mode error recovery:

```go
func (p *Parser) synchronize() {
    p.panicMode = false
    
    for !p.isAtEnd() {
        if p.previous().Type == lexing.SEMICOLON {
            return
        }

        switch p.peek().Type {
        case lexing.CLASS, lexing.FUNC, lexing.VAR, lexing.FOR, 
             lexing.IF, lexing.WHILE, lexing.RETURN:
            return
        }

        p.advance()
    }
}

func (p *Parser) error(token lexing.Token, message string) {
    if p.panicMode {
        return // Avoid cascading errors
    }
    
    p.errors = append(p.errors, common.SyntaxError{
        Message: message,
        Location: token.Location,
    })
    
    p.panicMode = true
}
```

Example error recovery:
```go
// Input: var x = 5 + * 3;
//                    ^ Error here
// Parser will:
// 1. Report "Unexpected '*' after '+'"
// 2. Enter panic mode
// 3. Skip tokens until semicolon
// 4. Continue parsing next statement
```

## Implementation Plan

### Phase 1: Basic Expression Parsing (Week 1)
- [x] Setup parser structure
- [ ] Implement literal parsing
- [ ] Add binary operations
- [ ] Handle unary operations
- [ ] Support grouping with parentheses

### Phase 2: Statement Parsing (Week 2)
- [ ] Variable declarations
- [ ] Basic control flow (if, while)
- [ ] Function declarations
- [ ] Return statements

### Phase 3: Advanced Features (Week 3)
- [ ] Class declarations
- [ ] Interface implementations
- [ ] Generic type parsing
- [ ] Match statements
- [ ] For loops with range support

### Phase 4: Error Handling & Recovery (Week 4)
- [ ] Implement error reporting
- [ ] Add synchronization points
- [ ] Enhance error messages
- [ ] Add recovery strategies

### Phase 5: Optimization & Refinement (Week 5)
- [ ] Optimize parser performance
- [ ] Enhance AST for better semantic analysis
- [ ] Add source location tracking
- [ ] Improve error messages

## Testing Strategy

### 1. Unit Tests

Test individual parser components:
```go
func TestIntegerLiteral(t *testing.T) {
    input := "42"
    parser := NewParser(lexer.Tokenize(input))
    
    expr := parser.expression()
    literal, ok := expr.(*IntegerLiteral)
    
    assert.True(t, ok)
    assert.Equal(t, int64(42), literal.Value)
}

func TestBinaryExpression(t *testing.T) {
    input := "1 + 2 * 3"
    parser := NewParser(lexer.Tokenize(input))
    
    expr := parser.expression()
    binary, ok := expr.(*BinaryExpression)
    
    assert.True(t, ok)
    assert.Equal(t, lexing.PLUS, binary.Operator.Type)
    
    right, ok := binary.Right.(*BinaryExpression)
    assert.True(t, ok)
    assert.Equal(t, lexing.MULTIPLY, right.Operator.Type)
}
```

### 2. Integration Tests

Test complete program parsing:
```go
func TestCompleteProgram(t *testing.T) {
    input := `
        func fibonacci(n:int):int {
            if n <= 1 {
                return n
            }
            return fibonacci(n-1) + fibonacci(n-2)
        }
    `
    parser := NewParser(lexer.Tokenize(input))
    program := parser.parse()
    
    assert.Nil(t, parser.errors)
    assert.Equal(t, 1, len(program.Declarations))
    
    funcDecl, ok := program.Declarations[0].(*FunctionDeclaration)
    assert.True(t, ok)
    assert.Equal(t, "fibonacci", funcDecl.Name)
}
```

### 3. Edge Cases

Test complex scenarios:
```go
func TestComplexGenericType(t *testing.T) {
    input := "var matrix:Array<Array<int, 3>, 2>"
    parser := NewParser(lexer.Tokenize(input))
    
    decl := parser.varDeclaration()
    assert.Nil(t, parser.errors)
    
    // Verify nested generic type structure
}

func TestErrorRecovery(t *testing.T) {
    input := "var x = 1 + * 2; var y = 3;"
    parser := NewParser(lexer.Tokenize(input))
    
    program := parser.parse()
    assert.Equal(t, 1, len(parser.errors))
    assert.Equal(t, 2, len(program.Declarations))
}
```

## Conclusion

This parser design provides a robust foundation for the Zen language, supporting all required language features while maintaining good error handling and recovery capabilities. The recursive descent approach allows for clear implementation of the grammar rules while the AST structure enables proper semantic analysis in later compilation phases.

Key features:
1. Complete grammar specification in EBNF
2. Comprehensive AST node hierarchy
3. Efficient expression parsing with precedence climbing
4. Robust error recovery mechanism
5. Clear testing strategy with examples
6. Phased implementation plan

