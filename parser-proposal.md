# Zen Parser Implementation Proposal

## Overview
This document outlines the implementation strategy for the Zen language parser using the recursive descent algorithm. The parser will take an array of tokens as input and produce an Abstract Syntax Tree (AST) with ProgramNode as the root.

## Parser Architecture

### Parser Class Structure
```go
type Parser struct {
    tokens []Token
    current int
    stopAtFirstError bool
    errors []SyntaxError
}

func NewParser(tokens []Token, stopAtFirstError bool) *Parser {
    return &Parser{
        tokens: tokens,
        current: 0,
        stopAtFirstError: stopAtFirstError,
        errors: make([]SyntaxError, 0),
    }
}
```

### Main Parse Method
```go
func (p *Parser) Parse() (*ProgramNode, []SyntaxError) {
    statements := make([]Node, 0)
    
    for !p.isAtEnd() {
        stmt := p.parseStatement()
        if stmt != nil {
            statements = append(statements, stmt)
        }
    }
    
    return NewProgramNode(statements), p.errors
}
```

## AST Node Hierarchy

### Base Node Interface (lang/parsing/Node.go)
```go
type Node interface {
    Accept(visitor Visitor) interface{}
    GetLocation() SourceLocation
}
```

### Expression Nodes (lang/parsing/expression/)
```go
// Expression.go - Base interface
type Expression interface {
    Node
    isExpression() // Marker method
}

// Specific expression nodes:
- BinaryExpression (left, operator, right)
- UnaryExpression (operator, expression)
- LiteralExpression (value)
- IdentifierExpression (name)
- CallExpression (callee, arguments)
- ArrayExpression (elements)
- MapExpression (entries)
- TupleExpression (elements)
- MemberExpression (object, property)
- LambdaExpression (parameters, body)
```

### Statement Nodes (lang/parsing/statement/)
```go
// Statement.go - Base interface
type Statement interface {
    Node
    isStatement() // Marker method
}

// Specific statement nodes:
- VarDeclarationStatement (name, type, initializer)
- AssignmentStatement (target, operator, value)
- IfStatement (condition, thenBranch, elseBranch)
- ForStatement (init, condition, increment, body)
- WhileStatement (condition, body)
- ReturnStatement (value)
- MatchStatement (value, cases, elseBranch)
- FunctionDeclaration (name, parameters, returnType, body)
- ClassDeclaration (name, implements, members)
- InterfaceDeclaration (name, members)
- ImportStatement (path, alias)
```

## Operator Precedence
Following standard precedence rules:

1. Primary (literals, identifiers, grouping)
2. Postfix (member access, array access, function calls)
3. Unary (!, -, not)
4. Multiplicative (*, /)
5. Additive (+, -)
6. Relational (<, >, <=, >=)
7. Equality (==, !=)
8. Logical AND (and)
9. Logical OR (or)
10. Assignment (=, +=, -=, *=, /=)

Implementation:
```go
func (p *Parser) parseExpression() Expression {
    return p.parseAssignment()
}

func (p *Parser) parseAssignment() Expression {
    expr := p.parseOr()
    
    if p.match(EQUAL, PLUS_EQUAL, MINUS_EQUAL, STAR_EQUAL, SLASH_EQUAL) {
        operator := p.previous()
        value := p.parseAssignment()
        // Handle assignment
        return NewAssignmentExpression(expr, operator, value)
    }
    
    return expr
}

// Similar methods for each precedence level
func (p *Parser) parseOr() Expression
func (p *Parser) parseAnd() Expression
func (p *Parser) parseEquality() Expression
// etc...
```

## Error Handling and Recovery

### Error Types
```go
type SyntaxError struct {
    Message string
    Location SourceLocation
    ExpectedToken TokenType
    ActualToken Token
}
```

### Recovery Strategy
The parser uses structural elements (braces, keywords) rather than whitespace for synchronization points. Each token carries source location information for error reporting.

