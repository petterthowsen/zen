package types

import "fmt"

// Type represents the type of a Value
type Type int

const (
	TypeInt   Type = iota // 32 bits
	TypeFloat             // 32 bits
	TypeInt64
	TypeFloat64

	TypeString

	TypeBool

	TypeNull
	TypeVoid

	// TypeFunction denotes a named or an anonymous function
	TypeFunction
	TypeBuiltinFunction

	// TypeLambda denotes an anonymous closure
	TypeLambda

	// TypeClass denotes a class
	TypeClass

	// TypeObject denotes an object instance (e.g., an instance of a class)
	TypeObject
)

// String returns the string representation of a Type
func (t Type) String() string {
	switch t {
	case TypeInt:
		return "int"
	case TypeFloat:
		return "float"
	case TypeInt64:
		return "int64"
	case TypeFloat64:
		return "float64"
	case TypeString:
		return "string"
	case TypeBool:
		return "bool"
	case TypeNull:
		return "null"
	case TypeVoid:
		return "void"
	case TypeFunction:
		return "function"
	case TypeLambda:
		return "lambda"
	case TypeClass:
		return "class"
	case TypeObject:
		return "object"
	default:
		return "unknown"
	}
}

// Value represents any value in the Zen language
type Value interface {
	// Type returns the Type of the value
	Type() Type

	// String returns a string representation of the value
	String() string

	// Equals returns true if this value equals another value
	Equals(other Value) bool

	// IsTruthy returns true if this value is considered true in a boolean context
	IsTruthy() bool

	// Clone returns a copy of this value: For primitives, it's a deep copy
	// For objects, it's a shallow copy
	Clone() Value
}

// TypeError represents an error during type operations
type TypeError struct {
	Message string
}

func (e *TypeError) Error() string {
	return fmt.Sprintf("Type error: %s", e.Message)
}

// NewTypeError creates a new TypeError with the given message
func NewTypeError(format string, args ...interface{}) error {
	return &TypeError{Message: fmt.Sprintf(format, args...)}
}

// IsNumeric returns true if the given type is numeric (int, float, int64, float64)
func IsNumeric(t Type) bool {
	switch t {
	case TypeInt, TypeFloat, TypeInt64, TypeFloat64:
		return true
	default:
		return false
	}
}

// CanCoerce returns true if a value of type 'from' can be coerced to type 'to'
func CanCoerce(from, to Type) bool {
	// Same type is always coercible
	if from == to {
		return true
	}

	// Numeric types can be coerced between each other
	if IsNumeric(from) && IsNumeric(to) {
		return true
	}

	// Special cases
	switch from {
	case TypeNull:
		// Null can be coerced to any type
		return true
	case TypeInt:
		// Int can be coerced to float types or string
		return to == TypeFloat || to == TypeFloat64 || to == TypeString
	case TypeFloat:
		// Float can be coerced to float64 or string
		return to == TypeFloat64 || to == TypeString
	case TypeInt64:
		// Int64 can be coerced to float64 or string
		return to == TypeFloat64 || to == TypeString
	case TypeFloat64:
		// Float64 can only be coerced to string
		return to == TypeString
	case TypeString:
		// String can be coerced to numeric types if it contains a valid number
		return IsNumeric(to)
	case TypeBool:
		// Bool can be coerced to string or numeric types
		return to == TypeString || IsNumeric(to)
	}

	return false
}

// IsValidBinaryOp returns true if the given types can be used in a binary operation
func IsValidBinaryOp(left, right Type, op string) bool {
	switch op {
	case "+":
		// Addition works with numeric types and strings
		if left == TypeString && right == TypeString {
			return true
		}
		return IsNumeric(left) && IsNumeric(right)
	case "-", "*", "/":
		// Arithmetic operations only work with numeric types
		return IsNumeric(left) && IsNumeric(right)
	case "<", "<=", ">", ">=":
		// Comparison works with numeric types and strings
		if left == TypeString && right == TypeString {
			return true
		}
		return IsNumeric(left) && IsNumeric(right)
	case "==", "!=":
		// Equality works with any types
		return true
	case "and", "or":
		// Logical operations only work with booleans
		return left == TypeBool && right == TypeBool
	default:
		return false
	}
}

// IsValidUnaryOp returns true if the given type can be used in a unary operation
func IsValidUnaryOp(t Type, op string) bool {
	switch op {
	case "-":
		return IsNumeric(t)
	case "not":
		return t == TypeBool
	default:
		return false
	}
}
