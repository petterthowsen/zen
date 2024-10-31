package parsing_test

import (
	"path/filepath"
	"runtime"
	"zen/lang/common"
	"zen/lang/lexing"
	"zen/lang/parsing"
	"zen/lang/parsing/ast"
)

// ParseString parses source code string and returns the AST
func ParseString(source string) (*ast.ProgramNode, []*common.SyntaxError) {
	sourceCode := common.NewInlineSourceCode(source)
	lexer := lexing.NewLexer(sourceCode)
	tokens, err := lexer.Scan()
	if err != nil {
		return nil, []*common.SyntaxError{err.(*common.SyntaxError)}
	}

	parser := parsing.NewParser(tokens, false)
	return parser.Parse()
}

// ParseFile parses a source file and returns the AST
func ParseFile(path string, content string) (*ast.ProgramNode, []*common.SyntaxError) {
	sourceCode := common.NewFileSourceCode(path, content)
	lexer := lexing.NewLexer(sourceCode)
	tokens, err := lexer.Scan()
	if err != nil {
		return nil, []*common.SyntaxError{err.(*common.SyntaxError)}
	}

	parser := parsing.NewParser(tokens, false)
	return parser.Parse()
}

// getTestDataPath returns the absolute path to the test data file
func getTestDataPath(filename string) string {
	_, currentFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(currentFile)
	return filepath.Join(dir, filename)
}
