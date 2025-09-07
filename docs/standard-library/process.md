# PROCESS Module - Process Management Operations

The PROCESS module provides process management functionality through the MINION class for launching and controlling processes and the PIPE class for inter-process communication.

## Importing PROCESS Module

```lol
BTW Import entire module
I CAN HAS PROCESS?

BTW Selective import
I CAN HAS MINION FROM PROCESS?
I CAN HAS PIPE FROM PROCESS?
```

**Note:** The PROCESS module automatically imports the IO module classes (READER, WRITER) when PIPE is imported.

## MINION Class

The MINION class represents a child process that can be launched, monitored, and controlled. It provides process lifecycle management with configurable working directory and environment variables.

### Constructor

```lol
I HAS A VARIABLE CMD TEH BUKKIT ITZ NEW BUKKIT
CMD DO PUSH WIT "ls"
CMD DO PUSH WIT "-la"
I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMD
```

The constructor creates a process definition with:
- **Command Line**: Command and arguments as BUKKIT
- **Working Directory**: Current directory (configurable)
- **Environment**: Current environment variables (configurable)

### Properties

- **CMDLINE**: BUKKIT (read-only) - Command and arguments
- **RUNNING**: BOOL (read-only) - Whether process is currently running
- **FINISHED**: BOOL (read-only) - Whether process has completed
- **EXIT_CODE**: INTEGR (read-only) - Process exit code
- **PID**: INTEGR (read-only) - Process ID (-1 if not started)
- **STDIN**: PIPE (read-only) - Stdin pipe (available after START)
- **STDOUT**: PIPE (read-only) - Stdout pipe (available after START)
- **STDERR**: PIPE (read-only) - Stderr pipe (available after START)

### Methods

#### Configuration Methods

##### SET_WORKDIR - Set Working Directory

Sets the working directory for the process before starting.

```lol
proc DO SET_WORKDIR WIT "/path/to/directory"
```

**Parameters:**
- **dir**: STRIN - The working directory path

**Note:** Cannot be changed after process starts

##### SET_ENV - Set Environment Variables

Replaces all environment variables with the provided BASKIT.

```lol
I HAS A VARIABLE ENV TEH BASKIT ITZ NEW BASKIT
ENV DO PUT WIT "PATH" AN WIT "/usr/bin"
proc DO SET_ENV WIT ENV
```

**Parameters:**
- **env**: BASKIT - Environment variables as key-value pairs

**Note:** Cannot be changed after process starts

##### ADD_ENV - Add Environment Variable

Adds or updates a single environment variable.

```lol
proc DO ADD_ENV WIT "MY_VAR" AN WIT "value"
```

**Parameters:**
- **key**: STRIN - Environment variable name
- **value**: STRIN - Environment variable value

**Note:** Cannot be changed after process starts

#### Process Control Methods

##### START - Start Process

Launches the process and creates stdin, stdout, and stderr pipes.

```lol
proc DO START
```

**Throws:** Exception if process fails to start or is already running

##### WAIT - Wait for Completion

Waits for the process to complete and returns the exit code.

```lol
I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ proc DO WAIT
```

**Returns:** INTEGR - Process exit code

**Throws:** Exception if process has not been started

##### KILL - Terminate Process

Forcefully terminates the running process.

```lol
proc DO KILL
```

**Throws:** Exception if process is not running

##### SIGNAL - Send Signal

Sends a signal to the running process (Unix/Linux systems).

```lol
proc DO SIGNAL WIT 15  BTW SIGTERM
```

**Parameters:**
- **code**: INTEGR - Signal number to send

**Throws:** Exception if process is not running

##### IS_ALIVE - Check Process Status

Returns whether the process is currently running.

```lol
I HAS A VARIABLE ALIVE TEH BOOL ITZ proc DO IS_ALIVE
```

**Returns:** BOOL - True if process is running

## PIPE Class

The PIPE class represents a communication pipe connected to a process's stdin, stdout, or stderr. It implements both READER and WRITER interfaces from the IO module.

### Properties

