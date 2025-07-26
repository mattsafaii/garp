package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ValidateProjectName checks if a project name is valid
func ValidateProjectName(name string) error {
	if name == "" {
		return NewValidationError("project name cannot be empty")
	}

	if len(name) < 2 {
		return NewValidationError("project name must be at least 2 characters long")
	}

	if len(name) > 50 {
		return NewValidationError("project name must be less than 50 characters")
	}

	// Check for valid characters (alphanumeric, hyphens, underscores)
	validName := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)
	if !validName.MatchString(name) {
		return NewValidationError("project name can only contain letters, numbers, hyphens, and underscores, and must start with a letter or number")
	}

	// Check for reserved names
	reserved := []string{"garp", "bin", "src", "cmd", "internal", "test", "tests", "node_modules", ".git"}
	for _, r := range reserved {
		if strings.EqualFold(name, r) {
			return NewValidationError(fmt.Sprintf("'%s' is a reserved name and cannot be used", name))
		}
	}

	return nil
}

// ValidatePort checks if a port number is valid
func ValidatePort(port int) error {
	if port < 1 || port > 65535 {
		return NewValidationError("port must be between 1 and 65535")
	}

	if port < 1024 {
		return NewValidationError("port numbers below 1024 require elevated privileges")
	}

	return nil
}

// ValidateHost checks if a host address is valid
func ValidateHost(host string) error {
	if host == "" {
		return NewValidationError("host cannot be empty")
	}

	// Allow common host patterns
	validHosts := []string{"localhost", "127.0.0.1", "0.0.0.0", "::1"}
	for _, validHost := range validHosts {
		if host == validHost {
			return nil
		}
	}

	// Basic hostname validation (simplified)
	if len(host) > 253 {
		return NewValidationError("host name too long")
	}

	return nil
}

// ValidateDirectory checks if a directory path is valid and accessible
func ValidateDirectory(path string) error {
	if path == "" {
		return NewValidationError("directory path cannot be empty")
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return NewFileSystemError("invalid directory path", err)
	}

	// Check if directory exists
	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return NewFileSystemError(fmt.Sprintf("directory does not exist: %s", absPath), err)
		}
		return NewFileSystemError(fmt.Sprintf("cannot access directory: %s", absPath), err)
	}

	// Check if it's actually a directory
	if !info.IsDir() {
		return NewFileSystemError(fmt.Sprintf("path is not a directory: %s", absPath), nil)
	}

	return nil
}

// ValidateFile checks if a file path is valid and accessible
func ValidateFile(path string) error {
	if path == "" {
		return NewValidationError("file path cannot be empty")
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return NewFileSystemError("invalid file path", err)
	}

	// Check if file exists
	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return NewFileSystemError(fmt.Sprintf("file does not exist: %s", absPath), err)
		}
		return NewFileSystemError(fmt.Sprintf("cannot access file: %s", absPath), err)
	}

	// Check if it's actually a file
	if info.IsDir() {
		return NewFileSystemError(fmt.Sprintf("path is a directory, not a file: %s", absPath), nil)
	}

	return nil
}

// ValidateGarpProject checks if the current directory is a valid Garp project
func ValidateGarpProject() error {
	// Check for essential Garp project files
	requiredFiles := []string{
		"site/",
		"site/docs/",
		"input.css",
		"site/Caddyfile",
	}

	for _, file := range requiredFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return NewConfigurationError(fmt.Sprintf("not a Garp project: missing %s (run 'garp init' first)", file))
		}
	}

	return nil
}

// ValidateWritableDirectory checks if a directory is writable
func ValidateWritableDirectory(path string) error {
	if err := ValidateDirectory(path); err != nil {
		return err
	}

	// Try creating a temporary file to test write permissions
	tempFile := filepath.Join(path, ".garp_write_test")
	file, err := os.Create(tempFile)
	if err != nil {
		return NewFileSystemError(fmt.Sprintf("directory is not writable: %s", path), err)
	}
	file.Close()
	os.Remove(tempFile)

	return nil
}