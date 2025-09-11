# FILE Module - File System Operations

The FILE module provides file system operations through the DOCUMENT class for file I/O and the CABINET class for directory operations.

## Importing FILE Module

```lol
BTW Import entire module
I CAN HAS FILE?

BTW Selective import
I CAN HAS DOCUMENT FROM FILE?
I CAN HAS CABINET FROM FILE?
I CAN HAS SEP FROM FILE?
```

**Note:** The FILE module automatically imports the IO module since DOCUMENT inherits from READWRITER.

## Path Constants

### SEP - Path Separator

The platform-specific path separator character.

**Type:** STRIN
**Value:** `"/"` on Unix/Linux/macOS, `"\"` on Windows

```lol
I CAN HAS SEP FROM FILE?

I HAS A VARIABLE PATH TEH STRIN ITZ "home" SMOOSH SEP SMOOSH "user" SMOOSH SEP SMOOSH "documents"
I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT PATH AN WIT "R"
```

## DOCUMENT Class

The DOCUMENT class represents a file on the file system and provides methods for file I/O operations. It inherits from IO.READWRITER, providing both reading and writing capabilities.

### Constructor

```lol
I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT <path> AN WIT <mode>
```

**Parameters:**
- **path**: STRIN - The file path (absolute or relative)
- **mode**: STRIN - The file access mode

**File Access Modes:**

| Mode | Description | Operations |
|------|-------------|------------|
| `"R"` | Read-only | Can read from existing file |
| `"W"` | Write-only | Creates new file or overwrites existing |
| `"RW"` | Read-write | Creates if needed, can read and write |
| `"A"` | Append | Creates if needed, writes append to end |

### Properties

- **PATH**: STRIN (read-only) - The file path
- **MODE**: STRIN (read-only) - The access mode
- **IS_OPEN**: BOOL - True if file is currently open
- **SIZ**: INTEGR (read-only) - File size in bytes
- **RWX**: INTEGR - File permissions (read/write/execute bits)

### Methods

#### File Lifecycle

##### OPEN - Open File

Opens the file for I/O operations according to the specified mode.

```lol
document DO OPEN
```

- Creates the file if it doesn't exist (for write/append modes)
- Throws exception if file cannot be opened
- Sets IS_OPEN to YEZ

##### CLOSE - Close File

Closes the file and releases system resources.

```lol
document DO CLOSE
```

- Flushes any pending writes
- Sets IS_OPEN to NO
- Safe to call multiple times

#### Reading Operations (Mode R or RW)

##### READ - Read Data

Reads up to the specified number of characters from the file.

```lol
I HAS A VARIABLE DATA TEH STRIN ITZ document DO READ WIT <size>
```

**Parameters:**
- **size**: INTEGR - Maximum number of characters to read

**Returns:** STRIN - The data read (may be shorter than requested at end of file)

#### Writing Operations (Mode W, RW, or A)

##### WRITE - Write Data

Writes string data to the file.

```lol
I HAS A VARIABLE BYTES_WRITTEN TEH INTEGR ITZ document DO WRITE WIT <data>
```

**Parameters:**
- **data**: STRIN - The data to write

**Returns:** INTEGR - Number of characters written

##### FLUSH - Force Write

Forces any buffered data to be written to disk immediately.

```lol
document DO FLUSH
```

#### File Position Operations

##### SEEK - Set Position

Sets the file position for next read/write operation.

```lol
document DO SEEK WIT <position>
```

**Parameters:**
- **position**: INTEGR - Byte position from start of file (0-based)

##### TELL - Get Position

Gets the current file position.

```lol
I HAS A VARIABLE POSITION TEH INTEGR ITZ document DO TELL
```

**Returns:** INTEGR - Current byte position in file

#### File Information


##### EXISTS - Check File Existence

Checks if the file exists on disk.

```lol
I HAS A VARIABLE FILE_EXISTS TEH BOOL ITZ document DO EXISTS
```

**Returns:** BOOL - YEZ if file exists, NO otherwise

#### File Management

##### DELETE - Delete File

Deletes the file from disk.

```lol
document DO DELETE
```

- Automatically closes file if open
- Throws exception if deletion fails
- Sets IS_OPEN to NO

## CABINET Class

The CABINET class provides directory operations for working with directories and their contents.

### Constructor

```lol
I HAS A VARIABLE DIR TEH CABINET ITZ NEW CABINET WIT <path>
```

**Parameters:**
- **path**: STRIN - The directory path

