package stdlib

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bjia56/objective-lol/pkg/ast"
	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/interpreter"
	"github.com/bjia56/objective-lol/pkg/types"
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

// Global IO class definitions - created once and reused
var ioClassesOnce = sync.Once{}
var ioClasses map[string]*environment.Class

func getIoClasses() map[string]*environment.Class {
	ioClassesOnce.Do(func() {
		ioClasses = map[string]*environment.Class{
			"READER": {
				Name: "READER",
				PublicFunctions: map[string]*environment.Function{
					"READ": {
						Name:       "READ",
						ReturnType: "STRIN",
						Parameters: []environment.Parameter{
							{Name: "size", Type: "INTEGR"},
						},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							return types.NOTHIN, ast.Exception{Message: "Not implemented"}
						},
					},
					"CLOSE": {
						Name: "CLOSE",
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							return types.NOTHIN, nil
						},
					},
				},
			},
			"WRITER": {
				Name: "WRITER",
				PublicFunctions: map[string]*environment.Function{
					"WRITE": {
						Name:       "WRITE",
						ReturnType: "INTEGR",
						Parameters: []environment.Parameter{
							{Name: "data", Type: "STRIN"},
						},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							return types.NOTHIN, ast.Exception{Message: "Not implemented"}
						},
					},
					"CLOSE": {
						Name: "CLOSE",
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							return types.NOTHIN, nil
						},
					},
				},
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
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							reader := args[0]

							// Validate that the argument is an object with READ and CLOSE methods
							readerObj, ok := reader.(types.ObjectValue)
							if !ok {
								return types.NOTHIN, fmt.Errorf("BUFFERED_READER constructor expects READER object, got %s", args[0].Type())
							}

							// Type assert to get the concrete ObjectInstance
							readerInstance, ok := readerObj.Instance.(*environment.ObjectInstance)
							if !ok {
								return types.NOTHIN, fmt.Errorf("BUFFERED_READER constructor: invalid object instance type")
							}

							// Check if reader has READ method
							_, err := readerInstance.GetMemberFunction("READ", "BUFFERED_READER", nil)
							if err != nil {
								return types.NOTHIN, fmt.Errorf("BUFFERED_READER constructor: provided object does not have READ method: %v", err)
							}

							// Check if reader has CLOSE method
							_, err = readerInstance.GetMemberFunction("CLOSE", "BUFFERED_READER", nil)
							if err != nil {
								return types.NOTHIN, fmt.Errorf("BUFFERED_READER constructor: provided object does not have CLOSE method: %v", err)
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

							return types.NOTHIN, nil
						},
					},
					// SET_SIZ method
					"SET_SIZ": {
						Name: "SET_SIZ",
						Parameters: []environment.Parameter{
							{Name: "new_size", Type: "INTEGR"},
						},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							if len(args) != 1 {
								return types.NOTHIN, fmt.Errorf("SET_SIZ expects 1 argument, got %d", len(args))
							}

							newSizeVal, ok := args[0].(types.IntegerValue)
							if !ok {
								return types.NOTHIN, fmt.Errorf("SET_SIZ expects INTEGR, got %s", args[0].Type())
							}

							newSize := int(newSizeVal)
							if newSize <= 0 {
								return types.NOTHIN, fmt.Errorf("SET_SIZ: buffer size must be positive, got %d", newSize)
							}

							if bufferData, ok := this.NativeData.(*BufferedReaderData); ok {
								bufferData.BufferSize = newSize
								// Clear the buffer when size changes
								bufferData.Buffer = ""
								bufferData.Position = 0
								bufferData.EOF = false

								// Update SIZ variable
								if sizVar, exists := this.Variables["SIZ"]; exists {
									sizVar.Value = types.IntegerValue(newSize)
								}

								return types.NOTHIN, nil
							}
							return types.NOTHIN, fmt.Errorf("SET_SIZ: invalid context")
						},
					},
					// READ method
					"READ": {
						Name:       "READ",
						ReturnType: "STRIN",
						Parameters: []environment.Parameter{
							{Name: "size", Type: "INTEGR"},
						},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							if len(args) != 1 {
								return types.NOTHIN, fmt.Errorf("READ expects 1 argument, got %d", len(args))
							}

							sizeVal, ok := args[0].(types.IntegerValue)
							if !ok {
								return types.NOTHIN, fmt.Errorf("read expects INTEGR size, got %s", args[0].Type())
							}

							requestedSize := int(sizeVal)
							if requestedSize <= 0 {
								return types.StringValue(""), nil
							}
							bufferData, ok := this.NativeData.(*BufferedReaderData)
							if !ok {
								return types.NOTHIN, fmt.Errorf("READ: invalid context")
							}

							// If we've reached EOF, return empty string
							if bufferData.EOF {
								return types.StringValue(""), nil
							}

							functionCtx, ok := ctx.(*interpreter.FunctionContext)
							if !ok {
								return types.NOTHIN, fmt.Errorf("READ: invalid function context")
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

								readResult, err := functionCtx.CallMethod(bufferData.Reader, "READ", "BUFFERED_READER", []types.Value{types.IntegerValue(toRequest)})
								if err != nil {
									return types.NOTHIN, fmt.Errorf("read: error reading from underlying reader: %v", err)
								}

								// Convert result to string
								readStr, ok := readResult.(types.StringValue)
								if !ok {
									return types.NOTHIN, fmt.Errorf("read: underlying reader returned non-string value: %s", readResult.Type())
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

							return types.StringValue(result), nil
						},
					},
					// CLOSE method
					"CLOSE": {
						Name: "CLOSE",
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							bufferData, ok := this.NativeData.(*BufferedReaderData)
							if !ok {
								return types.NOTHIN, fmt.Errorf("CLOSE: invalid context")
							}

							functionCtx, ok := ctx.(*interpreter.FunctionContext)
							if !ok {
								return types.NOTHIN, fmt.Errorf("CLOSE: invalid function context")
							}

							// Call the CLOSE method on the underlying reader
							_, err := functionCtx.CallMethod(bufferData.Reader, "CLOSE", "BUFFERED_READER", []types.Value{})
							if err != nil {
								return types.NOTHIN, fmt.Errorf("CLOSE: error closing underlying reader: %v", err)
							}

							// Clear buffer data
							bufferData.Buffer = ""
							bufferData.Position = 0
							bufferData.EOF = true

							return types.NOTHIN, nil
						},
					},
				},
				PublicVariables: map[string]*environment.Variable{
					"SIZ": {
						Name:     "SIZ",
						Type:     "INTEGR",
						Value:    types.IntegerValue(defaultBufferSize),
						IsLocked: false,
						IsPublic: true,
					},
				},
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
