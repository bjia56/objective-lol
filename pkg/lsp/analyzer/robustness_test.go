package analyzer

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"testing"
	"time"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

// The LSP server runs in a single process, so any panic in the analyzer kills
// the editor integration. These tests feed it broken input and require that
// every entry point returns normally.

// exercise runs every analyzer entry point against content at a spread of
// positions (including out-of-range ones) and returns a description of any panic.
func exercise(content string) (panicInfo string) {
	defer func() {
		if r := recover(); r != nil {
			stack := string(debug.Stack())
			// keep only the first frames that matter
			lines := strings.Split(stack, "\n")
			if len(lines) > 24 {
				lines = lines[:24]
			}
			panicInfo = fmt.Sprintf("panic: %v\n%s", r, strings.Join(lines, "\n"))
		}
	}()

	// Use SemanticAnalyzer directly: Analyzer recovers from panics, which would
	// hide the underlying bugs from these tests.
	uri := "file:///tmp/robustness.olol"
	sa := NewSemanticAnalyzer(uri, content)
	_ = sa.AnalyzeDocument(context.Background())
	_ = sa.GetDiagnostics()

	lines := strings.Split(content, "\n")
	positions := []protocol.Position{
		{Line: 0, Character: 0},
		{Line: 0, Character: 1000},
		{Line: 1 << 20, Character: 0},
		{Line: 1 << 20, Character: 1 << 20},
	}
	for i := 0; i < len(lines); i += 1 + len(lines)/4 {
		l := uint32(i)
		positions = append(positions,
			protocol.Position{Line: l, Character: 0},
			protocol.Position{Line: l, Character: uint32(len(lines[i]) / 2)},
			protocol.Position{Line: l, Character: uint32(len(lines[i]))},
			protocol.Position{Line: l, Character: uint32(len(lines[i]) + 5)},
		)
	}
	for _, p := range positions {
		_ = sa.GetHoverInfo(p)
		_ = sa.GetCompletionItems(p)
		_ = sa.GetDefinitionLocation(p)
	}
	return ""
}

func loadCorpus(t *testing.T) map[string]string {
	t.Helper()
	corpus := map[string]string{}
	files, _ := filepath.Glob("../../integration/tests/*.olol")
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err == nil {
			corpus[filepath.Base(f)] = string(b)
		}
	}
	if len(corpus) == 0 {
		t.Skip("integration corpus not found")
	}
	return corpus
}

// dedupe panics by their top frame so one bug doesn't produce hundreds of lines.
type panicSet struct {
	seen  map[string]bool
	order []string
}

func (p *panicSet) add(input, info string) {
	if p.seen == nil {
		p.seen = map[string]bool{}
	}
	lines := strings.Split(info, "\n")
	key := lines[0]
	for _, l := range lines {
		if strings.Contains(l, "objective-lol/pkg") && !strings.Contains(l, "robustness_test") {
			key += " @ " + strings.TrimSpace(l)
			break
		}
	}
	if p.seen[key] {
		return
	}
	p.seen[key] = true
	snippet := input
	if len(snippet) > 300 {
		snippet = snippet[:300] + "..."
	}
	p.order = append(p.order, fmt.Sprintf("%s\ninput: %q", info, snippet))
}

func (p *panicSet) report(t *testing.T) {
	for _, m := range p.order {
		t.Errorf("%s", m)
	}
}

