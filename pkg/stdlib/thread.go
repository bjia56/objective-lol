package stdlib

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/interpreter"
	"github.com/bjia56/objective-lol/pkg/runtime"
	"github.com/bjia56/objective-lol/pkg/types"
)

// ThreadData holds the native Go threading constructs for YARN objects
type ThreadData struct {
	goroutineRunning bool
	finished         bool
	wg               sync.WaitGroup
	functionCtx      *interpreter.FunctionContext
	result           types.Value
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
	return &environment.ObjectInstance{
		Class:      class,
		NativeData: &ThreadData{},
		Hierarchy:  []string{"YARN"},
		Variables: map[string]*environment.Variable{
			"RUNNING": {
				Name:     "RUNNING",
				Type:     "BOOL",
				Value:    types.NO,
				IsLocked: true,
				IsPublic: true,
			},
			"FINISHED": {
				Name:     "FINISHED",
				Type:     "BOOL",
				Value:    types.NO,
				IsLocked: true,
				IsPublic: true,
			},
		},
	}
}

// NewKnotInstance creates a new KNOT mutex object instance
func NewKnotInstance() *environment.ObjectInstance {
	class := getThreadClasses()["KNOT"]
	return &environment.ObjectInstance{
		Class:      class,
		NativeData: &MutexData{},
		Hierarchy:  []string{"KNOT"},
		Variables: map[string]*environment.Variable{
			"LOCKED": {
				Name:     "LOCKED",
				Type:     "BOOL",
				Value:    types.NO,
				IsLocked: true,
				IsPublic: true,
			},
		},
	}
}

// updateYarnStatus updates the status variables of a YARN object
func updateYarnStatus(obj *environment.ObjectInstance, threadData *ThreadData) {
	if runningVar, exists := obj.Variables["RUNNING"]; exists {
		runningVar.Value = types.BoolValue(threadData.goroutineRunning)
	}
	if finishedVar, exists := obj.Variables["FINISHED"]; exists {
		finishedVar.Value = types.BoolValue(threadData.finished)
	}
}

// updateKnotStatus updates the status variables of a KNOT object
func updateKnotStatus(obj *environment.ObjectInstance, mutexData *MutexData) {
	if lockedVar, exists := obj.Variables["LOCKED"]; exists {
		lockedVar.Value = types.BoolValue(mutexData.locked)
	}
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
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							// Initialize thread data
							threadData := &ThreadData{}
							this.NativeData = threadData
							updateYarnStatus(this, threadData)
							return types.NOTHIN, nil
						},
					},
					// Abstract SPIN method - must be overridden
					"SPIN": {
						Name:       "SPIN",
						Parameters: []environment.Parameter{},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							return types.NOTHIN, runtime.Exception{Message: "SPIN method must be implemented by subclass"}
						},
					},
					// START method - launches the thread
					"START": {
						Name:       "START",
						Parameters: []environment.Parameter{},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							if _, ok := this.NativeData.(*ThreadData); !ok {
								// Initialize thread data
								threadData := &ThreadData{}
								this.NativeData = threadData
								updateYarnStatus(this, threadData)
							}

							if threadData, ok := this.NativeData.(*ThreadData); ok {
								if threadData.goroutineRunning {
									return types.NOTHIN, runtime.Exception{Message: "Thread already running"}
								}

								// Get the interpreter from the context
								if funCtx, ok := ctx.(*interpreter.FunctionContext); ok {
									// Create a forked interpreter for the thread
									threadData.functionCtx = funCtx.Fork()
									threadData.goroutineRunning = true
									threadData.finished = false
									threadData.wg.Add(1)

									updateYarnStatus(this, threadData)

									// Launch goroutine that calls the SPIN method
									go func() {
										defer threadData.wg.Done()
										defer func() {
											threadData.goroutineRunning = false
											threadData.finished = true
											updateYarnStatus(this, threadData)
										}()

										threadData.result, threadData.err = threadData.functionCtx.CallMethod(this, "SPIN", "YARN", []types.Value{})
									}()

									return types.NOTHIN, nil
								}
								return types.NOTHIN, fmt.Errorf("START: invalid interpreter context")
							}
							return types.NOTHIN, fmt.Errorf("START: invalid thread context")
						},
					},
					// JOIN method - waits for thread completion
					"JOIN": {
						Name:       "JOIN",
						Parameters: []environment.Parameter{},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							if threadData, ok := this.NativeData.(*ThreadData); ok {
								threadData.wg.Wait()
								updateYarnStatus(this, threadData)

								// Return any error from the SPIN method
								if threadData.err != nil {
									return types.NOTHIN, threadData.err
								}
								return threadData.result, nil
							}
							return types.NOTHIN, fmt.Errorf("JOIN: invalid thread context")
						},
					},
				},
				PublicVariables: map[string]*environment.Variable{
					"RUNNING": {
						Name:     "RUNNING",
						Type:     "BOOL",
						Value:    types.NO,
						IsLocked: true,
						IsPublic: true,
					},
					"FINISHED": {
						Name:     "FINISHED",
						Type:     "BOOL",
						Value:    types.NO,
						IsLocked: true,
						IsPublic: true,
					},
				},
				QualifiedName: "stdlib:THREAD.YARN",
				ModulePath:    "stdlib:THREAD",
				ParentClass:   "",
				PrivateVariables: make(map[string]*environment.Variable),
				PrivateFunctions: make(map[string]*environment.Function),
				SharedVariables:  make(map[string]*environment.Variable),
				SharedFunctions:  make(map[string]*environment.Function),
			},
			"KNOT": {
				Name:          "KNOT",
				QualifiedName: "stdlib:THREAD.KNOT",
				ModulePath:    "stdlib:THREAD",
				ParentClass:   "",
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"KNOT": {
						Name:       "KNOT",
						Parameters: []environment.Parameter{},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							// Initialize mutex data
							mutexData := &MutexData{}
							this.NativeData = mutexData
							updateKnotStatus(this, mutexData)
							return types.NOTHIN, nil
						},
					},
					// TIE method - locks the mutex
					"TIE": {
						Name:       "TIE",
						Parameters: []environment.Parameter{},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							if mutexData, ok := this.NativeData.(*MutexData); ok {
								mutexData.mutex.Lock()
								mutexData.locked = true
								updateKnotStatus(this, mutexData)
								return types.NOTHIN, nil
							}
							return types.NOTHIN, fmt.Errorf("TIE: invalid mutex context")
						},
					},
					// UNTIE method - unlocks the mutex
					"UNTIE": {
						Name:       "UNTIE",
						Parameters: []environment.Parameter{},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							if mutexData, ok := this.NativeData.(*MutexData); ok {
								if !mutexData.locked {
									return types.NOTHIN, runtime.Exception{Message: "Cannot unlock mutex that is not locked"}
								}
								mutexData.mutex.Unlock()
								mutexData.locked = false
								updateKnotStatus(this, mutexData)
								return types.NOTHIN, nil
							}
							return types.NOTHIN, fmt.Errorf("UNTIE: invalid mutex context")
						},
					},
				},
				PublicVariables: map[string]*environment.Variable{
					"LOCKED": {
						Name:     "LOCKED",
						Type:     "BOOL",
						Value:    types.NO,
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
