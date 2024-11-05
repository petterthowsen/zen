package types

// Callable represents either a function, a lambda or a class method
type Callable interface {
	IsCallable() bool
}
