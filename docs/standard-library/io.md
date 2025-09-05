# IO Module - Advanced Input/Output Classes

The IO module provides advanced input/output functionality through abstract base classes and buffered implementations for efficient data processing.

## Importing IO Module

```lol
BTW Import entire module
I CAN HAS IO?

BTW Selective imports
I CAN HAS READER FROM IO?
I CAN HAS WRITER FROM IO?
I CAN HAS BUFFERED_READER FROM IO?
I CAN HAS BUFFERED_WRITER FROM IO?
I CAN HAS READWRITER FROM IO?
```

## Abstract Base Classes

The IO module provides abstract base classes that define the interface for reading and writing operations.

### READER Class

The READER class is an abstract base class for objects that can read data.

#### Methods

- **READ WIT size**: Read up to `size` characters
- **CLOSE**: Close the reader

```lol
BTW READER is abstract - you cannot instantiate it directly
BTW I HAS A VARIABLE R TEH READER ITZ NEW READER  BTW This would fail

BTW Instead, use it as a parent class or use concrete implementations
```

### WRITER Class

The WRITER class is an abstract base class for objects that can write data.

#### Methods

- **WRITE WIT data**: Write string data
- **CLOSE**: Close the writer

```lol
BTW WRITER is abstract - you cannot instantiate it directly
BTW I HAS A VARIABLE W TEH WRITER ITZ NEW WRITER  BTW This would fail

BTW Instead, use it as a parent class or use concrete implementations
```

### READWRITER Class

The READWRITER class combines READER and WRITER interfaces, providing both reading and writing capabilities.

```lol
BTW READWRITER inherits from both READER and WRITER
BTW It provides all methods from both parent classes
```

## Buffered I/O Classes

The IO module provides buffered implementations that improve performance by reducing the number of actual I/O operations.

### BUFFERED_READER Class

The BUFFERED_READER class wraps another READER object and provides buffering for improved read performance.

#### Constructor

```lol
I CAN HAS IO?

BTW Wrap any object that has READ and CLOSE methods
I HAS A VARIABLE BUFFERED TEH BUFFERED_READER ITZ NEW BUFFERED_READER WIT some_reader
```

#### Properties

- **SIZ**: INTEGR - Current buffer size (default: 1024)

#### Methods

- **SET_SIZ WIT newSize**: Change the buffer size
- **READ WIT size**: Read up to `size` characters (buffered)
- **CLOSE**: Close the buffered reader and underlying reader

#### Example

```lol
I CAN HAS IO?
I CAN HAS STDIO?

BTW Create a mock reader class for demonstration
HAI ME TEH CLAS MOCK_READER
    EVRYONE
    DIS TEH VARIABLE DATA TEH STRIN ITZ "This is test data for buffered reading example. It contains multiple sentences to demonstrate buffering behavior."
    DIS TEH VARIABLE POSITION TEH INTEGR ITZ 0

    DIS TEH FUNCSHUN READ TEH STRIN WIT SIZE TEH INTEGR
        IZ POSITION BIGGR THAN LEN WIT DATA LES 1?
            GIVEZ ""  BTW EOF reached
        KTHX

        I HAS A VARIABLE REMAINING TEH INTEGR ITZ LEN WIT DATA LES POSITION
        I HAS A VARIABLE TO_READ TEH INTEGR ITZ SIZE

        IZ TO_READ BIGGR THAN REMAINING?
            TO_READ ITZ REMAINING
        KTHX

        BTW Simulate reading data (simplified substring operation)
        I HAS A VARIABLE RESULT TEH STRIN ITZ "mock data"  BTW In real implementation, would extract substring
        POSITION ITZ POSITION MOAR TO_READ
        GIVEZ RESULT
    KTHX

    DIS TEH FUNCSHUN CLOSE
        SAYZ WIT "Mock reader closed"
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN DEMO_BUFFERED_READER
    SAYZ WIT "=== BUFFERED_READER Demo ==="

    BTW Create a mock reader
    I HAS A VARIABLE MOCK TEH MOCK_READER ITZ NEW MOCK_READER

    BTW Wrap it in a buffered reader
    I HAS A VARIABLE BUFFERED TEH BUFFERED_READER ITZ NEW BUFFERED_READER WIT MOCK

    SAY WIT "Default buffer size: "
    SAYZ WIT BUFFERED SIZ

    BTW Change buffer size
    BUFFERED DO SET_SIZ WIT 512
    SAY WIT "New buffer size: "
    SAYZ WIT BUFFERED SIZ

    BTW Read data in chunks
    I HAS A VARIABLE CHUNK1 TEH STRIN ITZ BUFFERED DO READ WIT 10
    SAY WIT "Read chunk 1: "
    SAYZ WIT CHUNK1

    I HAS A VARIABLE CHUNK2 TEH STRIN ITZ BUFFERED DO READ WIT 20
    SAY WIT "Read chunk 2: "
    SAYZ WIT CHUNK2

    BTW Close the buffered reader
    BUFFERED DO CLOSE
    SAYZ WIT "Buffered reader closed"
KTHXBAI
```

