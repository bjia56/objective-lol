# TIME Module - Date and Time Functions

The TIME module provides date and time functionality through the DATE class and utility functions for time-based operations.

## Importing TIME Module

```lol
BTW Import entire module
I CAN HAS TIME?

BTW Selective imports
I CAN HAS DATE FROM TIME?
I CAN HAS SLEEP FROM TIME?
```

## DATE Class

The DATE class represents a moment in time and provides methods to access and format date/time components.

### Creating DATE Objects

```lol
I CAN HAS TIME?

BTW Create DATE object with current time
I HAS A VARIABLE NOW TEH DATE ITZ NEW DATE
```

### DATE Methods

#### Time Component Methods

All time component methods return INTEGR values:

```lol
I CAN HAS TIME?

I HAS A VARIABLE NOW TEH DATE ITZ NEW DATE

BTW Get year (e.g., 2024)
I HAS A VARIABLE CURRENT_YEAR TEH INTEGR ITZ NOW DO YEAR
SAYZ WIT CURRENT_YEAR

BTW Get month (1-12)
I HAS A VARIABLE CURRENT_MONTH TEH INTEGR ITZ NOW DO MONTH
SAYZ WIT CURRENT_MONTH

BTW Get day of month (1-31)
I HAS A VARIABLE CURRENT_DAY TEH INTEGR ITZ NOW DO DAY
SAYZ WIT CURRENT_DAY

BTW Get hour (0-23)
I HAS A VARIABLE CURRENT_HOUR TEH INTEGR ITZ NOW DO HOUR
SAYZ WIT CURRENT_HOUR

BTW Get minute (0-59)
I HAS A VARIABLE CURRENT_MINUTE TEH INTEGR ITZ NOW DO MINUTE
SAYZ WIT CURRENT_MINUTE

BTW Get second (0-59)
I HAS A VARIABLE CURRENT_SECOND TEH INTEGR ITZ NOW DO SECOND
SAYZ WIT CURRENT_SECOND

BTW Get millisecond (0-999)
I HAS A VARIABLE CURRENT_MS TEH INTEGR ITZ NOW DO MILLISECOND
SAYZ WIT CURRENT_MS

BTW Get nanosecond (0-999999999)
I HAS A VARIABLE CURRENT_NS TEH INTEGR ITZ NOW DO NANOSECOND
SAYZ WIT CURRENT_NS
```

#### Date Formatting

The FORMAT method allows custom date formatting using Go's time layout format:

```lol
I CAN HAS TIME?
I CAN HAS STDIO?

I HAS A VARIABLE NOW TEH DATE ITZ NEW DATE

BTW Common formats
I HAS A VARIABLE ISO_FORMAT TEH STRIN ITZ NOW DO FORMAT WIT "2006-01-02T15:04:05"
SAYZ WIT ISO_FORMAT

I HAS A VARIABLE READABLE_FORMAT TEH STRIN ITZ NOW DO FORMAT WIT "January 2, 2006 at 3:04 PM"
SAYZ WIT READABLE_FORMAT

I HAS A VARIABLE CUSTOM_FORMAT TEH STRIN ITZ NOW DO FORMAT WIT "Mon Jan 2 15:04:05 2006"
SAYZ WIT CUSTOM_FORMAT
```

**Common Format Layouts:**

| Layout | Description | Example Output |
|--------|-------------|----------------|
| `"2006-01-02"` | ISO date | `2024-03-15` |
| `"15:04:05"` | 24-hour time | `14:30:45` |
| `"3:04 PM"` | 12-hour time | `2:30 PM` |
| `"January 2, 2006"` | Full date | `March 15, 2024` |
| `"Jan 2, 2006"` | Abbreviated date | `Mar 15, 2024` |
| `"Monday"` | Day of week | `Friday` |
| `"Mon"` | Abbreviated day | `Fri` |

### Complete DATE Example

```lol
I CAN HAS TIME?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN DISPLAY_CURRENT_TIME
    I HAS A VARIABLE NOW TEH DATE ITZ NEW DATE

    BTW Display individual components
    SAY WIT "Current time: "
    SAY WIT NOW DO YEAR
    SAY WIT "-"
    SAY WIT NOW DO MONTH
    SAY WIT "-"
    SAY WIT NOW DO DAY
    SAY WIT " "
    SAY WIT NOW DO HOUR
    SAY WIT ":"
    SAY WIT NOW DO MINUTE
    SAY WIT ":"
    SAYZ WIT NOW DO SECOND

    BTW Display formatted versions
    I HAS A VARIABLE FORMATTED TEH STRIN ITZ NOW DO FORMAT WIT "Monday, January 2, 2006 at 3:04:05 PM"
    SAYZ WIT FORMATTED

    BTW Display timestamp
    I HAS A VARIABLE TIMESTAMP TEH STRIN ITZ NOW DO FORMAT WIT "2006-01-02T15:04:05.000Z"
    SAYZ WIT TIMESTAMP
KTHXBAI
```

## SLEEP Function

The SLEEP function pauses program execution for the specified number of seconds.

### Syntax

```lol
SLEEP WIT <seconds>
```

- **seconds**: INTEGR - Number of seconds to sleep