### Properties

- **PATH**: STRIN (read-only) - The directory path

### Methods

#### EXISTS - Check Directory Existence

Checks if the directory exists.

```lol
I HAS A VARIABLE DIR_EXISTS TEH BOOL ITZ cabinet DO EXISTS
```

**Returns:** BOOL - YEZ if directory exists, NO otherwise

#### LIST - List Directory Contents

Returns all files and subdirectories in the directory.

```lol
I HAS A VARIABLE FILES TEH BUKKIT ITZ cabinet DO LIST
```

**Returns:** BUKKIT - Array of filenames and directory names

#### CREATE - Create Directory

Creates the directory (including parent directories if needed).

```lol
cabinet DO CREATE
```

- Creates all necessary parent directories
- Throws exception if creation fails

#### DELETE - Delete Directory

Deletes an empty directory.

```lol
cabinet DO DELETE
```

- Only deletes empty directories
- Throws exception if directory is not empty or deletion fails

#### FIND - Find Files by Pattern

Searches for files matching a glob pattern.

```lol
I HAS A VARIABLE MATCHES TEH BUKKIT ITZ cabinet DO FIND WIT <pattern>
```

**Parameters:**
- **pattern**: STRIN - Glob pattern (e.g., "*.txt", "test*")

**Returns:** BUKKIT - Array of matching filenames

## Basic File Operations

### Reading a Text File

```lol
I CAN HAS FILE?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN READ_TEXT_FILE WIT FILENAME TEH STRIN
    SAYZ WIT "=== Reading Text File ==="

    BTW Create document for reading
    I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT FILENAME AN WIT "R"

    BTW Check if file exists before opening
    I HAS A VARIABLE EXISTS TEH BOOL ITZ DOC DO EXISTS
    IZ EXISTS?
        SAYZ WIT "File exists, opening..."
        DOC DO OPEN

        BTW Get file size
        I HAS A VARIABLE SIZE TEH INTEGR ITZ DOC SIZ
        SAY WIT "File size: "
        SAYZ WIT SIZE

        BTW Read entire file content
        I HAS A VARIABLE CONTENT TEH STRIN ITZ DOC DO READ WIT SIZE
        SAYZ WIT "File content:"
        SAYZ WIT CONTENT

        BTW Close file
        DOC DO CLOSE
        SAYZ WIT "File closed"
    NOPE
        SAYZ WIT "File does not exist!"
    KTHX
KTHXBAI
```

### Writing a Text File

```lol
I CAN HAS FILE?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN WRITE_TEXT_FILE WIT FILENAME TEH STRIN AN WIT CONTENT TEH STRIN
    SAYZ WIT "=== Writing Text File ==="

    BTW Create document for writing
    I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT FILENAME AN WIT "W"

    BTW Open file for writing
    DOC DO OPEN
    SAYZ WIT "File opened for writing"

    BTW Write content
    I HAS A VARIABLE BYTES_WRITTEN TEH INTEGR ITZ DOC DO WRITE WIT CONTENT
    SAY WIT "Wrote "
    SAY WIT BYTES_WRITTEN
    SAYZ WIT " characters"

    BTW Flush to ensure data is written
    DOC DO FLUSH
    SAYZ WIT "Data flushed to disk"

    BTW Close file
    DOC DO CLOSE
    SAYZ WIT "File closed"
KTHXBAI

HAI ME TEH FUNCSHUN DEMO_WRITE_FILE
    I HAS A VARIABLE TEXT TEH STRIN ITZ "Hello, World!\nThis is a test file.\nCreated with Objective-LOL!"
    WRITE_TEXT_FILE WIT "example.txt" AN WIT TEXT
KTHXBAI
```

### Appending to a File

```lol
I CAN HAS FILE?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN APPEND_TO_FILE WIT FILENAME TEH STRIN AN WIT NEW_CONTENT TEH STRIN
    SAYZ WIT "=== Appending to File ==="

    BTW Create document for appending
    I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT FILENAME AN WIT "A"

    BTW Open in append mode
    DOC DO OPEN
    SAYZ WIT "File opened for appending"

    BTW Write new content
    I HAS A VARIABLE BYTES_WRITTEN TEH INTEGR ITZ DOC DO WRITE WIT NEW_CONTENT
    SAY WIT "Appended "
    SAY WIT BYTES_WRITTEN
    SAYZ WIT " characters"

    BTW Close file
    DOC DO CLOSE
    SAYZ WIT "File closed"
KTHXBAI

HAI ME TEH FUNCSHUN DEMO_APPEND
    APPEND_TO_FILE WIT "example.txt" AN WIT "\nThis line was appended!"
KTHXBAI
```