func TestInvalidInputHandcrafted(t *testing.T) {
	cases := map[string]string{
		"empty":                    "",
		"whitespace":               "   \n\t\n  ",
		"only newline":             "\n",
		"crlf":                     "HAI ME TEH FUNCSHUN MAIN\r\nKTHXBAI\r\n",
		"null bytes":               "HAI\x00ME\x00TEH",
		"invalid utf8":             "\xff\xfe\xfd HAI ME TEH \xc3\x28",
		"unicode":                  "HAI ME TEH FUNCSHUN 日本語 \U0001F600 KTHXBAI",
		"unterminated string":      `I HAS A VARIABLE X TEH STRIN ITZ "abc`,
		"unterminated comment":     "OBTW this never ends\nHAI ME TEH FUNCSHUN MAIN",
		"lone keyword":             "HAI",
		"lone import":              "I CAN HAS",
		"import no module":         "I CAN HAS ?",
		"import unknown":           "I CAN HAS NOPE_NOT_REAL?",
		"import selective broken":  "I CAN HAS FROM STDIO?",
		"import selective unknown": "I CAN HAS NOPE FROM STDIO?",
		"import file missing":      `I CAN HAS "does/not/exist.olol"?`,
		"import self":              `I CAN HAS "robustness.olol"?`,
		"func no name":             "HAI ME TEH FUNCSHUN",
		"func no end":              "HAI ME TEH FUNCSHUN MAIN\n  SAYZ WIT 1",
		"func bad params":          "HAI ME TEH FUNCSHUN F WIT WIT AN WIT\nKTHXBAI",
		"func dup params":          "HAI ME TEH FUNCSHUN F TEH INTEGR WIT X TEH INTEGR AN WIT X TEH INTEGR\nKTHX\nKTHXBAI",
		"class no name":            "HAI ME TEH CLAS",
		"class no end":             "HAI ME TEH CLAS FOO\n EVRYONE\n DIS TEH VARIABLE X TEH INTEGR",
		"class extends self":       "HAI ME TEH CLAS FOO KINDA LIKE FOO\nKTHXBAI",
		"class extends missing":    "HAI ME TEH CLAS FOO KINDA LIKE NOPE\nKTHXBAI",
		"class cycle":              "HAI ME TEH CLAS A KINDA LIKE B\nKTHXBAI\nHAI ME TEH CLAS B KINDA LIKE A\nKTHXBAI",
		"class dup members":        "HAI ME TEH CLAS FOO\n EVRYONE\n DIS TEH VARIABLE X TEH INTEGR\n DIS TEH VARIABLE X TEH INTEGR\nKTHXBAI",
		"member no type":           "HAI ME TEH CLAS FOO\n EVRYONE\n DIS TEH VARIABLE X\nKTHXBAI",
		"member before access":     "HAI ME TEH CLAS FOO\n DIS TEH VARIABLE X TEH INTEGR\nKTHXBAI",
		"var no type":              "I HAS A VARIABLE X",
		"var no name":              "I HAS A VARIABLE",
		"var bad init":             "I HAS A VARIABLE X TEH INTEGR ITZ",
		"var self ref":             "I HAS A VARIABLE X TEH INTEGR ITZ X",
		"var new missing class":    "HAI ME TEH FUNCSHUN MAIN\n I HAS A VARIABLE X TEH NOPE ITZ NEW NOPE\nKTHXBAI",
		"var new no class":         "HAI ME TEH FUNCSHUN MAIN\n I HAS A VARIABLE X TEH NOPE ITZ NEW\nKTHXBAI",
		"member access dangling":   "HAI ME TEH FUNCSHUN MAIN\n X DO\nKTHXBAI",
		"member access chain":      "HAI ME TEH FUNCSHUN MAIN\n A DO B DO C DO D DO\nKTHXBAI",
		"call missing wit":         "HAI ME TEH FUNCSHUN MAIN\n SAYZ WIT\nKTHXBAI",
		"call unbalanced paren":    "HAI ME TEH FUNCSHUN MAIN\n SAYZ WIT ((1 MOAR 2)\nKTHXBAI",
		"call extra paren":         "HAI ME TEH FUNCSHUN MAIN\n SAYZ WIT 1)))\nKTHXBAI",
		"binary dangling":          "HAI ME TEH FUNCSHUN MAIN\n I HAS A VARIABLE X TEH INTEGR ITZ 1 MOAR\nKTHXBAI",
		"unary dangling":           "HAI ME TEH FUNCSHUN MAIN\n I HAS A VARIABLE X TEH BOOL ITZ NOT\nKTHXBAI",
		"if no cond":               "HAI ME TEH FUNCSHUN MAIN\n IZ\nKTHXBAI",
		"if no end":                "HAI ME TEH FUNCSHUN MAIN\n IZ YEZ?\n SAYZ WIT 1\nKTHXBAI",
		"else without if":          "HAI ME TEH FUNCSHUN MAIN\n NOPE\n KTHX\nKTHXBAI",
		"while no cond":            "HAI ME TEH FUNCSHUN MAIN\n WHILE\nKTHXBAI",
		"try no catch":             "HAI ME TEH FUNCSHUN MAIN\n MAYB\n SAYZ WIT 1\nKTHXBAI",
		"catch no ident":           "HAI ME TEH FUNCSHUN MAIN\n MAYB\n KTHX\n OOPSIE\n KTHX\nKTHXBAI",
		"oops no expr":             "HAI ME TEH FUNCSHUN MAIN\n OOPS\nKTHXBAI",
		"return outside":           "GIVEZ 1",
		"stray kthx":               "KTHX\nKTHX\nKTHXBAI\nKTHXBAI",
		"stray tokens":             "))) ((( ??? !!! ,,, ... ;;;",
		"number overflow":          "I HAS A VARIABLE X TEH INTEGR ITZ 999999999999999999999999999999999999",
		"bad float":                "I HAS A VARIABLE X TEH DUBBLE ITZ 1.2.3.4",
		"minus only":               "-",
		"dot only":                 ".",
		"quote only":               `"`,
		"deep parens":              "HAI ME TEH FUNCSHUN MAIN\n SAYZ WIT " + strings.Repeat("(", 5000) + "1" + strings.Repeat(")", 5000) + "\nKTHXBAI",
		"deep unary":               "HAI ME TEH FUNCSHUN MAIN\n SAYZ WIT " + strings.Repeat("NOT ", 5000) + "YEZ\nKTHXBAI",
		"long ident":               "I HAS A VARIABLE " + strings.Repeat("X", 100000) + " TEH INTEGR",
		"many lines":               strings.Repeat("HAI ME TEH FUNCSHUN F\nKTHXBAI\n", 2000),
	}
	var ps panicSet
	for name, src := range cases {
		if info := exercise(src); info != "" {
			ps.add(name+": "+src, info)
		}
	}
	ps.report(t)
}

