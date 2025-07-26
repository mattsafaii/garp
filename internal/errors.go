package internal

import (
	"fmt"
	"os"
)

// Exit codes following Unix conventions
const (
	ExitSuccess      = 0
	ExitGeneralError = 1
	ExitMisuse       = 2
	ExitNoInput      = 66
	ExitNoHost       = 68
	ExitUnavailable  = 69
	ExitSoftware     = 70
	ExitOSFile       = 72
	ExitIOErr        = 74
	ExitTempFail     = 75
	ExitProtocol     = 76
	ExitNoPerm       = 77
	ExitConfig       = 78
)

// ErrorType represents different categories of errors
type ErrorType int

const (
	ErrorTypeValidation ErrorType = iota
	ErrorTypeFileSystem
	ErrorTypeNetwork
	ErrorTypeConfiguration
	ErrorTypeExternal
)

// AppError represents a structured application error
type AppError struct {
	Type     ErrorType
	Message  string
	ExitCode int
	Cause    error
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Error creation functions
func NewValidationError(message string) *AppError {
	return &AppError{
		Type:     ErrorTypeValidation,
		Message:  message,
		ExitCode: ExitMisuse,
	}
}

func NewFileSystemError(message string, cause error) *AppError {
	return &AppError{
		Type:     ErrorTypeFileSystem,
		Message:  message,
		ExitCode: ExitIOErr,
		Cause:    cause,
	}
}

func NewConfigurationError(message string) *AppError {
	return &AppError{
		Type:     ErrorTypeConfiguration,
		Message:  message,
		ExitCode: ExitConfig,
	}
}

func NewExternalError(message string, cause error) *AppError {
	return &AppError{
		Type:     ErrorTypeExternal,
		Message:  message,
		ExitCode: ExitUnavailable,
		Cause:    cause,
	}
}

// HandleError provides consistent error handling and exit
func HandleError(err error) {
	if err == nil {
		return
	}

	var appErr *AppError
	if appError, ok := err.(*AppError); ok {
		appErr = appError
	} else {
		// Wrap unknown errors
		appErr = &AppError{
			Type:     ErrorTypeExternal,
			Message:  "An unexpected error occurred",
			ExitCode: ExitGeneralError,
			Cause:    err,
		}
	}

	fmt.Fprintf(os.Stderr, "Error: %s\n", appErr.Message)
	
	// Add helpful context for common error types
	switch appErr.Type {
	case ErrorTypeValidation:
		fmt.Fprintf(os.Stderr, "Please check your command arguments and try again.\n")
	case ErrorTypeFileSystem:
		fmt.Fprintf(os.Stderr, "Please check file permissions and paths.\n")
	case ErrorTypeConfiguration:
		fmt.Fprintf(os.Stderr, "Please check your project configuration.\n")
	}

	os.Exit(appErr.ExitCode)
}