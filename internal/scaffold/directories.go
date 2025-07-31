package scaffold

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mattsafaii/garp/internal"
)

// ProjectStructure defines the required directory structure for a Garp project
type ProjectStructure struct {
	ProjectName    string
	BasePath       string
	EnableForms    bool
	EnableSearch   bool
}

// CreateDirectories creates the complete directory structure for a new Garp project
func (ps *ProjectStructure) CreateDirectories() error {
	directories := []string{
		ps.ProjectName,
		filepath.Join(ps.ProjectName, "public"),
		filepath.Join(ps.ProjectName, "public", "css"),
		filepath.Join(ps.ProjectName, "public", "js"),
		filepath.Join(ps.ProjectName, "public", "images"),
		filepath.Join(ps.ProjectName, "public", "assets"),
		filepath.Join(ps.ProjectName, "bin"),
	}

	for _, dir := range directories {
		if err := ps.createDirectory(dir); err != nil {
			return err
		}
	}

	return nil
}

// createDirectory creates a single directory with proper error handling
func (ps *ProjectStructure) createDirectory(path string) error {
	// Check if directory already exists
	if _, err := os.Stat(path); err == nil {
		return internal.NewValidationError(fmt.Sprintf("directory already exists: %s", path))
	}

	// Create the directory with appropriate permissions
	if err := os.MkdirAll(path, 0755); err != nil {
		return internal.NewFileSystemError(
			fmt.Sprintf("failed to create directory: %s", path),
			err,
		)
	}

	fmt.Printf("Created directory: %s\n", path)
	return nil
}

// ValidateProjectPath checks if the project path is valid for creation
func (ps *ProjectStructure) ValidateProjectPath() error {
	projectPath := ps.ProjectName

	// Check if project directory already exists
	if _, err := os.Stat(projectPath); err == nil {
		return internal.NewValidationError(
			fmt.Sprintf("project directory '%s' already exists", projectPath),
		)
	}

	// Check if we can create the directory (test parent directory permissions)
	parentDir := filepath.Dir(projectPath)
	if parentDir == "." {
		parentDir, _ = os.Getwd()
	}

	// Test write permissions on parent directory
	if err := internal.ValidateWritableDirectory(parentDir); err != nil {
		return internal.NewFileSystemError(
			"cannot create project in current directory",
			err,
		)
	}

	return nil
}

// GetProjectStructure returns the expected structure after creation
func (ps *ProjectStructure) GetProjectStructure() []string {
	return []string{
		ps.ProjectName + "/",
		ps.ProjectName + "/public/",
		ps.ProjectName + "/public/css/",
		ps.ProjectName + "/public/js/",
		ps.ProjectName + "/public/images/",
		ps.ProjectName + "/public/assets/",
		ps.ProjectName + "/bin/",
	}
}

// NewProjectStructure creates a new ProjectStructure instance
func NewProjectStructure(projectName string) *ProjectStructure {
	return &ProjectStructure{
		ProjectName:    projectName,
		BasePath:       ".",
		EnableForms:    false,
		EnableSearch:   true,
	}
}
