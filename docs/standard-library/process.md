# PROCESS Module

## Import

```lol
BTW Full import
I CAN HAS PROCESS?

BTW Selective import examples
```

## Miscellaneous

### MINION Class

A child process that can be launched, monitored, and controlled with full I/O access.
Provides comprehensive process management including environment variables, working directory, and signal handling.

**Methods:**

#### KILL

Forcefully terminates the running process.
Sends SIGKILL signal to immediately stop the process without cleanup.

**Syntax:** `<minion> DO KILL`
**Example: Kill a running process**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "sleep"
CMD DO PUSH WIT "30"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
SAYZ WIT "Process started, waiting 2 seconds before killing..."
I HAS A VARIABLE START_TIME TEH INTEGR ITZ NOW
WHILE (NOW MINUSZ START_TIME) LIEKZ 2000
BTW Wait 2 seconds
KTHX
PROC DO KILL
I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC DO WAIT
SAYZ WIT "Process killed, exit code: "
SAYZ WIT EXIT_CODE
```

**Example: Kill process on timeout**

```lol
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
I HAS A VARIABLE TIMEOUT TEH INTEGR ITZ 10000
I HAS A VARIABLE START_TIME TEH INTEGR ITZ NOW
WHILE (PROC RUNNING) AN ((NOW MINUSZ START_TIME) LIEKZ TIMEOUT)
I HAS A VARIABLE CHUNK TEH STRIN ITZ PROC STDOUT DO READ WIT 256
IZ CHUNK SAEM AS ""?
OUTTA HERE BTW No more output
KTHX
SAYZ WIT CHUNK
KTHX
IZ PROC RUNNING?
PROC DO KILL
SAYZ WIT "Process timed out and was terminated"
NOPE
SAYZ WIT "Process completed normally"
KTHX
```

**Example: Cleanup in error handling**

```lol
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
MAYB
BTW Do some work with the process
I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 1024
SAYZ WIT OUTPUT
OOPSIE ERR
SAYZ WIT "Error occurred: "
SAYZ WIT ERR
IZ PROC RUNNING?
PROC DO KILL
SAYZ WIT "Process was killed due to error"
KTHX
KTHX
```

**Note:** Immediately terminates process without cleanup

**Note:** Process exit code will be -1 after killing

**Note:** May leave child processes running if not handled properly

**Note:** Use WAIT after KILL to ensure process cleanup

#### MINION

Creates a new process definition with command and arguments.
Initializes with current working directory and copies current environment variables.

**Syntax:** `NEW MINION WIT <cmdline>`
**Parameters:**
- `cmdline` (BUKKIT): Command and arguments as array of strings

**Example: Create process for simple command**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "ls"
CMD DO PUSH WIT "-la"
CMD DO PUSH WIT "/home"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
BTW Process created but not started yet
```

**Example: Create process for complex command with pipes**

```lol
I HAS A VARIABLE GREP_CMD TEH BUKKIT ITZ NEW BUKKIT
GREP_CMD DO PUSH WIT "grep"
GREP_CMD DO PUSH WIT "-n"
GREP_CMD DO PUSH WIT "error"
I HAS A VARIABLE GREP_PROC TEH MINION ITZ NEW MINION WIT GREP_CMD
```

**Example: Create interactive shell**

```lol
I HAS A VARIABLE SHELL_CMD TEH BUKKIT ITZ NEW BUKKIT
SHELL_CMD DO PUSH WIT "bash"
SHELL_CMD DO PUSH WIT "--login"
I HAS A VARIABLE SHELL TEH MINION ITZ NEW MINION WIT SHELL_CMD
```

**Note:** First element in BUKKIT must be the executable name or path

**Note:** All elements must be convertible to strings

**Note:** Process inherits current environment and working directory

**Note:** Use START method to actually launch the process

#### SIGNAL

Sends a signal to the running process (Unix/Linux systems).
Allows sending specific signals like SIGTERM, SIGINT, etc. to processes.

**Syntax:** `<minion> DO SIGNAL WIT <code>`
**Parameters:**
- `code` (INTEGR): Signal number (e.g., 15 for SIGTERM, 2 for SIGINT)

**Example: Send SIGTERM to gracefully stop process**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "sleep"
CMD DO PUSH WIT "30"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
SAYZ WIT "Sending SIGTERM (15) to process..."
PROC DO SIGNAL WIT 15
I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC DO WAIT
SAYZ WIT "Process terminated with signal, exit code: "
SAYZ WIT EXIT_CODE
```

**Example: Send SIGINT (Ctrl+C equivalent)**

```lol
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
SAYZ WIT "Sending SIGINT (2) to process..."
PROC DO SIGNAL WIT 2
I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC DO WAIT
SAYZ WIT "Process interrupted, exit code: "
SAYZ WIT EXIT_CODE
```

**Example: Graceful shutdown with timeout**

```lol
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
SAYZ WIT "Attempting graceful shutdown..."
PROC DO SIGNAL WIT 15
I HAS A VARIABLE START_TIME TEH INTEGR ITZ NOW
WHILE (PROC RUNNING) AN ((NOW MINUSZ START_TIME) LIEKZ 5000)
BTW Wait up to 5 seconds for graceful shutdown
KTHX
IZ PROC RUNNING?
SAYZ WIT "Process didn't respond to SIGTERM, killing..."
PROC DO KILL
NOPE
SAYZ WIT "Process shut down gracefully"
KTHX
```

**Note:** Signal numbers vary by Unix system

**Note:** Common signals: 1=SIGHUP, 2=SIGINT, 9=SIGKILL, 15=SIGTERM

**Note:** Process may ignore some signals

**Note:** Not available on Windows systems

#### START

Launches the child process and creates communication pipes.
Creates STDIN, STDOUT, and STDERR pipes for process communication.

**Syntax:** `<minion> DO START`
**Example: Start simple command**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "date"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
SAYZ WIT "Process started with PID: "
SAYZ WIT PROC PID
```

**Example: Start and immediately read output**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "echo"
CMD DO PUSH WIT "Hello!"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 1024
SAYZ WIT "Process output: "
SAYZ WIT OUTPUT
```

**Example: Start process with custom environment**

```lol
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC ENV DO PUT WIT "MY_VAR" AN WIT "custom_value"
PROC WORKDIR ITZ "/tmp"
PROC DO START
SAYZ WIT "Process started in directory: "
SAYZ WIT PROC WORKDIR
```

**Note:** Process must not be already running or finished

**Note:** After START, pipes become available for I/O operations

**Note:** RUNNING property becomes YEZ after successful start

**Note:** PID property is set to the actual process ID

#### WAIT

Waits for the process to complete and returns the exit code.
Blocks until the child process terminates, then returns its exit status.

**Syntax:** `<minion> DO WAIT`
**Example: Wait for process completion**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "echo"
CMD DO PUSH WIT "Hello, World!"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 1024
I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC DO WAIT
SAYZ WIT "Process completed with exit code: "
SAYZ WIT EXIT_CODE
```

