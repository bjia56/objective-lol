package stdlib

import (
	"testing"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterJSON(t *testing.T) {
	env := environment.NewEnvironment(nil)
	err := RegisterJSONInEnv(env)
	require.NoError(t, err)

	for _, name := range []string{"TO_JSON", "FROM_JSON"} {
		_, err := env.GetFunction(name)
		assert.NoError(t, err, "Function %s should be registered", name)
	}
}

func TestRegisterJSONSelective(t *testing.T) {
	env := environment.NewEnvironment(nil)
	err := RegisterJSONInEnv(env, "TO_JSON")
	require.NoError(t, err)

	_, err = env.GetFunction("TO_JSON")
	assert.NoError(t, err)

	_, err = env.GetFunction("FROM_JSON")
	assert.Error(t, err)
}

func TestRegisterJSONUnknownDeclaration(t *testing.T) {
	env := environment.NewEnvironment(nil)
	err := RegisterJSONInEnv(env, "NONEXISTENT")
	assert.Error(t, err)
}

// --- TO_JSON tests ---

func getJSONFn(t *testing.T, name string) *environment.Function {
	t.Helper()
	env := environment.NewEnvironment(nil)
	require.NoError(t, RegisterJSONInEnv(env))
	fn, err := env.GetFunction(name)
	require.NoError(t, err)
	return fn
}

func TestToJSONNothing(t *testing.T) {
	fn := getJSONFn(t, "TO_JSON")
	result, err := fn.NativeImpl(nil, nil, []environment.Value{environment.NOTHIN})
	require.NoError(t, err)
	assert.Equal(t, environment.StringValue("null"), result)
}

func TestToJSONBool(t *testing.T) {
	fn := getJSONFn(t, "TO_JSON")

	result, err := fn.NativeImpl(nil, nil, []environment.Value{environment.YEZ})
	require.NoError(t, err)
	assert.Equal(t, environment.StringValue("true"), result)

	result, err = fn.NativeImpl(nil, nil, []environment.Value{environment.NO})
	require.NoError(t, err)
	assert.Equal(t, environment.StringValue("false"), result)
}

func TestToJSONInteger(t *testing.T) {
	fn := getJSONFn(t, "TO_JSON")
	result, err := fn.NativeImpl(nil, nil, []environment.Value{environment.IntegerValue(42)})
	require.NoError(t, err)
	assert.Equal(t, environment.StringValue("42"), result)
}

func TestToJSONDouble(t *testing.T) {
	fn := getJSONFn(t, "TO_JSON")
	result, err := fn.NativeImpl(nil, nil, []environment.Value{environment.DoubleValue(3.14)})
	require.NoError(t, err)
	assert.Equal(t, environment.StringValue("3.14"), result)
}

func TestToJSONString(t *testing.T) {
	fn := getJSONFn(t, "TO_JSON")
	result, err := fn.NativeImpl(nil, nil, []environment.Value{environment.StringValue("hello")})
	require.NoError(t, err)
	assert.Equal(t, environment.StringValue(`"hello"`), result)
}

func TestToJSONBukkit(t *testing.T) {
	fn := getJSONFn(t, "TO_JSON")

	bukkit := NewBukkitInstance()
	bukkit.NativeData = BukkitSlice{
		environment.IntegerValue(1),
		environment.StringValue("two"),
		environment.BoolValue(false),
	}

	result, err := fn.NativeImpl(nil, nil, []environment.Value{bukkit})
	require.NoError(t, err)
	assert.Equal(t, environment.StringValue(`[1,"two",false]`), result)
}

func TestToJSONBaskit(t *testing.T) {
	fn := getJSONFn(t, "TO_JSON")

	baskit := NewBaskitInstance()
	baskit.NativeData = BaskitMap{
		"name": environment.StringValue("Alice"),
		"age":  environment.IntegerValue(30),
	}

	result, err := fn.NativeImpl(nil, nil, []environment.Value{baskit})
	require.NoError(t, err)

	// JSON object key order is not guaranteed, so parse and compare
	str := string(result.(environment.StringValue))
	assert.Contains(t, str, `"name":"Alice"`)
	assert.Contains(t, str, `"age":30`)
}

func TestToJSONNestedStructure(t *testing.T) {
	fn := getJSONFn(t, "TO_JSON")

	inner := NewBukkitInstance()
	inner.NativeData = BukkitSlice{environment.IntegerValue(1), environment.IntegerValue(2)}

	outer := NewBaskitInstance()
	outer.NativeData = BaskitMap{"items": inner}

	result, err := fn.NativeImpl(nil, nil, []environment.Value{outer})
	require.NoError(t, err)
	assert.Equal(t, environment.StringValue(`{"items":[1,2]}`), result)
}

// --- FROM_JSON tests ---

func TestFromJSONNull(t *testing.T) {
	fn := getJSONFn(t, "FROM_JSON")
	result, err := fn.NativeImpl(nil, nil, []environment.Value{environment.StringValue("null")})
	require.NoError(t, err)
	assert.True(t, result.IsNothing())
}

func TestFromJSONBool(t *testing.T) {
	fn := getJSONFn(t, "FROM_JSON")

	result, err := fn.NativeImpl(nil, nil, []environment.Value{environment.StringValue("true")})
	require.NoError(t, err)
	assert.Equal(t, environment.BoolValue(true), result)

	result, err = fn.NativeImpl(nil, nil, []environment.Value{environment.StringValue("false")})
	require.NoError(t, err)
	assert.Equal(t, environment.BoolValue(false), result)
}

func TestFromJSONNumber(t *testing.T) {
	fn := getJSONFn(t, "FROM_JSON")

	result, err := fn.NativeImpl(nil, nil, []environment.Value{environment.StringValue("3.14")})
	require.NoError(t, err)
	assert.Equal(t, environment.DoubleValue(3.14), result)

	result, err = fn.NativeImpl(nil, nil, []environment.Value{environment.StringValue("42")})
	require.NoError(t, err)
	assert.Equal(t, environment.DoubleValue(42), result)
}

func TestFromJSONString(t *testing.T) {
	fn := getJSONFn(t, "FROM_JSON")
	result, err := fn.NativeImpl(nil, nil, []environment.Value{environment.StringValue(`"hello"`)})
	require.NoError(t, err)
	assert.Equal(t, environment.StringValue("hello"), result)
}

func TestFromJSONArray(t *testing.T) {
	fn := getJSONFn(t, "FROM_JSON")
	result, err := fn.NativeImpl(nil, nil, []environment.Value{environment.StringValue(`[1, "two", false]`)})
	require.NoError(t, err)

	obj, ok := result.(*environment.ObjectInstance)
	require.True(t, ok)
	assert.Equal(t, "BUKKIT", obj.Class.Name)

	slice := obj.NativeData.(BukkitSlice)
	require.Len(t, slice, 3)
	assert.Equal(t, environment.DoubleValue(1), slice[0])
	assert.Equal(t, environment.StringValue("two"), slice[1])
	assert.Equal(t, environment.BoolValue(false), slice[2])
}

func TestFromJSONObject(t *testing.T) {
	fn := getJSONFn(t, "FROM_JSON")
	result, err := fn.NativeImpl(nil, nil, []environment.Value{environment.StringValue(`{"name":"Bob","score":99.5}`)})
	require.NoError(t, err)

	obj, ok := result.(*environment.ObjectInstance)
	require.True(t, ok)
	assert.Equal(t, "BASKIT", obj.Class.Name)

	m := obj.NativeData.(BaskitMap)
	assert.Equal(t, environment.StringValue("Bob"), m["name"])
	assert.Equal(t, environment.DoubleValue(99.5), m["score"])
}

func TestFromJSONInvalidInput(t *testing.T) {
	fn := getJSONFn(t, "FROM_JSON")
	_, err := fn.NativeImpl(nil, nil, []environment.Value{environment.StringValue("not json")})
	assert.Error(t, err)
}

func TestFromJSONNonStringArg(t *testing.T) {
	fn := getJSONFn(t, "FROM_JSON")
	_, err := fn.NativeImpl(nil, nil, []environment.Value{environment.IntegerValue(42)})
	assert.Error(t, err)
}

// --- Round-trip test ---

func TestJSONRoundTrip(t *testing.T) {
	toJSON := getJSONFn(t, "TO_JSON")
	fromJSON := getJSONFn(t, "FROM_JSON")

	// Build: {"key": [1, true, null]}
	inner := NewBukkitInstance()
	inner.NativeData = BukkitSlice{
		environment.DoubleValue(1),
		environment.BoolValue(true),
		environment.NOTHIN,
	}
	outer := NewBaskitInstance()
	outer.NativeData = BaskitMap{"key": inner}

	serialized, err := toJSON.NativeImpl(nil, nil, []environment.Value{outer})
	require.NoError(t, err)

	parsed, err := fromJSON.NativeImpl(nil, nil, []environment.Value{serialized})
	require.NoError(t, err)

	parsedObj, ok := parsed.(*environment.ObjectInstance)
	require.True(t, ok)
	assert.Equal(t, "BASKIT", parsedObj.Class.Name)

	m := parsedObj.NativeData.(BaskitMap)
	arrObj, ok := m["key"].(*environment.ObjectInstance)
	require.True(t, ok)
	assert.Equal(t, "BUKKIT", arrObj.Class.Name)

	slice := arrObj.NativeData.(BukkitSlice)
	require.Len(t, slice, 3)
	assert.Equal(t, environment.DoubleValue(1), slice[0])
	assert.Equal(t, environment.BoolValue(true), slice[1])
	assert.True(t, slice[2].IsNothing())
}
