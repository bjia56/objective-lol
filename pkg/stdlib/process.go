package stdlib

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	goRuntime "runtime"
	"strings"
	"sync"
	"syscall"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/runtime"
)

// moduleProcessCategories defines the order that categories should be rendered in documentation
var moduleProcessCategories = []string{
	"process-creation",
	"process-configuration",
	"process-lifecycle",
	"process-communication",
	"process-control",
	"pipe-operations",
	"pipe-properties",
}

// PipeData stores the internal state of a PIPE
type PipeData struct {
	FDType  string    // "STDIN", "STDOUT", or "STDERR"
	Pipe    io.Closer // The actual pipe (io.WriteCloser for stdin, io.ReadCloser for stdout/stderr)
	IsOpen  bool      // Whether pipe is still open
	IsEOF   bool      // Whether EOF reached (for read pipes)
	Process *exec.Cmd // Reference to the parent process
}

// MinionData stores the internal state of a MINION
type MinionData struct {
	CmdLine    *environment.ObjectInstance // Command and arguments (BUKKIT)
	WorkDir    string                      // Working directory
	Env        *environment.ObjectInstance // Environment variables (BASKIT)
	Process    *exec.Cmd                   // The actual process
	Running    bool                        // Whether process is running
	Finished   bool                        // Whether process has completed
	ExitCode   int                         // Process exit code
	PID        int                         // Process ID
	StdinPipe  *environment.ObjectInstance // PIPE object for stdin
	StdoutPipe *environment.ObjectInstance // PIPE object for stdout
	StderrPipe *environment.ObjectInstance // PIPE object for stderr
}