### BUFFERED_WRITER Class

The BUFFERED_WRITER class wraps another WRITER object and provides buffering for improved write performance.

#### Constructor

```lol
I CAN HAS IO?

BTW Wrap any object that has WRITE and CLOSE methods
I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT some_writer
```

#### Properties

- **SIZ**: INTEGR - Current buffer size (default: 1024)

#### Methods

- **SET_SIZ WIT newSize**: Change the buffer size (flushes existing buffer)
- **WRITE WIT data**: Write string data (buffered)
- **FLUSH**: Force write buffered data to underlying writer
- **CLOSE**: Flush and close the buffered writer and underlying writer

#### Example

```lol
I CAN HAS IO?
I CAN HAS STDIO?

BTW Create a mock writer class for demonstration
HAI ME TEH CLAS MOCK_WRITER
    EVRYONE
    DIS TEH VARIABLE WRITTEN_DATA TEH STRIN ITZ ""

    DIS TEH FUNCSHUN WRITE TEH INTEGR WIT DATA TEH STRIN
        WRITTEN_DATA ITZ WRITTEN_DATA MOAR DATA
        SAY WIT "[MOCK WRITE: "
        SAY WIT DATA
        SAYZ WIT "]"
        GIVEZ LEN WIT DATA  BTW Return number of characters written
    KTHX

    DIS TEH FUNCSHUN CLOSE
        SAYZ WIT "Mock writer closed"
    KTHX

    DIS TEH FUNCSHUN GET_WRITTEN_DATA TEH STRIN
        GIVEZ WRITTEN_DATA
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN DEMO_BUFFERED_WRITER
    SAYZ WIT "=== BUFFERED_WRITER Demo ==="

    BTW Create a mock writer
    I HAS A VARIABLE MOCK TEH MOCK_WRITER ITZ NEW MOCK_WRITER

    BTW Wrap it in a buffered writer
    I HAS A VARIABLE BUFFERED TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT MOCK

    SAY WIT "Default buffer size: "
    SAYZ WIT BUFFERED SIZ

    BTW Write small amounts (should be buffered)
    BUFFERED DO WRITE WIT "Hello, "
    BUFFERED DO WRITE WIT "World! "
    BUFFERED DO WRITE WIT "This is buffered writing. "
    SAYZ WIT "Small writes completed (data is buffered)"

    BTW Force flush
    BUFFERED DO FLUSH
    SAYZ WIT "Buffer flushed"

    BTW Write more data
    BUFFERED DO WRITE WIT "More data after flush. "

    BTW Change buffer size (triggers flush)
    BUFFERED DO SET_SIZ WIT 2048
    SAYZ WIT "Buffer size changed (triggers flush)"

    BTW Write final data
    BUFFERED DO WRITE WIT "Final data."

    BTW Close (triggers final flush)
    BUFFERED DO CLOSE
    SAYZ WIT "Buffered writer closed"
KTHXBAI
```

## Practical Examples

### File-like Object Pattern