- **FD_TYPE**: STRIN (read-only) - Pipe type: "STDIN", "STDOUT", or "STDERR"
- **IS_OPEN**: BOOL (read-only) - Whether pipe is open for operations
- **IS_EOF**: BOOL (read-only) - Whether end-of-file reached (read pipes only)

### Methods

##### READ - Read from Pipe

Reads data from stdout or stderr pipes.

```lol
I HAS A VARIABLE DATA TEH STRIN ITZ pipe DO READ WIT 1024
```

**Parameters:**
- **size**: INTEGR - Maximum bytes to read

**Returns:** STRIN - Data read from pipe

**Throws:** Exception if pipe is not open or is stdin pipe

##### WRITE - Write to Pipe

Writes data to stdin pipe.

```lol
I HAS A VARIABLE BYTES_WRITTEN TEH INTEGR ITZ pipe DO WRITE WIT "Hello\n"
```

**Parameters:**
- **data**: STRIN - Data to write to pipe

**Returns:** INTEGR - Number of bytes written

**Throws:** Exception if pipe is not open or is stdout/stderr pipe

##### CLOSE - Close Pipe

Closes the pipe connection.

```lol
pipe DO CLOSE
```

**Note:** Automatically closes when process terminates

## Basic Process Operations

### Simple Command Execution

```lol
I CAN HAS PROCESS?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN RUN_SIMPLE_COMMAND WIT CMD TEH STRIN
    SAYZ WIT "=== Simple Command Execution ==="

    BTW Create command line
    I HAS A VARIABLE CMDLINE TEH BUKKIT ITZ NEW BUKKIT
    CMDLINE DO PUSH WIT CMD

    BTW Create and start process
    I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMDLINE

    MAYB
        PROC DO START
        SAYZ WIT "Process started"

        BTW Wait for completion
        I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC DO WAIT
        SAY WIT "Process completed with exit code: "
        SAYZ WIT EXIT_CODE

    OOPSIE PROC_ERROR
        SAY WIT "Process failed: "
        SAYZ WIT PROC_ERROR
    KTHX
KTHXBAI
```

### Reading Process Output

```lol
I CAN HAS PROCESS?
I CAN HAS STDIO?
I CAN HAS STRING?

HAI ME TEH FUNCSHUN READ_COMMAND_OUTPUT
    SAYZ WIT "=== Reading Command Output ==="

    BTW Create command to list files
    I HAS A VARIABLE CMDLINE TEH BUKKIT ITZ NEW BUKKIT
    CMDLINE DO PUSH WIT "ls"
    CMDLINE DO PUSH WIT "-la"

    I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMDLINE

    MAYB
        PROC DO START

        BTW Read all output
        I HAS A VARIABLE OUTPUT TEH STRIN ITZ ""
        I HAS A VARIABLE STDOUT_PIPE TEH PIPE ITZ PROC STDOUT
        I HAS A VARIABLE READING TEH BOOL ITZ YEZ

        WHILE READING
            MAYB
                I HAS A VARIABLE CHUNK TEH STRIN ITZ STDOUT_PIPE DO READ WIT 1024
                IZ LEN WIT CHUNK SAEM AS 0?
                    READING ITZ NO  BTW EOF reached
                NOPE
                    OUTPUT ITZ CONCAT WIT OUTPUT AN WIT CHUNK
                KTHX
            OOPSIE READ_ERROR
                READING ITZ NO  BTW Stop on read error
            KTHX
        KTHX

        BTW Wait for process completion
        I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC DO WAIT

        SAY WIT "Command output:\n"
        SAYZ WIT OUTPUT
        SAY WIT "Exit code: "
        SAYZ WIT EXIT_CODE

    OOPSIE PROC_ERROR
        SAY WIT "Process failed: "
        SAYZ WIT PROC_ERROR
    KTHX
KTHXBAI
```

### Interactive Process Communication

