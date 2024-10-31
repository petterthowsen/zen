package parsing

// Node represents a statement in the parse tree
type Node interface {
	// String returns the string representation of the Node for debugging purposes
	String() string
}
