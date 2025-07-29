package internal

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// LogLevel represents different logging levels
type LogLevel int

const (
	LogLevelError LogLevel = iota
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case LogLevelError:
		return "ERROR"
	case LogLevelWarn:
		return "WARN"
	case LogLevelInfo:
		return "INFO"
	case LogLevelDebug:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}

// ColorString returns the colored string representation for terminal output
func (l LogLevel) ColorString() string {
	switch l {
	case LogLevelError:
		return colorize(ColorRed, "ERROR")
	case LogLevelWarn:
		return colorize(ColorYellow, "WARN")
	case LogLevelInfo:
		return colorize(ColorBlue, "INFO")
	case LogLevelDebug:
		return colorize(ColorGray, "DEBUG")
	default:
		return "UNKNOWN"
	}
}

// Logger provides structured logging with different levels and file rotation
type Logger struct {
	level      LogLevel
	fileLogger *log.Logger
	logFile    *os.File
	logDir     string
	verbose    bool
}

// LoggerConfig holds configuration for the logger
type LoggerConfig struct {
	Level     LogLevel
	LogDir    string
	Verbose   bool
	MaxFiles  int
	MaxSizeMB int
}

// DefaultLoggerConfig returns sensible defaults for logging
func DefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Level:     LogLevelInfo,
		LogDir:    ".garp/logs",
		Verbose:   false,
		MaxFiles:  7,  // Keep 1 week of logs
		MaxSizeMB: 10, // 10MB max file size
	}
}

