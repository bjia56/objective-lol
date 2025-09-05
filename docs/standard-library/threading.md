# THREAD Module - Concurrency and Synchronization

The THREAD module provides concurrency support through thread objects (YARN) and mutex synchronization (KNOT) for building multi-threaded applications.

## Importing THREAD Module

```lol
BTW Import entire module
I CAN HAS THREAD?

BTW Selective imports
I CAN HAS YARN FROM THREAD?
I CAN HAS KNOT FROM THREAD?
```

## YARN Class - Thread Objects

The YARN class provides an abstract thread interface for creating concurrent execution. YARN is an abstract class that must be subclassed to implement the SPIN method.

### YARN Properties

- **RUNNING**: BOOL (read-only) - True if the thread is currently running
- **FINISHED**: BOOL (read-only) - True if the thread has completed execution

### YARN Methods

#### Constructor

```lol
I CAN HAS THREAD?

BTW YARN is abstract - create a subclass instead
HAI ME TEH CLAS MY_THREAD KITTEH OF YARN
    EVRYONE
    DIS TEH FUNCSHUN SPIN
        BTW Your thread code goes here
        SAYZ WIT "Thread is running!"
    KTHX
KTHXBAI

I HAS A VARIABLE THREAD TEH MY_THREAD ITZ NEW MY_THREAD
```

#### START - Launch Thread

Starts the thread execution by calling the SPIN method in a separate goroutine.

```lol
thread DO START
```

- Throws exception if thread is already running
- Sets RUNNING to YEZ
- Calls SPIN method concurrently

#### JOIN - Wait for Completion

Waits for the thread to finish execution.

```lol
thread DO JOIN
```

- Blocks until thread completes
- Returns any value returned by SPIN method
- Throws any exception thrown by SPIN method

#### SPIN - Thread Implementation (Abstract)

The SPIN method must be implemented by subclasses and contains the code that runs in the thread.

```lol
DIS TEH FUNCSHUN SPIN
    BTW Your thread implementation
KTHX
```

### Basic YARN Example

```lol
I CAN HAS THREAD?
I CAN HAS STDIO?

BTW Create a simple counter thread
HAI ME TEH CLAS COUNTER_THREAD KITTEH OF YARN
    EVRYONE
    DIS TEH VARIABLE COUNT TEH INTEGR ITZ 0
    DIS TEH VARIABLE MAX_COUNT TEH INTEGR ITZ 5

    DIS TEH FUNCSHUN COUNTER_THREAD WIT MAX TEH INTEGR
        MAX_COUNT ITZ MAX
    KTHX

    DIS TEH FUNCSHUN SPIN
        WHILE COUNT SMALLR THAN MAX_COUNT
            COUNT ITZ COUNT MOAR 1
            SAY WIT "Counter: "
            SAYZ WIT COUNT
            SLEEP WIT 1  BTW Simulate work (would need TIME module)
        KTHX
        SAYZ WIT "Counter thread finished!"
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN DEMO_BASIC_THREAD
    SAYZ WIT "=== Basic Thread Demo ==="

    BTW Create thread
    I HAS A VARIABLE COUNTER TEH COUNTER_THREAD ITZ NEW COUNTER_THREAD WIT 3

    BTW Check initial state
    SAY WIT "Before start - Running: "
    SAY WIT COUNTER RUNNING
    SAY WIT ", Finished: "
    SAYZ WIT COUNTER FINISHED

    BTW Start the thread
    COUNTER DO START
    SAYZ WIT "Thread started!"

    BTW Check running state
    SAY WIT "After start - Running: "
    SAYZ WIT COUNTER RUNNING

    BTW Wait for completion
    COUNTER DO JOIN
    SAYZ WIT "Thread joined!"

    BTW Check final state
    SAY WIT "After join - Running: "
    SAY WIT COUNTER RUNNING
    SAY WIT ", Finished: "
    SAYZ WIT COUNTER FINISHED
KTHXBAI
```

