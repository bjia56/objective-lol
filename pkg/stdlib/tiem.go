package stdlib

import (
	"fmt"
	"time"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/types"
)

// RegisterTIEM registers all TIEM (time) functions with the runtime environment
func RegisterTIEM(runtime *environment.RuntimeEnvironment) {
	// NOW function - current Unix timestamp
	now := &environment.Function{
		Name:       "NOW",
		ReturnType: "INTEGR",
		IsNative:   true,
		Parameters: []environment.Parameter{},
		NativeImpl: func(args []types.Value) (types.Value, error) {
			return types.IntegerValue(time.Now().Unix()), nil
		},
	}
	runtime.RegisterNative("NOW", now)

	// MILLIS function - current Unix timestamp in milliseconds
	millis := &environment.Function{
		Name:       "MILLIS",
		ReturnType: "INTEGR",
		IsNative:   true,
		Parameters: []environment.Parameter{},
		NativeImpl: func(args []types.Value) (types.Value, error) {
			return types.IntegerValue(time.Now().UnixMilli()), nil
		},
	}
	runtime.RegisterNative("MILLIS", millis)

	// SLEEP function - sleep for specified seconds
	sleep := &environment.Function{
		Name:     "SLEEP",
		IsNative: true,
		Parameters: []environment.Parameter{
			{Name: "duration", Type: "DUBBLE"},
		},
		NativeImpl: func(args []types.Value) (types.Value, error) {
			duration := args[0]

			if doubleVal, ok := duration.(types.DoubleValue); ok {
				if float64(doubleVal) < 0 {
					return types.NOTHIN, fmt.Errorf("SLEEP: negative duration")
				}
				duration := time.Duration(float64(doubleVal) * float64(time.Second))
				time.Sleep(duration)
				return types.NOTHIN, nil
			}

			return types.NOTHIN, fmt.Errorf("SLEEP: invalid numeric type")
		},
	}
	runtime.RegisterNative("SLEEP", sleep)

	// YEAR function - current year
	year := &environment.Function{
		Name:       "YEAR",
		ReturnType: "INTEGR",
		IsNative:   true,
		Parameters: []environment.Parameter{},
		NativeImpl: func(args []types.Value) (types.Value, error) {
			return types.IntegerValue(int64(time.Now().Year())), nil
		},
	}
	runtime.RegisterNative("YEAR", year)

	// MONTH function - current month (1-12)
	month := &environment.Function{
		Name:       "MONTH",
		ReturnType: "INTEGR",
		IsNative:   true,
		Parameters: []environment.Parameter{},
		NativeImpl: func(args []types.Value) (types.Value, error) {
			return types.IntegerValue(int64(time.Now().Month())), nil
		},
	}
	runtime.RegisterNative("MONTH", month)

	// DAY function - current day of month
	day := &environment.Function{
		Name:       "DAY",
		ReturnType: "INTEGR",
		IsNative:   true,
		Parameters: []environment.Parameter{},
		NativeImpl: func(args []types.Value) (types.Value, error) {
			return types.IntegerValue(int64(time.Now().Day())), nil
		},
	}
	runtime.RegisterNative("DAY", day)

	// HOUR function - current hour (0-23)
	hour := &environment.Function{
		Name:       "HOUR",
		ReturnType: "INTEGR",
		IsNative:   true,
		Parameters: []environment.Parameter{},
		NativeImpl: func(args []types.Value) (types.Value, error) {
			return types.IntegerValue(int64(time.Now().Hour())), nil
		},
	}
	runtime.RegisterNative("HOUR", hour)

	// MINUTE function - current minute (0-59)
	minute := &environment.Function{
		Name:       "MINUTE",
		ReturnType: "INTEGR",
		IsNative:   true,
		Parameters: []environment.Parameter{},
		NativeImpl: func(args []types.Value) (types.Value, error) {
			return types.IntegerValue(int64(time.Now().Minute())), nil
		},
	}
	runtime.RegisterNative("MINUTE", minute)

	// SECOND function - current second (0-59)
	second := &environment.Function{
		Name:       "SECOND",
		ReturnType: "INTEGR",
		IsNative:   true,
		Parameters: []environment.Parameter{},
		NativeImpl: func(args []types.Value) (types.Value, error) {
			return types.IntegerValue(int64(time.Now().Second())), nil
		},
	}
	runtime.RegisterNative("SECOND", second)

	// FORMAT_TIME function - format timestamp as string
	formatTime := &environment.Function{
		Name:       "FORMAT_TIME",
		ReturnType: "STRIN",
		IsNative:   true,
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
	runtime.RegisterNative("FORMAT_TIME", formatTime)
}