**Example: Handle process errors**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "false"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC DO WAIT
IZ EXIT_CODE SAEM AS 0?
SAYZ WIT "Process succeeded"
NOPE
SAYZ WIT "Process failed with code: "
SAYZ WIT EXIT_CODE
KTHX
```

**Example: Wait with timeout pattern**

```lol
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
I HAS A VARIABLE TIMEOUT TEH INTEGR ITZ 5000
I HAS A VARIABLE START_TIME TEH INTEGR ITZ NOW
WHILE (PROC RUNNING) AN ((NOW LES START_TIME) SAEM AS TIMEOUT)
I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC DO WAIT
OUTTA HERE BTW Process finished
KTHX
IZ PROC RUNNING?
PROC DO KILL
SAYZ WIT "Process timed out and was killed"
KTHX
```

**Note:** Blocks the current thread until process completes

**Note:** Exit code 0 typically indicates success

**Note:** Negative exit codes may indicate process was killed

**Member Variables:**

#### CMDLINE

Read-only property containing the command line arguments.
Contains the executable name and arguments passed to the process.


**Example: Access command line arguments**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "ls"
CMD DO PUSH WIT "-la"
CMD DO PUSH WIT "/home"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
I HAS A VARIABLE ARGS TEH BUKKIT ITZ PROC CMDLINE
SAYZ WIT "Command: "
SAYZ WIT ARGS 0
SAYZ WIT "Arguments: "
IM OUTTA UR ARGS NERFIN ARG
SAYZ WIT ARG
IM IN UR ARGS
```

**Example: Verify command before starting**

```lol
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
I HAS A VARIABLE CMD_ARGS TEH BUKKIT ITZ PROC CMDLINE
IZ (CMD_ARGS LENGTH) BIGGR THAN 0?
SAYZ WIT "Will execute: "
SAYZ WIT CMD_ARGS 0
PROC DO START
NOPE
SAYZ WIT "No command specified"
KTHX
```

**Example: Log process command for debugging**

```lol
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
SAYZ WIT "Starting process with command: "
I HAS A VARIABLE CMD_STR TEH STRIN ITZ ""
I HAS A VARIABLE ARGS TEH BUKKIT ITZ PROC CMDLINE
IM OUTTA UR ARGS NERFIN ARG
IZ CMD_STR SAEM AS ""?
CMD_STR ITZ ARG
NOPE
CMD_STR ITZ CMD_STR MOAR " " MOAR ARG
KTHX
IM IN UR ARGS
SAYZ WIT CMD_STR
```

**Note:** First element is the executable name or path

**Note:** Remaining elements are command arguments

**Note:** Cannot be modified after process creation

**Note:** Arguments are stored as strings

#### ENV

Environment variables for the process.
Key-value pairs that will be set as environment variables for the child process.


**Example: Set custom environment variables**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "printenv"
CMD DO PUSH WIT "MY_VAR"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC ENV DO PUT WIT "MY_VAR" AN WIT "Hello from environment!"
PROC ENV DO PUT WIT "ANOTHER_VAR" AN WIT "Another value"
PROC DO START
I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 1024
SAYZ WIT "Environment variable value: "
SAYZ WIT OUTPUT
```

**Example: Inherit and modify parent environment**

```lol
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
I HAS A VARIABLE ENV TEH BASKIT ITZ PROC ENV
ENV DO PUT WIT "PATH" AN WIT "/custom/path:/usr/bin:/bin"
ENV DO PUT WIT "HOME" AN WIT "/tmp"
PROC DO START
SAYZ WIT "Process started with modified environment"
```

**Example: Clear environment and set minimal vars**

```lol
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
I HAS A VARIABLE NEW_ENV TEH BASKIT ITZ NEW BASKIT
NEW_ENV DO PUT WIT "PATH" AN WIT "/bin:/usr/bin"
NEW_ENV DO PUT WIT "HOME" AN WIT "/tmp"
PROC ENV ITZ NEW_ENV
PROC DO START
SAYZ WIT "Process started with clean environment"
```

**Example: Handle environment variable errors**

```lol
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC ENV DO PUT WIT "INVALID_KEY" AN WIT 123
MAYB
PROC DO START
OOPSIE ERR
SAYZ WIT "Environment setup error: "
SAYZ WIT ERR
PROC ENV DO PUT WIT "INVALID_KEY" AN WIT "123"
PROC DO START
KTHX
```

**Note:** Inherits parent process environment by default

**Note:** Must be set before calling START

**Note:** All values must be convertible to strings

**Note:** Cannot be changed while process is running

**Note:** Use BASKIT operations (PUT, GET, HAS) to modify

#### EXIT_CODE

Read-only property containing the process exit code.
Available after process completion, indicates success (0) or failure (non-zero).


**Example: Check process success**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "ls"
CMD DO PUSH WIT "/nonexistent"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC DO WAIT
IZ EXIT_CODE SAEM AS 0?
SAYZ WIT "Command succeeded"
NOPE
SAYZ WIT "Command failed with code: "
SAYZ WIT EXIT_CODE
KTHX
```

**Example: Access exit code without waiting**

```lol
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
PROC DO WAIT
SAYZ WIT "Process exit code: "
SAYZ WIT PROC EXIT_CODE
```

**Example: Handle killed processes**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "sleep"
CMD DO PUSH WIT "30"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
PROC DO KILL
I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC DO WAIT
SAYZ WIT "Killed process exit code: "
SAYZ WIT EXIT_CODE
```

**Example: Exit code before process finishes**

```lol
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
SAYZ WIT "Exit code before completion: "
SAYZ WIT PROC EXIT_CODE
PROC DO WAIT
SAYZ WIT "Exit code after completion: "
SAYZ WIT PROC EXIT_CODE
```

**Example: Categorize exit codes**

```lol
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
PROC DO WAIT
I HAS A VARIABLE CODE TEH INTEGR ITZ PROC EXIT_CODE
IZ CODE SAEM AS 0?
SAYZ WIT "Success"
NOPE
IZ CODE SAEM AS -1?
SAYZ WIT "Process was killed"
NOPE
IZ CODE BIGGR THAN 128?
SAYZ WIT "Process terminated by signal: "
SAYZ WIT CODE MINUSZ 128
NOPE
SAYZ WIT "Process failed with code: "
SAYZ WIT CODE
KTHX
KTHX
KTHX
```

**Note:** 0 typically indicates successful execution

**Note:** Non-zero values indicate various types of failures

**Note:** -1 indicates process was killed

**Note:** Values > 128 often indicate termination by signal

**Note:** Only meaningful after process has finished

#### FINISHED

Read-only property indicating whether the process has completed.
True after process termination, remains true until process is restarted.


**Example: Check process completion status**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "echo"
CMD DO PUSH WIT "Hello"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
SAYZ WIT "Before start - Finished: "
SAYZ WIT PROC FINISHED
PROC DO START
SAYZ WIT "After start - Finished: "
SAYZ WIT PROC FINISHED
I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC DO WAIT
SAYZ WIT "After wait - Finished: "
SAYZ WIT PROC FINISHED
```

**Example: Poll for completion**

```lol
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
WHILE NO SAEM AS (PROC FINISHED)
I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 256
IZ NO SAEM AS (OUTPUT SAEM AS "")?
SAYZ WIT OUTPUT
KTHX
I HAS A VARIABLE SLEEP_TIME TEH INTEGR ITZ 100
BTW Sleep implementation would go here
KTHX
SAYZ WIT "Process has completed"
```

