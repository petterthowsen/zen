# Zen Programming Language

Zen is a high-level interpreted programming language implemented in Go. It features a clean syntax, strong typing, and modern programming concepts.

See spec.zen for example of zen code.

## Project Structure

```
zenlang/
├── lang/                   # Core language implementation
│   ├── common/            # Common utilities and types
│   │   ├── SourceCode.go      # Source code handling
│   │   ├── SourceLocation.go  # Location tracking
│   │   └── SyntaxError.go     # Error handling
│   ├── lexing/            # Lexical analysis
│   │   ├── Lexer.go          # Token generation
│   │   └── Token.go          # Token definitions
│   └── parsing/           # Syntax analysis
│       ├── ast/              # Abstract Syntax Tree definitions
│       ├── expression/       # Expression nodes
│       ├── statement/        # Statement nodes
│       └── Parser.go         # Parser implementation
├── tests/                 # Test suite
│   ├── lexing/           # Lexer tests
│   └── parsing/          # Parser tests
├── spec.zen              # Language specification
└── grammar.ebnf          # Formal grammar definition
```

## Running Tests

Tests are organized by component. Most test file have a corresponding `.zen` file containing the test cases.

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

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for your changes
4. Ensure all tests pass
5. Submit a pull request

## License

[License information here]
