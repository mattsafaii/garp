package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// ValidateProjectName checks if a project name is valid
func ValidateProjectName(name string) error {
	if name == "" {
		return NewValidationErrorWithSuggestions(
			"project name cannot be empty",
			[]string{
				"Provide a project name: garp init my-project",
				"Use lowercase letters, numbers, and hyphens",
				"Example: garp init my-blog",
			},
		)
	}

	if len(name) < 2 {
		return NewValidationErrorWithSuggestions(
			"project name must be at least 2 characters long",
			[]string{
				"Use a longer name like 'my-site' or 'blog'",
				"Single character names are not allowed",
			},
		)
	}

	if len(name) > 50 {
		return NewValidationErrorWithSuggestions(
			"project name must be less than 50 characters",
			[]string{
				"Choose a shorter, more concise name",
				"Long names can cause issues with file systems",
			},
		)
	}

	// Check for valid characters (alphanumeric, hyphens, underscores)
	validName := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)
	if !validName.MatchString(name) {
		return NewValidationErrorWithSuggestions(
			fmt.Sprintf("invalid project name '%s'", name),
			[]string{
				"Use only letters, numbers, hyphens, and underscores",
				"Start with a letter or number",
				"Examples: my-project, blog_site, site2024",
			},
		)
	}

	// Check for reserved names
	reserved := []string{"garp", "bin", "src", "cmd", "internal", "test", "tests", "node_modules", ".git"}
	for _, r := range reserved {
		if strings.EqualFold(name, r) {
			return NewValidationErrorWithSuggestions(
				fmt.Sprintf("'%s' is a reserved name and cannot be used", name),
				[]string{
					"Choose a different project name",
					"Reserved names can conflict with system directories",
					"Examples: my-blog, company-site, documentation",
				},
			)
		}
	}

	return nil
}

