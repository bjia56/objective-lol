package parser

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

// parseNoPanic parses src and fails the test on panic.
func parseNoPanic(t *testing.T, label, src string) (errs []string) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			s := src
			if len(s) > 200 {
				s = s[:200] + "..."
			}
			t.Errorf("%s: parser panicked: %v\ninput: %q", label, r, s)
		}
	}()
	p := NewParser(NewLexer(src))
	prog := p.ParseProgram()
	if prog != nil {
		// A typed-nil declaration would crash every AST consumer.
		for i, d := range prog.Declarations {
			if d == nil {
				t.Errorf("%s: nil declaration at index %d", label, i)
			}
		}
	}
	return p.Errors()
}

// Regression: failed declarations used to be stored as typed-nil ast.Node
// values (non-nil interface, nil pointer), crashing the LSP analyzer.
func TestParseErrorsDoNotProduceTypedNilNodes(t *testing.T) {
	for name, src := range map[string]string{
		"class no name":    "HAI ME TEH CLAS",
		"func no name":     "HAI ME TEH FUNCSHUN",
		"var no name":      "HAI ME TEH VARIABLE",
		"import no module": "I CAN HAS",
		"bad statement":    "HAI ME TEH FUNCSHUN MAIN\n I HAS A VARIABLE X TEH NOPE ITZ NEW NOPE\nKTHXBAI",
		"dangling binary":  "HAI ME TEH FUNCSHUN MAIN\n SAYZ WIT 1 MOAR\nKTHXBAI",
		"dangling or":      "HAI ME TEH FUNCSHUN MAIN\n I HAS A VARIABLE X TEH BOOL ITZ YEZ OR\nKTHXBAI",
		"dangling cast":    "HAI ME TEH FUNCSHUN MAIN\n I HAS A VARIABLE X TEH INTEGR ITZ AS INTEGR\nKTHXBAI",
	} {
		t.Run(name, func(t *testing.T) {
			errs := parseNoPanic(t, name, src)
			if len(errs) == 0 {
				t.Error("expected parse errors for invalid input")
			}
		})
	}
}

func TestParserInvalidInputs(t *testing.T) {
	deepParens := strings.Repeat("(", 3000) + "1" + strings.Repeat(")", 3000)
	for name, src := range map[string]string{
		"empty":         "",
		"nul bytes":     "\x00\x00\x00",
		"invalid utf8":  "\xff\xfe HAI \xc3\x28",
		"open string":   `"abc`,
		"open comment":  "OBTW never closed",
		"stray tokens":  ")))((( ??? ,,, ...",
		"deep parens":   "HAI ME TEH FUNCSHUN MAIN\n SAYZ WIT " + deepParens + "\nKTHXBAI",
		"deep unary":    "HAI ME TEH FUNCSHUN MAIN\n SAYZ WIT " + strings.Repeat("NOT ", 3000) + "YEZ\nKTHXBAI",
		"huge number":   "I HAS A VARIABLE X TEH INTEGR ITZ " + strings.Repeat("9", 400),
		"many decls":    strings.Repeat("HAI ME TEH FUNCSHUN F\nKTHXBAI\n", 5000),
		"only keywords": strings.Repeat("HAI ME TEH ", 1000),
	} {
		t.Run(name, func(t *testing.T) { parseNoPanic(t, name, src) })
	}
}

// Truncating a valid program at every byte must never panic.
func TestParserTruncation(t *testing.T) {
	src := `I CAN HAS STDIO?
HAI ME TEH VARIABLE G TEH INTEGR ITZ 5
HAI ME TEH CLAS FOO KINDA LIKE BAR
    EVRYONE
    DIS TEH VARIABLE X TEH INTEGR ITZ 1
    DIS TEH FUNCSHUN GET TEH INTEGR WIT Y TEH INTEGR
        IZ Y BIGGR THAN 0?
            GIVEZ X MOAR Y TIEMZ 2
        NOPE
            GIVEZ NOT YEZ AN NO OR YEZ
        KTHX
    KTHX
KTHXBAI
HAI ME TEH FUNCSHUN MAIN
    MAYB
        I HAS A VARIABLE F TEH FOO ITZ NEW FOO
        SAYZ WIT F DO GET WIT 3 AS STRIN
        WHILE YEZ
            OUTTA HERE
        KTHX
    OOPSIE ERR
        OOPS "bad"
    KTHX
KTHXBAI
`
	for i := 0; i <= len(src); i++ {
		parseNoPanic(t, fmt.Sprintf("trunc[:%d]", i), src[:i])
	}
}

func TestParserRandomTokenSoup(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	words := strings.Fields(`HAI ME TEH FUNCSHUN CLAS VARIABLE LOCKD I HAS A ITZ WIT AN KTHX KTHXBAI IZ NOPE MAYB OOPSIE OOPS GIVEZ
		NEW DO EVRYONE MYSELF DIS SHARD KINDA LIKE INTEGR STRIN BUKKIT BASKIT X Y MAIN 1 2.5 "s" ( ) ? MOAR LES TIEMZ DIVIDEZ NOT OR
		SAEM AS BIGGR SMALLR THAN WHILE OUTTA CAN BTW OBTW TLDR`)
	words = append(words, "\n", "\n", "\x00", "-", ".")
	for i := 0; i < 2000; i++ {
		var sb strings.Builder
		for j, n := 0, 1+rng.Intn(60); j < n; j++ {
			sb.WriteString(words[rng.Intn(len(words))])
			sb.WriteByte(' ')
		}
		parseNoPanic(t, "soup", sb.String())
	}
}

// The lexer is also fed directly; it must terminate and not panic on junk.
func TestLexerTerminatesOnGarbage(t *testing.T) {
	rng := rand.New(rand.NewSource(3))
	for i := 0; i < 500; i++ {
		b := make([]byte, rng.Intn(200))
		rng.Read(b)
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("lexer panicked on %q: %v", b, r)
				}
			}()
			l := NewLexer(string(b))
			for n := 0; n < len(b)+10; n++ {
				tok, _ := l.NextToken()
				if tok.Type == EOF {
					return
				}
			}
			t.Errorf("lexer did not reach EOF on %q", b)
		}()
	}
}

func FuzzParser(f *testing.F) {
	f.Add("HAI ME TEH FUNCSHUN MAIN\n SAYZ WIT 1 MOAR 2\nKTHXBAI")
	f.Add("HAI ME TEH CLAS FOO\n EVRYONE\n DIS TEH VARIABLE X TEH INTEGR\nKTHXBAI")
	f.Add("I CAN HAS STDIO?")
	f.Fuzz(func(t *testing.T, src string) {
		parseNoPanic(t, "fuzz", src)
	})
}
