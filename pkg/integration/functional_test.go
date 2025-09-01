package integration

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/bjia56/objective-lol/cmd/olol"
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

	for _, file := range testFiles {
		file := file // capture range variable
		t.Run(path.Base(file), func(t *testing.T) {
			exitCode := olol.Run([]string{file})
			if exitCode != 0 {
				t.Errorf("Test failed for %s: exit code %d", file, exitCode)
			}
		})
	}
}
