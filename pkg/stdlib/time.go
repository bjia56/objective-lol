package stdlib

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/runtime"
)

// Global TIME class definitions - created once and reused
var timeClassesOnce = sync.Once{}
var timeClasses map[string]*environment.Class

func getTimeClasses() map[string]*environment.Class {
	timeClassesOnce.Do(func() {
		timeClasses = map[string]*environment.Class{
			"DATE": {
				Name: "DATE",
				Documentation: []string{
					"Represents a date and time with methods for accessing components.",
					"Provides year, month, day, hour, minute, second, and formatting capabilities.",
				},
				PublicFunctions: map[string]*environment.Function{
					"DATE": {
						Name: "DATE",
						Documentation: []string{
							"Initializes a DATE object with the current date and time.",
						},
						Parameters: []environment.Parameter{},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							this.NativeData = time.Now()
							return environment.NOTHIN, nil
						},
					},
					"YEAR": {
						Name: "YEAR",
						Documentation: []string{
							"Returns the year component of the date.",
						},
						ReturnType: "INTEGR",
						Parameters: []environment.Parameter{},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							if date, ok := this.NativeData.(time.Time); ok {
								return environment.IntegerValue(date.Year()), nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "YEAR: invalid context"}
						},
					},
					"MONTH": {
						Name: "MONTH",
						Documentation: []string{
							"Returns the month component of the date (1-12).",
							"January = 1, February = 2, ..., December = 12.",
						},
						ReturnType: "INTEGR",
						Parameters: []environment.Parameter{},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							if date, ok := this.NativeData.(time.Time); ok {
								return environment.IntegerValue(date.Month()), nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "MONTH: invalid context"}
						},
					},
					"DAY": {
						Name: "DAY",
						Documentation: []string{
							"Returns the day of the month component (1-31).",
							"Range depends on the specific month and year.",
						},
						ReturnType: "INTEGR",
						Parameters: []environment.Parameter{},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							if date, ok := this.NativeData.(time.Time); ok {
								return environment.IntegerValue(date.Day()), nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "DAY: invalid context"}
						},
					},
					"HOUR": {
						Name: "HOUR",
						Documentation: []string{
							"Returns the hour component in 24-hour format (0-23).",
							"0 = midnight, 12 = noon, 23 = 11 PM.",
						},
						ReturnType: "INTEGR",
						Parameters: []environment.Parameter{},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							if date, ok := this.NativeData.(time.Time); ok {
								return environment.IntegerValue(date.Hour()), nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "HOUR: invalid context"}
						},
					},
					"MINUTE": {
						Name: "MINUTE",
						Documentation: []string{
							"Returns the minute component (0-59).",
							"Minutes past the hour.",
						},
						ReturnType: "INTEGR",
						Parameters: []environment.Parameter{},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							if date, ok := this.NativeData.(time.Time); ok {
								return environment.IntegerValue(date.Minute()), nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "MINUTE: invalid context"}
						},
					},
					"SECOND": {
						Name: "SECOND",
						Documentation: []string{
							"Returns the second component (0-59).",
							"Seconds past the minute.",
						},
						ReturnType: "INTEGR",
						Parameters: []environment.Parameter{},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							if date, ok := this.NativeData.(time.Time); ok {
								return environment.IntegerValue(date.Second()), nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "SECOND: invalid context"}
						},
					},
					"MILLISECOND": {
						Name: "MILLISECOND",
						Documentation: []string{
							"Returns the millisecond component (0-999).",
							"Milliseconds within the current second.",
						},
						ReturnType: "INTEGR",
						Parameters: []environment.Parameter{},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							if date, ok := this.NativeData.(time.Time); ok {
								return environment.IntegerValue(date.Nanosecond() / 1e6), nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "MILLISECOND: invalid context"}
						},
					},
					"NANOSECOND": {
						Name: "NANOSECOND",
						Documentation: []string{
							"Returns the nanosecond component (0-999999999).",
							"Nanoseconds within the current second for high precision timing.",
						},
						ReturnType: "INTEGR",
						Parameters: []environment.Parameter{},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							if date, ok := this.NativeData.(time.Time); ok {
								return environment.IntegerValue(date.Nanosecond()), nil
							}
							return environment.NOTHIN, runtime.Exception{Message: "NANOSECOND: invalid context"}
						},
					},
					"FORMAT": {
						Name: "FORMAT",
						Documentation: []string{
							"Formats the date according to the specified layout string.",
							"Uses Go's time formatting with reference time 'Mon Jan 2 15:04:05 MST 2006'.",
						},
						ReturnType: "STRIN",
						Parameters: []environment.Parameter{
							{Name: "layout", Type: "STRIN"},
						},
						NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
							layout := args[0]
							if stringVal, ok := layout.(environment.StringValue); ok {
								if date, ok := this.NativeData.(time.Time); ok {
									return environment.StringValue(date.Format(string(stringVal))), nil
								}
							}
							return environment.NOTHIN, runtime.Exception{Message: "FORMAT: invalid context"}
						},
					},
				},
				QualifiedName:    "stdlib:TIME.DATE",
				ModulePath:       "stdlib:TIME",
				ParentClasses:    []string{},
				MRO:              []string{"stdlib:TIME.DATE"},
				PrivateVariables: make(map[string]*environment.MemberVariable),
				PrivateFunctions: make(map[string]*environment.Function),
				SharedVariables:  make(map[string]*environment.MemberVariable),
				SharedFunctions:  make(map[string]*environment.Function),
				PublicVariables:  make(map[string]*environment.MemberVariable),
			},
		}
	})
	return timeClasses
}

// Global TIME function definitions - created once and reused

var timeFunctionsOnce = sync.Once{}
var timeFunctions map[string]*environment.Function

func getTimeFunctions() map[string]*environment.Function {
	timeFunctionsOnce.Do(func() {
		timeFunctions = map[string]*environment.Function{
			"SLEEP": {
				Name: "SLEEP",
				Documentation: []string{
					"Pauses execution for the specified number of seconds.",
					"Blocks the current thread until the sleep duration expires.",
				},
				Parameters: []environment.Parameter{
					{Name: "seconds", Type: "INTEGR"},
				},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					seconds := args[0]

					if secondsVal, ok := seconds.(environment.IntegerValue); ok {
						time.Sleep(time.Duration(secondsVal) * time.Second)
						return environment.NOTHIN, nil
					}

					return environment.NOTHIN, runtime.Exception{Message: "SLEEP: invalid argument type"}
				},
			},
		}
	})
	return timeFunctions
}

// RegisterTIMEInEnv registers TIME classes and functions in the given environment
// declarations: empty slice means import all, otherwise import only specified declarations
func RegisterTIMEInEnv(env *environment.Environment, declarations ...string) error {
	timeClasses := getTimeClasses()
	timeFunctions := getTimeFunctions()

	// If declarations is empty, import all classes and functions
	if len(declarations) == 0 {
		for _, class := range timeClasses {
			env.DefineClass(class)
		}
		for _, fn := range timeFunctions {
			env.DefineFunction(fn)
		}
		return nil
	}

	// Otherwise, import only specified declarations
	for _, decl := range declarations {
		declUpper := strings.ToUpper(decl)

		// Check if it's a class
		if class, exists := timeClasses[declUpper]; exists {
			env.DefineClass(class)
		} else if fn, exists := timeFunctions[declUpper]; exists {
			// Check if it's a function
			env.DefineFunction(fn)
		} else {
			return runtime.Exception{Message: fmt.Sprintf("unknown TIME declaration: %s", decl)}
		}
	}

	return nil
}