// NewLogger creates a new logger with the specified configuration
func NewLogger(config LoggerConfig) (*Logger, error) {
	// Create log directory if it doesn't exist
	if err := os.MkdirAll(config.LogDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Generate log file name with timestamp
	timestamp := time.Now().Format("2006-01-02")
	logFileName := fmt.Sprintf("garp-%s.log", timestamp)
	logFilePath := filepath.Join(config.LogDir, logFileName)

	// Open log file for appending
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Create file logger
	fileLogger := log.New(logFile, "", 0) // We'll handle our own formatting

	logger := &Logger{
		level:      config.Level,
		fileLogger: fileLogger,
		logFile:    logFile,
		logDir:     config.LogDir,
		verbose:    config.Verbose,
	}

	// Clean up old log files
	if err := logger.rotateOldLogs(config.MaxFiles); err != nil {
		// Log rotation failure shouldn't be fatal
		logger.Warn("Failed to rotate old log files", "error", err.Error())
	}

	return logger, nil
}

// Close closes the logger and its associated file
func (l *Logger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}

// SetLevel sets the minimum log level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// SetVerbose enables or disables verbose console output
func (l *Logger) SetVerbose(verbose bool) {
	l.verbose = verbose
}

// shouldLog checks if a message should be logged based on the current level
func (l *Logger) shouldLog(level LogLevel) bool {
	return level <= l.level
}

// formatLogEntry creates a formatted log entry
func (l *Logger) formatLogEntry(level LogLevel, message string, fields ...string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	
	// Create field pairs
	var fieldPairs []string
	for i := 0; i < len(fields)-1; i += 2 {
		key := fields[i]
		value := fields[i+1]
		fieldPairs = append(fieldPairs, fmt.Sprintf("%s=%s", key, value))
	}
	
	entry := fmt.Sprintf("[%s] %s %s", timestamp, level.String(), message)
	if len(fieldPairs) > 0 {
		entry += " " + strings.Join(fieldPairs, " ")
	}
	
	return entry
}

// writeToFile writes a log entry to the file
func (l *Logger) writeToFile(level LogLevel, message string, fields ...string) {
	if l.fileLogger != nil {
		entry := l.formatLogEntry(level, message, fields...)
		l.fileLogger.Println(entry)
	}
}

// writeToConsole writes a log entry to the console if verbose mode is enabled
func (l *Logger) writeToConsole(level LogLevel, message string, fields ...string) {
	if !l.verbose {
		return
	}

	timestamp := time.Now().Format("15:04:05")
	levelStr := level.ColorString()
	
	// Create field pairs for console
	var fieldPairs []string
	for i := 0; i < len(fields)-1; i += 2 {
		key := fields[i]
		value := fields[i+1]
		fieldPairs = append(fieldPairs, colorize(ColorGray, fmt.Sprintf("%s=%s", key, value)))
	}
	
	entry := fmt.Sprintf("%s %s %s", 
		colorize(ColorGray, fmt.Sprintf("[%s]", timestamp)), 
		levelStr, 
		message)
	
	if len(fieldPairs) > 0 {
		entry += " " + strings.Join(fieldPairs, " ")
	}
	
	// Write to stderr for errors and warnings, stdout for info and debug
	var output io.Writer = os.Stdout
	if level <= LogLevelWarn {
		output = os.Stderr
	}
	
	fmt.Fprintln(output, entry)
}

// log is the internal logging method
func (l *Logger) log(level LogLevel, message string, fields ...string) {
	if !l.shouldLog(level) {
		return
	}

	l.writeToFile(level, message, fields...)
	l.writeToConsole(level, message, fields...)
}

// Error logs an error message
func (l *Logger) Error(message string, fields ...string) {
	l.log(LogLevelError, message, fields...)
}

// Warn logs a warning message
func (l *Logger) Warn(message string, fields ...string) {
	l.log(LogLevelWarn, message, fields...)
}

// Info logs an info message
func (l *Logger) Info(message string, fields ...string) {
	l.log(LogLevelInfo, message, fields...)
}

// Debug logs a debug message
func (l *Logger) Debug(message string, fields ...string) {
	l.log(LogLevelDebug, message, fields...)
}

// ErrorWithError logs an error with an associated error object
func (l *Logger) ErrorWithError(message string, err error, fields ...string) {
	allFields := append([]string{"error", err.Error()}, fields...)
	l.Error(message, allFields...)
}

// rotateOldLogs removes old log files to keep only the specified number
func (l *Logger) rotateOldLogs(maxFiles int) error {
	if maxFiles <= 0 {
		return nil
	}

	// Find all garp log files
	files, err := filepath.Glob(filepath.Join(l.logDir, "garp-*.log"))
	if err != nil {
		return fmt.Errorf("failed to find log files: %w", err)
	}

	// If we have fewer files than the max, nothing to do
	if len(files) <= maxFiles {
		return nil
	}

	// Sort files by modification time (newest first)
	sort.Slice(files, func(i, j int) bool {
		infoI, errI := os.Stat(files[i])
		infoJ, errJ := os.Stat(files[j])
		if errI != nil || errJ != nil {
			return false
		}
		return infoI.ModTime().After(infoJ.ModTime())
	})

	// Remove excess files
	for i := maxFiles; i < len(files); i++ {
		if err := os.Remove(files[i]); err != nil {
			l.Warn("Failed to remove old log file", "file", files[i], "error", err.Error())
		} else {
			l.Debug("Removed old log file", "file", files[i])
		}
	}

	return nil
}

// Global logger instance
var globalLogger *Logger

// InitializeGlobalLogger initializes the global logger
func InitializeGlobalLogger(config LoggerConfig) error {
	var err error
	globalLogger, err = NewLogger(config)
	return err
}

// GetGlobalLogger returns the global logger instance
func GetGlobalLogger() *Logger {
	return globalLogger
}

// CloseGlobalLogger closes the global logger
func CloseGlobalLogger() error {
	if globalLogger != nil {
		return globalLogger.Close()
	}
	return nil
}

// Convenience functions for global logging
func LogError(message string, fields ...string) {
	if globalLogger != nil {
		globalLogger.Error(message, fields...)
	}
}

func LogWarn(message string, fields ...string) {
	if globalLogger != nil {
		globalLogger.Warn(message, fields...)
	}
}

func LogInfo(message string, fields ...string) {
	if globalLogger != nil {
		globalLogger.Info(message, fields...)
	}
}

func LogDebug(message string, fields ...string) {
	if globalLogger != nil {
		globalLogger.Debug(message, fields...)
	}
}

func LogErrorWithError(message string, err error, fields ...string) {
	if globalLogger != nil {
		globalLogger.ErrorWithError(message, err, fields...)
	}
}