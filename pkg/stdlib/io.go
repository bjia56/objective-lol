package stdlib

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/interpreter"
	"github.com/bjia56/objective-lol/pkg/runtime"
)

const defaultBufferSize = 1024

// BufferedReaderData stores the internal state of a BUFFERED_READER
type BufferedReaderData struct {
	Reader     *environment.ObjectInstance
	Buffer     string
	Position   int
	BufferSize int
	EOF        bool
}

// BufferedWriterData stores the internal state of a BUFFERED_WRITER
type BufferedWriterData struct {
	Writer     *environment.ObjectInstance
	Buffer     string
	BufferSize int
}

// Global IO class definitions - created once and reused
var ioClassesOnce = sync.Once{}
var ioClasses map[string]*environment.Class

func getIoClasses() map[string]*environment.Class {
	ioClassesOnce.Do(func() {
		ioClasses = map[string]*environment.Class{
			"READER": {
				Name:          "READER",
				QualifiedName: "stdlib:IO.READER",
				ModulePath:    "stdlib:IO",
				ParentClasses: []string{},
				MRO:           []string{"stdlib:IO.READER"},
				PublicFunctions: map[string]*environment.Function{
					"READ": {
						Name:       "READ",
						ReturnType: "STRIN",
						Parameters: []environment.Parameter{
							{Name: "size", Type: "INTEGR"},
						},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							return environment.NOTHIN, runtime.Exception{Message: "Not implemented"}
						},
					},
					"CLOSE": {
						Name: "CLOSE",
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							return environment.NOTHIN, nil
						},
					},
				},
				PrivateVariables: make(map[string]*environment.Variable),
				PrivateFunctions: make(map[string]*environment.Function),
				SharedVariables:  make(map[string]*environment.Variable),
				SharedFunctions:  make(map[string]*environment.Function),
				PublicVariables:  make(map[string]*environment.Variable),
			},
			"WRITER": {
				Name:          "WRITER",
				QualifiedName: "stdlib:IO.WRITER",
				ModulePath:    "stdlib:IO",
				ParentClasses: []string{},
				MRO:           []string{"stdlib:IO.WRITER"},
				PublicFunctions: map[string]*environment.Function{
					"WRITE": {
						Name:       "WRITE",
						ReturnType: "INTEGR",
						Parameters: []environment.Parameter{
							{Name: "data", Type: "STRIN"},
						},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							return environment.NOTHIN, runtime.Exception{Message: "Not implemented"}
						},
					},
					"CLOSE": {
						Name: "CLOSE",
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							return environment.NOTHIN, nil
						},
					},
				},
				PrivateVariables: make(map[string]*environment.Variable),
				PrivateFunctions: make(map[string]*environment.Function),
				SharedVariables:  make(map[string]*environment.Variable),
				SharedFunctions:  make(map[string]*environment.Function),
				PublicVariables:  make(map[string]*environment.Variable),
			},
			"BUFFERED_READER": {
				Name: "BUFFERED_READER",
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"BUFFERED_READER": {
						Name: "BUFFERED_READER",
						Parameters: []environment.Parameter{
							{Name: "reader", Type: "READER"},
						},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							reader := args[0]

							// Validate that the argument is an object with READ and CLOSE methods
							// Type assert to get the concrete ObjectInstance
							readerInstance, ok := reader.(*environment.ObjectInstance)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("BUFFERED_READER constructor: invalid object instance type")
							}

							// Check if reader has READ method
							_, err := readerInstance.GetMemberFunction("READ", "BUFFERED_READER", nil)
							if err != nil {
								return environment.NOTHIN, fmt.Errorf("BUFFERED_READER constructor: provided object does not have READ method: %v", err)
							}

							// Check if reader has CLOSE method
							_, err = readerInstance.GetMemberFunction("CLOSE", "BUFFERED_READER", nil)
							if err != nil {
								return environment.NOTHIN, fmt.Errorf("BUFFERED_READER constructor: provided object does not have CLOSE method: %v", err)
							}

							// Initialize the buffered reader data
							bufferData := &BufferedReaderData{
								Reader:     readerInstance,
								Buffer:     "",
								Position:   0,
								BufferSize: defaultBufferSize,
								EOF:        false,
							}
							this.NativeData = bufferData

							return environment.NOTHIN, nil
						},
					},
					// SET_SIZ method
					"SET_SIZ": {
						Name: "SET_SIZ",
						Parameters: []environment.Parameter{
							{Name: "newSize", Type: "INTEGR"},
						},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							newSize := args[0]

							newSizeVal, ok := newSize.(environment.IntegerValue)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("SET_SIZ expects INTEGR, got %s", newSize.Type())
							}

							size := int(newSizeVal)
							if size <= 0 {
								return environment.NOTHIN, fmt.Errorf("SET_SIZ: buffer size must be positive, got %d", size)
							}

							if bufferData, ok := this.NativeData.(*BufferedReaderData); ok {
								bufferData.BufferSize = size
								// Clear the buffer when size changes
								bufferData.Buffer = ""
								bufferData.Position = 0
								bufferData.EOF = false

								// Update SIZ variable
								if sizVar, exists := this.Variables["SIZ"]; exists {
									sizVar.Value = environment.IntegerValue(size)
								}

								return environment.NOTHIN, nil
							}
							return environment.NOTHIN, fmt.Errorf("SET_SIZ: invalid context")
						},
					},
					// READ method
					"READ": {
						Name:       "READ",
						ReturnType: "STRIN",
						Parameters: []environment.Parameter{
							{Name: "size", Type: "INTEGR"},
						},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							size := args[0]

							sizeVal, ok := size.(environment.IntegerValue)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("read expects INTEGR size, got %s", size.Type())
							}

							requestedSize := int(sizeVal)
							if requestedSize <= 0 {
								return environment.StringValue(""), nil
							}
							bufferData, ok := this.NativeData.(*BufferedReaderData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("READ: invalid context")
							}

							// If we've reached EOF, return empty string
							if bufferData.EOF {
								return environment.StringValue(""), nil
							}

							functionCtx, ok := ctx.(*interpreter.FunctionContext)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("READ: invalid function context")
							}

							result := ""

							// Keep reading until we have enough data or reach EOF
							for len(result) < requestedSize && !bufferData.EOF {
								// If buffer has available data, use it first
								availableInBuffer := len(bufferData.Buffer) - bufferData.Position
								if availableInBuffer > 0 {
									// Take what we need from buffer
									needed := requestedSize - len(result)
									toTake := needed
									if toTake > availableInBuffer {
										toTake = availableInBuffer
									}

									result += bufferData.Buffer[bufferData.Position : bufferData.Position+toTake]
									bufferData.Position += toTake

									// If we have enough, return it
									if len(result) >= requestedSize {
										break
									}
								}

								// Need more data - read from underlying reader
								// Calculate how much to request from underlying reader
								// Request at least buffer size or remaining needed, whichever is larger
								remaining := requestedSize - len(result)
								toRequest := bufferData.BufferSize
								if remaining > toRequest {
									toRequest = remaining
								}

								readResult, err := functionCtx.CallMethod(bufferData.Reader, "READ", "BUFFERED_READER", []environment.Value{environment.IntegerValue(toRequest)})
								if err != nil {
									return environment.NOTHIN, fmt.Errorf("read: error reading from underlying reader: %v", err)
								}

								// Convert result to string
								readStr, ok := readResult.(environment.StringValue)
								if !ok {
									return environment.NOTHIN, fmt.Errorf("read: underlying reader returned non-string value: %s", readResult.Type())
								}

								// If we got empty string, we've reached EOF
								if len(readStr) == 0 {
									bufferData.EOF = true
									break
								}

								// Store new data in buffer and reset position
								bufferData.Buffer = string(readStr)
								bufferData.Position = 0
							}

							return environment.StringValue(result), nil
						},
					},
					// CLOSE method
					"CLOSE": {
						Name: "CLOSE",
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							bufferData, ok := this.NativeData.(*BufferedReaderData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("CLOSE: invalid context")
							}

							functionCtx, ok := ctx.(*interpreter.FunctionContext)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("CLOSE: invalid function context")
							}

							// Call the CLOSE method on the underlying reader
							_, err := functionCtx.CallMethod(bufferData.Reader, "CLOSE", "BUFFERED_READER", []environment.Value{})
							if err != nil {
								return environment.NOTHIN, fmt.Errorf("CLOSE: error closing underlying reader: %v", err)
							}

							// Clear buffer data
							bufferData.Buffer = ""
							bufferData.Position = 0
							bufferData.EOF = true

							return environment.NOTHIN, nil
						},
					},
				},
				PublicVariables: map[string]*environment.Variable{
					"SIZ": {
						Name:     "SIZ",
						Type:     "INTEGR",
						Value:    environment.IntegerValue(defaultBufferSize),
						IsLocked: false,
						IsPublic: true,
					},
				},
				QualifiedName:    "stdlib:IO.BUFFERED_READER",
				ModulePath:       "stdlib:IO",
				ParentClasses:    []string{"stdlib:IO.READER"},
				MRO:              []string{"stdlib:IO.BUFFERED_READER", "stdlib:IO.READER"},
				PrivateVariables: make(map[string]*environment.Variable),
				PrivateFunctions: make(map[string]*environment.Function),
				SharedVariables:  make(map[string]*environment.Variable),
				SharedFunctions:  make(map[string]*environment.Function),
			},
			"BUFFERED_WRITER": {
				Name:          "BUFFERED_WRITER",
				QualifiedName: "stdlib:IO.BUFFERED_WRITER",
				ModulePath:    "stdlib:IO",
				ParentClasses: []string{"stdlib:IO.WRITER"},
				MRO:           []string{"stdlib:IO.BUFFERED_WRITER", "stdlib:IO.WRITER"},
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"BUFFERED_WRITER": {
						Name: "BUFFERED_WRITER",
						Parameters: []environment.Parameter{
							{Name: "writer", Type: "WRITER"},
						},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							writer := args[0]

							// Validate that the argument is an object with WRITE and CLOSE methods
							writerInstance, ok := writer.(*environment.ObjectInstance)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("BUFFERED_WRITER constructor expects WRITER object, got %s", args[0].Type())
							}

							// Check if writer has WRITE method
							_, err := writerInstance.GetMemberFunction("WRITE", "BUFFERED_WRITER", nil)
							if err != nil {
								return environment.NOTHIN, fmt.Errorf("BUFFERED_WRITER constructor: provided object does not have WRITE method: %v", err)
							}

							// Check if writer has CLOSE method
							_, err = writerInstance.GetMemberFunction("CLOSE", "BUFFERED_WRITER", nil)
							if err != nil {
								return environment.NOTHIN, fmt.Errorf("BUFFERED_WRITER constructor: provided object does not have CLOSE method: %v", err)
							}

							// Initialize the buffered writer data
							bufferData := &BufferedWriterData{
								Writer:     writerInstance,
								Buffer:     "",
								BufferSize: defaultBufferSize,
							}
							this.NativeData = bufferData

							return environment.NOTHIN, nil
						},
					},
					// SET_SIZ method
					"SET_SIZ": {
						Name: "SET_SIZ",
						Parameters: []environment.Parameter{
							{Name: "newSize", Type: "INTEGR"},
						},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							newSize := args[0]

							newSizeVal, ok := newSize.(environment.IntegerValue)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("SET_SIZ expects INTEGR, got %s", args[0].Type())
							}

							size := int(newSizeVal)
							if size <= 0 {
								return environment.NOTHIN, fmt.Errorf("SET_SIZ: buffer size must be positive, got %d", size)
							}

							if bufferData, ok := this.NativeData.(*BufferedWriterData); ok {
								functionCtx, ok := ctx.(*interpreter.FunctionContext)
								if !ok {
									return environment.NOTHIN, fmt.Errorf("SET_SIZ: invalid function context")
								}

								// Flush existing buffer before changing size
								if len(bufferData.Buffer) > 0 {
									_, err := functionCtx.CallMethod(bufferData.Writer, "WRITE", "BUFFERED_WRITER", []environment.Value{environment.StringValue(bufferData.Buffer)})
									if err != nil {
										return environment.NOTHIN, fmt.Errorf("SET_SIZ: error flushing buffer: %v", err)
									}
									bufferData.Buffer = ""
								}

								bufferData.BufferSize = size

								// Update SIZ variable
								if sizVar, exists := this.Variables["SIZ"]; exists {
									sizVar.Value = environment.IntegerValue(size)
								}

								return environment.NOTHIN, nil
							}
							return environment.NOTHIN, fmt.Errorf("SET_SIZ: invalid context")
						},
					},
					// WRITE method
					"WRITE": {
						Name:       "WRITE",
						ReturnType: "INTEGR",
						Parameters: []environment.Parameter{
							{Name: "data", Type: "STRIN"},
						},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							data := args[0]

							dataVal, ok := data.(environment.StringValue)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("WRITE expects STRIN data, got %s", data.Type())
							}

							dataBuffer := string(dataVal)
							originalLength := len(dataBuffer)

							bufferData, ok := this.NativeData.(*BufferedWriterData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("WRITE: invalid context")
							}

							functionCtx, ok := ctx.(*interpreter.FunctionContext)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("WRITE: invalid function context")
							}

							// If data fits in remaining buffer space, just buffer it
							if len(bufferData.Buffer)+len(dataBuffer) <= bufferData.BufferSize {
								bufferData.Buffer += dataBuffer
								return environment.IntegerValue(originalLength), nil
							}

							// Buffer is full or will be full, need to flush
							if len(bufferData.Buffer) > 0 {
								_, err := functionCtx.CallMethod(bufferData.Writer, "WRITE", "BUFFERED_WRITER", []environment.Value{environment.StringValue(bufferData.Buffer)})
								if err != nil {
									return environment.NOTHIN, fmt.Errorf("WRITE: error flushing buffer: %v", err)
								}
								bufferData.Buffer = ""
							}

							// If data is larger than buffer size, write it directly
							if len(dataBuffer) >= bufferData.BufferSize {
								_, err := functionCtx.CallMethod(bufferData.Writer, "WRITE", "BUFFERED_WRITER", []environment.Value{environment.StringValue(dataBuffer)})
								if err != nil {
									return environment.NOTHIN, fmt.Errorf("WRITE: error writing large data: %v", err)
								}
								return environment.IntegerValue(originalLength), nil
							}

							// Otherwise, buffer the data
							bufferData.Buffer = dataBuffer
							return environment.IntegerValue(originalLength), nil
						},
					},
					// FLUSH method
					"FLUSH": {
						Name: "FLUSH",
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							bufferData, ok := this.NativeData.(*BufferedWriterData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("FLUSH: invalid context")
							}

							functionCtx, ok := ctx.(*interpreter.FunctionContext)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("FLUSH: invalid function context")
							}

							// Flush buffer if it has data
							if len(bufferData.Buffer) > 0 {
								_, err := functionCtx.CallMethod(bufferData.Writer, "WRITE", "BUFFERED_WRITER", []environment.Value{environment.StringValue(bufferData.Buffer)})
								if err != nil {
									return environment.NOTHIN, fmt.Errorf("FLUSH: error flushing buffer: %v", err)
								}
								bufferData.Buffer = ""
							}

							return environment.NOTHIN, nil
						},
					},
					// CLOSE method
					"CLOSE": {
						Name: "CLOSE",
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							bufferData, ok := this.NativeData.(*BufferedWriterData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("CLOSE: invalid context")
							}

							functionCtx, ok := ctx.(*interpreter.FunctionContext)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("CLOSE: invalid function context")
							}

							// Flush buffer before closing
							if len(bufferData.Buffer) > 0 {
								_, err := functionCtx.CallMethod(bufferData.Writer, "WRITE", "BUFFERED_WRITER", []environment.Value{environment.StringValue(bufferData.Buffer)})
								if err != nil {
									return environment.NOTHIN, fmt.Errorf("CLOSE: error flushing buffer: %v", err)
								}
								bufferData.Buffer = ""
							}

							// Call the CLOSE method on the underlying writer
							_, err := functionCtx.CallMethod(bufferData.Writer, "CLOSE", "BUFFERED_WRITER", []environment.Value{})
							if err != nil {
								return environment.NOTHIN, fmt.Errorf("CLOSE: error closing underlying writer: %v", err)
							}

							return environment.NOTHIN, nil
						},
					},
				},
				PublicVariables: map[string]*environment.Variable{
					"SIZ": {
						Name:     "SIZ",
						Type:     "INTEGR",
						Value:    environment.IntegerValue(defaultBufferSize),
						IsLocked: false,
						IsPublic: true,
					},
				},
				PrivateVariables: make(map[string]*environment.Variable),
				PrivateFunctions: make(map[string]*environment.Function),
				SharedVariables:  make(map[string]*environment.Variable),
				SharedFunctions:  make(map[string]*environment.Function),
			},
		}
	})
	return ioClasses
}

// RegisterIOInEnv registers IO classes in the given environment
// declarations: empty slice means import all, otherwise import only specified classes
func RegisterIOInEnv(env *environment.Environment, declarations ...string) error {
	ioClasses := getIoClasses()

	// If declarations is empty, import all classes
	if len(declarations) == 0 {
		for _, class := range ioClasses {
			env.DefineClass(class)
		}
		return nil
	}

	// Otherwise, import only specified classes
	for _, decl := range declarations {
		declUpper := strings.ToUpper(decl)
		if class, exists := ioClasses[declUpper]; exists {
			env.DefineClass(class)
		} else {
			return fmt.Errorf("unknown IO class: %s", decl)
		}
	}

	return nil
}