## Advanced File Operations

### Reading File in Chunks

```lol
I CAN HAS FILE?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN READ_FILE_IN_CHUNKS WIT FILENAME TEH STRIN AN WIT CHUNK_SIZE TEH INTEGR
    SAYZ WIT "=== Reading File in Chunks ==="

    I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT FILENAME AN WIT "R"

    IZ DOC DO EXISTS?
        DOC DO OPEN

        I HAS A VARIABLE CHUNK_COUNT TEH INTEGR ITZ 0
        I HAS A VARIABLE TOTAL_READ TEH INTEGR ITZ 0

        WHILE YEZ  BTW Read until EOF
            I HAS A VARIABLE CHUNK TEH STRIN ITZ DOC DO READ WIT CHUNK_SIZE
            I HAS A VARIABLE CHUNK_LEN TEH INTEGR ITZ LEN WIT CHUNK

            IZ CHUNK_LEN SAEM AS 0?
                SAYZ WIT "End of file reached"
                BREAK  BTW Would use proper break in real implementation
            KTHX

            CHUNK_COUNT ITZ CHUNK_COUNT MOAR 1
            TOTAL_READ ITZ TOTAL_READ MOAR CHUNK_LEN

            SAY WIT "Chunk "
            SAY WIT CHUNK_COUNT
            SAY WIT " (size "
            SAY WIT CHUNK_LEN
            SAY WIT "): "
            SAYZ WIT CHUNK

            BTW Safety exit for demo
            IZ CHUNK_COUNT BIGGR THAN 5?
                SAYZ WIT "Demo limit reached"
                GIVEZ UP
            KTHX
        KTHX

        DOC DO CLOSE
        SAY WIT "Total chunks read: "
        SAY WIT CHUNK_COUNT
        SAY WIT ", Total characters: "
        SAYZ WIT TOTAL_READ
    NOPE
        SAYZ WIT "File does not exist!"
    KTHX
KTHXBAI
```

### File Position Operations

```lol
I CAN HAS FILE?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN DEMO_FILE_POSITIONING WIT FILENAME TEH STRIN
    SAYZ WIT "=== File Positioning Demo ==="

    I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT FILENAME AN WIT "RW"

    IZ DOC DO EXISTS SAEM AS NO?
        BTW Create file with sample content
        DOC DO OPEN
        DOC DO WRITE WIT "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
        DOC DO CLOSE
        SAYZ WIT "Created sample file"
    KTHX

    DOC DO OPEN

    BTW Show initial position
    I HAS A VARIABLE POS TEH INTEGR ITZ DOC DO TELL
    SAY WIT "Initial position: "
    SAYZ WIT POS

    BTW Read from beginning
    I HAS A VARIABLE DATA1 TEH STRIN ITZ DOC DO READ WIT 10
    SAY WIT "Read from start: "
    SAYZ WIT DATA1

    I HAS A VARIABLE POS2 TEH INTEGR ITZ DOC DO TELL
    SAY WIT "Position after read: "
    SAYZ WIT POS2

    BTW Seek to middle
    DOC DO SEEK WIT 15
    I HAS A VARIABLE POS3 TEH INTEGR ITZ DOC DO TELL
    SAY WIT "Position after seek to 15: "
    SAYZ WIT POS3

    BTW Read from new position
    I HAS A VARIABLE DATA2 TEH STRIN ITZ DOC DO READ WIT 5
    SAY WIT "Read from position 15: "
    SAYZ WIT DATA2

    BTW Seek back to beginning and write
    DOC DO SEEK WIT 0
    DOC DO WRITE WIT "START"

    DOC DO CLOSE
    SAYZ WIT "File positioning demo completed"
KTHXBAI
```

### File Management Operations

