# FILE Module

## Import

```lol
BTW Full import
I CAN HAS FILE?

BTW Selective import examples
I CAN HAS SEP FROM FILE?
```

## File Properties

### SEP

The platform-specific path separator character used to join file and directory paths.

**Type:** STRIN

```lol
I CAN HAS SEP FROM MODULE?

I HAS A VARIABLE DIR TEH STRIN ITZ "documents"
I HAS A VARIABLE FILENAME TEH STRIN ITZ "readme.txt"
I HAS A VARIABLE FULLPATH TEH STRIN ITZ DIR MOAR SEP MOAR FILENAME
SAYZ WIT FULLPATH BTW Prints: documents/readme.txt (Unix) or documents\\readme.txt (Windows)
```

```lol
I CAN HAS SEP FROM MODULE?

I HAS A VARIABLE PATH TEH STRIN ITZ "home" MOAR SEP MOAR "user" MOAR SEP MOAR "data"
I HAS A VARIABLE DIR TEH CABINET ITZ NEW CABINET WIT PATH
DIR DO CREATE
```

```lol
I CAN HAS SEP FROM MODULE?

I HAS A VARIABLE PARTS TEH BUKKIT ITZ FULLPATH DO SPLIT WIT SEP
I HAS A VARIABLE IDX TEH INTEGR ITZ 0
WHILE IDX SMALLR THAN PARTS SIZ
I HAS A VARIABLE PART TEH STRIN ITZ PARTS DO AT WIT IDX
SAYZ WIT "Path component: "
SAYZ WIT PART
IDX ITZ IDX MOAR 1
KTHX
```

## Miscellaneous

### CABINET Class

A directory on the filesystem that provides operations for working with directories and their contents.
Supports creating, listing, searching, and managing directories and files within them.

**Methods:**

#### CABINET

Initializes a CABINET instance with the specified directory path.
Creates a directory object for performing directory operations.

**Syntax:** `NEW CABINET WIT <path>`
**Parameters:**
- `path` (STRIN): The directory path to associate with this cabinet

**Example: Create cabinet for existing directory**

```lol
I HAS A VARIABLE DIR TEH CABINET ITZ NEW CABINET WIT "/home/user/documents"
```

**Example: Create cabinet for new directory**

```lol
I HAS A VARIABLE NEWDIR TEH CABINET ITZ NEW CABINET WIT "temp_folder"
BTW Directory doesn't need to exist yet, use CREATE method to make it
```

**Note:** Directory path can be relative or absolute

**Note:** Does not validate directory existence - use EXISTS method to check

**Note:** Does not create the directory - use CREATE method for that

#### CREATE

Creates the directory on the filesystem, including any necessary parent directories.
Similar to 'mkdir -p' command - creates the entire path if needed.

**Syntax:** `<cabinet> DO CREATE`
**Example: Create directory**

```lol
I HAS A VARIABLE DIR TEH CABINET ITZ NEW CABINET WIT "new_folder"
DIR DO CREATE
SAYZ WIT "Directory created successfully"
```

**Example: Create nested directory structure**

```lol
I HAS A VARIABLE NESTED TEH CABINET ITZ NEW CABINET WIT "parent/child/grandchild"
NESTED DO CREATE
SAYZ WIT "Nested directory structure created"
```

**Example: Safe directory creation**

```lol
IZ NO SAEM AS (DIR DO EXISTS)?
DIR DO CREATE
SAYZ WIT "Directory created"
NOPE
SAYZ WIT "Directory already exists"
KTHX
```

**Note:** Creates parent directories automatically if they don't exist

**Note:** No error if directory already exists

**Note:** May throw exception if lacking filesystem permissions

**Note:** Sets default directory permissions (usually 755)

#### DELETE

Deletes an empty directory from the filesystem.
Throws exception if directory is not empty - use DELETE_ALL for recursive deletion.

**Syntax:** `<cabinet> DO DELETE`
**Example: Delete empty directory**

```lol
I HAS A VARIABLE DIR TEH CABINET ITZ NEW CABINET WIT "empty_folder"
IZ DIR DO EXISTS?
DIR DO DELETE
SAYZ WIT "Empty directory deleted"
KTHX
```

**Example: Safe deletion with error handling**