// ValidatePort checks if a port number is valid
func ValidatePort(port int) error {
	if port < 1 || port > 65535 {
		return NewValidationErrorWithSuggestions(
			fmt.Sprintf("invalid port number: %d", port),
			[]string{
				"Port must be between 1 and 65535",
				"Use ports above 1024 to avoid requiring root privileges",
				"Common development ports: 3000, 4000, 8000, 8080",
			},
		)
	}

	if port < 1024 {
		return NewValidationErrorWithSuggestions(
			fmt.Sprintf("port %d requires elevated privileges", port),
			[]string{
				"Use ports above 1024 for development",
				"Run as root (not recommended) or choose a higher port",
				"Common development ports: 3000, 4000, 8000, 8080",
			},
		)
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
	// Check for essential Garp project files and directories
	requiredFiles := []string{
		"input.css",
		"tailwind.config.js",
	}
	
	requiredDirs := []string{
		"site",
		"bin",
	}
	
	requiredProjectFiles := []string{
		"site/Caddyfile",
	}

	missing := []string{}

	// Check required files in project root
	for _, file := range requiredFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			missing = append(missing, file)
		}
	}

	// Check required directories
	for _, dir := range requiredDirs {
		if stat, err := os.Stat(dir); os.IsNotExist(err) || !stat.IsDir() {
			missing = append(missing, dir+"/")
		}
	}
	
	// Check required files within the project
	for _, file := range requiredProjectFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			missing = append(missing, file)
		}
	}

	if len(missing) > 0 {
		suggestions := []string{
			"Run 'garp init' to create a new Garp project",
			"Make sure you're in the project root directory",
		}
		
		if len(missing) < len(requiredFiles)+len(requiredDirs)+len(requiredProjectFiles) {
			suggestions = append(suggestions, "Some files exist - this might be a corrupted project")
		}

		return NewConfigurationErrorWithSuggestions(
			fmt.Sprintf("not a valid Garp project (missing: %s)", strings.Join(missing, ", ")),
			suggestions,
		)
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

// ValidateExecutable checks if a required executable is available
func ValidateExecutable(name string) error {
	if _, err := exec.LookPath(name); err != nil {
		suggestions := getInstallationSuggestions(name)
		return NewDependencyErrorWithSuggestions(
			fmt.Sprintf("%s not found in PATH", name),
			err,
			suggestions,
		)
	}
	return nil
}

// getInstallationSuggestions returns platform-specific installation instructions
func getInstallationSuggestions(executable string) []string {
	switch executable {
	case "caddy":
		return []string{
			"macOS: brew install caddy",
			"Ubuntu: sudo apt install caddy",
			"Windows: Download from https://caddyserver.com/download",
			"Or install via Go: go install github.com/caddyserver/caddy/v2/cmd/caddy@latest",
			"Verify installation: caddy version",
		}
	case "tailwindcss":
		return []string{
			"Install via npm: npm install -D tailwindcss",
			"Or download standalone binary: https://tailwindcss.com/blog/standalone-cli",
			"Verify installation: npx tailwindcss --version",
		}
	case "ruby":
		return []string{
			"macOS: brew install ruby",
			"Ubuntu: sudo apt install ruby-full",
			"Windows: Download from https://rubyinstaller.org/",
			"Verify installation: ruby --version",
		}
	case "pagefind":
		return []string{
			"Install via npm: npm install -g pagefind",
			"Or download binary: https://github.com/CloudCannon/pagefind/releases",
			"Verify installation: pagefind --version",
		}
	default:
		return []string{
			fmt.Sprintf("Install %s and ensure it's in your PATH", executable),
			"Check the official documentation for installation instructions",
		}
	}
}

// ValidateCommandPrerequisites validates common prerequisites for Garp commands
func ValidateCommandPrerequisites(command string) error {
	switch command {
	case "serve":
		return ValidateExecutable("caddy")
	case "build":
		// Build command will validate dependencies internally via existing functions
		return nil
	case "form-server":
		return ValidateExecutable("ruby")
	default:
		return nil
	}
}

// ValidateAllDependencies checks all optional dependencies and provides a report
func ValidateAllDependencies() map[string]error {
	dependencies := map[string]error{
		"caddy":       ValidateExecutable("caddy"),
		"ruby":        ValidateExecutable("ruby"),
		"tailwindcss": ValidateExecutable("tailwindcss"),
		"pagefind":    ValidateExecutable("pagefind"),
	}
	return dependencies
}

// ValidateCaddyfile checks if the Caddyfile is valid and properly configured
func ValidateCaddyfile() error {
	caddyfilePath := "site/Caddyfile"
	
	// Check if file exists
	if _, err := os.Stat(caddyfilePath); os.IsNotExist(err) {
		return NewConfigurationErrorWithSuggestions(
			"Caddyfile not found",
			[]string{
				"Run 'garp init' to create a new project with Caddyfile",
				"Make sure you're in the project root directory",
			},
		)
	}
	
	// Basic content validation
	content, err := os.ReadFile(caddyfilePath)
	if err != nil {
		return NewFileSystemError("cannot read Caddyfile", err)
	}
	
	contentStr := string(content)
	
	// Check for basic required directives
	requiredPatterns := []string{
		"root", // Should have root directive
		"file_server", // Should serve files
	}
	
	missing := []string{}
	for _, pattern := range requiredPatterns {
		if !strings.Contains(contentStr, pattern) {
			missing = append(missing, pattern)
		}
	}
	
	if len(missing) > 0 {
		return NewConfigurationErrorWithSuggestions(
			fmt.Sprintf("Caddyfile missing required directives: %s", strings.Join(missing, ", ")),
			[]string{
				"Check the Caddyfile syntax and required directives",
				"Run 'garp init' to regenerate a proper Caddyfile",
				"Refer to Caddy documentation for proper syntax",
			},
		)
	}
	
	return nil
}

// ValidateTailwindConfig checks if the Tailwind configuration is valid
func ValidateTailwindConfig() error {
	configPath := "tailwind.config.js"
	
	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return NewConfigurationErrorWithSuggestions(
			"tailwind.config.js not found",
			[]string{
				"Run 'garp init' to create a new project with Tailwind config",
				"Create a tailwind.config.js file in the project root",
			},
		)
	}
	
	// Basic content validation
	content, err := os.ReadFile(configPath)
	if err != nil {
		return NewFileSystemError("cannot read tailwind.config.js", err)
	}
	
	contentStr := string(content)
	
	// Check for basic required sections
	requiredPatterns := []string{
		"content:", // Should specify content files
		"theme:", // Should have theme configuration
	}
	
	for _, pattern := range requiredPatterns {
		if !strings.Contains(contentStr, pattern) {
			return NewConfigurationErrorWithSuggestions(
				fmt.Sprintf("tailwind.config.js missing required section: %s", pattern),
				[]string{
					"Check the Tailwind configuration syntax",
					"Run 'garp init' to regenerate a proper config",
					"Refer to Tailwind CSS documentation",
				},
			)
		}
	}
	
	return nil
}

