package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bjia56/objective-lol/cmd/olol"
)

func TestFunctionalScripts(t *testing.T) {
	testFiles := []string{}
	// Walk through the tests directory to find all .olol files
	err := filepath.WalkDir("../../tests", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".olol") {
			testFiles = append(testFiles, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to walk tests directory: %v", err)
	}

	for _, file := range testFiles {
		file := file // capture range variable
		t.Run(file, func(t *testing.T) {
			exitCode := olol.Run([]string{file})
			if exitCode != 0 {
				t.Errorf("Test failed for %s: exit code %d", file, exitCode)
			}
		})
	}
}
