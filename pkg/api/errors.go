package api

import (
	"fmt"
	"strings"
	"time"
)

// ErrorType represents the type of error that occurred
type ErrorType string

const (
	// CompileErrorType indicates a parsing or compilation error
	CompileErrorType ErrorType = "compile"
	// RuntimeErrorType indicates an execution error
	RuntimeErrorType ErrorType = "runtime"
	// TimeoutErrorType indicates execution exceeded the timeout
	TimeoutErrorType ErrorType = "timeout"
	// ConversionErrorType indicates a type conversion error
	ConversionErrorType ErrorType = "conversion"
	// ConfigErrorType indicates a configuration error
	ConfigErrorType ErrorType = "config"
)

// VMError represents errors that can occur in the VM API
type VMError struct {
	Type     ErrorType
	Message  string
	Source   *SourceLocation
	Wrapped  error
	Duration time.Duration // For timeout errors
}

// SourceLocation represents a location in source code
type SourceLocation struct {
	Filename string
	Line     int
	Column   int
}

func (e *VMError) Error() string {
	var parts []string

	parts = append(parts, fmt.Sprintf("%s error", e.Type))

	if e.Source != nil {
		if e.Source.Filename != "" {
			parts = append(parts, fmt.Sprintf("at %s:%d:%d", e.Source.Filename, e.Source.Line, e.Source.Column))
		} else {
			parts = append(parts, fmt.Sprintf("at line %d:%d", e.Source.Line, e.Source.Column))
		}
	}

	parts = append(parts, e.Message, e.Wrapped.Error())

	if e.Duration > 0 {
		parts = append(parts, fmt.Sprintf("(after %v)", e.Duration))
	}

	return strings.Join(parts, ": ")
}

func (e *VMError) Unwrap() error {
	return e.Wrapped
}

// IsCompileError returns true if the error is a compilation error
func (e *VMError) IsCompileError() bool {
	return e.Type == CompileErrorType
}

// IsRuntimeError returns true if the error is a runtime error
func (e *VMError) IsRuntimeError() bool {
	return e.Type == RuntimeErrorType
}

// IsTimeoutError returns true if the error is a timeout error
func (e *VMError) IsTimeoutError() bool {
	return e.Type == TimeoutErrorType
}

// IsConversionError returns true if the error is a type conversion error
func (e *VMError) IsConversionError() bool {
	return e.Type == ConversionErrorType
}

// IsConfigError returns true if the error is a configuration error
func (e *VMError) IsConfigError() bool {
	return e.Type == ConfigErrorType
}

// NewCompileError creates a new compile error
func NewCompileError(message string, source *SourceLocation) *VMError {
	return &VMError{
		Type:    CompileErrorType,
		Message: message,
		Source:  source,
	}
}

// NewRuntimeError creates a new runtime error
func NewRuntimeError(message string, source *SourceLocation) *VMError {
	return &VMError{
		Type:    RuntimeErrorType,
		Message: message,
		Source:  source,
	}
}

// NewTimeoutError creates a new timeout error
func NewTimeoutError(duration time.Duration) *VMError {
	return &VMError{
		Type:     TimeoutErrorType,
		Message:  "execution exceeded timeout",
		Duration: duration,
	}
}

// NewConversionError creates a new type conversion error
func NewConversionError(message string, wrapped error) *VMError {
	return &VMError{
		Type:    ConversionErrorType,
		Message: message,
		Wrapped: wrapped,
	}
}

// NewConfigError creates a new configuration error
func NewConfigError(message string, wrapped error) *VMError {
	return &VMError{
		Type:    ConfigErrorType,
		Message: message,
		Wrapped: wrapped,
	}
}

// wrapError converts a generic error into a VMError if it isn't already one
func wrapError(err error, errorType ErrorType, message string) *VMError {
	if vmErr, ok := err.(*VMError); ok {
		return vmErr
	}

	return &VMError{
		Type:    errorType,
		Message: message,
		Wrapped: err,
	}
}