1. Synchronization points:
   - Opening braces of blocks
   - Statement-starting keywords (class, func, var, etc.)
   - Closing braces of blocks
   
2. Error Recovery Process:
   - On error, record the error with exact location from token
   - If stopAtFirstError is false, synchronize to next statement
   - Continue parsing from synchronized position

```go
func (p *Parser) synchronize() {
    p.advance()
    
    for !p.isAtEnd() {
        // Synchronize on block boundaries
        if p.previous().Type == RIGHT_BRACE {
            return
        }
        
        // Synchronize on statement-starting tokens
        switch p.peek().Type {
        case CLASS, FUNC, VAR, FOR, IF, WHILE, RETURN, LEFT_BRACE:
            return
        }
        
        p.advance()
    }
}

func (p *Parser) error(message string) {
    err := NewSyntaxError(message, p.peek().Location)
    p.errors = append(p.errors, err)
    
    if p.stopAtFirstError {
        panic(err) // Will be caught in Parse()
    }
}
```

## Implementation Strategy

### Helper Methods
```go
func (p *Parser) match(types ...TokenType) bool
func (p *Parser) advance() Token
func (p *Parser) peek() Token
func (p *Parser) previous() Token
func (p *Parser) check(type TokenType) bool
func (p *Parser) isAtEnd() bool
```

### Statement Parsing
Each statement type will have its own parsing method:
```go
func (p *Parser) parseStatement() Statement {
    switch {
    case p.match(VAR, CONST):
        return p.parseVarDeclaration()
    case p.match(IF):
        return p.parseIfStatement()
    case p.match(FOR):
        return p.parseForStatement()
    case p.match(WHILE):
        return p.parseWhileStatement()
    case p.match(FUNC):
        return p.parseFunctionDeclaration()
    case p.match(CLASS):
        return p.parseClassDeclaration()
    case p.match(RETURN):
        return p.parseReturnStatement()
    case p.match(MATCH):
        return p.parseMatchStatement()
    default:
        return p.parseExpressionStatement()
    }
}
```

### Type Parsing
Basic type checking during parsing:
```go
func (p *Parser) parseType() Type {
    if p.match(IDENTIFIER) {
        typeName := p.previous().Lexeme
        
        // Handle array types
        if p.match(LESS) {
            elementType := p.parseType()
            p.consume(GREATER, "Expect '>' after type parameters")
            return NewArrayType(elementType)
        }
        
        // Handle nullable types
        if p.match(QUESTION) {
            return NewNullableType(NewSimpleType(typeName))
        }
        
        return NewSimpleType(typeName)
    }
    
    p.error("Expect type name")
    return nil
}
```

## File Organization

```
lang/
├── parsing/
│   ├── Node.go             # Base node interface
│   ├── Parser.go           # Main parser implementation
│   ├── ProgramNode.go      # Root AST node
│   ├── Type.go            # Type system nodes
│   ├── expression/
│   │   ├── Expression.go   # Base expression interface
│   │   ├── Binary.go
│   │   ├── Unary.go
│   │   ├── Literal.go
│   │   ├── Identifier.go
│   │   └── ...
│   └── statement/
│       ├── Statement.go    # Base statement interface
│       ├── VarDecl.go
│       ├── FuncDecl.go
│       ├── ClassDecl.go
│       ├── If.go
│       └── ...
```

## Usage Example

```go
func ParseSource(source string) (*ProgramNode, []SyntaxError) {
    lexer := NewLexer(source)
    tokens := lexer.ScanTokens()
    
    parser := NewParser(tokens, false) // false = don't stop at first error
    return parser.Parse()
}
```

## Error Examples

```
Error: Unexpected token
Expected: '{'
Found: '('
At line 10, column 15

Error: Missing closing brace
Expected: '}'
Found: EOF
At line 20, column 1

Error: Invalid class member
Expected: function or variable declaration
Found: 'if'
At line 30, column 5
```

The error messages include exact source locations since each token maintains its position information, even though whitespace is not tokenized.
