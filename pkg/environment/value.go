package environment

import (
	"fmt"
	"strconv"
	"strings"
)

// Value represents the interface for all LOL values
type Value interface {
	// Type returns the type name of this value
	Type() string

	// String returns the string representation of this value
	String() string

	// Cast attempts to cast this value to the specified type
	Cast(targetType string) (Value, error)

	// EqualTo compares this value with another for equality
	EqualTo(other Value) (BoolValue, error)

	// Copy creates a deep copy of this value
	Copy() Value

	// ToBool returns the boolean representation of this value
	ToBool() BoolValue

	// IsNothing returns true if this value is NOTHIN
	IsNothing() bool
}

// NothingValue represents the NOTHIN value
type NothingValue struct{}

func (n NothingValue) Type() string      { return "NOTHIN" }
func (n NothingValue) String() string    { return "NOTHIN" }
func (n NothingValue) Copy() Value       { return n }
func (n NothingValue) ToBool() BoolValue { return NO }
func (n NothingValue) IsNothing() bool   { return true }

func (n NothingValue) Cast(targetType string) (Value, error) {
	switch targetType {
	case "NOTHIN", "":
		return n, nil
	case "BOOL":
		return NO, nil
	case "STRIN":
		return StringValue(""), nil
	case "INTEGR":
		return IntegerValue(0), nil
	case "DUBBLE":
		return DoubleValue(0.0), nil
	default:
		return nil, fmt.Errorf("cannot cast NOTHIN to %s", targetType)
	}
}

func (n NothingValue) EqualTo(other Value) (BoolValue, error) {
	return BoolValue(other.IsNothing()), nil
}

// BoolValue represents boolean values (YEZ/NO)
type BoolValue bool

const (
	YEZ BoolValue = true
	NO  BoolValue = false
)

func (b BoolValue) Type() string      { return "BOOL" }
func (b BoolValue) Copy() Value       { return b }
func (b BoolValue) ToBool() BoolValue { return b }
func (b BoolValue) IsNothing() bool   { return false }

func (b BoolValue) String() string {
	if b {
		return "YEZ"
	}
	return "NO"
}

func (b BoolValue) Cast(targetType string) (Value, error) {
	switch targetType {
	case "BOOL", "":
		return b, nil
	case "INTEGR":
		if b {
			return IntegerValue(1), nil
		}
		return IntegerValue(0), nil
	case "DUBBLE":
		if b {
			return DoubleValue(1.0), nil
		}
		return DoubleValue(0.0), nil
	case "STRIN":
		return StringValue(b.String()), nil
	default:
		return nil, fmt.Errorf("cannot cast BOOL to %s", targetType)
	}
}

func (b BoolValue) EqualTo(other Value) (BoolValue, error) {
	if other.IsNothing() {
		return NO, nil
	}
	otherBool := other.ToBool()
	return BoolValue(b == otherBool), nil
}

// IntegerValue represents integer values
type IntegerValue int64

func (i IntegerValue) Type() string    { return "INTEGR" }
func (i IntegerValue) String() string  { return strconv.FormatInt(int64(i), 10) }
func (i IntegerValue) Copy() Value     { return i }
func (i IntegerValue) IsNothing() bool { return false }

func (i IntegerValue) ToBool() BoolValue {
	return BoolValue(i != 0)
}

func (i IntegerValue) Cast(targetType string) (Value, error) {
	switch targetType {
	case "INTEGR", "":
		return i, nil
	case "DUBBLE":
		return DoubleValue(float64(i)), nil
	case "BOOL":
		return i.ToBool(), nil
	case "STRIN":
		return StringValue(i.String()), nil
	default:
		return nil, fmt.Errorf("cannot cast INTEGR to %s", targetType)
	}
}

func (i IntegerValue) EqualTo(other Value) (BoolValue, error) {
	if other.IsNothing() {
		return NO, nil
	}

	switch v := other.(type) {
	case IntegerValue:
		return BoolValue(i == v), nil
	case DoubleValue:
		return BoolValue(float64(i) == float64(v)), nil
	default:
		// Try casting other to INTEGR
		casted, err := other.Cast("INTEGR")
		if err != nil {
			return NO, nil
		}
		return BoolValue(i == casted.(IntegerValue)), nil
	}
}

