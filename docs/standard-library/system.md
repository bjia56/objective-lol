# SYSTEM Module - Environment Variables

The SYSTEM module provides system-level operations through environment variable management with the ENVBASKIT class.

## Importing SYSTEM Module

```lol
BTW Import entire module
I CAN HAS SYSTEM?

BTW Selective import
I CAN HAS ENVBASKIT FROM SYSTEM?
I CAN HAS ENV FROM SYSTEM?
```

**Note:** The SYSTEM module automatically imports the MAPS module since ENVBASKIT inherits from BASKIT.

## Global ENV Variable

The SYSTEM module provides a global `ENV` variable that is a pre-initialized ENVBASKIT instance containing all current environment variables.

```lol
I CAN HAS SYSTEM?

BTW Use global ENV to access environment variables
I HAS A VARIABLE HOME_PATH TEH STRIN ITZ ENV DO GET WIT "HOME"
SAYZ WIT HOME_PATH

BTW Set new environment variable
ENV DO PUT WIT "MY_APP_CONFIG" AN WIT "/path/to/config"
```

## ENVBASKIT Class

The ENVBASKIT class extends BASKIT (map/dictionary) functionality to provide seamless integration with system environment variables. It automatically syncs with the actual process environment.

### Constructor

The ENVBASKIT constructor is private and cannot be called directly from LOLCODE. Use the global `ENV` variable instead.

### Properties

- **SIZ**: INTEGR (read-only) - Number of environment variables

### Methods

#### Environment Access

##### GET - Get Environment Variable

Gets the value of an environment variable. Checks the internal map first, then the actual environment.

```lol
I HAS A VARIABLE VALUE TEH STRIN ITZ ENV DO GET WIT <key>
```

**Parameters:**
- **key**: STRIN - Environment variable name

**Returns:** STRIN - Environment variable value

**Throws exception if:** Variable does not exist

##### PUT - Set Environment Variable

Sets an environment variable both in the internal map and the actual process environment.

```lol
ENV DO PUT WIT <key> AN WIT <value>
```

**Parameters:**
- **key**: STRIN - Environment variable name
- **value**: STRIN - Environment variable value

##### CONTAINS - Check Variable Existence

Checks if an environment variable exists in either the internal map or actual environment.

```lol
I HAS A VARIABLE EXISTS TEH BOOL ITZ ENV DO CONTAINS WIT <key>
```

**Parameters:**
- **key**: STRIN - Environment variable name

**Returns:** BOOL - YEZ if variable exists, NO otherwise

##### REMOVE - Remove Environment Variable

Removes an environment variable from both the internal map and actual process environment.

```lol
I HAS A VARIABLE OLD_VALUE TEH STRIN ITZ ENV DO REMOVE WIT <key>
```

**Parameters:**
- **key**: STRIN - Environment variable name

**Returns:** STRIN - Previous value of the variable

**Throws exception if:** Variable does not exist

#### Synchronization

##### REFRESH - Sync with Environment

Refreshes the internal map with all current environment variables, discarding any previous state.

```lol
ENV DO REFRESH
```

##### CLEAR - Clear All Variables

Clears all environment variables that are tracked in the internal map and unsets them from the actual environment.

```lol
ENV DO CLEAR
```

**Warning:** This only removes variables that were previously accessed or set through this ENVBASKIT instance.

#### Collection Operations

##### KEYS - Get Variable Names

Returns a BUKKIT containing all environment variable names.

```lol
I HAS A VARIABLE KEYS TEH BUKKIT ITZ ENV DO KEYS
```

**Returns:** BUKKIT - Array of environment variable names

##### VALUES - Get Variable Values

Returns a BUKKIT containing all environment variable values.

```lol
I HAS A VARIABLE VALUES TEH BUKKIT ITZ ENV DO VALUES
```

**Returns:** BUKKIT - Array of environment variable values

##### PAIRS - Get Key-Value Pairs

Returns a BUKKIT containing BASKIT objects representing key-value pairs.

```lol
I HAS A VARIABLE PAIRS TEH BUKKIT ITZ ENV DO PAIRS
```

**Returns:** BUKKIT - Array of BASKIT objects with "KEY" and "VALUE" entries