// NewPipeInstance creates a new PIPE object instance
func NewPipeInstance(fdType string, process *exec.Cmd) *environment.ObjectInstance {
	class := getProcessClasses()["PIPE"]
	env := environment.NewEnvironment(nil)
	env.DefineClass(class)
	_ = RegisterIOInEnv(env, "READERWRITER", "READER", "WRITER") // Ensure IO interfaces are registered
	obj := &environment.ObjectInstance{
		Environment: env,
		Class:       class,
		NativeData: &PipeData{
			FDType:  fdType,
			Pipe:    nil,
			IsOpen:  false,
			IsEOF:   false,
			Process: process,
		},
		Variables: make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(obj)
	return obj
}

// Global PROCESS class definitions - created once and reused
var processClassesOnce = sync.Once{}
var processClasses map[string]*environment.Class

func getProcessClasses() map[string]*environment.Class {
	processClassesOnce.Do(func() {
		processClasses = map[string]*environment.Class{
			"PIPE": {
				Name:          "PIPE",
				QualifiedName: "stdlib:PROCESS.PIPE",
				ModulePath:    "stdlib:PROCESS",
				ParentClasses: []string{"stdlib:IO.READER", "stdlib:IO.WRITER"}, // Can act as both
				MRO:           []string{"stdlib:PROCESS.PIPE", "stdlib:IO.READWRITER", "stdlib:IO.READER", "stdlib:IO.WRITER"},
				Documentation: []string{
					"A communication pipe connected to a child process's standard input, output, or error streams.",
					"Provides read/write access to process streams and inherits from IO.READER and IO.WRITER interfaces.",
					"",
					"@class PIPE",
					"@inherits IO.READER, IO.WRITER",
					"@example Reading from process stdout",
					"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
					"CMD DO PUSH WIT \"echo\"",
					"CMD DO PUSH WIT \"Hello, World!\"",
					"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
					"PROC DO START",
					"I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 1024",
					"SAYZ WIT OUTPUT",
					"PROC DO WAIT",
					"@example Writing to process stdin",
					"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
					"CMD DO PUSH WIT \"cat\"",
					"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
					"PROC DO START",
					"PROC STDIN DO WRITE WIT \"Hello from parent process!\"",
					"PROC STDIN DO CLOSE",
					"I HAS A VARIABLE RESULT TEH STRIN ITZ PROC STDOUT DO READ WIT 1024",
					"SAYZ WIT RESULT",
					"@example Handle process stderr",
					"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
					"PROC DO START",
					"I HAS A VARIABLE ERROR_OUTPUT TEH STRIN ITZ PROC STDERR DO READ WIT 512",
					"IZ NO SAEM AS (ERROR_OUTPUT SAEM AS \"\")?",
					"    SAYZ WIT \"Process error: \"",
					"    SAYZ WIT ERROR_OUTPUT",
					"KTHX",
					"@note Created automatically when MINION.START is called",
					"@note STDIN pipes support WRITE operations, STDOUT/STDERR pipes support READ operations",
					"@note Pipes are automatically closed when the process terminates",
					"@see MINION",
				},
				PublicFunctions: map[string]*environment.Function{
					// READ method (for stdout/stderr pipes)
					"READ": {
						Name:       "READ",
						ReturnType: "STRIN",
						Parameters: []environment.Parameter{
							{Name: "size", Type: "INTEGR"},
						},
						Documentation: []string{
							"Reads up to the specified number of characters from stdout or stderr pipes.",
							"Blocks until data is available or the pipe reaches end-of-file.",
							"",
							"@syntax <pipe> DO READ WIT <size>",
							"@param {INTEGR} size - Maximum number of characters to read",
							"@returns {STRIN} The data read from the pipe",
							"@example Read process output",
							"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
							"CMD DO PUSH WIT \"ls\"",
							"CMD DO PUSH WIT \"-la\"",
							"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
							"PROC DO START",
							"I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 4096",
							"SAYZ WIT \"Directory listing:\"",
							"SAYZ WIT OUTPUT",
							"@example Read in chunks",
							"I HAS A VARIABLE BUFFER TEH STRIN ITZ \"\"",
							"WHILE NO SAEM AS (PROC STDOUT IS_EOF)",
							"    I HAS A VARIABLE CHUNK TEH STRIN ITZ PROC STDOUT DO READ WIT 256",
							"    IZ CHUNK SAEM AS \"\"?",
							"        GTFO BTW End of output",
							"    KTHX",
							"    BUFFER ITZ BUFFER MOAR CHUNK",
							"KTHX",
							"@example Read stderr separately",
							"I HAS A VARIABLE ERROR_DATA TEH STRIN ITZ PROC STDERR DO READ WIT 1024",
							"IZ NO SAEM AS (ERROR_DATA SAEM AS \"\")?",
							"    SAYZ WIT \"Process errors: \"",
							"    SAYZ WIT ERROR_DATA",
							"KTHX",
							"@note Only works on STDOUT and STDERR pipes, not STDIN",
							"@note Returns empty string when end-of-file is reached",
							"@note May return fewer characters than requested if less data is available",
							"@see WRITE, IS_EOF, IS_OPEN",
							"@category pipe-operations",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							size := args[0]

							sizeVal, ok := size.(environment.IntegerValue)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Errorf("READ expects INTEGR size, got %s", size.Type()).Error()}
							}

							pipeData, ok := this.NativeData.(*PipeData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "READ: invalid context"}
							}

							if !pipeData.IsOpen {
								return environment.NOTHIN, runtime.Exception{Message: "READ: pipe is not open"}
							}

							if pipeData.FDType == "STDIN" {
								return environment.NOTHIN, runtime.Exception{Message: "READ: cannot read from stdin pipe"}
							}

							if pipeData.IsEOF {
								return environment.StringValue(""), nil
							}

							readSize := int(sizeVal)
							if readSize <= 0 {
								return environment.StringValue(""), nil
							}

							buffer := make([]byte, readSize)
							reader, ok := pipeData.Pipe.(io.Reader)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Errorf("READ: pipe is not readable").Error()}
							}

							n, err := reader.Read(buffer)
							if err != nil {
								if err == io.EOF {
									pipeData.IsEOF = true
								} else {
									return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("READ: %v", err)}
								}
							}

							return environment.StringValue(string(buffer[:n])), nil
						},
					},
					// WRITE method (for stdin pipe)
					"WRITE": {
						Name:       "WRITE",
						ReturnType: "INTEGR",
						Parameters: []environment.Parameter{
							{Name: "data", Type: "STRIN"},
						},
						Documentation: []string{
							"Writes string data to the process's stdin pipe.",
							"Returns the number of bytes written, only works on STDIN pipes.",
							"",
							"@syntax <pipe> DO WRITE WIT <data>",
							"@param {STRIN} data - The string data to send to the process",
							"@returns {INTEGR} Number of bytes written",
							"@example Send input to interactive process",
							"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
							"CMD DO PUSH WIT \"python3\"",
							"CMD DO PUSH WIT \"-i\"",
							"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
							"PROC DO START",
							"PROC STDIN DO WRITE WIT \"print('Hello from Python!')\\n\"",
							"PROC STDIN DO WRITE WIT \"exit()\\n\"",
							"I HAS A VARIABLE RESULT TEH STRIN ITZ PROC STDOUT DO READ WIT 1024",
							"SAYZ WIT RESULT",
							"@example Send data to filter process",
							"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
							"CMD DO PUSH WIT \"grep\"",
							"CMD DO PUSH WIT \"error\"",
							"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
							"PROC DO START",
							"PROC STDIN DO WRITE WIT \"info: starting process\\n\"",
							"PROC STDIN DO WRITE WIT \"error: something went wrong\\n\"",
							"PROC STDIN DO WRITE WIT \"info: process completed\\n\"",
							"PROC STDIN DO CLOSE BTW Signal end of input",
							"I HAS A VARIABLE FILTERED TEH STRIN ITZ PROC STDOUT DO READ WIT 1024",
							"SAYZ WIT \"Filtered output: \"",
							"SAYZ WIT FILTERED",
							"@example Pipe data between processes",
							"PROC STDIN DO WRITE WIT \"line 1\\n\"",
							"PROC STDIN DO WRITE WIT \"line 2\\n\"",
							"I HAS A VARIABLE BYTES_WRITTEN TEH INTEGR ITZ PROC STDIN DO WRITE WIT \"line 3\\n\"",
							"SAYZ WIT \"Wrote \"",
							"SAYZ WIT BYTES_WRITTEN",
							"SAYZ WIT \" bytes\"",
							"@note Only works on STDIN pipes, throws exception on STDOUT/STDERR pipes",
							"@note Process must be started before writing to pipes",
							"@note Close stdin pipe when done writing to signal end of input",
							"@see READ, CLOSE, IS_OPEN",
							"@category pipe-operations",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							data := args[0]

							dataVal, ok := data.(environment.StringValue)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Errorf("WRITE expects STRIN data, got %s", data.Type()).Error()}
							}

							pipeData, ok := this.NativeData.(*PipeData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "WRITE: invalid context"}
							}

							if !pipeData.IsOpen {
								return environment.NOTHIN, runtime.Exception{Message: "WRITE: pipe is not open"}
							}

							if pipeData.FDType != "STDIN" {
								return environment.NOTHIN, runtime.Exception{Message: "WRITE: cannot write to stdout/stderr pipe"}
							}

							writer, ok := pipeData.Pipe.(io.Writer)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "WRITE: pipe is not writable"}
							}

							dataStr := string(dataVal)
							n, err := writer.Write([]byte(dataStr))
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("WRITE: %v", err)}
							}

							return environment.IntegerValue(n), nil
						},
					},
					// CLOSE method
					"CLOSE": {
						Name: "CLOSE",
						Documentation: []string{
							"Closes the pipe connection and releases associated resources.",
							"For stdin pipes, signals end-of-input to the child process.",
							"",
							"@syntax <pipe> DO CLOSE",
							"@returns {NOTHIN} No return value",
							"@example Close stdin to signal completion",
							"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
							"CMD DO PUSH WIT \"wc\"",
							"CMD DO PUSH WIT \"-l\"",
							"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
							"PROC DO START",
							"PROC STDIN DO WRITE WIT \"line 1\\n\"",
							"PROC STDIN DO WRITE WIT \"line 2\\n\"",
							"PROC STDIN DO WRITE WIT \"line 3\\n\"",
							"PROC STDIN DO CLOSE BTW Signal end of input",
							"I HAS A VARIABLE COUNT TEH STRIN ITZ PROC STDOUT DO READ WIT 100",
							"SAYZ WIT \"Line count: \"",
							"SAYZ WIT COUNT",
							"@example Cleanup pipes after process",
							"PROC DO WAIT",
							"PROC STDOUT DO CLOSE",
							"PROC STDERR DO CLOSE",
							"SAYZ WIT \"All pipes closed\"",
							"@example Close in error handling",
							"MAYB",
							"    PROC STDIN DO WRITE WIT \"some input\"",
							"    BTW More operations here",
							"OOPSIE ERR",
							"    SAYZ WIT \"Error occurred: \"",
							"    SAYZ WIT ERR",
							"FINALLY",
							"    IZ PROC STDIN IS_OPEN?",
							"        PROC STDIN DO CLOSE",
							"    KTHX",
							"KTHX",
							"@note Safe to call multiple times - no error if already closed",
							"@note For stdin pipes: signals EOF to child process",
							"@note For stdout/stderr pipes: prevents further reading",
							"@note Pipes are automatically closed when process terminates",
							"@see IS_OPEN, WRITE, READ",
							"@category pipe-operations",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							pipeData, ok := this.NativeData.(*PipeData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "CLOSE: invalid context"}
							}

							if !pipeData.IsOpen {
								return environment.NOTHIN, nil // Already closed
							}

							if pipeData.Pipe != nil {
								err := pipeData.Pipe.Close()
								if err != nil {
									return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("CLOSE: %v", err)}
								}
							}

							pipeData.IsOpen = false

							return environment.NOTHIN, nil
						},
					},
				},
				PublicVariables: map[string]*environment.MemberVariable{
					"FD_TYPE": {
						Variable: environment.Variable{
							Name:     "FD_TYPE",
							Type:     "STRIN",
							Value:    environment.StringValue(""),
							IsLocked: true,
							IsPublic: true,
							Documentation: []string{
								"Read-only property indicating the type of pipe connection.",
								"",
								"@property {STRIN} FD_TYPE - Pipe type (\"STDIN\", \"STDOUT\", or \"STDERR\")",
								"@example Check pipe type for appropriate operations",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC DO START",
								"IZ (PROC STDIN FD_TYPE) SAEM AS \"STDIN\"?",
								"    PROC STDIN DO WRITE WIT \"input data\"",
								"    PROC STDIN DO CLOSE",
								"KTHX",
								"@example Handle different pipe types",
								"I HAS A VARIABLE PIPES TEH BUKKIT ITZ NEW BUKKIT",
								"PIPES DO PUSH WIT PROC STDIN",
								"PIPES DO PUSH WIT PROC STDOUT",
								"PIPES DO PUSH WIT PROC STDERR",
								"IM OUTTA UR PIPES NERFIN PIPE",
								"    I HAS A VARIABLE TYPE TEH STRIN ITZ PIPE FD_TYPE",
								"    IZ TYPE SAEM AS \"STDIN\"?",
								"        SAYZ WIT \"Input pipe found\"",
								"    NOPE",
								"        SAYZ WIT \"Output pipe: \"",
								"        SAYZ WIT TYPE",
								"    KTHX",
								"IM IN UR PIPES",
								"@note Set automatically when pipe is created",
								"@note STDIN pipes support write operations",
								"@note STDOUT and STDERR pipes support read operations",
								"@see IS_OPEN, IS_EOF",
								"@category pipe-properties",
							},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							pipeData, ok := this.NativeData.(*PipeData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "FD_TYPE: invalid context"}
							}
							return environment.StringValue(pipeData.FDType), nil
						},
						NativeSet: nil, // Read-only
					},
					"IS_OPEN": {
						Variable: environment.Variable{
							Name:     "IS_OPEN",
							Type:     "BOOL",
							Value:    environment.NO,
							IsLocked: true,
							IsPublic: true,
							Documentation: []string{
								"Read-only property indicating whether the pipe is open for I/O operations.",
								"",
								"@property {BOOL} IS_OPEN - YEZ if pipe is open, NO if closed",
								"@example Check if pipe is available before operations",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC DO START",
								"IZ PROC STDIN IS_OPEN?",
								"    PROC STDIN DO WRITE WIT \"data\"",
								"    PROC STDIN DO CLOSE",
								"NOPE",
								"    SAYZ WIT \"Stdin pipe is not available\"",
								"KTHX",
								"@example Monitor pipe status during operations",
								"WHILE (PROC STDOUT IS_OPEN)",
								"    I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 256",
								"    IZ OUTPUT SAEM AS \"\"?",
								"        GTFO BTW End of output reached",
								"    KTHX",
								"    SAYZ WIT OUTPUT",
								"KTHX",
								"@example Safe pipe cleanup",
								"IZ PROC STDERR IS_OPEN?",
								"    I HAS A VARIABLE ERRORS TEH STRIN ITZ PROC STDERR DO READ WIT 1024",
								"    PROC STDERR DO CLOSE",
								"    IZ NO SAEM AS (ERRORS SAEM AS \"\")?",
								"        SAYZ WIT \"Process errors: \"",
								"        SAYZ WIT ERRORS",
								"    KTHX",
								"KTHX",
								"@note Automatically set to YEZ when process starts",
								"@note Becomes NO when pipe is closed or process terminates",
								"@note Use to avoid exceptions from operations on closed pipes",
								"@see FD_TYPE, IS_EOF, CLOSE",
								"@category pipe-properties",
							},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							pipeData, ok := this.NativeData.(*PipeData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "IS_OPEN: invalid context"}
							}
							return environment.BoolValue(pipeData.IsOpen), nil
						},
						NativeSet: nil, // Read-only
					},
					"IS_EOF": {
						Variable: environment.Variable{
							Name:     "IS_EOF",
							Type:     "BOOL",
							Value:    environment.NO,
							IsLocked: true,
							IsPublic: true,
							Documentation: []string{
								"Read-only property indicating whether end-of-file has been reached on read pipes.",
								"",
								"@property {BOOL} IS_EOF - YEZ if end-of-file reached, NO otherwise",
								"@example Read until EOF",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC DO START",
								"I HAS A VARIABLE ALL_OUTPUT TEH STRIN ITZ \"\"",
								"WHILE NO SAEM AS (PROC STDOUT IS_EOF)",
								"    I HAS A VARIABLE CHUNK TEH STRIN ITZ PROC STDOUT DO READ WIT 512",
								"    IZ CHUNK SAEM AS \"\"?",
								"        GTFO BTW No more data available",
								"    KTHX",
								"    ALL_OUTPUT ITZ ALL_OUTPUT MOAR CHUNK",
								"KTHX",
								"SAYZ WIT \"Complete output: \"",
								"SAYZ WIT ALL_OUTPUT",
								"@example Check EOF status after reading",
								"I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 1024",
								"IZ PROC STDOUT IS_EOF?",
								"    SAYZ WIT \"Reached end of process output\"",
								"NOPE",
								"    SAYZ WIT \"More output may be available\"",
								"KTHX",
								"@example Handle both stdout and stderr EOF",
								"WHILE NO SAEM AS ((PROC STDOUT IS_EOF) AN (PROC STDERR IS_EOF))",
								"    IZ NO SAEM AS (PROC STDOUT IS_EOF)?",
								"        I HAS A VARIABLE OUT TEH STRIN ITZ PROC STDOUT DO READ WIT 256",
								"        IZ NO SAEM AS (OUT SAEM AS \"\")?",
								"            SAYZ WIT \"STDOUT: \"",
								"            SAYZ WIT OUT",
								"        KTHX",
								"    KTHX",
								"    IZ NO SAEM AS (PROC STDERR IS_EOF)?",
								"        I HAS A VARIABLE ERR TEH STRIN ITZ PROC STDERR DO READ WIT 256",
								"        IZ NO SAEM AS (ERR SAEM AS \"\")?",
								"            SAYZ WIT \"STDERR: \"",
								"            SAYZ WIT ERR",
								"        KTHX",
								"    KTHX",
								"KTHX",
								"@note Only relevant for STDOUT and STDERR pipes, not STDIN",
								"@note Automatically set when READ operation encounters EOF",
								"@note Once EOF is reached, further READ calls return empty string",
								"@see IS_OPEN, READ, FD_TYPE",
								"@category pipe-properties",
							},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							pipeData, ok := this.NativeData.(*PipeData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "IS_EOF: invalid context"}
							}
							return environment.BoolValue(pipeData.IsEOF), nil
						},
						NativeSet: nil, // Read-only
					},
				},
				PrivateVariables: make(map[string]*environment.MemberVariable),
				PrivateFunctions: make(map[string]*environment.Function),
				SharedVariables:  make(map[string]*environment.MemberVariable),
				SharedFunctions:  make(map[string]*environment.Function),
			},
			"MINION": {
				Name:          "MINION",
				QualifiedName: "stdlib:PROCESS.MINION",
				ModulePath:    "stdlib:PROCESS",
				ParentClasses: []string{},
				MRO:           []string{"stdlib:PROCESS.MINION"},
				Documentation: []string{
					"A child process that can be launched, monitored, and controlled with full I/O access.",
					"Provides comprehensive process management including environment variables, working directory, and signal handling.",
					"",
					"@class MINION",
					"@example Basic command execution",
					"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
					"CMD DO PUSH WIT \"echo\"",
					"CMD DO PUSH WIT \"Hello, World!\"",
					"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
					"PROC DO START",
					"I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 1024",
					"I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC DO WAIT",
					"SAYZ WIT \"Output: \"",
					"SAYZ WIT OUTPUT",
					"SAYZ WIT \"Exit code: \"",
					"SAYZ WIT EXIT_CODE",
					"@example Interactive process with stdin/stdout",
					"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
					"CMD DO PUSH WIT \"python3\"",
					"CMD DO PUSH WIT \"-i\"",
					"I HAS A VARIABLE PYTHON TEH MINION ITZ NEW MINION WIT CMD",
					"PYTHON DO START",
					"PYTHON STDIN DO WRITE WIT \"print(2 + 2)\\n\"",
					"PYTHON STDIN DO WRITE WIT \"exit()\\n\"",
					"I HAS A VARIABLE RESULT TEH STRIN ITZ PYTHON STDOUT DO READ WIT 1024",
					"PYTHON DO WAIT",
					"SAYZ WIT \"Python result: \"",
					"SAYZ WIT RESULT",
					"@example Process with custom environment and working directory",
					"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
					"CMD DO PUSH WIT \"printenv\"",
					"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
					"PROC WORKDIR ITZ \"/tmp\"",
					"PROC ENV DO PUT WIT \"CUSTOM_VAR\" AN WIT \"Hello from environment!\"",
					"PROC DO START",
					"I HAS A VARIABLE ENV_OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 2048",
					"PROC DO WAIT",
					"SAYZ WIT ENV_OUTPUT",
					"@note Automatically creates STDIN, STDOUT, and STDERR pipes when started",
					"@note Supports environment variable customization and working directory changes",
					"@note Provides process lifecycle management (start, wait, kill, signal)",
					"@see PIPE",
				},
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"MINION": {
						Name: "MINION",
						Parameters: []environment.Parameter{
							{Name: "cmdline", Type: "BUKKIT"},
						},
						Documentation: []string{
							"Creates a new process definition with command and arguments.",
							"Initializes with current working directory and copies current environment variables.",
							"",
							"@syntax NEW MINION WIT <cmdline>",
							"@param {BUKKIT} cmdline - Command and arguments as array of strings",
							"@returns {NOTHIN} No return value (constructor)",
							"@example Create process for simple command",
							"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
							"CMD DO PUSH WIT \"ls\"",
							"CMD DO PUSH WIT \"-la\"",
							"CMD DO PUSH WIT \"/home\"",
							"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
							"BTW Process created but not started yet",
							"@example Create process for complex command with pipes",
							"I HAS A VARIABLE GREP_CMD TEH BUKKIT ITZ NEW BUKKIT",
							"GREP_CMD DO PUSH WIT \"grep\"",
							"GREP_CMD DO PUSH WIT \"-n\"",
							"GREP_CMD DO PUSH WIT \"error\"",
							"I HAS A VARIABLE GREP_PROC TEH MINION ITZ NEW MINION WIT GREP_CMD",
							"@example Create interactive shell",
							"I HAS A VARIABLE SHELL_CMD TEH BUKKIT ITZ NEW BUKKIT",
							"SHELL_CMD DO PUSH WIT \"bash\"",
							"SHELL_CMD DO PUSH WIT \"--login\"",
							"I HAS A VARIABLE SHELL TEH MINION ITZ NEW MINION WIT SHELL_CMD",
							"@note First element in BUKKIT must be the executable name or path",
							"@note All elements must be convertible to strings",
							"@note Process inherits current environment and working directory",
							"@note Use START method to actually launch the process",
							"@see START, CMDLINE, WORKDIR, ENV",
							"@category process-creation",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							cmdline := args[0]

							// Validate that the argument is a BUKKIT
							cmdlineInstance, ok := cmdline.(*environment.ObjectInstance)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("MINION constructor expects BUKKIT cmdline, got %s", cmdline.Type())}
							}

							// Extract command line from BUKKIT
							cmdLineSlice := []string{}
							if nativeData, ok := cmdlineInstance.NativeData.(BukkitSlice); ok {
								for _, val := range nativeData {
									if strVal, ok := val.(environment.StringValue); ok {
										cmdLineSlice = append(cmdLineSlice, string(strVal))
									} else {
										return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("MINION constructor: all cmdline elements must be STRING, got %s", val.Type())}
									}
								}
							} else {
								return environment.NOTHIN, runtime.Exception{Message: "MINION constructor: invalid BUKKIT data"}
							}

							if len(cmdLineSlice) == 0 {
								return environment.NOTHIN, runtime.Exception{Message: "MINION: command line cannot be empty"}
							}

							// Convert command line to BUKKIT object for storage
							cmdlineObj := NewBukkitInstance()
							cmdlineObjSlice := cmdlineObj.NativeData.(BukkitSlice)
							for _, arg := range cmdLineSlice {
								cmdlineObjSlice = append(cmdlineObjSlice, environment.StringValue(arg))
							}
							cmdlineObj.NativeData = cmdlineObjSlice

							// Initialize environment variable (BASKIT)
							baskitEnv := NewBaskitInstance()
							baskitMap := baskitEnv.NativeData.(BaskitMap)
							environ := os.Environ()
							for _, envVar := range environ {
								parts := strings.SplitN(envVar, "=", 2)
								if len(parts) == 2 {
									key := parts[0]
									value := parts[1]
									baskitMap[key] = environment.StringValue(value)
								}
							}
							baskitEnv.NativeData = baskitMap

							// Initialize the minion data
							minionData := &MinionData{
								CmdLine:  cmdlineObj,
								Env:      baskitEnv,
								WorkDir:  "", // Use current directory by default
								Running:  false,
								Finished: false,
								ExitCode: 0,
								PID:      -1,
							}
							this.NativeData = minionData

							return environment.NOTHIN, nil
						},
					},
					// START method
					"START": {
						Name: "START",
						Documentation: []string{
							"Launches the child process and creates communication pipes.",
							"Creates STDIN, STDOUT, and STDERR pipes for process communication.",
							"",
							"@syntax <minion> DO START",
							"@returns {NOTHIN} No return value",
							"@example Start simple command",
							"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
							"CMD DO PUSH WIT \"date\"",
							"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
							"PROC DO START",
							"SAYZ WIT \"Process started with PID: \"",
							"SAYZ WIT PROC PID",
							"@example Start and immediately read output",
							"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
							"CMD DO PUSH WIT \"echo\"",
							"CMD DO PUSH WIT \"Hello!\"",
							"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
							"PROC DO START",
							"I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 1024",
							"SAYZ WIT \"Process output: \"",
							"SAYZ WIT OUTPUT",
							"@example Start process with custom environment",
							"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
							"PROC ENV DO PUT WIT \"MY_VAR\" AN WIT \"custom_value\"",
							"PROC WORKDIR ITZ \"/tmp\"",
							"PROC DO START",
							"SAYZ WIT \"Process started in directory: \"",
							"SAYZ WIT PROC WORKDIR",
							"@note Process must not be already running or finished",
							"@note After START, pipes become available for I/O operations",
							"@note RUNNING property becomes YEZ after successful start",
							"@note PID property is set to the actual process ID",
							"@see WAIT, STDIN, STDOUT, STDERR, PID",
							"@category process-lifecycle",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							minionData, ok := this.NativeData.(*MinionData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "START: invalid minion context"}
							}

							if minionData.Running {
								return environment.NOTHIN, runtime.Exception{Message: "START: process is already running"}
							}

							if minionData.Finished {
								return environment.NOTHIN, runtime.Exception{Message: "START: process has already finished"}
							}

							// Create the command
							cmdline := minionData.CmdLine.NativeData.(BukkitSlice)
							if len(cmdline) == 0 {
								return environment.NOTHIN, runtime.Exception{Message: "START: command line cannot be empty"}
							}

							// Convert BUKKIT to []string
							cmdArgs := []string{}
							for _, arg := range cmdline {
								if strArg, err := arg.Cast("STRIN"); err == nil {
									cmdArgs = append(cmdArgs, string(strArg.(environment.StringValue)))
								} else {
									return environment.NOTHIN, runtime.Exception{Message: "START: all command line arguments must be convertible to STRIN"}
								}
							}

							// Convert BASKIT to []string for environment
							envMap := minionData.Env.NativeData.(BaskitMap)
							envVars := []string{}
							for key, value := range envMap {
								if strVal, err := value.Cast("STRIN"); err == nil {
									envVars = append(envVars, fmt.Sprintf("%s=%s", key, strVal.(environment.StringValue)))
								} else {
									return environment.NOTHIN, runtime.Exception{Message: "START: all environment variable values must be convertible to STRIN"}
								}
							}

							cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
							if minionData.WorkDir != "" {
								cmd.Dir = minionData.WorkDir
							}
							cmd.Env = envVars

							// Create pipes
							stdin, err := cmd.StdinPipe()
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("Failed to create stdin pipe: %v", err)}
							}

							stdout, err := cmd.StdoutPipe()
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("Failed to create stdout pipe: %v", err)}
							}

							stderr, err := cmd.StderrPipe()
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("Failed to create stderr pipe: %v", err)}
							}

							// Start the process
							err = cmd.Start()
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("Failed to start process: %v", err)}
							}

							// Update minion data
							minionData.Process = cmd
							minionData.Running = true
							minionData.PID = cmd.Process.Pid

							// Create PIPE objects
							stdinPipe := NewPipeInstance("STDIN", cmd)
							stdinPipe.NativeData.(*PipeData).Pipe = stdin
							stdinPipe.NativeData.(*PipeData).IsOpen = true

							stdoutPipe := NewPipeInstance("STDOUT", cmd)
							stdoutPipe.NativeData.(*PipeData).Pipe = stdout
							stdoutPipe.NativeData.(*PipeData).IsOpen = true

							stderrPipe := NewPipeInstance("STDERR", cmd)
							stderrPipe.NativeData.(*PipeData).Pipe = stderr
							stderrPipe.NativeData.(*PipeData).IsOpen = true

							minionData.StdinPipe = stdinPipe
							minionData.StdoutPipe = stdoutPipe
							minionData.StderrPipe = stderrPipe

							return environment.NOTHIN, nil
						},
					},
					// WAIT method
					"WAIT": {
						Name:       "WAIT",
						ReturnType: "INTEGR",
						Documentation: []string{
							"Waits for the process to complete and returns the exit code.",
							"Blocks until the child process terminates, then returns its exit status.",
							"",
							"@syntax <minion> DO WAIT",
							"@returns {INTEGR} The process exit code (0 for success, non-zero for errors)",
							"@example Wait for process completion",
							"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
							"CMD DO PUSH WIT \"echo\"",
							"CMD DO PUSH WIT \"Hello, World!\"",
							"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
							"PROC DO START",
							"I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 1024",
							"I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC DO WAIT",
							"SAYZ WIT \"Process completed with exit code: \"",
							"SAYZ WIT EXIT_CODE",
							"@example Handle process errors",
							"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
							"CMD DO PUSH WIT \"false\"", // Command that always fails
							"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
							"PROC DO START",
							"I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC DO WAIT",
							"IZ EXIT_CODE SAEM AS 0?",
							"    SAYZ WIT \"Process succeeded\"",
							"NOPE",
							"    SAYZ WIT \"Process failed with code: \"",
							"    SAYZ WIT EXIT_CODE",
							"KTHX",
							"@example Wait with timeout pattern",
							"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
							"PROC DO START",
							"I HAS A VARIABLE TIMEOUT TEH INTEGR ITZ 5000", // 5 seconds in milliseconds
							"I HAS A VARIABLE START_TIME TEH INTEGR ITZ NOW",
							"WHILE (PROC RUNNING) AN ((NOW LES START_TIME) SAEM AS TIMEOUT)",
							"    I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC DO WAIT",
							"    GTFO BTW Process finished",
							"KTHX",
							"IZ PROC RUNNING?",
							"    PROC DO KILL",
							"    SAYZ WIT \"Process timed out and was killed\"",
							"KTHX",
							"@throws Exception if process has not been started",
							"@note Blocks the current thread until process completes",
							"@note Exit code 0 typically indicates success",
							"@note Negative exit codes may indicate process was killed",
							"@see START, RUNNING, FINISHED, EXIT_CODE",
							"@category process-lifecycle",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							minionData, ok := this.NativeData.(*MinionData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "WAIT: invalid minion context"}
							}

							if !minionData.Running && !minionData.Finished {
								return environment.NOTHIN, runtime.Exception{Message: "WAIT: process has not been started"}
							}

							if minionData.Finished {
								return environment.IntegerValue(minionData.ExitCode), nil
							}

							// Wait for the process to complete
							err := minionData.Process.Wait()
							minionData.Running = false
							minionData.Finished = true

							if err != nil {
								if exitError, ok := err.(*exec.ExitError); ok {
									minionData.ExitCode = exitError.ExitCode()
								} else {
									return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("WAIT: %v", err)}
								}
							} else {
								minionData.ExitCode = 0
							}

							return environment.IntegerValue(minionData.ExitCode), nil
						},
					},
					// KILL method
					"KILL": {
						Name: "KILL",
						Documentation: []string{
							"Forcefully terminates the running process.",
							"Sends SIGKILL signal to immediately stop the process without cleanup.",
							"",
							"@syntax <minion> DO KILL",
							"@returns {NOTHIN} No return value",
							"@example Kill a running process",
							"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
							"CMD DO PUSH WIT \"sleep\"",
							"CMD DO PUSH WIT \"30\"", // Sleep for 30 seconds
							"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
							"PROC DO START",
							"SAYZ WIT \"Process started, waiting 2 seconds before killing...\"",
							"I HAS A VARIABLE START_TIME TEH INTEGR ITZ NOW",
							"WHILE (NOW MINUSZ START_TIME) LIEKZ 2000",
							"    BTW Wait 2 seconds",
							"KTHX",
							"PROC DO KILL",
							"I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC DO WAIT",
							"SAYZ WIT \"Process killed, exit code: \"",
							"SAYZ WIT EXIT_CODE",
							"@example Kill process on timeout",
							"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
							"PROC DO START",
							"I HAS A VARIABLE TIMEOUT TEH INTEGR ITZ 10000", // 10 seconds
							"I HAS A VARIABLE START_TIME TEH INTEGR ITZ NOW",
							"WHILE (PROC RUNNING) AN ((NOW MINUSZ START_TIME) LIEKZ TIMEOUT)",
							"    I HAS A VARIABLE CHUNK TEH STRIN ITZ PROC STDOUT DO READ WIT 256",
							"    IZ CHUNK SAEM AS \"\"?",
							"        GTFO BTW No more output",
							"    KTHX",
							"    SAYZ WIT CHUNK",
							"KTHX",
							"IZ PROC RUNNING?",
							"    PROC DO KILL",
							"    SAYZ WIT \"Process timed out and was terminated\"",
							"NOPE",
							"    SAYZ WIT \"Process completed normally\"",
							"KTHX",
							"@example Cleanup in error handling",
							"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
							"PROC DO START",
							"MAYB",
							"    BTW Do some work with the process",
							"    I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 1024",
							"    SAYZ WIT OUTPUT",
							"OOPSIE ERR",
							"    SAYZ WIT \"Error occurred: \"",
							"    SAYZ WIT ERR",
							"    IZ PROC RUNNING?",
							"        PROC DO KILL",
							"        SAYZ WIT \"Process was killed due to error\"",
							"    KTHX",
							"KTHX",
							"@throws Exception if process is not running",
							"@note Immediately terminates process without cleanup",
							"@note Process exit code will be -1 after killing",
							"@note May leave child processes running if not handled properly",
							"@note Use WAIT after KILL to ensure process cleanup",
							"@see SIGNAL, WAIT, RUNNING",
							"@category process-control",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							minionData, ok := this.NativeData.(*MinionData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "KILL: invalid minion context"}
							}

							if !minionData.Running {
								return environment.NOTHIN, runtime.Exception{Message: "KILL: process is not running"}
							}

							err := minionData.Process.Process.Kill()
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("KILL: %v", err)}
							}

							// Wait for process to finish after killing
							minionData.Process.Wait()
							minionData.Running = false
							minionData.Finished = true
							minionData.ExitCode = -1 // Indicate killed

							return environment.NOTHIN, nil
						},
					},
					// SIGNAL method
					"SIGNAL": {
						Name: "SIGNAL",
						Parameters: []environment.Parameter{
							{Name: "code", Type: "INTEGR"},
						},
						Documentation: []string{
							"Sends a signal to the running process (Unix/Linux systems).",
							"Allows sending specific signals like SIGTERM, SIGINT, etc. to processes.",
							"",
							"@syntax <minion> DO SIGNAL WIT <code>",
							"@param {INTEGR} code - Signal number (e.g., 15 for SIGTERM, 2 for SIGINT)",
							"@returns {NOTHIN} No return value",
							"@example Send SIGTERM to gracefully stop process",
							"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
							"CMD DO PUSH WIT \"sleep\"",
							"CMD DO PUSH WIT \"30\"",
							"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
							"PROC DO START",
							"SAYZ WIT \"Sending SIGTERM (15) to process...\"",
							"PROC DO SIGNAL WIT 15", // SIGTERM
							"I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC DO WAIT",
							"SAYZ WIT \"Process terminated with signal, exit code: \"",
							"SAYZ WIT EXIT_CODE",
							"@example Send SIGINT (Ctrl+C equivalent)",
							"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
							"PROC DO START",
							"SAYZ WIT \"Sending SIGINT (2) to process...\"",
							"PROC DO SIGNAL WIT 2", // SIGINT
							"I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC DO WAIT",
							"SAYZ WIT \"Process interrupted, exit code: \"",
							"SAYZ WIT EXIT_CODE",
							"@example Graceful shutdown with timeout",
							"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
							"PROC DO START",
							"SAYZ WIT \"Attempting graceful shutdown...\"",
							"PROC DO SIGNAL WIT 15", // SIGTERM first
							"I HAS A VARIABLE START_TIME TEH INTEGR ITZ NOW",
							"WHILE (PROC RUNNING) AN ((NOW MINUSZ START_TIME) LIEKZ 5000)",
							"    BTW Wait up to 5 seconds for graceful shutdown",
							"KTHX",
							"IZ PROC RUNNING?",
							"    SAYZ WIT \"Process didn't respond to SIGTERM, killing...\"",
							"    PROC DO KILL",
							"NOPE",
							"    SAYZ WIT \"Process shut down gracefully\"",
							"KTHX",
							"@throws Exception if process is not running",
							"@throws Exception if signal code is invalid",
							"@throws Exception on Windows (not supported)",
							"@note Signal numbers vary by Unix system",
							"@note Common signals: 1=SIGHUP, 2=SIGINT, 9=SIGKILL, 15=SIGTERM",
							"@note Process may ignore some signals",
							"@note Not available on Windows systems",
							"@see KILL, WAIT, RUNNING",
							"@category process-control",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							code := args[0]

							if goRuntime.GOOS == "windows" {
								return environment.NOTHIN, runtime.Exception{Message: "SIGNAL: not supported on Windows"}
							}

							codeVal, ok := code.(environment.IntegerValue)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("SIGNAL expects INTEGR code, got %s", code.Type())}
							}

							minionData, ok := this.NativeData.(*MinionData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "SIGNAL: invalid minion context"}
							}

							if !minionData.Running {
								return environment.NOTHIN, runtime.Exception{Message: "SIGNAL: process is not running"}
							}

							signal := syscall.Signal(int(codeVal))
							err := minionData.Process.Process.Signal(signal)
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("SIGNAL: %v", err)}
							}

							return environment.NOTHIN, nil
						},
					},
				},
				PublicVariables: map[string]*environment.MemberVariable{
					"CMDLINE": {
						Variable: environment.Variable{
							Name:     "CMDLINE",
							Type:     "BUKKIT",
							IsLocked: true,
							IsPublic: true,
							Documentation: []string{
								"Read-only property containing the command line arguments.",
								"Contains the executable name and arguments passed to the process.",
								"",
								"@property {BUKKIT} CMDLINE - Command and arguments as array of strings",
								"@example Access command line arguments",
								"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
								"CMD DO PUSH WIT \"ls\"",
								"CMD DO PUSH WIT \"-la\"",
								"CMD DO PUSH WIT \"/home\"",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"I HAS A VARIABLE ARGS TEH BUKKIT ITZ PROC CMDLINE",
								"SAYZ WIT \"Command: \"",
								"SAYZ WIT ARGS 0",
								"SAYZ WIT \"Arguments: \"",
								"IM OUTTA UR ARGS NERFIN ARG",
								"    SAYZ WIT ARG",
								"IM IN UR ARGS",
								"@example Verify command before starting",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"I HAS A VARIABLE CMD_ARGS TEH BUKKIT ITZ PROC CMDLINE",
								"IZ (CMD_ARGS LENGTH) BIGGR THAN 0?",
								"    SAYZ WIT \"Will execute: \"",
								"    SAYZ WIT CMD_ARGS 0",
								"    PROC DO START",
								"NOPE",
								"    SAYZ WIT \"No command specified\"",
								"KTHX",
								"@example Log process command for debugging",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"SAYZ WIT \"Starting process with command: \"",
								"I HAS A VARIABLE CMD_STR TEH STRIN ITZ \"\"",
								"I HAS A VARIABLE ARGS TEH BUKKIT ITZ PROC CMDLINE",
								"IM OUTTA UR ARGS NERFIN ARG",
								"    IZ CMD_STR SAEM AS \"\"?",
								"        CMD_STR ITZ ARG",
								"    NOPE",
								"        CMD_STR ITZ CMD_STR MOAR \" \" MOAR ARG",
								"    KTHX",
								"IM IN UR ARGS",
								"SAYZ WIT CMD_STR",
								"@note First element is the executable name or path",
								"@note Remaining elements are command arguments",
								"@note Cannot be modified after process creation",
								"@note Arguments are stored as strings",
								"@see START, WORKDIR, ENV",
								"@category process-configuration",
							},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							minionData, ok := this.NativeData.(*MinionData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "CMDLINE: invalid minion context"}
							}
							return minionData.CmdLine, nil
						},
						NativeSet: nil, // Read-only
					},
					"WORKDIR": {
						Variable: environment.Variable{
							Name:     "WORKDIR",
							Type:     "STRIN",
							IsLocked: false,
							IsPublic: true,
							Documentation: []string{
								"Working directory for the process.",
								"Directory where the process will be executed. Can be changed before starting.",
								"",
								"@property {STRIN} WORKDIR - Absolute or relative path to working directory",
								"@example Set working directory before starting",
								"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
								"CMD DO PUSH WIT \"ls\"",
								"CMD DO PUSH WIT \"-la\"",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC WORKDIR ITZ \"/tmp\"",
								"PROC DO START",
								"I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 1024",
								"SAYZ WIT \"Contents of /tmp:\"",
								"SAYZ WIT OUTPUT",
								"@example Use relative path",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC WORKDIR ITZ \"../parent_directory\"",
								"PROC DO START",
								"SAYZ WIT \"Process started in relative directory\"",
								"@example Change working directory dynamically",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC WORKDIR ITZ \"/home/user\"",
								"SAYZ WIT \"Initial workdir: \"",
								"SAYZ WIT PROC WORKDIR",
								"PROC WORKDIR ITZ \"/var/log\"",
								"SAYZ WIT \"Changed workdir: \"",
								"SAYZ WIT PROC WORKDIR",
								"@example Handle working directory errors",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC WORKDIR ITZ \"/nonexistent/directory\"",
								"MAYB",
								"    PROC DO START",
								"    SAYZ WIT \"Process started successfully\"",
								"OOPSIE ERR",
								"    SAYZ WIT \"Failed to start process: \"",
								"    SAYZ WIT ERR",
								"    PROC WORKDIR ITZ \"/tmp\"", // Fallback to safe directory
								"    PROC DO START",
								"KTHX",
								"@note Must be set before calling START",
								"@note Can be absolute or relative path",
								"@note Directory must exist and be accessible",
								"@note Cannot be changed while process is running",
								"@note Defaults to current working directory if not set",
								"@see START, CMDLINE, ENV",
								"@category process-configuration",
							},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							minionData, ok := this.NativeData.(*MinionData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "WORKDIR: invalid minion context"}
							}
							return environment.StringValue(minionData.WorkDir), nil
						},
						NativeSet: func(this *environment.ObjectInstance, newValue environment.Value) error {
							minionData, ok := this.NativeData.(*MinionData)
							if !ok {
								return runtime.Exception{Message: "WORKDIR: invalid minion context"}
							}
							if minionData.Running {
								return runtime.Exception{Message: "WORKDIR: cannot change working directory while process is running"}
							}
							strVal, err := newValue.Cast("STRIN")
							if err != nil {
								return runtime.Exception{Message: fmt.Sprintf("WORKDIR expects STRIN value, got %s", newValue.Type())}
							}
							minionData.WorkDir = string(strVal.(environment.StringValue))
							return nil
						},
					},
					"ENV": {
						Variable: environment.Variable{
							Name:     "ENV",
							Type:     "BASKIT",
							IsLocked: false,
							IsPublic: true,
							Documentation: []string{
								"Environment variables for the process.",
								"Key-value pairs that will be set as environment variables for the child process.",
								"",
								"@property {BASKIT} ENV - Environment variables as key-value map",
								"@example Set custom environment variables",
								"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
								"CMD DO PUSH WIT \"printenv\"",
								"CMD DO PUSH WIT \"MY_VAR\"",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC ENV DO PUT WIT \"MY_VAR\" AN WIT \"Hello from environment!\"",
								"PROC ENV DO PUT WIT \"ANOTHER_VAR\" AN WIT \"Another value\"",
								"PROC DO START",
								"I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 1024",
								"SAYZ WIT \"Environment variable value: \"",
								"SAYZ WIT OUTPUT",
								"@example Inherit and modify parent environment",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"I HAS A VARIABLE ENV TEH BASKIT ITZ PROC ENV",
								"ENV DO PUT WIT \"PATH\" AN WIT \"/custom/path:/usr/bin:/bin\"",
								"ENV DO PUT WIT \"HOME\" AN WIT \"/tmp\"",
								"PROC DO START",
								"SAYZ WIT \"Process started with modified environment\"",
								"@example Clear environment and set minimal vars",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"I HAS A VARIABLE NEW_ENV TEH BASKIT ITZ NEW BASKIT",
								"NEW_ENV DO PUT WIT \"PATH\" AN WIT \"/bin:/usr/bin\"",
								"NEW_ENV DO PUT WIT \"HOME\" AN WIT \"/tmp\"",
								"PROC ENV ITZ NEW_ENV",
								"PROC DO START",
								"SAYZ WIT \"Process started with clean environment\"",
								"@example Handle environment variable errors",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC ENV DO PUT WIT \"INVALID_KEY\" AN WIT 123", // Wrong type
								"MAYB",
								"    PROC DO START",
								"OOPSIE ERR",
								"    SAYZ WIT \"Environment setup error: \"",
								"    SAYZ WIT ERR",
								"    PROC ENV DO PUT WIT \"INVALID_KEY\" AN WIT \"123\"", // Fix type
								"    PROC DO START",
								"KTHX",
								"@note Inherits parent process environment by default",
								"@note Must be set before calling START",
								"@note All values must be convertible to strings",
								"@note Cannot be changed while process is running",
								"@note Use BASKIT operations (PUT, GET, HAS) to modify",
								"@see START, CMDLINE, WORKDIR",
								"@category process-configuration",
							},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							minionData, ok := this.NativeData.(*MinionData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "ENV: invalid minion context"}
							}
							return minionData.Env, nil
						},
						NativeSet: func(this *environment.ObjectInstance, newValue environment.Value) error {
							minionData, ok := this.NativeData.(*MinionData)
							if !ok {
								return runtime.Exception{Message: "ENV: invalid minion context"}
							}
							if minionData.Running {
								return runtime.Exception{Message: "ENV: cannot change environment while process is running"}
							}
							envInstance, err := newValue.Cast("BASKIT")
							if err != nil {
								return runtime.Exception{Message: fmt.Sprintf("ENV expects BASKIT value, got %s", newValue.Type())}
							}
							minionData.Env = envInstance.(*environment.ObjectInstance)
							return nil
						},
					},
					"RUNNING": {
						Variable: environment.Variable{
							Name:     "RUNNING",
							Type:     "BOOL",
							IsLocked: true,
							IsPublic: true,
							Documentation: []string{
								"Read-only property indicating whether the process is currently running.",
								"True from START until process termination, false otherwise.",
								"",
								"@property {BOOL} RUNNING - YEZ if process is running, NO if stopped",
								"@example Check if process is still running",
								"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
								"CMD DO PUSH WIT \"sleep\"",
								"CMD DO PUSH WIT \"5\"",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC DO START",
								"SAYZ WIT \"Process running: \"",
								"SAYZ WIT PROC RUNNING",
								"I HAS A VARIABLE START_TIME TEH INTEGR ITZ NOW",
								"WHILE (PROC RUNNING) AN ((NOW MINUSZ START_TIME) LIEKZ 3000)",
								"    BTW Wait up to 3 seconds",
								"KTHX",
								"SAYZ WIT \"Process still running: \"",
								"SAYZ WIT PROC RUNNING",
								"@example Wait for process to finish",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC DO START",
								"WHILE PROC RUNNING",
								"    I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 256",
								"    IZ NO SAEM AS (OUTPUT SAEM AS \"\")?",
								"        SAYZ WIT OUTPUT",
								"    KTHX",
								"    I HAS A VARIABLE SLEEP_TIME TEH INTEGR ITZ 100", // Sleep 100ms
								"    BTW Sleep implementation would go here",
								"KTHX",
								"SAYZ WIT \"Process has finished\"",
								"@example Monitor process lifecycle",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"SAYZ WIT \"Before start - Running: \"",
								"SAYZ WIT PROC RUNNING",
								"PROC DO START",
								"SAYZ WIT \"After start - Running: \"",
								"SAYZ WIT PROC RUNNING",
								"I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC DO WAIT",
								"SAYZ WIT \"After wait - Running: \"",
								"SAYZ WIT PROC RUNNING",
								"@example Handle process that exits quickly",
								"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
								"CMD DO PUSH WIT \"true\"", // Command that exits immediately
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC DO START",
								"IZ PROC RUNNING?",
								"    SAYZ WIT \"Process is running\"",
								"    I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC DO WAIT",
								"NOPE",
								"    SAYZ WIT \"Process already finished\"",
								"    SAYZ WIT \"Exit code: \"",
								"    SAYZ WIT PROC EXIT_CODE",
								"KTHX",
								"@note Becomes YEZ immediately after successful START",
								"@note Becomes NO when process terminates (normally or killed)",
								"@note Check after WAIT to confirm process has stopped",
								"@note Useful for polling or timeout implementations",
								"@see START, WAIT, FINISHED, EXIT_CODE",
								"@category process-lifecycle",
							},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							minionData, ok := this.NativeData.(*MinionData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "RUNNING: invalid context"}
							}
							return environment.BoolValue(minionData.Running), nil
						},
						NativeSet: nil, // Read-only
					},
					"FINISHED": {
						Variable: environment.Variable{
							Name:     "FINISHED",
							Type:     "BOOL",
							IsLocked: true,
							IsPublic: true,
							Documentation: []string{
								"Read-only property indicating whether the process has completed.",
								"True after process termination, remains true until process is restarted.",
								"",
								"@property {BOOL} FINISHED - YEZ if process has finished, NO if still running or not started",
								"@example Check process completion status",
								"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
								"CMD DO PUSH WIT \"echo\"",
								"CMD DO PUSH WIT \"Hello\"",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"SAYZ WIT \"Before start - Finished: \"",
								"SAYZ WIT PROC FINISHED",
								"PROC DO START",
								"SAYZ WIT \"After start - Finished: \"",
								"SAYZ WIT PROC FINISHED",
								"I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC DO WAIT",
								"SAYZ WIT \"After wait - Finished: \"",
								"SAYZ WIT PROC FINISHED",
								"@example Poll for completion",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC DO START",
								"WHILE NO SAEM AS (PROC FINISHED)",
								"    I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 256",
								"    IZ NO SAEM AS (OUTPUT SAEM AS \"\")?",
								"        SAYZ WIT OUTPUT",
								"    KTHX",
								"    I HAS A VARIABLE SLEEP_TIME TEH INTEGR ITZ 100", // Sleep 100ms
								"    BTW Sleep implementation would go here",
								"KTHX",
								"SAYZ WIT \"Process has completed\"",
								"@example Handle already finished process",
								"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
								"CMD DO PUSH WIT \"true\"", // Exits immediately
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC DO START",
								"I HAS A VARIABLE SHORT_WAIT TEH INTEGR ITZ 50", // Wait 50ms
								"BTW Sleep implementation would go here",
								"IZ PROC FINISHED?",
								"    SAYZ WIT \"Process finished quickly\"",
								"    SAYZ WIT \"Exit code: \"",
								"    SAYZ WIT PROC EXIT_CODE",
								"NOPE",
								"    SAYZ WIT \"Process still running\"",
								"    PROC DO WAIT",
								"KTHX",
								"@example Multiple WAIT calls on finished process",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC DO START",
								"I HAS A VARIABLE FIRST_WAIT TEH INTEGR ITZ PROC DO WAIT",
								"SAYZ WIT \"First wait result: \"",
								"SAYZ WIT FIRST_WAIT",
								"I HAS A VARIABLE SECOND_WAIT TEH INTEGR ITZ PROC DO WAIT",
								"SAYZ WIT \"Second wait result: \"",
								"SAYZ WIT SECOND_WAIT",
								"SAYZ WIT \"Finished status: \"",
								"SAYZ WIT PROC FINISHED",
								"@note Becomes YEZ after process terminates",
								"@note Remains YEZ even after multiple WAIT calls",
								"@note Useful for checking if WAIT will block or return immediately",
								"@note Different from RUNNING - FINISHED indicates completion, RUNNING indicates current state",
								"@see WAIT, RUNNING, EXIT_CODE",
								"@category process-lifecycle",
							},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							minionData, ok := this.NativeData.(*MinionData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "FINISHED: invalid minion context"}
							}
							return environment.BoolValue(minionData.Finished), nil
						},
						NativeSet: nil, // Read-only
					},
					"EXIT_CODE": {
						Variable: environment.Variable{
							Name:     "EXIT_CODE",
							Type:     "INTEGR",
							IsLocked: true,
							IsPublic: true,
							Documentation: []string{
								"Read-only property containing the process exit code.",
								"Available after process completion, indicates success (0) or failure (non-zero).",
								"",
								"@property {INTEGR} EXIT_CODE - Process exit status (-1 if killed, 0 if successful)",
								"@example Check process success",
								"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
								"CMD DO PUSH WIT \"ls\"",
								"CMD DO PUSH WIT \"/nonexistent\"",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC DO START",
								"I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC DO WAIT",
								"IZ EXIT_CODE SAEM AS 0?",
								"    SAYZ WIT \"Command succeeded\"",
								"NOPE",
								"    SAYZ WIT \"Command failed with code: \"",
								"    SAYZ WIT EXIT_CODE",
								"KTHX",
								"@example Access exit code without waiting",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC DO START",
								"PROC DO WAIT",
								"SAYZ WIT \"Process exit code: \"",
								"SAYZ WIT PROC EXIT_CODE",
								"@example Handle killed processes",
								"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
								"CMD DO PUSH WIT \"sleep\"",
								"CMD DO PUSH WIT \"30\"",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC DO START",
								"PROC DO KILL",
								"I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC DO WAIT",
								"SAYZ WIT \"Killed process exit code: \"",
								"SAYZ WIT EXIT_CODE",
								"@example Exit code before process finishes",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC DO START",
								"SAYZ WIT \"Exit code before completion: \"",
								"SAYZ WIT PROC EXIT_CODE", // Will be 0
								"PROC DO WAIT",
								"SAYZ WIT \"Exit code after completion: \"",
								"SAYZ WIT PROC EXIT_CODE", // Will be actual exit code
								"@example Categorize exit codes",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC DO START",
								"PROC DO WAIT",
								"I HAS A VARIABLE CODE TEH INTEGR ITZ PROC EXIT_CODE",
								"IZ CODE SAEM AS 0?",
								"    SAYZ WIT \"Success\"",
								"NOPE",
								"    IZ CODE SAEM AS -1?",
								"        SAYZ WIT \"Process was killed\"",
								"    NOPE",
								"        IZ CODE BIGGR THAN 128?",
								"            SAYZ WIT \"Process terminated by signal: \"",
								"            SAYZ WIT CODE MINUSZ 128",
								"        NOPE",
								"            SAYZ WIT \"Process failed with code: \"",
								"            SAYZ WIT CODE",
								"        KTHX",
								"    KTHX",
								"KTHX",
								"@note 0 typically indicates successful execution",
								"@note Non-zero values indicate various types of failures",
								"@note -1 indicates process was killed",
								"@note Values > 128 often indicate termination by signal",
								"@note Only meaningful after process has finished",
								"@see WAIT, FINISHED, RUNNING",
								"@category process-lifecycle",
							},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							minionData, ok := this.NativeData.(*MinionData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "EXIT_CODE: invalid minion context"}
							}
							if !minionData.Finished {
								return environment.IntegerValue(0), nil // Not finished yet
							}
							return environment.IntegerValue(minionData.ExitCode), nil
						},
						NativeSet: nil, // Read-only
					},
					"PID": {
						Variable: environment.Variable{
							Name:     "PID",
							Type:     "INTEGR",
							IsLocked: true,
							IsPublic: true,
							Documentation: []string{
								"Read-only property containing the process ID.",
								"System-assigned unique identifier for the running process.",
								"",
								"@property {INTEGR} PID - Process ID (-1 if not started)",
								"@example Get process ID after starting",
								"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
								"CMD DO PUSH WIT \"sleep\"",
								"CMD DO PUSH WIT \"10\"",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC DO START",
								"SAYZ WIT \"Process started with PID: \"",
								"SAYZ WIT PROC PID",
								"@example Monitor process by PID",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC DO START",
								"I HAS A VARIABLE PROCESS_ID TEH INTEGR ITZ PROC PID",
								"SAYZ WIT \"Monitoring process \"",
								"SAYZ WIT PROCESS_ID",
								"WHILE PROC RUNNING",
								"    SAYZ WIT \"Process \"",
								"    SAYZ WIT PROCESS_ID",
								"    SAYZ WIT \" is still running\"",
								"    I HAS A VARIABLE SLEEP_TIME TEH INTEGR ITZ 1000", // Sleep 1 second
								"    BTW Sleep implementation would go here",
								"KTHX",
								"@example PID before and after start",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"SAYZ WIT \"PID before start: \"",
								"SAYZ WIT PROC PID", // Will be -1
								"PROC DO START",
								"SAYZ WIT \"PID after start: \"",
								"SAYZ WIT PROC PID", // Will be actual PID
								"@example Use PID for external monitoring",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC DO START",
								"I HAS A VARIABLE PID TEH INTEGR ITZ PROC PID",
								"I HAS A VARIABLE MONITOR_CMD TEH BUKKIT ITZ NEW BUKKIT",
								"MONITOR_CMD DO PUSH WIT \"ps\"",
								"MONITOR_CMD DO PUSH WIT \"-p\"",
								"MONITOR_CMD DO PUSH WIT PID",
								"I HAS A VARIABLE MONITOR TEH MINION ITZ NEW MINION WIT MONITOR_CMD",
								"MONITOR DO START",
								"I HAS A VARIABLE PS_OUTPUT TEH STRIN ITZ MONITOR STDOUT DO READ WIT 1024",
								"SAYZ WIT \"Process status:\"",
								"SAYZ WIT PS_OUTPUT",
								"MONITOR DO WAIT",
								"@example Handle PID for process groups",
								"I HAS A VARIABLE PROC1 TEH MINION ITZ NEW MINION WIT CMD",
								"I HAS A VARIABLE PROC2 TEH MINION ITZ NEW MINION WIT CMD",
								"PROC1 DO START",
								"PROC2 DO START",
								"SAYZ WIT \"Started processes with PIDs: \"",
								"SAYZ WIT PROC1 PID",
								"SAYZ WIT \" and \"",
								"SAYZ WIT PROC2 PID",
								"@note Assigned by operating system when process starts",
								"@note Unique identifier for the process",
								"@note -1 indicates process has not been started",
								"@note Can be used with system tools like ps, kill, etc.",
								"@note Remains valid until process terminates",
								"@see START, RUNNING, KILL",
								"@category process-lifecycle",
							},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							minionData, ok := this.NativeData.(*MinionData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "PID: invalid minion context"}
							}
							return environment.IntegerValue(minionData.PID), nil
						},
						NativeSet: nil, // Read-only
					},
					"STDIN": {
						Variable: environment.Variable{
							Name:     "STDIN",
							Type:     "PIPE",
							IsLocked: true,
							IsPublic: true,
							Documentation: []string{
								"Read-only property providing access to the process's standard input pipe.",
								"Available after START, allows writing data to the child process.",
								"",
								"@property {PIPE} STDIN - Standard input pipe for writing to process",
								"@example Write to process stdin",
								"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
								"CMD DO PUSH WIT \"cat\"",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC DO START",
								"PROC STDIN DO WRITE WIT \"Hello, World!\\n\"",
								"PROC STDIN DO WRITE WIT \"This is input data\\n\"",
								"PROC STDIN DO CLOSE",
								"I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 1024",
								"SAYZ WIT OUTPUT",
								"@example Send commands to interactive shell",
								"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
								"CMD DO PUSH WIT \"python3\"",
								"CMD DO PUSH WIT \"-i\"",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC DO START",
								"PROC STDIN DO WRITE WIT \"print('Hello from Python')\\n\"",
								"PROC STDIN DO WRITE WIT \"x = 42\\n\"",
								"PROC STDIN DO WRITE WIT \"print(f'x = {x}')\\n\"",
								"PROC STDIN DO WRITE WIT \"exit()\\n\"",
								"I HAS A VARIABLE RESULT TEH STRIN ITZ PROC STDOUT DO READ WIT 2048",
								"SAYZ WIT RESULT",
								"@example Pipe data between processes",
								"I HAS A VARIABLE CMD1 TEH BUKKIT ITZ NEW BUKKIT",
								"CMD1 DO PUSH WIT \"echo\"",
								"CMD1 DO PUSH WIT \"line 1\\nline 2\\nline 3\"",
								"I HAS A VARIABLE PROC1 TEH MINION ITZ NEW MINION WIT CMD1",
								"I HAS A VARIABLE CMD2 TEH BUKKIT ITZ NEW BUKKIT",
								"CMD2 DO PUSH WIT \"grep\"",
								"CMD2 DO PUSH WIT \"line\"",
								"I HAS A VARIABLE PROC2 TEH MINION ITZ NEW MINION WIT CMD2",
								"PROC2 DO START",
								"PROC1 DO START",
								"I HAS A VARIABLE DATA TEH STRIN ITZ PROC1 STDOUT DO READ WIT 1024",
								"PROC2 STDIN DO WRITE WIT DATA",
								"PROC2 STDIN DO CLOSE",
								"I HAS A VARIABLE FILTERED TEH STRIN ITZ PROC2 STDOUT DO READ WIT 1024",
								"SAYZ WIT FILTERED",
								"@example Handle stdin pipe errors",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC DO START",
								"MAYB",
								"    PROC STDIN DO WRITE WIT \"some data\\n\"",
								"    PROC STDIN DO CLOSE",
								"OOPSIE ERR",
								"    SAYZ WIT \"Error writing to stdin: \"",
								"    SAYZ WIT ERR",
								"    IZ PROC STDIN IS_OPEN?",
								"        PROC STDIN DO CLOSE",
								"    KTHX",
								"KTHX",
								"@note Only available after calling START",
								"@note Supports WRITE operations to send data to process",
								"@note Should be closed when done writing to signal EOF",
								"@note Throws exception if accessed before START",
								"@note Pipe becomes unavailable after process termination",
								"@see STDOUT, STDERR, WRITE, CLOSE",
								"@category pipe-operations",
							},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							minionData, ok := this.NativeData.(*MinionData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "STDIN: invalid minion context"}
							}
							return minionData.StdinPipe, nil
						},
						NativeSet: nil, // Read-only
					},
					"STDOUT": {
						Variable: environment.Variable{
							Name:     "STDOUT",
							Type:     "PIPE",
							IsLocked: true,
							IsPublic: true,
							Documentation: []string{
								"Read-only property providing access to the process's standard output pipe.",
								"Available after START, allows reading data output by the child process.",
								"",
								"@property {PIPE} STDOUT - Standard output pipe for reading from process",
								"@example Read process output",
								"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
								"CMD DO PUSH WIT \"echo\"",
								"CMD DO PUSH WIT \"Hello, World!\"",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC DO START",
								"I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 1024",
								"SAYZ WIT \"Process output: \"",
								"SAYZ WIT OUTPUT",
								"@example Read output in chunks",
								"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
								"CMD DO PUSH WIT \"ls\"",
								"CMD DO PUSH WIT \"-la\"",
								"CMD DO PUSH WIT \"/usr/bin\"",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC DO START",
								"I HAS A VARIABLE ALL_OUTPUT TEH STRIN ITZ \"\"",
								"WHILE NO SAEM AS (PROC STDOUT IS_EOF)",
								"    I HAS A VARIABLE CHUNK TEH STRIN ITZ PROC STDOUT DO READ WIT 512",
								"    IZ CHUNK SAEM AS \"\"?",
								"        GTFO BTW No more data",
								"    KTHX",
								"    ALL_OUTPUT ITZ ALL_OUTPUT MOAR CHUNK",
								"KTHX",
								"SAYZ WIT \"Directory listing:\"",
								"SAYZ WIT ALL_OUTPUT",
								"@example Handle both stdout and stderr",
								"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
								"CMD DO PUSH WIT \"find\"",
								"CMD DO PUSH WIT \"/nonexistent\"",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC DO START",
								"I HAS A VARIABLE OUTPUT TEH STRIN ITZ \"\"",
								"I HAS A VARIABLE ERRORS TEH STRIN ITZ \"\"",
								"WHILE NO SAEM AS ((PROC STDOUT IS_EOF) AN (PROC STDERR IS_EOF))",
								"    IZ NO SAEM AS (PROC STDOUT IS_EOF)?",
								"        I HAS A VARIABLE CHUNK TEH STRIN ITZ PROC STDOUT DO READ WIT 256",
								"        OUTPUT ITZ OUTPUT MOAR CHUNK",
								"    KTHX",
								"    IZ NO SAEM AS (PROC STDERR IS_EOF)?",
								"        I HAS A VARIABLE ERR_CHUNK TEH STRIN ITZ PROC STDERR DO READ WIT 256",
								"        ERRORS ITZ ERRORS MOAR ERR_CHUNK",
								"    KTHX",
								"KTHX",
								"SAYZ WIT \"Output: \"",
								"SAYZ WIT OUTPUT",
								"SAYZ WIT \"Errors: \"",
								"SAYZ WIT ERRORS",
								"@example Process streaming output",
								"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
								"CMD DO PUSH WIT \"ping\"",
								"CMD DO PUSH WIT \"-c\"",
								"CMD DO PUSH WIT \"3\"",
								"CMD DO PUSH WIT \"8.8.8.8\"",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC DO START",
								"WHILE NO SAEM AS (PROC STDOUT IS_EOF)",
								"    I HAS A VARIABLE LINE TEH STRIN ITZ PROC STDOUT DO READ WIT 256",
								"    IZ NO SAEM AS (LINE SAEM AS \"\")?",
								"        SAYZ WIT \"Received: \"",
								"        SAYZ WIT LINE",
								"    KTHX",
								"KTHX",
								"PROC DO WAIT",
								"@note Only available after calling START",
								"@note Supports READ operations to get process output",
								"@note Reading returns empty string when no data available",
								"@note IS_EOF becomes true when process closes stdout",
								"@note Throws exception if accessed before START",
								"@see STDIN, STDERR, READ, IS_EOF",
								"@category pipe-operations",
							},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							minionData, ok := this.NativeData.(*MinionData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "STDOUT: invalid minion context"}
							}
							return minionData.StdoutPipe, nil
						},
						NativeSet: nil, // Read-only
					},
					"STDERR": {
						Variable: environment.Variable{
							Name:     "STDERR",
							Type:     "PIPE",
							Value:    environment.NOTHIN, // Set when process starts
							IsLocked: true,
							IsPublic: true,
							Documentation: []string{
								"Read-only property providing access to the process's standard error pipe.",
								"Available after START, allows reading error messages from the child process.",
								"",
								"@property {PIPE} STDERR - Standard error pipe for reading process errors",
								"@example Read process error output",
								"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
								"CMD DO PUSH WIT \"ls\"",
								"CMD DO PUSH WIT \"/nonexistent/directory\"",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC DO START",
								"I HAS A VARIABLE ERRORS TEH STRIN ITZ PROC STDERR DO READ WIT 1024",
								"IZ NO SAEM AS (ERRORS SAEM AS \"\")?",
								"    SAYZ WIT \"Process errors: \"",
								"    SAYZ WIT ERRORS",
								"NOPE",
								"    SAYZ WIT \"No errors reported\"",
								"KTHX",
								"PROC DO WAIT",
								"@example Separate stdout and stderr handling",
								"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
								"CMD DO PUSH WIT \"python3\"",
								"CMD DO PUSH WIT \"-c\"",
								"CMD DO PUSH WIT \"import sys; print('output'); print('error', file=sys.stderr)\"",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC DO START",
								"I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 512",
								"I HAS A VARIABLE ERRORS TEH STRIN ITZ PROC STDERR DO READ WIT 512",
								"SAYZ WIT \"STDOUT: \"",
								"SAYZ WIT OUTPUT",
								"SAYZ WIT \"STDERR: \"",
								"SAYZ WIT ERRORS",
								"PROC DO WAIT",
								"@example Monitor stderr for warnings",
								"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
								"CMD DO PUSH WIT \"gcc\"",
								"CMD DO PUSH WIT \"-Wall\"",
								"CMD DO PUSH WIT \"program.c\"",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC DO START",
								"I HAS A VARIABLE WARNINGS TEH STRIN ITZ \"\"",
								"WHILE NO SAEM AS (PROC STDERR IS_EOF)",
								"    I HAS A VARIABLE CHUNK TEH STRIN ITZ PROC STDERR DO READ WIT 256",
								"    WARNINGS ITZ WARNINGS MOAR CHUNK",
								"KTHX",
								"PROC DO WAIT",
								"IZ NO SAEM AS (WARNINGS SAEM AS \"\")?",
								"    SAYZ WIT \"Compiler warnings:\"",
								"    SAYZ WIT WARNINGS",
								"KTHX",
								"@example Check for errors after process completion",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC DO START",
								"PROC DO WAIT",
								"I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC EXIT_CODE",
								"IZ EXIT_CODE SAEM AS 0?",
								"    I HAS A VARIABLE ERRORS TEH STRIN ITZ PROC STDERR DO READ WIT 1024",
								"    IZ NO SAEM AS (ERRORS SAEM AS \"\")?",
								"        SAYZ WIT \"Process succeeded but produced warnings: \"",
								"        SAYZ WIT ERRORS",
								"    KTHX",
								"NOPE",
								"    I HAS A VARIABLE ERRORS TEH STRIN ITZ PROC STDERR DO READ WIT 1024",
								"    SAYZ WIT \"Process failed with errors: \"",
								"    SAYZ WIT ERRORS",
								"KTHX",
								"@example Real-time error monitoring",
								"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
								"CMD DO PUSH WIT \"long-running-command\"",
								"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
								"PROC DO START",
								"WHILE PROC RUNNING",
								"    I HAS A VARIABLE ERROR_CHUNK TEH STRIN ITZ PROC STDERR DO READ WIT 256",
								"    IZ NO SAEM AS (ERROR_CHUNK SAEM AS \"\")?",
								"        SAYZ WIT \"[ERROR] \"",
								"        SAYZ WIT ERROR_CHUNK",
								"    KTHX",
								"    I HAS A VARIABLE SLEEP_TIME TEH INTEGR ITZ 500", // Sleep 500ms
								"    BTW Sleep implementation would go here",
								"KTHX",
								"PROC DO WAIT",
								"@note Only available after calling START",
								"@note Supports READ operations to get error messages",
								"@note Many programs write diagnostic messages to stderr",
								"@note Should be checked even when exit code is 0",
								"@note Throws exception if accessed before START",
								"@see STDIN, STDOUT, READ, IS_EOF",
								"@category pipe-operations",
							},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							minionData, ok := this.NativeData.(*MinionData)
							if !ok {
								return environment.NOTHIN, runtime.Exception{Message: "STDERR: invalid minion context"}
							}
							return minionData.StderrPipe, nil
						},
						NativeSet: nil, // Read-only
					},
				},
				PrivateVariables: make(map[string]*environment.MemberVariable),
				PrivateFunctions: make(map[string]*environment.Function),
				SharedVariables:  make(map[string]*environment.MemberVariable),
				SharedFunctions:  make(map[string]*environment.Function),
			},
		}
	})
	return processClasses
}

// RegisterPROCESSInEnv registers PROCESS classes in the given environment
// declarations: empty slice means import all, otherwise import only specified classes
func RegisterPROCESSInEnv(env *environment.Environment, declarations ...string) error {
	// First ensure IO classes are available since PIPE uses IO interfaces
	err := RegisterIOInEnv(env, "READER", "WRITER")
	if err != nil {
		return runtime.Exception{Message: fmt.Sprintf("failed to register IO classes for PROCESS: %v", err)}
	}

	processClasses := getProcessClasses()

	// If declarations is empty, import all classes
	if len(declarations) == 0 {
		for _, class := range processClasses {
			env.DefineClass(class)
		}
		return nil
	}

	// Otherwise, import only specified classes
	for _, decl := range declarations {
		declUpper := strings.ToUpper(decl)
		if class, exists := processClasses[declUpper]; exists {
			env.DefineClass(class)
		} else {
			return runtime.Exception{Message: fmt.Sprintf("unknown PROCESS declaration: %s", decl)}
		}
	}

	return nil
}