```lol
I CAN HAS PROCESS?
I CAN HAS STDIO?
I CAN HAS STRING?

HAI ME TEH FUNCSHUN INTERACTIVE_PROCESS
    SAYZ WIT "=== Interactive Process ==="

    BTW Launch a shell process
    I HAS A VARIABLE CMDLINE TEH BUKKIT ITZ NEW BUKKIT
    CMDLINE DO PUSH WIT "cat"  BTW Echo input

    I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMDLINE

    MAYB
        PROC DO START

        I HAS A VARIABLE STDIN_PIPE TEH PIPE ITZ PROC STDIN
        I HAS A VARIABLE STDOUT_PIPE TEH PIPE ITZ PROC STDOUT

        BTW Send input to process
        STDIN_PIPE DO WRITE WIT "Hello, Process!\n"
        STDIN_PIPE DO WRITE WIT "This is interactive communication.\n"

        BTW Close stdin to signal end of input
        STDIN_PIPE DO CLOSE

        BTW Read response
        I HAS A VARIABLE RESPONSE TEH STRIN ITZ ""
        I HAS A VARIABLE READING TEH BOOL ITZ YEZ

        WHILE READING
            MAYB
                I HAS A VARIABLE CHUNK TEH STRIN ITZ STDOUT_PIPE DO READ WIT 1024
                IZ LEN WIT CHUNK SAEM AS 0?
                    READING ITZ NO
                NOPE
                    RESPONSE ITZ CONCAT WIT RESPONSE AN WIT CHUNK
                KTHX
            OOPSIE READ_ERROR
                READING ITZ NO
            KTHX
        KTHX

        I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC DO WAIT

        SAY WIT "Process response:\n"
        SAYZ WIT RESPONSE
        SAY WIT "Exit code: "
        SAYZ WIT EXIT_CODE

    OOPSIE PROC_ERROR
        SAY WIT "Interactive process failed: "
        SAYZ WIT PROC_ERROR
    KTHX
KTHXBAI
```

## Advanced Process Operations

### Process with Custom Environment

```lol
I CAN HAS PROCESS?
I CAN HAS STDIO?
I CAN HAS STRING?

HAI ME TEH FUNCSHUN CUSTOM_ENVIRONMENT_PROCESS
    SAYZ WIT "=== Custom Environment Process ==="

    BTW Create command to print environment variable
    I HAS A VARIABLE CMDLINE TEH BUKKIT ITZ NEW BUKKIT
    CMDLINE DO PUSH WIT "sh"
    CMDLINE DO PUSH WIT "-c"
    CMDLINE DO PUSH WIT "echo $MY_CUSTOM_VAR"

    I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMDLINE

    BTW Set custom environment variable
    PROC DO ADD_ENV WIT "MY_CUSTOM_VAR" AN WIT "Hello from custom environment!"

    BTW Set working directory
    PROC DO SET_WORKDIR WIT "/tmp"

    MAYB
        PROC DO START

        I HAS A VARIABLE STDOUT_PIPE TEH PIPE ITZ PROC STDOUT
        I HAS A VARIABLE OUTPUT TEH STRIN ITZ ""
        I HAS A VARIABLE READING TEH BOOL ITZ YEZ

        WHILE READING
            MAYB
                I HAS A VARIABLE CHUNK TEH STRIN ITZ STDOUT_PIPE DO READ WIT 1024
                IZ LEN WIT CHUNK SAEM AS 0?
                    READING ITZ NO
                NOPE
                    OUTPUT ITZ CONCAT WIT OUTPUT AN WIT CHUNK
                KTHX
            OOPSIE READ_ERROR
                READING ITZ NO
            KTHX
        KTHX

        I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC DO WAIT

        SAY WIT "Environment variable output: "
        SAYZ WIT OUTPUT
        SAY WIT "Exit code: "
        SAYZ WIT EXIT_CODE

    OOPSIE PROC_ERROR
        SAY WIT "Custom environment process failed: "
        SAYZ WIT PROC_ERROR
    KTHX
KTHXBAI
```

### Process Monitoring and Control

