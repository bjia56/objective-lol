package stdlib

import (
	"fmt"
	"os"
	"path/filepath"
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

// CabinetData stores the internal state of a CABINET
type CabinetData struct {
	DirPath string
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
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
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

							return environment.NOTHIN, nil
						},
					},
					// OPEN method
					"OPEN": {
						Name: "OPEN",
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
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
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
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
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
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
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
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
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
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
					// EXISTS method
					"EXISTS": {
						Name:       "EXISTS",
						ReturnType: "BOOL",
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
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
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
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
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
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

							return environment.NOTHIN, nil
						},
					},
					// DELETE method
					"DELETE": {
						Name: "DELETE",
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
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
				PublicVariables: map[string]*environment.MemberVariable{
					"PATH": {
						Variable: environment.Variable{
							Name:     "PATH",
							Type:     "STRIN",
							IsLocked: true,
							IsPublic: true,
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							if data, ok := this.NativeData.(*DocumentData); ok {
								return environment.StringValue(data.FilePath), nil
							}
							return environment.StringValue(""), fmt.Errorf("invalid context for PATH")
						},
					},
					"MODE": {
						Variable: environment.Variable{
							Name:     "MODE",
							Type:     "STRIN",
							IsLocked: true,
							IsPublic: true,
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							if data, ok := this.NativeData.(*DocumentData); ok {
								return environment.StringValue(data.FileMode), nil
							}
							return environment.StringValue(""), fmt.Errorf("invalid context for MODE")
						},
					},
					"IS_OPEN": {
						Variable: environment.Variable{
							Name:     "IS_OPEN",
							Type:     "BOOL",
							IsLocked: true,
							IsPublic: true,
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							if data, ok := this.NativeData.(*DocumentData); ok {
								if data.IsOpen {
									return environment.YEZ, nil
								}
								return environment.NO, nil
							}
							return environment.NO, fmt.Errorf("invalid context for IS_OPEN")
						},
					},
					"SIZ": {
						Variable: environment.Variable{
							Name:     "SIZ",
							Type:     "INTEGR",
							IsLocked: true,
							IsPublic: true,
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							if data, ok := this.NativeData.(*DocumentData); ok {
								stat, err := os.Stat(data.FilePath)
								if err != nil {
									return environment.IntegerValue(0), fmt.Errorf("failed to get file size: %v", err)
								}
								return environment.IntegerValue(int(stat.Size())), nil
							}
							return environment.IntegerValue(0), fmt.Errorf("invalid context for SIZ")
						},
					},
					"RWX": {
						Variable: environment.Variable{
							Name:     "RWX",
							Type:     "INTEGR",
							IsLocked: false,
							IsPublic: true,
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							if data, ok := this.NativeData.(*DocumentData); ok {
								stat, err := os.Stat(data.FilePath)
								if err != nil {
									return environment.IntegerValue(0), fmt.Errorf("failed to get file permissions: %v", err)
								}
								return environment.IntegerValue(int(stat.Mode().Perm())), nil
							}
							return environment.IntegerValue(0), fmt.Errorf("invalid context for RWX")
						},
						NativeSet: func(this *environment.ObjectInstance, val environment.Value) error {
							intVal, ok := val.(environment.IntegerValue)
							if !ok {
								return fmt.Errorf("RWX expects INTEGR value, got %s", val.Type())
							}
							if data, ok := this.NativeData.(*DocumentData); ok {
								err := os.Chmod(data.FilePath, os.FileMode(int(intVal)))
								if err != nil {
									return fmt.Errorf("failed to set file permissions: %v", err)
								}
								return nil
							}
							return fmt.Errorf("invalid context for RWX")
						},
					},
				},
				PrivateVariables: make(map[string]*environment.MemberVariable),
				PrivateFunctions: make(map[string]*environment.Function),
				SharedVariables:  make(map[string]*environment.MemberVariable),
				SharedFunctions:  make(map[string]*environment.Function),
			},
			"CABINET": {
				Name:          "CABINET",
				QualifiedName: "stdlib:FILE.CABINET",
				ModulePath:    "stdlib:FILE",
				ParentClasses: []string{},
				MRO:           []string{"stdlib:FILE.CABINET"},
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"CABINET": {
						Name: "CABINET",
						Parameters: []environment.Parameter{
							{Name: "path", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							path := args[0]

							pathVal, ok := path.(environment.StringValue)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("CABINET constructor expects STRIN path, got %s", path.Type())
							}

							// Initialize the cabinet data
							cabinetData := &CabinetData{
								DirPath: string(pathVal),
							}
							this.NativeData = cabinetData

							return environment.NOTHIN, nil
						},
					},
					// EXISTS method
					"EXISTS": {
						Name:       "EXISTS",
						ReturnType: "BOOL",
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							cabinetData, ok := this.NativeData.(*CabinetData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("EXISTS: invalid context")
							}

							stat, err := os.Stat(cabinetData.DirPath)
							if err != nil {
								return environment.NO, nil
							}
							exists := stat.IsDir()
							return environment.BoolValue(exists), nil
						},
					},
					// LIST method
					"LIST": {
						Name:       "LIST",
						ReturnType: "BUKKIT",
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							cabinetData, ok := this.NativeData.(*CabinetData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("LIST: invalid context")
							}

							entries, err := os.ReadDir(cabinetData.DirPath)
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("Failed to read directory: %v", err)}
							}

							// Create BUKKIT array of filenames
							bukkitObj := NewBukkitInstance()
							bukkitSlice := make(BukkitSlice, len(entries))

							for i, entry := range entries {
								bukkitSlice[i] = environment.StringValue(entry.Name())
							}

							bukkitObj.NativeData = bukkitSlice
							return bukkitObj, nil
						},
					},
					// CREATE method
					"CREATE": {
						Name: "CREATE",
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							cabinetData, ok := this.NativeData.(*CabinetData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("CREATE: invalid context")
							}

							err := os.MkdirAll(cabinetData.DirPath, 0755)
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("Failed to create directory: %v", err)}
							}

							return environment.NOTHIN, nil
						},
					},
					// DELETE method
					"DELETE": {
						Name: "DELETE",
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							cabinetData, ok := this.NativeData.(*CabinetData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("DELETE: invalid context")
							}

							// Remove empty directory only
							err := os.Remove(cabinetData.DirPath)
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("Failed to delete directory: %v", err)}
							}

							return environment.NOTHIN, nil
						},
					},
					// DELETE_ALL method
					"DELETE_ALL": {
						Name: "DELETE_ALL",
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							cabinetData, ok := this.NativeData.(*CabinetData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("DELETE_ALL: invalid context")
							}

							// Remove directory and all its contents
							err := os.RemoveAll(cabinetData.DirPath)
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("Failed to delete directory and its contents: %v", err)}
							}

							return environment.NOTHIN, nil
						},
					},
					// FIND method
					"FIND": {
						Name:       "FIND",
						ReturnType: "BUKKIT",
						Parameters: []environment.Parameter{
							{Name: "pattern", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							pattern := args[0]

							patternVal, ok := pattern.(environment.StringValue)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("FIND expects STRIN pattern, got %s", pattern.Type())
							}

							cabinetData, ok := this.NativeData.(*CabinetData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("FIND: invalid context")
							}

							entries, err := os.ReadDir(cabinetData.DirPath)
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("Failed to read directory: %v", err)}
							}

							var matches []environment.Value
							patternStr := string(patternVal)

							for _, entry := range entries {
								matched, err := filepath.Match(patternStr, entry.Name())
								if err != nil {
									return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("Invalid pattern: %v", err)}
								}
								if matched {
									matches = append(matches, environment.StringValue(entry.Name()))
								}
							}

							// Create BUKKIT array of matching filenames
							bukkitObj := NewBukkitInstance()
							bukkitObj.NativeData = BukkitSlice(matches)

							return bukkitObj, nil
						},
					},
				},
				PublicVariables: map[string]*environment.MemberVariable{
					"PATH": {
						Variable: environment.Variable{
							Name:     "PATH",
							Type:     "STRIN",
							IsLocked: true,
							IsPublic: true,
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							if data, ok := this.NativeData.(*CabinetData); ok {
								return environment.StringValue(data.DirPath), nil
							}
							return environment.StringValue(""), fmt.Errorf("invalid context for PATH")
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
	return fileClasses
}

// Global FILE variable definitions - created once and reused
var fileVarsOnce = sync.Once{}
var fileVariables map[string]*environment.Variable

func getFileVariables() map[string]*environment.Variable {
	fileVarsOnce.Do(func() {
		fileVariables = map[string]*environment.Variable{
			"SEP": {
				Name:     "SEP",
				Type:     "STRIN",
				Value:    environment.StringValue(string(filepath.Separator)),
				IsLocked: true,
				IsPublic: true,
			},
		}
	})
	return fileVariables
}

// RegisterFILEInEnv registers FILE classes and variables in the given environment
// declarations: empty slice means import all, otherwise import only specified classes/variables
func RegisterFILEInEnv(env *environment.Environment, declarations ...string) error {
	// First ensure IO classes are available since DOCUMENT inherits from IO.READWRITER
	err := RegisterIOInEnv(env, "READWRITER", "READER", "WRITER")
	if err != nil {
		return fmt.Errorf("failed to register IO classes: %v", err)
	}

	fileClasses := getFileClasses()
	fileVariables := getFileVariables()

	// If declarations is empty, import all classes and variables
	if len(declarations) == 0 {
		for _, class := range fileClasses {
			env.DefineClass(class)
		}
		for _, variable := range fileVariables {
			err := env.DefineVariable(variable.Name, variable.Type, variable.Value, variable.IsLocked, nil)
			if err != nil {
				return fmt.Errorf("failed to define FILE variable %s: %v", variable.Name, err)
			}
		}
		return nil
	}

	// Otherwise, import only specified classes and variables
	for _, decl := range declarations {
		declUpper := strings.ToUpper(decl)
		if class, exists := fileClasses[declUpper]; exists {
			env.DefineClass(class)
		} else if variable, exists := fileVariables[declUpper]; exists {
			err := env.DefineVariable(variable.Name, variable.Type, variable.Value, variable.IsLocked, nil)
			if err != nil {
				return fmt.Errorf("failed to define FILE variable %s: %v", variable.Name, err)
			}
		} else {
			return fmt.Errorf("unknown FILE declaration: %s", decl)
		}
	}

	return nil
}
