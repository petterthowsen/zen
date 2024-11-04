package parsing

import (
	"testing"
	"zen/lang/parsing/ast"
	"zen/lang/parsing/expression"
)

func TestParametricTypes(t *testing.T) {
	program := ParseTestFile(t, "parametric_types.zen")
	if program == nil {
		return
	}

	// Basic type annotations
	AssertVarDeclarationWithType(t, program.Statements[0], "a", false, false,
		func(t *testing.T, typ ast.Expression) {
			AssertBasicType(t, typ, "int")
		})

	AssertVarDeclarationWithType(t, program.Statements[1], "b", false, false,
		func(t *testing.T, typ ast.Expression) {
			AssertBasicType(t, typ, "string")
		})

	AssertVarDeclarationWithType(t, program.Statements[2], "c", false, false,
		func(t *testing.T, typ ast.Expression) {
			AssertBasicType(t, typ, "float64")
		})

	// Nullable types
	AssertVarDeclarationWithType(t, program.Statements[3], "d", false, true,
		func(t *testing.T, typ ast.Expression) {
			AssertBasicType(t, typ, "int")
		})

	AssertVarDeclarationWithType(t, program.Statements[4], "e", false, true,
		func(t *testing.T, typ ast.Expression) {
			AssertBasicType(t, typ, "string")
		})

	// Array types with size
	AssertVarDeclarationWithType(t, program.Statements[5], "arr1", false, false,
		func(t *testing.T, typ ast.Expression) {
			paramType := AssertParametricType(t, typ, "Array", 2)
			if paramType == nil {
				return
			}
			AssertTypeParameter(t, paramType.Parameters[0], "int")
			AssertValueParameter(t, paramType.Parameters[1], int64(5))
		})

	AssertVarDeclarationWithType(t, program.Statements[6], "arr2", false, false,
		func(t *testing.T, typ ast.Expression) {
			paramType := AssertParametricType(t, typ, "Array", 2)
			if paramType == nil {
				return
			}
			AssertTypeParameter(t, paramType.Parameters[0], "string")
			AssertValueParameter(t, paramType.Parameters[1], int64(10))
		})

	// Array types with type only
	AssertVarDeclarationWithType(t, program.Statements[7], "nums", false, false,
		func(t *testing.T, typ ast.Expression) {
			paramType := AssertParametricType(t, typ, "Array", 1)
			if paramType == nil {
				return
			}
			AssertTypeParameter(t, paramType.Parameters[0], "int")
		})

	AssertVarDeclarationWithType(t, program.Statements[8], "strs", false, false,
		func(t *testing.T, typ ast.Expression) {
			paramType := AssertParametricType(t, typ, "Array", 1)
			if paramType == nil {
				return
			}
			AssertTypeParameter(t, paramType.Parameters[0], "string")
		})

	// Map types
	AssertVarDeclarationWithType(t, program.Statements[9], "scores", false, false,
		func(t *testing.T, typ ast.Expression) {
			paramType := AssertParametricType(t, typ, "Map", 2)
			if paramType == nil {
				return
			}
			AssertTypeParameter(t, paramType.Parameters[0], "string")
			AssertTypeParameter(t, paramType.Parameters[1], "int")
		})

	AssertVarDeclarationWithType(t, program.Statements[10], "config", false, false,
		func(t *testing.T, typ ast.Expression) {
			paramType := AssertParametricType(t, typ, "Map", 2)
			if paramType == nil {
				return
			}
			AssertTypeParameter(t, paramType.Parameters[0], "string")
			AssertTypeParameter(t, paramType.Parameters[1], "any")
		})

	// Nested parametric types
	AssertVarDeclarationWithType(t, program.Statements[11], "matrix", false, false,
		func(t *testing.T, typ ast.Expression) {
			paramType := AssertParametricType(t, typ, "Array", 2)
			if paramType == nil {
				return
			}

			// First parameter is Array<int, 3>
			innerType := paramType.Parameters[0].Value.(*expression.ParametricType)
			if innerType.BaseType != "Array" {
				t.Errorf("Expected inner type 'Array', got '%s'", innerType.BaseType)
				return
			}
			AssertTypeParameter(t, innerType.Parameters[0], "int")
			AssertValueParameter(t, innerType.Parameters[1], int64(3))

			// Second parameter is size 2
			AssertValueParameter(t, paramType.Parameters[1], int64(2))
		})

	AssertVarDeclarationWithType(t, program.Statements[12], "nested", false, false,
		func(t *testing.T, typ ast.Expression) {
			paramType := AssertParametricType(t, typ, "Map", 2)
			if paramType == nil {
				return
			}

			// First parameter is string
			AssertTypeParameter(t, paramType.Parameters[0], "string")

			// Second parameter is Array<int, 5>
			innerType := paramType.Parameters[1].Value.(*expression.ParametricType)
			if innerType.BaseType != "Array" {
				t.Errorf("Expected inner type 'Array', got '%s'", innerType.BaseType)
				return
			}
			AssertTypeParameter(t, innerType.Parameters[0], "int")
			AssertValueParameter(t, innerType.Parameters[1], int64(5))
		})

	// Multiple parameters
	AssertVarDeclarationWithType(t, program.Statements[13], "grid", false, false,
		func(t *testing.T, typ ast.Expression) {
			paramType := AssertParametricType(t, typ, "Grid", 3)
			if paramType == nil {
				return
			}
			AssertTypeParameter(t, paramType.Parameters[0], "int")
			AssertValueParameter(t, paramType.Parameters[1], int64(3))
			AssertValueParameter(t, paramType.Parameters[2], int64(4))
		})

	AssertVarDeclarationWithType(t, program.Statements[14], "cube", false, false,
		func(t *testing.T, typ ast.Expression) {
			paramType := AssertParametricType(t, typ, "Grid", 4)
			if paramType == nil {
				return
			}
			AssertTypeParameter(t, paramType.Parameters[0], "float64")
			AssertValueParameter(t, paramType.Parameters[1], int64(2))
			AssertValueParameter(t, paramType.Parameters[2], int64(2))
			AssertValueParameter(t, paramType.Parameters[3], int64(2))
		})
}

func TestParametricTypeErrors(t *testing.T) {
	cases := []string{
		"var x : Array<>",        // Empty type parameters
		"var x : Array<int,>",    // Trailing comma
		"var x : Array<,int>",    // Leading comma
		"var x : Array<int int>", // Missing comma
		"var x : Grid<int, 2, 3", // Missing closing bracket
		"var x : Array<int 5>",   // Missing comma between parameters
	}

	for _, input := range cases {
		AssertParseError(t, input)
	}
}