```lol
I CAN HAS IO?
I CAN HAS STDIO?

BTW Simple in-memory "file" that implements READER/WRITER interface
HAI ME TEH CLAS MEMORY_FILE
    EVRYONE
    DIS TEH VARIABLE CONTENT TEH STRIN ITZ ""
    DIS TEH VARIABLE READ_POS TEH INTEGR ITZ 0

    BTW READER interface
    DIS TEH FUNCSHUN READ TEH STRIN WIT SIZE TEH INTEGR
        IZ READ_POS BIGGR THAN LEN WIT CONTENT LES 1?
            GIVEZ ""  BTW EOF
        KTHX

        I HAS A VARIABLE REMAINING TEH INTEGR ITZ LEN WIT CONTENT LES READ_POS
        I HAS A VARIABLE TO_READ TEH INTEGR ITZ SIZE

        IZ TO_READ BIGGR THAN REMAINING?
            TO_READ ITZ REMAINING
        KTHX

        BTW For simplicity, return a fixed-size chunk
        READ_POS ITZ READ_POS MOAR TO_READ
        GIVEZ "chunk"  BTW In real implementation, would return actual substring
    KTHX

    BTW WRITER interface
    DIS TEH FUNCSHUN WRITE TEH INTEGR WIT DATA TEH STRIN
        CONTENT ITZ CONTENT MOAR DATA
        GIVEZ LEN WIT DATA
    KTHX

    DIS TEH FUNCSHUN CLOSE
        SAYZ WIT "Memory file closed"
    KTHX

    DIS TEH FUNCSHUN RESET
        CONTENT ITZ ""
        READ_POS ITZ 0
    KTHX

    DIS TEH FUNCSHUN GET_CONTENT TEH STRIN
        GIVEZ CONTENT
    KTHX
KTHXBAI

HAI ME TEH FUNCSHUN DEMO_MEMORY_FILE_WITH_BUFFERING
    SAYZ WIT "=== Memory File with Buffering Demo ==="

    BTW Create memory file
    I HAS A VARIABLE MEM_FILE TEH MEMORY_FILE ITZ NEW MEMORY_FILE

    BTW Create buffered writer for efficient writing
    I HAS A VARIABLE BUF_WRITER TEH BUFFERED_WRITER ITZ NEW BUFFERED_WRITER WIT MEM_FILE

    BTW Write data efficiently through buffer
    BUF_WRITER DO WRITE WIT "Line 1: Introduction"
    BUF_WRITER DO WRITE WIT "\n"
    BUF_WRITER DO WRITE WIT "Line 2: Main content"
    BUF_WRITER DO WRITE WIT "\n"
    BUF_WRITER DO WRITE WIT "Line 3: Conclusion"

    BTW Flush and close
    BUF_WRITER DO CLOSE

    BTW Show what was written
    SAY WIT "Content written: "
    SAYZ WIT MEM_FILE DO GET_CONTENT

    BTW Now read it back with buffered reader
    MEM_FILE DO RESET  BTW Reset read position
    I HAS A VARIABLE BUF_READER TEH BUFFERED_READER ITZ NEW BUFFERED_READER WIT MEM_FILE

    BTW Read data in chunks
    I HAS A VARIABLE CHUNK1 TEH STRIN ITZ BUF_READER DO READ WIT 15
    SAY WIT "Read chunk: "
    SAYZ WIT CHUNK1

    BUF_READER DO CLOSE
KTHXBAI
```

### Stream Processing Pattern

```lol
I CAN HAS IO?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN PROCESS_STREAM WIT READER TEH BUFFERED_READER AN WIT WRITER TEH BUFFERED_WRITER
    SAYZ WIT "Processing stream data..."

    I HAS A VARIABLE TOTAL_PROCESSED TEH INTEGR ITZ 0

    WHILE YEZ  BTW Infinite loop, break on EOF
        I HAS A VARIABLE CHUNK TEH STRIN ITZ READER DO READ WIT 100

        BTW Check for EOF
        IZ LEN WIT CHUNK SAEM AS 0?
            SAYZ WIT "End of stream reached"
            BREAK  BTW Would use proper break syntax in real implementation
        KTHX

        BTW Process chunk (uppercase conversion simulation)
        BTW In real implementation, would do actual text processing
        I HAS A VARIABLE PROCESSED TEH STRIN ITZ "[PROCESSED: " MOAR CHUNK MOAR "]"

        BTW Write processed chunk
        WRITER DO WRITE WIT PROCESSED

        TOTAL_PROCESSED ITZ TOTAL_PROCESSED MOAR LEN WIT CHUNK

        BTW Exit condition for demo
        IZ TOTAL_PROCESSED BIGGR THAN 200?
            SAYZ WIT "Demo limit reached"
            GIVEZ UP  BTW Exit function
        KTHX
    KTHX

    SAY WIT "Processed "
    SAY WIT TOTAL_PROCESSED
    SAYZ WIT " characters"
KTHXBAI
```