### Advanced YARN Examples

#### Producer-Consumer Pattern

```lol
I CAN HAS THREAD?
I CAN HAS STDIO?

BTW Shared data structure (simplified)
HAI ME TEH VARIABLE SHARED_DATA TEH BUKKIT ITZ NEW BUKKIT

BTW Producer thread
HAI ME TEH CLAS PRODUCER_THREAD KITTEH OF YARN
    EVRYONE
    DIS TEH VARIABLE ITEM_COUNT TEH INTEGR ITZ 5

    DIS TEH FUNCSHUN SPIN
        I HAS A VARIABLE N TEH INTEGR ITZ 1
        WHILE N SMALLR THAN ITEM_COUNT MOAR 1
            SHARED_DATA DO PUSH WIT I
            SAY WIT "Produced item: "
            SAYZ WIT I
            N ITZ N MOAR 1
            SLEEP WIT 1  BTW Simulate production time
        KTHX
        SAYZ WIT "Producer finished!"
    KTHX
KTHXBAI

BTW Consumer thread
HAI ME TEH CLAS CONSUMER_THREAD KITTEH OF YARN
    EVRYONE
    DIS TEH VARIABLE CONSUMED_COUNT TEH INTEGR ITZ 0

    DIS TEH FUNCSHUN SPIN
        WHILE CONSUMED_COUNT SMALLR THAN 5
            IZ SHARED_DATA SIZ BIGGR THAN 0?
                I HAS A VARIABLE ITEM TEH INTEGR ITZ SHARED_DATA DO POP
                SAY WIT "Consumed item: "
                SAYZ WIT ITEM
                CONSUMED_COUNT ITZ CONSUMED_COUNT MOAR 1
            NOPE
                SLEEP WIT 1  BTW Wait for items
            KTHX
        KTHX
        SAYZ WIT "Consumer finished!"
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN DEMO_PRODUCER_CONSUMER
    SAYZ WIT "=== Producer-Consumer Demo ==="

    BTW Create threads
    I HAS A VARIABLE PRODUCER TEH PRODUCER_THREAD ITZ NEW PRODUCER_THREAD
    I HAS A VARIABLE CONSUMER TEH CONSUMER_THREAD ITZ NEW CONSUMER_THREAD

    BTW Start both threads
    PRODUCER DO START
    CONSUMER DO START
    SAYZ WIT "Both threads started"

    BTW Wait for both to complete
    PRODUCER DO JOIN
    SAYZ WIT "Producer completed"

    CONSUMER DO JOIN
    SAYZ WIT "Consumer completed"

    SAYZ WIT "Producer-Consumer demo finished"
KTHXBAI
```

## KNOT Class - Mutex Objects

The KNOT class provides mutual exclusion (mutex) functionality for synchronizing access to shared resources between threads.

### KNOT Properties

- **LOCKED**: BOOL (read-only) - True if the mutex is currently locked

### KNOT Methods

#### Constructor

```lol
I CAN HAS THREAD?

I HAS A VARIABLE MUTEX TEH KNOT ITZ NEW KNOT
```

#### TIE - Lock Mutex

Acquires the mutex lock. If already locked by another thread, blocks until available.

```lol
mutex DO TIE
```

#### UNTIE - Unlock Mutex

Releases the mutex lock.

```lol
mutex DO UNTIE
```

- Throws exception if mutex is not currently locked
- Only the thread that locked the mutex should unlock it

### Basic KNOT Example

```lol
I CAN HAS THREAD?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN DEMO_BASIC_MUTEX
    SAYZ WIT "=== Basic Mutex Demo ==="

    BTW Create mutex
    I HAS A VARIABLE MUTEX TEH KNOT ITZ NEW KNOT

    BTW Check initial state
    SAY WIT "Initial state - Locked: "
    SAYZ WIT MUTEX LOCKED

    BTW Lock the mutex
    MUTEX DO TIE
    SAY WIT "After TIE - Locked: "
    SAYZ WIT MUTEX LOCKED

    BTW Unlock the mutex
    MUTEX DO UNTIE
    SAY WIT "After UNTIE - Locked: "
    SAYZ WIT MUTEX LOCKED

    SAYZ WIT "Basic mutex demo completed"
KTHXBAI
```

