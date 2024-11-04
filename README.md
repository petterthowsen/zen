# Zen Programming Language

Zen is a high-level interpreted programming language implemented in Go. It strong typing with type inference where possible and OOP features.

See spec.zen for detailed examples of zen code.

## Project Structure

```
zenlang/
├── lang/                   # Core language implementation
│   ├── common/            # Common utilities and types
│   │   ├── SourceCode.go      # Source code handling
│   │   ├── SourceLocation.go  # Source code Location tracking (line, column and filename)
│   │   └── SyntaxError.go     # Error handling
│   ├── lexing/            # Lexical analysis
│   │   ├── Lexer.go          # Token generation
│   │   └── Token.go          # Token definitions
│   └── parsing/           # Syntax analysis
│       ├── ast/              # Abstract Syntax Tree definitions
│       ├── expression/       # Expression nodes
│       │   ├── BinaryExpression.go     # Binary operations
│       │   ├── CallExpression.go       # Function calls
│       │   ├── LiteralExpression.go    # Literals
│       │   ├── MemberAccessExpression.go # Object property access
│       │   └── UnaryExpression.go      # Unary operations
│       ├── statement/        # Statement nodes
│       │   ├── BreakStatement.go       # Break from loops
│       │   ├── ContinueStatement.go    # Continue in loops
│       │   ├── ForStatement.go         # For loops
│       │   ├── ForInStatement.go       # For-in loops
│       │   ├── FuncDeclaration.go      # Function declarations
│       │   ├── IfStatement.go          # If conditions
│       │   ├── ReturnStatement.go      # Return statements
│       │   ├── VarDeclaration.go       # Variable declarations
│       │   └── WhileStatement.go       # While loops
│       └── Parser.go         # Main Parser implementation
├── tests/                 # Test suite
│   ├── lexing/           # Lexer tests
│   └── parsing/          # Parser tests
├── spec.zen              # Language specification
└── grammar.ebnf          # Formal grammar definition (Out of date!)
```

## Features

### Control Flow
- If statements with complex conditions
- For loops with initialization, condition, and increment
- For-in loops for iteration
- While loops
- Break and continue statements
- Return statements

### Expressions
- Literals (string, integer, float, boolean, null)
- Binary operations (+, -, *, /, %, ==, !=, <, >, <=, >=, and, or)
- Unary operations (-, not)
- Member access (obj.prop, obj.nested.prop)
- Bracket access for arrays (myArray[5])
- Curly access for maps (myMap{"name"})
- Function calls

### Declarations
- Variable declarations with type annotations and nullability with question mark
- Function declarations with parameters and return types

## Running Tests

Tests are organized by component. Most test files have a corresponding `.zen` file containing the test cases.

To run all tests:
```bash
go test ./...
```

To run tests for a specific component:
```bash
go test ./tests/lexing    # Run lexer tests
go test ./tests/parsing   # Run parser tests
```

To run a specific test with verbose output:
```bash
go test ./tests/parsing -v -run TestVarDeclaration
```

### Test Output Files

The test suite generates output files to help with debugging and verification:
- `.tokens` files: Show the lexical analysis results
- `.ast` files: Show the parsed Abstract Syntax Tree

### Test Utilities

The project provides test utilities to simplify writing and maintaining tests:

#### Parser Test Utilities (`tests/parsing/parser_test_utils.go`)

1. File Handling:
```go
// Parse a test file and save its AST
program := ParseTestFile(t, "your_test.zen")

// Parse a string directly
program, errors := ParseString("var x = 42")
```

2. AST Node Assertions:
```go
// Assert variable declarations
varDecl := AssertVarDeclaration(t, stmt, "name", "type", isConst, isNullable)

// Assert expressions
literal := AssertLiteralExpression(t, expr, 42)
binary := AssertBinaryExpression(t, expr, "+")
unary := AssertUnaryExpression(t, expr, "-")
ident := AssertIdentifierExpression(t, expr, "varName")
call := AssertCallExpression(t, expr, expectedArgCount)
member := AssertMemberAccess(t, expr, "object", "property")

// Assert statements
ifStmt := AssertIfStatement(t, stmt)
funcDecl := AssertFuncDeclaration(t, stmt)

// Assert parsing errors
AssertParseError(t, "invalid { syntax")
```

## Implementing New Features

### 1. Update the Grammar
First, add your new feature to `grammar.ebnf`. This ensures the syntax is formally defined.

### 2. Add Lexer Support
If your feature requires new tokens:
1. Add token types in `lang/lexing/Token.go`
2. Update the lexer in `lang/lexing/Lexer.go`
3. Add test cases in `tests/lexing/`

Example token addition:
```go
// In Token.go
const (
    // ... existing tokens ...
    NEW_TOKEN
)

var tokenTypeNames = map[TokenType]string{
    // ... existing mappings ...
    NEW_TOKEN: "NewToken",
}
```

### 3. Add Parser Support

1. Create AST Nodes:
   - For expressions: Add to `lang/parsing/expression/`
   - For statements: Add to `lang/parsing/statement/`

