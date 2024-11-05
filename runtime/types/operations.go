package types

// BinaryOp performs a binary operation on two values
func BinaryOp(left, right Value, op string) (Value, error) {
	// First check if the operation is valid for these types
	if !IsValidBinaryOp(left.Type(), right.Type(), op) {
		return nil, NewTypeError("invalid operation: %s %s %s", left.Type(), op, right.Type())
	}

	// For logical operators, both operands must be boolean
	if op == "and" || op == "or" {
		if left.Type() != TypeBool || right.Type() != TypeBool {
			return nil, NewTypeError("logical operators require boolean operands, got %s and %s", left.Type(), right.Type())
		}
		if op == "and" {
			return NewBool(left.(*Bool).Value() && right.(*Bool).Value()), nil
		}
		return NewBool(left.(*Bool).Value() || right.(*Bool).Value()), nil
	}

	// For other operators, try to coerce operands to compatible types
	l, r, err := CoerceForOperation(left, right, op)
	if err != nil {
		return nil, err
	}

	switch op {
	case "+":
		return add(l, r)
	case "-":
		return subtract(l, r)
	case "*":
		return multiply(l, r)
	case "/":
		return divide(l, r)
	case "<":
		return lessThan(l, r)
	case "<=":
		return lessEqual(l, r)
	case ">":
		return greaterThan(l, r)
	case ">=":
		return greaterEqual(l, r)
	case "==":
		return equals(l, r)
	case "!=":
		result, err := equals(l, r)
		if err != nil {
			return nil, err
		}
		return NewBool(!result.(*Bool).Value()), nil
	default:
		return nil, NewTypeError("unknown operator %s", op)
	}
}

// UnaryOp performs a unary operation on a value
func UnaryOp(v Value, op string) (Value, error) {
	// First check if the operation is valid for this type
	if !IsValidUnaryOp(v.Type(), op) {
		return nil, NewTypeError("invalid operation: %s %s", op, v.Type())
	}

	switch op {
	case "-":
		return negate(v)
	case "not":
		if v.Type() != TypeBool {
			return nil, NewTypeError("logical NOT requires boolean operand, got %s", v.Type())
		}
		return NewBool(!v.(*Bool).Value()), nil
	default:
		return nil, NewTypeError("unknown unary operator %s", op)
	}
}

// Helper functions for operations

func add(l, r Value) (Value, error) {
	if l.Type() == TypeString && r.Type() == TypeString {
		return NewString(l.(*String).Value() + r.(*String).Value()), nil
	}

	switch l.Type() {
	case TypeInt:
		return NewInt(l.(*Int).Value() + r.(*Int).Value()), nil
	case TypeFloat:
		return NewFloat(l.(*Float).Value() + r.(*Float).Value()), nil
	case TypeInt64:
		return NewInt64(l.(*Int64).Value() + r.(*Int64).Value()), nil
	case TypeFloat64:
		return NewFloat64(l.(*Float64).Value() + r.(*Float64).Value()), nil
	default:
		return nil, NewTypeError("invalid types for addition: %s and %s", l.Type(), r.Type())
	}
}

func subtract(l, r Value) (Value, error) {
	switch l.Type() {
	case TypeInt:
		return NewInt(l.(*Int).Value() - r.(*Int).Value()), nil
	case TypeFloat:
		return NewFloat(l.(*Float).Value() - r.(*Float).Value()), nil
	case TypeInt64:
		return NewInt64(l.(*Int64).Value() - r.(*Int64).Value()), nil
	case TypeFloat64:
		return NewFloat64(l.(*Float64).Value() - r.(*Float64).Value()), nil
	default:
		return nil, NewTypeError("invalid types for subtraction: %s and %s", l.Type(), r.Type())
	}
}

func multiply(l, r Value) (Value, error) {
	switch l.Type() {
	case TypeInt:
		return NewInt(l.(*Int).Value() * r.(*Int).Value()), nil
	case TypeFloat:
		return NewFloat(l.(*Float).Value() * r.(*Float).Value()), nil
	case TypeInt64:
		return NewInt64(l.(*Int64).Value() * r.(*Int64).Value()), nil
	case TypeFloat64:
		return NewFloat64(l.(*Float64).Value() * r.(*Float64).Value()), nil
	default:
		return nil, NewTypeError("invalid types for multiplication: %s and %s", l.Type(), r.Type())
	}
}

