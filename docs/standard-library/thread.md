# THREAD Module

## Import

```lol
BTW Full import
I CAN HAS THREAD?

BTW Selective import examples
```

## Miscellaneous

### KNOT Class

Mutual exclusion (mutex) class for synchronizing access to shared resources between threads.
Prevents race conditions by ensuring only one thread can access a resource at a time.
Use TIE to lock and UNTIE to unlock the mutex.

**Methods:**

#### KNOT

Initializes a KNOT mutex instance.
Creates an unlocked mutex ready for synchronization.

**Syntax:** `NEW KNOT`
**Example: Create mutex**

```lol
I HAS A VARIABLE MUTEX TEH KNOT ITZ NEW KNOT
BTW Mutex starts unlocked
```

**Note:** Mutex starts in unlocked state

**Note:** Use TIE to acquire lock, UNTIE to release

#### TIE

Acquires the mutex lock for exclusive access to shared resources.
Blocks the calling thread if another thread already holds the lock.
Sets LOCKED to YEZ when lock is acquired.

**Syntax:** `<mutex> DO TIE`
**Example: Acquire lock**

```lol
I HAS A VARIABLE MUTEX TEH KNOT ITZ NEW KNOT
MUTEX DO TIE
BTW Now have exclusive access
```

**Example: Protect critical section**

```lol
MUTEX DO TIE
SHARED_VAR TEH INTEGR ITZ SHARED_VAR MOAR 1
MUTEX DO UNTIE
BTW Critical section protected
```

**Note:** Blocks if lock is already held

**Note:** Only one thread can hold lock at a time

**Note:** Must call UNTIE to release lock

**Note:** Same thread must call both TIE and UNTIE

#### UNTIE

Releases the mutex lock, allowing other threads to acquire it.
Sets LOCKED to NO and wakes up any waiting threads.
Throws exception if mutex is not currently locked.

**Syntax:** `<mutex> DO UNTIE`
**Example: Release lock**

```lol
I HAS A VARIABLE MUTEX TEH KNOT ITZ NEW KNOT
MUTEX DO TIE
BTW Do some work here
MUTEX DO UNTIE
BTW Lock released, other threads can now acquire
```

**Example: Complete critical section**

```lol
MUTEX DO TIE
SHARED_DATA ITZ "updated value"
MUTEX DO UNTIE
BTW Critical section complete
```

**Note:** Must be called by same thread that called TIE

**Note:** Throws exception if mutex not locked

**Note:** Allows waiting threads to proceed

**Note:** Sets LOCKED to NO

**Member Variables:**

#### LOCKED

Indicates whether the mutex is currently locked.
YEZ when a thread holds the lock, NO when available.

**Type:** {BOOL}

**Example: Check lock status**

```lol
I HAS A VARIABLE MUTEX TEH KNOT ITZ NEW KNOT
IZ MUTEX LOCKED?
SAYZ WIT "Mutex is locked!"
NOPE
SAYZ WIT "Mutex is available"
KTHX
```

**Example: Wait for lock to be available**

```lol
WHILE YEZ
IZ NO SAEM AS MUTEX LOCKED?
OUTTA HERE
KTHX
KTHX
MUTEX DO TIE
```

**Note:** Read-only property

**Note:** YEZ while TIE is held

**Note:** NO when UNTIE is called or initially

**Note:** Use for polling lock status

**Example: Basic mutex usage**

```lol
HAI ME TEH VARIABLE MUTEX TEH KNOT ITZ NEW KNOT
HAI ME TEH VARIABLE SHARED_DATA TEH INTEGR ITZ 0

HAI ME TEH CLASS WORKER KITTEH OF YARN
DIS TEH FUNCSHUN SPIN
MUTEX DO TIE
SHARED_DATA ITZ SHARED_DATA MOAR 1
MUTEX DO UNTIE
KTHX
KTHX
```

**Example: Protecting shared resources**

```lol
DIS TEH VARIABLE LOCK TEH KNOT ITZ NEW KNOT
DIS TEH VARIABLE COUNTER TEH INTEGR ITZ 0

DIS TEH FUNCSHUN SAFE_INCREMENT
LOCK DO TIE
COUNTER ITZ COUNTER MOAR 1
LOCK DO UNTIE
GIVEZ COUNTER
KTHX
```

**Example: Multiple threads with synchronization**

```lol
HAI ME TEH VARIABLE MUTEX TEH KNOT ITZ NEW KNOT
HAI ME TEH VARIABLE TOTAL TEH INTEGR ITZ 0

HAI ME TEH CLASS ADDER KITTEH OF YARN
DIS TEH FUNCSHUN SPIN
I HAS A VARIABLE IDX TEH INTEGR ITZ 0
WHILE IDX SMALLR THAN 100
MUTEX DO TIE
TOTAL ITZ TOTAL MOAR 1
MUTEX DO UNTIE
KTHX
KTHX
KTHX
```

### YARN Class

Abstract thread class for creating concurrent execution in Objective-LOL.
Must be subclassed to implement the SPIN method with thread logic.
Provides START and JOIN methods for thread lifecycle management.

**Methods:**

#### JOIN

Waits for thread completion and returns the value from SPIN method.
Blocks the calling thread until the target thread finishes execution.
Returns the value returned by SPIN or throws any exception from SPIN.

**Syntax:** `<thread> DO JOIN`
**Example: Wait for thread completion**

```lol
I HAS A VARIABLE THREAD TEH MY_THREAD ITZ NEW MY_THREAD
THREAD DO START
THREAD DO JOIN
BTW Now thread has completed
```

