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
					"I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC HAZ STDOUT DO READ WIT 1024",
					"SAYZ WIT OUTPUT",
					"PROC DO WAIT",
					"@example Writing to process stdin",
					"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
					"CMD DO PUSH WIT \"cat\"",
					"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
					"PROC DO START",
					"PROC HAZ STDIN DO WRITE WIT \"Hello from parent process!\"",
					"PROC HAZ STDIN DO CLOSE",
					"I HAS A VARIABLE RESULT TEH STRIN ITZ PROC HAZ STDOUT DO READ WIT 1024",
					"SAYZ WIT RESULT",
					"@example Handle process stderr",
					"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
					"PROC DO START",
					"I HAS A VARIABLE ERROR_OUTPUT TEH STRIN ITZ PROC HAZ STDERR DO READ WIT 512",
					"IZ NOT (ERROR_OUTPUT SAEM AS \"\")?",
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
							"I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC HAZ STDOUT DO READ WIT 4096",
							"SAYZ WIT \"Directory listing:\"",
							"SAYZ WIT OUTPUT",
							"@example Read in chunks",
							"I HAS A VARIABLE BUFFER TEH STRIN ITZ \"\"",
							"IM OUTTA UR LOOP WHILE NOT (PROC HAZ STDOUT HAZ IS_EOF)",
							"    I HAS A VARIABLE CHUNK TEH STRIN ITZ PROC HAZ STDOUT DO READ WIT 256",
							"    IZ CHUNK SAEM AS \"\"?",
							"        GTFO BTW End of output",
							"    KTHX",
							"    BUFFER ITZ BUFFER MOAR CHUNK",
							"IM IN UR LOOP",
							"@example Read stderr separately",
							"I HAS A VARIABLE ERROR_DATA TEH STRIN ITZ PROC HAZ STDERR DO READ WIT 1024",
							"IZ NOT (ERROR_DATA SAEM AS \"\")?",
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
							"PROC HAZ STDIN DO WRITE WIT \"print('Hello from Python!')\\n\"",
							"PROC HAZ STDIN DO WRITE WIT \"exit()\\n\"",
							"I HAS A VARIABLE RESULT TEH STRIN ITZ PROC HAZ STDOUT DO READ WIT 1024",
							"SAYZ WIT RESULT",
							"@example Send data to filter process",
							"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
							"CMD DO PUSH WIT \"grep\"",
							"CMD DO PUSH WIT \"error\"",
							"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
							"PROC DO START",
							"PROC HAZ STDIN DO WRITE WIT \"info: starting process\\n\"",
							"PROC HAZ STDIN DO WRITE WIT \"error: something went wrong\\n\"",
							"PROC HAZ STDIN DO WRITE WIT \"info: process completed\\n\"",
							"PROC HAZ STDIN DO CLOSE BTW Signal end of input",
							"I HAS A VARIABLE FILTERED TEH STRIN ITZ PROC HAZ STDOUT DO READ WIT 1024",
							"SAYZ WIT \"Filtered output: \"",
							"SAYZ WIT FILTERED",
							"@example Pipe data between processes",
							"PROC HAZ STDIN DO WRITE WIT \"line 1\\n\"",
							"PROC HAZ STDIN DO WRITE WIT \"line 2\\n\"",
							"I HAS A VARIABLE BYTES_WRITTEN TEH INTEGR ITZ PROC HAZ STDIN DO WRITE WIT \"line 3\\n\"",
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
							"PROC HAZ STDIN DO WRITE WIT \"line 1\\n\"",
							"PROC HAZ STDIN DO WRITE WIT \"line 2\\n\"",
							"PROC HAZ STDIN DO WRITE WIT \"line 3\\n\"",
							"PROC HAZ STDIN DO CLOSE BTW Signal end of input",
							"I HAS A VARIABLE COUNT TEH STRIN ITZ PROC HAZ STDOUT DO READ WIT 100",
							"SAYZ WIT \"Line count: \"",
							"SAYZ WIT COUNT",
							"@example Cleanup pipes after process",
							"PROC DO WAIT",
							"PROC HAZ STDOUT DO CLOSE",
							"PROC HAZ STDERR DO CLOSE",
							"SAYZ WIT \"All pipes closed\"",
							"@example Close in error handling",
							"MAYB",
							"    PROC HAZ STDIN DO WRITE WIT \"some input\"",
							"    BTW More operations here",
							"OOPSIE ERR",
							"    SAYZ WIT \"Error occurred: \"",
							"    SAYZ WIT ERR",
							"FINALLY",
							"    IZ PROC HAZ STDIN HAZ IS_OPEN?",
							"        PROC HAZ STDIN DO CLOSE",
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
								"IZ (PROC HAZ STDIN HAZ FD_TYPE) SAEM AS \"STDIN\"?",
								"    PROC HAZ STDIN DO WRITE WIT \"input data\"",
								"    PROC HAZ STDIN DO CLOSE",
								"KTHX",
								"@example Handle different pipe types",
								"I HAS A VARIABLE PIPES TEH BUKKIT ITZ NEW BUKKIT",
								"PIPES DO PUSH WIT PROC HAZ STDIN",
								"PIPES DO PUSH WIT PROC HAZ STDOUT",
								"PIPES DO PUSH WIT PROC HAZ STDERR",
								"IM OUTTA UR PIPES NERFIN PIPE",
								"    I HAS A VARIABLE TYPE TEH STRIN ITZ PIPE HAZ FD_TYPE",
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
								"IZ PROC HAZ STDIN HAZ IS_OPEN?",
								"    PROC HAZ STDIN DO WRITE WIT \"data\"",
								"    PROC HAZ STDIN DO CLOSE",
								"NOPE",
								"    SAYZ WIT \"Stdin pipe is not available\"",
								"KTHX",
								"@example Monitor pipe status during operations",
								"IM OUTTA UR LOOP WHILE (PROC HAZ STDOUT HAZ IS_OPEN)",
								"    I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC HAZ STDOUT DO READ WIT 256",
								"    IZ OUTPUT SAEM AS \"\"?",
								"        GTFO BTW End of output reached",
								"    KTHX",
								"    SAYZ WIT OUTPUT",
								"IM IN UR LOOP",
								"@example Safe pipe cleanup",
								"IZ PROC HAZ STDERR HAZ IS_OPEN?",
								"    I HAS A VARIABLE ERRORS TEH STRIN ITZ PROC HAZ STDERR DO READ WIT 1024",
								"    PROC HAZ STDERR DO CLOSE",
								"    IZ NOT (ERRORS SAEM AS \"\")?",
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
								"IM OUTTA UR LOOP WHILE NOT (PROC HAZ STDOUT HAZ IS_EOF)",
								"    I HAS A VARIABLE CHUNK TEH STRIN ITZ PROC HAZ STDOUT DO READ WIT 512",
								"    IZ CHUNK SAEM AS \"\"?",
								"        GTFO BTW No more data available",
								"    KTHX",
								"    ALL_OUTPUT ITZ ALL_OUTPUT MOAR CHUNK",
								"IM IN UR LOOP",
								"SAYZ WIT \"Complete output: \"",
								"SAYZ WIT ALL_OUTPUT",
								"@example Check EOF status after reading",
								"I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC HAZ STDOUT DO READ WIT 1024",
								"IZ PROC HAZ STDOUT HAZ IS_EOF?",
								"    SAYZ WIT \"Reached end of process output\"",
								"NOPE",
								"    SAYZ WIT \"More output may be available\"",
								"KTHX",
								"@example Handle both stdout and stderr EOF",
								"IM OUTTA UR LOOP WHILE NOT ((PROC HAZ STDOUT HAZ IS_EOF) AN (PROC HAZ STDERR HAZ IS_EOF))",
								"    IZ NOT (PROC HAZ STDOUT HAZ IS_EOF)?",
								"        I HAS A VARIABLE OUT TEH STRIN ITZ PROC HAZ STDOUT DO READ WIT 256",
								"        IZ NOT (OUT SAEM AS \"\")?",
								"            SAYZ WIT \"STDOUT: \"",
								"            SAYZ WIT OUT",
								"        KTHX",
								"    KTHX",
								"    IZ NOT (PROC HAZ STDERR HAZ IS_EOF)?",
								"        I HAS A VARIABLE ERR TEH STRIN ITZ PROC HAZ STDERR DO READ WIT 256",
								"        IZ NOT (ERR SAEM AS \"\")?",
								"            SAYZ WIT \"STDERR: \"",
								"            SAYZ WIT ERR",
								"        KTHX",
								"    KTHX",
								"IM IN UR LOOP",
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
					"I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC HAZ STDOUT DO READ WIT 1024",
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
					"PYTHON HAZ STDIN DO WRITE WIT \"print(2 + 2)\\n\"",
					"PYTHON HAZ STDIN DO WRITE WIT \"exit()\\n\"",
					"I HAS A VARIABLE RESULT TEH STRIN ITZ PYTHON HAZ STDOUT DO READ WIT 1024",
					"PYTHON DO WAIT",
					"SAYZ WIT \"Python result: \"",
					"SAYZ WIT RESULT",
					"@example Process with custom environment and working directory",
					"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
					"CMD DO PUSH WIT \"printenv\"",
					"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
					"PROC HAZ WORKDIR ITZ \"/tmp\"",
					"PROC HAZ ENV DO PUT WIT \"CUSTOM_VAR\" AN WIT \"Hello from environment!\"",
					"PROC DO START",
					"I HAS A VARIABLE ENV_OUTPUT TEH STRIN ITZ PROC HAZ STDOUT DO READ WIT 2048",
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
							"SAYZ WIT PROC HAZ PID",
							"@example Start and immediately read output",
							"I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT",
							"CMD DO PUSH WIT \"echo\"",
							"CMD DO PUSH WIT \"Hello!\"",
							"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
							"PROC DO START",
							"I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC HAZ STDOUT DO READ WIT 1024",
							"SAYZ WIT \"Process output: \"",
							"SAYZ WIT OUTPUT",
							"@example Start process with custom environment",
							"I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD",
							"PROC HAZ ENV DO PUT WIT \"MY_VAR\" AN WIT \"custom_value\"",
							"PROC HAZ WORKDIR ITZ \"/tmp\"",
							"PROC DO START",
							"SAYZ WIT \"Process started in directory: \"",
							"SAYZ WIT PROC HAZ WORKDIR",
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
							"Throws exception if process has not been started.",
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
							"Forcefully terminates the running process,",
							"Throws exception if process is not running",
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
							"Throws exception if process is not running.",
							"Not supported on Windows.",
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
								"Command line arguments.",
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
								"Working directory.",
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
								"Environment variables.",
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
								"Whether process is currently running.",
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
								"Process completion status.",
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
								"Process exit code.",
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
								"Process ID (-1 if not started).",
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
								"Stdin pipe (available after START).",
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
								"Stdout pipe (available after START).",
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
								"Stderr pipe (available after START).",
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