// Every prefix of every valid program is a realistic "user is mid-typing" state.
func TestTruncatedPrograms(t *testing.T) {
	var ps panicSet
	for name, src := range loadCorpus(t) {
		step := 1 + len(src)/80
		if len(src) > 4000 {
			step = len(src) / 60
		}
		for i := 0; i < len(src); i += step {
			if info := exercise(src[:i]); info != "" {
				ps.add(fmt.Sprintf("%s[:%d]\n%s", name, i, src[:i]), info)
			}
		}
	}
	ps.report(t)
}

// Mutations: delete a line / delete a token / swap lines / insert junk.
func TestMutatedPrograms(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	junk := []string{"KTHX", "KTHXBAI", "(", ")", "\"", "WIT", "AN WIT", "ITZ", "DO", "NEW", "TEH", "?", "\n", "\x00", "MAYB", "OOPS", "HAI ME TEH CLAS", "-", "."}
	var ps panicSet
	iterations := 25
	if testing.Short() {
		iterations = 5
	}
	for name, src := range loadCorpus(t) {
		lines := strings.Split(src, "\n")
		for n := 0; n < iterations; n++ {
			m := append([]string(nil), lines...)
			switch rng.Intn(4) {
			case 0: // delete a line
				i := rng.Intn(len(m))
				m = append(m[:i], m[i+1:]...)
			case 1: // delete a token
				i := rng.Intn(len(m))
				f := strings.Fields(m[i])
				if len(f) > 0 {
					j := rng.Intn(len(f))
					f = append(f[:j], f[j+1:]...)
				}
				m[i] = strings.Join(f, " ")
			case 2: // swap lines
				i, j := rng.Intn(len(m)), rng.Intn(len(m))
				m[i], m[j] = m[j], m[i]
			case 3: // insert junk
				i := rng.Intn(len(m))
				m[i] = m[i] + " " + junk[rng.Intn(len(junk))]
			}
			mutated := strings.Join(m, "\n")
			if info := exercise(mutated); info != "" {
				ps.add(name+"\n"+mutated, info)
			}
		}
	}
	ps.report(t)
}

// Random bytes and random token soup.
func TestRandomInput(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	words := []string{"HAI", "ME", "TEH", "FUNCSHUN", "CLAS", "VARIABLE", "I", "HAS", "A", "ITZ", "WIT", "AN", "KTHX", "KTHXBAI",
		"IZ", "NOPE", "MAYB", "OOPSIE", "OOPS", "GIVEZ", "NEW", "DO", "EVRYONE", "MYSELF", "DIS", "SHARD", "INTEGR", "STRIN",
		"BUKKIT", "BASKIT", "X", "Y", "MAIN", "1", "2.5", `"s"`, "(", ")", "?", "MOAR", "LES", "NOT", "BTW", "OBTW", "TLDR", "\n"}
	var ps panicSet
	n := 300
	if testing.Short() {
		n = 50
	}
	for i := 0; i < n; i++ {
		var sb strings.Builder
		for j := 0; j < 5+rng.Intn(80); j++ {
			if i%3 == 0 {
				sb.WriteByte(byte(rng.Intn(256)))
			} else {
				sb.WriteString(words[rng.Intn(len(words))])
				sb.WriteByte(' ')
			}
		}
		if info := exercise(sb.String()); info != "" {
			ps.add(sb.String(), info)
		}
	}
	ps.report(t)
}

