package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestObjectValue(t *testing.T) {
	instance := &MockObjectInstance{
		name:      "test",
		hierarchy: []string{"TestClass"},
	}
	obj := NewObjectValue(instance, "TestClass")

	assert.Equal(t, "TestClass", obj.Type())
	assert.Equal(t, "<TestClass object>", obj.String())
	assert.Equal(t, YEZ, obj.ToBool()) // Objects are always truthy
	assert.False(t, obj.IsNothing())
	assert.Equal(t, obj, obj.Copy()) // Objects are reference types
}

func TestObjectValueWithNilInstance(t *testing.T) {
	obj := NewObjectValue(nil, "TestClass")

	assert.Equal(t, "TestClass", obj.Type())
	assert.Equal(t, "<TestClass object>", obj.String())
	assert.Equal(t, YEZ, obj.ToBool()) // Even nil objects are truthy by design
	assert.True(t, obj.IsNothing())    // But they are "nothing"
}

func TestObjectValueCasting(t *testing.T) {
	instance := &MockObjectInstance{
		name:      "test",
		hierarchy: []string{"TestClass", "TestParent"},
	}
	obj := NewObjectValue(instance, "TestClass")

	tests := []struct {
		name       string
		targetType string
		expected   Value
		shouldErr  bool
	}{
		{"Same class", "TestClass", obj, false},
		{"Parent class", "TestParent", obj, false},
		{"Empty type", "", obj, false},
		{"To string", "STRIN", StringValue("<TestClass object>"), false},
		{"To bool", "BOOL", YEZ, false},
		{"Invalid type", "INTEGR", nil, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := obj.Cast(test.targetType)
			if test.shouldErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

func TestObjectValueEquality(t *testing.T) {
	instance1 := &MockObjectInstance{
		name:      "test1",
		hierarchy: []string{"TestClass"},
	}
	instance2 := &MockObjectInstance{
		name:      "test2",
		hierarchy: []string{"TestClass"},
	}

	obj1a := NewObjectValue(instance1, "TestClass")
	obj1b := NewObjectValue(instance1, "TestClass") // Same instance
	obj2 := NewObjectValue(instance2, "TestClass")  // Different instance
	nilObj := NewObjectValue(nil, "TestClass")

	tests := []struct {
		name     string
		obj1     ObjectValue
		obj2     Value
		expected BoolValue
	}{
		{"Same instance", obj1a, obj1b, YEZ},
		{"Different instance", obj1a, obj2, NO},
		{"Object vs nil", obj1a, nilObj, NO},
		{"Nil vs nil", nilObj, NewObjectValue(nil, "TestClass"), YEZ},
		{"Object vs nothing", obj1a, NothingValue{}, NO},
		{"Object vs other type", obj1a, IntegerValue(42), NO},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.obj1.EqualTo(test.obj2)
			require.NoError(t, err)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestNewObjectValue(t *testing.T) {
	instance := &MockObjectInstance{name: "test"}
	obj := NewObjectValue(instance, "TestClass")

	assert.Equal(t, instance, obj.Instance)
	assert.Equal(t, "TestClass", obj.ClassName)
}

// MockObjectInstance is a simple mock for testing
type MockObjectInstance struct {
	name      string
	hierarchy []string
}

func (m *MockObjectInstance) GetHierarchy() []string {
	return m.hierarchy
}

func (m *MockObjectInstance) GetName() string {
	return m.name
}