## Error Handling

### I/O Exception Handling

```lol
I CAN HAS IO?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN SAFE_IO_OPERATIONS
    SAYZ WIT "=== Safe I/O Operations ==="

    MAYB
        BTW Attempt to create buffered reader with invalid object
        I HAS A VARIABLE BAD_READER TEH BUFFERED_READER ITZ NEW BUFFERED_READER WIT "not_a_reader"
    OOPSIE READER_ERROR
        SAYZ WIT "Reader creation error (expected): "
        SAYZ WIT READER_ERROR
    KTHX

    MAYB
        BTW Attempt to set invalid buffer size
        I HAS A VARIABLE READER TEH BUFFERED_READER ITZ NEW BUFFERED_READER WIT some_reader
        READER DO SET_SIZ WIT -100  BTW Invalid negative size
    OOPSIE SIZE_ERROR
        SAYZ WIT "Buffer size error (expected): "
        SAYZ WIT SIZE_ERROR
    KTHX
KTHXBAI
```

### Resource Cleanup Pattern

```lol
I CAN HAS IO?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN PROPER_RESOURCE_CLEANUP
    SAYZ WIT "=== Resource Cleanup Demo ==="

    I HAS A VARIABLE READER TEH BUFFERED_READER ITZ NOTHIN
    I HAS A VARIABLE WRITER TEH BUFFERED_WRITER ITZ NOTHIN

    MAYB
        BTW Create resources
        READER ITZ NEW BUFFERED_READER WIT some_reader_source
        WRITER ITZ NEW BUFFERED_WRITER WIT some_writer_dest

        BTW Do work with resources
        I HAS A VARIABLE DATA TEH STRIN ITZ READER DO READ WIT 100
        WRITER DO WRITE WIT DATA

        SAYZ WIT "Work completed successfully"

    OOPSIE IO_ERROR
        SAYZ WIT "I/O error occurred: "
        SAYZ WIT IO_ERROR

    ALWAYZ
        BTW Always clean up resources
        IZ READER SAEM AS NOTHIN SAEM AS NO?
            READER DO CLOSE
            SAYZ WIT "Reader closed"
        KTHX

        IZ WRITER SAEM AS NOTHIN SAEM AS NO?
            WRITER DO CLOSE
            SAYZ WIT "Writer closed"
        KTHX

        SAYZ WIT "Resource cleanup completed"
    KTHX
KTHXBAI
```

## Quick Reference

### Classes

| Class | Type | Description |
|-------|------|-------------|
| `READER` | Abstract | Base class for reading operations |
| `WRITER` | Abstract | Base class for writing operations |
| `READWRITER` | Abstract | Combined reader/writer interface |
| `BUFFERED_READER` | Concrete | Buffered reading implementation |
| `BUFFERED_WRITER` | Concrete | Buffered writing implementation |

### BUFFERED_READER Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `BUFFERED_READER WIT reader` | reader: READER | - | Constructor |
| `SET_SIZ WIT size` | size: INTEGR | - | Change buffer size |
| `READ WIT size` | size: INTEGR | STRIN | Read data (buffered) |
| `CLOSE` | - | - | Close reader |

### BUFFERED_WRITER Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `BUFFERED_WRITER WIT writer` | writer: WRITER | - | Constructor |
| `SET_SIZ WIT size` | size: INTEGR | - | Change buffer size |
| `WRITE WIT data` | data: STRIN | INTEGR | Write data (buffered) |
| `FLUSH` | - | - | Force write buffer |
| `CLOSE` | - | - | Flush and close writer |

### Properties

| Property | Type | Description |
|----------|------|-------------|
| `SIZ` | INTEGR | Current buffer size |

## Related

- [STDIO Module](stdio.md) - Basic input/output functions
- [Collections](collections.md) - Working with data structures
- [Control Flow](../language-guide/control-flow.md) - Exception handling patterns