### Synchronized Access Pattern

```lol
I CAN HAS THREAD?
I CAN HAS STDIO?

BTW Shared counter with mutex protection
HAI ME TEH VARIABLE SHARED_COUNTER TEH INTEGR ITZ 0
HAI ME TEH VARIABLE COUNTER_MUTEX TEH KNOT ITZ NEW KNOT

BTW Thread that increments counter safely
HAI ME TEH CLAS SAFE_COUNTER_THREAD KITTEH OF YARN
    EVRYONE
    DIS TEH VARIABLE THREAD_ID TEH STRIN ITZ "Unknown"
    DIS TEH VARIABLE INCREMENT_COUNT TEH INTEGR ITZ 3

    DIS TEH FUNCSHUN SAFE_COUNTER_THREAD WIT ID TEH STRIN
        THREAD_ID ITZ ID
    KTHX

    DIS TEH FUNCSHUN SPIN
        I HAS A VARIABLE I TEH INTEGR ITZ 0
        WHILE I SMALLR THAN INCREMENT_COUNT
            BTW Lock mutex for thread-safe access
            COUNTER_MUTEX DO TIE

            BTW Critical section - modify shared resource
            I HAS A VARIABLE OLD_VALUE TEH INTEGR ITZ SHARED_COUNTER
            SHARED_COUNTER ITZ SHARED_COUNTER MOAR 1

            SAY WIT THREAD_ID
            SAY WIT " incremented counter from "
            SAY WIT OLD_VALUE
            SAY WIT " to "
            SAYZ WIT SHARED_COUNTER

            BTW Unlock mutex
            COUNTER_MUTEX DO UNTIE

            I ITZ I MOAR 1
            SLEEP WIT 1  BTW Simulate other work
        KTHX

        SAY WIT THREAD_ID
        SAYZ WIT " finished"
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN DEMO_SYNCHRONIZED_ACCESS
    SAYZ WIT "=== Synchronized Access Demo ==="

    BTW Reset shared counter
    SHARED_COUNTER ITZ 0

    BTW Create multiple threads
    I HAS A VARIABLE THREAD1 TEH SAFE_COUNTER_THREAD ITZ NEW SAFE_COUNTER_THREAD WIT "Thread-A"
    I HAS A VARIABLE THREAD2 TEH SAFE_COUNTER_THREAD ITZ NEW SAFE_COUNTER_THREAD WIT "Thread-B"
    I HAS A VARIABLE THREAD3 TEH SAFE_COUNTER_THREAD ITZ NEW SAFE_COUNTER_THREAD WIT "Thread-C"

    BTW Start all threads
    THREAD1 DO START
    THREAD2 DO START
    THREAD3 DO START
    SAYZ WIT "All threads started"

    BTW Wait for all to complete
    THREAD1 DO JOIN
    THREAD2 DO JOIN
    THREAD3 DO JOIN
    SAYZ WIT "All threads completed"

    SAY WIT "Final counter value: "
    SAYZ WIT SHARED_COUNTER
    SAYZ WIT "Should be 9 (3 threads × 3 increments each)"
KTHXBAI
```

## Exception Handling in Threading

### Thread Exception Handling

