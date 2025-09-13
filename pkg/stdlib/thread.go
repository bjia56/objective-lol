package stdlib

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/runtime"
)

// ThreadData holds the native Go threading constructs for YARN objects
type ThreadData struct {
	goroutineRunning bool
	finished         bool
	wg               sync.WaitGroup
	interpreter      environment.Interpreter
	result           environment.Value
	err              error
}

// MutexData holds the native Go mutex for KNOT objects
type MutexData struct {
	mutex  sync.Mutex
	locked bool
}

// NewYarnInstance creates a new YARN thread object instance
func NewYarnInstance() *environment.ObjectInstance {
	class := getThreadClasses()["YARN"]
	env := environment.NewEnvironment(nil)
	env.DefineClass(class)
	obj := &environment.ObjectInstance{
		Environment: env,
		Class:       class,
		NativeData:  &ThreadData{},
		Variables:   make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(obj)
	return obj
}

// NewKnotInstance creates a new KNOT mutex object instance
func NewKnotInstance() *environment.ObjectInstance {
	class := getThreadClasses()["KNOT"]
	env := environment.NewEnvironment(nil)
	env.DefineClass(class)
	obj := &environment.ObjectInstance{
		Environment: env,
		Class:       class,
		NativeData:  &MutexData{},
		Variables:   make(map[string]*environment.MemberVariable),
	}
	env.InitializeInstanceVariablesWithMRO(obj)
	return obj
}

// Global THREAD class definitions - created once and reused
var threadClassesOnce = sync.Once{}
var threadClasses map[string]*environment.Class

func getThreadClasses() map[string]*environment.Class {
	threadClassesOnce.Do(func() {
		threadClasses = map[string]*environment.Class{
			"YARN": {
				Name: "YARN",
				Documentation: []string{
					"Abstract thread interface for creating concurrent execution.",
					"Must be subclassed to implement the SPIN method.",
					"Subclasses must call the YARN constructor to initialize.",
				},
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"YARN": {
						Name:       "YARN",
						Parameters: []environment.Parameter{},
						Documentation: []string{
							"Initializes a YARN thread instance.",
							"Subclasses must call this constructor to set up the thread.",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							// Initialize thread data
							threadData := &ThreadData{}
							this.NativeData = threadData
							return environment.NOTHIN, nil
						},
					},
					// Abstract SPIN method - must be overridden
					"SPIN": {
						Name:       "SPIN",
						Parameters: []environment.Parameter{},
						Documentation: []string{
							"Thread implementation method that must be overridden by subclasses.",
							"Contains the code that runs in the thread.",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							return environment.NOTHIN, runtime.Exception{Message: "SPIN method must be implemented by subclass"}
						},
					},
					// START method - launches the thread
					"START": {
						Name:       "START",
						Parameters: []environment.Parameter{},
						Documentation: []string{
							"Starts the thread execution by calling the SPIN method in a separate thread.",
							"Throws exception if thread is already running, sets RUNNING to YEZ",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							if threadData, ok := this.NativeData.(*ThreadData); ok {
								if threadData.goroutineRunning {
									return environment.NOTHIN, runtime.Exception{Message: "START: thread already running"}
								}

								// Create a forked interpreter for the thread
								threadData.interpreter = interpreter.Fork()
								threadData.goroutineRunning = true
								threadData.finished = false
								threadData.wg.Add(1)

								// Launch goroutine that calls the SPIN method
								go func() {
									defer threadData.wg.Done()
									defer func() {
										threadData.goroutineRunning = false
										threadData.finished = true
									}()

									threadData.result, threadData.err = threadData.interpreter.CallMemberFunction(this, "SPIN", []environment.Value{})
								}()

								return environment.NOTHIN, nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "START: invalid thread context"}
						},
					},
					// JOIN method - waits for thread completion
					"JOIN": {
						Name:       "JOIN",
						Parameters: []environment.Parameter{},
						Documentation: []string{
							"Waits for the thread to finish execution.",
							"Blocks until thread completes, returns any value returned by SPIN method, throws any exception thrown by SPIN method.",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							if threadData, ok := this.NativeData.(*ThreadData); ok {
								threadData.wg.Wait()

								// Return any error from the SPIN method
								if threadData.err != nil {
									return environment.NOTHIN, threadData.err
								}
								return threadData.result, nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "JOIN: invalid thread context"}
						},
					},
				},
				PublicVariables: map[string]*environment.MemberVariable{
					"RUNNING": {
						Variable: environment.Variable{
							Name:          "RUNNING",
							Type:          "BOOL",
							IsLocked:      true,
							IsPublic:      true,
							Documentation: []string{"YEZ if the thread is currently running."},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							if threadData, ok := this.NativeData.(*ThreadData); ok {
								return environment.BoolValue(threadData.goroutineRunning), nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "RUNNING: invalid thread context"}
						},
						NativeSet: nil, // Read-only
					},
					"FINISHED": {
						Variable: environment.Variable{
							Name:          "FINISHED",
							Type:          "BOOL",
							IsLocked:      true,
							IsPublic:      true,
							Documentation: []string{"YEZ if the thread has completed execution."},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							if threadData, ok := this.NativeData.(*ThreadData); ok {
								return environment.BoolValue(threadData.finished), nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "FINISHED: invalid thread context"}
						},
						NativeSet: nil, // Read-only
					},
				},
				QualifiedName:    "stdlib:THREAD.YARN",
				ModulePath:       "stdlib:THREAD",
				ParentClasses:    []string{},
				MRO:              []string{"stdlib:THREAD.YARN"},
				PrivateVariables: make(map[string]*environment.MemberVariable),
				PrivateFunctions: make(map[string]*environment.Function),
				SharedVariables:  make(map[string]*environment.MemberVariable),
				SharedFunctions:  make(map[string]*environment.Function),
			},
			"KNOT": {
				Name:          "KNOT",
				QualifiedName: "stdlib:THREAD.KNOT",
				ModulePath:    "stdlib:THREAD",
				ParentClasses: []string{},
				MRO:           []string{"stdlib:THREAD.KNOT"},
				Documentation: []string{
					"Provides mutual exclusion (mutex) functionality for synchronizing access to shared resources between threads.",
				},
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"KNOT": {
						Name:          "KNOT",
						Parameters:    []environment.Parameter{},
						Documentation: []string{"Initializes a KNOT mutex instance."},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							// Initialize mutex data
							mutexData := &MutexData{}
							this.NativeData = mutexData
							return environment.NOTHIN, nil
						},
					},
					// TIE method - locks the mutex
					"TIE": {
						Name:       "TIE",
						Parameters: []environment.Parameter{},
						Documentation: []string{
							"Acquires the mutex lock.",
							"If already locked by another thread, blocks until available.",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							if mutexData, ok := this.NativeData.(*MutexData); ok {
								mutexData.mutex.Lock()
								mutexData.locked = true
								return environment.NOTHIN, nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "TIE: invalid mutex context"}
						},
					},
					// UNTIE method - unlocks the mutex
					"UNTIE": {
						Name:       "UNTIE",
						Parameters: []environment.Parameter{},
						Documentation: []string{
							"Releases the mutex lock.",
							"Throws exception if mutex is not currently locked, only the thread that locked the mutex should unlock it.",
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							if mutexData, ok := this.NativeData.(*MutexData); ok {
								if !mutexData.locked {
									return environment.NOTHIN, runtime.Exception{Message: "UNTIE: cannot unlock mutex that is not locked"}
								}
								mutexData.mutex.Unlock()
								mutexData.locked = false
								return environment.NOTHIN, nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "UNTIE: invalid mutex context"}
						},
					},
				},
				PublicVariables: map[string]*environment.MemberVariable{
					"LOCKED": {
						Variable: environment.Variable{
							Name:          "LOCKED",
							Type:          "BOOL",
							IsLocked:      true,
							IsPublic:      true,
							Documentation: []string{"YEZ if the mutex is currently locked."},
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							if mutexData, ok := this.NativeData.(*MutexData); ok {
								return environment.BoolValue(mutexData.locked), nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "LOCKED: invalid mutex context"}
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
	return threadClasses
}

// RegisterTHREADInEnv registers THREAD classes in the given environment
// declarations: empty slice means import all, otherwise import only specified declarations
func RegisterTHREADInEnv(env *environment.Environment, declarations ...string) error {
	threadClasses := getThreadClasses()

	// If declarations is empty, import all classes
	if len(declarations) == 0 {
		for _, class := range threadClasses {
			env.DefineClass(class)
		}
		return nil
	}

	// Otherwise, import only specified declarations
	for _, decl := range declarations {
		declUpper := strings.ToUpper(decl)

		// Check if it's a class
		if class, exists := threadClasses[declUpper]; exists {
			env.DefineClass(class)
		} else {
			return runtime.Exception{Message: fmt.Sprintf("unknown THREAD declaration: %s", decl)}
		}
	}

	return nil
}
