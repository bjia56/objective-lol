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
	Stderr io.Writer

	// Execution configuration
	WorkingDirectory string
	ModulePaths      []string
	Timeout          time.Duration
	Sandbox          bool

	// Custom stdlib modules
	CustomStdlib map[string]interpreter.StdlibInitializer

	// Debug and profiling
	EnableProfiling bool
	EnableDebug     bool

	// Resource limits
	MaxMemory int64 // Maximum memory usage in bytes (0 = unlimited)
	MaxStack  int   // Maximum call stack depth (0 = unlimited)
}

// DefaultConfig returns a default configuration
func DefaultConfig() *VMConfig {
	return &VMConfig{
		Stdout:           os.Stdout,
		Stdin:            os.Stdin,
		Stderr:           os.Stderr,
		WorkingDirectory: ".",
		ModulePaths:      []string{},
		Timeout:          0, // No timeout by default
		Sandbox:          false,
		CustomStdlib:     make(map[string]interpreter.StdlibInitializer),
		EnableProfiling:  false,
		EnableDebug:      false,
		MaxMemory:        0, // Unlimited by default
		MaxStack:         0, // Unlimited by default
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
func WithStderr(w io.Writer) VMOption {
	return func(cfg *VMConfig) error {
		if w == nil {
			return NewConfigError("stderr writer cannot be nil", nil)
		}
		cfg.Stderr = w
		return nil
	}
}

// WithWorkingDirectory sets the working directory for module resolution
func WithWorkingDirectory(dir string) VMOption {
	return func(cfg *VMConfig) error {
		if dir == "" {
			return NewConfigError("working directory cannot be empty", nil)
		}
		cfg.WorkingDirectory = dir
		return nil
	}
}

// WithModulePaths adds additional paths for module resolution
func WithModulePaths(paths []string) VMOption {
	return func(cfg *VMConfig) error {
		cfg.ModulePaths = make([]string, len(paths))
		copy(cfg.ModulePaths, paths)
		return nil
	}
}

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

// WithSandbox enables or disables sandbox mode
func WithSandbox(enabled bool) VMOption {
	return func(cfg *VMConfig) error {
		cfg.Sandbox = enabled
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

// WithProfiling enables or disables profiling
func WithProfiling(enabled bool) VMOption {
	return func(cfg *VMConfig) error {
		cfg.EnableProfiling = enabled
		return nil
	}
}

// WithDebug enables or disables debug mode
func WithDebug(enabled bool) VMOption {
	return func(cfg *VMConfig) error {
		cfg.EnableDebug = enabled
		return nil
	}
}

// WithMaxMemory sets the maximum memory limit
func WithMaxMemory(maxMemory int64) VMOption {
	return func(cfg *VMConfig) error {
		if maxMemory < 0 {
			return NewConfigError("max memory cannot be negative", nil)
		}
		cfg.MaxMemory = maxMemory
		return nil
	}
}

// WithMaxStack sets the maximum call stack depth
func WithMaxStack(maxStack int) VMOption {
	return func(cfg *VMConfig) error {
		if maxStack < 0 {
			return NewConfigError("max stack cannot be negative", nil)
		}
		cfg.MaxStack = maxStack
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
	if cfg.Stderr == nil {
		return NewConfigError("stderr is required", nil)
	}
	if cfg.WorkingDirectory == "" {
		return NewConfigError("working directory is required", nil)
	}
	return nil
}