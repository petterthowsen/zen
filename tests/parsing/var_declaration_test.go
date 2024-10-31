package parsing_test

import (
	"os"
	"strings"
	"testing"
	"zen/lang/parsing/statement"
)

func TestVarDeclarationFile(t *testing.T) {
	// Load and parse the test file
	path := getTestDataPath("var_declaration.zen")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", path, err)
	}

	program, errors := ParseFile(path, string(content))
	if len(errors) > 0 {
		t.Errorf("Parser errors:")
		for _, err := range errors {
			t.Errorf("  %s", err.Error())
		}
		return
	}

	if program == nil {
		t.Error("Expected program node, got nil")
		return
	}

	// Save AST to file
	astStr := program.String(0)
	astPath := strings.TrimSuffix(path, ".zen") + ".ast"
	if err := os.WriteFile(astPath, []byte(astStr), 0644); err != nil {
		t.Logf("Warning: Failed to write AST file: %v", err)
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

		// Type assert to VarDeclarationNode
		varDecl, ok := stmt.(*statement.VarDeclarationNode)
		if !ok {
			t.Errorf("Statement %d: expected VarDeclarationNode, got %T", i, stmt)
			continue
		}

		if varDecl.Name != exp.name {
			t.Errorf("Declaration %d: expected name %q, got %q", i, exp.name, varDecl.Name)
		}

		expectedType := exp.typ
		if exp.nullable {
			expectedType += "?"
		}
		if varDecl.Type != expectedType {
			t.Errorf("Declaration %d: expected type %q, got %q", i, expectedType, varDecl.Type)
		}

		if varDecl.IsConstant != exp.isConst {
			t.Errorf("Declaration %d: expected const=%v, got %v", i, exp.isConst, varDecl.IsConstant)
		}

		// TODO: Add initializer checks once expression parsing is implemented
	}
}

func TestParserError(t *testing.T) {
	program, errors := ParseString("var")
	if len(errors) == 0 {
		t.Error("Expected parsing error, got none")
	}
	if program != nil {
		t.Error("Expected nil program for error case")
	}
}