```lol
I CAN HAS FILE?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN DEMO_FILE_MANAGEMENT
    SAYZ WIT "=== File Management Demo ==="

    I HAS A VARIABLE FILENAME TEH STRIN ITZ "temp_file.txt"
    I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT FILENAME AN WIT "W"

    BTW Check if file exists before creation
    I HAS A VARIABLE EXISTS_BEFORE TEH BOOL ITZ DOC DO EXISTS
    SAY WIT "File exists before creation: "
    SAYZ WIT EXISTS_BEFORE

    BTW Create and write to file
    DOC DO OPEN
    DOC DO WRITE WIT "This is a temporary file for testing."
    DOC DO CLOSE

    BTW Check existence and size after creation
    I HAS A VARIABLE EXISTS_AFTER TEH BOOL ITZ DOC DO EXISTS
    SAY WIT "File exists after creation: "
    SAYZ WIT EXISTS_AFTER

    I HAS A VARIABLE SIZE TEH INTEGR ITZ DOC SIZ
    SAY WIT "File size: "
    SAYZ WIT SIZE

    BTW Delete the file
    DOC DO DELETE
    SAYZ WIT "File deleted"

    BTW Verify deletion
    I HAS A VARIABLE EXISTS_FINAL TEH BOOL ITZ DOC DO EXISTS
    SAY WIT "File exists after deletion: "
    SAYZ WIT EXISTS_FINAL
KTHXBAI
```

## Working with Buffered I/O

The DOCUMENT class can be used with the IO module's buffered classes for improved performance:

```lol
I CAN HAS FILE?
I CAN HAS IO?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN DEMO_BUFFERED_FILE_IO
    SAYZ WIT "=== Buffered File I/O Demo ==="

    I HAS A VARIABLE FILENAME TEH STRIN ITZ "buffered_test.txt"

    BTW Write with buffered writer
    I HAS A VARIABLE WRITE_DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT FILENAME AN WIT "W"
    WRITE_DOC DO OPEN

    I HAS A VARIABLE BUF_WRITER TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT WRITE_DOC

    BTW Write multiple small pieces efficiently
    BUF_WRITER DO WRITE WIT "Line 1\n"
    BUF_WRITER DO WRITE WIT "Line 2\n"
    BUF_WRITER DO WRITE WIT "Line 3\n"
    BUF_WRITER DO WRITE WIT "Line 4\n"

    BTW Flush and close
    BUF_WRITER DO CLOSE  BTW This also closes the underlying document
    SAYZ WIT "Buffered writing completed"

    BTW Read with buffered reader
    I HAS A VARIABLE READ_DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT FILENAME AN WIT "R"
    READ_DOC DO OPEN

    I HAS A VARIABLE BUF_READER TEH BUFFERED_READER ITZ NEW BUFFERED_READER WIT READ_DOC

    BTW Read data efficiently
    I HAS A VARIABLE CONTENT TEH STRIN ITZ BUF_READER DO READ WIT 1000
    SAYZ WIT "Buffered read content:"
    SAYZ WIT CONTENT

    BUF_READER DO CLOSE  BTW This also closes the underlying document
    SAYZ WIT "Buffered reading completed"

    BTW Clean up
    I HAS A VARIABLE CLEANUP_DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT FILENAME AN WIT "R"
    CLEANUP_DOC DO DELETE
    SAYZ WIT "Test file deleted"
KTHXBAI
```

## Error Handling

### File Operation Error Handling

```lol
I CAN HAS FILE?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN SAFE_FILE_OPERATIONS WIT FILENAME TEH STRIN
    SAYZ WIT "=== Safe File Operations ==="

    I HAS A VARIABLE DOC TEH DOCUMENT ITZ NOTHIN

    MAYB
        BTW Try to open file
        DOC ITZ NEW DOCUMENT WIT FILENAME AN WIT "R"

        IZ DOC DO EXISTS?
            DOC DO OPEN

            BTW Try to read data
            I HAS A VARIABLE CONTENT TEH STRIN ITZ DOC DO READ WIT 100
            SAYZ WIT "Successfully read file content"

        NOPE
            SAYZ WIT "File does not exist"
        KTHX

    OOPSIE FILE_ERROR
        SAYZ WIT "File operation error: "
        SAYZ WIT FILE_ERROR

    ALWAYZ
        BTW Always close file if it was opened
        IZ DOC SAEM AS NOTHIN SAEM AS NO?
            IZ DOC IS_OPEN?
                DOC DO CLOSE
                SAYZ WIT "File closed in cleanup"
            KTHX
        KTHX
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN DEMO_FILE_ERRORS
    BTW Test with non-existent file
    SAFE_FILE_OPERATIONS WIT "nonexistent_file.txt"

    BTW Test with invalid mode
    MAYB
        I HAS A VARIABLE BAD_DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT "test.txt" AN WIT "INVALID"
    OOPSIE MODE_ERROR
        SAYZ WIT "Invalid mode error (expected): "
        SAYZ WIT MODE_ERROR
    KTHX
KTHXBAI
```