**Example: Handle already finished process**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "true"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
I HAS A VARIABLE SHORT_WAIT TEH INTEGR ITZ 50
BTW Sleep implementation would go here
IZ PROC FINISHED?
SAYZ WIT "Process finished quickly"
SAYZ WIT "Exit code: "
SAYZ WIT PROC EXIT_CODE
NOPE
SAYZ WIT "Process still running"
PROC DO WAIT
KTHX
```

**Example: Multiple WAIT calls on finished process**

```lol
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
I HAS A VARIABLE FIRST_WAIT TEH INTEGR ITZ PROC DO WAIT
SAYZ WIT "First wait result: "
SAYZ WIT FIRST_WAIT
I HAS A VARIABLE SECOND_WAIT TEH INTEGR ITZ PROC DO WAIT
SAYZ WIT "Second wait result: "
SAYZ WIT SECOND_WAIT
SAYZ WIT "Finished status: "
SAYZ WIT PROC FINISHED
```

**Note:** Becomes YEZ after process terminates

**Note:** Remains YEZ even after multiple WAIT calls

**Note:** Useful for checking if WAIT will block or return immediately

**Note:** Different from RUNNING - FINISHED indicates completion, RUNNING indicates current state

#### PID

Read-only property containing the process ID.
System-assigned unique identifier for the running process.


**Example: Get process ID after starting**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "sleep"
CMD DO PUSH WIT "10"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
SAYZ WIT "Process started with PID: "
SAYZ WIT PROC PID
```

**Example: Monitor process by PID**

```lol
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
I HAS A VARIABLE PROCESS_ID TEH INTEGR ITZ PROC PID
SAYZ WIT "Monitoring process "
SAYZ WIT PROCESS_ID
WHILE PROC RUNNING
SAYZ WIT "Process "
SAYZ WIT PROCESS_ID
SAYZ WIT " is still running"
I HAS A VARIABLE SLEEP_TIME TEH INTEGR ITZ 1000
BTW Sleep implementation would go here
KTHX
```

**Example: PID before and after start**

```lol
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
SAYZ WIT "PID before start: "
SAYZ WIT PROC PID
PROC DO START
SAYZ WIT "PID after start: "
SAYZ WIT PROC PID
```

**Example: Use PID for external monitoring**

```lol
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
I HAS A VARIABLE PID TEH INTEGR ITZ PROC PID
I HAS A VARIABLE MONITOR_CMD TEH BUKKIT ITZ NEW BUKKIT
MONITOR_CMD DO PUSH WIT "ps"
MONITOR_CMD DO PUSH WIT "-p"
MONITOR_CMD DO PUSH WIT PID
I HAS A VARIABLE MONITOR TEH MINION ITZ NEW MINION WIT MONITOR_CMD
MONITOR DO START
I HAS A VARIABLE PS_OUTPUT TEH STRIN ITZ MONITOR STDOUT DO READ WIT 1024
SAYZ WIT "Process status:"
SAYZ WIT PS_OUTPUT
MONITOR DO WAIT
```

**Example: Handle PID for process groups**

```lol
I HAS A VARIABLE PROC1 TEH MINION ITZ NEW MINION WIT CMD
I HAS A VARIABLE PROC2 TEH MINION ITZ NEW MINION WIT CMD
PROC1 DO START
PROC2 DO START
SAYZ WIT "Started processes with PIDs: "
SAYZ WIT PROC1 PID
SAYZ WIT " and "
SAYZ WIT PROC2 PID
```

**Note:** Assigned by operating system when process starts

**Note:** Unique identifier for the process

**Note:** -1 indicates process has not been started

**Note:** Can be used with system tools like ps, kill, etc.

**Note:** Remains valid until process terminates

#### RUNNING

Read-only property indicating whether the process is currently running.
True from START until process termination, false otherwise.


**Example: Check if process is still running**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "sleep"
CMD DO PUSH WIT "5"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
SAYZ WIT "Process running: "
SAYZ WIT PROC RUNNING
I HAS A VARIABLE START_TIME TEH INTEGR ITZ NOW
WHILE (PROC RUNNING) AN ((NOW MINUSZ START_TIME) LIEKZ 3000)
BTW Wait up to 3 seconds
KTHX
SAYZ WIT "Process still running: "
SAYZ WIT PROC RUNNING
```

**Example: Wait for process to finish**

```lol
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
WHILE PROC RUNNING
I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 256
IZ NO SAEM AS (OUTPUT SAEM AS "")?
SAYZ WIT OUTPUT
KTHX
I HAS A VARIABLE SLEEP_TIME TEH INTEGR ITZ 100
BTW Sleep implementation would go here
KTHX
SAYZ WIT "Process has finished"
```

**Example: Monitor process lifecycle**

```lol
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
SAYZ WIT "Before start - Running: "
SAYZ WIT PROC RUNNING
PROC DO START
SAYZ WIT "After start - Running: "
SAYZ WIT PROC RUNNING
I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC DO WAIT
SAYZ WIT "After wait - Running: "
SAYZ WIT PROC RUNNING
```

**Example: Handle process that exits quickly**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "true"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
IZ PROC RUNNING?
SAYZ WIT "Process is running"
I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC DO WAIT
NOPE
SAYZ WIT "Process already finished"
SAYZ WIT "Exit code: "
SAYZ WIT PROC EXIT_CODE
KTHX
```

**Note:** Becomes YEZ immediately after successful START

**Note:** Becomes NO when process terminates (normally or killed)

**Note:** Check after WAIT to confirm process has stopped

**Note:** Useful for polling or timeout implementations

#### STDERR

Read-only property providing access to the process's standard error pipe.
Available after START, allows reading error messages from the child process.


**Example: Read process error output**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "ls"
CMD DO PUSH WIT "/nonexistent/directory"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
I HAS A VARIABLE ERRORS TEH STRIN ITZ PROC STDERR DO READ WIT 1024
IZ NO SAEM AS (ERRORS SAEM AS "")?
SAYZ WIT "Process errors: "
SAYZ WIT ERRORS
NOPE
SAYZ WIT "No errors reported"
KTHX
PROC DO WAIT
```

**Example: Separate stdout and stderr handling**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "python3"
CMD DO PUSH WIT "-c"
CMD DO PUSH WIT "import sys; print('output'); print('error', file=sys.stderr)"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 512
I HAS A VARIABLE ERRORS TEH STRIN ITZ PROC STDERR DO READ WIT 512
SAYZ WIT "STDOUT: "
SAYZ WIT OUTPUT
SAYZ WIT "STDERR: "
SAYZ WIT ERRORS
PROC DO WAIT
```