```lol
MAYB
DIR DO DELETE
SAYZ WIT "Directory deleted successfully"
OOPSIE ERR
SAYZ WIT "Failed to delete directory: "
SAYZ WIT ERR
SAYZ WIT "(Directory may not be empty)"
KTHX
```

**Note:** Only deletes empty directories - fails if directory contains files or subdirectories

**Note:** Use LIST to check directory contents before deletion

**Note:** Use DELETE_ALL for recursive deletion of non-empty directories

**Note:** No error if directory doesn't exist

#### DELETE_ALL

Removes directory and all its contents recursively.
Deletes all files and subdirectories within the directory, then deletes the directory itself.

**Syntax:** `<cabinet> DO DELETE_ALL`
**Example: Delete directory tree**

```lol
I HAS A VARIABLE TMPDIR TEH CABINET ITZ NEW CABINET WIT "temp_data"
IZ TMPDIR DO EXISTS?
TMPDIR DO DELETE_ALL
SAYZ WIT "Directory and all contents deleted"
KTHX
```

**Example: Cleanup with confirmation**

```lol
I HAS A VARIABLE FILES TEH BUKKIT ITZ TMPDIR DO LIST
SAYZ WIT "About to delete directory with "
SAYZ WIT FILES SIZ
SAYZ WIT " items. Continue? (y/n)"
I HAS A VARIABLE CONFIRM TEH STRIN ITZ GIMME
IZ CONFIRM SAEM AS "y"?
TMPDIR DO DELETE_ALL
SAYZ WIT "Directory tree deleted"
NOPE
SAYZ WIT "Operation cancelled"
KTHX
```

**Note:** DANGEROUS OPERATION - permanently deletes all contents

**Note:** Use with extreme caution - there is no undo

**Note:** Equivalent to 'rm -rf' command

**Note:** No error if directory doesn't exist

#### EXISTS

Checks if the directory exists on the filesystem.
Returns YEZ if directory exists and is actually a directory, NO otherwise.

**Syntax:** `<cabinet> DO EXISTS`
**Example: Check directory before operations**

```lol
I HAS A VARIABLE DIR TEH CABINET ITZ NEW CABINET WIT "my_folder"
IZ DIR DO EXISTS?
SAYZ WIT "Directory exists"
I HAS A VARIABLE FILES TEH BUKKIT ITZ DIR DO LIST
NOPE
SAYZ WIT "Directory doesn't exist"
KTHX
```

**Example: Conditional directory creation**

```lol
IZ NO SAEM AS (DIR DO EXISTS)?
DIR DO CREATE
SAYZ WIT "Directory created"
NOPE
SAYZ WIT "Directory already exists"
KTHX
```

