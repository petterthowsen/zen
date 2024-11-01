package parsing

import (
	"testing"
)

func TestVarDeclarationFile(t *testing.T) {
	program := ParseTestFile(t, "var_declaration.zen")
	if program == nil {
		return
	}

	// Verify basic structure
	expectedDecls := []struct {
		name     string
		typ      string
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
		t.Logf("Checking declaration %d: %s", i, stmt.String(0))
		AssertVarDeclaration(t, stmt, exp.name, exp.typ, exp.isConst, exp.nullable)
	}
}

func TestParserError(t *testing.T) {
	AssertParseError(t, "var")
}
