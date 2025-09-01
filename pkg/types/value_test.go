package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNothingValue(t *testing.T) {
	nothing := NothingValue{}

	assert.Equal(t, "NOTHIN", nothing.Type())
	assert.Equal(t, "NOTHIN", nothing.String())
	assert.Equal(t, NO, nothing.ToBool())
	assert.True(t, nothing.IsNothing())
	assert.Equal(t, nothing, nothing.Copy())
}

func TestNothingValueCasting(t *testing.T) {
	nothing := NothingValue{}

	tests := []struct {
		targetType string
		expected   Value
	}{
		{"NOTHIN", nothing},
		{"", nothing},
		{"BOOL", NO},
		{"STRIN", StringValue("")},
		{"INTEGR", IntegerValue(0)},
		{"DUBBLE", DoubleValue(0.0)},
	}

	for _, test := range tests {
		t.Run("Cast to "+test.targetType, func(t *testing.T) {
			result, err := nothing.Cast(test.targetType)
			require.NoError(t, err)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestNothingValueEquality(t *testing.T) {
	nothing1 := NothingValue{}
	nothing2 := NothingValue{}

	result, err := nothing1.EqualTo(nothing2)
	require.NoError(t, err)
	assert.Equal(t, YEZ, result)

	// Test inequality with other types
	result, err = nothing1.EqualTo(IntegerValue(0))
	require.NoError(t, err)
	assert.Equal(t, NO, result)
}

func TestBoolValue(t *testing.T) {
	tests := []struct {
		value    BoolValue
		typeName string
		str      string
		toBool   BoolValue
		nothing  bool
	}{
		{YEZ, "BOOL", "YEZ", YEZ, false},
		{NO, "BOOL", "NO", NO, false},
	}

	for _, test := range tests {
		t.Run(test.str, func(t *testing.T) {
			assert.Equal(t, test.typeName, test.value.Type())
			assert.Equal(t, test.str, test.value.String())
			assert.Equal(t, test.toBool, test.value.ToBool())
			assert.Equal(t, test.nothing, test.value.IsNothing())
			assert.Equal(t, test.value, test.value.Copy())
		})
	}
}

func TestBoolValueCasting(t *testing.T) {
	tests := []struct {
		name       string
		value      BoolValue
		targetType string
		expected   Value
		shouldErr  bool
	}{
		{"YEZ to BOOL", YEZ, "BOOL", YEZ, false},
		{"NO to BOOL", NO, "BOOL", NO, false},
		{"YEZ to STRIN", YEZ, "STRIN", StringValue("YEZ"), false},
		{"NO to STRIN", NO, "STRIN", StringValue("NO"), false},
		{"YEZ to INTEGR", YEZ, "INTEGR", IntegerValue(1), false},
		{"NO to INTEGR", NO, "INTEGR", IntegerValue(0), false},
		{"YEZ to DUBBLE", YEZ, "DUBBLE", DoubleValue(1.0), false},
		{"NO to DUBBLE", NO, "DUBBLE", DoubleValue(0.0), false},
		{"Invalid type", YEZ, "INVALID", nil, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.value.Cast(test.targetType)
			if test.shouldErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

func TestIntegerValue(t *testing.T) {
	val := IntegerValue(42)

	assert.Equal(t, "INTEGR", val.Type())
	assert.Equal(t, "42", val.String())
	assert.Equal(t, YEZ, val.ToBool()) // Non-zero is truthy
	assert.False(t, val.IsNothing())

	val0 := IntegerValue(0)
	assert.Equal(t, NO, val0.ToBool()) // Zero is falsy
}

func TestIntegerValueCasting(t *testing.T) {
	tests := []struct {
		name       string
		value      IntegerValue
		targetType string
		expected   Value
		shouldErr  bool
	}{
		{"42 to INTEGR", IntegerValue(42), "INTEGR", IntegerValue(42), false},
		{"42 to STRIN", IntegerValue(42), "STRIN", StringValue("42"), false},
		{"42 to DUBBLE", IntegerValue(42), "DUBBLE", DoubleValue(42.0), false},
		{"42 to BOOL", IntegerValue(42), "BOOL", YEZ, false},
		{"0 to BOOL", IntegerValue(0), "BOOL", NO, false},
		{"Invalid type", IntegerValue(42), "INVALID", nil, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.value.Cast(test.targetType)
			if test.shouldErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

func TestDoubleValue(t *testing.T) {
	val := DoubleValue(3.14)

	assert.Equal(t, "DUBBLE", val.Type())
	assert.Equal(t, "3.14", val.String())
	assert.Equal(t, YEZ, val.ToBool())
	assert.False(t, val.IsNothing())

	val0 := DoubleValue(0.0)
	assert.Equal(t, NO, val0.ToBool())
}

func TestDoubleValueCasting(t *testing.T) {
	tests := []struct {
		name       string
		value      DoubleValue
		targetType string
		expected   Value
		shouldErr  bool
	}{
		{"3.14 to DUBBLE", DoubleValue(3.14), "DUBBLE", DoubleValue(3.14), false},
		{"3.14 to STRIN", DoubleValue(3.14), "STRIN", StringValue("3.14"), false},
		{"3.0 to INTEGR", DoubleValue(3.0), "INTEGR", IntegerValue(3), false},
		{"3.14 to INTEGR", DoubleValue(3.14), "INTEGR", IntegerValue(3), false},
		{"3.14 to BOOL", DoubleValue(3.14), "BOOL", YEZ, false},
		{"0.0 to BOOL", DoubleValue(0.0), "BOOL", NO, false},
		{"Invalid type", DoubleValue(3.14), "INVALID", nil, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.value.Cast(test.targetType)
			if test.shouldErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

func TestStringValue(t *testing.T) {
	val := StringValue("hello")

	assert.Equal(t, "STRIN", val.Type())
	assert.Equal(t, "hello", val.String())
	assert.Equal(t, YEZ, val.ToBool()) // Non-empty is truthy
	assert.False(t, val.IsNothing())

	empty := StringValue("")
	assert.Equal(t, NO, empty.ToBool()) // Empty is falsy
}

func TestStringValueCasting(t *testing.T) {
	tests := []struct {
		name       string
		value      StringValue
		targetType string
		expected   Value
		shouldErr  bool
	}{
		{"hello to STRIN", StringValue("hello"), "STRIN", StringValue("hello"), false},
		{"hello to BOOL", StringValue("hello"), "BOOL", YEZ, false},
		{"empty to BOOL", StringValue(""), "BOOL", NO, false},
		{"42 to INTEGR", StringValue("42"), "INTEGR", IntegerValue(42), false},
		{"3.14 to DUBBLE", StringValue("3.14"), "DUBBLE", DoubleValue(3.14), false},
		{"invalid int", StringValue("abc"), "INTEGR", nil, true},
		{"invalid double", StringValue("abc"), "DUBBLE", nil, true},
		{"Invalid type", StringValue("hello"), "INVALID", nil, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.value.Cast(test.targetType)
			if test.shouldErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

func TestValueEquality(t *testing.T) {
	tests := []struct {
		name     string
		val1     Value
		val2     Value
		expected BoolValue
	}{
		{"Integer equality", IntegerValue(42), IntegerValue(42), YEZ},
		{"Integer inequality", IntegerValue(42), IntegerValue(43), NO},
		{"Double equality", DoubleValue(3.14), DoubleValue(3.14), YEZ},
		{"String equality", StringValue("hello"), StringValue("hello"), YEZ},
		{"Bool equality", YEZ, YEZ, YEZ},
		{"Cross-type comparison", IntegerValue(42), DoubleValue(42.0), YEZ},
		{"Cross-type inequality", IntegerValue(42), StringValue("hello"), NO},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.val1.EqualTo(test.val2)
			require.NoError(t, err)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestValueCopy(t *testing.T) {
	tests := []Value{
		NothingValue{},
		YEZ,
		IntegerValue(42),
		DoubleValue(3.14),
		StringValue("hello"),
	}

	for _, val := range tests {
		t.Run(val.Type(), func(t *testing.T) {
			copy := val.Copy()
			assert.Equal(t, val, copy)

			// Ensure it's a deep copy for strings
			if str, ok := val.(StringValue); ok {
				copyStr := copy.(StringValue)
				assert.Equal(t, str, copyStr)
			}
		})
	}
}
