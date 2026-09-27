//go:build js

package stdlib

import (
	"fmt"
	"os"
)

// chmodFile is the GOOS=js counterpart of chmod_other.go's version: there
// is no os.Chmod under wasm (browser or Node), so RWX's setter fails with
// an ordinary Objective-LOL exception instead of failing to compile.
func chmodFile(path string, mode os.FileMode) error {
	return fmt.Errorf("chmod is not supported in this environment")
}