// Arithmetic operations for IntegerValue
func (i IntegerValue) Add(other NumberValue) NumberValue {
	switch v := other.(type) {
	case IntegerValue:
		return IntegerValue(i + v)
	case DoubleValue:
		return DoubleValue(float64(i) + float64(v))
	}
	return i
}

func (i IntegerValue) Subtract(other NumberValue) NumberValue {
	switch v := other.(type) {
	case IntegerValue:
		return IntegerValue(i - v)
	case DoubleValue:
		return DoubleValue(float64(i) - float64(v))
	}
	return i
}

func (i IntegerValue) Multiply(other NumberValue) NumberValue {
	switch v := other.(type) {
	case IntegerValue:
		return IntegerValue(i * v)
	case DoubleValue:
		return DoubleValue(float64(i) * float64(v))
	}
	return i
}

func (i IntegerValue) Divide(other NumberValue) NumberValue {
	switch v := other.(type) {
	case IntegerValue:
		if v == 0 {
			return DoubleValue(0) // Handle division by zero
		}
		return DoubleValue(float64(i) / float64(v))
	case DoubleValue:
		if v == 0 {
			return DoubleValue(0) // Handle division by zero
		}
		return DoubleValue(float64(i) / float64(v))
	}
	return i
}

func (i IntegerValue) GreaterThan(other NumberValue) BoolValue {
	switch v := other.(type) {
	case IntegerValue:
		return BoolValue(i > v)
	case DoubleValue:
		return BoolValue(float64(i) > float64(v))
	}
	return NO
}

func (i IntegerValue) LessThan(other NumberValue) BoolValue {
	switch v := other.(type) {
	case IntegerValue:
		return BoolValue(i < v)
	case DoubleValue:
		return BoolValue(float64(i) < float64(v))
	}
	return NO
}

// DoubleValue represents floating-point values
type DoubleValue float64

func (d DoubleValue) Type() string    { return "DUBBLE" }
func (d DoubleValue) String() string  { return strconv.FormatFloat(float64(d), 'f', -1, 64) }
func (d DoubleValue) Copy() Value     { return d }
func (d DoubleValue) IsNothing() bool { return false }

func (d DoubleValue) ToBool() BoolValue {
	return BoolValue(d != 0.0)
}

func (d DoubleValue) Cast(targetType string) (Value, error) {
	switch targetType {
	case "DUBBLE", "":
		return d, nil
	case "INTEGR":
		return IntegerValue(int64(d)), nil
	case "BOOL":
		return d.ToBool(), nil
	case "STRIN":
		return StringValue(d.String()), nil
	default:
		return nil, fmt.Errorf("cannot cast DUBBLE to %s", targetType)
	}
}

func (d DoubleValue) EqualTo(other Value) (BoolValue, error) {
	if other.IsNothing() {
		return NO, nil
	}

	switch v := other.(type) {
	case DoubleValue:
		return BoolValue(d == v), nil
	case IntegerValue:
		return BoolValue(float64(d) == float64(v)), nil
	default:
		// Try casting other to DUBBLE
		casted, err := other.Cast("DUBBLE")
		if err != nil {
			return NO, nil
		}
		return BoolValue(d == casted.(DoubleValue)), nil
	}
}

// Arithmetic operations for DoubleValue
func (d DoubleValue) Add(other NumberValue) NumberValue {
	switch v := other.(type) {
	case IntegerValue:
		return DoubleValue(float64(d) + float64(v))
	case DoubleValue:
		return DoubleValue(d + v)
	}
	return d
}

func (d DoubleValue) Subtract(other NumberValue) NumberValue {
	switch v := other.(type) {
	case IntegerValue:
		return DoubleValue(float64(d) - float64(v))
	case DoubleValue:
		return DoubleValue(d - v)
	}
	return d
}

