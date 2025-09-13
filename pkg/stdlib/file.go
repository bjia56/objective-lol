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
				Name: "DOCUMENT",
				Documentation: []string{
					"A file on the file system. Provides methods for file I/O operations.",
					"Supports multiple access modes: R (read-only), W (write-only), RW (read-write), A (append).",
				},
				QualifiedName: "stdlib:FILE.DOCUMENT",
				ModulePath:    "stdlib:FILE",
				ParentClasses: []string{"stdlib:IO.READWRITER"},
				MRO:           []string{"stdlib:FILE.DOCUMENT", "stdlib:IO.READWRITER", "stdlib:IO.READER", "stdlib:IO.WRITER"},
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"DOCUMENT": {
						Name: "DOCUMENT",
						Documentation: []string{
							"Initializes a DOCUMENT instance with a file path.",
						},
						Parameters: []environment.Parameter{
							{Name: "path", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							path := args[0]

							pathVal, ok := path.(environment.StringValue)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("DOCUMENT constructor expects STRIN path, got %s", path.Type())}
							}

							// Initialize the document data
							docData := &DocumentData{
								FilePath: string(pathVal),
								FileMode: "",
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
						Documentation: []string{
							"Opens the file for I/O operations according to the specified mode.",
							"Creates the file if it doesn't exist (for write/append modes).",
							"Sets IS_OPEN to YEZ.",
							"",
							"Modes: R (read-only), W (write-only, overwrites), RW (read-write, creates if needed), A (append).",
						},
						Parameters: []environment.Parameter{
							{Name: "mode", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							mode := args[0]
							modeVal, ok := mode.(environment.StringValue)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("OPEN expects STRIN mode, got %s", mode.Type())}
							}

							docData, ok := this.NativeData.(*DocumentData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "OPEN: invalid context"}
							}

							if docData.IsOpen {
								return environment.NOTHIN, runtime.Exception{Message: "OPEN: file is already open"}
							}

							docData.FileMode = strings.ToUpper(string(modeVal))
							if docData.FileMode != "R" && docData.FileMode != "W" && docData.FileMode != "RW" && docData.FileMode != "A" {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("OPEN: invalid file mode %s", modeVal)}
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
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("OPEN: %v", err)}
							}

							docData.File = file
							docData.IsOpen = true

							return environment.NOTHIN, nil
						},
					},
					// READ method (implements READER interface)
					"READ": {
						Name: "READ",
						Documentation: []string{
							"Reads up to the specified number of characters from the file.",
							"Returns the data read (may be shorter than requested at end of file).",
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

							docData, ok := this.NativeData.(*DocumentData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("READ: invalid context")}
							}

							if !docData.IsOpen {
								return environment.NOTHIN, runtime.Exception{Message: "READ: file is not open"}
							}

							if docData.FileMode != "R" && docData.FileMode != "RW" {
								return environment.NOTHIN, runtime.Exception{Message: "READ: file is not open for reading"}
							}

							readSize := int(sizeVal)
							if readSize <= 0 {
								return environment.StringValue(""), nil
							}

							buffer := make([]byte, readSize)
							n, err := docData.File.Read(buffer)
							if err != nil && err.Error() != "EOF" {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("READ: %v", err)}
							}

							return environment.StringValue(string(buffer[:n])), nil
						},
					},
					// WRITE method (implements WRITER interface)
					"WRITE": {
						Name: "WRITE",
						Documentation: []string{
							"Writes string data to the file.",
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
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("WRITE expects STRIN data, got %s", data.Type())}
							}

							docData, ok := this.NativeData.(*DocumentData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "WRITE: invalid context"}
							}

							if !docData.IsOpen {
								return environment.NOTHIN, runtime.Exception{Message: "WRITE: file is not open"}
							}

							if docData.FileMode != "W" && docData.FileMode != "RW" && docData.FileMode != "A" {
								return environment.NOTHIN, runtime.Exception{Message: "WRITE: file is not open for writing"}
							}

							dataStr := string(dataVal)
							n, err := docData.File.WriteString(dataStr)
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("WRITE: %v", err)}
							}

							return environment.IntegerValue(n), nil
						},
					},
					// SEEK method
					"SEEK": {
						Name: "SEEK",
						Documentation: []string{
							"Sets the file position for next read/write operation.",
							"Position is byte offset from start of file (0-based).",
						},
						Parameters: []environment.Parameter{
							{Name: "position", Type: "INTEGR"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							position := args[0]

							posVal, ok := position.(environment.IntegerValue)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("SEEK expects INTEGR position, got %s", position.Type())}
							}

							docData, ok := this.NativeData.(*DocumentData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "SEEK: invalid context"}
							}

							if !docData.IsOpen {
								return environment.NOTHIN, runtime.Exception{Message: "SEEK: file is not open"}
							}

							_, err := docData.File.Seek(int64(posVal), 0)
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("SEEK: %v", err)}
							}

							return environment.NOTHIN, nil
						},
					},
					// TELL method
					"TELL": {
						Name: "TELL",
						Documentation: []string{
							"Gets the current file position.",
							"Returns current byte position in file.",
						},
						ReturnType: "INTEGR",
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							docData, ok := this.NativeData.(*DocumentData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("TELL: invalid context")
							}

							if !docData.IsOpen {
								return environment.NOTHIN, runtime.Exception{Message: "TELL: file is not open"}
							}

							pos, err := docData.File.Seek(0, 1) // Get current position
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("TELL: %v", err)}
							}

							return environment.IntegerValue(int(pos)), nil
						},
					},
					// EXISTS method
					"EXISTS": {
						Name: "EXISTS",
						Documentation: []string{
							"Checks if the file exists on disk.",
							"Returns YEZ if file exists, NO otherwise.",
						},
						ReturnType: "BOOL",
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							docData, ok := this.NativeData.(*DocumentData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "EXISTS: invalid context"}
							}

							_, err := os.Stat(docData.FilePath)
							exists := err == nil
							return environment.BoolValue(exists), nil
						},
					},
					// FLUSH method
					"FLUSH": {
						Name: "FLUSH",
						Documentation: []string{
							"Flushes the file's contents to disk.",
							"Ensures all buffered data is written.",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							docData, ok := this.NativeData.(*DocumentData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "FLUSH: invalid context"}
							}

							if !docData.IsOpen {
								return environment.NOTHIN, runtime.Exception{Message: "FLUSH: file is not open"}
							}

							err := docData.File.Sync()
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("FLUSH: %v", err)}
							}

							return environment.NOTHIN, nil
						},
					},
					// CLOSE method (implements both READER and WRITER interfaces)
					"CLOSE": {
						Name: "CLOSE",
						Documentation: []string{
							"Closes the file.",
							"Releases any resources associated with the file.",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							docData, ok := this.NativeData.(*DocumentData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "CLOSE: invalid context"}
							}

							if !docData.IsOpen {
								return environment.NOTHIN, nil // Already closed, no error
							}

							err := docData.File.Close()
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("CLOSE: %v", err)}
							}

							docData.IsOpen = false
							docData.File = nil
							docData.FileMode = ""

							return environment.NOTHIN, nil
						},
					},
					// DELETE method
					"DELETE": {
						Name: "DELETE",
						Documentation: []string{
							"Deletes the file from disk.",
							"Automatically closes file if open and sets IS_OPEN to NO.",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							docData, ok := this.NativeData.(*DocumentData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "DELETE: invalid context"}
							}

							// Close file if it's open before deleting
							if docData.IsOpen {
								err := docData.File.Close()
								if err != nil {
									return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("DELETE: error closing file before deletion: %v", err)}
								}
								docData.IsOpen = false
								docData.File = nil
								docData.FileMode = ""
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
							Documentation: []string{
								"Read-only property containing the file path.",
							},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							if data, ok := this.NativeData.(*DocumentData); ok {
								return environment.StringValue(data.FilePath), nil
							}
							return environment.StringValue(""), runtime.Exception{Message: "PATH: invalid context"}
						},
					},
					"MODE": {
						Variable: environment.Variable{
							Name:     "MODE",
							Type:     "STRIN",
							IsLocked: true,
							IsPublic: true,
							Documentation: []string{
								"Read-only property containing the access mode (R, W, RW, or A).",
							},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							if data, ok := this.NativeData.(*DocumentData); ok {
								return environment.StringValue(data.FileMode), nil
							}
							return environment.StringValue(""), runtime.Exception{Message: "MODE: invalid context"}
						},
					},
					"IS_OPEN": {
						Variable: environment.Variable{
							Name:     "IS_OPEN",
							Type:     "BOOL",
							IsLocked: true,
							IsPublic: true,
							Documentation: []string{
								"Read-only property that returns YEZ if file is currently open, NO otherwise.",
							},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							if data, ok := this.NativeData.(*DocumentData); ok {
								if data.IsOpen {
									return environment.YEZ, nil
								}
								return environment.NO, nil
							}
							return environment.NO, runtime.Exception{Message: "IS_OPEN: invalid context"}
						},
					},
					"SIZ": {
						Variable: environment.Variable{
							Name:     "SIZ",
							Type:     "INTEGR",
							IsLocked: true,
							IsPublic: true,
							Documentation: []string{
								"Read-only property that returns the file size in bytes.",
							},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							if data, ok := this.NativeData.(*DocumentData); ok {
								stat, err := os.Stat(data.FilePath)
								if err != nil {
									return environment.IntegerValue(0), runtime.Exception{Message: fmt.Sprintf("SIZ: failed to get file size: %v", err)}
								}
								return environment.IntegerValue(int(stat.Size())), nil
							}
							return environment.IntegerValue(0), runtime.Exception{Message: "SIZ: invalid context"}
						},
					},
					"RWX": {
						Variable: environment.Variable{
							Name:     "RWX",
							Type:     "INTEGR",
							IsLocked: false,
							IsPublic: true,
							Documentation: []string{
								"File permissions (read/write/execute bits).",
								"Can be read to get current permissions or written to change them.",
							},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							if data, ok := this.NativeData.(*DocumentData); ok {
								stat, err := os.Stat(data.FilePath)
								if err != nil {
									return environment.IntegerValue(0), runtime.Exception{Message: fmt.Sprintf("RWX: failed to get file permissions: %v", err)}
								}
								return environment.IntegerValue(int(stat.Mode().Perm())), nil
							}
							return environment.IntegerValue(0), runtime.Exception{Message: "RWX: invalid context"}
						},
						NativeSet: func(this *environment.ObjectInstance, val environment.Value) error {
							intValue, err := val.Cast("INTEGR")
							if err != nil {
								return runtime.Exception{Message: fmt.Sprintf("RWX: value must be castable to INTEGR, got %s", val.Type())}
							}
							intVal := intValue.(environment.IntegerValue)
							if data, ok := this.NativeData.(*DocumentData); ok {
								err := os.Chmod(data.FilePath, os.FileMode(int(intVal)))
								if err != nil {
									return runtime.Exception{Message: fmt.Sprintf("RWX: failed to set file permissions: %v", err)}
								}
								return nil
							}
							return fmt.Errorf("RWX: invalid context")
						},
					},
				},
				PrivateVariables: make(map[string]*environment.MemberVariable),
				PrivateFunctions: make(map[string]*environment.Function),
				SharedVariables:  make(map[string]*environment.MemberVariable),
				SharedFunctions:  make(map[string]*environment.Function),
			},
			"CABINET": {
				Name: "CABINET",
				Documentation: []string{
					"A directory on the filesystem. Provides directory operations for working with directories and their contents.",
				},
				QualifiedName: "stdlib:FILE.CABINET",
				ModulePath:    "stdlib:FILE",
				ParentClasses: []string{},
				MRO:           []string{"stdlib:FILE.CABINET"},
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"CABINET": {
						Name: "CABINET",
						Documentation: []string{
							"Initializes a CABINET instance with the specified path.",
						},
						Parameters: []environment.Parameter{
							{Name: "path", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							path := args[0]

							pathVal, ok := path.(environment.StringValue)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("CABINET constructor expects STRIN path, got %s", path.Type())}
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
						Name: "EXISTS",
						Documentation: []string{
							"Checks if the directory exists, returns YEZ if directory exists, NO otherwise.",
						},
						ReturnType: "BOOL",
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							cabinetData, ok := this.NativeData.(*CabinetData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "EXISTS: invalid context"}
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
						Name: "LIST",
						Documentation: []string{
							"Returns all files and subdirectories in the directory as a BUKKIT.",
						},
						ReturnType: "BUKKIT",
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							cabinetData, ok := this.NativeData.(*CabinetData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "LIST: invalid context"}
							}

							entries, err := os.ReadDir(cabinetData.DirPath)
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("LIST: %v", err)}
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
						Documentation: []string{
							"Creates the directory (including parent directories if needed).",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							cabinetData, ok := this.NativeData.(*CabinetData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "CREATE: invalid context"}
							}

							err := os.MkdirAll(cabinetData.DirPath, 0755)
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("CREATE: %v", err)}
							}

							return environment.NOTHIN, nil
						},
					},
					// DELETE method
					"DELETE": {
						Name: "DELETE",
						Documentation: []string{
							"Deletes an empty directory, throws exception if directory is not empty.",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							cabinetData, ok := this.NativeData.(*CabinetData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "DELETE: invalid context"}
							}

							// Remove empty directory only
							err := os.Remove(cabinetData.DirPath)
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("DELETE: %v", err)}
							}

							return environment.NOTHIN, nil
						},
					},
					// DELETE_ALL method
					"DELETE_ALL": {
						Name: "DELETE_ALL",
						Documentation: []string{
							"Removes directory and all its contents recursively.",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							cabinetData, ok := this.NativeData.(*CabinetData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "DELETE_ALL: invalid context"}
							}

							// Remove directory and all its contents
							err := os.RemoveAll(cabinetData.DirPath)
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("DELETE_ALL: %v", err)}
							}

							return environment.NOTHIN, nil
						},
					},
					// FIND method
					"FIND": {
						Name: "FIND",
						Documentation: []string{
							"Searches for files matching a glob pattern, returns BUKKIT of matching filenames.",
						},
						ReturnType: "BUKKIT",
						Parameters: []environment.Parameter{
							{Name: "pattern", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							pattern := args[0]

							patternVal, ok := pattern.(environment.StringValue)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("FIND expects STRIN pattern, got %s", pattern.Type())}
							}

							cabinetData, ok := this.NativeData.(*CabinetData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "FIND: invalid context"}
							}

							entries, err := os.ReadDir(cabinetData.DirPath)
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("FIND: failed to read directory: %v", err)}
							}

							var matches []environment.Value
							patternStr := string(patternVal)

							for _, entry := range entries {
								matched, err := filepath.Match(patternStr, entry.Name())
								if err != nil {
									return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("FIND: invalid pattern: %v", err)}
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
							Documentation: []string{
								"Read-only property containing the directory path.",
							},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							if data, ok := this.NativeData.(*CabinetData); ok {
								return environment.StringValue(data.DirPath), nil
							}
							return environment.StringValue(""), runtime.Exception{Message: "PATH: invalid context"}
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
				Documentation: []string{
					"The platform-specific path separator character (/ on Unix, \\ on Windows).",
				},
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
			env.DefineVariable(variable.Name, variable.Type, variable.Value, variable.IsLocked, variable.Documentation)
		}
		return nil
	}

	// Otherwise, import only specified classes and variables
	for _, decl := range declarations {
		declUpper := strings.ToUpper(decl)
		if class, exists := fileClasses[declUpper]; exists {
			env.DefineClass(class)
		} else if variable, exists := fileVariables[declUpper]; exists {
			env.DefineVariable(variable.Name, variable.Type, variable.Value, variable.IsLocked, variable.Documentation)
		} else {
			return runtime.Exception{Message: fmt.Sprintf("unknown FILE declaration: %s", decl)}
		}
	}

	return nil
}
