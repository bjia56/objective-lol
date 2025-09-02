package api

import (
	"io"
	"os"
	"time"

	"github.com/bjia56/objective-lol/pkg/interpreter"
)

// VMConfig holds configuration options for the VM
type VMConfig struct {
	// I/O configuration
	Stdout io.Writer
	Stdin  io.Reader
	//Stderr io.Writer

	// Execution configuration
	Timeout          time.Duration
	WorkingDirectory string

	// Custom stdlib modules
	CustomStdlib map[string]interpreter.StdlibInitializer
}

// DefaultConfig returns a default configuration
func DefaultConfig() *VMConfig {
	return &VMConfig{
		Stdout: os.Stdout,
		Stdin:  os.Stdin,
		//Stderr:           os.Stderr,
		Timeout:          0, // No timeout by default
		WorkingDirectory: ".",
		CustomStdlib:     make(map[string]interpreter.StdlibInitializer),
	}
}

// VMOption represents a configuration option
type VMOption func(*VMConfig) error

// WithStdout sets the standard output stream
func WithStdout(w io.Writer) VMOption {
	return func(cfg *VMConfig) error {
		if w == nil {
			return NewConfigError("stdout writer cannot be nil", nil)
		}
		cfg.Stdout = w
		return nil
	}
}

// WithStdin sets the standard input stream
func WithStdin(r io.Reader) VMOption {
	return func(cfg *VMConfig) error {
		if r == nil {
			return NewConfigError("stdin reader cannot be nil", nil)
		}
		cfg.Stdin = r
		return nil
	}
}

// WithStderr sets the standard error stream
/*
func WithStderr(w io.Writer) VMOption {
	return func(cfg *VMConfig) error {
		if w == nil {
			return NewConfigError("stderr writer cannot be nil", nil)
		}
		cfg.Stderr = w
		return nil
	}
}
*/

// WithTimeout sets the execution timeout
func WithTimeout(timeout time.Duration) VMOption {
	return func(cfg *VMConfig) error {
		if timeout < 0 {
			return NewConfigError("timeout cannot be negative", nil)
		}
		cfg.Timeout = timeout
		return nil
	}
}

// WithWorkingDirectory sets the working directory
func WithWorkingDirectory(dir string) VMOption {
	return func(cfg *VMConfig) error {
		if dir == "" {
			return NewConfigError("working directory cannot be empty", nil)
		}
		cfg.WorkingDirectory = dir
		return nil
	}
}

// WithCustomStdlib adds custom standard library modules
func WithCustomStdlib(stdlib map[string]interpreter.StdlibInitializer) VMOption {
	return func(cfg *VMConfig) error {
		if stdlib == nil {
			return NewConfigError("custom stdlib map cannot be nil", nil)
		}
		cfg.CustomStdlib = make(map[string]interpreter.StdlibInitializer)
		for k, v := range stdlib {
			cfg.CustomStdlib[k] = v
		}
		return nil
	}
}

// Validate checks if the configuration is valid
func (cfg *VMConfig) Validate() error {
	if cfg.Stdout == nil {
		return NewConfigError("stdout is required", nil)
	}
	if cfg.Stdin == nil {
		return NewConfigError("stdin is required", nil)
	}
	/*
		if cfg.Stderr == nil {
			return NewConfigError("stderr is required", nil)
		}
	*/
	if cfg.WorkingDirectory == "" {
		return NewConfigError("working directory is required", nil)
	}
	return nil
}