func (d DoubleValue) Multiply(other NumberValue) NumberValue {
	switch v := other.(type) {
	case IntegerValue:
		return DoubleValue(float64(d) * float64(v))
	case DoubleValue:
		return DoubleValue(d * v)
	}
	return d
}

func (d DoubleValue) Divide(other NumberValue) NumberValue {
	switch v := other.(type) {
	case IntegerValue:
		if v == 0 {
			return DoubleValue(0) // Handle division by zero
		}
		return DoubleValue(float64(d) / float64(v))
	case DoubleValue:
		if v == 0 {
			return DoubleValue(0) // Handle division by zero
		}
		return DoubleValue(d / v)
	}
	return d
}

func (d DoubleValue) GreaterThan(other NumberValue) BoolValue {
	switch v := other.(type) {
	case IntegerValue:
		return BoolValue(float64(d) > float64(v))
	case DoubleValue:
		return BoolValue(d > v)
	}
	return NO
}

func (d DoubleValue) LessThan(other NumberValue) BoolValue {
	switch v := other.(type) {
	case IntegerValue:
		return BoolValue(float64(d) < float64(v))
	case DoubleValue:
		return BoolValue(d < v)
	}
	return NO
}

// StringValue represents string values
type StringValue string

func (s StringValue) Type() string    { return "STRIN" }
func (s StringValue) String() string  { return string(s) }
func (s StringValue) Copy() Value     { return s }
func (s StringValue) IsNothing() bool { return false }

func (s StringValue) ToBool() BoolValue {
	return BoolValue(len(s) > 0)
}

func (s StringValue) Cast(targetType string) (Value, error) {
	str := string(s)

	switch targetType {
	case "STRIN", "":
		return s, nil
	case "BOOL":
		upper := strings.ToUpper(str)
		switch upper {
		case "YEZ":
			return YEZ, nil
		case "NO":
			return NO, nil
		}
		return s.ToBool(), nil
	case "INTEGR":
		// Try parsing as integer
		if val, err := strconv.ParseInt(str, 10, 64); err == nil {
			return IntegerValue(val), nil
		}
		// Try parsing as hex
		if strings.HasPrefix(strings.ToUpper(str), "0X") {
			if val, err := strconv.ParseInt(str[2:], 16, 64); err == nil {
				return IntegerValue(val), nil
			}
		}
		return nil, fmt.Errorf("cannot cast string '%s' to INTEGR", str)
	case "DUBBLE":
		if val, err := strconv.ParseFloat(str, 64); err == nil {
			return DoubleValue(val), nil
		}
		return nil, fmt.Errorf("cannot cast string '%s' to DUBBLE", str)
	default:
		return nil, fmt.Errorf("cannot cast STRIN to %s", targetType)
	}
}

func (s StringValue) EqualTo(other Value) (BoolValue, error) {
	if other.IsNothing() {
		return NO, nil
	}

	switch v := other.(type) {
	case StringValue:
		return BoolValue(s == v), nil
	default:
		// Try casting other to STRIN
		casted, err := other.Cast("STRIN")
		if err != nil {
			return NO, nil
		}
		return BoolValue(s == casted.(StringValue)), nil
	}
}

// NumberValue interface for arithmetic operations
type NumberValue interface {
	Value
	Add(other NumberValue) NumberValue
	Subtract(other NumberValue) NumberValue
	Multiply(other NumberValue) NumberValue
	Divide(other NumberValue) NumberValue
	GreaterThan(other NumberValue) BoolValue
	LessThan(other NumberValue) BoolValue
}

// Global NOTHIN singleton
var NOTHIN = NothingValue{}

// ValueOf creates a LOL value from a Go value
func ValueOf(v interface{}) Value {
	switch val := v.(type) {
	case Value:
		return val
	case bool:
		return BoolValue(val)
	case int:
		return IntegerValue(int64(val))
	case int32:
		return IntegerValue(int64(val))
	case int64:
		return IntegerValue(val)
	case float32:
		return DoubleValue(float64(val))
	case float64:
		return DoubleValue(val)
	case string:
		return StringValue(val)
	case nil:
		return NOTHIN
	default:
		return StringValue(fmt.Sprintf("%v", val))
	}
}