**Example: Get thread result**

```lol
I HAS A VARIABLE CALC_THREAD TEH CALCULATOR_THREAD ITZ NEW CALCULATOR_THREAD
CALC_THREAD DO START
I HAS A VARIABLE RESULT TEH INTEGR ITZ CALC_THREAD DO JOIN
SAYZ WIT "Calculation result: " MOAR RESULT
```

**Example: Handle thread exceptions**

```lol
I HAS A VARIABLE THREAD TEH MY_THREAD ITZ NEW MY_THREAD
THREAD DO START
BTW If SPIN throws exception, JOIN will throw it here
THREAD DO JOIN
```

**Note:** Blocks until thread completes

**Note:** Returns value from SPIN method

**Note:** Propagates exceptions from SPIN method

**Note:** Thread becomes FINISHED after JOIN

#### SPIN

Abstract method that must be implemented by thread subclasses.
Contains the code that executes in the separate thread.
Return value is available to calling thread via JOIN method.

**Syntax:** `<thread> DO SPIN`
**Example: Simple thread logic**

```lol
DIS TEH FUNCSHUN SPIN
SAYZ WIT "Thread is running!"
KTHX
```

**Example: Thread with computation**

```lol
DIS TEH FUNCSHUN SPIN
I HAS A VARIABLE RESULT TEH INTEGR ITZ 10 MOAR 20
GIVEZ RESULT
KTHX
```

**Example: Long-running thread**

```lol
DIS TEH FUNCSHUN SPIN
I HAS A VARIABLE IDX TEH INTEGR ITZ 0
WHILE IDX SMALLR THAN 100
BTW Do some work here
I HAS A VARIABLE PROGRESS TEH INTEGR ITZ IDX
IDX ITZ IDX MOAR 1
KTHX
GIVEZ "Work complete"
KTHX
```

**Note:** Called automatically when thread starts

**Note:** Runs in separate goroutine from main thread

**Note:** Return value available via JOIN method

**Note:** Exceptions thrown here are propagated to JOIN

#### START

Starts thread execution by calling SPIN method in a separate goroutine.
Thread begins running concurrently with the calling thread.
Sets RUNNING to YEZ and FINISHED to NO.

**Syntax:** `<thread> DO START`
**Example: Start a thread**

```lol
I HAS A VARIABLE THREAD TEH MY_THREAD ITZ NEW MY_THREAD
THREAD DO START
BTW Thread is now running concurrently
```

**Example: Check thread status**

```lol
THREAD DO START
IZ THREAD RUNNING?
SAYZ WIT "Thread is running!"
KTHX
```

**Note:** Cannot start thread that's already running

**Note:** Thread execution begins immediately

**Note:** Use JOIN to wait for completion

**Note:** RUNNING becomes YEZ, FINISHED becomes NO

#### YARN

Initializes a YARN thread instance.
Must be called manually by thread subclasses.

**Syntax:** `NEW <ThreadClass>`
**Example: Thread creation**

```lol
I HAS A VARIABLE THREAD TEH MY_THREAD ITZ NEW MY_THREAD
BTW Constructor called automatically
```

**Note:** Usually called automatically by NEW operator

**Note:** Initializes internal thread state

**Member Variables:**

#### FINISHED

Indicates whether the thread has completed execution.
YEZ after thread finishes, NO while running or before start.

**Type:** {BOOL}

**Example: Wait for completion**

```lol
I HAS A VARIABLE THREAD TEH MY_THREAD ITZ NEW MY_THREAD
THREAD DO START
WHILE YEZ
IZ THREAD FINISHED?
OUTTA HERE
KTHX
KTHX
SAYZ WIT "Thread completed!"
```

**Note:** Read-only property

**Note:** YEZ after SPIN method completes

**Note:** NO during execution or before START

**Note:** Use JOIN for blocking wait

#### RUNNING

Indicates whether the thread is currently executing.
YEZ while thread is running, NO when stopped or finished.

**Type:** {BOOL}

**Example: Check thread status**

```lol
I HAS A VARIABLE THREAD TEH MY_THREAD ITZ NEW MY_THREAD
THREAD DO START
IZ THREAD RUNNING?
SAYZ WIT "Thread is active!"
KTHX
```

**Note:** Read-only property

**Note:** YEZ during SPIN method execution

**Note:** NO before START or after completion

**Example: Basic thread subclass**

```lol
HAI ME TEH CLASS MY_THREAD KITTEH OF YARN
DIS TEH FUNCSHUN MY_THREAD
YARN
KTHX
DIS TEH FUNCSHUN SPIN
SAYZ WIT "Hello from thread!"
KTHX
KTHX

I HAS A VARIABLE THREAD TEH MY_THREAD ITZ NEW MY_THREAD
THREAD DO START
THREAD DO JOIN
```

**Example: Thread with return value**

```lol
HAI ME TEH CLASS CALCULATOR_THREAD KITTEH OF YARN
DIS TEH FUNCSHUN CALCULATOR_THREAD
YARN
KTHX
DIS TEH FUNCSHUN SPIN
GIVEZ SUM WIT 10 AN WIT 20
KTHX
KTHX

I HAS A VARIABLE CALC TEH CALCULATOR_THREAD ITZ NEW CALCULATOR_THREAD
CALC DO START
I HAS A VARIABLE RESULT TEH INTEGR ITZ CALC DO JOIN
SAYZ WIT "Result: " MOAR RESULT
```

