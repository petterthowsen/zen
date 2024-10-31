package lexing

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"zen/lang/common"
	"zen/lang/lexing"
)

// TokenAssert holds expected token information for testing
type TokenAssert struct {
	Type    lexing.TokenType
	Literal string
}

func verifyTokens(t *testing.T, tokens []lexing.Token, expected []TokenAssert) bool {
	// Verify token count (including EOF)
	if len(tokens) != len(expected)+1 {
		t.Errorf("Expected %d tokens, got %d", len(expected)+1, len(tokens))
		t.Logf("First few tokens received:")
		for i := 0; i < min(3, len(tokens)); i++ {
			t.Logf("%d: %s", i, tokenString(tokens[i]))
		}
	}

	// Verify each token
	for i, exp := range expected {
		actual := tokens[i]
		if actual.Type != exp.Type || actual.Literal != exp.Literal {
			t.Errorf("Token mismatch at position %d:\nExpected: Type=%v, Literal='%s'\nGot:      %s",
				i, lexing.TokenName(exp.Type), exp.Literal, tokenString(actual))
		}
	}

	// Verify EOF is last token
	lastToken := tokens[len(tokens)-1]
	if lastToken.Type != lexing.EOF {
		t.Errorf("Expected last token to be EOF, got %v", lastToken.Type)
		return false
	}

	return true
}

// AssertTokens verifies that the lexer produces the expected sequence of tokens
func AssertTokens(t *testing.T, source string, expected []TokenAssert) {
	sourceCode := common.NewInlineSourceCode(source)
	lexer := lexing.NewLexer(sourceCode)
	tokens, err := lexer.Scan()

	if err != nil {
		t.Log("Lexer errors:")
		for _, err := range lexer.Errors {
			t.Log(err.Error())
		}
		t.FailNow()
	}

	verifyTokens(t, tokens, expected)
}

// getTestDataPath returns the absolute path to the test data file
func getTestDataPath(filename string) string {
	_, currentFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(currentFile)
	return filepath.Join(dir, filename)
}

// LoadAndAssertTokens loads a .zen file and verifies its tokens
func LoadAndAssertTokens(t *testing.T, filename string, expected []TokenAssert) {
	path := getTestDataPath(filename)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", path, err)
	}

	sourceCode := common.NewFileSourceCode(path, string(content))
	lexer := lexing.NewLexer(sourceCode)
	tokens, err := lexer.Scan()

	if err != nil {
		t.Log("Lexer errors:")
		for _, err := range lexer.Errors {
			t.Log(err.Error())
		}
		t.FailNow()
	}

	// Save tokens to .tokens file for any .zen file
	if strings.HasSuffix(filename, ".zen") {
		var tokenStrings []string
		for i, tok := range tokens {
			tokenStrings = append(tokenStrings, fmt.Sprintf("%d: %s", i, tokenString(tok)))
		}

		// Replace .zen extension with .tokens
		tokensPath := strings.TrimSuffix(path, ".zen") + ".tokens"
		if err := os.WriteFile(tokensPath, []byte(strings.Join(tokenStrings, "\n")), 0644); err != nil {
			t.Logf("Warning: Failed to write tokens file: %v", err)
		}
	}

	if !verifyTokens(t, tokens, expected) {
		t.Log("\nFile contents:")
		t.Log(string(content))
	}
}

// Helper function to print token type and literal
func tokenString(t lexing.Token) string {
	return fmt.Sprintf("Type=%v, Literal='%s'", lexing.TokenName(t.Type), t.Literal)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
