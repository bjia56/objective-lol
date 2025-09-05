package stdlib

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterSTDIO(t *testing.T) {
	env := environment.NewEnvironment(nil)

	err := RegisterSTDIOInEnv(env)
	require.NoError(t, err)

	// Test that stdio functions are registered
	stdioFunctions := []string{"SAY", "SAYZ", "GIMME"}

	for _, funcName := range stdioFunctions {
		_, err := env.GetFunction(funcName)
		assert.NoError(t, err, "Function %s should be registered", funcName)
	}
}

func TestStdioSAY(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterSTDIOInEnv(env)

	sayFunc, err := env.GetFunction("SAY")
	require.NoError(t, err)

	// Set up buffer to capture output
	var buf bytes.Buffer
	SetOutput(&buf)
	defer ResetToStandardStreams() // Clean up

	tests := []struct {
		name     string
		input    environment.Value
		expected string
	}{
		{
			"Print string",
			environment.StringValue("Hello"),
			"Hello",
		},
		{
			"Print integer",
			environment.IntegerValue(42),
			"42",
		},
		{
			"Print double",
			environment.DoubleValue(3.14),
			"3.14",
		},
		{
			"Print boolean",
			environment.YEZ,
			"YEZ",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := []environment.Value{test.input}
			result, err := sayFunc.NativeImpl(nil, nil, args)
			require.NoError(t, err)
			assert.Equal(t, environment.NOTHIN, result)
		})
	}

	// Verify all outputs are concatenated
	output := buf.String()
	expected := "Hello423.14YEZ"
	assert.Equal(t, expected, output)
}

func TestStdioSAYZ(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterSTDIOInEnv(env)

	sayzFunc, err := env.GetFunction("SAYZ")
	require.NoError(t, err)

	// Set up buffer to capture output
	var buf bytes.Buffer
	SetOutput(&buf)
	defer ResetToStandardStreams() // Clean up

	args := []environment.Value{environment.StringValue("Hello World")}
	result, err := sayzFunc.NativeImpl(nil, nil, args)
	require.NoError(t, err)
	assert.Equal(t, environment.NOTHIN, result)

	// SAYZ should add a newline
	output := buf.String()
	assert.Equal(t, "Hello World\n", output)
}

func TestStdioGIMME(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterSTDIOInEnv(env)

	gimmeFunc, err := env.GetFunction("GIMME")
	require.NoError(t, err)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"Simple input",
			"hello\n",
			"hello",
		},
		{
			"Input with spaces",
			"hello world\n",
			"hello world",
		},
		{
			"Windows line ending",
			"test\r\n",
			"test",
		},
		{
			"Empty input",
			"\n",
			"",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Set up string reader for input
			inputReader := strings.NewReader(test.input)
			SetInput(inputReader)
			defer ResetToStandardStreams() // Clean up

			args := []environment.Value{}
			result, err := gimmeFunc.NativeImpl(nil, nil, args)
			require.NoError(t, err)

			stringResult, ok := result.(environment.StringValue)
			require.True(t, ok, "GIMME should return StringValue")
			assert.Equal(t, test.expected, string(stringResult))
		})
	}
}

func TestStdioSelectiveImport(t *testing.T) {
	env := environment.NewEnvironment(nil)

	// Import only SAY function
	err := RegisterSTDIOInEnv(env, "SAY")
	require.NoError(t, err)

	// SAY should be available
	_, err = env.GetFunction("SAY")
	assert.NoError(t, err)

	// SAYZ should not be available
	_, err = env.GetFunction("SAYZ")
	assert.Error(t, err)

	// GIMME should not be available
	_, err = env.GetFunction("GIMME")
	assert.Error(t, err)
}

func TestStdioInvalidFunction(t *testing.T) {
	env := environment.NewEnvironment(nil)

	// Try to import non-existent function
	err := RegisterSTDIOInEnv(env, "INVALID")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown STDIO function: INVALID")
}

func TestStdioCaseInsensitive(t *testing.T) {
	env := environment.NewEnvironment(nil)

	// Import using lowercase
	err := RegisterSTDIOInEnv(env, "say", "sayz")
	require.NoError(t, err)

	// Functions should be available (stored in uppercase)
	_, err = env.GetFunction("SAY")
	assert.NoError(t, err)

	_, err = env.GetFunction("SAYZ")
	assert.NoError(t, err)
}
