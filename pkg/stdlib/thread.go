package stdlib

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/runtime"
)

// moduleThreadCategories defines the order that categories should be rendered in documentation
var moduleThreadCategories = []string{
	"threading",
	"synchronization",
}

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
					"Abstract thread class for creating concurrent execution in Objective-LOL.",
					"Must be subclassed to implement the SPIN method with thread logic.",
					"Provides START and JOIN methods for thread lifecycle management.",
					"",
					"@class YARN",
					"@abstract",
					"@example Basic thread subclass",
					"HAI ME TEH CLASS MY_THREAD KITTEH OF YARN",
					"    DIS TEH FUNCSHUN MY_THREAD",
					"        YARN",
					"    KTHX",
					"    DIS TEH FUNCSHUN SPIN",
					"        SAYZ WIT \"Hello from thread!\"",
					"    KTHX",
					"KTHX",
					"",
					"I HAS A VARIABLE THREAD TEH MY_THREAD ITZ NEW MY_THREAD",
					"THREAD DO START",
					"THREAD DO JOIN",
					"@example Thread with return value",
					"HAI ME TEH CLASS CALCULATOR_THREAD KITTEH OF YARN",
					"    DIS TEH FUNCSHUN CALCULATOR_THREAD",
					"        YARN",
					"    KTHX",
					"    DIS TEH FUNCSHUN SPIN",
					"        GIVEZ SUM WIT 10 AN WIT 20",
					"    KTHX",
					"KTHX",
					"",
					"I HAS A VARIABLE CALC TEH CALCULATOR_THREAD ITZ NEW CALCULATOR_THREAD",
					"CALC DO START",
					"I HAS A VARIABLE RESULT TEH INTEGR ITZ CALC DO JOIN",
					"SAYZ WIT \"Result: \" MOAR RESULT",
					"@note Must implement SPIN method in subclasses",
					"@note Call START to begin execution, JOIN to wait for completion",
					"@note Thread runs concurrently with main program",
					"@note Use KNOT for synchronizing access to shared resources",
					"@see KNOT, START, JOIN, SPIN",
				},
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"YARN": {
						Name:       "YARN",
						Parameters: []environment.Parameter{},
						Documentation: []string{
							"Initializes a YARN thread instance.",
							"Must be called manually by thread subclasses.",
							"",
							"@syntax NEW <ThreadClass>",
							"@returns {NOTHIN} No return value (constructor)",
							"@example Thread creation",
							"I HAS A VARIABLE THREAD TEH MY_THREAD ITZ NEW MY_THREAD",
							"BTW Constructor called automatically",
							"@note Usually called automatically by NEW operator",
							"@note Initializes internal thread state",
							"@category threading",
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
							"Abstract method that must be implemented by thread subclasses.",
							"Contains the code that executes in the separate thread.",
							"Return value is available to calling thread via JOIN method.",
							"",
							"@syntax <thread> DO SPIN",
							"@returns {} Any value to return to joining thread",
							"@example Simple thread logic",
							"DIS TEH FUNCSHUN SPIN",
							"    SAYZ WIT \"Thread is running!\"",
							"KTHX",
							"@example Thread with computation",
							"DIS TEH FUNCSHUN SPIN",
							"    I HAS A VARIABLE RESULT TEH INTEGR ITZ 10 MOAR 20",
							"    GIVEZ RESULT",
							"KTHX",
							"@example Long-running thread",
							"DIS TEH FUNCSHUN SPIN",
							"    I HAS A VARIABLE IDX TEH INTEGR ITZ 0",
							"    WHILE IDX SMALLR THAN 100",
							"        BTW Do some work here",
							"        I HAS A VARIABLE PROGRESS TEH INTEGR ITZ IDX",
							"        IDX ITZ IDX MOAR 1",
							"    KTHX",
							"    GIVEZ \"Work complete\"",
							"KTHX",
							"@note Called automatically when thread starts",
							"@note Runs in separate goroutine from main thread",
							"@note Return value available via JOIN method",
							"@note Exceptions thrown here are propagated to JOIN",
							"@category threading",
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
							"Starts thread execution by calling SPIN method in a separate goroutine.",
							"Thread begins running concurrently with the calling thread.",
							"Sets RUNNING to YEZ and FINISHED to NO.",
							"",
							"@syntax <thread> DO START",
							"@returns {NOTHIN} No return value",
							"@example Start a thread",
							"I HAS A VARIABLE THREAD TEH MY_THREAD ITZ NEW MY_THREAD",
							"THREAD DO START",
							"BTW Thread is now running concurrently",
							"@example Check thread status",
							"THREAD DO START",
							"IZ THREAD RUNNING?",
							"    SAYZ WIT \"Thread is running!\"",
							"KTHX",
							"@note Cannot start thread that's already running",
							"@note Thread execution begins immediately",
							"@note Use JOIN to wait for completion",
							"@note RUNNING becomes YEZ, FINISHED becomes NO",
							"@category threading",
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
							"Waits for thread completion and returns the value from SPIN method.",
							"Blocks the calling thread until the target thread finishes execution.",
							"Returns the value returned by SPIN or throws any exception from SPIN.",
							"",
							"@syntax <thread> DO JOIN",
							"@returns {*} Value returned by SPIN method",
							"@example Wait for thread completion",
							"I HAS A VARIABLE THREAD TEH MY_THREAD ITZ NEW MY_THREAD",
							"THREAD DO START",
							"THREAD DO JOIN",
							"BTW Now thread has completed",
							"@example Get thread result",
							"I HAS A VARIABLE CALC_THREAD TEH CALCULATOR_THREAD ITZ NEW CALCULATOR_THREAD",
							"CALC_THREAD DO START",
							"I HAS A VARIABLE RESULT TEH INTEGR ITZ CALC_THREAD DO JOIN",
							"SAYZ WIT \"Calculation result: \" MOAR RESULT",
							"@example Handle thread exceptions",
							"I HAS A VARIABLE THREAD TEH MY_THREAD ITZ NEW MY_THREAD",
							"THREAD DO START",
							"BTW If SPIN throws exception, JOIN will throw it here",
							"THREAD DO JOIN",
							"@note Blocks until thread completes",
							"@note Returns value from SPIN method",
							"@note Propagates exceptions from SPIN method",
							"@note Thread becomes FINISHED after JOIN",
							"@category threading",
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
							Name:     "RUNNING",
							Type:     "BOOL",
							IsLocked: true,
							IsPublic: true,
							Documentation: []string{
								"Indicates whether the thread is currently executing.",
								"YEZ while thread is running, NO when stopped or finished.",
								"",
								"@type {BOOL}",
								"@readonly",
								"@example Check thread status",
								"I HAS A VARIABLE THREAD TEH MY_THREAD ITZ NEW MY_THREAD",
								"THREAD DO START",
								"IZ THREAD RUNNING?",
								"    SAYZ WIT \"Thread is active!\"",
								"KTHX",
								"@note Read-only property",
								"@note YEZ during SPIN method execution",
								"@note NO before START or after completion",
								"@category threading",
							},
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
							Name:     "FINISHED",
							Type:     "BOOL",
							IsLocked: true,
							IsPublic: true,
							Documentation: []string{
								"Indicates whether the thread has completed execution.",
								"YEZ after thread finishes, NO while running or before start.",
								"",
								"@type {BOOL}",
								"@readonly",
								"@example Wait for completion",
								"I HAS A VARIABLE THREAD TEH MY_THREAD ITZ NEW MY_THREAD",
								"THREAD DO START",
								"WHILE YEZ",
								"    IZ THREAD FINISHED?",
								"        OUTTA HERE",
								"    KTHX",
								"KTHX",
								"SAYZ WIT \"Thread completed!\"",
								"@note Read-only property",
								"@note YEZ after SPIN method completes",
								"@note NO during execution or before START",
								"@note Use JOIN for blocking wait",
								"@category threading",
							},
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
					"Mutual exclusion (mutex) class for synchronizing access to shared resources between threads.",
					"Prevents race conditions by ensuring only one thread can access a resource at a time.",
					"Use TIE to lock and UNTIE to unlock the mutex.",
					"",
					"@class KNOT",
					"@example Basic mutex usage",
					"HAI ME TEH VARIABLE MUTEX TEH KNOT ITZ NEW KNOT",
					"HAI ME TEH VARIABLE SHARED_DATA TEH INTEGR ITZ 0",
					"",
					"HAI ME TEH CLASS WORKER KITTEH OF YARN",
					"    DIS TEH FUNCSHUN SPIN",
					"        MUTEX DO TIE",
					"        SHARED_DATA ITZ SHARED_DATA MOAR 1",
					"        MUTEX DO UNTIE",
					"    KTHX",
					"KTHX",
					"@example Protecting shared resources",
					"DIS TEH VARIABLE LOCK TEH KNOT ITZ NEW KNOT",
					"DIS TEH VARIABLE COUNTER TEH INTEGR ITZ 0",
					"",
					"DIS TEH FUNCSHUN SAFE_INCREMENT",
					"    LOCK DO TIE",
					"    COUNTER ITZ COUNTER MOAR 1",
					"    LOCK DO UNTIE",
					"    GIVEZ COUNTER",
					"KTHX",
					"@example Multiple threads with synchronization",
					"HAI ME TEH VARIABLE MUTEX TEH KNOT ITZ NEW KNOT",
					"HAI ME TEH VARIABLE TOTAL TEH INTEGR ITZ 0",
					"",
					"HAI ME TEH CLASS ADDER KITTEH OF YARN",
					"    DIS TEH FUNCSHUN SPIN",
					"        I HAS A VARIABLE IDX TEH INTEGR ITZ 0",
					"        WHILE IDX SMALLR THAN 100",
					"            MUTEX DO TIE",
					"            TOTAL ITZ TOTAL MOAR 1",
					"            MUTEX DO UNTIE",
					"        KTHX",
					"    KTHX",
					"KTHX",
					"@note Only one thread can hold the lock at a time",
					"@note TIE blocks if another thread holds the lock",
					"@note Same thread must call UNTIE that called TIE",
					"@note Use for protecting shared variables and resources",
					"@see YARN, TIE, UNTIE",
				},
				PublicFunctions: map[string]*environment.Function{
					// Constructor
					"KNOT": {
						Name:       "KNOT",
						Parameters: []environment.Parameter{},
						Documentation: []string{
							"Initializes a KNOT mutex instance.",
							"Creates an unlocked mutex ready for synchronization.",
							"",
							"@syntax NEW KNOT",
							"@returns {NOTHIN} No return value (constructor)",
							"@example Create mutex",
							"I HAS A VARIABLE MUTEX TEH KNOT ITZ NEW KNOT",
							"BTW Mutex starts unlocked",
							"@note Mutex starts in unlocked state",
							"@note Use TIE to acquire lock, UNTIE to release",
							"@category synchronization",
						},
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
							"Acquires the mutex lock for exclusive access to shared resources.",
							"Blocks the calling thread if another thread already holds the lock.",
							"Sets LOCKED to YEZ when lock is acquired.",
							"",
							"@syntax <mutex> DO TIE",
							"@returns {NOTHIN} No return value",
							"@example Acquire lock",
							"I HAS A VARIABLE MUTEX TEH KNOT ITZ NEW KNOT",
							"MUTEX DO TIE",
							"BTW Now have exclusive access",
							"@example Protect critical section",
							"MUTEX DO TIE",
							"SHARED_VAR TEH INTEGR ITZ SHARED_VAR MOAR 1",
							"MUTEX DO UNTIE",
							"BTW Critical section protected",
							"@note Blocks if lock is already held",
							"@note Only one thread can hold lock at a time",
							"@note Must call UNTIE to release lock",
							"@note Same thread must call both TIE and UNTIE",
							"@category synchronization",
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
							"Releases the mutex lock, allowing other threads to acquire it.",
							"Sets LOCKED to NO and wakes up any waiting threads.",
							"Throws exception if mutex is not currently locked.",
							"",
							"@syntax <mutex> DO UNTIE",
							"@returns {NOTHIN} No return value",
							"@example Release lock",
							"I HAS A VARIABLE MUTEX TEH KNOT ITZ NEW KNOT",
							"MUTEX DO TIE",
							"BTW Do some work here",
							"MUTEX DO UNTIE",
							"BTW Lock released, other threads can now acquire",
							"@example Complete critical section",
							"MUTEX DO TIE",
							"SHARED_DATA ITZ \"updated value\"",
							"MUTEX DO UNTIE",
							"BTW Critical section complete",
							"@note Must be called by same thread that called TIE",
							"@note Throws exception if mutex not locked",
							"@note Allows waiting threads to proceed",
							"@note Sets LOCKED to NO",
							"@category synchronization",
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
							Name:     "LOCKED",
							Type:     "BOOL",
							IsLocked: true,
							IsPublic: true,
							Documentation: []string{
								"Indicates whether the mutex is currently locked.",
								"YEZ when a thread holds the lock, NO when available.",
								"",
								"@type {BOOL}",
								"@readonly",
								"@example Check lock status",
								"I HAS A VARIABLE MUTEX TEH KNOT ITZ NEW KNOT",
								"IZ MUTEX LOCKED?",
								"    SAYZ WIT \"Mutex is locked!\"",
								"NOPE",
								"    SAYZ WIT \"Mutex is available\"",
								"KTHX",
								"@example Wait for lock to be available",
								"WHILE YEZ",
								"    IZ NO SAEM AS MUTEX LOCKED?",
								"        OUTTA HERE",
								"    KTHX",
								"KTHX",
								"MUTEX DO TIE",
								"@note Read-only property",
								"@note YEZ while TIE is held",
								"@note NO when UNTIE is called or initially",
								"@note Use for polling lock status",
								"@category synchronization",
							},
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
