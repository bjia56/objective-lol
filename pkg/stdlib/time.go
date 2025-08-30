package stdlib

import (
	"fmt"
	"strings"
	"time"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/types"
)

// Global TIME class and function definitions - created once and reused
var timeClasses = map[string]*environment.Class{
	"DATE": {
		Name: "DATE",
		PublicFunctions: map[string]*environment.Function{
			"DATE": {
				Name:       "DATE",
				Parameters: []environment.Parameter{},
				NativeImpl: func(currentObject *environment.ObjectInstance, args []types.Value) (types.Value, error) {
					currentObject.NativeData = time.Now()
					return types.NOTHIN, nil
				},
			},
			"YEAR": {
				Name:       "YEAR",
				Parameters: []environment.Parameter{},
				NativeImpl: func(currentObject *environment.ObjectInstance, args []types.Value) (types.Value, error) {
					if date, ok := currentObject.NativeData.(time.Time); ok {
						return types.IntegerValue(date.Year()), nil
					}
					return types.NOTHIN, fmt.Errorf("YEAR: invalid context")
				},
			},
			"MONTH": {
				Name:       "MONTH",
				Parameters: []environment.Parameter{},
				NativeImpl: func(currentObject *environment.ObjectInstance, args []types.Value) (types.Value, error) {
					if date, ok := currentObject.NativeData.(time.Time); ok {
						return types.IntegerValue(date.Month()), nil
					}
					return types.NOTHIN, fmt.Errorf("MONTH: invalid context")
				},
			},
			"DAY": {
				Name:       "DAY",
				Parameters: []environment.Parameter{},
				NativeImpl: func(currentObject *environment.ObjectInstance, args []types.Value) (types.Value, error) {
					if date, ok := currentObject.NativeData.(time.Time); ok {
						return types.IntegerValue(date.Day()), nil
					}
					return types.NOTHIN, fmt.Errorf("DAY: invalid context")
				},
			},
			"HOUR": {
				Name:       "HOUR",
				Parameters: []environment.Parameter{},
				NativeImpl: func(currentObject *environment.ObjectInstance, args []types.Value) (types.Value, error) {
					if date, ok := currentObject.NativeData.(time.Time); ok {
						return types.IntegerValue(date.Hour()), nil
					}
					return types.NOTHIN, fmt.Errorf("HOUR: invalid context")
				},
			},
			"MINUTE": {
				Name:       "MINUTE",
				Parameters: []environment.Parameter{},
				NativeImpl: func(currentObject *environment.ObjectInstance, args []types.Value) (types.Value, error) {
					if date, ok := currentObject.NativeData.(time.Time); ok {
						return types.IntegerValue(date.Minute()), nil
					}
					return types.NOTHIN, fmt.Errorf("MINUTE: invalid context")
				},
			},
			"SECOND": {
				Name:       "SECOND",
				Parameters: []environment.Parameter{},
				NativeImpl: func(currentObject *environment.ObjectInstance, args []types.Value) (types.Value, error) {
					if date, ok := currentObject.NativeData.(time.Time); ok {
						return types.IntegerValue(date.Second()), nil
					}
					return types.NOTHIN, fmt.Errorf("SECOND: invalid context")
				},
			},
			"MILLISECOND": {
				Name:       "MILLISECOND",
				Parameters: []environment.Parameter{},
				NativeImpl: func(currentObject *environment.ObjectInstance, args []types.Value) (types.Value, error) {
					if date, ok := currentObject.NativeData.(time.Time); ok {
						return types.IntegerValue(date.Nanosecond() / 1e6), nil
					}
					return types.NOTHIN, fmt.Errorf("MILLISECOND: invalid context")
				},
			},
			"NANOSECOND": {
				Name:       "NANOSECOND",
				Parameters: []environment.Parameter{},
				NativeImpl: func(currentObject *environment.ObjectInstance, args []types.Value) (types.Value, error) {
					if date, ok := currentObject.NativeData.(time.Time); ok {
						return types.IntegerValue(date.Nanosecond()), nil
					}
					return types.NOTHIN, fmt.Errorf("NANOSECOND: invalid context")
				},
			},
			"FORMAT": {
				Name: "FORMAT",
				Parameters: []environment.Parameter{
					{Name: "layout", Type: "STRIN"},
				},
				NativeImpl: func(currentObject *environment.ObjectInstance, args []types.Value) (types.Value, error) {
					layout := args[0]
					if stringVal, ok := layout.(types.StringValue); ok {
						if date, ok := currentObject.NativeData.(time.Time); ok {
							return types.StringValue(date.Format(string(stringVal))), nil
						}
					}
					return types.NOTHIN, fmt.Errorf("FORMAT: invalid context")
				},
			},
		},
	},
}

var timeFunctions = map[string]*environment.Function{
	"SLEEP": {
		Name: "SLEEP",
		Parameters: []environment.Parameter{
			{Name: "seconds", Type: "INTEGR"},
		},
		NativeImpl: func(_ *environment.ObjectInstance, args []types.Value) (types.Value, error) {
			seconds := args[0]

			if secondsVal, ok := seconds.(types.IntegerValue); ok {
				time.Sleep(time.Duration(secondsVal) * time.Second)
				return types.NOTHIN, nil
			}

			return types.NOTHIN, fmt.Errorf("SLEEP: invalid argument type")
		},
	},
}

// RegisterTIMEInEnv registers TIME classes and functions in the given environment
// declarations: empty slice means import all, otherwise import only specified declarations
func RegisterTIMEInEnv(env *environment.Environment, declarations []string) error {
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
			return fmt.Errorf("unknown TIME declaration: %s", decl)
		}
	}

	return nil
}