```lol
I CAN HAS PROCESS?
I CAN HAS STDIO?
I CAN HAS TIME?

HAI ME TEH FUNCSHUN MONITORED_PROCESS
    SAYZ WIT "=== Process Monitoring ==="

    BTW Launch a long-running process
    I HAS A VARIABLE CMDLINE TEH BUKKIT ITZ NEW BUKKIT
    CMDLINE DO PUSH WIT "sleep"
    CMDLINE DO PUSH WIT "5"

    I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMDLINE

    MAYB
        PROC DO START

        SAY WIT "Process started with PID: "
        SAYZ WIT PROC PID

        BTW Monitor process status
        I HAS A VARIABLE COUNT TEH INTEGR ITZ 0

        WHILE PROC DO IS_ALIVE
            SAY WIT "Process running... ("
            SAY WIT COUNT
            SAYZ WIT ")"
            SLEEP WIT 1000  BTW Wait 1 second
            COUNT ITZ COUNT MOAR 1

            BTW Kill after 3 seconds for demo
            IZ COUNT SAEM AS 3?
                SAYZ WIT "Terminating process..."
                PROC DO KILL
            KTHX
        KTHX

        I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC DO WAIT
        SAY WIT "Process finished with exit code: "
        SAYZ WIT EXIT_CODE

    OOPSIE PROC_ERROR
        SAY WIT "Process monitoring failed: "
        SAYZ WIT PROC_ERROR
    KTHX
KTHXBAI
```

### Error Handling with Stderr

```lol
I CAN HAS PROCESS?
I CAN HAS STDIO?
I CAN HAS STRING?

HAI ME TEH FUNCSHUN PROCESS_WITH_STDERR
    SAYZ WIT "=== Process Error Handling ==="

    BTW Command that produces both stdout and stderr
    I HAS A VARIABLE CMDLINE TEH BUKKIT ITZ NEW BUKKIT
    CMDLINE DO PUSH WIT "sh"
    CMDLINE DO PUSH WIT "-c"
    CMDLINE DO PUSH WIT "echo 'Success message'; echo 'Error message' >&2; exit 1"

    I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMDLINE

    MAYB
        PROC DO START

        I HAS A VARIABLE STDOUT_PIPE TEH PIPE ITZ PROC STDOUT
        I HAS A VARIABLE STDERR_PIPE TEH PIPE ITZ PROC STDERR

        BTW Read stdout
        I HAS A VARIABLE STDOUT_OUTPUT TEH STRIN ITZ ""
        I HAS A VARIABLE READING_OUT TEH BOOL ITZ YEZ

        WHILE READING_OUT
            MAYB
                I HAS A VARIABLE CHUNK TEH STRIN ITZ STDOUT_PIPE DO READ WIT 1024
                IZ LEN WIT CHUNK SAEM AS 0?
                    READING_OUT ITZ NO
                NOPE
                    STDOUT_OUTPUT ITZ CONCAT WIT STDOUT_OUTPUT AN WIT CHUNK
                KTHX
            OOPSIE READ_ERROR
                READING_OUT ITZ NO
            KTHX
        KTHX

        BTW Read stderr
        I HAS A VARIABLE STDERR_OUTPUT TEH STRIN ITZ ""
        I HAS A VARIABLE READING_ERR TEH BOOL ITZ YEZ

        WHILE READING_ERR
            MAYB
                I HAS A VARIABLE CHUNK TEH STRIN ITZ STDERR_PIPE DO READ WIT 1024
                IZ LEN WIT CHUNK SAEM AS 0?
                    READING_ERR ITZ NO
                NOPE
                    STDERR_OUTPUT ITZ CONCAT WIT STDERR_OUTPUT AN WIT CHUNK
                KTHX
            OOPSIE READ_ERROR
                READING_ERR ITZ NO
            KTHX
        KTHX

        I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC DO WAIT

        SAY WIT "STDOUT: "
        SAYZ WIT STDOUT_OUTPUT
        SAY WIT "STDERR: "
        SAYZ WIT STDERR_OUTPUT
        SAY WIT "Exit code: "
        SAYZ WIT EXIT_CODE

    OOPSIE PROC_ERROR
        SAY WIT "Process with stderr handling failed: "
        SAYZ WIT PROC_ERROR
    KTHX
KTHXBAI
```