// Many documents/versions through one Analyzer, concurrently (the LSP server
// handles requests on multiple goroutines).
func TestConcurrentAnalysis(t *testing.T) {
	corpus := loadCorpus(t)
	a := NewAnalyzer()
	done := make(chan string, 16)
	for g := 0; g < 8; g++ {
		go func(g int) {
			defer func() {
				if r := recover(); r != nil {
					done <- fmt.Sprintf("goroutine %d panic: %v\n%s", g, r, debug.Stack())
					return
				}
				done <- ""
			}()
			for name, src := range corpus {
				for i := 0; i < len(src); i += 1 + len(src)/10 {
					uri := fmt.Sprintf("file:///tmp/doc%d.olol", (g+i)%3)
					content := src[:i]
					_, _ = a.AnalyzeDocument(uri, content)
					_, _ = a.GetHoverInfo(uri, content, protocol.Position{Line: 1, Character: 3})
					_, _ = a.GetCompletions(uri, content, protocol.Position{Line: 2, Character: 0})
					_, _ = a.GetDefinition(uri, content, protocol.Position{Line: 3, Character: 2})
				}
				_ = name
			}
		}(g)
	}
	for g := 0; g < 8; g++ {
		select {
		case msg := <-done:
			if msg != "" {
				t.Error(msg)
			}
		case <-time.After(120 * time.Second):
			t.Fatal("timed out: possible deadlock or hang")
		}
	}
}

// Large but valid-ish input must finish in reasonable time.
func TestLargeDocumentPerformance(t *testing.T) {
	var sb strings.Builder
	sb.WriteString("I CAN HAS STDIO?\n")
	for i := 0; i < 1500; i++ {
		fmt.Fprintf(&sb, "HAI ME TEH FUNCSHUN F%d TEH INTEGR WIT X TEH INTEGR\n I HAS A VARIABLE Y TEH INTEGR ITZ X MOAR %d\n GIVEZ Y\nKTHX\n", i, i)
	}
	start := time.Now()
	done := make(chan string, 1)
	go func() { done <- exercise(sb.String()) }()
	select {
	case info := <-done:
		if info != "" {
			t.Error(info)
		}
	case <-time.After(60 * time.Second):
		t.Fatal("analysis of large document did not finish in 60s")
	}
	t.Logf("large document analyzed in %v", time.Since(start))
}

// Analyzer is the entry point used by the server; it must convert internal
// failures to errors/diagnostics rather than panic.
func TestAnalyzerAPIOnInvalidInput(t *testing.T) {
	a := NewAnalyzer()
	for _, src := range []string{"", "HAI ME TEH CLAS", "I CAN HAS", "HAI ME TEH FUNCSHUN MAIN\n SAYZ WIT 1 MOAR\nKTHXBAI"} {
		diags, err := a.AnalyzeDocument("file:///x.olol", src)
		if err == nil && len(diags) == 0 && src != "" {
			t.Errorf("expected diagnostics for invalid source %q", src)
		}
		p := protocol.Position{Line: 0, Character: 3}
		_ = exercise("\x00\xff")
		_, _ = a.GetHoverInfo("file:///x.olol", src, p)
		_, _ = a.GetCompletions("file:///x.olol", src, p)
		_, _ = a.GetDefinition("file:///x.olol", src, p)
	}
}

func FuzzAnalyzer(f *testing.F) {
	f.Add("HAI ME TEH FUNCSHUN MAIN\n SAYZ WIT 1\nKTHXBAI")
	f.Add("HAI ME TEH CLAS FOO\n EVRYONE\n DIS TEH VARIABLE X TEH INTEGR\nKTHXBAI")
	f.Fuzz(func(t *testing.T, src string) {
		if info := exercise(src); info != "" {
			t.Fatal(info)
		}
	})
}