**Example: Monitor stderr for warnings**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "gcc"
CMD DO PUSH WIT "-Wall"
CMD DO PUSH WIT "program.c"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
I HAS A VARIABLE WARNINGS TEH STRIN ITZ ""
WHILE NO SAEM AS (PROC STDERR IS_EOF)
I HAS A VARIABLE CHUNK TEH STRIN ITZ PROC STDERR DO READ WIT 256
WARNINGS ITZ WARNINGS MOAR CHUNK
KTHX
PROC DO WAIT
IZ NO SAEM AS (WARNINGS SAEM AS "")?
SAYZ WIT "Compiler warnings:"
SAYZ WIT WARNINGS
KTHX
```

**Example: Check for errors after process completion**

```lol
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
PROC DO WAIT
I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC EXIT_CODE
IZ EXIT_CODE SAEM AS 0?
I HAS A VARIABLE ERRORS TEH STRIN ITZ PROC STDERR DO READ WIT 1024
IZ NO SAEM AS (ERRORS SAEM AS "")?
SAYZ WIT "Process succeeded but produced warnings: "
SAYZ WIT ERRORS
KTHX
NOPE
I HAS A VARIABLE ERRORS TEH STRIN ITZ PROC STDERR DO READ WIT 1024
SAYZ WIT "Process failed with errors: "
SAYZ WIT ERRORS
KTHX
```

**Example: Real-time error monitoring**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "long-running-command"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
WHILE PROC RUNNING
I HAS A VARIABLE ERROR_CHUNK TEH STRIN ITZ PROC STDERR DO READ WIT 256
IZ NO SAEM AS (ERROR_CHUNK SAEM AS "")?
SAYZ WIT "[ERROR] "
SAYZ WIT ERROR_CHUNK
KTHX
I HAS A VARIABLE SLEEP_TIME TEH INTEGR ITZ 500
BTW Sleep implementation would go here
KTHX
PROC DO WAIT
```

**Note:** Only available after calling START

**Note:** Supports READ operations to get error messages

**Note:** Many programs write diagnostic messages to stderr

**Note:** Should be checked even when exit code is 0

**Note:** Throws exception if accessed before START

#### STDIN

Read-only property providing access to the process's standard input pipe.
Available after START, allows writing data to the child process.


**Example: Write to process stdin**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "cat"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
PROC STDIN DO WRITE WIT "Hello, World!\n"
PROC STDIN DO WRITE WIT "This is input data\n"
PROC STDIN DO CLOSE
I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 1024
SAYZ WIT OUTPUT
```

**Example: Send commands to interactive shell**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "python3"
CMD DO PUSH WIT "-i"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
PROC STDIN DO WRITE WIT "print('Hello from Python')\n"
PROC STDIN DO WRITE WIT "x = 42\n"
PROC STDIN DO WRITE WIT "print(f'x = {x}')\n"
PROC STDIN DO WRITE WIT "exit()\n"
I HAS A VARIABLE RESULT TEH STRIN ITZ PROC STDOUT DO READ WIT 2048
SAYZ WIT RESULT
```

**Example: Pipe data between processes**

```lol
I HAS A VARIABLE CMD1 TEH BUKKIT ITZ NEW BUKKIT
CMD1 DO PUSH WIT "echo"
CMD1 DO PUSH WIT "line 1\nline 2\nline 3"
I HAS A VARIABLE PROC1 TEH MINION ITZ NEW MINION WIT CMD1
I HAS A VARIABLE CMD2 TEH BUKKIT ITZ NEW BUKKIT
CMD2 DO PUSH WIT "grep"
CMD2 DO PUSH WIT "line"
I HAS A VARIABLE PROC2 TEH MINION ITZ NEW MINION WIT CMD2
PROC2 DO START
PROC1 DO START
I HAS A VARIABLE DATA TEH STRIN ITZ PROC1 STDOUT DO READ WIT 1024
PROC2 STDIN DO WRITE WIT DATA
PROC2 STDIN DO CLOSE
I HAS A VARIABLE FILTERED TEH STRIN ITZ PROC2 STDOUT DO READ WIT 1024
SAYZ WIT FILTERED
```

**Example: Handle stdin pipe errors**

```lol
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
MAYB
PROC STDIN DO WRITE WIT "some data\n"
PROC STDIN DO CLOSE
OOPSIE ERR
SAYZ WIT "Error writing to stdin: "
SAYZ WIT ERR
IZ PROC STDIN IS_OPEN?
PROC STDIN DO CLOSE
KTHX
KTHX
```

**Note:** Only available after calling START

**Note:** Supports WRITE operations to send data to process

**Note:** Should be closed when done writing to signal EOF

**Note:** Throws exception if accessed before START

**Note:** Pipe becomes unavailable after process termination

#### STDOUT

Read-only property providing access to the process's standard output pipe.
Available after START, allows reading data output by the child process.


**Example: Read process output**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "echo"
CMD DO PUSH WIT "Hello, World!"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 1024
SAYZ WIT "Process output: "
SAYZ WIT OUTPUT
```

**Example: Read output in chunks**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "ls"
CMD DO PUSH WIT "-la"
CMD DO PUSH WIT "/usr/bin"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
I HAS A VARIABLE ALL_OUTPUT TEH STRIN ITZ ""
WHILE NO SAEM AS (PROC STDOUT IS_EOF)
I HAS A VARIABLE CHUNK TEH STRIN ITZ PROC STDOUT DO READ WIT 512
IZ CHUNK SAEM AS ""?
OUTTA HERE BTW No more data
KTHX
ALL_OUTPUT ITZ ALL_OUTPUT MOAR CHUNK
KTHX
SAYZ WIT "Directory listing:"
SAYZ WIT ALL_OUTPUT
```

**Example: Handle both stdout and stderr**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "find"
CMD DO PUSH WIT "/nonexistent"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
I HAS A VARIABLE OUTPUT TEH STRIN ITZ ""
I HAS A VARIABLE ERRORS TEH STRIN ITZ ""
WHILE NO SAEM AS ((PROC STDOUT IS_EOF) AN (PROC STDERR IS_EOF))
IZ NO SAEM AS (PROC STDOUT IS_EOF)?
I HAS A VARIABLE CHUNK TEH STRIN ITZ PROC STDOUT DO READ WIT 256
OUTPUT ITZ OUTPUT MOAR CHUNK
KTHX
IZ NO SAEM AS (PROC STDERR IS_EOF)?
I HAS A VARIABLE ERR_CHUNK TEH STRIN ITZ PROC STDERR DO READ WIT 256
ERRORS ITZ ERRORS MOAR ERR_CHUNK
KTHX
KTHX
SAYZ WIT "Output: "
SAYZ WIT OUTPUT
SAYZ WIT "Errors: "
SAYZ WIT ERRORS
```

**Example: Process streaming output**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "ping"
CMD DO PUSH WIT "-c"
CMD DO PUSH WIT "3"
CMD DO PUSH WIT "8.8.8.8"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
WHILE NO SAEM AS (PROC STDOUT IS_EOF)
I HAS A VARIABLE LINE TEH STRIN ITZ PROC STDOUT DO READ WIT 256
IZ NO SAEM AS (LINE SAEM AS "")?
SAYZ WIT "Received: "
SAYZ WIT LINE
KTHX
KTHX
PROC DO WAIT
```

**Note:** Only available after calling START

**Note:** Supports READ operations to get process output

**Note:** Reading returns empty string when no data available

**Note:** IS_EOF becomes true when process closes stdout

**Note:** Throws exception if accessed before START

#### WORKDIR

Working directory for the process.
Directory where the process will be executed. Can be changed before starting.


