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

// moduleSTDIOCategories defines the order that categories should be rendered in documentation
var moduleSTDIOCategories = []string{
	"output",
	"input",
}

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
					"",
					"@syntax SAY WIT <value>",
					"@param {ANY} value - Any value to print (INTEGR, DUBBLE, STRIN, BOOL, etc.)",
					"@returns {NOTHIN} No return value",
					"@example Print string without newline",
					"SAY WIT \"Hello \"",
					"SAY WIT \"World\"",
					"SAY WIT \"!\"",
					"BTW Output: Hello World!",
					"@example Print numbers",
					"SAY WIT 42",
					"SAY WIT \" is the answer\"",
					"BTW Output: 42 is the answer",
					"@example Print boolean values",
					"SAY WIT YEZ",
					"SAY WIT \" and \"",
					"SAY WIT NO",
					"BTW Output: YEZ and NO",
					"@note Does not add a newline character",
					"@note Accepts any type and converts to string representation",
					"@see SAYZ",
					"@category output",
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
					"",
					"@syntax SAYZ WIT <value>",
					"@param {ANY} value - Any value to print (INTEGR, DUBBLE, STRIN, BOOL, etc.)",
					"@returns {NOTHIN} No return value",
					"@example Print lines of text",
					"SAYZ WIT \"First line\"",
					"SAYZ WIT \"Second line\"",
					"SAYZ WIT 42",
					"BTW Output:",
					"BTW First line",
					"BTW Second line",
					"BTW 42",
					"@example Print variables",
					"I HAS A VARIABLE NAME TEH STRIN ITZ \"Alice\"",
					"I HAS A VARIABLE AGE TEH INTEGR ITZ 25",
					"SAYZ WIT NAME",
					"SAYZ WIT AGE",
					"BTW Output:",
					"BTW Alice",
					"BTW 25",
					"@note Automatically adds a newline character",
					"@note Accepts any type and converts to string representation",
					"@see SAY",
					"@category output",
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
					"",
					"@syntax GIMME",
					"@returns {STRIN} The input line with trailing newline removed",
					"@example Read user input",
					"SAYZ WIT \"Enter your name: \"",
					"I HAS A VARIABLE USER_NAME TEH STRIN ITZ GIMME",
					"SAY WIT \"Hello, \"",
					"SAYZ WIT USER_NAME",
					"BTW If user enters \"Alice\", output: Hello, Alice",
					"@example Interactive calculator",
					"SAYZ WIT \"Enter first number: \"",
					"I HAS A VARIABLE NUM1_STR TEH STRIN ITZ GIMME",
					"I HAS A VARIABLE NUM1 TEH INTEGR ITZ NUM1_STR AS INTEGR",
					"SAYZ WIT \"Enter second number: \"",
					"I HAS A VARIABLE NUM2_STR TEH STRIN ITZ GIMME",
					"I HAS A VARIABLE NUM2 TEH INTEGR ITZ NUM2_STR AS INTEGR",
					"I HAS A VARIABLE SUM TEH INTEGR ITZ NUM1 MOAR NUM2",
					"SAY WIT \"Sum: \"",
					"SAYZ WIT SUM",
					"@example Simple quiz",
					"SAYZ WIT \"What is 5 + 3?\"",
					"I HAS A VARIABLE ANSWER TEH STRIN ITZ GIMME",
					"IZ ANSWER SAEM AS \"8\"?",
					"    SAYZ WIT \"Correct!\"",
					"NOPE",
					"    SAYZ WIT \"Wrong! The answer is 8.\"",
					"KTHX",
					"@note Waits for user to press Enter",
					"@note Removes trailing newline and carriage return characters",
					"@note Returns empty string on EOF",
					"@see SAY, SAYZ",
					"@category input",
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
