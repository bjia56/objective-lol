package stdlib

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/runtime"
)

// DocumentData stores the internal state of a DOCUMENT
type DocumentData struct {
	FilePath string
	FileMode string
	File     *os.File
	IsOpen   bool
}

// Global FILE class definitions - created once and reused
var fileClassesOnce = sync.Once{}
var fileClasses map[string]*environment.Class

func getFileClasses() map[string]*environment.Class {
	fileClassesOnce.Do(func() {
		fileClasses = map[string]*environment.Class{
			"DOCUMENT": {
				Name:          "DOCUMENT",
				QualifiedName: "stdlib:FILE.DOCUMENT",
				ModulePath:    "stdlib:FILE",
				ParentClasses: []string{"stdlib:IO.READWRITER"},
				MRO:           []string{"stdlib:FILE.DOCUMENT", "stdlib:IO.READWRITER", "stdlib:IO.READER", "stdlib:IO.WRITER"},
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"DOCUMENT": {
						Name: "DOCUMENT",
						Parameters: []environment.Parameter{
							{Name: "path", Type: "STRIN"},
							{Name: "mode", Type: "STRIN"},
						},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							path := args[0]
							mode := args[1]

							pathVal, ok := path.(environment.StringValue)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("DOCUMENT constructor expects STRIN path, got %s", path.Type())
							}

							modeVal, ok := mode.(environment.StringValue)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("DOCUMENT constructor expects STRIN mode, got %s", mode.Type())
							}

							// Validate mode
							modeStr := strings.ToUpper(string(modeVal))
							if modeStr != "R" && modeStr != "W" && modeStr != "RW" && modeStr != "A" {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("Invalid file mode: %s. Valid modes: R, W, RW, A", modeStr)}
							}

							// Initialize the document data
							docData := &DocumentData{
								FilePath: string(pathVal),
								FileMode: modeStr,
								File:     nil,
								IsOpen:   false,
							}
							this.NativeData = docData

							// Set public variables
							this.Variables["PATH"] = &environment.Variable{
								Name:     "PATH",
								Type:     "STRIN",
								Value:    pathVal,
								IsLocked: true,
								IsPublic: true,
							}
							this.Variables["MODE"] = &environment.Variable{
								Name:     "MODE",
								Type:     "STRIN",
								Value:    environment.StringValue(modeStr),
								IsLocked: true,
								IsPublic: true,
							}
							this.Variables["IS_OPEN"] = &environment.Variable{
								Name:     "IS_OPEN",
								Type:     "BOOL",
								Value:    environment.NO,
								IsLocked: false,
								IsPublic: true,
							}

							return environment.NOTHIN, nil
						},
					},
					// OPEN method
					"OPEN": {
						Name: "OPEN",
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							docData, ok := this.NativeData.(*DocumentData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("OPEN: invalid context")
							}

							if docData.IsOpen {
								return environment.NOTHIN, runtime.Exception{Message: "File is already open"}
							}

							var flag int
							switch docData.FileMode {
							case "R":
								flag = os.O_RDONLY
							case "W":
								flag = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
							case "RW":
								flag = os.O_RDWR | os.O_CREATE
							case "A":
								flag = os.O_WRONLY | os.O_CREATE | os.O_APPEND
							}

							file, err := os.OpenFile(docData.FilePath, flag, 0644)
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("Failed to open file: %v", err)}
							}

							docData.File = file
							docData.IsOpen = true

							// Update IS_OPEN variable
							if isOpenVar, exists := this.Variables["IS_OPEN"]; exists {
								isOpenVar.Value = environment.YEZ
							}

							return environment.NOTHIN, nil
						},
					},
					// READ method (implements READER interface)
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
								return environment.NOTHIN, fmt.Errorf("READ expects INTEGR size, got %s", size.Type())
							}

							docData, ok := this.NativeData.(*DocumentData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("READ: invalid context")
							}

							if !docData.IsOpen {
								return environment.NOTHIN, runtime.Exception{Message: "File is not open"}
							}

							if docData.FileMode != "R" && docData.FileMode != "RW" {
								return environment.NOTHIN, runtime.Exception{Message: "File is not open for reading"}
							}

							readSize := int(sizeVal)
							if readSize <= 0 {
								return environment.StringValue(""), nil
							}

							buffer := make([]byte, readSize)
							n, err := docData.File.Read(buffer)
							if err != nil && err.Error() != "EOF" {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("Read error: %v", err)}
							}

							return environment.StringValue(string(buffer[:n])), nil
						},
					},
					// WRITE method (implements WRITER interface)
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

							docData, ok := this.NativeData.(*DocumentData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("WRITE: invalid context")
							}

							if !docData.IsOpen {
								return environment.NOTHIN, runtime.Exception{Message: "File is not open"}
							}

							if docData.FileMode != "W" && docData.FileMode != "RW" && docData.FileMode != "A" {
								return environment.NOTHIN, runtime.Exception{Message: "File is not open for writing"}
							}

							dataStr := string(dataVal)
							n, err := docData.File.WriteString(dataStr)
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("Write error: %v", err)}
							}

							return environment.IntegerValue(n), nil
						},
					},
					// SEEK method
					"SEEK": {
						Name: "SEEK",
						Parameters: []environment.Parameter{
							{Name: "position", Type: "INTEGR"},
						},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							position := args[0]

							posVal, ok := position.(environment.IntegerValue)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("SEEK expects INTEGR position, got %s", position.Type())
							}

							docData, ok := this.NativeData.(*DocumentData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("SEEK: invalid context")
							}

							if !docData.IsOpen {
								return environment.NOTHIN, runtime.Exception{Message: "File is not open"}
							}

							_, err := docData.File.Seek(int64(posVal), 0)
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("Seek error: %v", err)}
							}

							return environment.NOTHIN, nil
						},
					},
					// TELL method
					"TELL": {
						Name:       "TELL",
						ReturnType: "INTEGR",
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							docData, ok := this.NativeData.(*DocumentData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("TELL: invalid context")
							}

							if !docData.IsOpen {
								return environment.NOTHIN, runtime.Exception{Message: "File is not open"}
							}

							pos, err := docData.File.Seek(0, 1) // Get current position
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("Tell error: %v", err)}
							}

							return environment.IntegerValue(int(pos)), nil
						},
					},
					// SIZE method
					"SIZE": {
						Name:       "SIZE",
						ReturnType: "INTEGR",
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							docData, ok := this.NativeData.(*DocumentData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("SIZE: invalid context")
							}

							stat, err := os.Stat(docData.FilePath)
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("Failed to get file size: %v", err)}
							}

							return environment.IntegerValue(int(stat.Size())), nil
						},
					},
					// EXISTS method
					"EXISTS": {
						Name:       "EXISTS",
						ReturnType: "BOOL",
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							docData, ok := this.NativeData.(*DocumentData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("EXISTS: invalid context")
							}

							_, err := os.Stat(docData.FilePath)
							exists := err == nil
							return environment.BoolValue(exists), nil
						},
					},
					// FLUSH method
					"FLUSH": {
						Name: "FLUSH",
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							docData, ok := this.NativeData.(*DocumentData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("FLUSH: invalid context")
							}

							if !docData.IsOpen {
								return environment.NOTHIN, runtime.Exception{Message: "File is not open"}
							}

							err := docData.File.Sync()
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("Flush error: %v", err)}
							}

							return environment.NOTHIN, nil
						},
					},
					// CLOSE method (implements both READER and WRITER interfaces)
					"CLOSE": {
						Name: "CLOSE",
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							docData, ok := this.NativeData.(*DocumentData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("CLOSE: invalid context")
							}

							if !docData.IsOpen {
								return environment.NOTHIN, nil // Already closed, no error
							}

							err := docData.File.Close()
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("Close error: %v", err)}
							}

							docData.IsOpen = false
							docData.File = nil

							// Update IS_OPEN variable
							if isOpenVar, exists := this.Variables["IS_OPEN"]; exists {
								isOpenVar.Value = environment.NO
							}

							return environment.NOTHIN, nil
						},
					},
					// DELETE method
					"DELETE": {
						Name: "DELETE",
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							docData, ok := this.NativeData.(*DocumentData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("DELETE: invalid context")
							}

							// Close file if it's open before deleting
							if docData.IsOpen {
								err := docData.File.Close()
								if err != nil {
									return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("DELETE: error closing file before deletion: %v", err)}
								}
								docData.IsOpen = false
								docData.File = nil

								// Update IS_OPEN variable
								if isOpenVar, exists := this.Variables["IS_OPEN"]; exists {
									isOpenVar.Value = environment.NO
								}
							}

							// Delete the file
							err := os.Remove(docData.FilePath)
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("DELETE: failed to delete file: %v", err)}
							}

							return environment.NOTHIN, nil
						},
					},
				},
				PublicVariables: map[string]*environment.Variable{
					"PATH": {
						Name:     "PATH",
						Type:     "STRIN",
						Value:    environment.StringValue(""),
						IsLocked: true,
						IsPublic: true,
					},
					"MODE": {
						Name:     "MODE",
						Type:     "STRIN",
						Value:    environment.StringValue(""),
						IsLocked: true,
						IsPublic: true,
					},
					"IS_OPEN": {
						Name:     "IS_OPEN",
						Type:     "BOOL",
						Value:    environment.NO,
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
	return fileClasses
}

// RegisterFILEInEnv registers FILE classes in the given environment
// declarations: empty slice means import all, otherwise import only specified classes
func RegisterFILEInEnv(env *environment.Environment, declarations ...string) error {
	// First ensure IO classes are available since DOCUMENT inherits from IO.READWRITER
	err := RegisterIOInEnv(env, "READWRITER")
	if err != nil {
		return fmt.Errorf("failed to register IO classes: %v", err)
	}

	fileClasses := getFileClasses()

	// If declarations is empty, import all classes
	if len(declarations) == 0 {
		for _, class := range fileClasses {
			env.DefineClass(class)
		}
		return nil
	}

	// Otherwise, import only specified classes
	for _, decl := range declarations {
		declUpper := strings.ToUpper(decl)
		if class, exists := fileClasses[declUpper]; exists {
			env.DefineClass(class)
		} else {
			return fmt.Errorf("unknown FILE class: %s", decl)
		}
	}

	return nil
}
