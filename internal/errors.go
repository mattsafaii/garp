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

// Color constants for terminal output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorGreen  = "\033[32m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[37m"
	ColorBold   = "\033[1m"
)

// colorEnabled checks if color output should be enabled
func colorEnabled() bool {
	// Disable colors if NO_COLOR environment variable is set
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	// Check if stdout is a terminal
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		return true
	}
	return false
}

// colorize applies color to text if colors are enabled
func colorize(color, text string) string {
	if colorEnabled() {
		return color + text + ColorReset
	}
	return text
}

// ErrorType represents different categories of errors
type ErrorType int

const (
	ErrorTypeValidation ErrorType = iota
	ErrorTypeFileSystem
	ErrorTypeNetwork
	ErrorTypeConfiguration
	ErrorTypeExternal
)

// String returns the string representation of the error type
func (e ErrorType) String() string {
	switch e {
	case ErrorTypeValidation:
		return "validation"
	case ErrorTypeFileSystem:
		return "filesystem"
	case ErrorTypeNetwork:
		return "network"
	case ErrorTypeConfiguration:
		return "configuration"
	case ErrorTypeExternal:
		return "external"
	default:
		return "unknown"
	}
}

// AppError represents a structured application error
type AppError struct {
	Type        ErrorType
	Message     string
	ExitCode    int
	Cause       error
	Suggestions []string // Helpful suggestions for the user
	Context     string   // Additional context information
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Error creation functions with enhanced support for context and suggestions

func NewValidationError(message string) *AppError {
	return &AppError{
		Type:     ErrorTypeValidation,
		Message:  message,
		ExitCode: ExitMisuse,
	}
}

func NewValidationErrorWithSuggestions(message string, suggestions []string) *AppError {
	return &AppError{
		Type:        ErrorTypeValidation,
		Message:     message,
		ExitCode:    ExitMisuse,
		Suggestions: suggestions,
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

func NewFileSystemErrorWithContext(message, context string, cause error) *AppError {
	return &AppError{
		Type:     ErrorTypeFileSystem,
		Message:  message,
		ExitCode: ExitIOErr,
		Cause:    cause,
		Context:  context,
	}
}

func NewConfigurationError(message string) *AppError {
	return &AppError{
		Type:     ErrorTypeConfiguration,
		Message:  message,
		ExitCode: ExitConfig,
	}
}

func NewConfigurationErrorWithSuggestions(message string, suggestions []string) *AppError {
	return &AppError{
		Type:        ErrorTypeConfiguration,
		Message:     message,
		ExitCode:    ExitConfig,
		Suggestions: suggestions,
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

func NewDependencyError(message string, cause error) *AppError {
	return NewDependencyErrorWithSuggestions(message, cause, nil)
}

func NewDependencyErrorWithSuggestions(message string, cause error, suggestions []string) *AppError {
	return &AppError{
		Type:        ErrorTypeExternal,
		Message:     message,
		ExitCode:    ExitUnavailable,
		Cause:       cause,
		Suggestions: suggestions,
	}
}

// HandleError provides consistent error handling and exit with enhanced user messaging
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

	// Log the error with context for debugging
	logFields := []string{
		"type", appErr.Type.String(),
		"exit_code", fmt.Sprintf("%d", appErr.ExitCode),
	}
	if appErr.Context != "" {
		logFields = append(logFields, "context", appErr.Context)
	}
	if len(appErr.Suggestions) > 0 {
		logFields = append(logFields, "suggestions", fmt.Sprintf("%d_provided", len(appErr.Suggestions)))
	}

	if appErr.Cause != nil {
		LogErrorWithError(appErr.Message, appErr.Cause, logFields...)
	} else {
		LogError(appErr.Message, logFields...)
	}

	// Print the main error message with color
	errorIcon := colorize(ColorRed, "❌")
	errorLabel := colorize(ColorRed+ColorBold, "Error:")
	fmt.Fprintf(os.Stderr, "%s %s %s\n", errorIcon, errorLabel, appErr.Message)

	// Print cause if available and different from main message
	if appErr.Cause != nil && appErr.Cause.Error() != appErr.Message {
		causeText := colorize(ColorGray, fmt.Sprintf("   Cause: %v", appErr.Cause))
		fmt.Fprintf(os.Stderr, "%s\n", causeText)
	}

	// Print context if provided
	if appErr.Context != "" {
		contextText := colorize(ColorBlue, fmt.Sprintf("   Context: %s", appErr.Context))
		fmt.Fprintf(os.Stderr, "%s\n", contextText)
	}

	// Print suggestions if available
	if len(appErr.Suggestions) > 0 {
		suggestionIcon := colorize(ColorYellow, "💡")
		suggestionLabel := colorize(ColorYellow+ColorBold, "Suggestions:")
		fmt.Fprintf(os.Stderr, "\n%s %s\n", suggestionIcon, suggestionLabel)
		for _, suggestion := range appErr.Suggestions {
			bulletPoint := colorize(ColorYellow, "   •")
			fmt.Fprintf(os.Stderr, "%s %s\n", bulletPoint, suggestion)
		}
	} else {
		// Add default helpful context for common error types
		suggestions := getDefaultSuggestions(appErr.Type)
		if len(suggestions) > 0 {
			suggestionIcon := colorize(ColorYellow, "💡")
			suggestionLabel := colorize(ColorYellow+ColorBold, "Try:")
			fmt.Fprintf(os.Stderr, "\n%s %s\n", suggestionIcon, suggestionLabel)
			for _, suggestion := range suggestions {
				bulletPoint := colorize(ColorYellow, "   •")
				fmt.Fprintf(os.Stderr, "%s %s\n", bulletPoint, suggestion)
			}
		}
	}

	// Add help hint for validation errors
	if appErr.Type == ErrorTypeValidation {
		helpHint := colorize(ColorCyan, "\n💭 Run 'garp [command] --help' for usage information")
		fmt.Fprintf(os.Stderr, "%s\n", helpHint)
	}

	os.Exit(appErr.ExitCode)
}

// getDefaultSuggestions returns helpful suggestions based on error type
func getDefaultSuggestions(errorType ErrorType) []string {
	switch errorType {
	case ErrorTypeValidation:
		return []string{
			"Check your command arguments and options",
			"Ensure required parameters are provided",
			"Verify the command syntax matches expected format",
		}
	case ErrorTypeFileSystem:
		return []string{
			"Check that the file or directory exists",
			"Verify you have the necessary read/write permissions",
			"Ensure the path is correct and accessible",
		}
	case ErrorTypeConfiguration:
		return []string{
			"Verify your project configuration files",
			"Check that you're in a valid Garp project directory",
			"Run 'garp init' if you need to initialize a new project",
		}
	case ErrorTypeExternal:
		return []string{
			"Check that required dependencies are installed",
			"Verify network connectivity if needed",
			"Try running the command again",
		}
	default:
		return []string{}
	}
}
