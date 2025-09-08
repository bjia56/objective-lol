package stdlib

import (
	"fmt"
	"io"
	"os"
	"os/exec"
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
	CmdLine    []string                    // Command and arguments
	WorkDir    string                      // Working directory
	Env        []string                    // Environment variables
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
	return &environment.ObjectInstance{
		Environment: env,
		Class:       class,
		NativeData: &PipeData{
			FDType:  fdType,
			Pipe:    nil,
			IsOpen:  false,
			IsEOF:   false,
			Process: process,
		},
		MRO: class.MRO,
		Variables: map[string]*environment.Variable{
			"FD_TYPE": {
				Name:     "FD_TYPE",
				Type:     "STRIN",
				Value:    environment.StringValue(fdType),
				IsLocked: true,
				IsPublic: true,
			},
			"IS_OPEN": {
				Name:     "IS_OPEN",
				Type:     "BOOL",
				Value:    environment.NO,
				IsLocked: true,
				IsPublic: true,
			},
			"IS_EOF": {
				Name:     "IS_EOF",
				Type:     "BOOL",
				Value:    environment.NO,
				IsLocked: true,
				IsPublic: true,
			},
		},
	}
}

// updatePipeStatus updates the status variables of a PIPE object
func updatePipeStatus(obj *environment.ObjectInstance, pipeData *PipeData) {
	if isOpenVar, exists := obj.Variables["IS_OPEN"]; exists {
		isOpenVar.Value = environment.BoolValue(pipeData.IsOpen)
	}
	if isEOFVar, exists := obj.Variables["IS_EOF"]; exists {
		isEOFVar.Value = environment.BoolValue(pipeData.IsEOF)
	}
}

