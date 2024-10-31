# Parser Implementation TODO

## Statement Types
- [x] Variable declarations
  - [x] Basic declarations (var name = value)
  - [x] Type annotations (var name:type = value)
  - [x] Nullable types (var name:type? = value)
  - [x] Support for string, int, boolean, and null literals
- [ ] Function declarations
  - [ ] Parameters
  - [ ] Return type
  - [ ] Function body
- [ ] Class declarations
  - [ ] Properties
  - [ ] Methods
  - [ ] Inheritance
- [ ] Interface declarations
  - [ ] Method signatures
  - [ ] Property signatures
- [ ] Control flow
  - [ ] If statements
  - [ ] For loops
  - [ ] While loops
  - [ ] When expressions
  - [ ] Return statements

## Expression Types
- [x] Literals
  - [x] String literals
  - [x] Integer literals
  - [x] Boolean literals
  - [x] Null literal
  - [ ] Float literals
- [ ] Binary expressions
  - [ ] Arithmetic operators (+, -, *, /, %)
  - [ ] Comparison operators (==, !=, <, >, <=, >=)
  - [ ] Logical operators (and, or)
- [ ] Unary expressions
  - [ ] Negation (-)
  - [ ] Not (!)
- [ ] Member access (obj.prop)
- [ ] Function calls
- [ ] Array literals and access
- [ ] Map literals and access
- [ ] Lambda expressions

## Error Handling
- [x] Basic error recovery
- [x] Source location tracking
- [x] Error synchronization
- [ ] Multiple error reporting
- [ ] Better error messages
  - [ ] Expected vs actual token
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
- [ ] Function declaration tests
- [ ] Class declaration tests
- [ ] Interface declaration tests
- [ ] Expression tests
- [ ] Error recovery tests
- [ ] Edge case tests
