# TIME Module

## Import

```lol
BTW Full import
I CAN HAS TIME?

BTW Selective import examples
I CAN HAS SLEEP FROM TIME?
```

### SLEEP

Pauses execution for the specified number of seconds.
Blocks the current thread until the sleep duration expires.

**Returns:** 

### DATE Class

Represents a date and time with methods for accessing components.
Provides year, month, day, hour, minute, second, and formatting capabilities.

**Methods:**

#### DATE

Initializes a DATE object with the current date and time.

#### DAY

Returns the day of the month component (1-31).
Range depends on the specific month and year.

#### FORMAT

Formats the date according to the specified layout string.
Uses Go's time formatting with reference time 'Mon Jan 2 15:04:05 MST 2006'.

#### HOUR

Returns the hour component in 24-hour format (0-23).
0 = midnight, 12 = noon, 23 = 11 PM.

#### MILLISECOND

Returns the millisecond component (0-999).
Milliseconds within the current second.

#### MINUTE

Returns the minute component (0-59).
Minutes past the hour.

#### MONTH

Returns the month component of the date (1-12).
January = 1, February = 2, ..., December = 12.

#### NANOSECOND

Returns the nanosecond component (0-999999999).
Nanoseconds within the current second for high precision timing.

#### SECOND

Returns the second component (0-59).
Seconds past the minute.

#### YEAR

Returns the year component of the date.