## Error Handling

### Process Launch Failures

```lol
I CAN HAS PROCESS?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN ROBUST_PROCESS_LAUNCH WIT COMMAND TEH STRIN
    SAYZ WIT "=== Robust Process Launch ==="

    I HAS A VARIABLE CMDLINE TEH BUKKIT ITZ NEW BUKKIT
    CMDLINE DO PUSH WIT COMMAND

    I HAS A VARIABLE PROC TEH MINION ITZ NEW MINION WIT CMDLINE

    MAYB
        SAYZ WIT "Attempting to start process..."
        PROC DO START

        SAY WIT "Process started successfully with PID: "
        SAYZ WIT PROC PID

        BTW Wait for completion
        I HAS A VARIABLE EXIT_CODE TEH INTEGR ITZ PROC DO WAIT

        IZ EXIT_CODE SAEM AS 0?
            SAYZ WIT "Process completed successfully"
        NOPE
            SAY WIT "Process failed with exit code: "
            SAYZ WIT EXIT_CODE
        KTHX

    OOPSIE LAUNCH_ERROR
        SAYZ WIT "Failed to launch process: "
        SAYZ WIT LAUNCH_ERROR
        SAYZ WIT "This could be due to:"
        SAYZ WIT "- Command not found"
        SAYZ WIT "- Permission denied"
        SAYZ WIT "- Invalid working directory"
        SAYZ WIT "- System resource limits"
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN DEMO_ERROR_HANDLING
    BTW Test with valid command
    ROBUST_PROCESS_LAUNCH WIT "echo"

    BTW Test with invalid command
    ROBUST_PROCESS_LAUNCH WIT "nonexistent-command"
KTHXBAI
```

## Quick Reference

### Constructor

| Usage | Description |
|-------|-------------|
| `NEW MINION WIT cmdline` | Create process with command line BUKKIT |

### Configuration Methods

| Method | Parameters | Description |
|--------|------------|-------------|
| `SET_WORKDIR WIT dir` | dir: STRIN | Set working directory |
| `SET_ENV WIT env` | env: BASKIT | Set environment variables |
| `ADD_ENV WIT key AN WIT value` | key: STRIN, value: STRIN | Add environment variable |

### Process Control Methods

| Method | Returns | Description |
|--------|---------|-------------|
| `START` | - | Launch the process |
| `WAIT` | INTEGR | Wait for completion, return exit code |
| `KILL` | - | Terminate process forcefully |
| `SIGNAL WIT code` | - | Send signal to process |
| `IS_ALIVE` | BOOL | Check if process is running |

### Process Properties

| Property | Type | Description |
|----------|------|-------------|
| `CMDLINE` | BUKKIT | Command line arguments |
| `RUNNING` | BOOL | Process running status |
| `FINISHED` | BOOL | Process completion status |
| `EXIT_CODE` | INTEGR | Process exit code |
| `PID` | INTEGR | Process ID |
| `STDIN` | PIPE | Stdin pipe |
| `STDOUT` | PIPE | Stdout pipe |
| `STDERR` | PIPE | Stderr pipe |

### PIPE Properties

| Property | Type | Description |
|----------|------|-------------|
| `FD_TYPE` | STRIN | Pipe type (STDIN/STDOUT/STDERR) |
| `IS_OPEN` | BOOL | Pipe open status |
| `IS_EOF` | BOOL | End-of-file status |

### PIPE Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `READ WIT size` | size: INTEGR | STRIN | Read data from pipe |
| `WRITE WIT data` | data: STRIN | INTEGR | Write data to pipe |
| `CLOSE` | - | - | Close pipe |

## Related

- [IO Module](io.md) - READER and WRITER interfaces for pipe operations
- [STDIO Module](stdio.md) - Console input/output for process interaction
- [String Module](string.md) - String manipulation for command line processing
- [Collections](collections.md) - BUKKIT and BASKIT for command arguments and environment
- [Control Flow](../language-guide/control-flow.md) - Exception handling patterns