```lol
I CAN HAS THREAD?
I CAN HAS STDIO?

BTW Thread that may throw exceptions
HAI ME TEH CLAS RISKY_THREAD KITTEH OF YARN
    EVRYONE
    DIS TEH VARIABLE SHOULD_FAIL TEH BOOL ITZ NO

    DIS TEH FUNCSHUN RISKY_THREAD WIT FAIL TEH BOOL
        SHOULD_FAIL ITZ FAIL
    KTHX

    DIS TEH FUNCSHUN SPIN
        IZ SHOULD_FAIL?
            OOPS "Thread encountered an error!"
        NOPE
            SAYZ WIT "Thread completed successfully"
        KTHX
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN DEMO_THREAD_EXCEPTIONS
    SAYZ WIT "=== Thread Exception Handling ==="

    BTW Test successful thread
    I HAS A VARIABLE GOOD_THREAD TEH RISKY_THREAD ITZ NEW RISKY_THREAD WIT NO
    GOOD_THREAD DO START

    MAYB
        GOOD_THREAD DO JOIN
        SAYZ WIT "✓ Good thread completed normally"
    OOPSIE THREAD_ERROR
        SAYZ WIT "❌ Unexpected error: "
        SAYZ WIT THREAD_ERROR
    KTHX

    BTW Test failing thread
    I HAS A VARIABLE BAD_THREAD TEH RISKY_THREAD ITZ NEW RISKY_THREAD WIT YEZ
    BAD_THREAD DO START

    MAYB
        BAD_THREAD DO JOIN
        SAYZ WIT "❌ This should not print"
    OOPSIE THREAD_ERROR
        SAYZ WIT "✓ Caught thread exception: "
        SAYZ WIT THREAD_ERROR
    KTHX
KTHXBAI
```

### Mutex Exception Handling

```lol
I CAN HAS THREAD?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN DEMO_MUTEX_EXCEPTIONS
    SAYZ WIT "=== Mutex Exception Handling ==="

    I HAS A VARIABLE MUTEX TEH KNOT ITZ NEW KNOT

    BTW Test double unlock error
    MUTEX DO TIE
    MUTEX DO UNTIE  BTW First unlock - should work

    MAYB
        MUTEX DO UNTIE  BTW Second unlock - should fail
        SAYZ WIT "❌ This should not print"
    OOPSIE MUTEX_ERROR
        SAYZ WIT "✓ Caught expected mutex error: "
        SAYZ WIT MUTEX_ERROR
    KTHX
KTHXBAI
```

## Quick Reference

### YARN Methods

| Method | Description | Throws |
|--------|-------------|---------|
| `START` | Launch thread execution | If already running |
| `JOIN` | Wait for thread completion | Thread exceptions |
| `SPIN` | Thread implementation (abstract) | User-defined |

### YARN Properties

| Property | Type | Description |
|----------|------|-------------|
| `RUNNING` | BOOL | True if thread is running |
| `FINISHED` | BOOL | True if thread has completed |

### KNOT Methods

| Method | Description | Throws |
|--------|-------------|---------|
| `TIE` | Acquire mutex lock | - |
| `UNTIE` | Release mutex lock | If not locked |

### KNOT Properties

| Property | Type | Description |
|----------|------|-------------|
| `LOCKED` | BOOL | True if mutex is locked |

### Threading Patterns

```lol
BTW Basic thread creation
HAI ME TEH CLAS MY_THREAD KITTEH OF YARN
    DIS TEH FUNCSHUN SPIN
        BTW Thread code here
    KTHX
KTHXBAI

BTW Thread lifecycle
I HAS A VARIABLE T TEH MY_THREAD ITZ NEW MY_THREAD
T DO START
T DO JOIN

BTW Mutex synchronization
I HAS A VARIABLE M TEH KNOT ITZ NEW KNOT
M DO TIE
BTW Critical section
M DO UNTIE

BTW Safe resource access pattern
MAYB
    MUTEX DO TIE
    BTW Use shared resource
ALWAYZ
    MUTEX DO UNTIE
KTHX
```

## Related

- [Control Flow](../language-guide/control-flow.md) - Exception handling in threads
- [Classes](../language-guide/classes.md) - Understanding inheritance for thread classes
- [TIME Module](time.md) - Using SLEEP for thread timing