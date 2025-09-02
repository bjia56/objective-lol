package stdlib

import (
	"testing"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/interpreter"
	"github.com/bjia56/objective-lol/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterIO(t *testing.T) {
	env := environment.NewEnvironment(nil)

	err := RegisterIOInEnv(env)
	require.NoError(t, err)

	// Test that IO classes are registered
	_, err = env.GetClass("READER")
	assert.NoError(t, err, "READER class should be registered")

	_, err = env.GetClass("WRITER")
	assert.NoError(t, err, "WRITER class should be registered")

	_, err = env.GetClass("BUFFERED_READER")
	assert.NoError(t, err, "BUFFERED_READER class should be registered")
}

func TestRegisterIOSelective(t *testing.T) {
	env := environment.NewEnvironment(nil)

	// Test selective import
	err := RegisterIOInEnv(env, "BUFFERED_READER")
	require.NoError(t, err)

	// Should have BUFFERED_READER
	_, err = env.GetClass("BUFFERED_READER")
	assert.NoError(t, err, "BUFFERED_READER class should be registered")

	// Should not have READER or WRITER
	_, err = env.GetClass("READER")
	assert.Error(t, err, "READER class should not be registered")

	_, err = env.GetClass("WRITER")
	assert.Error(t, err, "WRITER class should not be registered")
}

// MockReader implements a simple reader for testing
type MockReader struct {
	data        string
	position    int
	readCount   int
	closeCount  int
	maxReadSize int
}

func (m *MockReader) setupAsReader(t *testing.T, env *environment.Environment) *environment.ObjectInstance {
	// Create a mock reader class that extends READER
	readerClass := &environment.Class{
		Name:        "MockReader",
		ParentClass: "READER",
		PublicFunctions: map[string]*environment.Function{
			"READ": {
				Name:       "READ",
				ReturnType: "STRIN",
				Parameters: []environment.Parameter{{Name: "size", Type: "INTEGR"}},
				NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
					m.readCount++
					if len(args) != 1 {
						return types.NOTHIN, assert.AnError
					}
					
					size := int(args[0].(types.IntegerValue))
					if size <= 0 {
						return types.StringValue(""), nil
					}
					
					if m.maxReadSize > 0 && size > m.maxReadSize {
						size = m.maxReadSize
					}
					
					remaining := len(m.data) - m.position
					if remaining <= 0 {
						return types.StringValue(""), nil // EOF
					}
					
					if size > remaining {
						size = remaining
					}
					
					result := m.data[m.position : m.position+size]
					m.position += size
					return types.StringValue(result), nil
				},
			},
			"CLOSE": {
				Name: "CLOSE",
				NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
					m.closeCount++
					return types.NOTHIN, nil
				},
			},
		},
	}
	
	env.DefineClass(readerClass)
	instanceInterface, err := env.NewObjectInstance("MockReader")
	require.NoError(t, err)
	
	instance, ok := instanceInterface.(*environment.ObjectInstance)
	require.True(t, ok)
	
	return instance
}

func TestBufferedReaderClass(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterIOInEnv(env)

	bufferedClass, err := env.GetClass("BUFFERED_READER")
	require.NoError(t, err)

	// Check that required methods exist
	_, exists := bufferedClass.PublicFunctions["BUFFERED_READER"]
	assert.True(t, exists, "BUFFERED_READER constructor should exist")

	_, exists = bufferedClass.PublicFunctions["READ"]
	assert.True(t, exists, "READ method should exist")

	_, exists = bufferedClass.PublicFunctions["SET_SIZ"]
	assert.True(t, exists, "SET_SIZ method should exist")

	_, exists = bufferedClass.PublicFunctions["CLOSE"]
	assert.True(t, exists, "CLOSE method should exist")

	// Check SIZ variable
	_, exists = bufferedClass.PublicVariables["SIZ"]
	assert.True(t, exists, "SIZ variable should exist")
}

func TestBufferedReaderConstructor(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterIOInEnv(env)
	
	mockReader := &MockReader{data: "test data", maxReadSize: 100}
	readerInstance := mockReader.setupAsReader(t, env)
	
	// Create BUFFERED_READER instance
	instanceInterface, err := env.NewObjectInstance("BUFFERED_READER")
	require.NoError(t, err)
	
	instance, ok := instanceInterface.(*environment.ObjectInstance)
	require.True(t, ok)
	
	// Call constructor
	bufferedClass, err := env.GetClass("BUFFERED_READER")
	require.NoError(t, err)
	
	constructor, exists := bufferedClass.PublicFunctions["BUFFERED_READER"]
	require.True(t, exists)
	
	readerValue := types.NewObjectValue(readerInstance, "MockReader")
	_, err = constructor.NativeImpl(nil, instance, []types.Value{readerValue})
	require.NoError(t, err)
	
	// Verify that BufferedReaderData was initialized
	bufferData, ok := instance.NativeData.(*BufferedReaderData)
	require.True(t, ok, "NativeData should be BufferedReaderData")
	assert.Equal(t, readerInstance, bufferData.Reader)
	assert.Equal(t, "", bufferData.Buffer)
	assert.Equal(t, 0, bufferData.Position)
	assert.Equal(t, defaultBufferSize, bufferData.BufferSize)
	assert.False(t, bufferData.EOF)
	
	// Verify SIZ variable was set
	sizVar, exists := instance.Variables["SIZ"]
	require.True(t, exists, "SIZ variable should exist")
	assert.Equal(t, types.IntegerValue(defaultBufferSize), sizVar.Value)
}

