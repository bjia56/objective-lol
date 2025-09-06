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