##### COPY - Create Copy

Creates a new ENVBASKIT instance with the same data as the current instance.

```lol
I HAS A VARIABLE ENV_COPY TEH ENVBASKIT ITZ ENV DO COPY
```

**Returns:** ENVBASKIT - New instance with copied data

**Note:** The copy is independent and does not sync changes back to the original.

##### MERGE - Merge with Another Map

Merges environment variables from another BASKIT or ENVBASKIT.

```lol
ENV DO MERGE WIT <other_baskit>
```

**Parameters:**
- **other_baskit**: BASKIT or ENVBASKIT - Source to merge from

## Basic Environment Operations

### Reading Environment Variables

```lol
I CAN HAS SYSTEM?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN READ_ENV_VARS
    SAYZ WIT "=== Reading Environment Variables ==="

    BTW Read common environment variables
    MAYB
        I HAS A VARIABLE HOME TEH STRIN ITZ ENV DO GET WIT "HOME"
        SAY WIT "Home directory: "
        SAYZ WIT HOME
    OOPSIE HOME_ERROR
        SAYZ WIT "HOME variable not found"
    KTHX

    MAYB
        I HAS A VARIABLE USER TEH STRIN ITZ ENV DO GET WIT "USER"
        SAY WIT "User: "
        SAYZ WIT USER
    OOPSIE USER_ERROR
        SAYZ WIT "USER variable not found"
    KTHX

    BTW Check if variable exists before reading
    IZ ENV DO CONTAINS WIT "PATH"?
        I HAS A VARIABLE PATH TEH STRIN ITZ ENV DO GET WIT "PATH"
        SAYZ WIT "PATH is set"
    NOPE
        SAYZ WIT "PATH variable not found"
    KTHX
KTHXBAI
```

### Setting Environment Variables

```lol
I CAN HAS SYSTEM?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN SET_ENV_VARS
    SAYZ WIT "=== Setting Environment Variables ==="

    BTW Set application configuration
    ENV DO PUT WIT "MY_APP_NAME" AN WIT "ObjectiveLOL App"
    ENV DO PUT WIT "MY_APP_VERSION" AN WIT "1.0.0"
    ENV DO PUT WIT "MY_APP_DEBUG" AN WIT "true"

    SAYZ WIT "Environment variables set"

    BTW Verify they were set
    I HAS A VARIABLE APP_NAME TEH STRIN ITZ ENV DO GET WIT "MY_APP_NAME"
    SAY WIT "App name: "
    SAYZ WIT APP_NAME

    I HAS A VARIABLE DEBUG_MODE TEH STRIN ITZ ENV DO GET WIT "MY_APP_DEBUG"
    SAY WIT "Debug mode: "
    SAYZ WIT DEBUG_MODE
KTHXBAI
```

## Quick Reference

### Global Variables

| Variable | Type | Description |
|----------|------|-------------|
| `ENV` | ENVBASKIT | Global environment variable manager |

### ENVBASKIT Methods

| Method | Parameters | Returns | Description |
|--------|------------|---------|-------------|
| `GET WIT key` | key: STRIN | STRIN | Get environment variable |
| `PUT WIT key AN WIT value` | key: STRIN, value: STRIN | - | Set environment variable |
| `CONTAINS WIT key` | key: STRIN | BOOL | Check if variable exists |
| `REMOVE WIT key` | key: STRIN | STRIN | Remove environment variable |
| `REFRESH` | - | - | Sync with current environment |
| `CLEAR` | - | - | Clear tracked variables |
| `KEYS` | - | BUKKIT | Get all variable names |
| `VALUES` | - | BUKKIT | Get all variable values |
| `PAIRS` | - | BUKKIT | Get key-value pairs |
| `COPY` | - | ENVBASKIT | Create independent copy |
| `MERGE WIT other` | other: BASKIT | - | Merge from another map |

### Properties

| Property | Type | Description |
|----------|------|-------------|
| `SIZ` | INTEGR | Number of environment variables (read-only) |

## Related

- [MAPS Module](collections.md) - BASKIT parent class
- [STDIO Module](stdio.md) - Console input/output
- [Control Flow](../language-guide/control-flow.md) - Exception handling patterns