func TestBufferedReaderRead(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterIOInEnv(env)
	
	mockReader := &MockReader{
		data:        "Hello, World! This is a test string.",
		maxReadSize: 20, // Simulate small underlying reads
	}
	readerInstance := mockReader.setupAsReader(t, env)
	
	// Create and initialize BUFFERED_READER
	instanceInterface, err := env.NewObjectInstance("BUFFERED_READER")
	require.NoError(t, err)
	
	instance, ok := instanceInterface.(*environment.ObjectInstance)
	require.True(t, ok)
	
	bufferedClass, err := env.GetClass("BUFFERED_READER")
	require.NoError(t, err)
	
	constructor := bufferedClass.PublicFunctions["BUFFERED_READER"]
	readerValue := types.NewObjectValue(readerInstance, "MockReader")
	_, err = constructor.NativeImpl(nil, instance, []types.Value{readerValue})
	require.NoError(t, err)
	
	// Create function context for method calls  
	interp := interpreter.NewInterpreter(
		map[string]interpreter.StdlibInitializer{
			"IO": RegisterIOInEnv,
		},
		DefaultGlobalInitializers()...,
	)
	ctx := interpreter.NewFunctionContext(interp, env)
	
	readMethod := bufferedClass.PublicFunctions["READ"]
	
	// Test 1: Read 5 characters
	result, err := readMethod.NativeImpl(ctx, instance, []types.Value{types.IntegerValue(5)})
	require.NoError(t, err)
	assert.Equal(t, types.StringValue("Hello"), result)
	assert.Equal(t, 1, mockReader.readCount, "Should have called underlying reader once")
	
	// Test 2: Read 7 more characters (should use buffer)
	result, err = readMethod.NativeImpl(ctx, instance, []types.Value{types.IntegerValue(7)})
	require.NoError(t, err)
	assert.Equal(t, types.StringValue(", World"), result)
	assert.Equal(t, 1, mockReader.readCount, "Should still have called underlying reader only once")
	
	// Test 3: Read more to exhaust buffer
	result, err = readMethod.NativeImpl(ctx, instance, []types.Value{types.IntegerValue(10)})
	require.NoError(t, err)
	assert.Equal(t, types.StringValue("! This is "), result)
	assert.Equal(t, 2, mockReader.readCount, "Should have called underlying reader twice now")
}

func TestBufferedReaderSetSiz(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterIOInEnv(env)
	
	mockReader := &MockReader{data: "test data", maxReadSize: 100}
	readerInstance := mockReader.setupAsReader(t, env)
	
	// Create and initialize BUFFERED_READER
	instanceInterface, err := env.NewObjectInstance("BUFFERED_READER")
	require.NoError(t, err)
	
	instance, ok := instanceInterface.(*environment.ObjectInstance)
	require.True(t, ok)
	
	bufferedClass, err := env.GetClass("BUFFERED_READER")
	require.NoError(t, err)
	
	constructor := bufferedClass.PublicFunctions["BUFFERED_READER"]
	readerValue := types.NewObjectValue(readerInstance, "MockReader")
	_, err = constructor.NativeImpl(nil, instance, []types.Value{readerValue})
	require.NoError(t, err)
	
	// Change buffer size
	setSizMethod := bufferedClass.PublicFunctions["SET_SIZ"]
	_, err = setSizMethod.NativeImpl(nil, instance, []types.Value{types.IntegerValue(512)})
	require.NoError(t, err)
	
	// Verify buffer size changed
	bufferData, ok := instance.NativeData.(*BufferedReaderData)
	require.True(t, ok)
	assert.Equal(t, 512, bufferData.BufferSize)
	
	// Verify SIZ variable updated
	sizVar, exists := instance.Variables["SIZ"]
	require.True(t, exists)
	assert.Equal(t, types.IntegerValue(512), sizVar.Value)
}

func TestBufferedReaderClose(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterIOInEnv(env)
	
	mockReader := &MockReader{data: "test data", maxReadSize: 100}
	readerInstance := mockReader.setupAsReader(t, env)
	
	// Create and initialize BUFFERED_READER
	instanceInterface, err := env.NewObjectInstance("BUFFERED_READER")
	require.NoError(t, err)
	
	instance, ok := instanceInterface.(*environment.ObjectInstance)
	require.True(t, ok)
	
	bufferedClass, err := env.GetClass("BUFFERED_READER")
	require.NoError(t, err)
	
	constructor := bufferedClass.PublicFunctions["BUFFERED_READER"]
	readerValue := types.NewObjectValue(readerInstance, "MockReader")
	_, err = constructor.NativeImpl(nil, instance, []types.Value{readerValue})
	require.NoError(t, err)
	
	// Create function context for close method
	interp := interpreter.NewInterpreter(
		map[string]interpreter.StdlibInitializer{
			"IO": RegisterIOInEnv,
		},
		DefaultGlobalInitializers()...,
	)
	ctx := interpreter.NewFunctionContext(interp, env)
	
	// Close the buffered reader
	closeMethod := bufferedClass.PublicFunctions["CLOSE"]
	_, err = closeMethod.NativeImpl(ctx, instance, []types.Value{})
	require.NoError(t, err)
	
	// Verify underlying reader was closed
	assert.Equal(t, 1, mockReader.closeCount, "Underlying reader should have been closed")
	
	// Verify buffer was cleared and EOF set
	bufferData, ok := instance.NativeData.(*BufferedReaderData)
	require.True(t, ok)
	assert.Equal(t, "", bufferData.Buffer)
	assert.Equal(t, 0, bufferData.Position)
	assert.True(t, bufferData.EOF)
}