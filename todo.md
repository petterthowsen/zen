# Parser Implementation TODO

## Statement Types
- [x] Variable declarations
  - [x] Basic declarations (var name = value)
  - [x] Type annotations (var name:type = value)
  - [x] Nullable types (var name:type? = value)
  - [x] Support for string, int, boolean, and null literals
- [x] Function declarations
  - [x] Parameters
  - [x] Return type
  - [x] Function body
- [ ] Class declarations
  - [ ] Properties
  - [ ] Methods
  - [ ] Inheritance
- [ ] Interface declarations
  - [ ] Method signatures
  - [ ] Property signatures
- [ ] Control flow
  - [x] If statements with complex conditions
  - [ ] For loops
  - [ ] While loops
  - [ ] When statements
  - [x] Return statements

## Expression Types
- [x] Literals
  - [x] String literals
  - [x] Integer literals
  - [x] Boolean literals
  - [x] Null literal
  - [x] Float literals
- [x] Binary expressions
  - [x] Arithmetic operators (+, -, *, /, %)
  - [x] Comparison operators (==, !=, <, >, <=, >=)
  - [x] Logical operators (and, or)
- [x] Unary expressions
  - [x] Negation (-)
  - [x] Not (not)
- [ ] Member access (obj.prop)
- [x] Function calls
- [ ] Array literals and access
- [ ] Map literals and access
- [ ] Lambda expressions

## Error Handling
- [x] Basic error recovery
- [x] Source location tracking
- [ ] Error synchronization
- [x] Multiple error reporting (but synchronization isn't working)
- [ ] Better error messages
  - [x] Expected vs actual token
  - [ ] Suggestions for common mistakes
  - [ ] Context-aware hints

## Type System
- [ ] Basic type checking
- [ ] Type inference
- [ ] Generic types
- [ ] Union types
- [ ] Nullable type handling
- [ ] Type compatibility rules

## Testing
- [x] Variable declaration tests
- [x] Function declaration tests
- [ ] Class declaration tests
- [ ] Interface declaration tests
- [x] Expression tests
- [ ] Error recovery tests
- [ ] Edge case tests