func divide(l, r Value) (Value, error) {
	// Check for division by zero
	if !r.IsTruthy() {
		return nil, NewTypeError("division by zero")
	}

	switch l.Type() {
	case TypeInt:
		return NewInt(l.(*Int).Value() / r.(*Int).Value()), nil
	case TypeFloat:
		return NewFloat(l.(*Float).Value() / r.(*Float).Value()), nil
	case TypeInt64:
		return NewInt64(l.(*Int64).Value() / r.(*Int64).Value()), nil
	case TypeFloat64:
		return NewFloat64(l.(*Float64).Value() / r.(*Float64).Value()), nil
	default:
		return nil, NewTypeError("invalid types for division: %s and %s", l.Type(), r.Type())
	}
}

func lessThan(l, r Value) (Value, error) {
	switch l.Type() {
	case TypeInt:
		return NewBool(l.(*Int).Value() < r.(*Int).Value()), nil
	case TypeFloat:
		return NewBool(l.(*Float).Value() < r.(*Float).Value()), nil
	case TypeInt64:
		return NewBool(l.(*Int64).Value() < r.(*Int64).Value()), nil
	case TypeFloat64:
		return NewBool(l.(*Float64).Value() < r.(*Float64).Value()), nil
	case TypeString:
		return NewBool(l.(*String).Value() < r.(*String).Value()), nil
	default:
		return nil, NewTypeError("invalid types for comparison: %s and %s", l.Type(), r.Type())
	}
}

func lessEqual(l, r Value) (Value, error) {
	switch l.Type() {
	case TypeInt:
		return NewBool(l.(*Int).Value() <= r.(*Int).Value()), nil
	case TypeFloat:
		return NewBool(l.(*Float).Value() <= r.(*Float).Value()), nil
	case TypeInt64:
		return NewBool(l.(*Int64).Value() <= r.(*Int64).Value()), nil
	case TypeFloat64:
		return NewBool(l.(*Float64).Value() <= r.(*Float64).Value()), nil
	case TypeString:
		return NewBool(l.(*String).Value() <= r.(*String).Value()), nil
	default:
		return nil, NewTypeError("invalid types for comparison: %s and %s", l.Type(), r.Type())
	}
}

func greaterThan(l, r Value) (Value, error) {
	switch l.Type() {
	case TypeInt:
		return NewBool(l.(*Int).Value() > r.(*Int).Value()), nil
	case TypeFloat:
		return NewBool(l.(*Float).Value() > r.(*Float).Value()), nil
	case TypeInt64:
		return NewBool(l.(*Int64).Value() > r.(*Int64).Value()), nil
	case TypeFloat64:
		return NewBool(l.(*Float64).Value() > r.(*Float64).Value()), nil
	case TypeString:
		return NewBool(l.(*String).Value() > r.(*String).Value()), nil
	default:
		return nil, NewTypeError("invalid types for comparison: %s and %s", l.Type(), r.Type())
	}
}

func greaterEqual(l, r Value) (Value, error) {
	switch l.Type() {
	case TypeInt:
		return NewBool(l.(*Int).Value() >= r.(*Int).Value()), nil
	case TypeFloat:
		return NewBool(l.(*Float).Value() >= r.(*Float).Value()), nil
	case TypeInt64:
		return NewBool(l.(*Int64).Value() >= r.(*Int64).Value()), nil
	case TypeFloat64:
		return NewBool(l.(*Float64).Value() >= r.(*Float64).Value()), nil
	case TypeString:
		return NewBool(l.(*String).Value() >= r.(*String).Value()), nil
	default:
		return nil, NewTypeError("invalid types for comparison: %s and %s", l.Type(), r.Type())
	}
}

func equals(l, r Value) (Value, error) {
	// Handle null equality
	if l.Type() == TypeNull || r.Type() == TypeNull {
		return NewBool(l.Type() == TypeNull && r.Type() == TypeNull), nil
	}

	// For all other types, use the Value's Equals method
	return NewBool(l.Equals(r)), nil
}

func negate(v Value) (Value, error) {
	switch v.Type() {
	case TypeInt:
		return NewInt(-v.(*Int).Value()), nil
	case TypeFloat:
		return NewFloat(-v.(*Float).Value()), nil
	case TypeInt64:
		return NewInt64(-v.(*Int64).Value()), nil
	case TypeFloat64:
		return NewFloat64(-v.(*Float64).Value()), nil
	default:
		return nil, NewTypeError("cannot negate %s", v.Type())
	}
}