### SLEEP Examples

```lol
I CAN HAS TIME?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN COUNTDOWN
    I HAS A VARIABLE COUNT TEH INTEGR ITZ 5

    WHILE COUNT BIGGR THAN 0
        SAYZ WIT COUNT
        SLEEP WIT 1         BTW Wait 1 second
        COUNT ITZ COUNT LES 1
    KTHX

    SAYZ WIT "Time's up!"
KTHXBAI

HAI ME TEH FUNCSHUN PROGRESS_SIMULATION
    I HAS A VARIABLE I TEH INTEGR ITZ 1

    WHILE I SMALLR THAN 6
        SAY WIT "Step "
        SAY WIT I
        SAYZ WIT " processing..."
        SLEEP WIT 2         BTW Wait 2 seconds between steps
        I ITZ I MOAR 1
    KTHX

    SAYZ WIT "Processing complete!"
KTHXBAI
```

## Practical Examples

### Digital Clock Display

```lol
I CAN HAS TIME?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN DIGITAL_CLOCK
    I HAS A VARIABLE COUNTER TEH INTEGR ITZ 10

    WHILE COUNTER BIGGR THAN 0
        I HAS A VARIABLE NOW TEH DATE ITZ NEW DATE
        I HAS A VARIABLE TIME_STR TEH STRIN ITZ NOW DO FORMAT WIT "15:04:05"

        SAY WIT "Current time: "
        SAYZ WIT TIME_STR

        SLEEP WIT 1
        COUNTER ITZ COUNTER LES 1
    KTHX
KTHXBAI
```

### Date Calculator

```lol
I CAN HAS TIME?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN DATE_INFO
    I HAS A VARIABLE NOW TEH DATE ITZ NEW DATE

    I HAS A VARIABLE YEAR TEH INTEGR ITZ NOW DO YEAR
    I HAS A VARIABLE MONTH TEH INTEGR ITZ NOW DO MONTH
    I HAS A VARIABLE DAY TEH INTEGR ITZ NOW DO DAY

    SAY WIT "Today is "
    I HAS A VARIABLE DATE_STR TEH STRIN ITZ NOW DO FORMAT WIT "Monday, January 2, 2006"
    SAYZ WIT DATE_STR

    BTW Check if it's a weekend (simplified - just checks day name)
    I HAS A VARIABLE DAY_NAME TEH STRIN ITZ NOW DO FORMAT WIT "Monday"
    IZ DAY_NAME SAEM AS "Saturday"?
        SAYZ WIT "It's the weekend!"
    NOPE
        IZ DAY_NAME SAEM AS "Sunday"?
            SAYZ WIT "It's the weekend!"
        NOPE
            SAYZ WIT "It's a weekday."
        KTHX
    KTHX
KTHXBAI
```

### Timer Application

```lol
I CAN HAS TIME?
I CAN HAS STDIO?

HAI ME TEH FUNCSHUN TIMER WIT DURATION TEH INTEGR
    SAYZ WIT "Timer started!"

    I HAS A VARIABLE START_TIME TEH DATE ITZ NEW DATE
    I HAS A VARIABLE START_STR TEH STRIN ITZ START_TIME DO FORMAT WIT "15:04:05"
    SAY WIT "Start time: "
    SAYZ WIT START_STR

    SLEEP WIT DURATION

    I HAS A VARIABLE END_TIME TEH DATE ITZ NEW DATE
    I HAS A VARIABLE END_STR TEH STRIN ITZ END_TIME DO FORMAT WIT "15:04:05"
    SAY WIT "End time: "
    SAYZ WIT END_STR

    SAYZ WIT "Timer finished!"
KTHXBAI

HAI ME TEH FUNCSHUN POMODORO_TIMER
    SAYZ WIT "Starting 25-minute Pomodoro session"

    BTW Work session (25 minutes = 1500 seconds, but using 5 for demo)
    TIMER WIT 5
    SAYZ WIT "Work session complete! Take a break."

    BTW Break session (5 minutes = 300 seconds, but using 2 for demo)
    TIMER WIT 2
    SAYZ WIT "Break complete! Back to work."
KTHXBAI
```


## Quick Reference

### DATE Methods

| Method | Return Type | Description |
|--------|-------------|-------------|
| `YEAR` | INTEGR | Get year (e.g., 2024) |
| `MONTH` | INTEGR | Get month (1-12) |
| `DAY` | INTEGR | Get day of month (1-31) |
| `HOUR` | INTEGR | Get hour (0-23) |
| `MINUTE` | INTEGR | Get minute (0-59) |
| `SECOND` | INTEGR | Get second (0-59) |
| `MILLISECOND` | INTEGR | Get millisecond (0-999) |
| `NANOSECOND` | INTEGR | Get nanosecond (0-999999999) |
| `FORMAT WIT layout` | STRIN | Format date using layout |

### Functions

| Function | Parameters | Description |
|----------|------------|-------------|
| `SLEEP WIT seconds` | seconds: INTEGR | Pause execution |

## Related

- [STDIO Module](stdio.md) - For displaying time information
- [Control Flow](../language-guide/control-flow.md) - Using time in loops and conditions