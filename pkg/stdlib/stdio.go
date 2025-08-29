package stdlib

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/types"
)

// RegisterSTDIOInEnv registers all STDIO functions directly in the given environment
func RegisterSTDIOInEnv(env *environment.Environment) {
	// SAY function - prints value to stdout
	say := &environment.Function{
		Name:       "SAY",
		Parameters: []environment.Parameter{{Name: "VALUE", Type: ""}}, // Accept any type
		NativeImpl: func(args []types.Value) (types.Value, error) {
			fmt.Print(args[0].String())
			return types.NOTHIN, nil
		},
	}
	env.DefineFunction(say)

	// SAYZ function - prints value with newline
	sayz := &environment.Function{
		Name:       "SAYZ",
		Parameters: []environment.Parameter{{Name: "VALUE", Type: ""}}, // Accept any type
		NativeImpl: func(args []types.Value) (types.Value, error) {
			fmt.Println(args[0].String())
			return types.NOTHIN, nil
		},
	}
	env.DefineFunction(sayz)

	// GIMME function - reads line from stdin
	gimme := &environment.Function{
		Name:       "GIMME",
		ReturnType: "STRIN",
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
	env.DefineFunction(gimme)
}
