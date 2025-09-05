package integration

import (
	"io"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/bjia56/objective-lol/pkg/cli"
	"github.com/bjia56/objective-lol/pkg/stdlib"
)

func TestFunctionalScripts(t *testing.T) {
	testFiles := []string{}

	dir := "tests"
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("Failed to read tests directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".olol") {
			testFiles = append(testFiles, path.Join(dir, entry.Name()))
		}
	}

	oldstdout := os.Stdout
	stdlib.SetOutput(io.Discard)
	defer stdlib.SetOutput(oldstdout)

	for _, file := range testFiles {
		file := file // capture range variable
		t.Run(path.Base(file), func(t *testing.T) {
			err := cli.Run([]string{file})
			if err != nil {
				t.Errorf("Test failed for %s: %v", file, err)
			}
		})
	}
}