// updateMinionStatus updates the status variables of a MINION object
func updateMinionStatus(obj *environment.ObjectInstance, minionData *MinionData) {
	if runningVar, exists := obj.Variables["RUNNING"]; exists {
		runningVar.Value = environment.BoolValue(minionData.Running)
	}
	if finishedVar, exists := obj.Variables["FINISHED"]; exists {
		finishedVar.Value = environment.BoolValue(minionData.Finished)
	}
	if exitCodeVar, exists := obj.Variables["EXIT_CODE"]; exists {
		exitCodeVar.Value = environment.IntegerValue(minionData.ExitCode)
	}
	if pidVar, exists := obj.Variables["PID"]; exists {
		pidVar.Value = environment.IntegerValue(minionData.PID)
	}
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
				PublicFunctions: map[string]*environment.Function{
					// READ method (for stdout/stderr pipes)
					"READ": {
						Name:       "READ",
						ReturnType: "STRIN",
						Parameters: []environment.Parameter{
							{Name: "size", Type: "INTEGR"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this environment.GenericObject, args []environment.Value) (environment.Value, error) {
							size := args[0]

							sizeVal, ok := size.(environment.IntegerValue)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("READ expects INTEGR size, got %s", size.Type())
							}

							thisObj := this.(*environment.ObjectInstance)
							pipeData, ok := thisObj.NativeData.(*PipeData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("READ: invalid pipe context")
							}

							if !pipeData.IsOpen {
								return environment.NOTHIN, runtime.Exception{Message: "Pipe is not open"}
							}

							if pipeData.FDType == "STDIN" {
								return environment.NOTHIN, runtime.Exception{Message: "Cannot read from stdin pipe"}
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
								return environment.NOTHIN, fmt.Errorf("READ: pipe is not readable")
							}

							n, err := reader.Read(buffer)
							if err != nil {
								if err == io.EOF {
									pipeData.IsEOF = true
									updatePipeStatus(thisObj, pipeData)
								} else {
									return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("Read error: %v", err)}
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
						NativeImpl: func(interpreter environment.Interpreter, this environment.GenericObject, args []environment.Value) (environment.Value, error) {
							data := args[0]

							dataVal, ok := data.(environment.StringValue)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("WRITE expects STRIN data, got %s", data.Type())
							}

							thisObj := this.(*environment.ObjectInstance)
							pipeData, ok := thisObj.NativeData.(*PipeData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("WRITE: invalid pipe context")
							}

							if !pipeData.IsOpen {
								return environment.NOTHIN, runtime.Exception{Message: "Pipe is not open"}
							}

							if pipeData.FDType != "STDIN" {
								return environment.NOTHIN, runtime.Exception{Message: "Cannot write to stdout/stderr pipe"}
							}

							writer, ok := pipeData.Pipe.(io.Writer)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("WRITE: pipe is not writable")
							}

							dataStr := string(dataVal)
							n, err := writer.Write([]byte(dataStr))
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("Write error: %v", err)}
							}

							return environment.IntegerValue(n), nil
						},
					},
					// CLOSE method
					"CLOSE": {
						Name: "CLOSE",
						NativeImpl: func(interpreter environment.Interpreter, this environment.GenericObject, args []environment.Value) (environment.Value, error) {
							thisObj := this.(*environment.ObjectInstance)
							pipeData, ok := thisObj.NativeData.(*PipeData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("CLOSE: invalid pipe context")
							}

							if !pipeData.IsOpen {
								return environment.NOTHIN, nil // Already closed
							}

							if pipeData.Pipe != nil {
								err := pipeData.Pipe.Close()
								if err != nil {
									return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("Close error: %v", err)}
								}
							}

							pipeData.IsOpen = false
							updatePipeStatus(thisObj, pipeData)

							return environment.NOTHIN, nil
						},
					},
				},
				PublicVariables: map[string]*environment.Variable{
					"FD_TYPE": {
						Name:     "FD_TYPE",
						Type:     "STRIN",
						Value:    environment.StringValue(""),
						IsLocked: true,
						IsPublic: true,
					},
					"IS_OPEN": {
						Name:     "IS_OPEN",
						Type:     "BOOL",
						Value:    environment.NO,
						IsLocked: true,
						IsPublic: true,
					},
					"IS_EOF": {
						Name:     "IS_EOF",
						Type:     "BOOL",
						Value:    environment.NO,
						IsLocked: true,
						IsPublic: true,
					},
				},
				PrivateVariables: make(map[string]*environment.Variable),
				PrivateFunctions: make(map[string]*environment.Function),
				SharedVariables:  make(map[string]*environment.Variable),
				SharedFunctions:  make(map[string]*environment.Function),
			},
			"MINION": {
				Name:          "MINION",
				QualifiedName: "stdlib:PROCESS.MINION",
				ModulePath:    "stdlib:PROCESS",
				ParentClasses: []string{},
				MRO:           []string{"stdlib:PROCESS.MINION"},
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"MINION": {
						Name: "MINION",
						Parameters: []environment.Parameter{
							{Name: "cmdline", Type: "BUKKIT"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this environment.GenericObject, args []environment.Value) (environment.Value, error) {
							cmdline := args[0]

							// Validate that the argument is a BUKKIT
							cmdlineInstance, ok := cmdline.(*environment.ObjectInstance)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("MINION constructor expects BUKKIT cmdline, got %s", cmdline.Type())
							}

							if cmdlineInstance.Class.Name != "BUKKIT" {
								return environment.NOTHIN, fmt.Errorf("MINION constructor expects BUKKIT cmdline, got %s", cmdlineInstance.Class.Name)
							}

							// Extract command line from BUKKIT
							cmdLineSlice := []string{}
							if nativeData, ok := cmdlineInstance.NativeData.(BukkitSlice); ok {
								for _, val := range nativeData {
									if strVal, ok := val.(environment.StringValue); ok {
										cmdLineSlice = append(cmdLineSlice, string(strVal))
									} else {
										return environment.NOTHIN, fmt.Errorf("MINION constructor: all cmdline elements must be STRIN, got %s", val.Type())
									}
								}
							} else {
								return environment.NOTHIN, fmt.Errorf("MINION constructor: invalid BUKKIT data")
							}

							if len(cmdLineSlice) == 0 {
								return environment.NOTHIN, runtime.Exception{Message: "Command line cannot be empty"}
							}

							// Initialize the minion data
							minionData := &MinionData{
								CmdLine:  cmdLineSlice,
								WorkDir:  "",           // Use current directory by default
								Env:      os.Environ(), // Use current environment by default
								Running:  false,
								Finished: false,
								ExitCode: 0,
								PID:      -1,
							}
							thisObj := this.(*environment.ObjectInstance)
							thisObj.NativeData = minionData

							// Set CMDLINE variable
							thisObj.Variables["CMDLINE"] = &environment.Variable{
								Name:     "CMDLINE",
								Type:     "BUKKIT",
								Value:    cmdline,
								IsLocked: true,
								IsPublic: true,
							}

							// Initialize all status variables
							thisObj.Variables["RUNNING"] = &environment.Variable{
								Name:     "RUNNING",
								Type:     "BOOL",
								Value:    environment.NO,
								IsLocked: true,
								IsPublic: true,
							}
							thisObj.Variables["FINISHED"] = &environment.Variable{
								Name:     "FINISHED",
								Type:     "BOOL",
								Value:    environment.NO,
								IsLocked: true,
								IsPublic: true,
							}
							thisObj.Variables["EXIT_CODE"] = &environment.Variable{
								Name:     "EXIT_CODE",
								Type:     "INTEGR",
								Value:    environment.IntegerValue(0),
								IsLocked: true,
								IsPublic: true,
							}
							thisObj.Variables["PID"] = &environment.Variable{
								Name:     "PID",
								Type:     "INTEGR",
								Value:    environment.IntegerValue(-1),
								IsLocked: true,
								IsPublic: true,
							}

							updateMinionStatus(thisObj, minionData)

							return environment.NOTHIN, nil
						},
					},
					// SET_WORKDIR method
					"SET_WORKDIR": {
						Name: "SET_WORKDIR",
						Parameters: []environment.Parameter{
							{Name: "dir", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this environment.GenericObject, args []environment.Value) (environment.Value, error) {
							dir := args[0]

							dirVal, ok := dir.(environment.StringValue)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("SET_WORKDIR expects STRIN dir, got %s", dir.Type())
							}

							minionData, ok := this.(*environment.ObjectInstance).NativeData.(*MinionData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("SET_WORKDIR: invalid minion context")
							}

							if minionData.Running {
								return environment.NOTHIN, runtime.Exception{Message: "Cannot change working directory after process started"}
							}

							minionData.WorkDir = string(dirVal)
							return environment.NOTHIN, nil
						},
					},
					// SET_ENV method
					"SET_ENV": {
						Name: "SET_ENV",
						Parameters: []environment.Parameter{
							{Name: "env", Type: "BASKIT"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this environment.GenericObject, args []environment.Value) (environment.Value, error) {
							env := args[0]

							// Validate that the argument is a BASKIT
							envInstance, ok := env.(*environment.ObjectInstance)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("SET_ENV expects BASKIT env, got %s", env.Type())
							}

							if envInstance.Class.Name != "BASKIT" {
								return environment.NOTHIN, fmt.Errorf("SET_ENV expects BASKIT env, got %s", envInstance.Class.Name)
							}

							minionData, ok := this.(*environment.ObjectInstance).NativeData.(*MinionData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("SET_ENV: invalid minion context")
							}

							if minionData.Running {
								return environment.NOTHIN, runtime.Exception{Message: "Cannot change environment after process started"}
							}

							// Extract environment variables from BASKIT
							envSlice := []string{}
							if nativeData, ok := envInstance.NativeData.(BaskitMap); ok {
								for key, val := range nativeData {
									if strVal, ok := val.(environment.StringValue); ok {
										envSlice = append(envSlice, key+"="+string(strVal))
									} else {
										return environment.NOTHIN, fmt.Errorf("SET_ENV: all environment values must be STRIN, got %s for key %s", val.Type(), key)
									}
								}
							} else {
								return environment.NOTHIN, fmt.Errorf("SET_ENV: invalid BASKIT data")
							}

							minionData.Env = envSlice
							return environment.NOTHIN, nil
						},
					},
					// ADD_ENV method
					"ADD_ENV": {
						Name: "ADD_ENV",
						Parameters: []environment.Parameter{
							{Name: "key", Type: "STRIN"},
							{Name: "value", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this environment.GenericObject, args []environment.Value) (environment.Value, error) {
							key := args[0]
							value := args[1]

							keyVal, ok := key.(environment.StringValue)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("ADD_ENV expects STRIN key, got %s", key.Type())
							}

							valueVal, ok := value.(environment.StringValue)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("ADD_ENV expects STRIN value, got %s", value.Type())
							}

							minionData, ok := this.(*environment.ObjectInstance).NativeData.(*MinionData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("ADD_ENV: invalid minion context")
							}

							if minionData.Running {
								return environment.NOTHIN, runtime.Exception{Message: "Cannot change environment after process started"}
							}

							// Add or update the environment variable
							envVar := string(keyVal) + "=" + string(valueVal)

							// Remove existing entry with same key
							keyPrefix := string(keyVal) + "="
							newEnv := []string{}
							for _, env := range minionData.Env {
								if !strings.HasPrefix(env, keyPrefix) {
									newEnv = append(newEnv, env)
								}
							}

							// Add the new entry
							newEnv = append(newEnv, envVar)
							minionData.Env = newEnv

							return environment.NOTHIN, nil
						},
					},
					// START method
					"START": {
						Name: "START",
						NativeImpl: func(interpreter environment.Interpreter, this environment.GenericObject, args []environment.Value) (environment.Value, error) {
							minionData, ok := this.(*environment.ObjectInstance).NativeData.(*MinionData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("START: invalid minion context")
							}

							if minionData.Running {
								return environment.NOTHIN, runtime.Exception{Message: "Process is already running"}
							}

							if minionData.Finished {
								return environment.NOTHIN, runtime.Exception{Message: "Process has already finished"}
							}

							// Create the command
							cmd := exec.Command(minionData.CmdLine[0], minionData.CmdLine[1:]...)
							if minionData.WorkDir != "" {
								cmd.Dir = minionData.WorkDir
							}
							cmd.Env = minionData.Env

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
							updatePipeStatus(stdinPipe, stdinPipe.NativeData.(*PipeData))

							stdoutPipe := NewPipeInstance("STDOUT", cmd)
							stdoutPipe.NativeData.(*PipeData).Pipe = stdout
							stdoutPipe.NativeData.(*PipeData).IsOpen = true
							updatePipeStatus(stdoutPipe, stdoutPipe.NativeData.(*PipeData))

							stderrPipe := NewPipeInstance("STDERR", cmd)
							stderrPipe.NativeData.(*PipeData).Pipe = stderr
							stderrPipe.NativeData.(*PipeData).IsOpen = true
							updatePipeStatus(stderrPipe, stderrPipe.NativeData.(*PipeData))

							minionData.StdinPipe = stdinPipe
							minionData.StdoutPipe = stdoutPipe
							minionData.StderrPipe = stderrPipe

							// Set PIPE variables
							thisObj := this.(*environment.ObjectInstance)
							thisObj.Variables["STDIN"] = &environment.Variable{
								Name:     "STDIN",
								Type:     "PIPE",
								Value:    stdinPipe,
								IsLocked: true,
								IsPublic: true,
							}
							thisObj.Variables["STDOUT"] = &environment.Variable{
								Name:     "STDOUT",
								Type:     "PIPE",
								Value:    stdoutPipe,
								IsLocked: true,
								IsPublic: true,
							}
							thisObj.Variables["STDERR"] = &environment.Variable{
								Name:     "STDERR",
								Type:     "PIPE",
								Value:    stderrPipe,
								IsLocked: true,
								IsPublic: true,
							}

							updateMinionStatus(thisObj, minionData)

							return environment.NOTHIN, nil
						},
					},
					// WAIT method
					"WAIT": {
						Name:       "WAIT",
						ReturnType: "INTEGR",
						NativeImpl: func(interpreter environment.Interpreter, this environment.GenericObject, args []environment.Value) (environment.Value, error) {
							thisObj := this.(*environment.ObjectInstance)
							minionData, ok := this.(*environment.ObjectInstance).NativeData.(*MinionData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("WAIT: invalid minion context")
							}

							if !minionData.Running && !minionData.Finished {
								return environment.NOTHIN, runtime.Exception{Message: "Process has not been started"}
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
									return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("Process wait error: %v", err)}
								}
							} else {
								minionData.ExitCode = 0
							}

							updateMinionStatus(thisObj, minionData)

							return environment.IntegerValue(minionData.ExitCode), nil
						},
					},
					// KILL method
					"KILL": {
						Name: "KILL",
						NativeImpl: func(interpreter environment.Interpreter, this environment.GenericObject, args []environment.Value) (environment.Value, error) {
							thisObj := this.(*environment.ObjectInstance)
							minionData, ok := thisObj.NativeData.(*MinionData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("KILL: invalid minion context")
							}

							if !minionData.Running {
								return environment.NOTHIN, runtime.Exception{Message: "Process is not running"}
							}

							err := minionData.Process.Process.Kill()
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("Failed to kill process: %v", err)}
							}

							// Wait for process to finish after killing
							minionData.Process.Wait()
							minionData.Running = false
							minionData.Finished = true
							minionData.ExitCode = -1 // Indicate killed

							updateMinionStatus(thisObj, minionData)

							return environment.NOTHIN, nil
						},
					},
					// SIGNAL method
					"SIGNAL": {
						Name: "SIGNAL",
						Parameters: []environment.Parameter{
							{Name: "code", Type: "INTEGR"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this environment.GenericObject, args []environment.Value) (environment.Value, error) {
							code := args[0]

							codeVal, ok := code.(environment.IntegerValue)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("SIGNAL expects INTEGR code, got %s", code.Type())
							}

							minionData, ok := this.(*environment.ObjectInstance).NativeData.(*MinionData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("SIGNAL: invalid minion context")
							}

							if !minionData.Running {
								return environment.NOTHIN, runtime.Exception{Message: "Process is not running"}
							}

							signal := syscall.Signal(int(codeVal))
							err := minionData.Process.Process.Signal(signal)
							if err != nil {
								return environment.NOTHIN, runtime.Exception{Message: fmt.Sprintf("Failed to send signal: %v", err)}
							}

							return environment.NOTHIN, nil
						},
					},
					// IS_ALIVE method
					"IS_ALIVE": {
						Name:       "IS_ALIVE",
						ReturnType: "BOOL",
						NativeImpl: func(interpreter environment.Interpreter, this environment.GenericObject, args []environment.Value) (environment.Value, error) {
							minionData, ok := this.(*environment.ObjectInstance).NativeData.(*MinionData)
							if !ok {
								return environment.NOTHIN, fmt.Errorf("IS_ALIVE: invalid minion context")
							}

							return environment.BoolValue(minionData.Running), nil
						},
					},
				},
				PublicVariables: map[string]*environment.Variable{
					"CMDLINE": {
						Name:     "CMDLINE",
						Type:     "BUKKIT",
						Value:    environment.NOTHIN, // Set in constructor
						IsLocked: true,
						IsPublic: true,
					},
					"RUNNING": {
						Name:     "RUNNING",
						Type:     "BOOL",
						Value:    environment.NO,
						IsLocked: true,
						IsPublic: true,
					},
					"FINISHED": {
						Name:     "FINISHED",
						Type:     "BOOL",
						Value:    environment.NO,
						IsLocked: true,
						IsPublic: true,
					},
					"EXIT_CODE": {
						Name:     "EXIT_CODE",
						Type:     "INTEGR",
						Value:    environment.IntegerValue(0),
						IsLocked: true,
						IsPublic: true,
					},
					"PID": {
						Name:     "PID",
						Type:     "INTEGR",
						Value:    environment.IntegerValue(-1),
						IsLocked: true,
						IsPublic: true,
					},
					"STDIN": {
						Name:     "STDIN",
						Type:     "PIPE",
						Value:    environment.NOTHIN, // Set when process starts
						IsLocked: true,
						IsPublic: true,
					},
					"STDOUT": {
						Name:     "STDOUT",
						Type:     "PIPE",
						Value:    environment.NOTHIN, // Set when process starts
						IsLocked: true,
						IsPublic: true,
					},
					"STDERR": {
						Name:     "STDERR",
						Type:     "PIPE",
						Value:    environment.NOTHIN, // Set when process starts
						IsLocked: true,
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
	return processClasses
}

// RegisterPROCESSInEnv registers PROCESS classes in the given environment
// declarations: empty slice means import all, otherwise import only specified classes
func RegisterPROCESSInEnv(env *environment.Environment, declarations ...string) error {
	// First ensure IO classes are available since PIPE uses IO interfaces
	err := RegisterIOInEnv(env, "READER", "WRITER")
	if err != nil {
		return fmt.Errorf("failed to register IO classes: %v", err)
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
			return fmt.Errorf("unknown PROCESS class: %s", decl)
		}
	}

	return nil
}
