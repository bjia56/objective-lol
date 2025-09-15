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

// moduleFileCategories defines the order that categories should be rendered in documentation
var moduleFileCategories = []string{
	"file-creation",
	"file-io",
	"file-properties",
	"file-positioning",
	"file-management",
	"directory-operations",
	"directory-properties",
}

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
					"A file on the file system that provides methods for file I/O operations.",
					"Supports multiple access modes and inherits from IO.READWRITER for compatibility.",
					"",
					"@class DOCUMENT",
					"@inherits IO.READWRITER",
					"@example Create and write to a file",
					"I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT \"data.txt\"",
					"DOC DO OPEN WIT \"W\"",
					"DOC DO WRITE WIT \"Hello, World!\"",
					"DOC DO CLOSE",
					"@example Read from a file",
					"I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT \"input.txt\"",
					"DOC DO OPEN WIT \"R\"",
					"I HAS A VARIABLE CONTENT TEH STRIN ITZ DOC DO READ WIT 1024",
					"DOC DO CLOSE",
					"@example Check file properties",
					"I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT \"test.txt\"",
					"IZ DOC DO EXISTS?",
					"    SAYZ WIT \"File exists!\"",
					"KTHX",
					"SAYZ WIT DOC PATH",
					"SAYZ WIT DOC SIZ",
					"@note Supports modes: R (read), W (write/overwrite), RW (read-write), A (append)",
					"@note File is automatically created for W, RW, and A modes if it doesn't exist",
					"@see CABINET, IO.READWRITER",
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
							"Creates a new file object but does not open the file yet.",
							"",
							"@syntax NEW DOCUMENT WIT <path>",
							"@param {STRIN} path - The file path to associate with this document",
							"@returns {NOTHIN} No return value (constructor)",
							"@example Create document instance",
							"I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT \"myfile.txt\"",
							"BTW File is not opened yet, use OPEN method to access it",
							"@note File path can be relative or absolute",
							"@note Does not validate file existence - use EXISTS method to check",
							"@see OPEN, EXISTS",
							"@category file-creation",
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
							"Creates the file if it doesn't exist for write/append modes and sets IS_OPEN to YEZ.",
							"",
							"@syntax <document> DO OPEN WIT <mode>",
							"@param {STRIN} mode - Access mode: R (read), W (write/overwrite), RW (read-write), A (append)",
							"@returns {NOTHIN} No return value",
							"@example Open file for reading",
							"I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT \"data.txt\"",
							"DOC DO OPEN WIT \"R\"",
							"BTW File is now open for reading",
							"@example Open file for writing (creates if needed)",
							"I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT \"output.txt\"",
							"DOC DO OPEN WIT \"W\"",
							"BTW File is created and open for writing, existing content is overwritten",
							"@example Open file for append",
							"DOC DO OPEN WIT \"A\"",
							"BTW File is open for appending, writes go to end of file",
							"@note R mode requires file to exist, throws exception if not found",
							"@note W and RW modes create the file if it doesn't exist and overwrite existing content",
							"@note A mode creates the file if it doesn't exist and preserves existing content",
							"@note File must be closed before opening with a different mode",
							"@see CLOSE, IS_OPEN",
							"@category file-io",
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
							"Returns the data read as a string (may be shorter than requested at end of file).",
							"",
							"@syntax <document> DO READ WIT <size>",
							"@param {INTEGR} size - Maximum number of characters to read",
							"@returns {STRIN} The data read from the file",
							"@example Read entire small file",
							"I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT \"config.txt\"",
							"DOC DO OPEN WIT \"R\"",
							"I HAS A VARIABLE CONTENT TEH STRIN ITZ DOC DO READ WIT 1024",
							"DOC DO CLOSE",
							"SAYZ WIT CONTENT",
							"@example Read file in chunks",
							"DOC DO OPEN WIT \"R\"",
							"I HAS A VARIABLE CHUNK TEH STRIN ITZ DOC DO READ WIT 256",
							"WHILE NO SAEM AS (CHUNK SAEM AS \"\")",
							"    SAYZ WIT CHUNK",
							"    CHUNK ITZ DOC DO READ WIT 256",
							"KTHX",
							"DOC DO CLOSE",
							"@note File must be opened in R or RW mode before reading",
							"@note Returns empty string when end of file is reached",
							"@note Use file position methods (SEEK, TELL) for random access",
							"@see WRITE, SEEK, TELL",
							"@category file-io",
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
								return environment.NOTHIN, runtime.Exception{Message: "READ: invalid context"}
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
							"Writes string data to the file at the current position.",
							"Returns the number of characters actually written.",
							"",
							"@syntax <document> DO WRITE WIT <data>",
							"@param {STRIN} data - The string data to write to the file",
							"@returns {INTEGR} Number of characters written",
							"@example Write text to file",
							"I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT \"output.txt\"",
							"DOC DO OPEN WIT \"W\"",
							"I HAS A VARIABLE BYTES_WRITTEN TEH INTEGR ITZ DOC DO WRITE WIT \"Hello, World!\"",
							"DOC DO CLOSE",
							"SAYZ WIT BYTES_WRITTEN BTW Should print 13",
							"@example Append to existing file",
							"DOC DO OPEN WIT \"A\"",
							"DOC DO WRITE WIT \"\\nNew line added\"",
							"DOC DO CLOSE",
							"@example Write multiple lines",
							"DOC DO OPEN WIT \"W\"",
							"DOC DO WRITE WIT \"Line 1\\n\"",
							"DOC DO WRITE WIT \"Line 2\\n\"",
							"DOC DO WRITE WIT \"Line 3\"",
							"DOC DO CLOSE",
							"@note File must be opened in W, RW, or A mode before writing",
							"@note In A mode, data is always written to end of file",
							"@note Use FLUSH to ensure data is written to disk immediately",
							"@see READ, FLUSH, SEEK",
							"@category file-io",
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
							"Sets the file position for the next read/write operation.",
							"Position is specified as a byte offset from the start of the file (0-based).",
							"",
							"@syntax <document> DO SEEK WIT <position>",
							"@param {INTEGR} position - Byte offset from start of file (0 = beginning)",
							"@returns {NOTHIN} No return value",
							"@example Seek to beginning of file",
							"I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT \"data.txt\"",
							"DOC DO OPEN WIT \"RW\"",
							"DOC DO SEEK WIT 0 BTW Go to start of file",
							"@example Seek to specific position",
							"DOC DO SEEK WIT 100 BTW Go to byte 100",
							"I HAS A VARIABLE DATA TEH STRIN ITZ DOC DO READ WIT 50",
							"@example Random access file operations",
							"DOC DO SEEK WIT 0",
							"DOC DO WRITE WIT \"Header\"",
							"DOC DO SEEK WIT 50",
							"DOC DO WRITE WIT \"Middle\"",
							"DOC DO CLOSE",
							"@note File must be open before seeking",
							"@note Position beyond end of file is allowed but behavior is undefined",
							"@note Use TELL to get current position",
							"@see TELL, READ, WRITE",
							"@category file-positioning",
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
							"Gets the current file position as a byte offset from the start of the file.",
							"Returns the current position where the next read or write will occur.",
							"",
							"@syntax <document> DO TELL",
							"@returns {INTEGR} Current byte position in the file (0-based)",
							"@example Check current position",
							"I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT \"data.txt\"",
							"DOC DO OPEN WIT \"R\"",
							"I HAS A VARIABLE POS TEH INTEGR ITZ DOC DO TELL",
							"SAYZ WIT POS BTW Should be 0 at start",
							"@example Track position while reading",
							"DOC DO READ WIT 10 BTW Read 10 characters",
							"POS ITZ DOC DO TELL",
							"SAYZ WIT POS BTW Should be 10",
							"@example Save and restore position",
							"I HAS A VARIABLE SAVED_POS TEH INTEGR ITZ DOC DO TELL",
							"DOC DO READ WIT 100 BTW Read ahead",
							"DOC DO SEEK WIT SAVED_POS BTW Return to saved position",
							"DOC DO CLOSE",
							"@note File must be open before getting position",
							"@note Position is measured in bytes, not characters",
							"@note Useful for implementing random access patterns",
							"@see SEEK, READ, WRITE",
							"@category file-positioning",
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
							"",
							"@syntax <document> DO EXISTS",
							"@returns {BOOL} YEZ if file exists, NO if it doesn't",
							"@example Check file existence before reading",
							"I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT \"config.txt\"",
							"IZ DOC DO EXISTS?",
							"    DOC DO OPEN WIT \"R\"",
							"    BTW Safe to read the file",
							"NOPE",
							"    SAYZ WIT \"File not found!\"",
							"KTHX",
							"@example Conditional file creation",
							"IZ NO SAEM AS (DOC DO EXISTS)?",
							"    DOC DO OPEN WIT \"W\"",
							"    DOC DO WRITE WIT \"Default content\"",
							"    DOC DO CLOSE",
							"    SAYZ WIT \"Created new file with defaults\"",
							"KTHX",
							"@note Does not require file to be open",
							"@note Works with any file path, not just open files",
							"@note Use before opening files in read mode to avoid exceptions",
							"@see OPEN, DELETE",
							"@category file-properties",
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
							"Flushes the file's buffered contents to disk.",
							"Ensures all pending writes are physically written to storage.",
							"",
							"@syntax <document> DO FLUSH",
							"@returns {NOTHIN} No return value",
							"@example Ensure data is saved",
							"I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT \"important.txt\"",
							"DOC DO OPEN WIT \"W\"",
							"DOC DO WRITE WIT \"Critical data\"",
							"DOC DO FLUSH BTW Force write to disk immediately",
							"SAYZ WIT \"Data is now safely on disk\"",
							"DOC DO CLOSE",
							"@example Frequent flush for logging",
							"DOC DO OPEN WIT \"A\"",
							"DOC DO WRITE WIT \"Log entry 1\\n\"",
							"DOC DO FLUSH BTW Ensure log is written",
							"DOC DO WRITE WIT \"Log entry 2\\n\"",
							"DOC DO FLUSH BTW Flush again",
							"DOC DO CLOSE",
							"@note File must be open before flushing",
							"@note Useful for ensuring data persistence in case of crashes",
							"@note May impact performance if called frequently",
							"@note CLOSE automatically flushes, so explicit FLUSH is optional before closing",
							"@see WRITE, CLOSE",
							"@category file-io",
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
							"Closes the file and releases any resources associated with it.",
							"Automatically flushes any buffered data and sets IS_OPEN to NO.",
							"",
							"@syntax <document> DO CLOSE",
							"@returns {NOTHIN} No return value",
							"@example Basic file operations",
							"I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT \"data.txt\"",
							"DOC DO OPEN WIT \"W\"",
							"DOC DO WRITE WIT \"Some content\"",
							"DOC DO CLOSE BTW File is now closed and data is saved",
							"@example Always close files in exception handling",
							"MAYB",
							"    DOC DO OPEN WIT \"R\"",
							"    I HAS A VARIABLE DATA TEH STRIN ITZ DOC DO READ WIT 1024",
							"    BTW Process data here",
							"OOPSIE ERR",
							"    SAYZ WIT ERR",
							"ALWAYZ",
							"    IZ DOC IS_OPEN?",
							"        DOC DO CLOSE BTW Ensure file is always closed",
							"    KTHX",
							"KTHX",
							"@note Safe to call multiple times - no error if already closed",
							"@note Automatically flushes buffered data before closing",
							"@note File cannot be used for I/O operations after closing",
							"@note Always close files to prevent resource leaks",
							"@see OPEN, FLUSH, IS_OPEN",
							"@category file-io",
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
							"Deletes the file from disk permanently.",
							"Automatically closes the file if open and sets IS_OPEN to NO.",
							"",
							"@syntax <document> DO DELETE",
							"@returns {NOTHIN} No return value",
							"@example Delete a file safely",
							"I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT \"temp.txt\"",
							"IZ DOC DO EXISTS?",
							"    DOC DO DELETE",
							"    SAYZ WIT \"File deleted successfully\"",
							"NOPE",
							"    SAYZ WIT \"File doesn't exist\"",
							"KTHX",
							"@example Cleanup after processing",
							"DOC DO OPEN WIT \"W\"",
							"DOC DO WRITE WIT \"Temporary data\"",
							"DOC DO CLOSE",
							"BTW Process the file here",
							"DOC DO DELETE BTW Clean up temporary file",
							"@example Delete with error handling",
							"MAYB",
							"    DOC DO DELETE",
							"OOPSIE ERR",
							"    SAYZ WIT \"Failed to delete file: \"",
							"    SAYZ WIT ERR",
							"KTHX",
							"@note File is automatically closed before deletion if it was open",
							"@note Throws exception if file cannot be deleted (permissions, etc.)",
							"@note Operation is irreversible - file cannot be recovered",
							"@note Use EXISTS to check before deletion to avoid exceptions",
							"@see EXISTS, CLOSE",
							"@category file-management",
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
								"",
								"@property {STRIN} PATH - The file path used when creating this document",
								"@example Access file path",
								"I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT \"/home/user/data.txt\"",
								"SAYZ WIT DOC PATH BTW Prints: /home/user/data.txt",
								"@note This is the original path provided to the constructor",
								"@note Path may be relative or absolute depending on how document was created",
								"@see MODE, IS_OPEN",
								"@category file-properties",
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
								"Read-only property containing the current access mode.",
								"",
								"@property {STRIN} MODE - Current file access mode (R, W, RW, or A)",
								"@example Check file mode",
								"I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT \"data.txt\"",
								"DOC DO OPEN WIT \"RW\"",
								"SAYZ WIT DOC MODE BTW Prints: RW",
								"@example Conditional operations based on mode",
								"IZ (DOC MODE) SAEM AS \"R\"?",
								"    SAYZ WIT \"File is read-only\"",
								"NOPE IZ (DOC MODE) SAEM AS \"W\"?",
								"    SAYZ WIT \"File is write-only\"",
								"NOPE",
								"    SAYZ WIT \"File allows both read and write\"",
								"KTHX",
								"@note Returns empty string if file has never been opened",
								"@note Mode is set when OPEN is called and cleared when CLOSE is called",
								"@see OPEN, CLOSE, IS_OPEN",
								"@category file-properties",
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
								"Read-only property that indicates whether the file is currently open for I/O operations.",
								"",
								"@property {BOOL} IS_OPEN - YEZ if file is open, NO if closed",
								"@example Check if file is open before operations",
								"I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT \"data.txt\"",
								"IZ DOC IS_OPEN?",
								"    SAYZ WIT \"File is already open\"",
								"NOPE",
								"    DOC DO OPEN WIT \"R\"",
								"    SAYZ WIT \"File is now open\"",
								"KTHX",
								"@example Safe file operations",
								"IZ DOC IS_OPEN?",
								"    I HAS A VARIABLE CONTENT TEH STRIN ITZ DOC DO READ WIT 100",
								"NOPE",
								"    SAYZ WIT \"Cannot read - file is not open\"",
								"KTHX",
								"@note Returns NO for newly created documents (before OPEN is called)",
								"@note Automatically set to YEZ by OPEN and NO by CLOSE or DELETE",
								"@note Use this property to avoid exceptions from I/O operations on closed files",
								"@see OPEN, CLOSE, DELETE",
								"@category file-properties",
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
								"Read-only property that returns the current file size in bytes.",
								"",
								"@property {INTEGR} SIZ - Current file size in bytes",
								"@example Check file size",
								"I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT \"data.txt\"",
								"IZ DOC DO EXISTS?",
								"    I HAS A VARIABLE FILESIZE TEH INTEGR ITZ DOC SIZ",
								"    SAYZ WIT \"File size: \"",
								"    SAYZ WIT FILESIZE",
								"    SAYZ WIT \" bytes\"",
								"KTHX",
								"@example Size-based file processing",
								"IZ (DOC SIZ) BIGGR DAN 1024?",
								"    SAYZ WIT \"Large file - processing in chunks\"",
								"NOPE",
								"    SAYZ WIT \"Small file - loading entirely\"",
								"KTHX",
								"@note Returns current size on disk, even if file is not open",
								"@note Size may change if file is modified by other processes",
								"@note Throws exception if file doesn't exist",
								"@see EXISTS, PATH",
								"@category file-properties",
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
								"File permissions property that controls read/write/execute access.",
								"Can be read to get current permissions or assigned to change them.",
								"",
								"@property {INTEGR} RWX - File permission bits (Unix-style octal)",
								"@example Check file permissions",
								"I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT \"script.sh\"",
								"I HAS A VARIABLE PERMS TEH INTEGR ITZ DOC RWX",
								"SAYZ WIT \"File permissions: \"",
								"SAYZ WIT PERMS",
								"@example Make file executable",
								"DOC RWX ITZ 755 BTW rwxr-xr-x",
								"SAYZ WIT \"File is now executable\"",
								"@example Set read-only permissions",
								"DOC RWX ITZ 444 BTW r--r--r--",
								"SAYZ WIT \"File is now read-only\"",
								"@example Common permission values",
								"BTW 644 = rw-r--r-- (owner read/write, others read)",
								"BTW 755 = rwxr-xr-x (owner all, others read/execute)",
								"BTW 600 = rw------- (owner read/write only)",
								"@note Uses Unix-style octal permission notation",
								"@note Changes take effect immediately on the filesystem",
								"@note May throw exception if user lacks permission to change file permissions",
								"@see PATH, EXISTS",
								"@category file-properties",
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
					"A directory on the filesystem that provides operations for working with directories and their contents.",
					"Supports creating, listing, searching, and managing directories and files within them.",
					"",
					"@class CABINET",
					"@example Create and work with directory",
					"I HAS A VARIABLE DIR TEH CABINET ITZ NEW CABINET WIT \"my_folder\"",
					"IZ NO SAEM AS (DIR DO EXISTS)?",
					"    DIR DO CREATE",
					"    SAYZ WIT \"Directory created\"",
					"KTHX",
					"@example List directory contents",
					"I HAS A VARIABLE FILES TEH BUKKIT ITZ DIR DO LIST",
					"I HAS A VARIABLE IDX TEH INTEGR ITZ 0",
					"WHILE IDX SMALLR THAN FILES SIZ",
					"    I HAS A VARIABLE FILENAME TEH STRIN ITZ FILES DO AT WIT IDX",
					"    SAYZ WIT FILENAME",
					"    IDX ITZ IDX MOAR 1",
					"KTHX",
					"@example Find specific files",
					"I HAS A VARIABLE TXTFILES TEH BUKKIT ITZ DIR DO FIND WIT \"*.txt\"",
					"SAYZ WIT \"Found \"",
					"SAYZ WIT TXTFILES SIZ",
					"SAYZ WIT \" text files\"",
					"@note Works with both relative and absolute directory paths",
					"@note Directory operations may require appropriate filesystem permissions",
					"@see DOCUMENT, SEP",
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
							"Initializes a CABINET instance with the specified directory path.",
							"Creates a directory object for performing directory operations.",
							"",
							"@syntax NEW CABINET WIT <path>",
							"@param {STRIN} path - The directory path to associate with this cabinet",
							"@returns {NOTHIN} No return value (constructor)",
							"@example Create cabinet for existing directory",
							"I HAS A VARIABLE DIR TEH CABINET ITZ NEW CABINET WIT \"/home/user/documents\"",
							"@example Create cabinet for new directory",
							"I HAS A VARIABLE NEWDIR TEH CABINET ITZ NEW CABINET WIT \"temp_folder\"",
							"BTW Directory doesn't need to exist yet, use CREATE method to make it",
							"@note Directory path can be relative or absolute",
							"@note Does not validate directory existence - use EXISTS method to check",
							"@note Does not create the directory - use CREATE method for that",
							"@see CREATE, EXISTS",
							"@category directory-operations",
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
							"Checks if the directory exists on the filesystem.",
							"Returns YEZ if directory exists and is actually a directory, NO otherwise.",
							"",
							"@syntax <cabinet> DO EXISTS",
							"@returns {BOOL} YEZ if directory exists, NO if it doesn't",
							"@example Check directory before operations",
							"I HAS A VARIABLE DIR TEH CABINET ITZ NEW CABINET WIT \"my_folder\"",
							"IZ DIR DO EXISTS?",
							"    SAYZ WIT \"Directory exists\"",
							"    I HAS A VARIABLE FILES TEH BUKKIT ITZ DIR DO LIST",
							"NOPE",
							"    SAYZ WIT \"Directory doesn't exist\"",
							"KTHX",
							"@example Conditional directory creation",
							"IZ NO SAEM AS (DIR DO EXISTS)?",
							"    DIR DO CREATE",
							"    SAYZ WIT \"Directory created\"",
							"NOPE",
							"    SAYZ WIT \"Directory already exists\"",
							"KTHX",
							"@note Returns NO if path exists but is not a directory (e.g., it's a file)",
							"@note Use this before LIST or other operations to avoid exceptions",
							"@see CREATE, LIST",
							"@category directory-properties",
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
							"Returns all files and subdirectories in the directory as a BUKKIT of strings.",
							"Each entry is just the filename or directory name, not the full path.",
							"",
							"@syntax <cabinet> DO LIST",
							"@returns {BUKKIT} Array of filenames and directory names",
							"@example List all directory contents",
							"I HAS A VARIABLE DIR TEH CABINET ITZ NEW CABINET WIT \"documents\"",
							"I HAS A VARIABLE FILES TEH BUKKIT ITZ DIR DO LIST",
							"SAYZ WIT \"Directory contains \"",
							"SAYZ WIT FILES SIZ",
							"SAYZ WIT \" items:\"",
							"I HAS A VARIABLE IDX TEH INTEGR ITZ 0",
							"WHILE IDX SMALLR THAN FILES SIZ",
							"    I HAS A VARIABLE FILENAME TEH STRIN ITZ FILES DO AT WIT IDX",
							"    SAYZ WIT \"  \"",
							"    SAYZ WIT FILENAME",
							"    IDX ITZ IDX MOAR 1",
							"KTHX",
							"@example Check for specific files",
							"FILES ITZ DIR DO LIST",
							"IDX ITZ 0",
							"WHILE IDX SMALLR THAN FILES SIZ",
							"    I HAS A VARIABLE FILENAME TEH STRIN ITZ FILES DO AT WIT IDX",
							"    IZ FILENAME SAEM AS \"config.txt\"?",
							"        SAYZ WIT \"Found config file!\"",
							"    KTHX",
							"    IDX ITZ IDX MOAR 1",
							"KTHX",
							"@note Returns only names, not full paths - combine with PATH property for full paths",
							"@note Directory must exist before listing - check with EXISTS first",
							"@note Returns empty BUKKIT for empty directories",
							"@note Order of entries is not guaranteed",
							"@see EXISTS, FIND",
							"@category directory-operations",
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
							"Creates the directory on the filesystem, including any necessary parent directories.",
							"Similar to 'mkdir -p' command - creates the entire path if needed.",
							"",
							"@syntax <cabinet> DO CREATE",
							"@returns {NOTHIN} No return value",
							"@example Create directory",
							"I HAS A VARIABLE DIR TEH CABINET ITZ NEW CABINET WIT \"new_folder\"",
							"DIR DO CREATE",
							"SAYZ WIT \"Directory created successfully\"",
							"@example Create nested directory structure",
							"I HAS A VARIABLE NESTED TEH CABINET ITZ NEW CABINET WIT \"parent/child/grandchild\"",
							"NESTED DO CREATE",
							"SAYZ WIT \"Nested directory structure created\"",
							"@example Safe directory creation",
							"IZ NO SAEM AS (DIR DO EXISTS)?",
							"    DIR DO CREATE",
							"    SAYZ WIT \"Directory created\"",
							"NOPE",
							"    SAYZ WIT \"Directory already exists\"",
							"KTHX",
							"@note Creates parent directories automatically if they don't exist",
							"@note No error if directory already exists",
							"@note May throw exception if lacking filesystem permissions",
							"@note Sets default directory permissions (usually 755)",
							"@see EXISTS, DELETE",
							"@category directory-operations",
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
							"Deletes an empty directory from the filesystem.",
							"Throws exception if directory is not empty - use DELETE_ALL for recursive deletion.",
							"",
							"@syntax <cabinet> DO DELETE",
							"@returns {NOTHIN} No return value",
							"@example Delete empty directory",
							"I HAS A VARIABLE DIR TEH CABINET ITZ NEW CABINET WIT \"empty_folder\"",
							"IZ DIR DO EXISTS?",
							"    DIR DO DELETE",
							"    SAYZ WIT \"Empty directory deleted\"",
							"KTHX",
							"@example Safe deletion with error handling",
							"MAYB",
							"    DIR DO DELETE",
							"    SAYZ WIT \"Directory deleted successfully\"",
							"OOPSIE ERR",
							"    SAYZ WIT \"Failed to delete directory: \"",
							"    SAYZ WIT ERR",
							"    SAYZ WIT \"(Directory may not be empty)\"",
							"KTHX",
							"@note Only deletes empty directories - fails if directory contains files or subdirectories",
							"@note Use LIST to check directory contents before deletion",
							"@note Use DELETE_ALL for recursive deletion of non-empty directories",
							"@note No error if directory doesn't exist",
							"@see DELETE_ALL, EXISTS, LIST",
							"@category directory-operations",
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
							"Deletes all files and subdirectories within the directory, then deletes the directory itself.",
							"",
							"@syntax <cabinet> DO DELETE_ALL",
							"@returns {NOTHIN} No return value",
							"@example Delete directory tree",
							"I HAS A VARIABLE TMPDIR TEH CABINET ITZ NEW CABINET WIT \"temp_data\"",
							"IZ TMPDIR DO EXISTS?",
							"    TMPDIR DO DELETE_ALL",
							"    SAYZ WIT \"Directory and all contents deleted\"",
							"KTHX",
							"@example Cleanup with confirmation",
							"I HAS A VARIABLE FILES TEH BUKKIT ITZ TMPDIR DO LIST",
							"SAYZ WIT \"About to delete directory with \"",
							"SAYZ WIT FILES SIZ",
							"SAYZ WIT \" items. Continue? (y/n)\"",
							"I HAS A VARIABLE CONFIRM TEH STRIN ITZ GIMME",
							"IZ CONFIRM SAEM AS \"y\"?",
							"    TMPDIR DO DELETE_ALL",
							"    SAYZ WIT \"Directory tree deleted\"",
							"NOPE",
							"    SAYZ WIT \"Operation cancelled\"",
							"KTHX",
							"@note DANGEROUS OPERATION - permanently deletes all contents",
							"@note Use with extreme caution - there is no undo",
							"@note Equivalent to 'rm -rf' command",
							"@note No error if directory doesn't exist",
							"@see DELETE, LIST, EXISTS",
							"@category directory-operations",
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
							"Searches for files and directories matching a glob pattern.",
							"Returns a BUKKIT containing names of items that match the pattern.",
							"",
							"@syntax <cabinet> DO FIND WIT <pattern>",
							"@param {STRIN} pattern - Glob pattern to match against filenames",
							"@returns {BUKKIT} Array of matching filenames and directory names",
							"@example Find text files",
							"I HAS A VARIABLE DIR TEH CABINET ITZ NEW CABINET WIT \"documents\"",
							"I HAS A VARIABLE TXTFILES TEH BUKKIT ITZ DIR DO FIND WIT \"*.txt\"",
							"I HAS A VARIABLE IDX TEH INTEGR ITZ 0",
							"WHILE IDX SMALLR THAN TXTFILES SIZ",
							"    I HAS A VARIABLE FILENAME TEH STRIN ITZ TXTFILES DO AT WIT IDX",
							"    SAYZ WIT \"Found text file: \"",
							"    SAYZ WIT FILENAME",
							"    IDX ITZ IDX MOAR 1",
							"KTHX",
							"@example Find files with specific prefix",
							"I HAS A VARIABLE LOGFILES TEH BUKKIT ITZ DIR DO FIND WIT \"log_*\"",
							"SAYZ WIT \"Found \"",
							"SAYZ WIT LOGFILES SIZ",
							"SAYZ WIT \" log files\"",
							"@example Find files with specific patterns",
							"I HAS A VARIABLE IMAGES TEH BUKKIT ITZ DIR DO FIND WIT \"*.{jpg,png,gif}\"",
							"BTW Note: Some glob implementations may not support brace expansion",
							"@note Uses standard glob patterns: * (any chars), ? (single char), [] (char class)",
							"@note Pattern matching is case-sensitive on most systems",
							"@note Returns empty BUKKIT if no matches found",
							"@note Directory must exist before searching",
							"@see LIST, EXISTS",
							"@category directory-operations",
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
								"",
								"@property {STRIN} PATH - The directory path used when creating this cabinet",
								"@example Access directory path",
								"I HAS A VARIABLE DIR TEH CABINET ITZ NEW CABINET WIT \"/home/user/documents\"",
								"SAYZ WIT DIR PATH BTW Prints: /home/user/documents",
								"@example Build full file paths",
								"I HAS A VARIABLE FILES TEH BUKKIT ITZ DIR DO LIST",
								"I HAS A VARIABLE IDX TEH INTEGR ITZ 0",
								"WHILE IDX SMALLR THAN FILES SIZ",
								"    I HAS A VARIABLE FILENAME TEH STRIN ITZ FILES DO AT WIT IDX",
								"    I HAS A VARIABLE FULLPATH TEH STRIN ITZ DIR PATH",
								"    FULLPATH ITZ FULLPATH MOAR SEP MOAR FILENAME",
								"    SAYZ WIT \"Full path: \"",
								"    SAYZ WIT FULLPATH",
								"    IDX ITZ IDX MOAR 1",
								"KTHX",
								"@note This is the original path provided to the constructor",
								"@note Path may be relative or absolute depending on how cabinet was created",
								"@note Combine with SEP to build full file paths",
								"@see SEP",
								"@category directory-properties",
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
					"The platform-specific path separator character used to join file and directory paths.",
					"",
					"@variable {STRIN} SEP - Platform path separator (/ on Unix, \\\\ on Windows)",
					"@example Build file paths",
					"I HAS A VARIABLE DIR TEH STRIN ITZ \"documents\"",
					"I HAS A VARIABLE FILENAME TEH STRIN ITZ \"readme.txt\"",
					"I HAS A VARIABLE FULLPATH TEH STRIN ITZ DIR MOAR SEP MOAR FILENAME",
					"SAYZ WIT FULLPATH BTW Prints: documents/readme.txt (Unix) or documents\\\\readme.txt (Windows)",
					"@example Cross-platform directory creation",
					"I HAS A VARIABLE PATH TEH STRIN ITZ \"home\" MOAR SEP MOAR \"user\" MOAR SEP MOAR \"data\"",
					"I HAS A VARIABLE DIR TEH CABINET ITZ NEW CABINET WIT PATH",
					"DIR DO CREATE",
					"@example Parse directory paths",
					"I HAS A VARIABLE PARTS TEH BUKKIT ITZ FULLPATH DO SPLIT WIT SEP",
					"I HAS A VARIABLE IDX TEH INTEGR ITZ 0",
					"WHILE IDX SMALLR THAN PARTS SIZ",
					"    I HAS A VARIABLE PART TEH STRIN ITZ PARTS DO AT WIT IDX",
					"    SAYZ WIT \"Path component: \"",
					"    SAYZ WIT PART",
					"    IDX ITZ IDX MOAR 1",
					"KTHX",
					"@note Value is / on Unix-like systems and \\\\ on Windows",
					"@note Use this instead of hardcoded separators for cross-platform compatibility",
					"@note Combine with string concatenation to build file paths",
					"@see DOCUMENT, CABINET",
					"@category file-properties",
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