### Read-only File Writing Error

```lol
I CAN HAS FILE?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN DEMO_READ_ONLY_ERROR
    SAYZ WIT "=== Read-Only File Error Demo ==="

    BTW Try to write to a read-only file
    I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT "readonly_test.txt" AN WIT "R"

    MAYB
        DOC DO OPEN
        DOC DO WRITE WIT "This should fail"
    OOPSIE WRITE_ERROR
        SAYZ WIT "Expected write error (read-only file): "
        SAYZ WIT WRITE_ERROR
    ALWAYZ
        IZ DOC IS_OPEN?
            DOC DO CLOSE
        KTHX
    KTHX
KTHXBAI
```

## Data Processing Examples

### Log File Analysis

```lol
I CAN HAS FILE?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN ANALYZE_LOG_FILE WIT LOG_FILE TEH STRIN
    SAYZ WIT "=== Log File Analysis ==="

    I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT LOG_FILE AN WIT "R"

    IZ DOC DO EXISTS?
        DOC DO OPEN

        I HAS A VARIABLE LINE_COUNT TEH INTEGR ITZ 0
        I HAS A VARIABLE ERROR_COUNT TEH INTEGR ITZ 0
        I HAS A VARIABLE WARNING_COUNT TEH INTEGR ITZ 0

        BTW Read file content
        I HAS A VARIABLE SIZE TEH INTEGR ITZ DOC SIZ
        I HAS A VARIABLE CONTENT TEH STRIN ITZ DOC DO READ WIT SIZE

        BTW Simple analysis (would need string processing functions)
        BTW This is a simplified example
        LINE_COUNT ITZ 10  BTW Mock count
        ERROR_COUNT ITZ 2  BTW Mock count
        WARNING_COUNT ITZ 5  BTW Mock count

        DOC DO CLOSE

        BTW Report results
        SAYZ WIT "Log Analysis Results:"
        SAY WIT "Total lines: "
        SAYZ WIT LINE_COUNT
        SAY WIT "Errors: "
        SAYZ WIT ERROR_COUNT
        SAY WIT "Warnings: "
        SAYZ WIT WARNING_COUNT

    NOPE
        SAYZ WIT "Log file not found!"
    KTHX
KTHXBAI
```

### Configuration File Manager

```lol
I CAN HAS FILE?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN CREATE_CONFIG_FILE WIT CONFIG_FILE TEH STRIN
    SAYZ WIT "=== Creating Configuration File ==="

    I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT CONFIG_FILE AN WIT "W"
    DOC DO OPEN

    BTW Write configuration data
    DOC DO WRITE WIT "# Application Configuration\n"
    DOC DO WRITE WIT "server_port=8080\n"
    DOC DO WRITE WIT "debug_mode=true\n"
    DOC DO WRITE WIT "database_url=localhost:5432\n"
    DOC DO WRITE WIT "max_connections=100\n"
    DOC DO WRITE WIT "log_level=info\n"

    DOC DO CLOSE
    SAYZ WIT "Configuration file created"
KTHXBAI

HAI ME TEH FUNCSHUN READ_CONFIG_FILE WIT CONFIG_FILE TEH STRIN
    SAYZ WIT "=== Reading Configuration File ==="

    I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT CONFIG_FILE AN WIT "R"

    IZ DOC DO EXISTS?
        DOC DO OPEN

        I HAS A VARIABLE SIZE TEH INTEGR ITZ DOC SIZ
        I HAS A VARIABLE CONFIG_CONTENT TEH STRIN ITZ DOC DO READ WIT SIZE

        SAYZ WIT "Configuration file contents:"
        SAYZ WIT CONFIG_CONTENT

        DOC DO CLOSE
    NOPE
        SAYZ WIT "Configuration file not found!"
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN DEMO_CONFIG_MANAGER
    I HAS A VARIABLE CONFIG_FILE TEH STRIN ITZ "app_config.txt"

    CREATE_CONFIG_FILE WIT CONFIG_FILE
    READ_CONFIG_FILE WIT CONFIG_FILE

    BTW Clean up
    I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT CONFIG_FILE AN WIT "R"
    DOC DO DELETE
    SAYZ WIT "Configuration file deleted"
KTHXBAI
```

## Directory Operations Examples

### Basic Directory Operations

