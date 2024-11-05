package types

import (
	"strconv"
)

// Convert attempts to convert a value to the specified type
func Convert(v Value, to Type) (Value, error) {
	// Same type, just return a clone
	if v.Type() == to {
		return v.Clone(), nil
	}

	// Handle null conversions
	if v.Type() == TypeNull {
		return convertNull(to)
	}

	switch to {
	case TypeInt:
		return convertToInt(v)
	case TypeFloat:
		return convertToFloat(v)
	case TypeInt64:
		return convertToInt64(v)
	case TypeFloat64:
		return convertToFloat64(v)
	case TypeString:
		return convertToString(v)
	case TypeBool:
		return convertToBool(v)
	default:
		return nil, NewTypeError("cannot convert %s to %s", v.Type(), to)
	}
}

func convertNull(to Type) (Value, error) {
	switch to {
	case TypeString:
		return NewString("null"), nil
	case TypeBool:
		return NewBool(false), nil
	case TypeInt:
		return NewInt(0), nil
	case TypeFloat:
		return NewFloat(0), nil
	case TypeInt64:
		return NewInt64(0), nil
	case TypeFloat64:
		return NewFloat64(0), nil
	default:
		return nil, NewTypeError("cannot convert null to %s", to)
	}
}

func convertToInt(v Value) (Value, error) {
	switch val := v.(type) {
	case *Int:
		return val.Clone(), nil
	case *Float:
		return NewInt(int32(val.Value())), nil
	case *Int64:
		if val.Value() > int64(^int32(0)) || val.Value() < int64(-int32(^uint32(0)>>1)-1) {
			return nil, NewTypeError("int64 value %d out of range for int32", val.Value())
		}
		return NewInt(int32(val.Value())), nil
	case *Float64:
		if val.Value() > float64(^int32(0)) || val.Value() < float64(-int32(^uint32(0)>>1)-1) {
			return nil, NewTypeError("float64 value %f out of range for int32", val.Value())
		}
		return NewInt(int32(val.Value())), nil
	case *String:
		i, err := strconv.ParseInt(val.Value(), 10, 32)
		if err != nil {
			return nil, NewTypeError("cannot convert string '%s' to int: %v", val.Value(), err)
		}
		return NewInt(int32(i)), nil
	case *Bool:
		if val.Value() {
			return NewInt(1), nil
		}
		return NewInt(0), nil
	default:
		return nil, NewTypeError("cannot convert %s to int", v.Type())
	}
}

func convertToFloat(v Value) (Value, error) {
	switch val := v.(type) {
	case *Int:
		return NewFloat(float32(val.Value())), nil
	case *Float:
		return val.Clone(), nil
	case *Int64:
		return NewFloat(float32(val.Value())), nil
	case *Float64:
		return NewFloat(float32(val.Value())), nil
	case *String:
		f, err := strconv.ParseFloat(val.Value(), 32)
		if err != nil {
			return nil, NewTypeError("cannot convert string '%s' to float: %v", val.Value(), err)
		}
		return NewFloat(float32(f)), nil
	case *Bool:
		if val.Value() {
			return NewFloat(1), nil
		}
		return NewFloat(0), nil
	default:
		return nil, NewTypeError("cannot convert %s to float", v.Type())
	}
}

func convertToInt64(v Value) (Value, error) {
	switch val := v.(type) {
	case *Int:
		return NewInt64(int64(val.Value())), nil
	case *Float:
		return NewInt64(int64(val.Value())), nil
	case *Int64:
		return val.Clone(), nil
	case *Float64:
		return NewInt64(int64(val.Value())), nil
	case *String:
		i, err := strconv.ParseInt(val.Value(), 10, 64)
		if err != nil {
			return nil, NewTypeError("cannot convert string '%s' to int64: %v", val.Value(), err)
		}
		return NewInt64(i), nil
	case *Bool:
		if val.Value() {
			return NewInt64(1), nil
		}
		return NewInt64(0), nil
	default:
		return nil, NewTypeError("cannot convert %s to int64", v.Type())
	}
}

func convertToFloat64(v Value) (Value, error) {
	switch val := v.(type) {
	case *Int:
		return NewFloat64(float64(val.Value())), nil
	case *Float:
		return NewFloat64(float64(val.Value())), nil
	case *Int64:
		return NewFloat64(float64(val.Value())), nil
	case *Float64:
		return val.Clone(), nil
	case *String:
		f, err := strconv.ParseFloat(val.Value(), 64)
		if err != nil {
			return nil, NewTypeError("cannot convert string '%s' to float64: %v", val.Value(), err)
		}
		return NewFloat64(f), nil
	case *Bool:
		if val.Value() {
			return NewFloat64(1), nil
		}
		return NewFloat64(0), nil
	default:
		return nil, NewTypeError("cannot convert %s to float64", v.Type())
	}
}

func convertToString(v Value) (Value, error) {
	return NewString(v.String()), nil
}

func convertToBool(v Value) (Value, error) {
	switch val := v.(type) {
	case *Bool:
		return val.Clone(), nil
	case *String:
		switch val.Value() {
		case "true", "1":
			return NewBool(true), nil
		case "false", "0", "":
			return NewBool(false), nil
		default:
			return nil, NewTypeError("cannot convert string '%s' to bool", val.Value())
		}
	case *Int:
		return NewBool(val.Value() != 0), nil
	case *Float:
		return NewBool(val.Value() != 0), nil
	case *Int64:
		return NewBool(val.Value() != 0), nil
	case *Float64:
		return NewBool(val.Value() != 0), nil
	default:
		return nil, NewTypeError("cannot convert %s to bool", v.Type())
	}
}
