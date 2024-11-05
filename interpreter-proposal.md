# Proposal for Execution System
The following is a very general and loosely sketched proposal for the interpreter/runtime system for Zen.

It should:
- Run in a single-threaded event loop and supports the `async` and `await` keywords to handle asynchronous programming.
- Support exceptions with `throw`, `try`, `catch` statements.
- Execution "environments" or "scopes" that can be created and forked from parent environments (functions have their own scope, as do for loops etc. but might have access to their parent scope in some contexts.)
- Error reporting with stack traces
- Robust type system with support for parametric types (generics such as `Array<int, 5>`)
- Flexible namespace/import system

At some point, we should consider adding support for true concurrency, probably via goroutines.

## Type system
Zen is strictly typed, with type inference. We also need to support type casting/conversion.
Probably by providing built in functions `string()`, `ìnt()` etc.

Zen's type system has a number of primitive types supported:
- int (signed 32 bit integer)
- int64 (signed 64 bit integer)
- float (signed 32 bit float)
- float64 (signed 64 bit float)
- bool
- null

Additionally, the 'any' type-hint denotes an untyped/mixed-type variable.

Nullability is enabled with QMARK token, I.E `var name = string?`

### Tokens
These are supported by the lexer's:
KEYWORD tokens for type-hints:
- these KEYWORD tokens (token.literal): "int", "int64", "float", "float64", "string", "bool"

and these KEYWORD tokens for literal values:
- "true", "false", "null", "void"

and number and string literals have their own tokens types:
- STRING for "some strings"
- INT for "123" or "42"
- FLOAT for "3.14" or "-9.0"

Other types include Array, Map and any other custom classes.

### Parsing

At the AST level, types are represented by either the BasicType or ParametricType nodes.

# Semantic analysis, Type-checking
After parsing, we should run through the AST and:
- perform type inference
- symbol resolution and import resolving
- constant folding (5+3 becomes 7 etc.)
- error reporting based on the above steps

# Namespaces, Built ins & interfacing with Go 
Zen's runtime/execution system should support the ability to easily register built in functions at the Go level that become available in Zen, allowing interfacing with Go at a low level and handle namespacing and importing of other zen code.

Off the top of my head, something like this:
```go
printFunc := runtime.createFunction("print")
printFunc.addVariadicParameter()
printFunc.setFunction(func(runtime *runtime, ...params))

runtime.registerFunction("sys", printFunc)
```

And similarly for packages, namespaces etc.

This system would need to be thought out properly in detail.

## Namespaces & Imports
```go (actually zen but no syntax highlighting for that)
//-----------------------------------//
// Namespaces, Packages, and Imports
//-----------------------------------//
// Zen organizes code into namespaces based on folder structure.
// - Folders represent namespaces.
// - Each `.zen` file in a given namespace is a "module".
// - Modules can define one or more symbols (functions, classes, etc.).
// - Zen projects define their root namespace through a `package.zen` file.
// - When importing, Zen automatically treats modules as namespaces if they
//   contain more than one symbol, simplifying import syntax and scoping.
//
// Note: Zen ignores any files starting with an underscore.
// This is to allow developers to create .zen files without in their package for other purposes.

//-----------------------------------
// Example Directory Structure:

// GameEngine/
// - package.zen
// - Core/
// -- Vector2.zen
// -- Utils.zen
// -- Nested/
// --- Special.zen

// where `package.zen` contains
package GameEngine

// where `Utils.zen` contains:
func hello() {}
func world() {}
 
// where `Vector2.zen` contains:
class Vector2 {
    x: int
    y: int
}

// where `Special.zen` contains:
func specialFunc() {}

//-----------------------------------
// Import Syntax

// Zen supports multiple ways to import modules and symbols within namespaces, 
// with automatic handling for single-symbol vs. multi-symbol modules:

// 1. Import Namespace: Imports all symbols from all modules directly under
//    the specified namespace if the module has multiple symbols.

import GameEngine/Core
hello()         // Available directly
world()         // Available directly
var vec2 = Vector2() // Only if Vector2 is defined as part of GameEngine/Core

// 2. Import Namespace with Scoped Alias: Imports all symbols but keeps them
//    scoped under an alias.

import GameEngine/Core as gec
gec.hello()     // Scoped access
gec.world()     // Scoped access

// 3. Importing a Single-symbol Module Directly: If a module defines only
//    a single symbol, it is imported directly by default.

import GameEngine/Core/Vector2
var vec2 = Vector2() // Direct access without specifying symbol name

// 4. Importing a Single-symbol Module with Alias: You can alias a single-symbol
//    module for clarity or disambiguation.

import GameEngine/Core/Vector2 as CoreVector2
var vec2 = CoreVector2()

// 5. Importing Specific Symbols from a Multi-symbol Module: If a module defines 
//    multiple symbols, you can selectively import symbols using the `from` syntax.

from GameEngine/Core/Utils import hello
hello()        // Directly accessible
world()        // Error: not imported

// Notes:
// - Modules with more than one symbol are treated as namespaces, 
//   allowing flexible import and scoping.
// - Single-symbol modules are directly accessible upon import, 
//   reducing verbosity and improving readability.
// - Aliasing with `as` works for both single-symbol and multi-symbol modules.
```

## Exposing low-level Go-related functions to zen

Much of the barebones functionality such as I/O will be supported by custom implementations that bridge Zen -> Go.
However, we might expose a set of low-level functions that will allow Zen developers to interface more directly with Go, for the purpose of making some Go libraries available in Zen. I believe this is called "Bindings".

We may also provide Reflection capabilities and the like.