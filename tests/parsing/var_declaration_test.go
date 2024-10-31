package parsing_test

import (
	"os"
	"strings"
	"testing"
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
	if len(program.Statements) == 0 {
		t.Error("Expected statements in program, got none")
		return
	}

	// TODO: Add more specific AST structure verification
}

func TestParserError(t *testing.T) {
	program, errors := ParseString("var")
	if len(errors) == 0 {
		t.Error("Expected parsing error, got none")
	}
	if program != nil && len(program.Statements) > 0 {
		t.Error("Expected no statements for error case")
	}
}