Example AST node:
```go
// In statement/NewFeature.go
type NewFeatureNode struct {
    // Node properties
    Location *common.SourceLocation
}

func (n *NewFeatureNode) Accept(visitor ast.Visitor) interface{} {
    return visitor.VisitNewFeature(n)
}

func (n *NewFeatureNode) GetLocation() *common.SourceLocation {
    return n.Location
}

func (n *NewFeatureNode) IsStatement() {}

func (n *NewFeatureNode) String(indent int) string {
    return strings.Repeat("  ", indent) + "NewFeature\n"
}
```

2. Update the parser in `lang/parsing/Parser.go`:
```go
func (p *Parser) parseStatement() ast.Statement {
    if p.matchKeyword("newfeature") {
        return p.parseNewFeature()
    }
    // ... existing statement parsing ...
}

func (p *Parser) parseNewFeature() ast.Statement {
    // Implementation
}
```

### 4. Testing Strategy

1. Create test files:
   - `tests/parsing/new_feature.zen`: Example code using the new feature
   - `tests/parsing/new_feature_test.go`: Test cases using test utilities

2. Write comprehensive tests:
```go
func TestNewFeature(t *testing.T) {
    // Test successful parsing using utilities
    program := ParseTestFile(t, "new_feature.zen")
    if program == nil {
        return
    }

    // Use assertion helpers for validation
    stmt := program.Statements[0]
    feature := AssertNewFeature(t, stmt, expectedName)
    
    // Validate complex structures
    if feature.SomeExpression != nil {
        AssertBinaryExpression(t, feature.SomeExpression, "+")
    }

    // Test error cases
    AssertParseError(t, "incomplete newfeature")
    AssertParseError(t, "newfeature with { invalid syntax")
}
```

3. Verify AST output:
   - Check the generated `.ast` file for correct structure
   - Ensure proper indentation and node relationships
   - Validate error messages and locations

## Error Handling

The project uses `SyntaxError` for error reporting, which includes:
- Error message
- Source location
- Pretty printing with source context

When implementing new features, use the error system:
```go
p.error("Expected { after new feature declaration")
```

## Parser Implementation

The Parser in Zen follows a recursive descent parsing strategy, transforming a sequence of tokens into an Abstract Syntax Tree (AST). Here's how it works:

### Core Components

1. **Parser State Management**
   - Maintains current position in token stream
   - Tracks syntax errors
   - Provides utilities for token consumption and lookahead

2. **Expression Parsing**
   - Implements operator precedence through recursive descent
   - Handles binary operations (+, -, *, /, etc.)
   - Supports unary operations (-, not)
   - Processes function calls and member access
   - Parses primary expressions (literals, identifiers, parenthesized expressions)

3. **Statement Parsing**
   - Recognizes and delegates to specific parsers for:
     - Variable declarations (var/const)
     - Function declarations
     - Control flow (if, for, while)
     - Break/Continue statements
     - Return statements
     - Expression statements

### Parsing Process

1. **Token Stream Processing**
   ```go
   func (p *Parser) Parse() (*ast.ProgramNode, []*common.SyntaxError) {
       statements := make([]ast.Statement, 0)
       for !p.isAtEnd() {
           stmt := p.parseStatement()
           if stmt != nil {
               statements = append(statements, stmt)
           }
       }
       return ast.NewProgramNode(statements), p.errors
   }
   ```

2. **Expression Precedence**
   The parser handles operator precedence through a series of recursive functions:
   ```
   parseExpression
   └── parseAssignment
       └── parseLogicalOr
           └── parseLogicalAnd
               └── parseEquality
                   └── parseComparison
                       └── parseAdditive
                           └── parseMultiplicative
                               └── parseUnary
                                   └── parsePostfix
                                       └── parseCall
                                           └── parsePrimary
   ```

3. **Error Recovery**
   - Implements synchronization points for error recovery
   - Continues parsing after encountering errors
   - Maintains error context for debugging

### Key Features

1. **Flexible Statement Parsing**
   ```go
   func (p *Parser) parseStatement() ast.Statement {
       if p.matchKeyword("var", "const") {
           return p.parseVarDeclaration()
       }
       if p.matchKeyword("func") {
           return p.parseFuncDeclaration()
       }
       // ... other statement types
   }
   ```

2. **Expression Handling**
   ```go
   func (p *Parser) parseAdditive() ast.Expression {
       expr := p.parseMultiplicative()
       for p.match(lexing.PLUS, lexing.MINUS) {
           operator := p.previous().Literal
           right := p.parseMultiplicative()
           expr = expression.NewBinaryExpression(expr, operator, right, p.previous().Location)
       }
       return expr
   }
   ```

3. **Error Handling**
   ```go
   func (p *Parser) synchronize() bool {
       for !p.isAtEnd() {
           if p.checkKeyword("var") || p.checkKeyword("func") {
               return true // Found a synchronization point
           }
           p.advance()
       }
       return false
   }
   ```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for your changes
4. Ensure all tests pass
5. Submit a pull request

## License

MIT