**Example: Set working directory before starting**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "ls"
CMD DO PUSH WIT "-la"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC WORKDIR ITZ "/tmp"
PROC DO START
I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 1024
SAYZ WIT "Contents of /tmp:"
SAYZ WIT OUTPUT
```

**Example: Use relative path**

```lol
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC WORKDIR ITZ "../parent_directory"
PROC DO START
SAYZ WIT "Process started in relative directory"
```

**Example: Change working directory dynamically**

```lol
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC WORKDIR ITZ "/home/user"
SAYZ WIT "Initial workdir: "
SAYZ WIT PROC WORKDIR
PROC WORKDIR ITZ "/var/log"
SAYZ WIT "Changed workdir: "
SAYZ WIT PROC WORKDIR
```

**Example: Handle working directory errors**

```lol
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC WORKDIR ITZ "/nonexistent/directory"
MAYB
PROC DO START
SAYZ WIT "Process started successfully"
OOPSIE ERR
SAYZ WIT "Failed to start process: "
SAYZ WIT ERR
PROC WORKDIR ITZ "/tmp"
PROC DO START
KTHX
```

**Note:** Must be set before calling START

**Note:** Can be absolute or relative path

**Note:** Directory must exist and be accessible

**Note:** Cannot be changed while process is running

**Note:** Defaults to current working directory if not set

**Example: Basic command execution**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "echo"
CMD DO PUSH WIT "Hello, World!"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 1024
I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC DO WAIT
SAYZ WIT "Output: "
SAYZ WIT OUTPUT
SAYZ WIT "Exit code: "
SAYZ WIT EXIT_CODE
```

**Example: Interactive process with stdin/stdout**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "python3"
CMD DO PUSH WIT "-i"
I HAS A VARIABLE PYTHON TEH MINION ITZ NEW MINION WIT CMD
PYTHON DO START
PYTHON STDIN DO WRITE WIT "print(2 + 2)\n"
PYTHON STDIN DO WRITE WIT "exit()\n"
I HAS A VARIABLE RESULT TEH STRIN ITZ PYTHON STDOUT DO READ WIT 1024
PYTHON DO WAIT
SAYZ WIT "Python result: "
SAYZ WIT RESULT
```

**Example: Process with custom environment and working directory**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "printenv"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC WORKDIR ITZ "/tmp"
PROC ENV DO PUT WIT "CUSTOM_VAR" AN WIT "Hello from environment!"
PROC DO START
I HAS A VARIABLE ENV_OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 2048
PROC DO WAIT
SAYZ WIT ENV_OUTPUT
```

### PIPE Class

A communication pipe connected to a child process's standard input, output, or error streams.
Provides read/write access to process streams and inherits from IO.READER and IO.WRITER interfaces.

**Methods:**

#### CLOSE

Closes the pipe connection and releases associated resources.
For stdin pipes, signals end-of-input to the child process.

**Syntax:** `<pipe> DO CLOSE`
**Example: Close stdin to signal completion**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "wc"
CMD DO PUSH WIT "-l"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
PROC STDIN DO WRITE WIT "line 1\n"
PROC STDIN DO WRITE WIT "line 2\n"
PROC STDIN DO WRITE WIT "line 3\n"
PROC STDIN DO CLOSE BTW Signal end of input
I HAS A VARIABLE COUNT TEH STRIN ITZ PROC STDOUT DO READ WIT 100
SAYZ WIT "Line count: "
SAYZ WIT COUNT
```

**Example: Cleanup pipes after process**

```lol
PROC DO WAIT
PROC STDOUT DO CLOSE
PROC STDERR DO CLOSE
SAYZ WIT "All pipes closed"
```

**Example: Close in error handling**

```lol
MAYB
PROC STDIN DO WRITE WIT "some input"
BTW More operations here
OOPSIE ERR
SAYZ WIT "Error occurred: "
SAYZ WIT ERR
ALWAYZ
IZ PROC STDIN IS_OPEN?
PROC STDIN DO CLOSE
KTHX
KTHX
```

**Note:** Safe to call multiple times - no error if already closed

**Note:** For stdin pipes: signals EOF to child process

**Note:** For stdout/stderr pipes: prevents further reading

**Note:** Pipes are automatically closed when process terminates

#### READ

Reads up to the specified number of characters from stdout or stderr pipes.
Blocks until data is available or the pipe reaches end-of-file.

**Syntax:** `<pipe> DO READ WIT <size>`
**Parameters:**
- `size` (INTEGR): Maximum number of characters to read

**Example: Read process output**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "ls"
CMD DO PUSH WIT "-la"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 4096
SAYZ WIT "Directory listing:"
SAYZ WIT OUTPUT
```

**Example: Read in chunks**

```lol
I HAS A VARIABLE BUFFER TEH STRIN ITZ ""
WHILE NO SAEM AS (PROC STDOUT IS_EOF)
I HAS A VARIABLE CHUNK TEH STRIN ITZ PROC STDOUT DO READ WIT 256
IZ CHUNK SAEM AS ""?
OUTTA HERE BTW End of output
KTHX
BUFFER ITZ BUFFER MOAR CHUNK
KTHX
```

**Example: Read stderr separately**

```lol
I HAS A VARIABLE ERROR_DATA TEH STRIN ITZ PROC STDERR DO READ WIT 1024
IZ NO SAEM AS (ERROR_DATA SAEM AS "")?
SAYZ WIT "Process errors: "
SAYZ WIT ERROR_DATA
KTHX
```

**Note:** Only works on STDOUT and STDERR pipes, not STDIN

**Note:** Returns empty string when end-of-file is reached

**Note:** May return fewer characters than requested if less data is available

#### WRITE

Writes string data to the process's stdin pipe.
Returns the number of bytes written, only works on STDIN pipes.

**Syntax:** `<pipe> DO WRITE WIT <data>`
**Parameters:**
- `data` (STRIN): The string data to send to the process

**Example: Send input to interactive process**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "python3"
CMD DO PUSH WIT "-i"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
PROC STDIN DO WRITE WIT "print('Hello from Python!')\n"
PROC STDIN DO WRITE WIT "exit()\n"
I HAS A VARIABLE RESULT TEH STRIN ITZ PROC STDOUT DO READ WIT 1024
SAYZ WIT RESULT
```

**Example: Send data to filter process**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "grep"
CMD DO PUSH WIT "error"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
PROC STDIN DO WRITE WIT "info: starting process\n"
PROC STDIN DO WRITE WIT "error: something went wrong\n"
PROC STDIN DO WRITE WIT "info: process completed\n"
PROC STDIN DO CLOSE BTW Signal end of input
I HAS A VARIABLE FILTERED TEH STRIN ITZ PROC STDOUT DO READ WIT 1024
SAYZ WIT "Filtered output: "
SAYZ WIT FILTERED
```

**Example: Pipe data between processes**

```lol
PROC STDIN DO WRITE WIT "line 1\n"
PROC STDIN DO WRITE WIT "line 2\n"
I HAS A VARIABLE BYTES_WRITTEN TEH INTEGR ITZ PROC STDIN DO WRITE WIT "line 3\n"
SAYZ WIT "Wrote "
SAYZ WIT BYTES_WRITTEN
SAYZ WIT " bytes"
```

**Note:** Only works on STDIN pipes, throws exception on STDOUT/STDERR pipes

