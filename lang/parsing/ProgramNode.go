package parsing

import "zen/lang/lexing"

// ProgramNode represents the root node of the parse tree
type ProgramNode struct {
	Tokens []*lexing.Token
}
