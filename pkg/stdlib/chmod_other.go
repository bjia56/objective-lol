//go:build !js

package stdlib

import "os"

// chmodFile changes a file's permission bits. Split out behind a build
// tag because os.Chmod is unimplemented for GOOS=js (both the standard
// wasm port and TinyGo), which otherwise breaks compiling this package
// for the browser.
func chmodFile(path string, mode os.FileMode) error {
	return os.Chmod(path, mode)
}