```lol
I CAN HAS FILE?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN DEMO_DIRECTORY_OPERATIONS
    SAYZ WIT "=== Directory Operations Demo ==="

    BTW Create a cabinet for a directory
    I HAS A VARIABLE DIR TEH CABINET ITZ NEW CABINET WIT "test_directory"

    BTW Check if directory exists
    I HAS A VARIABLE EXISTS_BEFORE TEH BOOL ITZ DIR DO EXISTS
    SAY WIT "Directory exists before creation: "
    SAYZ WIT EXISTS_BEFORE

    BTW Create the directory
    DIR DO CREATE
    SAYZ WIT "Directory created"

    BTW Check existence after creation
    I HAS A VARIABLE EXISTS_AFTER TEH BOOL ITZ DIR DO EXISTS
    SAY WIT "Directory exists after creation: "
    SAYZ WIT EXISTS_AFTER

    BTW List contents (should be empty)
    I HAS A VARIABLE CONTENTS TEH BUKKIT ITZ DIR DO LIST
    SAY WIT "Directory contents count: "
    SAYZ WIT CONTENTS SIZ

    BTW Find files matching a pattern
    I HAS A VARIABLE TXT_FILES TEH BUKKIT ITZ DIR DO FIND WIT "*.txt"
    SAY WIT "Text files found: "
    SAYZ WIT TXT_FILES SIZ

    BTW Clean up - delete the directory
    DIR DO DELETE
    SAYZ WIT "Directory deleted"
KTHXBAI
```

### Working with File Permissions

```lol
I CAN HAS FILE?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN DEMO_FILE_PERMISSIONS WIT FILENAME TEH STRIN
    SAYZ WIT "=== File Permissions Demo ==="

    BTW Create a file
    I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT FILENAME AN WIT "W"
    DOC DO OPEN
    DOC DO WRITE WIT "Test file for permissions"
    DOC DO CLOSE

    BTW Check current permissions
    I HAS A VARIABLE CURRENT_PERMS TEH INTEGR ITZ DOC RWX
    SAY WIT "Current permissions: "
    SAYZ WIT CURRENT_PERMS

    BTW Change permissions (make read-only: 0644 octal)
    DOC RWX ITZ 0644
    SAYZ WIT "Permissions changed to read-only"

    BTW Verify new permissions
    I HAS A VARIABLE NEW_PERMS TEH INTEGR ITZ DOC RWX
    SAY WIT "New permissions: "
    SAYZ WIT NEW_PERMS

    BTW Clean up
    DOC RWX ITZ 0666  BTW Allow deletion
    DOC DO DELETE
    SAYZ WIT "File deleted"
KTHXBAI
```

## Quick Reference

### Constants

| Constant | Type | Description |
|----------|------|-------------|
| `SEP` | STRIN | Platform-specific path separator (`/` or `\`) |

### DOCUMENT Constructor

| Usage | Description |
|-------|-------------|
| `NEW DOCUMENT WIT path AN WIT "R"` | Read-only access |
| `NEW DOCUMENT WIT path AN WIT "W"` | Write access (overwrites) |
| `NEW DOCUMENT WIT path AN WIT "RW"` | Read-write access |
| `NEW DOCUMENT WIT path AN WIT "A"` | Append access |

### File Lifecycle

| Method | Description |
|--------|-------------|
| `OPEN` | Open file for operations |
| `CLOSE` | Close file and release resources |

### I/O Operations

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `READ WIT size` | size: INTEGR | STRIN | Read up to size characters |
| `WRITE WIT data` | data: STRIN | INTEGR | Write string data |
| `FLUSH` | - | - | Force write to disk |

### Position Operations

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `SEEK WIT position` | position: INTEGR | - | Set file position |
| `TELL` | - | INTEGR | Get current position |

### File Information

| Method | Returns | Description |
|--------|---------|-------------|
| `EXISTS` | BOOL | Check if file exists |

### File Management

| Method | Description |
|--------|-------------|
| `DELETE` | Delete file from disk |

### Properties

| Property | Type | Description |
|----------|------|-------------|
| `PATH` | STRIN | File path (read-only) |
| `MODE` | STRIN | Access mode (read-only) |
| `IS_OPEN` | BOOL | True if file is open |
| `SIZ` | INTEGR | File size in bytes (read-only) |
| `RWX` | INTEGR | File permissions (read/write/execute bits) |

## Related

- [IO Module](io.md) - Advanced I/O with buffering
- [STDIO Module](stdio.md) - Console input/output
- [Control Flow](../language-guide/control-flow.md) - Exception handling patterns