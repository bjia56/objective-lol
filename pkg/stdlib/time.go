package stdlib

import (
	"fmt"
	"time"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/types"
)

// RegisterTIMEInEnv registers all TIME functions directly in the given environment
func RegisterTIMEInEnv(env *environment.Environment) {
	// NOW function - current timestamp in seconds
	now := &environment.Function{
		Name:       "NOW",
		ReturnType: "INTEGR",
		Parameters: []environment.Parameter{},
		NativeImpl: func(args []types.Value) (types.Value, error) {
			return types.IntegerValue(time.Now().Unix()), nil
		},
	}
	env.DefineFunction(now)

	// MILLIS function - current timestamp in milliseconds
	millis := &environment.Function{
		Name:       "MILLIS",
		ReturnType: "INTEGR",
		Parameters: []environment.Parameter{},
		NativeImpl: func(args []types.Value) (types.Value, error) {
			return types.IntegerValue(time.Now().UnixMilli()), nil
		},
	}
	env.DefineFunction(millis)

	// SLEEP function - sleep for specified seconds
	sleep := &environment.Function{
		Name: "SLEEP",
		Parameters: []environment.Parameter{
			{Name: "seconds", Type: "INTEGR"},
		},
		NativeImpl: func(args []types.Value) (types.Value, error) {
			seconds := args[0]

			if secondsVal, ok := seconds.(types.IntegerValue); ok {
				time.Sleep(time.Duration(secondsVal) * time.Second)
				return types.NOTHIN, nil
			}

			return types.NOTHIN, fmt.Errorf("SLEEP: invalid argument type")
		},
	}
	env.DefineFunction(sleep)

	// YEAR function - current year
	year := &environment.Function{
		Name:       "YEAR",
		ReturnType: "INTEGR",
		Parameters: []environment.Parameter{},
		NativeImpl: func(args []types.Value) (types.Value, error) {
			return types.IntegerValue(time.Now().Year()), nil
		},
	}
	env.DefineFunction(year)

	// MONTH function - current month (1-12)
	month := &environment.Function{
		Name:       "MONTH",
		ReturnType: "INTEGR",
		Parameters: []environment.Parameter{},
		NativeImpl: func(args []types.Value) (types.Value, error) {
			return types.IntegerValue(int(time.Now().Month())), nil
		},
	}
	env.DefineFunction(month)

	// DAY function - current day of month
	day := &environment.Function{
		Name:       "DAY",
		ReturnType: "INTEGR",
		Parameters: []environment.Parameter{},
		NativeImpl: func(args []types.Value) (types.Value, error) {
			return types.IntegerValue(time.Now().Day()), nil
		},
	}
	env.DefineFunction(day)

	// HOUR function - current hour (0-23)
	hour := &environment.Function{
		Name:       "HOUR",
		ReturnType: "INTEGR",
		Parameters: []environment.Parameter{},
		NativeImpl: func(args []types.Value) (types.Value, error) {
			return types.IntegerValue(time.Now().Hour()), nil
		},
	}
	env.DefineFunction(hour)

	// MINUTE function - current minute (0-59)
	minute := &environment.Function{
		Name:       "MINUTE",
		ReturnType: "INTEGR",
		Parameters: []environment.Parameter{},
		NativeImpl: func(args []types.Value) (types.Value, error) {
			return types.IntegerValue(time.Now().Minute()), nil
		},
	}
	env.DefineFunction(minute)

	// SECOND function - current second (0-59)
	second := &environment.Function{
		Name:       "SECOND",
		ReturnType: "INTEGR",
		Parameters: []environment.Parameter{},
		NativeImpl: func(args []types.Value) (types.Value, error) {
			return types.IntegerValue(time.Now().Second()), nil
		},
	}
	env.DefineFunction(second)

	// FORMAT_TIME function - format timestamp using Go time format
	formatTime := &environment.Function{
		Name:       "FORMAT_TIME",
		ReturnType: "STRIN",
		Parameters: []environment.Parameter{
			{Name: "timestamp", Type: "INTEGR"},
			{Name: "format", Type: "STRIN"},
		},
		NativeImpl: func(args []types.Value) (types.Value, error) {
			timestamp, format := args[0], args[1]

			if timestampVal, ok := timestamp.(types.IntegerValue); ok {
				if formatVal, ok := format.(types.StringValue); ok {
					t := time.Unix(int64(timestampVal), 0)
					formatted := t.Format(formatVal.String())
					return types.StringValue(formatted), nil
				}
			}

			return types.NOTHIN, fmt.Errorf("FORMAT_TIME: invalid argument types")
		},
	}
	env.DefineFunction(formatTime)
}