**Note:** Process must be started before writing to pipes

**Note:** Close stdin pipe when done writing to signal end of input

**Member Variables:**

#### FD_TYPE

Read-only property indicating the type of pipe connection.


**Example: Check pipe type for appropriate operations**

```lol
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
IZ (PROC STDIN FD_TYPE) SAEM AS "STDIN"?
PROC STDIN DO WRITE WIT "input data"
PROC STDIN DO CLOSE
KTHX
```

**Example: Handle different pipe types**

```lol
I HAS A VARIABLE PIPES TEH BUKKIT ITZ NEW BUKKIT
PIPES DO PUSH WIT PROC STDIN
PIPES DO PUSH WIT PROC STDOUT
PIPES DO PUSH WIT PROC STDERR
IM OUTTA UR PIPES NERFIN PIPE
I HAS A VARIABLE TYPE TEH STRIN ITZ PIPE FD_TYPE
IZ TYPE SAEM AS "STDIN"?
SAYZ WIT "Input pipe found"
NOPE
SAYZ WIT "Output pipe: "
SAYZ WIT TYPE
KTHX
IM IN UR PIPES
```

**Note:** Set automatically when pipe is created

**Note:** STDIN pipes support write operations

**Note:** STDOUT and STDERR pipes support read operations

#### IS_EOF

Read-only property indicating whether end-of-file has been reached on read pipes.


**Example: Read until EOF**

```lol
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
I HAS A VARIABLE ALL_OUTPUT TEH STRIN ITZ ""
WHILE NO SAEM AS (PROC STDOUT IS_EOF)
I HAS A VARIABLE CHUNK TEH STRIN ITZ PROC STDOUT DO READ WIT 512
IZ CHUNK SAEM AS ""?
OUTTA HERE BTW No more data available
KTHX
ALL_OUTPUT ITZ ALL_OUTPUT MOAR CHUNK
KTHX
SAYZ WIT "Complete output: "
SAYZ WIT ALL_OUTPUT
```

**Example: Check EOF status after reading**

```lol
I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 1024
IZ PROC STDOUT IS_EOF?
SAYZ WIT "Reached end of process output"
NOPE
SAYZ WIT "More output may be available"
KTHX
```

**Example: Handle both stdout and stderr EOF**

```lol
WHILE NO SAEM AS ((PROC STDOUT IS_EOF) AN (PROC STDERR IS_EOF))
IZ NO SAEM AS (PROC STDOUT IS_EOF)?
I HAS A VARIABLE OUT TEH STRIN ITZ PROC STDOUT DO READ WIT 256
IZ NO SAEM AS (OUT SAEM AS "")?
SAYZ WIT "STDOUT: "
SAYZ WIT OUT
KTHX
KTHX
IZ NO SAEM AS (PROC STDERR IS_EOF)?
I HAS A VARIABLE ERR TEH STRIN ITZ PROC STDERR DO READ WIT 256
IZ NO SAEM AS (ERR SAEM AS "")?
SAYZ WIT "STDERR: "
SAYZ WIT ERR
KTHX
KTHX
KTHX
```

**Note:** Only relevant for STDOUT and STDERR pipes, not STDIN

**Note:** Automatically set when READ operation encounters EOF

**Note:** Once EOF is reached, further READ calls return empty string

#### IS_OPEN

Read-only property indicating whether the pipe is open for I/O operations.


**Example: Check if pipe is available before operations**

```lol
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
IZ PROC STDIN IS_OPEN?
PROC STDIN DO WRITE WIT "data"
PROC STDIN DO CLOSE
NOPE
SAYZ WIT "Stdin pipe is not available"
KTHX
```

**Example: Monitor pipe status during operations**

```lol
WHILE (PROC STDOUT IS_OPEN)
I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 256
IZ OUTPUT SAEM AS ""?
OUTTA HERE BTW End of output reached
KTHX
SAYZ WIT OUTPUT
KTHX
```

**Example: Safe pipe cleanup**

```lol
IZ PROC STDERR IS_OPEN?
I HAS A VARIABLE ERRORS TEH STRIN ITZ PROC STDERR DO READ WIT 1024
PROC STDERR DO CLOSE
IZ NO SAEM AS (ERRORS SAEM AS "")?
SAYZ WIT "Process errors: "
SAYZ WIT ERRORS
KTHX
KTHX
```

**Note:** Automatically set to YEZ when process starts

**Note:** Becomes NO when pipe is closed or process terminates

**Note:** Use to avoid exceptions from operations on closed pipes

**Example: Reading from process stdout**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "echo"
CMD DO PUSH WIT "Hello, World!"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
I HAS A VARIABLE OUTPUT TEH STRIN ITZ PROC STDOUT DO READ WIT 1024
SAYZ WIT OUTPUT
PROC DO WAIT
```

**Example: Writing to process stdin**

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "cat"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
PROC STDIN DO WRITE WIT "Hello from parent process!"
PROC STDIN DO CLOSE
I HAS A VARIABLE RESULT TEH STRIN ITZ PROC STDOUT DO READ WIT 1024
SAYZ WIT RESULT
```

**Example: Handle process stderr**

```lol
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
PROC DO START
I HAS A VARIABLE ERROR_OUTPUT TEH STRIN ITZ PROC STDERR DO READ WIT 512
IZ NO SAEM AS (ERROR_OUTPUT SAEM AS "")?
SAYZ WIT "Process error: "
SAYZ WIT ERROR_OUTPUT
KTHX
```

### READER Class

Abstract base class for objects that can read data.
Defines the interface for reading operations with READ and CLOSE methods.

**Methods:**

#### CLOSE

Closes the reader and releases any associated resources.
Should be called when done reading to ensure proper cleanup.

**Syntax:** `<reader> DO CLOSE`
**Example: Basic cleanup**

```lol
I HAS A VARIABLE READER TEH READER ITZ GET_FILE_READER
I HAS A VARIABLE DATA TEH STRIN ITZ READER DO READ WIT 1024
WHILE NO SAEM AS (DATA LENGTH SAEM AS 0)
PROCESS_DATA WIT DATA
DATA ITZ READER DO READ WIT 1024
KTHX
READER DO CLOSE
SAYZ WIT "Reader closed successfully"
```

**Example: Close in error handling**

```lol
I HAS A VARIABLE READER TEH READER ITZ GET_NETWORK_READER
MAYB
I HAS A VARIABLE DATA TEH STRIN ITZ READER DO READ WIT 512
PROCESS_DATA WIT DATA
OOPSIE ERR
SAYZ WIT "Error during processing: "
SAYZ WIT ERR
READER DO CLOSE
KTHX
```

**Example: Multiple readers cleanup**

```lol
I HAS A VARIABLE READERS TEH BUKKIT ITZ NEW BUKKIT
READERS DO PUSH WIT GET_FILE_READER_1
READERS DO PUSH WIT GET_FILE_READER_2
READERS DO PUSH WIT GET_FILE_READER_3
WHILE NO SAEM AS (READERS LENGTH SAEM AS 0)
I HAS A VARIABLE READER TEH READER ITZ READERS DO POP
I HAS A VARIABLE DATA TEH STRIN ITZ READER DO READ WIT 1024
WHILE NO SAEM AS (DATA LENGTH SAEM AS 0)
PROCESS_DATA WIT DATA
DATA ITZ READER DO READ WIT 1024
KTHX
READER DO CLOSE
KTHX
SAYZ WIT "All readers closed"
```

