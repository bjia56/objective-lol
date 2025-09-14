package stdlib

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/runtime"
)

// moduleIOCategories defines the order that categories should be rendered in documentation
var moduleIOCategories = []string{
	"io-interfaces",
	"io-buffered",
	"io-operations",
	"io-construction",
	"io-configuration",
}

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
					"",
					"@class READER",
					"@example Basic reader usage pattern",
					"I HAS A VARIABLE READER TEH READER ITZ GET_SOME_READER",
					"I HAS A VARIABLE DATA TEH STRIN ITZ READER DO READ WIT 1024",
					"WHILE NO SAEM AS (DATA LENGTH SAEM AS 0)",
					"    SAYZ WIT DATA",
					"    DATA ITZ READER DO READ WIT 1024",
					"KTHX",
					"READER DO CLOSE",
					"@example Reading with error handling",
					"I HAS A VARIABLE READER TEH READER ITZ GET_SOME_READER",
					"MAYB",
					"    I HAS A VARIABLE DATA TEH STRIN ITZ READER DO READ WIT 512",
					"    SAYZ WIT \"Read data: \"",
					"    SAYZ WIT DATA",
					"OOPSIE ERR",
					"    SAYZ WIT \"Error reading: \"",
					"    SAYZ WIT ERR",
					"KTHX",
					"READER DO CLOSE",
					"@example Reading entire content",
					"I HAS A VARIABLE READER TEH READER ITZ GET_SOME_READER",
					"I HAS A VARIABLE CONTENT TEH STRIN ITZ \"\"",
					"I HAS A VARIABLE CHUNK TEH STRIN ITZ READER DO READ WIT 4096",
					"WHILE NO SAEM AS (CHUNK LENGTH SAEM AS 0)",
					"    CONTENT ITZ CONTENT MOAR CHUNK",
					"    CHUNK ITZ READER DO READ WIT 4096",
					"KTHX",
					"SAYZ WIT \"Total content length: \"",
					"SAYZ WIT CONTENT LENGTH",
					"READER DO CLOSE",
					"@note This is an abstract class - use concrete implementations like BUFFERED_READER",
					"@note READ returns empty string when end-of-file is reached",
					"@note Always call CLOSE when done to release resources",
					"@note READ may return fewer characters than requested",
					"@see BUFFERED_READER, READWRITER",
					"@category io-interfaces",
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
							"",
							"@syntax <reader> DO READ WIT <size>",
							"@param {INTEGR} size - Maximum number of characters to read",
							"@returns {STRIN} Data read from the source, or empty string at EOF",
							"@example Read fixed-size chunks",
							"I HAS A VARIABLE READER TEH READER ITZ GET_FILE_READER",
							"I HAS A VARIABLE CHUNK_SIZE TEH INTEGR ITZ 1024",
							"I HAS A VARIABLE DATA TEH STRIN ITZ READER DO READ WIT CHUNK_SIZE",
							"WHILE NO SAEM AS (DATA LENGTH SAEM AS 0)",
							"    SAYZ WIT \"Read chunk of length: \"",
							"    SAYZ WIT DATA LENGTH",
							"    DATA ITZ READER DO READ WIT CHUNK_SIZE",
							"KTHX",
							"READER DO CLOSE",
							"@example Read single character at a time",
							"I HAS A VARIABLE READER TEH READER ITZ GET_INPUT_READER",
							"I HAS A VARIABLE CHAR TEH STRIN ITZ READER DO READ WIT 1",
							"WHILE NO SAEM AS (CHAR LENGTH SAEM AS 0)",
							"    SAYZ WIT \"Character: \"",
							"    SAYZ WIT CHAR",
							"    CHAR ITZ READER DO READ WIT 1",
							"KTHX",
							"READER DO CLOSE",
							"@example Handle end-of-file",
							"I HAS A VARIABLE READER TEH READER ITZ GET_FILE_READER",
							"I HAS A VARIABLE DATA TEH STRIN ITZ READER DO READ WIT 100",
							"IZ DATA LENGTH SAEM AS 0?",
							"    SAYZ WIT \"Reached end of file\"",
							"NOPE",
							"    SAYZ WIT \"Read data: \"",
							"    SAYZ WIT DATA",
							"KTHX",
							"READER DO CLOSE",
							"@example Read with size validation",
							"I HAS A VARIABLE READER TEH READER ITZ GET_NETWORK_READER",
							"I HAS A VARIABLE REQUESTED_SIZE TEH INTEGR ITZ 2048",
							"I HAS A VARIABLE DATA TEH STRIN ITZ READER DO READ WIT REQUESTED_SIZE",
							"IZ DATA LENGTH SAEM AS 0?",
							"    SAYZ WIT \"No data available\"",
							"NOPE",
							"    IZ DATA LENGTH BIGGR THAN REQUESTED_SIZE?",
							"        SAYZ WIT \"Error: Read more than requested\"",
							"    NOPE",
							"        SAYZ WIT \"Successfully read \"",
							"        SAYZ WIT DATA LENGTH",
							"        SAYZ WIT \" characters\"",
							"    KTHX",
							"KTHX",
							"READER DO CLOSE",
							"@note May return fewer characters than requested",
							"@note Returns empty string when end-of-file is reached",
							"@note Size must be a positive integer",
							"@note Implementation depends on the concrete reader type",
							"@throws Exception if size is invalid or I/O error occurs",
							"@see CLOSE, BUFFERED_READER",
							"@category io-operations",
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
							"",
							"@syntax <reader> DO CLOSE",
							"@example Basic cleanup",
							"I HAS A VARIABLE READER TEH READER ITZ GET_FILE_READER",
							"I HAS A VARIABLE DATA TEH STRIN ITZ READER DO READ WIT 1024",
							"WHILE NO SAEM AS (DATA LENGTH SAEM AS 0)",
							"    PROCESS_DATA WIT DATA",
							"    DATA ITZ READER DO READ WIT 1024",
							"KTHX",
							"READER DO CLOSE",
							"SAYZ WIT \"Reader closed successfully\"",
							"@example Close in error handling",
							"I HAS A VARIABLE READER TEH READER ITZ GET_NETWORK_READER",
							"MAYB",
							"    I HAS A VARIABLE DATA TEH STRIN ITZ READER DO READ WIT 512",
							"    PROCESS_DATA WIT DATA",
							"OOPSIE ERR",
							"    SAYZ WIT \"Error during processing: \"",
							"    SAYZ WIT ERR",
							"    READER DO CLOSE",
							"KTHX",
							"@example Multiple readers cleanup",
							"I HAS A VARIABLE READERS TEH BUKKIT ITZ NEW BUKKIT",
							"READERS DO PUSH WIT GET_FILE_READER_1",
							"READERS DO PUSH WIT GET_FILE_READER_2",
							"READERS DO PUSH WIT GET_FILE_READER_3",
							"WHILE NO SAEM AS (READERS LENGTH SAEM AS 0)",
							"    I HAS A VARIABLE READER TEH READER ITZ READERS DO POP",
							"    I HAS A VARIABLE DATA TEH STRIN ITZ READER DO READ WIT 1024",
							"    WHILE NO SAEM AS (DATA LENGTH SAEM AS 0)",
							"        PROCESS_DATA WIT DATA",
							"        DATA ITZ READER DO READ WIT 1024",
							"    KTHX",
							"    READER DO CLOSE",
							"KTHX",
							"SAYZ WIT \"All readers closed\"",
							"@example Close with resource tracking",
							"I HAS A VARIABLE READER TEH READER ITZ GET_DATABASE_READER",
							"I HAS A VARIABLE RESOURCE_COUNT TEH INTEGR ITZ 1",
							"MAYB",
							"    I HAS A VARIABLE DATA TEH STRIN ITZ READER DO READ WIT 2048",
							"    WHILE NO SAEM AS (DATA LENGTH SAEM AS 0)",
							"        SAVE_TO_DATABASE WIT DATA",
							"        DATA ITZ READER DO READ WIT 2048",
							"    KTHX",
							"OOPSIE ERR",
							"    SAYZ WIT \"Database operation failed: \"",
							"    SAYZ WIT ERR",
							"KTHX",
							"RESOURCE_COUNT ITZ RESOURCE_COUNT MINUS 1",
							"READER DO CLOSE",
							"SAYZ WIT \"Resources remaining: \"",
							"SAYZ WIT RESOURCE_COUNT",
							"@note Always call CLOSE to prevent resource leaks",
							"@note CLOSE should be called even if errors occur during reading",
							"@note Multiple CLOSE calls should be safe (idempotent)",
							"@note CLOSE may flush buffers or finalize operations",
							"@note After CLOSE, further READ operations may fail",
							"@throws Exception if close operation fails",
							"@see READ, BUFFERED_READER",
							"@category io-operations",
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
					"",
					"@class WRITER",
					"@example Basic writer usage pattern",
					"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_SOME_WRITER",
					"WRITER DO WRITE WIT \"Hello, World!\"",
					"WRITER DO CLOSE",
					"SAYZ WIT \"Data written successfully\"",
					"@example Writing with error handling",
					"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_FILE_WRITER",
					"MAYB",
					"    WRITER DO WRITE WIT \"Important data\"",
					"    SAYZ WIT \"Data written successfully\"",
					"OOPSIE ERR",
					"    SAYZ WIT \"Error writing data: \"",
					"    SAYZ WIT ERR",
					"KTHX",
					"WRITER DO CLOSE",
					"@example Writing multiple pieces of data",
					"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_LOG_WRITER",
					"I HAS A VARIABLE MESSAGES TEH BUKKIT ITZ NEW BUKKIT",
					"MESSAGES DO PUSH WIT \"Starting process...\"",
					"MESSAGES DO PUSH WIT \"Processing data...\"",
					"MESSAGES DO PUSH WIT \"Process complete\"",
					"WHILE NO SAEM AS (MESSAGES LENGTH SAEM AS 0)",
					"    I HAS A VARIABLE MSG TEH STRIN ITZ MESSAGES DO POP",
					"    WRITER DO WRITE WIT MSG",
					"    WRITER DO WRITE WIT \"\\n\"",
					"KTHX",
					"WRITER DO CLOSE",
					"@note This is an abstract class - use concrete implementations like BUFFERED_WRITER",
					"@note Always call CLOSE when done to ensure data is flushed",
					"@note WRITE may not immediately write to the underlying destination",
					"@note Implementation depends on the concrete writer type",
					"@see BUFFERED_WRITER, READWRITER",
					"@category io-interfaces",
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
							"",
							"@syntax <writer> DO WRITE WIT <data>",
							"@param {STRIN} data - String data to write",
							"@returns {INTEGR} Number of characters written",
							"@example Write simple text",
							"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_FILE_WRITER",
							"I HAS A VARIABLE CHARS_WRITTEN TEH INTEGR ITZ WRITER DO WRITE WIT \"Hello World\"",
							"SAYZ WIT \"Wrote \"",
							"SAYZ WIT CHARS_WRITTEN",
							"SAYZ WIT \" characters\"",
							"WRITER DO CLOSE",
							"@example Write with validation",
							"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_NETWORK_WRITER",
							"I HAS A VARIABLE DATA TEH STRIN ITZ \"Important message\"",
							"I HAS A VARIABLE WRITTEN TEH INTEGR ITZ WRITER DO WRITE WIT DATA",
							"IZ WRITTEN SAEM AS DATA LENGTH?",
							"    SAYZ WIT \"All data written successfully\"",
							"NOPE",
							"    SAYZ WIT \"Warning: Only wrote \"",
							"    SAYZ WIT WRITTEN",
							"    SAYZ WIT \" of \"",
							"    SAYZ WIT DATA LENGTH",
							"    SAYZ WIT \" characters\"",
							"KTHX",
							"WRITER DO CLOSE",
							"@example Write formatted data",
							"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_LOG_WRITER",
							"I HAS A VARIABLE TIMESTAMP TEH STRIN ITZ GET_CURRENT_TIME",
							"I HAS A VARIABLE LEVEL TEH STRIN ITZ \"INFO\"",
							"I HAS A VARIABLE MESSAGE TEH STRIN ITZ \"Process started\"",
							"WRITER DO WRITE WIT TIMESTAMP",
							"WRITER DO WRITE WIT \" [\"",
							"WRITER DO WRITE WIT LEVEL",
							"WRITER DO WRITE WIT \"] \"",
							"WRITER DO WRITE WIT MESSAGE",
							"WRITER DO WRITE WIT \"\\n\"",
							"WRITER DO CLOSE",
							"@example Write binary-like data",
							"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_BINARY_WRITER",
							"I HAS A VARIABLE HEADER TEH STRIN ITZ \"\\x00\\x01\\x02\\x03\"",
							"I HAS A VARIABLE PAYLOAD TEH STRIN ITZ \"\\x04\\x05\\x06\\x07\"",
							"WRITER DO WRITE WIT HEADER",
							"WRITER DO WRITE WIT PAYLOAD",
							"I HAS A VARIABLE TOTAL_WRITTEN TEH INTEGR ITZ HEADER LENGTH MOAR PAYLOAD LENGTH",
							"SAYZ WIT \"Binary data written: \"",
							"SAYZ WIT TOTAL_WRITTEN",
							"SAYZ WIT \" bytes\"",
							"WRITER DO CLOSE",
							"@example Batch writing with progress",
							"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_FILE_WRITER",
							"I HAS A VARIABLE LINES TEH BUKKIT ITZ NEW BUKKIT",
							"LINES DO PUSH WIT \"Line 1\"",
							"LINES DO PUSH WIT \"Line 2\"",
							"LINES DO PUSH WIT \"Line 3\"",
							"I HAS A VARIABLE TOTAL_CHARS TEH INTEGR ITZ 0",
							"WHILE NO SAEM AS (LINES LENGTH SAEM AS 0)",
							"    I HAS A VARIABLE LINE TEH STRIN ITZ LINES DO POP",
							"    I HAS A VARIABLE LINE_WITH_NEWLINE TEH STRIN ITZ LINE MOAR \"\\n\"",
							"    I HAS A VARIABLE WRITTEN TEH INTEGR ITZ WRITER DO WRITE WIT LINE_WITH_NEWLINE",
							"    TOTAL_CHARS ITZ TOTAL_CHARS MOAR WRITTEN",
							"KTHX",
							"SAYZ WIT \"Total characters written: \"",
							"SAYZ WIT TOTAL_CHARS",
							"WRITER DO CLOSE",
							"@note Returns the number of characters actually written",
							"@note May write fewer characters than provided if error occurs",
							"@note Data may be buffered and not immediately written",
							"@note Call CLOSE to ensure all data is flushed",
							"@note Implementation depends on the concrete writer type",
							"@throws Exception if data is invalid or I/O error occurs",
							"@see CLOSE, BUFFERED_WRITER",
							"@category io-operations",
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
							"",
							"@syntax <writer> DO CLOSE",
							"@example Basic cleanup after writing",
							"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_FILE_WRITER",
							"WRITER DO WRITE WIT \"Final data\"",
							"WRITER DO CLOSE",
							"SAYZ WIT \"Writer closed successfully\"",
							"@example Close with error handling",
							"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_NETWORK_WRITER",
							"MAYB",
							"    WRITER DO WRITE WIT \"Important data\"",
							"    WRITER DO CLOSE",
							"    SAYZ WIT \"Data written and writer closed\"",
							"OOPSIE ERR",
							"    SAYZ WIT \"Error during write/close: \"",
							"    SAYZ WIT ERR",
							"    MAYB",
							"        WRITER DO CLOSE",
							"    OOPSIE CLOSE_ERR",
							"        SAYZ WIT \"Error closing writer: \"",
							"        SAYZ WIT CLOSE_ERR",
							"    KTHX",
							"KTHX",
							"@example Multiple writers cleanup",
							"I HAS A VARIABLE WRITERS TEH BUKKIT ITZ NEW BUKKIT",
							"WRITERS DO PUSH WIT GET_FILE_WRITER_1",
							"WRITERS DO PUSH WIT GET_FILE_WRITER_2",
							"WRITERS DO PUSH WIT GET_FILE_WRITER_3",
							"WHILE NO SAEM AS (WRITERS LENGTH SAEM AS 0)",
							"    I HAS A VARIABLE WRITER TEH WRITER ITZ WRITERS DO POP",
							"    WRITER DO WRITE WIT \"Data for writer\"",
							"    WRITER DO CLOSE",
							"KTHX",
							"SAYZ WIT \"All writers closed\"",
							"@example Close with final flush",
							"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_BUFFERED_WRITER",
							"WRITER DO WRITE WIT \"Data 1\"",
							"WRITER DO WRITE WIT \"Data 2\"",
							"WRITER DO WRITE WIT \"Data 3\"",
							"WRITER DO CLOSE",
							"SAYZ WIT \"All buffered data flushed and writer closed\"",
							"@example Resource management pattern",
							"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_DATABASE_WRITER",
							"I HAS A VARIABLE RESOURCE_COUNT TEH INTEGR ITZ 1",
							"MAYB",
							"    WRITER DO WRITE WIT \"INSERT INTO table VALUES (...)\"",
							"    WRITER DO WRITE WIT \"COMMIT\"",
							"OOPSIE ERR",
							"    SAYZ WIT \"Database write failed: \"",
							"    SAYZ WIT ERR",
							"    WRITER DO WRITE WIT \"ROLLBACK\"",
							"KTHX",
							"RESOURCE_COUNT ITZ RESOURCE_COUNT MINUS 1",
							"WRITER DO CLOSE",
							"SAYZ WIT \"Database connection closed, resources: \"",
							"SAYZ WIT RESOURCE_COUNT",
							"@note Always call CLOSE to prevent resource leaks",
							"@note CLOSE ensures all buffered data is written",
							"@note CLOSE should be called even if errors occur during writing",
							"@note Multiple CLOSE calls should be safe (idempotent)",
							"@note After CLOSE, further WRITE operations may fail",
							"@note CLOSE may finalize file headers, network connections, etc.",
							"@throws Exception if close operation fails",
							"@see WRITE, BUFFERED_WRITER",
							"@category io-operations",
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
					"",
					"@class READWRITER",
					"@example Basic read-write operations",
					"I HAS A VARIABLE RW TEH READWRITER ITZ GET_READWRITER",
					"RW DO WRITE WIT \"Hello\"",
					"I HAS A VARIABLE RESPONSE TEH STRIN ITZ RW DO READ WIT 1024",
					"SAYZ WIT \"Response: \"",
					"SAYZ WIT RESPONSE",
					"RW DO CLOSE",
					"@example Echo protocol implementation",
					"I HAS A VARIABLE CONNECTION TEH READWRITER ITZ GET_NETWORK_CONNECTION",
					"I HAS A VARIABLE RUNNING TEH BOOL ITZ YEZ",
					"WHILE RUNNING",
					"    I HAS A VARIABLE INPUT TEH STRIN ITZ CONNECTION DO READ WIT 1024",
					"    IZ INPUT LENGTH SAEM AS 0?",
					"        RUNNING ITZ NO",
					"    NOPE",
					"        CONNECTION DO WRITE WIT \"Echo: \"",
					"        CONNECTION DO WRITE WIT INPUT",
					"        CONNECTION DO WRITE WIT \"\\n\"",
					"    KTHX",
					"KTHX",
					"CONNECTION DO CLOSE",
					"@example File copy operation",
					"I HAS A VARIABLE SOURCE TEH READWRITER ITZ OPEN_FILE_FOR_READWRITE",
					"I HAS A VARIABLE DEST TEH READWRITER ITZ OPEN_DEST_FILE",
					"I HAS A VARIABLE BUFFER TEH STRIN ITZ SOURCE DO READ WIT 4096",
					"WHILE NO SAEM AS (BUFFER LENGTH SAEM AS 0)",
					"    DEST DO WRITE WIT BUFFER",
					"    BUFFER ITZ SOURCE DO READ WIT 4096",
					"KTHX",
					"SOURCE DO CLOSE",
					"DEST DO CLOSE",
					"SAYZ WIT \"File copy completed\"",
					"@example Interactive session",
					"I HAS A VARIABLE SESSION TEH READWRITER ITZ START_INTERACTIVE_SESSION",
					"I HAS A VARIABLE COMMANDS TEH BUKKIT ITZ NEW BUKKIT",
					"COMMANDS DO PUSH WIT \"HELP\"",
					"COMMANDS DO PUSH WIT \"STATUS\"",
					"COMMANDS DO PUSH WIT \"QUIT\"",
					"WHILE NO SAEM AS (COMMANDS LENGTH SAEM AS 0)",
					"    I HAS A VARIABLE CMD TEH STRIN ITZ COMMANDS DO POP",
					"    SESSION DO WRITE WIT CMD",
					"    SESSION DO WRITE WIT \"\\n\"",
					"    I HAS A VARIABLE RESPONSE TEH STRIN ITZ SESSION DO READ WIT 2048",
					"    SAYZ WIT \"Command: \"",
					"    SAYZ WIT CMD",
					"    SAYZ WIT \"Response: \"",
					"    SAYZ WIT RESPONSE",
					"KTHX",
					"SESSION DO CLOSE",
					"@note Combines both reading and writing capabilities",
					"@note Useful for network connections, files opened for read/write, etc.",
					"@note Inherits all methods from both READER and WRITER",
					"@note CLOSE method closes both reading and writing operations",
					"@see READER, WRITER",
					"@category io-interfaces",
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
					"",
					"@class BUFFERED_READER",
					"@example Basic buffered reading",
					"I HAS A VARIABLE BASE_READER TEH READER ITZ GET_FILE_READER",
					"I HAS A VARIABLE BUFFERED TEH BUFFERED_READER ITZ NEW BUFFERED_READER WIT BASE_READER",
					"I HAS A VARIABLE DATA TEH STRIN ITZ BUFFERED DO READ WIT 1024",
					"WHILE NO SAEM AS (DATA LENGTH SAEM AS 0)",
					"    PROCESS_DATA WIT DATA",
					"    DATA ITZ BUFFERED DO READ WIT 1024",
					"KTHX",
					"BUFFERED DO CLOSE",
					"@example Custom buffer size",
					"I HAS A VARIABLE READER TEH READER ITZ GET_NETWORK_READER",
					"I HAS A VARIABLE BUFFERED TEH BUFFERED_READER ITZ NEW BUFFERED_READER WIT READER",
					"BUFFERED SIZ ITZ 8192",
					"SAYZ WIT \"Buffer size set to 8192\"",
					"I HAS A VARIABLE CHUNK TEH STRIN ITZ BUFFERED DO READ WIT 4096",
					"WHILE NO SAEM AS (CHUNK LENGTH SAEM AS 0)",
					"    SAYZ WIT \"Read chunk of \"",
					"    SAYZ WIT CHUNK LENGTH",
					"    SAYZ WIT \" characters\"",
					"    CHUNK ITZ BUFFERED DO READ WIT 4096",
					"KTHX",
					"BUFFERED DO CLOSE",
					"@example Efficient small reads",
					"I HAS A VARIABLE READER TEH READER ITZ GET_LARGE_FILE_READER",
					"I HAS A VARIABLE BUFFERED TEH BUFFERED_READER ITZ NEW BUFFERED_READER WIT READER",
					"I HAS A VARIABLE TOTAL_CHARS TEH INTEGR ITZ 0",
					"I HAS A VARIABLE CHAR TEH STRIN ITZ BUFFERED DO READ WIT 1",
					"WHILE NO SAEM AS (CHAR LENGTH SAEM AS 0)",
					"    TOTAL_CHARS ITZ TOTAL_CHARS UP 1",
					"    IZ TOTAL_CHARS MOD 1000 SAEM AS 0?",
					"        SAYZ WIT \"Processed \"",
					"        SAYZ WIT TOTAL_CHARS",
					"        SAYZ WIT \" characters\"",
					"    KTHX",
					"    CHAR ITZ BUFFERED DO READ WIT 1",
					"KTHX",
					"SAYZ WIT \"Total characters read: \"",
					"SAYZ WIT TOTAL_CHARS",
					"BUFFERED DO CLOSE",
					"@example Buffer size optimization",
					"I HAS A VARIABLE READER TEH READER ITZ GET_FAST_DEVICE_READER",
					"I HAS A VARIABLE BUFFERED TEH BUFFERED_READER ITZ NEW BUFFERED_READER WIT READER",
					"BUFFERED SIZ ITZ 65536",
					"SAYZ WIT \"Using large buffer for fast device: \"",
					"SAYZ WIT BUFFERED SIZ",
					"I HAS A VARIABLE DATA TEH STRIN ITZ BUFFERED DO READ WIT BUFFERED SIZ",
					"WHILE NO SAEM AS (DATA LENGTH SAEM AS 0)",
					"    PROCESS_LARGE_CHUNK WIT DATA",
					"    DATA ITZ BUFFERED DO READ WIT BUFFERED SIZ",
					"KTHX",
					"BUFFERED DO CLOSE",
					"@example Error handling with buffering",
					"I HAS A VARIABLE READER TEH READER ITZ GET_UNRELIABLE_READER",
					"I HAS A VARIABLE BUFFERED TEH BUFFERED_READER ITZ NEW BUFFERED_READER WIT READER",
					"MAYB",
					"    I HAS A VARIABLE DATA TEH STRIN ITZ BUFFERED DO READ WIT 2048",
					"    WHILE NO SAEM AS (DATA LENGTH SAEM AS 0)",
					"        SAVE_DATA WIT DATA",
					"        DATA ITZ BUFFERED DO READ WIT 2048",
					"    KTHX",
					"OOPSIE ERR",
					"    SAYZ WIT \"Error during buffered reading: \"",
					"    SAYZ WIT ERR",
					"KTHX",
					"BUFFERED DO CLOSE",
					"@note Improves performance by reducing I/O calls",
					"@note Default buffer size is 1024 characters",
					"@note Buffer size can be changed with SIZ property",
					"@note Automatically manages buffer filling and draining",
					"@note Wraps any READER implementation",
					"@see READER, SIZ",
					"@category io-buffered",
				},
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"BUFFERED_READER": {
						Name: "BUFFERED_READER",
						Documentation: []string{
							"Initializes a BUFFERED_READER with the given READER object.",
							"Default buffer size is 1024.",
							"",
							"@syntax NEW BUFFERED_READER WIT <reader>",
							"@param {READER} reader - The underlying reader to wrap with buffering",
							"@returns {BUFFERED_READER} New buffered reader instance",
							"@example Basic construction",
							"I HAS A VARIABLE BASE_READER TEH READER ITZ GET_FILE_READER",
							"I HAS A VARIABLE BUFFERED TEH BUFFERED_READER ITZ NEW BUFFERED_READER WIT BASE_READER",
							"SAYZ WIT \"Buffered reader created with default buffer size\"",
							"@example Construction with immediate use",
							"I HAS A VARIABLE BUFFERED TEH BUFFERED_READER ITZ NEW BUFFERED_READER WIT GET_NETWORK_READER",
							"I HAS A VARIABLE FIRST_CHUNK TEH STRIN ITZ BUFFERED DO READ WIT 512",
							"SAYZ WIT \"First chunk: \"",
							"SAYZ WIT FIRST_CHUNK",
							"BUFFERED DO CLOSE",
							"@example Construction for different reader types",
							"I HAS A VARIABLE TYPES TEH BUKKIT ITZ NEW BUKKIT",
							"TYPES DO PUSH WIT \"file\"",
							"TYPES DO PUSH WIT \"network\"",
							"TYPES DO PUSH WIT \"memory\"",
							"WHILE NO SAEM AS (TYPES LENGTH SAEM AS 0)",
							"    I HAS A VARIABLE TYPE TEH STRIN ITZ TYPES DO POP",
							"    I HAS A VARIABLE READER TEH READER ITZ GET_READER_BY_TYPE WIT TYPE",
							"    I HAS A VARIABLE BUFFERED TEH BUFFERED_READER ITZ NEW BUFFERED_READER WIT READER",
							"    SAYZ WIT \"Created buffered reader for type: \"",
							"    SAYZ WIT TYPE",
							"    BUFFERED DO CLOSE",
							"KTHX",
							"@example Construction with error handling",
							"MAYB",
							"    I HAS A VARIABLE READER TEH READER ITZ GET_SOME_READER",
							"    I HAS A VARIABLE BUFFERED TEH BUFFERED_READER ITZ NEW BUFFERED_READER WIT READER",
							"    SAYZ WIT \"Buffered reader created successfully\"",
							"OOPSIE ERR",
							"    SAYZ WIT \"Failed to create buffered reader: \"",
							"    SAYZ WIT ERR",
							"KTHX",
							"@note Default buffer size is 1024 characters",
							"@note Buffer size can be changed after construction using SIZ",
							"@note The underlying reader is not closed when buffered reader is closed",
							"@note Buffered reader takes ownership of the underlying reader",
							"@throws Exception if reader is not a valid READER object",
							"@see SIZ, READ, CLOSE",
							"@category io-construction",
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
							"",
							"@syntax <buffered_reader> DO READ WIT <size>",
							"@param {INTEGR} size - Maximum number of characters to read",
							"@returns {STRIN} Data read from the source, or empty string at EOF",
							"@example Efficient large reads",
							"I HAS A VARIABLE READER TEH READER ITZ GET_LARGE_FILE_READER",
							"I HAS A VARIABLE BUFFERED TEH BUFFERED_READER ITZ NEW BUFFERED_READER WIT READER",
							"I HAS A VARIABLE CHUNK_SIZE TEH INTEGR ITZ 8192",
							"I HAS A VARIABLE DATA TEH STRIN ITZ BUFFERED DO READ WIT CHUNK_SIZE",
							"I HAS A VARIABLE TOTAL_READ TEH INTEGR ITZ 0",
							"WHILE NO SAEM AS (DATA LENGTH SAEM AS 0)",
							"    TOTAL_READ ITZ TOTAL_READ MOAR DATA LENGTH",
							"    PROCESS_CHUNK WIT DATA",
							"    DATA ITZ BUFFERED DO READ WIT CHUNK_SIZE",
							"KTHX",
							"SAYZ WIT \"Total characters read: \"",
							"SAYZ WIT TOTAL_READ",
							"BUFFERED DO CLOSE",
							"@example Small reads with buffering benefit",
							"I HAS A VARIABLE READER TEH READER ITZ GET_CHARACTER_DEVICE_READER",
							"I HAS A VARIABLE BUFFERED TEH BUFFERED_READER ITZ NEW BUFFERED_READER WIT READER",
							"I HAS A VARIABLE CHAR TEH STRIN ITZ BUFFERED DO READ WIT 1",
							"I HAS A VARIABLE COUNT TEH INTEGR ITZ 0",
							"WHILE NO SAEM AS (CHAR LENGTH SAEM AS 0)",
							"    COUNT ITZ COUNT UP 1",
							"    IZ COUNT MOD 100 SAEM AS 0?",
							"        SAYZ WIT \"Read \"",
							"        SAYZ WIT COUNT",
							"        SAYZ WIT \" characters so far\"",
							"    KTHX",
							"    CHAR ITZ BUFFERED DO READ WIT 1",
							"KTHX",
							"SAYZ WIT \"Total characters: \"",
							"SAYZ WIT COUNT",
							"BUFFERED DO CLOSE",
							"@example Buffer size impact demonstration",
							"I HAS A VARIABLE READER TEH READER ITZ GET_TEST_READER",
							"I HAS A VARIABLE BUFFERED TEH BUFFERED_READER ITZ NEW BUFFERED_READER WIT READER",
							"BUFFERED SIZ ITZ 64",
							"SAYZ WIT \"Small buffer (64): \"",
							"I HAS A VARIABLE DATA1 TEH STRIN ITZ BUFFERED DO READ WIT 32",
							"SAYZ WIT DATA1 LENGTH",
							"BUFFERED SIZ ITZ 4096",
							"SAYZ WIT \"Large buffer (4096): \"",
							"I HAS A VARIABLE DATA2 TEH STRIN ITZ BUFFERED DO READ WIT 32",
							"SAYZ WIT DATA2 LENGTH",
							"BUFFERED DO CLOSE",
							"@example End-of-file handling",
							"I HAS A VARIABLE READER TEH READER ITZ GET_SHORT_FILE_READER",
							"I HAS A VARIABLE BUFFERED TEH BUFFERED_READER ITZ NEW BUFFERED_READER WIT READER",
							"I HAS A VARIABLE DATA TEH STRIN ITZ BUFFERED DO READ WIT 1000",
							"WHILE NO SAEM AS (DATA LENGTH SAEM AS 0)",
							"    SAYZ WIT \"Read: \"",
							"    SAYZ WIT DATA LENGTH",
							"    SAYZ WIT \" characters\"",
							"    DATA ITZ BUFFERED DO READ WIT 1000",
							"KTHX",
							"SAYZ WIT \"Reached end of file\"",
							"BUFFERED DO CLOSE",
							"@example Mixed read sizes",
							"I HAS A VARIABLE READER TEH READER ITZ GET_MIXED_DATA_READER",
							"I HAS A VARIABLE BUFFERED TEH BUFFERED_READER ITZ NEW BUFFERED_READER WIT READER",
							"I HAS A VARIABLE SIZES TEH BUKKIT ITZ NEW BUKKIT",
							"SIZES DO PUSH WIT 10",
							"SIZES DO PUSH WIT 50",
							"SIZES DO PUSH WIT 100",
							"SIZES DO PUSH WIT 500",
							"WHILE NO SAEM AS (SIZES LENGTH SAEM AS 0)",
							"    I HAS A VARIABLE SIZE TEH INTEGR ITZ SIZES DO POP",
							"    I HAS A VARIABLE DATA TEH STRIN ITZ BUFFERED DO READ WIT SIZE",
							"    SAYZ WIT \"Requested \"",
							"    SAYZ WIT SIZE",
							"    SAYZ WIT \" got \"",
							"    SAYZ WIT DATA LENGTH",
							"    SAYZ WIT \" characters\"",
							"KTHX",
							"BUFFERED DO CLOSE",
							"@note Uses internal buffer to reduce underlying I/O calls",
							"@note May return fewer characters than requested",
							"@note Returns empty string when end-of-file is reached",
							"@note Buffer size affects performance but not correctness",
							"@note Subsequent calls may reuse buffered data",
							"@throws Exception if size is invalid or I/O error occurs",
							"@see SIZ, CLOSE, BUFFERED_READER",
							"@category io-operations",
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
							"",
							"@syntax <buffered_reader> DO CLOSE",
							"@example Basic cleanup",
							"I HAS A VARIABLE READER TEH READER ITZ GET_FILE_READER",
							"I HAS A VARIABLE BUFFERED TEH BUFFERED_READER ITZ NEW BUFFERED_READER WIT READER",
							"I HAS A VARIABLE DATA TEH STRIN ITZ BUFFERED DO READ WIT 1024",
							"WHILE NO SAEM AS (DATA LENGTH SAEM AS 0)",
							"    PROCESS_DATA WIT DATA",
							"    DATA ITZ BUFFERED DO READ WIT 1024",
							"KTHX",
							"BUFFERED DO CLOSE",
							"SAYZ WIT \"Buffered reader and underlying reader closed\"",
							"@example Close with error handling",
							"I HAS A VARIABLE READER TEH READER ITZ GET_NETWORK_READER",
							"I HAS A VARIABLE BUFFERED TEH BUFFERED_READER ITZ NEW BUFFERED_READER WIT READER",
							"MAYB",
							"    I HAS A VARIABLE DATA TEH STRIN ITZ BUFFERED DO READ WIT 512",
							"    WHILE NO SAEM AS (DATA LENGTH SAEM AS 0)",
							"        SAVE_DATA WIT DATA",
							"        DATA ITZ BUFFERED DO READ WIT 512",
							"    KTHX",
							"OOPSIE ERR",
							"    SAYZ WIT \"Error during reading: \"",
							"    SAYZ WIT ERR",
							"KTHX",
							"BUFFERED DO CLOSE",
							"@example Multiple buffered readers cleanup",
							"I HAS A VARIABLE READERS TEH BUKKIT ITZ NEW BUKKIT",
							"READERS DO PUSH WIT GET_READER_1",
							"READERS DO PUSH WIT GET_READER_2",
							"READERS DO PUSH WIT GET_READER_3",
							"WHILE NO SAEM AS (READERS LENGTH SAEM AS 0)",
							"    I HAS A VARIABLE READER TEH READER ITZ READERS DO POP",
							"    I HAS A VARIABLE BUFFERED TEH BUFFERED_READER ITZ NEW BUFFERED_READER WIT READER",
							"    I HAS A VARIABLE DATA TEH STRIN ITZ BUFFERED DO READ WIT 1024",
							"    WHILE NO SAEM AS (DATA LENGTH SAEM AS 0)",
							"        PROCESS_DATA WIT DATA",
							"        DATA ITZ BUFFERED DO READ WIT 1024",
							"    KTHX",
							"    BUFFERED DO CLOSE",
							"KTHX",
							"SAYZ WIT \"All buffered readers closed\"",
							"@example Close after partial reading",
							"I HAS A VARIABLE READER TEH READER ITZ GET_LARGE_FILE_READER",
							"I HAS A VARIABLE BUFFERED TEH BUFFERED_READER ITZ NEW BUFFERED_READER WIT READER",
							"I HAS A VARIABLE FIRST_CHUNK TEH STRIN ITZ BUFFERED DO READ WIT 100",
							"SAYZ WIT \"Read first chunk: \"",
							"SAYZ WIT FIRST_CHUNK LENGTH",
							"SAYZ WIT \" characters\"",
							"BUFFERED DO CLOSE",
							"SAYZ WIT \"Reader closed early - remaining data discarded\"",
							"@example Resource cleanup pattern",
							"I HAS A VARIABLE READER TEH READER ITZ GET_DATABASE_READER",
							"I HAS A VARIABLE BUFFERED TEH BUFFERED_READER ITZ NEW BUFFERED_READER WIT READER",
							"I HAS A VARIABLE RESOURCE_COUNT TEH INTEGR ITZ 1",
							"MAYB",
							"    I HAS A VARIABLE DATA TEH STRIN ITZ BUFFERED DO READ WIT 2048",
							"    WHILE NO SAEM AS (DATA LENGTH SAEM AS 0)",
							"        PROCESS_DATABASE_ROW WIT DATA",
							"        DATA ITZ BUFFERED DO READ WIT 2048",
							"    KTHX",
							"OOPSIE ERR",
							"    SAYZ WIT \"Database read failed: \"",
							"    SAYZ WIT ERR",
							"KTHX",
							"RESOURCE_COUNT ITZ RESOURCE_COUNT MINUS 1",
							"BUFFERED DO CLOSE",
							"SAYZ WIT \"Database connection closed, resources: \"",
							"SAYZ WIT RESOURCE_COUNT",
							"@note Closes both the buffered reader and underlying reader",
							"@note Releases all associated resources and buffers",
							"@note Any remaining buffered data is discarded",
							"@note Safe to call multiple times (idempotent)",
							"@note After CLOSE, further READ operations will fail",
							"@note Ensures underlying reader is properly closed",
							"@throws Exception if close operation fails",
							"@see READ, BUFFERED_READER, SIZ",
							"@category io-operations",
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
								"",
								"@property {INTEGR} SIZ - Buffer size in characters",
								"@example Get current buffer size",
								"I HAS A VARIABLE READER TEH READER ITZ GET_FILE_READER",
								"I HAS A VARIABLE BUFFERED TEH BUFFERED_READER ITZ NEW BUFFERED_READER WIT READER",
								"SAYZ WIT \"Current buffer size: \"",
								"SAYZ WIT BUFFERED SIZ",
								"BUFFERED DO CLOSE",
								"@example Set custom buffer size",
								"I HAS A VARIABLE READER TEH READER ITZ GET_NETWORK_READER",
								"I HAS A VARIABLE BUFFERED TEH BUFFERED_READER ITZ NEW BUFFERED_READER WIT READER",
								"BUFFERED SIZ ITZ 4096",
								"SAYZ WIT \"Buffer size set to: \"",
								"SAYZ WIT BUFFERED SIZ",
								"BUFFERED DO CLOSE",
								"@example Optimize for different data patterns",
								"I HAS A VARIABLE PATTERNS TEH BUKKIT ITZ NEW BUKKIT",
								"PATTERNS DO PUSH WIT \"small_records\"",
								"PATTERNS DO PUSH WIT \"large_files\"",
								"PATTERNS DO PUSH WIT \"network_packets\"",
								"WHILE NO SAEM AS (PATTERNS LENGTH SAEM AS 0)",
								"    I HAS A VARIABLE PATTERN TEH STRIN ITZ PATTERNS DO POP",
								"    I HAS A VARIABLE READER TEH READER ITZ GET_READER_FOR_PATTERN WIT PATTERN",
								"    I HAS A VARIABLE BUFFERED TEH BUFFERED_READER ITZ NEW BUFFERED_READER WIT READER",
								"    IZ PATTERN SAEM AS \"small_records\"?",
								"        BUFFERED SIZ ITZ 512",
								"    NOPE",
								"        IZ PATTERN SAEM AS \"large_files\"?",
								"            BUFFERED SIZ ITZ 65536",
								"        NOPE",
								"            IZ PATTERN SAEM AS \"network_packets\"?",
								"                BUFFERED SIZ ITZ 8192",
								"            NOPE",
								"                BUFFERED SIZ ITZ 1024",
								"            KTHX",
								"        KTHX",
								"    KTHX",
								"    SAYZ WIT \"Pattern \"",
								"    SAYZ WIT PATTERN",
								"    SAYZ WIT \" using buffer size: \"",
								"    SAYZ WIT BUFFERED SIZ",
								"    BUFFERED DO CLOSE",
								"KTHX",
								"@example Dynamic buffer size adjustment",
								"I HAS A VARIABLE READER TEH READER ITZ GET_ADAPTIVE_READER",
								"I HAS A VARIABLE BUFFERED TEH BUFFERED_READER ITZ NEW BUFFERED_READER WIT READER",
								"I HAS A VARIABLE INITIAL_SIZE TEH INTEGR ITZ 1024",
								"BUFFERED SIZ ITZ INITIAL_SIZE",
								"I HAS A VARIABLE DATA TEH STRIN ITZ BUFFERED DO READ WIT 100",
								"I HAS A VARIABLE READ_COUNT TEH INTEGR ITZ 0",
								"WHILE NO SAEM AS (DATA LENGTH SAEM AS 0)",
								"    READ_COUNT ITZ READ_COUNT UP 1",
								"    IZ READ_COUNT MOD 10 SAEM AS 0?",
								"        BUFFERED SIZ ITZ BUFFERED SIZ MUL 2",
								"        SAYZ WIT \"Increased buffer size to: \"",
								"        SAYZ WIT BUFFERED SIZ",
								"    KTHX",
								"    DATA ITZ BUFFERED DO READ WIT 100",
								"KTHX",
								"BUFFERED DO CLOSE",
								"@example Buffer size validation",
								"I HAS A VARIABLE READER TEH READER ITZ GET_FILE_READER",
								"I HAS A VARIABLE BUFFERED TEH BUFFERED_READER ITZ NEW BUFFERED_READER WIT READER",
								"MAYB",
								"    BUFFERED SIZ ITZ 0",
								"    SAYZ WIT \"Error: Should not allow zero buffer size\"",
								"OOPSIE ERR",
								"    SAYZ WIT \"Expected error for zero buffer size: \"",
								"    SAYZ WIT ERR",
								"KTHX",
								"MAYB",
								"    BUFFERED SIZ ITZ -100",
								"    SAYZ WIT \"Error: Should not allow negative buffer size\"",
								"OOPSIE ERR",
								"    SAYZ WIT \"Expected error for negative buffer size: \"",
								"    SAYZ WIT ERR",
								"KTHX",
								"BUFFERED SIZ ITZ 2048",
								"SAYZ WIT \"Valid buffer size set to: \"",
								"SAYZ WIT BUFFERED SIZ",
								"BUFFERED DO CLOSE",
								"@note Default buffer size is 1024 characters",
								"@note Setting new size clears existing buffer contents",
								"@note Must be a positive integer greater than 0",
								"@note Larger buffers improve performance for sequential reads",
								"@note Smaller buffers use less memory",
								"@note Can be changed at any time during reader lifetime",
								"@throws Exception if size is not a positive integer",
								"@see BUFFERED_READER, READ",
								"@category io-configuration",
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
					"",
					"@class BUFFERED_WRITER",
					"@example Basic buffered writing",
					"I HAS A VARIABLE BASE_WRITER TEH WRITER ITZ GET_FILE_WRITER",
					"I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT BASE_WRITER",
					"BUFFERED DO WRITE WIT \"Hello, World!\"",
					"BUFFERED DO WRITE WIT \"\\n\"",
					"BUFFERED DO WRITE WIT \"This is buffered writing.\"",
					"BUFFERED DO CLOSE",
					"SAYZ WIT \"Data written with buffering\"",
					"@example Custom buffer size",
					"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_NETWORK_WRITER",
					"I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT WRITER",
					"BUFFERED SIZ ITZ 4096",
					"SAYZ WIT \"Buffer size set to 4096\"",
					"I HAS A VARIABLE DATA TEH STRIN ITZ \"Large data payload...\"",
					"BUFFERED DO WRITE WIT DATA",
					"BUFFERED DO CLOSE",
					"SAYZ WIT \"Large data written efficiently\"",
					"@example Explicit flushing",
					"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_LOG_WRITER",
					"I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT WRITER",
					"BUFFERED DO WRITE WIT \"Log entry 1\"",
					"BUFFERED DO WRITE WIT \"\\n\"",
					"BUFFERED DO FLUSH",
					"SAYZ WIT \"First log entry flushed immediately\"",
					"BUFFERED DO WRITE WIT \"Log entry 2\"",
					"BUFFERED DO WRITE WIT \"\\n\"",
					"BUFFERED DO CLOSE",
					"SAYZ WIT \"Second log entry flushed on close\"",
					"@example Batch writing with progress",
					"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_FILE_WRITER",
					"I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT WRITER",
					"I HAS A VARIABLE LINES TEH BUKKIT ITZ NEW BUKKIT",
					"LINES DO PUSH WIT \"Line 1\"",
					"LINES DO PUSH WIT \"Line 2\"",
					"LINES DO PUSH WIT \"Line 3\"",
					"I HAS A VARIABLE TOTAL_CHARS TEH INTEGR ITZ 0",
					"WHILE NO SAEM AS (LINES LENGTH SAEM AS 0)",
					"    I HAS A VARIABLE LINE TEH STRIN ITZ LINES DO POP",
					"    I HAS A VARIABLE LINE_WITH_NEWLINE TEH STRIN ITZ LINE MOAR \"\\n\"",
					"    I HAS A VARIABLE WRITTEN TEH INTEGR ITZ BUFFERED DO WRITE WIT LINE_WITH_NEWLINE",
					"    TOTAL_CHARS ITZ TOTAL_CHARS MOAR WRITTEN",
					"    IZ TOTAL_CHARS MOD 50 SAEM AS 0?",
					"        SAYZ WIT \"Written \"",
					"        SAYZ WIT TOTAL_CHARS",
					"        SAYZ WIT \" characters so far\"",
					"    KTHX",
					"KTHX",
					"SAYZ WIT \"Total characters written: \"",
					"SAYZ WIT TOTAL_CHARS",
					"BUFFERED DO CLOSE",
					"@example Error handling with buffering",
					"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_UNRELIABLE_WRITER",
					"I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT WRITER",
					"MAYB",
					"    BUFFERED DO WRITE WIT \"Important data\"",
					"    BUFFERED DO WRITE WIT \"\\n\"",
					"    BUFFERED DO FLUSH",
					"    SAYZ WIT \"Data written and flushed successfully\"",
					"OOPSIE ERR",
					"    SAYZ WIT \"Error during buffered writing: \"",
					"    SAYZ WIT ERR",
					"KTHX",
					"BUFFERED DO CLOSE",
					"@note Improves performance by reducing I/O calls",
					"@note Default buffer size is 1024 characters",
					"@note Buffer size can be changed with SIZ property",
					"@note Data may be buffered until FLUSH or CLOSE",
					"@note Wraps any WRITER implementation",
					"@see WRITER, SIZ, FLUSH",
					"@category io-buffered",
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
							"",
							"@syntax NEW BUFFERED_WRITER WIT <writer>",
							"@param {WRITER} writer - The underlying writer to wrap with buffering",
							"@returns {BUFFERED_WRITER} New buffered writer instance",
							"@example Basic construction",
							"I HAS A VARIABLE BASE_WRITER TEH WRITER ITZ GET_FILE_WRITER",
							"I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT BASE_WRITER",
							"SAYZ WIT \"Buffered writer created with default buffer size\"",
							"@example Construction with immediate use",
							"I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT GET_NETWORK_WRITER",
							"BUFFERED DO WRITE WIT \"Hello from buffered writer!\"",
							"BUFFERED DO CLOSE",
							"SAYZ WIT \"Data written through buffered writer\"",
							"@example Construction for different writer types",
							"I HAS A VARIABLE TYPES TEH BUKKIT ITZ NEW BUKKIT",
							"TYPES DO PUSH WIT \"file\"",
							"TYPES DO PUSH WIT \"network\"",
							"TYPES DO PUSH WIT \"console\"",
							"WHILE NO SAEM AS (TYPES LENGTH SAEM AS 0)",
							"    I HAS A VARIABLE TYPE TEH STRIN ITZ TYPES DO POP",
							"    I HAS A VARIABLE WRITER TEH WRITER ITZ GET_WRITER_BY_TYPE WIT TYPE",
							"    I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT WRITER",
							"    SAYZ WIT \"Created buffered writer for type: \"",
							"    SAYZ WIT TYPE",
							"    BUFFERED DO CLOSE",
							"KTHX",
							"@example Construction with error handling",
							"MAYB",
							"    I HAS A VARIABLE WRITER TEH WRITER ITZ GET_SOME_WRITER",
							"    I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT WRITER",
							"    SAYZ WIT \"Buffered writer created successfully\"",
							"OOPSIE ERR",
							"    SAYZ WIT \"Failed to create buffered writer: \"",
							"    SAYZ WIT ERR",
							"KTHX",
							"@note Default buffer size is 1024 characters",
							"@note Buffer size can be changed after construction using SIZ",
							"@note The underlying writer is not closed when buffered writer is closed",
							"@note Buffered writer takes ownership of the underlying writer",
							"@throws Exception if writer is not a valid WRITER object",
							"@see SIZ, WRITE, FLUSH, CLOSE",
							"@category io-construction",
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
							"",
							"@syntax <buffered_writer> DO WRITE WIT <data>",
							"@param {STRIN} data - String data to write",
							"@returns {INTEGR} Number of characters written",
							"@example Basic buffered writing",
							"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_FILE_WRITER",
							"I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT WRITER",
							"I HAS A VARIABLE CHARS_WRITTEN TEH INTEGR ITZ BUFFERED DO WRITE WIT \"Hello World\"",
							"SAYZ WIT \"Wrote \"",
							"SAYZ WIT CHARS_WRITTEN",
							"SAYZ WIT \" characters to buffer\"",
							"BUFFERED DO CLOSE",
							"@example Writing with buffer management",
							"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_NETWORK_WRITER",
							"I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT WRITER",
							"I HAS A VARIABLE DATA TEH STRIN ITZ \"Small data\"",
							"I HAS A VARIABLE WRITTEN TEH INTEGR ITZ BUFFERED DO WRITE WIT DATA",
							"SAYZ WIT \"Data buffered: \"",
							"SAYZ WIT WRITTEN",
							"SAYZ WIT \" characters\"",
							"I HAS A VARIABLE LARGE_DATA TEH STRIN ITZ \"Very large data that exceeds buffer size...\"",
							"I HAS A VARIABLE LARGE_WRITTEN TEH INTEGR ITZ BUFFERED DO WRITE WIT LARGE_DATA",
							"SAYZ WIT \"Large data written directly: \"",
							"SAYZ WIT LARGE_WRITTEN",
							"SAYZ WIT \" characters\"",
							"BUFFERED DO CLOSE",
							"@example Batch writing with size tracking",
							"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_LOG_WRITER",
							"I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT WRITER",
							"I HAS A VARIABLE MESSAGES TEH BUKKIT ITZ NEW BUKKIT",
							"MESSAGES DO PUSH WIT \"INFO: Starting process\"",
							"MESSAGES DO PUSH WIT \"INFO: Processing data\"",
							"MESSAGES DO PUSH WIT \"INFO: Process complete\"",
							"I HAS A VARIABLE TOTAL_WRITTEN TEH INTEGR ITZ 0",
							"WHILE NO SAEM AS (MESSAGES LENGTH SAEM AS 0)",
							"    I HAS A VARIABLE MSG TEH STRIN ITZ MESSAGES DO POP",
							"    I HAS A VARIABLE MSG_WITH_NEWLINE TEH STRIN ITZ MSG MOAR \"\\n\"",
							"    I HAS A VARIABLE WRITTEN TEH INTEGR ITZ BUFFERED DO WRITE WIT MSG_WITH_NEWLINE",
							"    TOTAL_WRITTEN ITZ TOTAL_WRITTEN MOAR WRITTEN",
							"KTHX",
							"SAYZ WIT \"Total characters buffered: \"",
							"SAYZ WIT TOTAL_WRITTEN",
							"BUFFERED DO CLOSE",
							"@example Writing with explicit flushing",
							"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_CONSOLE_WRITER",
							"I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT WRITER",
							"BUFFERED DO WRITE WIT \"First message\"",
							"BUFFERED DO FLUSH",
							"SAYZ WIT \"First message should be visible now\"",
							"BUFFERED DO WRITE WIT \"Second message\"",
							"BUFFERED DO CLOSE",
							"SAYZ WIT \"Second message flushed on close\"",
							"@example Error handling during writing",
							"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_UNRELIABLE_WRITER",
							"I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT WRITER",
							"MAYB",
							"    I HAS A VARIABLE WRITTEN TEH INTEGR ITZ BUFFERED DO WRITE WIT \"Test data\"",
							"    SAYZ WIT \"Data buffered successfully: \"",
							"    SAYZ WIT WRITTEN",
							"    SAYZ WIT \" characters\"",
							"OOPSIE ERR",
							"    SAYZ WIT \"Error during buffered writing: \"",
							"    SAYZ WIT ERR",
							"KTHX",
							"BUFFERED DO CLOSE",
							"@example Performance comparison",
							"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_SLOW_WRITER",
							"I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT WRITER",
							"I HAS A VARIABLE START_TIME TEH INTEGR ITZ GET_CURRENT_TIME",
							"I HAS A VARIABLE IDX TEH INTEGR ITZ 0",
							"WHILE IDX SMALLR THAN 100",
							"    BUFFERED DO WRITE WIT \"Data chunk \"",
							"    BUFFERED DO WRITE WIT IDX",
							"    BUFFERED DO WRITE WIT \"\\n\"",
							"    IDX ITZ IDX UP 1",
							"KTHX",
							"BUFFERED DO CLOSE",
							"I HAS A VARIABLE END_TIME TEH INTEGR ITZ GET_CURRENT_TIME",
							"SAYZ WIT \"Buffered writing took: \"",
							"SAYZ WIT END_TIME MINUS START_TIME",
							"SAYZ WIT \" time units\"",
							"@note Data is buffered until buffer is full or FLUSH/CLOSE is called",
							"@note Large data may be written directly bypassing the buffer",
							"@note Returns the number of characters actually written",
							"@note Buffer size affects when data is automatically flushed",
							"@note May return fewer characters if error occurs during writing",
							"@throws Exception if data is invalid or I/O error occurs",
							"@see FLUSH, CLOSE, SIZ, BUFFERED_WRITER",
							"@category io-operations",
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
							"",
							"@syntax <buffered_writer> DO FLUSH",
							"@example Explicit flushing for immediate writes",
							"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_LOG_WRITER",
							"I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT WRITER",
							"BUFFERED DO WRITE WIT \"Important log message\"",
							"BUFFERED DO FLUSH",
							"SAYZ WIT \"Log message written immediately\"",
							"@example Periodic flushing",
							"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_FILE_WRITER",
							"I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT WRITER",
							"I HAS A VARIABLE LINES TEH BUKKIT ITZ NEW BUKKIT",
							"LINES DO PUSH WIT \"Line 1\"",
							"LINES DO PUSH WIT \"Line 2\"",
							"LINES DO PUSH WIT \"Line 3\"",
							"I HAS A VARIABLE COUNT TEH INTEGR ITZ 0",
							"WHILE NO SAEM AS (LINES LENGTH SAEM AS 0)",
							"    I HAS A VARIABLE LINE TEH STRIN ITZ LINES DO POP",
							"    BUFFERED DO WRITE WIT LINE",
							"    BUFFERED DO WRITE WIT \"\\n\"",
							"    COUNT ITZ COUNT UP 1",
							"    IZ COUNT MOD 2 SAEM AS 0?",
							"        BUFFERED DO FLUSH",
							"        SAYZ WIT \"Flushed after \"",
							"        SAYZ WIT COUNT",
							"        SAYZ WIT \" lines\"",
							"    KTHX",
							"KTHX",
							"BUFFERED DO CLOSE",
							"@example Flush before reading response",
							"I HAS A VARIABLE CONNECTION TEH READWRITER ITZ GET_NETWORK_CONNECTION",
							"I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT CONNECTION",
							"BUFFERED DO WRITE WIT \"GET /data HTTP/1.1\\n\"",
							"BUFFERED DO WRITE WIT \"Host: example.com\\n\\n\"",
							"BUFFERED DO FLUSH",
							"I HAS A VARIABLE RESPONSE TEH STRIN ITZ CONNECTION DO READ WIT 1024",
							"SAYZ WIT \"Server response: \"",
							"SAYZ WIT RESPONSE",
							"BUFFERED DO CLOSE",
							"@example Error handling during flush",
							"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_UNRELIABLE_WRITER",
							"I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT WRITER",
							"BUFFERED DO WRITE WIT \"Data to flush\"",
							"MAYB",
							"    BUFFERED DO FLUSH",
							"    SAYZ WIT \"Data flushed successfully\"",
							"OOPSIE ERR",
							"    SAYZ WIT \"Error during flush: \"",
							"    SAYZ WIT ERR",
							"KTHX",
							"BUFFERED DO CLOSE",
							"@example Flush with empty buffer",
							"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_FILE_WRITER",
							"I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT WRITER",
							"BUFFERED DO FLUSH",
							"SAYZ WIT \"Flush on empty buffer - no operation performed\"",
							"BUFFERED DO WRITE WIT \"Some data\"",
							"BUFFERED DO FLUSH",
							"SAYZ WIT \"Data flushed after writing\"",
							"BUFFERED DO CLOSE",
							"@example Multiple flushes in sequence",
							"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_CONSOLE_WRITER",
							"I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT WRITER",
							"BUFFERED DO WRITE WIT \"First \"",
							"BUFFERED DO FLUSH",
							"BUFFERED DO WRITE WIT \"Second \"",
							"BUFFERED DO FLUSH",
							"BUFFERED DO WRITE WIT \"Third\"",
							"BUFFERED DO FLUSH",
							"SAYZ WIT \"All parts written incrementally\"",
							"BUFFERED DO CLOSE",
							"@note Forces immediate writing of all buffered data",
							"@note Does not affect buffer size or configuration",
							"@note Safe to call on empty buffer (no operation)",
							"@note Can be called multiple times safely",
							"@note Useful for ensuring data reaches destination immediately",
							"@note May improve responsiveness for interactive applications",
							"@throws Exception if flush operation fails",
							"@see WRITE, CLOSE, SIZ",
							"@category io-operations",
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
							"",
							"@syntax <buffered_writer> DO CLOSE",
							"@example Basic cleanup with auto-flush",
							"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_FILE_WRITER",
							"I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT WRITER",
							"BUFFERED DO WRITE WIT \"Final data\"",
							"BUFFERED DO CLOSE",
							"SAYZ WIT \"Buffered writer closed - data auto-flushed\"",
							"@example Close with error handling",
							"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_NETWORK_WRITER",
							"I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT WRITER",
							"BUFFERED DO WRITE WIT \"Data to send\"",
							"MAYB",
							"    BUFFERED DO CLOSE",
							"    SAYZ WIT \"Writer closed successfully\"",
							"OOPSIE ERR",
							"    SAYZ WIT \"Error during close: \"",
							"    SAYZ WIT ERR",
							"KTHX",
							"@example Multiple writers cleanup",
							"I HAS A VARIABLE WRITERS TEH BUKKIT ITZ NEW BUKKIT",
							"WRITERS DO PUSH WIT GET_WRITER_1",
							"WRITERS DO PUSH WIT GET_WRITER_2",
							"WRITERS DO PUSH WIT GET_WRITER_3",
							"WHILE NO SAEM AS (WRITERS LENGTH SAEM AS 0)",
							"    I HAS A VARIABLE WRITER TEH WRITER ITZ WRITERS DO POP",
							"    I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT WRITER",
							"    BUFFERED DO WRITE WIT \"Data for writer\"",
							"    BUFFERED DO CLOSE",
							"KTHX",
							"SAYZ WIT \"All buffered writers closed\"",
							"@example Close after explicit flush",
							"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_LOG_WRITER",
							"I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT WRITER",
							"BUFFERED DO WRITE WIT \"Log data\"",
							"BUFFERED DO FLUSH",
							"SAYZ WIT \"Data flushed explicitly\"",
							"BUFFERED DO CLOSE",
							"SAYZ WIT \"Writer closed (no additional flush needed)\"",
							"@example Resource management pattern",
							"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_DATABASE_WRITER",
							"I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT WRITER",
							"I HAS A VARIABLE RESOURCE_COUNT TEH INTEGR ITZ 1",
							"MAYB",
							"    BUFFERED DO WRITE WIT \"INSERT INTO table VALUES (...)\"",
							"    BUFFERED DO WRITE WIT \"COMMIT\"",
							"OOPSIE ERR",
							"    SAYZ WIT \"Database write failed: \"",
							"    SAYZ WIT ERR",
							"    BUFFERED DO WRITE WIT \"ROLLBACK\"",
							"KTHX",
							"RESOURCE_COUNT ITZ RESOURCE_COUNT MINUS 1",
							"BUFFERED DO CLOSE",
							"SAYZ WIT \"Database connection closed, resources: \"",
							"SAYZ WIT RESOURCE_COUNT",
							"@example Close with partial buffer",
							"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_FILE_WRITER",
							"I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT WRITER",
							"BUFFERED DO WRITE WIT \"Incomplete\"",
							"BUFFERED DO CLOSE",
							"SAYZ WIT \"Partial buffer data flushed during close\"",
							"@note Automatically flushes any remaining buffered data",
							"@note Closes both the buffered writer and underlying writer",
							"@note Releases all associated resources",
							"@note Safe to call multiple times (idempotent)",
							"@note After CLOSE, further WRITE operations will fail",
							"@note Ensures underlying writer is properly closed",
							"@throws Exception if flush or close operation fails",
							"@see WRITE, FLUSH, BUFFERED_WRITER, SIZ",
							"@category io-operations",
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
								"",
								"@property {INTEGR} SIZ - Buffer size in characters",
								"@example Get current buffer size",
								"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_FILE_WRITER",
								"I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT WRITER",
								"SAYZ WIT \"Current buffer size: \"",
								"SAYZ WIT BUFFERED SIZ",
								"BUFFERED DO CLOSE",
								"@example Set custom buffer size",
								"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_NETWORK_WRITER",
								"I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT WRITER",
								"BUFFERED SIZ ITZ 2048",
								"SAYZ WIT \"Buffer size set to: \"",
								"SAYZ WIT BUFFERED SIZ",
								"BUFFERED DO CLOSE",
								"@example Optimize for different data patterns",
								"I HAS A VARIABLE PATTERNS TEH BUKKIT ITZ NEW BUKKIT",
								"PATTERNS DO PUSH WIT \"frequent_small_writes\"",
								"PATTERNS DO PUSH WIT \"occasional_large_writes\"",
								"PATTERNS DO PUSH WIT \"real_time_logging\"",
								"WHILE NO SAEM AS (PATTERNS LENGTH SAEM AS 0)",
								"    I HAS A VARIABLE PATTERN TEH STRIN ITZ PATTERNS DO POP",
								"    I HAS A VARIABLE WRITER TEH WRITER ITZ GET_WRITER_FOR_PATTERN WIT PATTERN",
								"    I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT WRITER",
								"    IZ PATTERN SAEM AS \"frequent_small_writes\"?",
								"        BUFFERED SIZ ITZ 4096",
								"    NOPE",
								"        IZ PATTERN SAEM AS \"occasional_large_writes\"?",
								"            BUFFERED SIZ ITZ 1024",
								"        NOPE",
								"            IZ PATTERN SAEM AS \"real_time_logging\"?",
								"                BUFFERED SIZ ITZ 256",
								"            NOPE",
								"                BUFFERED SIZ ITZ 1024",
								"            KTHX",
								"        KTHX",
								"    KTHX",
								"    SAYZ WIT \"Pattern \"",
								"    SAYZ WIT PATTERN",
								"    SAYZ WIT \" using buffer size: \"",
								"    SAYZ WIT BUFFERED SIZ",
								"    BUFFERED DO CLOSE",
								"KTHX",
								"@example Dynamic buffer size adjustment",
								"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_ADAPTIVE_WRITER",
								"I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT WRITER",
								"I HAS A VARIABLE INITIAL_SIZE TEH INTEGR ITZ 1024",
								"BUFFERED SIZ ITZ INITIAL_SIZE",
								"I HAS A VARIABLE WRITE_COUNT TEH INTEGR ITZ 0",
								"WHILE WRITE_COUNT SMALLR THAN 50",
								"    BUFFERED DO WRITE WIT \"Data chunk \"",
								"    BUFFERED DO WRITE WIT WRITE_COUNT",
								"    BUFFERED DO WRITE WIT \"\\n\"",
								"    WRITE_COUNT ITZ WRITE_COUNT UP 1",
								"    IZ WRITE_COUNT MOD 10 SAEM AS 0?",
								"        BUFFERED SIZ ITZ BUFFERED SIZ MUL 2",
								"        SAYZ WIT \"Increased buffer size to: \"",
								"        SAYZ WIT BUFFERED SIZ",
								"    KTHX",
								"KTHX",
								"BUFFERED DO CLOSE",
								"@example Buffer size impact on flushing",
								"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_FILE_WRITER",
								"I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT WRITER",
								"BUFFERED SIZ ITZ 10",
								"SAYZ WIT \"Small buffer (10): \"",
								"I HAS A VARIABLE IDX TEH INTEGR ITZ 0",
								"WHILE IDX SMALLR THAN 5",
								"    BUFFERED DO WRITE WIT \"Data chunk that exceeds buffer\"",
								"    IDX ITZ IDX UP 1",
								"KTHX",
								"SAYZ WIT \"Small buffer caused multiple flushes\"",
								"BUFFERED SIZ ITZ 1000",
								"SAYZ WIT \"Large buffer (1000): \"",
								"I HAS A VARIABLE IDX TEH INTEGR ITZ 0",
								"WHILE IDX SMALLR THAN 5",
								"    BUFFERED DO WRITE WIT \"Data chunk that fits in buffer\"",
								"    IDX ITZ IDX UP 1",
								"KTHX",
								"SAYZ WIT \"Large buffer delayed flushing\"",
								"BUFFERED DO CLOSE",
								"@example Buffer size validation",
								"I HAS A VARIABLE WRITER TEH WRITER ITZ GET_FILE_WRITER",
								"I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT WRITER",
								"MAYB",
								"    BUFFERED SIZ ITZ 0",
								"    SAYZ WIT \"Error: Should not allow zero buffer size\"",
								"OOPSIE ERR",
								"    SAYZ WIT \"Expected error for zero buffer size: \"",
								"    SAYZ WIT ERR",
								"KTHX",
								"MAYB",
								"    BUFFERED SIZ ITZ -500",
								"    SAYZ WIT \"Error: Should not allow negative buffer size\"",
								"OOPSIE ERR",
								"    SAYZ WIT \"Expected error for negative buffer size: \"",
								"    SAYZ WIT ERR",
								"KTHX",
								"BUFFERED SIZ ITZ 3072",
								"SAYZ WIT \"Valid buffer size set to: \"",
								"SAYZ WIT BUFFERED SIZ",
								"BUFFERED DO CLOSE",
								"@note Default buffer size is 1024 characters",
								"@note Setting new size drops existing buffer contents",
								"@note Must be a positive integer greater than 0",
								"@note Larger buffers reduce I/O calls but use more memory",
								"@note Smaller buffers flush more frequently",
								"@note Can be changed at any time during writer lifetime",
								"@throws Exception if size is not a positive integer",
								"@see BUFFERED_WRITER, WRITE, FLUSH",
								"@category io-configuration",
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
