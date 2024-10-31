# Zen Programming Language

Zen is a high-level interpreted programming language implemented in Go. It features a clean syntax, strong typing, and modern programming concepts.

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

Tests are organized by component and include both unit tests and integration tests. Each test file has a corresponding `.zen` file containing the test cases.

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
    NEW_TOKEN TokenType = iota
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

3. Add test cases:
   - Create `tests/parsing/new_feature.zen` with example code
   - Create `tests/parsing/new_feature_test.go` with test cases

### 4. Testing Strategy

1. Start with lexer tests to ensure proper tokenization
2. Add parser tests to verify AST construction
3. Use the `.tokens` and `.ast` files to verify output
4. Include error cases to test proper error handling

Example test structure:
```go
func TestNewFeature(t *testing.T) {
    // Test successful parsing
    program, errors := ParseString("newfeature { }")
    if len(errors) > 0 {
        t.Error("Expected no errors")
    }

    // Test error cases
    program, errors = ParseString("newfeature")
    if len(errors) == 0 {
        t.Error("Expected error for incomplete feature")
    }
}
```

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