// ValidateInputCSS checks if the input CSS file is valid
func ValidateInputCSS() error {
	inputPath := "input.css"
	
	// Check if file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return NewConfigurationErrorWithSuggestions(
			"input.css not found",
			[]string{
				"Run 'garp init' to create a new project with input.css",
				"Create an input.css file in the project root",
			},
		)
	}
	
	// Basic content validation
	content, err := os.ReadFile(inputPath)
	if err != nil {
		return NewFileSystemError("cannot read input.css", err)
	}
	
	contentStr := string(content)
	
	// Check for Tailwind directives (either individual @tailwind directives or @import)
	hasImport := strings.Contains(contentStr, "@import \"tailwindcss\"") || 
				 strings.Contains(contentStr, "@import 'tailwindcss'")
	
	if !hasImport {
		// Check for individual directives
		tailwindDirectives := []string{
			"@tailwind base",
			"@tailwind components", 
			"@tailwind utilities",
		}
		
		missing := []string{}
		for _, directive := range tailwindDirectives {
			if !strings.Contains(contentStr, directive) {
				missing = append(missing, directive)
			}
		}
		
		if len(missing) > 0 {
			return NewConfigurationErrorWithSuggestions(
				fmt.Sprintf("input.css missing Tailwind setup - neither @import nor individual directives found"),
				[]string{
					"Add '@import \"tailwindcss\";' to the top of input.css",
					"Or add individual directives: @tailwind base; @tailwind components; @tailwind utilities;",
					"Run 'garp init' to regenerate a proper input.css",
					"Refer to Tailwind CSS documentation for setup",
				},
			)
		}
	}
	
	return nil
}

// ValidateBuildScripts checks if the required build scripts exist and are executable
func ValidateBuildScripts() error {
	scripts := []string{
		"bin/build-css",
		"bin/build-search-index",
	}
	
	missing := []string{}
	notExecutable := []string{}
	
	for _, script := range scripts {
		stat, err := os.Stat(script)
		if os.IsNotExist(err) {
			missing = append(missing, script)
			continue
		}
		if err != nil {
			return NewFileSystemError(fmt.Sprintf("cannot access %s", script), err)
		}
		
		// Check if executable
		if stat.Mode()&0111 == 0 {
			notExecutable = append(notExecutable, script)
		}
	}
	
	if len(missing) > 0 {
		return NewConfigurationErrorWithSuggestions(
			fmt.Sprintf("missing build scripts: %s", strings.Join(missing, ", ")),
			[]string{
				"Run 'garp init' to create a new project with build scripts",
				"Make sure you're in the project root directory",
			},
		)
	}
	
	if len(notExecutable) > 0 {
		return NewConfigurationErrorWithSuggestions(  
			fmt.Sprintf("build scripts not executable: %s", strings.Join(notExecutable, ", ")),
			[]string{
				"Run 'chmod +x bin/*' to make scripts executable",
				"Check file permissions on build scripts",
			},
		)
	}
	
	return nil
}

// ValidateProjectConfiguration runs comprehensive configuration validation
func ValidateProjectConfiguration() []error {
	var errors []error
	
	validations := []func() error{
		ValidateGarpProject,
		ValidateCaddyfile,
		ValidateTailwindConfig,
		ValidateInputCSS,
		ValidateBuildScripts,
	}
	
	for _, validation := range validations {
		if err := validation(); err != nil {
			errors = append(errors, err)
		}
	}
	
	return errors
}