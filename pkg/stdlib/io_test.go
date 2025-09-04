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

	_, err = env.GetClass("BUFFERED_WRITER")
	assert.NoError(t, err, "BUFFERED_WRITER class should be registered")
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
		ParentClasses: []string{"READER"},
		MRO:           []string{},
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

// MockWriter implements a simple writer for testing
type MockWriter struct {
	data         string
	writeCount   int
	closeCount   int
	maxWriteSize int
}

func (m *MockWriter) setupAsWriter(t *testing.T, env *environment.Environment) *environment.ObjectInstance {
	// Create a mock writer class that extends WRITER
	writerClass := &environment.Class{
		Name:        "MockWriter",
		ParentClasses: []string{"WRITER"},
		MRO:           []string{},
		PublicFunctions: map[string]*environment.Function{
			"WRITE": {
				Name:       "WRITE",
				ReturnType: "INTEGR",
				Parameters: []environment.Parameter{{Name: "data", Type: "STRIN"}},
				NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
					m.writeCount++
					if len(args) != 1 {
						return types.NOTHIN, assert.AnError
					}

					data := string(args[0].(types.StringValue))
					actualWritten := len(data)

					if m.maxWriteSize > 0 && actualWritten > m.maxWriteSize {
						actualWritten = m.maxWriteSize
						data = data[:actualWritten]
					}

					m.data += data
					return types.IntegerValue(actualWritten), nil
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

	env.DefineClass(writerClass)
	instanceInterface, err := env.NewObjectInstance("MockWriter")
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

func TestBufferedWriterClass(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterIOInEnv(env)

	bufferedClass, err := env.GetClass("BUFFERED_WRITER")
	require.NoError(t, err)

	// Check that required methods exist
	_, exists := bufferedClass.PublicFunctions["BUFFERED_WRITER"]
	assert.True(t, exists, "BUFFERED_WRITER constructor should exist")

	_, exists = bufferedClass.PublicFunctions["WRITE"]
	assert.True(t, exists, "WRITE method should exist")

	_, exists = bufferedClass.PublicFunctions["FLUSH"]
	assert.True(t, exists, "FLUSH method should exist")

	_, exists = bufferedClass.PublicFunctions["SET_SIZ"]
	assert.True(t, exists, "SET_SIZ method should exist")

	_, exists = bufferedClass.PublicFunctions["CLOSE"]
	assert.True(t, exists, "CLOSE method should exist")

	// Check SIZ variable
	_, exists = bufferedClass.PublicVariables["SIZ"]
	assert.True(t, exists, "SIZ variable should exist")
}

func TestBufferedWriterConstructor(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterIOInEnv(env)

	mockWriter := &MockWriter{}
	writerInstance := mockWriter.setupAsWriter(t, env)

	// Create BUFFERED_WRITER instance
	instanceInterface, err := env.NewObjectInstance("BUFFERED_WRITER")
	require.NoError(t, err)

	instance, ok := instanceInterface.(*environment.ObjectInstance)
	require.True(t, ok)

	// Call constructor
	bufferedClass, err := env.GetClass("BUFFERED_WRITER")
	require.NoError(t, err)

	constructor, exists := bufferedClass.PublicFunctions["BUFFERED_WRITER"]
	require.True(t, exists)

	writerValue := types.NewObjectValue(writerInstance, "MockWriter")
	_, err = constructor.NativeImpl(nil, instance, []types.Value{writerValue})
	require.NoError(t, err)

	// Verify that BufferedWriterData was initialized
	bufferData, ok := instance.NativeData.(*BufferedWriterData)
	require.True(t, ok, "NativeData should be BufferedWriterData")
	assert.Equal(t, writerInstance, bufferData.Writer)
	assert.Equal(t, "", bufferData.Buffer)
	assert.Equal(t, defaultBufferSize, bufferData.BufferSize)

	// Verify SIZ variable was set
	sizVar, exists := instance.Variables["SIZ"]
	require.True(t, exists, "SIZ variable should exist")
	assert.Equal(t, types.IntegerValue(defaultBufferSize), sizVar.Value)
}

func TestBufferedWriterWrite(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterIOInEnv(env)

	mockWriter := &MockWriter{}
	writerInstance := mockWriter.setupAsWriter(t, env)

	// Create and initialize BUFFERED_WRITER
	instanceInterface, err := env.NewObjectInstance("BUFFERED_WRITER")
	require.NoError(t, err)

	instance, ok := instanceInterface.(*environment.ObjectInstance)
	require.True(t, ok)

	bufferedClass, err := env.GetClass("BUFFERED_WRITER")
	require.NoError(t, err)

	constructor := bufferedClass.PublicFunctions["BUFFERED_WRITER"]
	writerValue := types.NewObjectValue(writerInstance, "MockWriter")
	_, err = constructor.NativeImpl(nil, instance, []types.Value{writerValue})
	require.NoError(t, err)

	// Create function context for method calls
	interp := interpreter.NewInterpreter(
		map[string]interpreter.StdlibInitializer{
			"IO": RegisterIOInEnv,
		},
		DefaultGlobalInitializers()...,
	)
	ctx := interpreter.NewFunctionContext(interp, env)

	writeMethod := bufferedClass.PublicFunctions["WRITE"]

	// Test 1: Write small data that should be buffered
	result, err := writeMethod.NativeImpl(ctx, instance, []types.Value{types.StringValue("Hello")})
	require.NoError(t, err)
	assert.Equal(t, types.IntegerValue(5), result)
	assert.Equal(t, 0, mockWriter.writeCount, "Should not have called underlying writer yet")
	assert.Equal(t, "", mockWriter.data, "Mock writer should not have received data yet")

	// Test 2: Write more small data (should still be buffered)
	result, err = writeMethod.NativeImpl(ctx, instance, []types.Value{types.StringValue(", World!")})
	require.NoError(t, err)
	assert.Equal(t, types.IntegerValue(8), result)
	assert.Equal(t, 0, mockWriter.writeCount, "Should still not have called underlying writer")

	// Check buffer contents
	bufferData, ok := instance.NativeData.(*BufferedWriterData)
	require.True(t, ok)
	assert.Equal(t, "Hello, World!", bufferData.Buffer)

	// Test 3: Write large data that exceeds buffer size
	largeData := make([]byte, 2000) // Larger than default buffer size (1024)
	for i := range largeData {
		largeData[i] = 'A'
	}
	result, err = writeMethod.NativeImpl(ctx, instance, []types.Value{types.StringValue(string(largeData))})
	require.NoError(t, err)
	assert.Equal(t, types.IntegerValue(2000), result)
	assert.Equal(t, 2, mockWriter.writeCount, "Should have flushed buffer and written large data")
	expectedData := "Hello, World!" + string(largeData)
	assert.Equal(t, expectedData, mockWriter.data)
	assert.Equal(t, "", bufferData.Buffer, "Buffer should be empty after large write")
}

func TestBufferedWriterFlush(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterIOInEnv(env)

	mockWriter := &MockWriter{}
	writerInstance := mockWriter.setupAsWriter(t, env)

	// Create and initialize BUFFERED_WRITER
	instanceInterface, err := env.NewObjectInstance("BUFFERED_WRITER")
	require.NoError(t, err)

	instance, ok := instanceInterface.(*environment.ObjectInstance)
	require.True(t, ok)

	bufferedClass, err := env.GetClass("BUFFERED_WRITER")
	require.NoError(t, err)

	constructor := bufferedClass.PublicFunctions["BUFFERED_WRITER"]
	writerValue := types.NewObjectValue(writerInstance, "MockWriter")
	_, err = constructor.NativeImpl(nil, instance, []types.Value{writerValue})
	require.NoError(t, err)

	// Create function context
	interp := interpreter.NewInterpreter(
		map[string]interpreter.StdlibInitializer{
			"IO": RegisterIOInEnv,
		},
		DefaultGlobalInitializers()...,
	)
	ctx := interpreter.NewFunctionContext(interp, env)

	writeMethod := bufferedClass.PublicFunctions["WRITE"]
	flushMethod := bufferedClass.PublicFunctions["FLUSH"]

	// Write some data
	_, err = writeMethod.NativeImpl(ctx, instance, []types.Value{types.StringValue("Buffered data")})
	require.NoError(t, err)
	assert.Equal(t, 0, mockWriter.writeCount, "Should not have written yet")

	// Flush the buffer
	_, err = flushMethod.NativeImpl(ctx, instance, []types.Value{})
	require.NoError(t, err)
	assert.Equal(t, 1, mockWriter.writeCount, "Should have written after flush")
	assert.Equal(t, "Buffered data", mockWriter.data)

	// Check buffer is empty
	bufferData, ok := instance.NativeData.(*BufferedWriterData)
	require.True(t, ok)
	assert.Equal(t, "", bufferData.Buffer)
}

func TestBufferedWriterSetSiz(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterIOInEnv(env)

	mockWriter := &MockWriter{}
	writerInstance := mockWriter.setupAsWriter(t, env)

	// Create and initialize BUFFERED_WRITER
	instanceInterface, err := env.NewObjectInstance("BUFFERED_WRITER")
	require.NoError(t, err)

	instance, ok := instanceInterface.(*environment.ObjectInstance)
	require.True(t, ok)

	bufferedClass, err := env.GetClass("BUFFERED_WRITER")
	require.NoError(t, err)

	constructor := bufferedClass.PublicFunctions["BUFFERED_WRITER"]
	writerValue := types.NewObjectValue(writerInstance, "MockWriter")
	_, err = constructor.NativeImpl(nil, instance, []types.Value{writerValue})
	require.NoError(t, err)

	// Create function context
	interp := interpreter.NewInterpreter(
		map[string]interpreter.StdlibInitializer{
			"IO": RegisterIOInEnv,
		},
		DefaultGlobalInitializers()...,
	)
	ctx := interpreter.NewFunctionContext(interp, env)

	writeMethod := bufferedClass.PublicFunctions["WRITE"]
	setSizMethod := bufferedClass.PublicFunctions["SET_SIZ"]

	// Write some data to buffer
	_, err = writeMethod.NativeImpl(ctx, instance, []types.Value{types.StringValue("Test data")})
	require.NoError(t, err)

	// Change buffer size (should flush existing buffer)
	_, err = setSizMethod.NativeImpl(ctx, instance, []types.Value{types.IntegerValue(512)})
	require.NoError(t, err)

	// Verify buffer size changed
	bufferData, ok := instance.NativeData.(*BufferedWriterData)
	require.True(t, ok)
	assert.Equal(t, 512, bufferData.BufferSize)
	assert.Equal(t, "", bufferData.Buffer, "Buffer should have been flushed")
	assert.Equal(t, 1, mockWriter.writeCount, "Should have flushed existing buffer")
	assert.Equal(t, "Test data", mockWriter.data)

	// Verify SIZ variable updated
	sizVar, exists := instance.Variables["SIZ"]
	require.True(t, exists)
	assert.Equal(t, types.IntegerValue(512), sizVar.Value)
}

func TestBufferedWriterClose(t *testing.T) {
	env := environment.NewEnvironment(nil)
	RegisterIOInEnv(env)

	mockWriter := &MockWriter{}
	writerInstance := mockWriter.setupAsWriter(t, env)

	// Create and initialize BUFFERED_WRITER
	instanceInterface, err := env.NewObjectInstance("BUFFERED_WRITER")
	require.NoError(t, err)

	instance, ok := instanceInterface.(*environment.ObjectInstance)
	require.True(t, ok)

	bufferedClass, err := env.GetClass("BUFFERED_WRITER")
	require.NoError(t, err)

	constructor := bufferedClass.PublicFunctions["BUFFERED_WRITER"]
	writerValue := types.NewObjectValue(writerInstance, "MockWriter")
	_, err = constructor.NativeImpl(nil, instance, []types.Value{writerValue})
	require.NoError(t, err)

	// Create function context
	interp := interpreter.NewInterpreter(
		map[string]interpreter.StdlibInitializer{
			"IO": RegisterIOInEnv,
		},
		DefaultGlobalInitializers()...,
	)
	ctx := interpreter.NewFunctionContext(interp, env)

	writeMethod := bufferedClass.PublicFunctions["WRITE"]
	closeMethod := bufferedClass.PublicFunctions["CLOSE"]

	// Write some data to buffer
	_, err = writeMethod.NativeImpl(ctx, instance, []types.Value{types.StringValue("Final data")})
	require.NoError(t, err)

	// Close the buffered writer
	_, err = closeMethod.NativeImpl(ctx, instance, []types.Value{})
	require.NoError(t, err)

	// Verify buffer was flushed and underlying writer was closed
	assert.Equal(t, 1, mockWriter.writeCount, "Should have flushed buffer")
	assert.Equal(t, 1, mockWriter.closeCount, "Underlying writer should have been closed")
	assert.Equal(t, "Final data", mockWriter.data)

	// Verify buffer was cleared
	bufferData, ok := instance.NativeData.(*BufferedWriterData)
	require.True(t, ok)
	assert.Equal(t, "", bufferData.Buffer)
}
