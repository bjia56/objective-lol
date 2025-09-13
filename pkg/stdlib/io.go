package stdlib

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bjia56/objective-lol/pkg/environment"
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
				Name: "READER",
				Documentation: []string{
					"Abstract base class for objects that can read data.",
					"Defines the interface for reading operations with READ and CLOSE methods.",
				},
				QualifiedName: "stdlib:IO.READER",
				ModulePath:    "stdlib:IO",
				ParentClasses: []string{},
				MRO:           []string{"stdlib:IO.READER"},
				PublicFunctions: map[string]*environment.Function{
					"READ": {
						Name: "READ",
						Documentation: []string{
							"Reads up to the specified number of characters from the input source.",
							"Returns empty string when end-of-file is reached.",
						},
						ReturnType: "STRIN",
						Parameters: []environment.Parameter{
							{Name: "size", Type: "INTEGR"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							return environment.NOTHIN, runtime.Exception{Message: "not implemented"}
						},
					},
					"CLOSE": {
						Name: "CLOSE",
						Documentation: []string{
							"Closes the reader and releases any associated resources.",
							"Should be called when done reading to ensure proper cleanup.",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							return environment.NOTHIN, nil
						},
					},
				},
				PrivateVariables: make(map[string]*environment.MemberVariable),
				PrivateFunctions: make(map[string]*environment.Function),
				SharedVariables:  make(map[string]*environment.MemberVariable),
				SharedFunctions:  make(map[string]*environment.Function),
				PublicVariables:  make(map[string]*environment.MemberVariable),
			},
			"WRITER": {
				Name: "WRITER",
				Documentation: []string{
					"Abstract base class for objects that can write data.",
					"Defines the interface for writing operations with WRITE and CLOSE methods.",
				},
				QualifiedName: "stdlib:IO.WRITER",
				ModulePath:    "stdlib:IO",
				ParentClasses: []string{},
				MRO:           []string{"stdlib:IO.WRITER"},
				PublicFunctions: map[string]*environment.Function{
					"WRITE": {
						Name: "WRITE",
						Documentation: []string{
							"Writes string data to the output destination.",
							"Returns the number of characters written.",
						},
						ReturnType: "INTEGR",
						Parameters: []environment.Parameter{
							{Name: "data", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							return environment.NOTHIN, runtime.Exception{Message: "not implemented"}
						},
					},
					"CLOSE": {
						Name: "CLOSE",
						Documentation: []string{
							"Closes the writer and releases any associated resources.",
							"Should be called when done writing to ensure proper cleanup.",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							return environment.NOTHIN, nil
						},
					},
				},
				PrivateVariables: make(map[string]*environment.MemberVariable),
				PrivateFunctions: make(map[string]*environment.Function),
				SharedVariables:  make(map[string]*environment.MemberVariable),
				SharedFunctions:  make(map[string]*environment.Function),
				PublicVariables:  make(map[string]*environment.MemberVariable),
			},
			"READWRITER": {
				Name: "READWRITER",
				Documentation: []string{
					"Interface that combines both READER and WRITER capabilities.",
					"Inherits READ and CLOSE from READER, and WRITE from WRITER.",
				},
				QualifiedName:    "stdlib:IO.READWRITER",
				ModulePath:       "stdlib:IO",
				ParentClasses:    []string{},
				MRO:              []string{"stdlib:IO.READWRITER", "stdlib:IO.READER", "stdlib:IO.WRITER"},
				PublicFunctions:  make(map[string]*environment.Function),
				PrivateVariables: make(map[string]*environment.MemberVariable),
				PrivateFunctions: make(map[string]*environment.Function),
				SharedVariables:  make(map[string]*environment.MemberVariable),
				SharedFunctions:  make(map[string]*environment.Function),
				PublicVariables:  make(map[string]*environment.MemberVariable),
			},
			"BUFFERED_READER": {
				Name: "BUFFERED_READER",
				Documentation: []string{
					"Buffered reader that wraps another READER object for improved performance.",
					"Reduces the number of actual I/O operations by buffering data internally.",
				},
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"BUFFERED_READER": {
						Name: "BUFFERED_READER",
						Documentation: []string{
							"Initializes a BUFFERED_READER with the given READER object.",
							"Default buffer size is 1024.",
						},
						Parameters: []environment.Parameter{
							{Name: "reader", Type: "READER"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							reader := args[0]

							// Validate that the argument is an object with READ and CLOSE methods
							// Type assert to get the concrete ObjectInstance
							readerInstance, err := reader.Cast("READER")
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("BUFFERED_READER constructor expects READER object, got %s", reader.Type())}
							}

							// Initialize the buffered reader data
							bufferData := &BufferedReaderData{
								Reader:     readerInstance.(*environment.ObjectInstance),
								Buffer:     "",
								Position:   0,
								BufferSize: defaultBufferSize,
								EOF:        false,
							}
							this.NativeData = bufferData

							return environment.NOTHIN, nil
						},
					},
					// READ method
					"READ": {
						Name: "READ",
						Documentation: []string{
							"Reads up to the specified number of characters from the buffered reader.",
							"Utilizes internal buffer to minimize calls to the underlying READER.",
							"Returns empty string when end-of-file is reached.",
						},
						ReturnType: "STRIN",
						Parameters: []environment.Parameter{
							{Name: "size", Type: "INTEGR"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							size := args[0]

							sizeVal, ok := size.(environment.IntegerValue)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("READ expects INTEGR size, got %s", size.Type())}
							}

							requestedSize := int(sizeVal)
							if requestedSize <= 0 {
								return environment.StringValue(""), nil
							}
							bufferData, ok := this.NativeData.(*BufferedReaderData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "READ: invalid context"}
							}

							// If we've reached EOF, return empty string
							if bufferData.EOF {
								return environment.StringValue(""), nil
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

								readResult, err := interpreter.CallMemberFunction(bufferData.Reader, "READ", []environment.Value{environment.IntegerValue(toRequest)})
								if err != nil {
									return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("READ: error reading from underlying reader: %v", err)}
								}

								// Convert result to string
								readStr, ok := readResult.(environment.StringValue)
								if !ok {
									return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("READ: underlying reader returned non-string value: %s", readResult.Type())}
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
						Documentation: []string{
							"Closes the buffered reader and the underlying READER.",
							"Releases any associated resources.",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							bufferData, ok := this.NativeData.(*BufferedReaderData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "CLOSE: invalid context"}
							}

							// Call the CLOSE method on the underlying reader
							_, err := interpreter.CallMemberFunction(bufferData.Reader, "CLOSE", []environment.Value{})
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("CLOSE: error closing underlying reader: %v", err)}
							}

							// Clear buffer data
							bufferData.Buffer = ""
							bufferData.Position = 0
							bufferData.EOF = true

							return environment.NOTHIN, nil
						},
					},
				},
				PublicVariables: map[string]*environment.MemberVariable{
					"SIZ": {
						Variable: environment.Variable{
							Name: "SIZ",
							Documentation: []string{
								"Gets or sets the size of the internal buffer used for reading.",
								"Setting a new size will clear the existing buffer.",
								"Default is 1024. Must be a positive integer.",
							},
							Type:     "INTEGR",
							IsLocked: false,
							IsPublic: true,
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							if bufferData, ok := this.NativeData.(*BufferedReaderData); ok {
								return environment.IntegerValue(bufferData.BufferSize), nil
							}
							return environment.NOTHIN, fmt.Errorf("SIZ: invalid context")
						},
						NativeSet: func(this *environment.ObjectInstance, newValue environment.Value) error {
							val, err := newValue.Cast("INTEGR")
							if err != nil {
								return runtime.Exception{Message: fmt.Sprintf("SIZ expects INTEGR, got %s", newValue.Type())}
							}

							size := int(val.(environment.IntegerValue))
							if size <= 0 {
								return runtime.Exception{Message: fmt.Sprintf("SIZ: buffer size must be positive, got %d", size)}
							}

							if bufferData, ok := this.NativeData.(*BufferedReaderData); ok {
								bufferData.BufferSize = size
								// Clear the buffer when size changes
								bufferData.Buffer = ""
								bufferData.Position = 0
								bufferData.EOF = false
								return nil
							}
							return runtime.Exception{Message: "SIZ: invalid context"}
						},
					},
				},
				QualifiedName:    "stdlib:IO.BUFFERED_READER",
				ModulePath:       "stdlib:IO",
				ParentClasses:    []string{"stdlib:IO.READER"},
				MRO:              []string{"stdlib:IO.BUFFERED_READER", "stdlib:IO.READER"},
				PrivateVariables: make(map[string]*environment.MemberVariable),
				PrivateFunctions: make(map[string]*environment.Function),
				SharedVariables:  make(map[string]*environment.MemberVariable),
				SharedFunctions:  make(map[string]*environment.Function),
			},
			"BUFFERED_WRITER": {
				Name: "BUFFERED_WRITER",
				Documentation: []string{
					"Buffered writer that wraps another WRITER object for improved performance.",
					"Reduces the number of actual I/O operations by buffering data internally.",
				},
				QualifiedName: "stdlib:IO.BUFFERED_WRITER",
				ModulePath:    "stdlib:IO",
				ParentClasses: []string{"stdlib:IO.WRITER"},
				MRO:           []string{"stdlib:IO.BUFFERED_WRITER", "stdlib:IO.WRITER"},
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"BUFFERED_WRITER": {
						Name: "BUFFERED_WRITER",
						Documentation: []string{
							"Initializes a BUFFERED_WRITER with the given WRITER object.",
							"Default buffer size is 1024.",
						},
						Parameters: []environment.Parameter{
							{Name: "writer", Type: "WRITER"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							writer := args[0]

							// Validate that the argument is an object with WRITE and CLOSE methods
							writerInstance, err := writer.Cast("WRITER")
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("BUFFERED_WRITER constructor expects WRITER object, got %s", writer.Type())}
							}

							// Initialize the buffered writer data
							bufferData := &BufferedWriterData{
								Writer:     writerInstance.(*environment.ObjectInstance),
								Buffer:     "",
								BufferSize: defaultBufferSize,
							}
							this.NativeData = bufferData

							return environment.NOTHIN, nil
						},
					},
					// WRITE method
					"WRITE": {
						Name: "WRITE",
						Documentation: []string{
							"Writes string data to the buffered writer.",
							"Data is buffered internally and flushed to the underlying WRITER when the buffer is full or when FLUSH/CLOSE is called.",
							"Returns the number of characters written.",
						},
						ReturnType: "INTEGR",
						Parameters: []environment.Parameter{
							{Name: "data", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							data := args[0]

							dataVal, ok := data.(environment.StringValue)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("WRITE expects STRING data, got %s", data.Type())}
							}

							dataBuffer := string(dataVal)
							originalLength := len(dataBuffer)

							bufferData, ok := this.NativeData.(*BufferedWriterData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "WRITE: invalid context"}
							}

							// If data fits in remaining buffer space, just buffer it
							if len(bufferData.Buffer)+len(dataBuffer) <= bufferData.BufferSize {
								bufferData.Buffer += dataBuffer
								return environment.IntegerValue(originalLength), nil
							}

							// Buffer is full or will be full, need to flush
							if len(bufferData.Buffer) > 0 {
								_, err := interpreter.CallMemberFunction(bufferData.Writer, "WRITE", []environment.Value{environment.StringValue(bufferData.Buffer)})
								if err != nil {
									return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("WRITE: error flushing buffer: %v", err)}
								}
								bufferData.Buffer = ""
							}

							// If data is larger than buffer size, write it directly
							if len(dataBuffer) >= bufferData.BufferSize {
								_, err := interpreter.CallMemberFunction(bufferData.Writer, "WRITE", []environment.Value{environment.StringValue(dataBuffer)})
								if err != nil {
									return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("WRITE: error writing large data: %v", err)}
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
						Documentation: []string{
							"Flushes the internal buffer to the underlying WRITER.",
							"Should be called to ensure all buffered data is written out.",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							bufferData, ok := this.NativeData.(*BufferedWriterData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "FLUSH: invalid context"}
							}

							// Flush buffer if it has data
							if len(bufferData.Buffer) > 0 {
								_, err := interpreter.CallMemberFunction(bufferData.Writer, "WRITE", []environment.Value{environment.StringValue(bufferData.Buffer)})
								if err != nil {
									return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("FLUSH: error flushing buffer: %v", err)}
								}
								bufferData.Buffer = ""
							}

							return environment.NOTHIN, nil
						},
					},
					// CLOSE method
					"CLOSE": {
						Name: "CLOSE",
						Documentation: []string{
							"Closes the buffered writer and the underlying WRITER.",
							"Flushes any remaining buffered data before closing.",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							bufferData, ok := this.NativeData.(*BufferedWriterData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "CLOSE: invalid context"}
							}

							// Flush buffer before closing
							if len(bufferData.Buffer) > 0 {
								_, err := interpreter.CallMemberFunction(bufferData.Writer, "WRITE", []environment.Value{environment.StringValue(bufferData.Buffer)})
								if err != nil {
									return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("CLOSE: error flushing buffer: %v", err)}
								}
								bufferData.Buffer = ""
							}

							// Call the CLOSE method on the underlying writer
							_, err := interpreter.CallMemberFunction(bufferData.Writer, "CLOSE", []environment.Value{})
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("CLOSE: error closing underlying writer: %v", err)}
							}

							return environment.NOTHIN, nil
						},
					},
				},
				PublicVariables: map[string]*environment.MemberVariable{
					"SIZ": {
						Variable: environment.Variable{
							Name: "SIZ",
							Documentation: []string{
								"Gets or sets the size of the internal buffer used for writing.",
								"Setting a new size will drop the existing buffer.",
								"Default is 1024. Must be a positive integer.",
							},
							Type:     "INTEGR",
							IsLocked: false,
							IsPublic: true,
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							if bufferData, ok := this.NativeData.(*BufferedWriterData); ok {
								return environment.IntegerValue(bufferData.BufferSize), nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "SIZ: invalid context"}
						},
						NativeSet: func(this *environment.ObjectInstance, newValue environment.Value) error {
							val, err := newValue.Cast("INTEGR")
							if err != nil {
								return runtime.Exception{Message: fmt.Sprintf("SIZ expects INTEGR, got %s", newValue.Type())}
							}

							size := int(val.(environment.IntegerValue))
							if size <= 0 {
								return runtime.Exception{Message: fmt.Sprintf("SIZ: buffer size must be positive, got %d", size)}
							}

							if bufferData, ok := this.NativeData.(*BufferedWriterData); ok {
								// Flush existing buffer before changing size
								if len(bufferData.Buffer) > 0 {
									// TODO: pass interpreter into NativeGet/Set?
									// interpreter.CallMemberFunction(bufferData.Writer, "WRITE", []environment.Value{environment.StringValue(bufferData.Buffer)})
									bufferData.Buffer = ""
								}

								bufferData.BufferSize = size
								return nil
							}
							return runtime.Exception{Message: "SIZ: invalid context"}
						},
					},
				},
				PrivateVariables: make(map[string]*environment.MemberVariable),
				PrivateFunctions: make(map[string]*environment.Function),
				SharedVariables:  make(map[string]*environment.MemberVariable),
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
			return runtime.Exception{Message: fmt.Sprintf("unknown IO declaration: %s", decl)}
		}
	}

	return nil
}