**Note:** Returns NO if path exists but is not a directory (e.g., it's a file)

**Note:** Use this before LIST or other operations to avoid exceptions

#### FIND

Searches for files and directories matching a glob pattern.
Returns a BUKKIT containing names of items that match the pattern.

**Syntax:** `<cabinet> DO FIND WIT <pattern>`
**Parameters:**
- `pattern` (STRIN): Glob pattern to match against filenames

**Example: Find text files**

```lol
I HAS A VARIABLE DIR TEH CABINET ITZ NEW CABINET WIT "documents"
I HAS A VARIABLE TXTFILES TEH BUKKIT ITZ DIR DO FIND WIT "*.txt"
I HAS A VARIABLE IDX TEH INTEGR ITZ 0
WHILE IDX SMALLR THAN TXTFILES SIZ
I HAS A VARIABLE FILENAME TEH STRIN ITZ TXTFILES DO AT WIT IDX
SAYZ WIT "Found text file: "
SAYZ WIT FILENAME
IDX ITZ IDX MOAR 1
KTHX
```

**Example: Find files with specific prefix**

```lol
I HAS A VARIABLE LOGFILES TEH BUKKIT ITZ DIR DO FIND WIT "log_*"
SAYZ WIT "Found "
SAYZ WIT LOGFILES SIZ
SAYZ WIT " log files"
```

**Example: Find files with specific patterns**

```lol
I HAS A VARIABLE IMAGES TEH BUKKIT ITZ DIR DO FIND WIT "*.{jpg,png,gif}"
BTW Note: Some glob implementations may not support brace expansion
```

**Note:** Uses standard glob patterns: * (any chars), ? (single char), [] (char class)

**Note:** Pattern matching is case-sensitive on most systems

**Note:** Returns empty BUKKIT if no matches found

**Note:** Directory must exist before searching

#### LIST

Returns all files and subdirectories in the directory as a BUKKIT of strings.
Each entry is just the filename or directory name, not the full path.

**Syntax:** `<cabinet> DO LIST`
**Example: List all directory contents**

```lol
I HAS A VARIABLE DIR TEH CABINET ITZ NEW CABINET WIT "documents"
I HAS A VARIABLE FILES TEH BUKKIT ITZ DIR DO LIST
SAYZ WIT "Directory contains "
SAYZ WIT FILES SIZ
SAYZ WIT " items:"
I HAS A VARIABLE IDX TEH INTEGR ITZ 0
WHILE IDX SMALLR THAN FILES SIZ
I HAS A VARIABLE FILENAME TEH STRIN ITZ FILES DO AT WIT IDX
SAYZ WIT "  "
SAYZ WIT FILENAME
IDX ITZ IDX MOAR 1
KTHX
```

**Example: Check for specific files**

```lol
FILES ITZ DIR DO LIST
IDX ITZ 0
WHILE IDX SMALLR THAN FILES SIZ
I HAS A VARIABLE FILENAME TEH STRIN ITZ FILES DO AT WIT IDX
IZ FILENAME SAEM AS "config.txt"?
SAYZ WIT "Found config file!"
KTHX
IDX ITZ IDX MOAR 1
KTHX
```

**Note:** Returns only names, not full paths - combine with PATH property for full paths

**Note:** Directory must exist before listing - check with EXISTS first

**Note:** Returns empty BUKKIT for empty directories

**Note:** Order of entries is not guaranteed

**Member Variables:**

#### PATH

Read-only property containing the directory path.


**Example: Access directory path**

```lol
I HAS A VARIABLE DIR TEH CABINET ITZ NEW CABINET WIT "/home/user/documents"
SAYZ WIT DIR PATH BTW Prints: /home/user/documents
```

**Example: Build full file paths**

```lol
I HAS A VARIABLE FILES TEH BUKKIT ITZ DIR DO LIST
I HAS A VARIABLE IDX TEH INTEGR ITZ 0
WHILE IDX SMALLR THAN FILES SIZ
I HAS A VARIABLE FILENAME TEH STRIN ITZ FILES DO AT WIT IDX
I HAS A VARIABLE FULLPATH TEH STRIN ITZ DIR PATH
FULLPATH ITZ FULLPATH MOAR SEP MOAR FILENAME
SAYZ WIT "Full path: "
SAYZ WIT FULLPATH
IDX ITZ IDX MOAR 1
KTHX
```

**Note:** This is the original path provided to the constructor

**Note:** Path may be relative or absolute depending on how cabinet was created

**Note:** Combine with SEP to build full file paths

**Example: Create and work with directory**

```lol
I HAS A VARIABLE DIR TEH CABINET ITZ NEW CABINET WIT "my_folder"
IZ NO SAEM AS (DIR DO EXISTS)?
DIR DO CREATE
SAYZ WIT "Directory created"
KTHX
```

**Example: List directory contents**

```lol
I HAS A VARIABLE FILES TEH BUKKIT ITZ DIR DO LIST
I HAS A VARIABLE IDX TEH INTEGR ITZ 0
WHILE IDX SMALLR THAN FILES SIZ
I HAS A VARIABLE FILENAME TEH STRIN ITZ FILES DO AT WIT IDX
SAYZ WIT FILENAME
IDX ITZ IDX MOAR 1
KTHX
```

**Example: Find specific files**

```lol
I HAS A VARIABLE TXTFILES TEH BUKKIT ITZ DIR DO FIND WIT "*.txt"
SAYZ WIT "Found "
SAYZ WIT TXTFILES SIZ
SAYZ WIT " text files"
```

### DOCUMENT Class

A file on the file system that provides methods for file I/O operations.
Supports multiple access modes and inherits from IO.READWRITER for compatibility.

**Methods:**

#### CLOSE

Closes the file and releases any resources associated with it.
Automatically flushes any buffered data and sets IS_OPEN to NO.

**Syntax:** `<document> DO CLOSE`
**Example: Basic file operations**

```lol
I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT "data.txt"
DOC DO OPEN WIT "W"
DOC DO WRITE WIT "Some content"
DOC DO CLOSE BTW File is now closed and data is saved
```

**Example: Always close files in exception handling**

```lol
MAYB
DOC DO OPEN WIT "R"
I HAS A VARIABLE DATA TEH STRIN ITZ DOC DO READ WIT 1024
BTW Process data here
OOPSIE ERR
SAYZ WIT ERR
ALWAYZ
IZ DOC IS_OPEN?
DOC DO CLOSE BTW Ensure file is always closed
KTHX
KTHX
```

**Note:** Safe to call multiple times - no error if already closed

**Note:** Automatically flushes buffered data before closing

**Note:** File cannot be used for I/O operations after closing

**Note:** Always close files to prevent resource leaks

#### DELETE

Deletes the file from disk permanently.
Automatically closes the file if open and sets IS_OPEN to NO.

**Syntax:** `<document> DO DELETE`
**Example: Delete a file safely**

```lol
I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT "temp.txt"
IZ DOC DO EXISTS?
DOC DO DELETE
SAYZ WIT "File deleted successfully"
NOPE
SAYZ WIT "File doesn't exist"
KTHX
```

**Example: Cleanup after processing**

```lol
DOC DO OPEN WIT "W"
DOC DO WRITE WIT "Temporary data"
DOC DO CLOSE
BTW Process the file here
DOC DO DELETE BTW Clean up temporary file
```

**Example: Delete with error handling**

```lol
MAYB
DOC DO DELETE
OOPSIE ERR
SAYZ WIT "Failed to delete file: "
SAYZ WIT ERR
KTHX
```

**Note:** File is automatically closed before deletion if it was open

**Note:** Throws exception if file cannot be deleted (permissions, etc.)

**Note:** Operation is irreversible - file cannot be recovered

**Note:** Use EXISTS to check before deletion to avoid exceptions

#### DOCUMENT

Initializes a DOCUMENT instance with a file path.
Creates a new file object but does not open the file yet.

**Syntax:** `NEW DOCUMENT WIT <path>`
**Parameters:**
- `path` (STRIN): The file path to associate with this document

**Example: Create document instance**

```lol
I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT "myfile.txt"
BTW File is not opened yet, use OPEN method to access it
```

**Note:** File path can be relative or absolute

**Note:** Does not validate file existence - use EXISTS method to check

#### EXISTS

Checks if the file exists on disk.
Returns YEZ if file exists, NO otherwise.

**Syntax:** `<document> DO EXISTS`
**Example: Check file existence before reading**

```lol
I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT "config.txt"
IZ DOC DO EXISTS?
DOC DO OPEN WIT "R"
BTW Safe to read the file
NOPE
SAYZ WIT "File not found!"
KTHX
```

**Example: Conditional file creation**

```lol
IZ NO SAEM AS (DOC DO EXISTS)?
DOC DO OPEN WIT "W"
DOC DO WRITE WIT "Default content"
DOC DO CLOSE
SAYZ WIT "Created new file with defaults"
KTHX
```

**Note:** Does not require file to be open

**Note:** Works with any file path, not just open files

**Note:** Use before opening files in read mode to avoid exceptions

#### FLUSH

Flushes the file's buffered contents to disk.
Ensures all pending writes are physically written to storage.

**Syntax:** `<document> DO FLUSH`
**Example: Ensure data is saved**

```lol
I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT "important.txt"
DOC DO OPEN WIT "W"
DOC DO WRITE WIT "Critical data"
DOC DO FLUSH BTW Force write to disk immediately
SAYZ WIT "Data is now safely on disk"
DOC DO CLOSE
```

**Example: Frequent flush for logging**

```lol
DOC DO OPEN WIT "A"
DOC DO WRITE WIT "Log entry 1\n"
DOC DO FLUSH BTW Ensure log is written
DOC DO WRITE WIT "Log entry 2\n"
DOC DO FLUSH BTW Flush again
DOC DO CLOSE
```

**Note:** File must be open before flushing

**Note:** Useful for ensuring data persistence in case of crashes

**Note:** May impact performance if called frequently

**Note:** CLOSE automatically flushes, so explicit FLUSH is optional before closing

#### OPEN

Opens the file for I/O operations according to the specified mode.
Creates the file if it doesn't exist for write/append modes and sets IS_OPEN to YEZ.

**Syntax:** `<document> DO OPEN WIT <mode>`
**Parameters:**
- `mode` (STRIN): Access mode: R (read), W (write/overwrite), RW (read-write), A (append)

**Example: Open file for reading**

```lol
I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT "data.txt"
DOC DO OPEN WIT "R"
BTW File is now open for reading
```

**Example: Open file for writing (creates if needed)**

```lol
I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT "output.txt"
DOC DO OPEN WIT "W"
BTW File is created and open for writing, existing content is overwritten
```

**Example: Open file for append**

```lol
DOC DO OPEN WIT "A"
BTW File is open for appending, writes go to end of file
```

**Note:** R mode requires file to exist, throws exception if not found

**Note:** W and RW modes create the file if it doesn't exist and overwrite existing content

**Note:** A mode creates the file if it doesn't exist and preserves existing content

**Note:** File must be closed before opening with a different mode

#### READ

Reads up to the specified number of characters from the file.
Returns the data read as a string (may be shorter than requested at end of file).

**Syntax:** `<document> DO READ WIT <size>`
**Parameters:**
- `size` (INTEGR): Maximum number of characters to read

**Example: Read entire small file**

```lol
I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT "config.txt"
DOC DO OPEN WIT "R"
I HAS A VARIABLE CONTENT TEH STRIN ITZ DOC DO READ WIT 1024
DOC DO CLOSE
SAYZ WIT CONTENT
```

**Example: Read file in chunks**

```lol
DOC DO OPEN WIT "R"
I HAS A VARIABLE CHUNK TEH STRIN ITZ DOC DO READ WIT 256
WHILE NO SAEM AS (CHUNK SAEM AS "")
SAYZ WIT CHUNK
CHUNK ITZ DOC DO READ WIT 256
KTHX
DOC DO CLOSE
```

**Note:** File must be opened in R or RW mode before reading

**Note:** Returns empty string when end of file is reached

**Note:** Use file position methods (SEEK, TELL) for random access

#### SEEK

Sets the file position for the next read/write operation.
Position is specified as a byte offset from the start of the file (0-based).

**Syntax:** `<document> DO SEEK WIT <position>`
**Parameters:**
- `position` (INTEGR): Byte offset from start of file (0 = beginning)

**Example: Seek to beginning of file**

```lol
I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT "data.txt"
DOC DO OPEN WIT "RW"
DOC DO SEEK WIT 0 BTW Go to start of file
```

**Example: Seek to specific position**

```lol
DOC DO SEEK WIT 100 BTW Go to byte 100
I HAS A VARIABLE DATA TEH STRIN ITZ DOC DO READ WIT 50
```

**Example: Random access file operations**

```lol
DOC DO SEEK WIT 0
DOC DO WRITE WIT "Header"
DOC DO SEEK WIT 50
DOC DO WRITE WIT "Middle"
DOC DO CLOSE
```

**Note:** File must be open before seeking

**Note:** Position beyond end of file is allowed but behavior is undefined

**Note:** Use TELL to get current position

#### TELL

Gets the current file position as a byte offset from the start of the file.
Returns the current position where the next read or write will occur.

**Syntax:** `<document> DO TELL`
**Example: Check current position**

```lol
I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT "data.txt"
DOC DO OPEN WIT "R"
I HAS A VARIABLE POS TEH INTEGR ITZ DOC DO TELL
SAYZ WIT POS BTW Should be 0 at start
```

**Example: Track position while reading**

```lol
DOC DO READ WIT 10 BTW Read 10 characters
POS ITZ DOC DO TELL
SAYZ WIT POS BTW Should be 10
```

**Example: Save and restore position**

```lol
I HAS A VARIABLE SAVED_POS TEH INTEGR ITZ DOC DO TELL
DOC DO READ WIT 100 BTW Read ahead
DOC DO SEEK WIT SAVED_POS BTW Return to saved position
DOC DO CLOSE
```

**Note:** File must be open before getting position

**Note:** Position is measured in bytes, not characters

**Note:** Useful for implementing random access patterns

#### WRITE

Writes string data to the file at the current position.
Returns the number of characters actually written.

**Syntax:** `<document> DO WRITE WIT <data>`
**Parameters:**
- `data` (STRIN): The string data to write to the file

**Example: Write text to file**

```lol
I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT "output.txt"
DOC DO OPEN WIT "W"
I HAS A VARIABLE BYTES_WRITTEN TEH INTEGR ITZ DOC DO WRITE WIT "Hello, World!"
DOC DO CLOSE
SAYZ WIT BYTES_WRITTEN BTW Should print 13
```

**Example: Append to existing file**

```lol
DOC DO OPEN WIT "A"
DOC DO WRITE WIT "\nNew line added"
DOC DO CLOSE
```

**Example: Write multiple lines**

```lol
DOC DO OPEN WIT "W"
DOC DO WRITE WIT "Line 1\n"
DOC DO WRITE WIT "Line 2\n"
DOC DO WRITE WIT "Line 3"
DOC DO CLOSE
```

**Note:** File must be opened in W, RW, or A mode before writing

**Note:** In A mode, data is always written to end of file

**Note:** Use FLUSH to ensure data is written to disk immediately

**Member Variables:**

#### IS_OPEN

Read-only property that indicates whether the file is currently open for I/O operations.


**Example: Check if file is open before operations**

```lol
I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT "data.txt"
IZ DOC IS_OPEN?
SAYZ WIT "File is already open"
NOPE
DOC DO OPEN WIT "R"
SAYZ WIT "File is now open"
KTHX
```

**Example: Safe file operations**

```lol
IZ DOC IS_OPEN?
I HAS A VARIABLE CONTENT TEH STRIN ITZ DOC DO READ WIT 100
NOPE
SAYZ WIT "Cannot read - file is not open"
KTHX
```

**Note:** Returns NO for newly created documents (before OPEN is called)

**Note:** Automatically set to YEZ by OPEN and NO by CLOSE or DELETE

**Note:** Use this property to avoid exceptions from I/O operations on closed files

#### MODE

Read-only property containing the current access mode.


**Example: Check file mode**

```lol
I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT "data.txt"
DOC DO OPEN WIT "RW"
SAYZ WIT DOC MODE BTW Prints: RW
```

**Example: Conditional operations based on mode**

```lol
IZ (DOC MODE) SAEM AS "R"?
SAYZ WIT "File is read-only"
NOPE IZ (DOC MODE) SAEM AS "W"?
SAYZ WIT "File is write-only"
NOPE
SAYZ WIT "File allows both read and write"
KTHX
```

**Note:** Returns empty string if file has never been opened

**Note:** Mode is set when OPEN is called and cleared when CLOSE is called

#### PATH

Read-only property containing the file path.


**Example: Access file path**

```lol
I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT "/home/user/data.txt"
SAYZ WIT DOC PATH BTW Prints: /home/user/data.txt
```

**Note:** This is the original path provided to the constructor

**Note:** Path may be relative or absolute depending on how document was created

#### RWX

File permissions property that controls read/write/execute access.
Can be read to get current permissions or assigned to change them.


**Example: Check file permissions**

```lol
I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT "script.sh"
I HAS A VARIABLE PERMS TEH INTEGR ITZ DOC RWX
SAYZ WIT "File permissions: "
SAYZ WIT PERMS
```

**Example: Make file executable**

```lol
DOC RWX ITZ 0o755 BTW rwxr-xr-x
SAYZ WIT "File is now executable"
```

**Example: Set read-only permissions**

```lol
DOC RWX ITZ 0o444 BTW r--r--r--
SAYZ WIT "File is now read-only"
```

**Example: Common permission values**

```lol
BTW 644 = rw-r--r-- (owner read/write, others read)
BTW 755 = rwxr-xr-x (owner all, others read/execute)
BTW 600 = rw------- (owner read/write only)
```

**Note:** Uses Unix-style octal permission notation

**Note:** Changes take effect immediately on the filesystem

**Note:** May throw exception if user lacks permission to change file permissions

#### SIZ

Read-only property that returns the current file size in bytes.


**Example: Check file size**

```lol
I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT "data.txt"
IZ DOC DO EXISTS?
I HAS A VARIABLE FILESIZE TEH INTEGR ITZ DOC SIZ
SAYZ WIT "File size: "
SAYZ WIT FILESIZE
SAYZ WIT " bytes"
KTHX
```

**Example: Size-based file processing**

```lol
IZ (DOC SIZ) BIGGR DAN 1024?
SAYZ WIT "Large file - processing in chunks"
NOPE
SAYZ WIT "Small file - loading entirely"
KTHX
```

**Note:** Returns current size on disk, even if file is not open

**Note:** Size may change if file is modified by other processes

**Note:** Throws exception if file doesn't exist

**Example: Create and write to a file**

```lol
I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT "data.txt"
DOC DO OPEN WIT "W"
DOC DO WRITE WIT "Hello, World!"
DOC DO CLOSE
```

**Example: Read from a file**

```lol
I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT "input.txt"
DOC DO OPEN WIT "R"
I HAS A VARIABLE CONTENT TEH STRIN ITZ DOC DO READ WIT 1024
DOC DO CLOSE
```

**Example: Check file properties**

```lol
I HAS A VARIABLE DOC TEH DOCUMENT ITZ NEW DOCUMENT WIT "test.txt"
IZ DOC DO EXISTS?
SAYZ WIT "File exists!"
KTHX
SAYZ WIT DOC PATH
SAYZ WIT DOC SIZ
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

### READWRITER Class

Interface that combines both READER and WRITER capabilities.
Inherits READ and CLOSE from READER, and WRITE from WRITER.

**Example: Basic read-write operations**

```lol
I HAS A VARIABLE RW TEH READWRITER ITZ GET_READWRITER
RW DO WRITE WIT "Hello"
I HAS A VARIABLE RESPONSE TEH STRIN ITZ RW DO READ WIT 1024
SAYZ WIT "Response: "
SAYZ WIT RESPONSE
RW DO CLOSE
```

**Example: Echo protocol implementation**

```lol
I HAS A VARIABLE CONNECTION TEH READWRITER ITZ GET_NETWORK_CONNECTION
I HAS A VARIABLE RUNNING TEH BOOL ITZ YEZ
WHILE RUNNING
I HAS A VARIABLE INPUT TEH STRIN ITZ CONNECTION DO READ WIT 1024
IZ INPUT LENGTH SAEM AS 0?
RUNNING ITZ NO
NOPE
CONNECTION DO WRITE WIT "Echo: "
CONNECTION DO WRITE WIT INPUT
CONNECTION DO WRITE WIT "\n"
KTHX
KTHX
CONNECTION DO CLOSE
```

**Example: File copy operation**

```lol
I HAS A VARIABLE SOURCE TEH READWRITER ITZ OPEN_FILE_FOR_READWRITE
I HAS A VARIABLE DEST TEH READWRITER ITZ OPEN_DEST_FILE
I HAS A VARIABLE BUFFER TEH STRIN ITZ SOURCE DO READ WIT 4096
WHILE NO SAEM AS (BUFFER LENGTH SAEM AS 0)
DEST DO WRITE WIT BUFFER
BUFFER ITZ SOURCE DO READ WIT 4096
KTHX
SOURCE DO CLOSE
DEST DO CLOSE
SAYZ WIT "File copy completed"
```

**Example: Interactive session**

```lol
I HAS A VARIABLE SESSION TEH READWRITER ITZ START_INTERACTIVE_SESSION
I HAS A VARIABLE COMMANDS TEH BUKKIT ITZ NEW BUKKIT
COMMANDS DO PUSH WIT "HELP"
COMMANDS DO PUSH WIT "STATUS"
COMMANDS DO PUSH WIT "QUIT"
WHILE NO SAEM AS (COMMANDS LENGTH SAEM AS 0)
I HAS A VARIABLE CMD TEH STRIN ITZ COMMANDS DO POP
SESSION DO WRITE WIT CMD
SESSION DO WRITE WIT "\n"
I HAS A VARIABLE RESPONSE TEH STRIN ITZ SESSION DO READ WIT 2048
SAYZ WIT "Command: "
SAYZ WIT CMD
SAYZ WIT "Response: "
SAYZ WIT RESPONSE
KTHX
SESSION DO CLOSE
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

