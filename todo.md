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
  - [x] For loops
  - [x] While loops
  - [ ] When statements
  - [x] Return statements
- [ ] Exceptions
  - [ ] Throw statements
  - [ ] Try / Catch statements

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
- [x] Member access (obj.prop)
- [x] Function calls
- [x] await expressions
- [x] Array literals
  - [ ] Bracket access (array[0])
- [x] Map literals
  - [ ] Curly access (map{"key"})
- [ ] Tuple Literals (("status": 404, "body": "not found"))
  - [ ] Tuple destructuring (var status, body <= resultTuple)
  - [ ] Tuple accessors just use regular object access
  - [ ] Tuple inference in functions (func() : status: int, body: string) inferred as Tuple<int, string> with names "status" and "body"
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
- [x] Generic types
  - [x] Simple parametric types (Array<int>)
  - [x] Multiple parameters (Grid<int, 3, 4>)
  - [x] Nested types (Array<Array<int, 3>, 2>)
  - [x] Mixed type and value parameters (Array<string, 10>)
- [ ] Union types
- [ ] Nullable type handling

## Type checking (Should be done after parsing stage)
- [ ] Type compatibility rules
- [ ] Basic type checking
- [ ] Type inference

## Testing
- [x] Variable declaration tests
- [x] Function declaration tests
- [ ] Class declaration tests
- [ ] Interface declaration tests
- [x] Expression tests
- [ ] Error recovery tests
- [ ] Edge case tests
