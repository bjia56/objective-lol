package stdlib

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/types"
)

// Global TIME class definitions - created once and reused
var timeClassesOnce = sync.Once{}
var timeClasses map[string]*environment.Class

func getTimeClasses() map[string]*environment.Class {
	timeClassesOnce.Do(func() {
		timeClasses = map[string]*environment.Class{
			"DATE": {
				Name: "DATE",
				PublicFunctions: map[string]*environment.Function{
					"DATE": {
						Name:       "DATE",
						Parameters: []environment.Parameter{},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							this.NativeData = time.Now()
							return types.NOTHIN, nil
						},
					},
					"YEAR": {
						Name:       "YEAR",
						ReturnType: "INTEGR",
						Parameters: []environment.Parameter{},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							if date, ok := this.NativeData.(time.Time); ok {
								return types.IntegerValue(date.Year()), nil
							}
							return types.NOTHIN, fmt.Errorf("YEAR: invalid context")
						},
					},
					"MONTH": {
						Name:       "MONTH",
						ReturnType: "INTEGR",
						Parameters: []environment.Parameter{},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							if date, ok := this.NativeData.(time.Time); ok {
								return types.IntegerValue(date.Month()), nil
							}
							return types.NOTHIN, fmt.Errorf("MONTH: invalid context")
						},
					},
					"DAY": {
						Name:       "DAY",
						ReturnType: "INTEGR",
						Parameters: []environment.Parameter{},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							if date, ok := this.NativeData.(time.Time); ok {
								return types.IntegerValue(date.Day()), nil
							}
							return types.NOTHIN, fmt.Errorf("DAY: invalid context")
						},
					},
					"HOUR": {
						Name:       "HOUR",
						ReturnType: "INTEGR",
						Parameters: []environment.Parameter{},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							if date, ok := this.NativeData.(time.Time); ok {
								return types.IntegerValue(date.Hour()), nil
							}
							return types.NOTHIN, fmt.Errorf("HOUR: invalid context")
						},
					},
					"MINUTE": {
						Name:       "MINUTE",
						ReturnType: "INTEGR",
						Parameters: []environment.Parameter{},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							if date, ok := this.NativeData.(time.Time); ok {
								return types.IntegerValue(date.Minute()), nil
							}
							return types.NOTHIN, fmt.Errorf("MINUTE: invalid context")
						},
					},
					"SECOND": {
						Name:       "SECOND",
						ReturnType: "INTEGR",
						Parameters: []environment.Parameter{},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							if date, ok := this.NativeData.(time.Time); ok {
								return types.IntegerValue(date.Second()), nil
							}
							return types.NOTHIN, fmt.Errorf("SECOND: invalid context")
						},
					},
					"MILLISECOND": {
						Name:       "MILLISECOND",
						ReturnType: "INTEGR",
						Parameters: []environment.Parameter{},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							if date, ok := this.NativeData.(time.Time); ok {
								return types.IntegerValue(date.Nanosecond() / 1e6), nil
							}
							return types.NOTHIN, fmt.Errorf("MILLISECOND: invalid context")
						},
					},
					"NANOSECOND": {
						Name:       "NANOSECOND",
						ReturnType: "INTEGR",
						Parameters: []environment.Parameter{},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							if date, ok := this.NativeData.(time.Time); ok {
								return types.IntegerValue(date.Nanosecond()), nil
							}
							return types.NOTHIN, fmt.Errorf("NANOSECOND: invalid context")
						},
					},
					"FORMAT": {
						Name:       "FORMAT",
						ReturnType: "STRIN",
						Parameters: []environment.Parameter{
							{Name: "layout", Type: "STRIN"},
						},
						NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
							layout := args[0]
							if stringVal, ok := layout.(types.StringValue); ok {
								if date, ok := this.NativeData.(time.Time); ok {
									return types.StringValue(date.Format(string(stringVal))), nil
								}
							}
							return types.NOTHIN, fmt.Errorf("FORMAT: invalid context")
						},
					},
				},
				QualifiedName: "stdlib:TIME.DATE",
				ModulePath:    "stdlib:TIME",
				ParentClass:   "",
				PrivateVariables: make(map[string]*environment.Variable),
				PrivateFunctions: make(map[string]*environment.Function),
				SharedVariables:  make(map[string]*environment.Variable),
				SharedFunctions:  make(map[string]*environment.Function),
				PublicVariables:  make(map[string]*environment.Variable),
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
				Parameters: []environment.Parameter{
					{Name: "seconds", Type: "INTEGR"},
				},
				NativeImpl: func(ctx interface{}, this *environment.ObjectInstance, args []types.Value) (types.Value, error) {
					seconds := args[0]

					if secondsVal, ok := seconds.(types.IntegerValue); ok {
						time.Sleep(time.Duration(secondsVal) * time.Second)
						return types.NOTHIN, nil
					}

					return types.NOTHIN, fmt.Errorf("SLEEP: invalid argument type")
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
			return fmt.Errorf("unknown TIME declaration: %s", decl)
		}
	}

	return nil
}
