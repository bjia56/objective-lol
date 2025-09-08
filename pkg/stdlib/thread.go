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
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"YARN": {
						Name:       "YARN",
						Parameters: []environment.Parameter{},
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
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							return environment.NOTHIN, runtime.Exception{Message: "SPIN method must be implemented by subclass"}
						},
					},
					// START method - launches the thread
					"START": {
						Name:       "START",
						Parameters: []environment.Parameter{},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							if threadData, ok := this.NativeData.(*ThreadData); ok {
								if threadData.goroutineRunning {
									return environment.NOTHIN, runtime.Exception{Message: "Thread already running"}
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
							return environment.NOTHIN, fmt.Errorf("START: invalid thread context")
						},
					},
					// JOIN method - waits for thread completion
					"JOIN": {
						Name:       "JOIN",
						Parameters: []environment.Parameter{},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							if threadData, ok := this.NativeData.(*ThreadData); ok {
								threadData.wg.Wait()

								// Return any error from the SPIN method
								if threadData.err != nil {
									return environment.NOTHIN, threadData.err
								}
								return threadData.result, nil
							}
							return environment.NOTHIN, fmt.Errorf("JOIN: invalid thread context")
						},
					},
				},
				PublicVariables: map[string]*environment.MemberVariable{
					"RUNNING": {
						Variable: environment.Variable{
							Name:     "RUNNING",
							Type:     "BOOL",
							IsLocked: true,
							IsPublic: true,
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							if threadData, ok := this.NativeData.(*ThreadData); ok {
								return environment.BoolValue(threadData.goroutineRunning), nil
							}
							return environment.NOTHIN, fmt.Errorf("RUNNING: invalid thread context")
						},
						NativeSet: nil, // Read-only
					},
					"FINISHED": {
						Variable: environment.Variable{
							Name:     "FINISHED",
							Type:     "BOOL",
							IsLocked: true,
							IsPublic: true,
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							if threadData, ok := this.NativeData.(*ThreadData); ok {
								return environment.BoolValue(threadData.finished), nil
							}
							return environment.NOTHIN, fmt.Errorf("FINISHED: invalid thread context")
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
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"KNOT": {
						Name:       "KNOT",
						Parameters: []environment.Parameter{},
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
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							if mutexData, ok := this.NativeData.(*MutexData); ok {
								mutexData.mutex.Lock()
								mutexData.locked = true
								return environment.NOTHIN, nil
							}
							return environment.NOTHIN, fmt.Errorf("TIE: invalid mutex context")
						},
					},
					// UNTIE method - unlocks the mutex
					"UNTIE": {
						Name:       "UNTIE",
						Parameters: []environment.Parameter{},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							if mutexData, ok := this.NativeData.(*MutexData); ok {
								if !mutexData.locked {
									return environment.NOTHIN, runtime.Exception{Message: "Cannot unlock mutex that is not locked"}
								}
								mutexData.mutex.Unlock()
								mutexData.locked = false
								return environment.NOTHIN, nil
							}
							return environment.NOTHIN, fmt.Errorf("UNTIE: invalid mutex context")
						},
					},
				},
				PublicVariables: map[string]*environment.MemberVariable{
					"LOCKED": {
						Variable: environment.Variable{
							Name:     "LOCKED",
							Type:     "BOOL",
							IsLocked: true,
							IsPublic: true,
						},
						NativeGet: func(this *environment.ObjectInstance) (environment.Value, error) {
							if mutexData, ok := this.NativeData.(*MutexData); ok {
								return environment.BoolValue(mutexData.locked), nil
							}
							return environment.NOTHIN, fmt.Errorf("LOCKED: invalid mutex context")
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
			return fmt.Errorf("unknown THREAD declaration: %s", decl)
		}
	}

	return nil
}
