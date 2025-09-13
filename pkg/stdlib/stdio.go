package stdlib

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/bjia56/objective-lol/pkg/environment"
	"github.com/bjia56/objective-lol/pkg/runtime"
)

// Global I/O configuration - defaults to standard streams
var (
	StdoutWriter io.Writer = os.Stdout
	StdinReader  io.Reader = os.Stdin
)

// SetOutput sets the output writer for SAY and SAYZ functions
func SetOutput(w io.Writer) {
	StdoutWriter = w
}

// SetInput sets the input reader for GIMME function
func SetInput(r io.Reader) {
	StdinReader = r
}

// ResetToStandardStreams resets to os.Stdout and os.Stdin
func ResetToStandardStreams() {
	StdoutWriter = os.Stdout
	StdinReader = os.Stdin
}

// Global STDIO function definitions - created once and reused
var stdioFunctionsOnce = sync.Once{}
var stdioFunctions map[string]*environment.Function

func getStdioFunctions() map[string]*environment.Function {
	stdioFunctionsOnce.Do(func() {
		stdioFunctions = map[string]*environment.Function{
			"SAY": {
				Name: "SAY",
				Documentation: []string{
					"Prints a value to standard output without a newline.",
					"Accepts any type and converts it to STRIN representation.",
				},
				Parameters: []environment.Parameter{{Name: "VALUE", Type: ""}}, // Accept any type
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					fmt.Fprint(StdoutWriter, args[0].String())
					return environment.NOTHIN, nil
				},
			},
			"SAYZ": {
				Name: "SAYZ",
				Documentation: []string{
					"Prints a value to standard output followed by a newline.",
					"Accepts any type and converts it to STRIN representation.",
				},
				Parameters: []environment.Parameter{{Name: "VALUE", Type: ""}}, // Accept any type
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					fmt.Fprintln(StdoutWriter, args[0].String())
					return environment.NOTHIN, nil
				},
			},
			"GIMME": {
				Name: "GIMME",
				Documentation: []string{
					"Reads a line of input from standard input.",
					"Returns the input as a STRIN with trailing newline removed.",
				},
				ReturnType: "STRIN",
				Parameters: []environment.Parameter{},
				NativeImpl: func(interpreter environment.Interpreter, this *environment.ObjectInstance, args []environment.Value) (environment.Value, error) {
					reader := bufio.NewReader(StdinReader)
					line, err := reader.ReadString('\n')
					if err != nil {
						return environment.StringValue(""), nil
					}

					// Remove trailing newline
					line = strings.TrimSuffix(line, "\n")
					line = strings.TrimSuffix(line, "\r")

					return environment.StringValue(line), nil
				},
			},
		}
	})
	return stdioFunctions
}

// RegisterSTDIOInEnv registers STDIO functions in the given environment
// declarations: empty slice means import all, otherwise import only specified functions
func RegisterSTDIOInEnv(env *environment.Environment, declarations ...string) error {
	stdioFunctions := getStdioFunctions()

	// If declarations is empty, import all functions
	if len(declarations) == 0 {
		for _, fn := range stdioFunctions {
			env.DefineFunction(fn)
		}
		return nil
	}

	// Otherwise, import only specified functions
	for _, decl := range declarations {
		declUpper := strings.ToUpper(decl)
		if fn, exists := stdioFunctions[declUpper]; exists {
			env.DefineFunction(fn)
		} else {
			return runtime.Exception{Message: fmt.Sprintf("unknown STDIO declaration: %s", decl)}
		}
	}

	return nil
}
