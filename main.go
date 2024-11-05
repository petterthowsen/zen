package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"zen/interpreter"
	"zen/lang/common"
	"zen/lang/lexing"
	"zen/lang/parsing"
)

var DEBUG bool
var REPL bool
var DEBUG_TOKENS bool
var DEBUG_PARSE_TREE bool
var VERBOSE bool

// Main Executes zen code from stdin, from a file ending with .zen or starts the REPL.
func main() {
	flag.BoolVar(&DEBUG, "debug", false, "enable debug mode")

	flag.BoolVar(&REPL, "interactive", false, "Run a Read-Eval-Print Loop (REPL)")
	flag.BoolVar(&REPL, "i", false, "enable interactive mode (REPL)")

	flag.BoolVar(&DEBUG_TOKENS, "tokens", false, "show tokenized input")
	flag.BoolVar(&DEBUG_PARSE_TREE, "ast", false, "show parse tree")
	flag.BoolVar(&VERBOSE, "verbose", false, "show tokens and parse tree")
	flag.BoolVar(&VERBOSE, "v", false, "show tokens and parse tree")
	flag.Parse()

	if VERBOSE {
		DEBUG_TOKENS = true
		DEBUG_PARSE_TREE = true
	}

	if DEBUG_TOKENS {
		println("Debugging Tokens Enabled")
	}

	// Check if stdin has input or if we should start REPL
	stat, _ := os.Stdin.Stat()
	if stat.Mode()&os.ModeCharDevice == 0 {
		// Data is being piped into stdin
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println("Error reading input:", err)
			os.Exit(1)
		}

		var sourceCode common.SourceCode

		// Check if input is a file path or code
		if strings.HasSuffix(string(input), ".zen") {
			sourceCode = loadFileSource(strings.TrimSpace(string(input)))
		} else {
			sourceCode = common.NewInlineSourceCode(string(input))
		}
		execute(sourceCode)
	} else if REPL {
		// Start REPL if no stdin input and interactive flag is set
		startREPL()
	} else {
		fmt.Println("No input detected. Use -interactive for REPL or provide code.")
	}
}

func loadFileSource(filename string) common.SourceCode {

	code, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	return common.NewFileSourceCode(filename, string(code))
}

func execute(code common.SourceCode) {
	// tokenize
	lexer := lexing.NewLexer(code)

	tokens, err := lexer.Scan()
	if err != nil {
		if len(lexer.Errors) > 0 {
			syntaxErrors := make([]*common.SyntaxError, len(lexer.Errors))
			for i, err := range lexer.Errors {
				syntaxErrors[i] = &err
			}
			printSyntaxErrors(syntaxErrors)
		}
	}

	if DEBUG_TOKENS {
		fmt.Println("Scanned tokens:")
		for _, token := range tokens {
			fmt.Println(token.String())
		}
	}

	// stop?
	if err != nil {
		return
	}

	// parse
	parser := parsing.NewParser(tokens, true)
	program, syntaxErrors := parser.Parse()

	if DEBUG_PARSE_TREE {
		if DEBUG_TOKENS {
			fmt.Println("-----------------------------------------")
		}
		fmt.Println("AST:")
		fmt.Println(program.String(2))
	}

	if len(syntaxErrors) > 0 {
		printSyntaxErrors(syntaxErrors)
	} else {
		// execute
		i := interpreter.NewInterpreter()
		if err := i.Execute(program); err != nil {
			fmt.Println("Interpreter error:", err)
		}
	}

}

// printSyntaxErrors prints the list of syntax errors to the standard output.
// Each error is printed using the Error method of the common.SyntaxError type.
func printSyntaxErrors(errors []*common.SyntaxError) {
	fmt.Println("Whoops! Syntax Error(s):")

	for _, err := range errors {
		fmt.Println(err.Error())
	}
}

// startREPL initializes and starts a Read-Eval-Print Loop (REPL) environment allowing interactive code execution.
func startREPL() {
	fmt.Println("Zen REPL:")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(">> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		if input == "exit" || input == "quit" {
			break
		}

		sourceCode := common.NewInlineSourceCode(input)
		execute(sourceCode)
	}
}
