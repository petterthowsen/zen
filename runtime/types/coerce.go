package types

// CoerceForOperation attempts to coerce two values to compatible types for a binary operation
func CoerceForOperation(left, right Value, op string) (Value, Value, error) {
	// String concatenation
	if op == "+" && (left.Type() == TypeString || right.Type() == TypeString) {
		if left.Type() != TypeString || right.Type() != TypeString {
			return nil, nil, NewTypeError("cannot concatenate %s and %s", left.Type(), right.Type())
		}
		return left, right, nil
	}

	// Logical operations require boolean operands
	if op == "and" || op == "or" {
		if left.Type() != TypeBool || right.Type() != TypeBool {
			return nil, nil, NewTypeError("logical operators require boolean operands, got %s and %s", left.Type(), right.Type())
		}
		return left, right, nil
	}

	// For numeric operations, coerce to the highest precision type
	if IsNumeric(left.Type()) && IsNumeric(right.Type()) {
		return coerceNumeric(left, right)
	}

	// For comparison operations between different types
	if op == "==" || op == "!=" {
		// Null can be compared with any type
		if left.Type() == TypeNull || right.Type() == TypeNull {
			return left, right, nil
		}
		// Otherwise, try to coerce to the same type
		return coerceToSameType(left, right)
	}

	// For other comparison operations
	if op == "<" || op == "<=" || op == ">" || op == ">=" {
		// String comparisons
		if left.Type() == TypeString && right.Type() == TypeString {
			return left, right, nil
		}
		// Numeric comparisons
		if IsNumeric(left.Type()) && IsNumeric(right.Type()) {
			return coerceNumeric(left, right)
		}
		return nil, nil, NewTypeError("cannot compare %s and %s", left.Type(), right.Type())
	}

	return nil, nil, NewTypeError("cannot coerce %s and %s for operation %s", left.Type(), right.Type(), op)
}

// coerceNumeric coerces two numeric values to the highest precision type
func coerceNumeric(left, right Value) (Value, Value, error) {
	// Determine target type based on highest precision
	targetType := highestNumericType(left.Type(), right.Type())

	// Convert both values to target type
	l, err := Convert(left, targetType)
	if err != nil {
		return nil, nil, err
	}

	r, err := Convert(right, targetType)
	if err != nil {
		return nil, nil, err
	}

	return l, r, nil
}

// coerceToSameType attempts to coerce two values to the same type
func coerceToSameType(left, right Value) (Value, Value, error) {
	// If types are the same, no coercion needed
	if left.Type() == right.Type() {
		return left, right, nil
	}

	// Try to coerce to the "higher" type
	targetType := preferredType(left.Type(), right.Type())
	if targetType == TypeNull {
		return nil, nil, NewTypeError("cannot coerce %s and %s to same type", left.Type(), right.Type())
	}

	l, err := Convert(left, targetType)
	if err != nil {
		return nil, nil, err
	}

	r, err := Convert(right, targetType)
	if err != nil {
		return nil, nil, err
	}

	return l, r, nil
}

// highestNumericType returns the highest precision numeric type between two types
func highestNumericType(a, b Type) Type {
	if !IsNumeric(a) || !IsNumeric(b) {
		return TypeNull // Invalid combination
	}

	// Order of precedence: float64 > int64 > float > int
	switch {
	case a == TypeFloat64 || b == TypeFloat64:
		return TypeFloat64
	case a == TypeInt64 || b == TypeInt64:
		return TypeInt64
	case a == TypeFloat || b == TypeFloat:
		return TypeFloat
	default:
		return TypeInt
	}
}

// preferredType returns the preferred type for coercion between two types
func preferredType(a, b Type) Type {
	// If either type is null, can't determine preferred type
	if a == TypeNull || b == TypeNull {
		return TypeNull
	}

	// If both types are numeric, use highest precision
	if IsNumeric(a) && IsNumeric(b) {
		return highestNumericType(a, b)
	}

	// If either type is string and the operation allows string coercion,
	// prefer string
	if a == TypeString || b == TypeString {
		return TypeString
	}

	// If both types are boolean, use boolean
	if a == TypeBool && b == TypeBool {
		return TypeBool
	}

	// Can't determine preferred type
	return TypeNull
}
