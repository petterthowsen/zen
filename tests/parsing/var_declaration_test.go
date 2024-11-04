package parsing

import (
	"testing"
	"zen/lang/parsing/ast"
)

func TestVarDeclarationFile(t *testing.T) {
	program := ParseTestFile(t, "var_declaration.zen")
	if program == nil {
		return
	}

	// Verify basic structure
	expectedDecls := []struct {
		name     string
		typ      string // empty string means no type annotation
		isConst  bool
		nullable bool
	}{
		{"name", "", false, false},
		{"age", "", false, false},
		{"isValid", "", false, false},
		{"title", "string", false, false},
		{"count", "int", false, false},
		{"enabled", "bool", false, false},
		{"description", "string", false, true},
		{"quantity", "int", false, true},
	}

	if len(program.Statements) != len(expectedDecls) {
		t.Errorf("Expected %d declarations, got %d", len(expectedDecls), len(program.Statements))
		return
	}

	for i, exp := range expectedDecls {
		stmt := program.Statements[i]
		if exp.typ == "" {
			// No type annotation
			AssertVarDeclaration(t, stmt, exp.name, exp.isConst, exp.nullable)
		} else {
			// Has type annotation
			AssertVarDeclarationWithType(t, stmt, exp.name, exp.isConst, exp.nullable,
				func(t *testing.T, typ ast.Expression) {
					AssertBasicType(t, typ, exp.typ)
				})
		}
	}
}

func TestParserError(t *testing.T) {
	AssertParseError(t, "var")
}
