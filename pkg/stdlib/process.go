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
					"A communication pipe connected to a process's stdin, stdout, or stderr.",
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
							"Reads data from stdout or stderr pipes.",
							"Cannot be used on stdin pipes.",
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
							"Writes data to stdin pipe.",
							"Returns number of bytes written, cannot be used on stdout/stderr pipes.",
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
							"Closes the pipe connection.",
							"Automatically closes when process terminates.",
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
								"Pipe type. STDIN, STDOUT, or STDERR.",
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
								"Whether pipe is open for operations.",
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
								"Whether end-of-file reached (read pipes only).",
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
					"A child process that can be launched, monitored, and controlled.",
				},
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"MINION": {
						Name: "MINION",
						Parameters: []environment.Parameter{
							{Name: "cmdline", Type: "BUKKIT"},
						},
						Documentation: []string{
							"Creates a process definition with command and arguments as BUKKIT.",
							"Sets working directory to current directory and environment to current environment by default.",
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
							"Launches the process and creates stdin, stdout, and stderr pipes.",
							"Throws exception if process fails to start or is already running.",
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
