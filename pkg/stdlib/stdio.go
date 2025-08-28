package stdlib

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/types"
)

// RegisterSTDIO registers all STDIO functions with the runtime environment
func RegisterSTDIO(runtime *environment.RuntimeEnvironment) {
	// VISIBLE function - prints value to stdout
	visible := &environment.Function{
		Name:       "VISIBLE",
		IsNative:   true,
		Parameters: []environment.Parameter{{Name: "VALUE", Type: ""}}, // Accept any type
		NativeImpl: func(args []types.Value) (types.Value, error) {
			fmt.Print(args[0].String())
			return types.NOTHIN, nil
		},
	}
	runtime.RegisterNative("VISIBLE", visible)

	// VISIBLEZ function - prints value with newline
	visiblez := &environment.Function{
		Name:       "VISIBLEZ",
		IsNative:   true,
		Parameters: []environment.Parameter{{Name: "VALUE", Type: ""}}, // Accept any type
		NativeImpl: func(args []types.Value) (types.Value, error) {
			fmt.Println(args[0].String())
			return types.NOTHIN, nil
		},
	}
	runtime.RegisterNative("VISIBLEZ", visiblez)

	// GIMME function - reads line from stdin
	gimme := &environment.Function{
		Name:       "GIMME",
		ReturnType: "STRIN",
		IsNative:   true,
		Parameters: []environment.Parameter{},
		NativeImpl: func(args []types.Value) (types.Value, error) {
			reader := bufio.NewReader(os.Stdin)
			line, err := reader.ReadString('\n')
			if err != nil {
				return types.StringValue(""), nil
			}

			// Remove trailing newline
			line = strings.TrimSuffix(line, "\n")
			line = strings.TrimSuffix(line, "\r")

			return types.StringValue(line), nil
		},
	}
	runtime.RegisterNative("GIMME", gimme)
}
