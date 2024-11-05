package types

import (
	"strconv"
)

// Int represents a 32-bit integer value
type Int struct {
	value int32
}

func NewInt(value int32) *Int {
	return &Int{value: value}
}

func (i *Int) Type() Type     { return TypeInt }
func (i *Int) String() string { return strconv.FormatInt(int64(i.value), 10) }
func (i *Int) IsTruthy() bool { return i.value != 0 }
func (i *Int) Clone() Value   { return NewInt(i.value) }
func (i *Int) Value() int32   { return i.value }
func (i *Int) Equals(other Value) bool {
	if other.Type() != TypeInt {
		return false
	}
	return i.value == other.(*Int).value
}

// Float represents a 32-bit floating-point value
type Float struct {
	value float32
}

func NewFloat(value float32) *Float {
	return &Float{value: value}
}

func (f *Float) Type() Type     { return TypeFloat }
func (f *Float) String() string { return strconv.FormatFloat(float64(f.value), 'g', -1, 32) }
func (f *Float) IsTruthy() bool { return f.value != 0 }
func (f *Float) Clone() Value   { return NewFloat(f.value) }
func (f *Float) Value() float32 { return f.value }
func (f *Float) Equals(other Value) bool {
	if other.Type() != TypeFloat {
		return false
	}
	return f.value == other.(*Float).value
}

// Int64 represents a 64-bit integer value
type Int64 struct {
	value int64
}

func NewInt64(value int64) *Int64 {
	return &Int64{value: value}
}

func (i *Int64) Type() Type     { return TypeInt64 }
func (i *Int64) String() string { return strconv.FormatInt(i.value, 10) }
func (i *Int64) IsTruthy() bool { return i.value != 0 }
func (i *Int64) Clone() Value   { return NewInt64(i.value) }
func (i *Int64) Value() int64   { return i.value }
func (i *Int64) Equals(other Value) bool {
	if other.Type() != TypeInt64 {
		return false
	}
	return i.value == other.(*Int64).value
}

// Float64 represents a 64-bit floating-point value
type Float64 struct {
	value float64
}

func NewFloat64(value float64) *Float64 {
	return &Float64{value: value}
}

func (f *Float64) Type() Type     { return TypeFloat64 }
func (f *Float64) String() string { return strconv.FormatFloat(f.value, 'g', -1, 64) }
func (f *Float64) IsTruthy() bool { return f.value != 0 }
func (f *Float64) Clone() Value   { return NewFloat64(f.value) }
func (f *Float64) Value() float64 { return f.value }
func (f *Float64) Equals(other Value) bool {
	if other.Type() != TypeFloat64 {
		return false
	}
	return f.value == other.(*Float64).value
}

// String represents a string value
type String struct {
	value string
}

func NewString(value string) *String {
	return &String{value: value}
}

func (s *String) Type() Type     { return TypeString }
func (s *String) String() string { return s.value }
func (s *String) IsTruthy() bool { return len(s.value) > 0 }
func (s *String) Clone() Value   { return NewString(s.value) }
func (s *String) Value() string  { return s.value }
func (s *String) Equals(other Value) bool {
	if other.Type() != TypeString {
		return false
	}
	return s.value == other.(*String).value
}

// Bool represents a boolean value
type Bool struct {
	value bool
}

func NewBool(value bool) *Bool {
	return &Bool{value: value}
}

func (b *Bool) Type() Type     { return TypeBool }
func (b *Bool) String() string { return strconv.FormatBool(b.value) }
func (b *Bool) IsTruthy() bool { return b.value }
func (b *Bool) Clone() Value   { return NewBool(b.value) }
func (b *Bool) Value() bool    { return b.value }
func (b *Bool) Equals(other Value) bool {
	if other.Type() != TypeBool {
		return false
	}
	return b.value == other.(*Bool).value
}

// Null represents a null value
type Null struct{}

var nullValue = &Null{}

func NewNull() *Null {
	return nullValue
}

func (n *Null) Type() Type              { return TypeNull }
func (n *Null) String() string          { return "null" }
func (n *Null) IsTruthy() bool          { return false }
func (n *Null) Clone() Value            { return n }
func (n *Null) Equals(other Value) bool { return other.Type() == TypeNull }

// FromGoValue converts a Go value to a Zen value
func FromGoValue(v interface{}) (Value, error) {
	switch val := v.(type) {
	case nil:
		return NewNull(), nil
	case bool:
		return NewBool(val), nil
	case int:
		return NewInt(int32(val)), nil // Always use int32 for plain ints
	case int32:
		return NewInt(val), nil
	case int64:
		return NewInt64(val), nil
	case float32:
		return NewFloat(val), nil
	case float64:
		return NewFloat64(val), nil
	case string:
		return NewString(val), nil
	default:
		return nil, NewTypeError("cannot convert Go type %T to Zen value", v)
	}
}

// ToGoValue converts a Zen value to a Go value
func ToGoValue(v Value) interface{} {
	switch val := v.(type) {
	case *Int:
		return val.Value()
	case *Float:
		return val.Value()
	case *Int64:
		return val.Value()
	case *Float64:
		return val.Value()
	case *String:
		return val.Value()
	case *Bool:
		return val.Value()
	case *Null:
		return nil
	default:
		return nil
	}
}