**Example: Close with resource tracking**

```lol
I HAS A VARIABLE READER TEH READER ITZ GET_DATABASE_READER
I HAS A VARIABLE RESOURCE_COUNT TEH INTEGR ITZ 1
MAYB
I HAS A VARIABLE DATA TEH STRIN ITZ READER DO READ WIT 2048
WHILE NO SAEM AS (DATA LENGTH SAEM AS 0)
SAVE_TO_DATABASE WIT DATA
DATA ITZ READER DO READ WIT 2048
KTHX
OOPSIE ERR
SAYZ WIT "Database operation failed: "
SAYZ WIT ERR
KTHX
RESOURCE_COUNT ITZ RESOURCE_COUNT MINUS 1
READER DO CLOSE
SAYZ WIT "Resources remaining: "
SAYZ WIT RESOURCE_COUNT
```

**Note:** Always call CLOSE to prevent resource leaks

**Note:** CLOSE should be called even if errors occur during reading

**Note:** Multiple CLOSE calls should be safe (idempotent)

**Note:** CLOSE may flush buffers or finalize operations

**Note:** After CLOSE, further READ operations may fail

#### READ

Reads up to the specified number of characters from the input source.
Returns empty string when end-of-file is reached.

**Syntax:** `<reader> DO READ WIT <size>`
**Parameters:**
- `size` (INTEGR): Maximum number of characters to read

**Example: Read fixed-size chunks**

```lol
I HAS A VARIABLE READER TEH READER ITZ GET_FILE_READER
I HAS A VARIABLE CHUNK_SIZE TEH INTEGR ITZ 1024
I HAS A VARIABLE DATA TEH STRIN ITZ READER DO READ WIT CHUNK_SIZE
WHILE NO SAEM AS (DATA LENGTH SAEM AS 0)
SAYZ WIT "Read chunk of length: "
SAYZ WIT DATA LENGTH
DATA ITZ READER DO READ WIT CHUNK_SIZE
KTHX
READER DO CLOSE
```

**Example: Read single character at a time**

```lol
I HAS A VARIABLE READER TEH READER ITZ GET_INPUT_READER
I HAS A VARIABLE CHAR TEH STRIN ITZ READER DO READ WIT 1
WHILE NO SAEM AS (CHAR LENGTH SAEM AS 0)
SAYZ WIT "Character: "
SAYZ WIT CHAR
CHAR ITZ READER DO READ WIT 1
KTHX
READER DO CLOSE
```

**Example: Handle end-of-file**

```lol
I HAS A VARIABLE READER TEH READER ITZ GET_FILE_READER
I HAS A VARIABLE DATA TEH STRIN ITZ READER DO READ WIT 100
IZ DATA LENGTH SAEM AS 0?
SAYZ WIT "Reached end of file"
NOPE
SAYZ WIT "Read data: "
SAYZ WIT DATA
KTHX
READER DO CLOSE
```

**Example: Read with size validation**

```lol
I HAS A VARIABLE READER TEH READER ITZ GET_NETWORK_READER
I HAS A VARIABLE REQUESTED_SIZE TEH INTEGR ITZ 2048
I HAS A VARIABLE DATA TEH STRIN ITZ READER DO READ WIT REQUESTED_SIZE
IZ DATA LENGTH SAEM AS 0?
SAYZ WIT "No data available"
NOPE
IZ DATA LENGTH BIGGR THAN REQUESTED_SIZE?
SAYZ WIT "Error: Read more than requested"
NOPE
SAYZ WIT "Successfully read "
SAYZ WIT DATA LENGTH
SAYZ WIT " characters"
KTHX
KTHX
READER DO CLOSE
```

**Note:** May return fewer characters than requested

**Note:** Returns empty string when end-of-file is reached

**Note:** Size must be a positive integer

**Note:** Implementation depends on the concrete reader type

**Example: Basic reader usage pattern**

```lol
I HAS A VARIABLE READER TEH READER ITZ GET_SOME_READER
I HAS A VARIABLE DATA TEH STRIN ITZ READER DO READ WIT 1024
WHILE NO SAEM AS (DATA LENGTH SAEM AS 0)
SAYZ WIT DATA
DATA ITZ READER DO READ WIT 1024
KTHX
READER DO CLOSE
```

**Example: Reading with error handling**

```lol
I HAS A VARIABLE READER TEH READER ITZ GET_SOME_READER
MAYB
I HAS A VARIABLE DATA TEH STRIN ITZ READER DO READ WIT 512
SAYZ WIT "Read data: "
SAYZ WIT DATA
OOPSIE ERR
SAYZ WIT "Error reading: "
SAYZ WIT ERR
KTHX
READER DO CLOSE
```

**Example: Reading entire content**

```lol
I HAS A VARIABLE READER TEH READER ITZ GET_SOME_READER
I HAS A VARIABLE CONTENT TEH STRIN ITZ ""
I HAS A VARIABLE CHUNK TEH STRIN ITZ READER DO READ WIT 4096
WHILE NO SAEM AS (CHUNK LENGTH SAEM AS 0)
CONTENT ITZ CONTENT MOAR CHUNK
CHUNK ITZ READER DO READ WIT 4096
KTHX
SAYZ WIT "Total content length: "
SAYZ WIT CONTENT LENGTH
READER DO CLOSE
```

### WRITER Class

Abstract base class for objects that can write data.
Defines the interface for writing operations with WRITE and CLOSE methods.

**Methods:**

#### CLOSE

Closes the writer and releases any associated resources.
Should be called when done writing to ensure proper cleanup.

**Syntax:** `<writer> DO CLOSE`
**Example: Basic cleanup after writing**

```lol
I HAS A VARIABLE WRITER TEH WRITER ITZ GET_FILE_WRITER
WRITER DO WRITE WIT "Final data"
WRITER DO CLOSE
SAYZ WIT "Writer closed successfully"
```

**Example: Close with error handling**

```lol
I HAS A VARIABLE WRITER TEH WRITER ITZ GET_NETWORK_WRITER
MAYB
WRITER DO WRITE WIT "Important data"
WRITER DO CLOSE
SAYZ WIT "Data written and writer closed"
OOPSIE ERR
SAYZ WIT "Error during write/close: "
SAYZ WIT ERR
MAYB
WRITER DO CLOSE
OOPSIE CLOSE_ERR
SAYZ WIT "Error closing writer: "
SAYZ WIT CLOSE_ERR
KTHX
KTHX
```

**Example: Multiple writers cleanup**

```lol
I HAS A VARIABLE WRITERS TEH BUKKIT ITZ NEW BUKKIT
WRITERS DO PUSH WIT GET_FILE_WRITER_1
WRITERS DO PUSH WIT GET_FILE_WRITER_2
WRITERS DO PUSH WIT GET_FILE_WRITER_3
WHILE NO SAEM AS (WRITERS LENGTH SAEM AS 0)
I HAS A VARIABLE WRITER TEH WRITER ITZ WRITERS DO POP
WRITER DO WRITE WIT "Data for writer"
WRITER DO CLOSE
KTHX
SAYZ WIT "All writers closed"
```

**Example: Close with final flush**

```lol
I HAS A VARIABLE WRITER TEH WRITER ITZ GET_BUFFERED_WRITER
WRITER DO WRITE WIT "Data 1"
WRITER DO WRITE WIT "Data 2"
WRITER DO WRITE WIT "Data 3"
WRITER DO CLOSE
SAYZ WIT "All buffered data flushed and writer closed"
```

**Example: Resource management pattern**

```lol
I HAS A VARIABLE WRITER TEH WRITER ITZ GET_DATABASE_WRITER
I HAS A VARIABLE RESOURCE_COUNT TEH INTEGR ITZ 1
MAYB
WRITER DO WRITE WIT "INSERT INTO table VALUES (...)"
WRITER DO WRITE WIT "COMMIT"
OOPSIE ERR
SAYZ WIT "Database write failed: "
SAYZ WIT ERR
WRITER DO WRITE WIT "ROLLBACK"
KTHX
RESOURCE_COUNT ITZ RESOURCE_COUNT MINUS 1
WRITER DO CLOSE
SAYZ WIT "Database connection closed, resources: "
SAYZ WIT RESOURCE_COUNT
```

**Note:** Always call CLOSE to prevent resource leaks

**Note:** CLOSE ensures all buffered data is written

**Note:** CLOSE should be called even if errors occur during writing

**Note:** Multiple CLOSE calls should be safe (idempotent)

**Note:** After CLOSE, further WRITE operations may fail

**Note:** CLOSE may finalize file headers, network connections, etc.

#### WRITE

Writes string data to the output destination.
Returns the number of characters written.

**Syntax:** `<writer> DO WRITE WIT <data>`
**Parameters:**
- `data` (STRIN): String data to write

**Example: Write simple text**

```lol
I HAS A VARIABLE WRITER TEH WRITER ITZ GET_FILE_WRITER
I HAS A VARIABLE CHARS_WRITTEN TEH INTEGR ITZ WRITER DO WRITE WIT "Hello World"
SAYZ WIT "Wrote "
SAYZ WIT CHARS_WRITTEN
SAYZ WIT " characters"
WRITER DO CLOSE
```

**Example: Write with validation**

```lol
I HAS A VARIABLE WRITER TEH WRITER ITZ GET_NETWORK_WRITER
I HAS A VARIABLE DATA TEH STRIN ITZ "Important message"
I HAS A VARIABLE WRITTEN TEH INTEGR ITZ WRITER DO WRITE WIT DATA
IZ WRITTEN SAEM AS DATA LENGTH?
SAYZ WIT "All data written successfully"
NOPE
SAYZ WIT "Warning: Only wrote "
SAYZ WIT WRITTEN
SAYZ WIT " of "
SAYZ WIT DATA LENGTH
SAYZ WIT " characters"
KTHX
WRITER DO CLOSE
```

**Example: Write formatted data**

```lol
I HAS A VARIABLE WRITER TEH WRITER ITZ GET_LOG_WRITER
I HAS A VARIABLE TIMESTAMP TEH STRIN ITZ GET_CURRENT_TIME
I HAS A VARIABLE LEVEL TEH STRIN ITZ "INFO"
I HAS A VARIABLE MESSAGE TEH STRIN ITZ "Process started"
WRITER DO WRITE WIT TIMESTAMP
WRITER DO WRITE WIT " ["
WRITER DO WRITE WIT LEVEL
WRITER DO WRITE WIT "] "
WRITER DO WRITE WIT MESSAGE
WRITER DO WRITE WIT "\n"
WRITER DO CLOSE
```

**Example: Write binary-like data**

```lol
I HAS A VARIABLE WRITER TEH WRITER ITZ GET_BINARY_WRITER
I HAS A VARIABLE HEADER TEH STRIN ITZ "\x00\x01\x02\x03"
I HAS A VARIABLE PAYLOAD TEH STRIN ITZ "\x04\x05\x06\x07"
WRITER DO WRITE WIT HEADER
WRITER DO WRITE WIT PAYLOAD
I HAS A VARIABLE TOTAL_WRITTEN TEH INTEGR ITZ HEADER LENGTH MOAR PAYLOAD LENGTH
SAYZ WIT "Binary data written: "
SAYZ WIT TOTAL_WRITTEN
SAYZ WIT " bytes"
WRITER DO CLOSE
```

**Example: Batch writing with progress**

```lol
I HAS A VARIABLE WRITER TEH WRITER ITZ GET_FILE_WRITER
I HAS A VARIABLE LINES TEH BUKKIT ITZ NEW BUKKIT
LINES DO PUSH WIT "Line 1"
LINES DO PUSH WIT "Line 2"
LINES DO PUSH WIT "Line 3"
I HAS A VARIABLE TOTAL_CHARS TEH INTEGR ITZ 0
WHILE NO SAEM AS (LINES LENGTH SAEM AS 0)
I HAS A VARIABLE LINE TEH STRIN ITZ LINES DO POP
I HAS A VARIABLE LINE_WITH_NEWLINE TEH STRIN ITZ LINE MOAR "\n"
I HAS A VARIABLE WRITTEN TEH INTEGR ITZ WRITER DO WRITE WIT LINE_WITH_NEWLINE
TOTAL_CHARS ITZ TOTAL_CHARS MOAR WRITTEN
KTHX
SAYZ WIT "Total characters written: "
SAYZ WIT TOTAL_CHARS
WRITER DO CLOSE
```

**Note:** Returns the number of characters actually written

**Note:** May write fewer characters than provided if error occurs

**Note:** Data may be buffered and not immediately written

**Note:** Call CLOSE to ensure all data is flushed

**Note:** Implementation depends on the concrete writer type

**Example: Basic writer usage pattern**

```lol
I HAS A VARIABLE WRITER TEH WRITER ITZ GET_SOME_WRITER
WRITER DO WRITE WIT "Hello, World!"
WRITER DO CLOSE
SAYZ WIT "Data written successfully"
```

**Example: Writing with error handling**

```lol
I HAS A VARIABLE WRITER TEH WRITER ITZ GET_FILE_WRITER
MAYB
WRITER DO WRITE WIT "Important data"
SAYZ WIT "Data written successfully"
OOPSIE ERR
SAYZ WIT "Error writing data: "
SAYZ WIT ERR
KTHX
WRITER DO CLOSE
```

**Example: Writing multiple pieces of data**

```lol
I HAS A VARIABLE WRITER TEH WRITER ITZ GET_LOG_WRITER
I HAS A VARIABLE MESSAGES TEH BUKKIT ITZ NEW BUKKIT
MESSAGES DO PUSH WIT "Starting process..."
MESSAGES DO PUSH WIT "Processing data..."
MESSAGES DO PUSH WIT "Process complete"
WHILE NO SAEM AS (MESSAGES LENGTH SAEM AS 0)
I HAS A VARIABLE MSG TEH STRIN ITZ MESSAGES DO POP
WRITER DO WRITE WIT MSG
WRITER DO WRITE WIT "\n"
KTHX
WRITER DO CLOSE
```

