package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// PagefindInfo contains information about Pagefind installation
type PagefindInfo struct {
	IsInstalled    bool
	Version        string
	ExecutablePath string
	IsExtended     bool
}

// DetectPagefind attempts to find and verify Pagefind installation
func DetectPagefind() (*PagefindInfo, error) {
	info := &PagefindInfo{}

	// Try different command names based on installation method
	possibleCommands := []string{
		"pagefind",          // Direct binary or cargo install
		"pagefind_extended", // Extended binary
	}

	for _, cmd := range possibleCommands {
		if checkPagefindCommand(cmd, info) {
			return info, nil
		}
	}

	// Check for NPX availability (fallback method)
	if checkNpxPagefind(info) {
		return info, nil
	}

	// Check for Python module availability
	if checkPythonPagefind(info) {
		return info, nil
	}

	// Check common installation paths
	commonPaths := getPagefindCommonPaths()
	for _, path := range commonPaths {
		if checkPagefindPath(path, info) {
			return info, nil
		}
	}

	return info, nil
}

// checkPagefindCommand tests if a command works and gets version info
func checkPagefindCommand(cmd string, info *PagefindInfo) bool {
	execCmd := exec.Command(cmd, "--version")
	output, err := execCmd.Output()
	if err != nil {
		return false
	}

	version := strings.TrimSpace(string(output))
	if version != "" {
		info.IsInstalled = true
		info.Version = version
		info.ExecutablePath = cmd
		info.IsExtended = strings.Contains(cmd, "extended")
		return true
	}

	return false
}

// checkNpxPagefind tests if NPX can run Pagefind
func checkNpxPagefind(info *PagefindInfo) bool {
	// Check if npx is available
	if _, err := exec.LookPath("npx"); err != nil {
		return false
	}

	// Test npx pagefind --version
	cmd := exec.Command("npx", "pagefind", "--version")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	version := strings.TrimSpace(string(output))
	if version != "" {
		info.IsInstalled = true
		info.Version = version
		info.ExecutablePath = "npx pagefind"
		info.IsExtended = true // NPX uses extended by default
		return true
	}

	return false
}

// checkPythonPagefind tests if Python module can run Pagefind
func checkPythonPagefind(info *PagefindInfo) bool {
	// Check if python3 is available
	if _, err := exec.LookPath("python3"); err != nil {
		return false
	}

	// Test python3 -m pagefind --version
	cmd := exec.Command("python3", "-m", "pagefind", "--version")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	version := strings.TrimSpace(string(output))
	if version != "" {
		info.IsInstalled = true
		info.Version = version
		info.ExecutablePath = "python3 -m pagefind"
		info.IsExtended = true // Python package uses extended by default
		return true
	}

	return false
}

// checkPagefindPath tests if a specific path contains Pagefind
func checkPagefindPath(path string, info *PagefindInfo) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}

	cmd := exec.Command(path, "--version")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	version := strings.TrimSpace(string(output))
	if version != "" {
		info.IsInstalled = true
		info.Version = version
		info.ExecutablePath = path
		info.IsExtended = strings.Contains(filepath.Base(path), "extended")
		return true
	}

	return false
}

// getPagefindCommonPaths returns platform-specific common installation paths
func getPagefindCommonPaths() []string {
	var paths []string

	switch runtime.GOOS {
	case "windows":
		// Windows paths
		paths = append(paths,
			filepath.Join(os.Getenv("APPDATA"), "npm", "pagefind.exe"),
			filepath.Join(os.Getenv("APPDATA"), "npm", "pagefind_extended.exe"),
			filepath.Join(os.Getenv("LOCALAPPDATA"), "npm", "pagefind.exe"),
			filepath.Join(os.Getenv("LOCALAPPDATA"), "npm", "pagefind_extended.exe"),
			"C:\\Program Files\\nodejs\\pagefind.exe",
			"C:\\Program Files\\nodejs\\pagefind_extended.exe",
		)
	case "darwin":
		// macOS paths
		homeDir, _ := os.UserHomeDir()
		paths = append(paths,
			"/usr/local/bin/pagefind",
			"/usr/local/bin/pagefind_extended",
			"/opt/homebrew/bin/pagefind",
			"/opt/homebrew/bin/pagefind_extended",
			filepath.Join(homeDir, ".local", "bin", "pagefind"),
			filepath.Join(homeDir, ".local", "bin", "pagefind_extended"),
			filepath.Join(homeDir, ".cargo", "bin", "pagefind"),
			filepath.Join(homeDir, ".cargo", "bin", "pagefind_extended"),
		)
	case "linux":
		// Linux paths
		homeDir, _ := os.UserHomeDir()
		paths = append(paths,
			"/usr/local/bin/pagefind",
			"/usr/local/bin/pagefind_extended",
			"/usr/bin/pagefind",
			"/usr/bin/pagefind_extended",
			filepath.Join(homeDir, ".local", "bin", "pagefind"),
			filepath.Join(homeDir, ".local", "bin", "pagefind_extended"),
			filepath.Join(homeDir, ".cargo", "bin", "pagefind"),
			filepath.Join(homeDir, ".cargo", "bin", "pagefind_extended"),
		)
	}

	return paths
}

// GetPagefindInstallationInstructions returns platform-specific installation instructions
func GetPagefindInstallationInstructions() string {
	var instructions strings.Builder

	instructions.WriteString("Pagefind not found. Please install it using one of the following methods:\n\n")

	switch runtime.GOOS {
	case "windows":
		instructions.WriteString("For Windows:\n")
		instructions.WriteString("1. Using NPX (recommended - no installation needed):\n")
		instructions.WriteString("   npx pagefind --site site\n\n")
		instructions.WriteString("2. Using npm globally:\n")
		instructions.WriteString("   npm install -g pagefind\n\n")
		instructions.WriteString("3. Using Python:\n")
		instructions.WriteString("   pip install 'pagefind[extended]'\n\n")
		instructions.WriteString("4. Download standalone binary:\n")
		instructions.WriteString("   Visit: https://github.com/CloudCannon/pagefind/releases\n")
		instructions.WriteString("   Download: pagefind-windows-x64.exe\n")
		instructions.WriteString("   Rename to: pagefind.exe\n")
		instructions.WriteString("   Add to your PATH\n\n")

	case "darwin":
		instructions.WriteString("For macOS:\n")
		instructions.WriteString("1. Using NPX (recommended - no installation needed):\n")
		instructions.WriteString("   npx pagefind --site site\n\n")
		instructions.WriteString("2. Using Homebrew:\n")
		instructions.WriteString("   brew install pagefind\n\n")
		instructions.WriteString("3. Using npm globally:\n")
		instructions.WriteString("   npm install -g pagefind\n\n")
		instructions.WriteString("4. Using Cargo (Rust):\n")
		instructions.WriteString("   cargo install pagefind --features extended\n\n")
		instructions.WriteString("5. Using Python:\n")
		instructions.WriteString("   pip3 install 'pagefind[extended]'\n\n")
		instructions.WriteString("6. Download standalone binary:\n")
		instructions.WriteString("   Visit: https://github.com/CloudCannon/pagefind/releases\n")
		instructions.WriteString("   Download: pagefind-macos-x64 (Intel) or pagefind-macos-arm64 (Apple Silicon)\n")
		instructions.WriteString("   chmod +x pagefind-macos-*\n")
		instructions.WriteString("   mv pagefind-macos-* /usr/local/bin/pagefind\n\n")

	case "linux":
		instructions.WriteString("For Linux:\n")
		instructions.WriteString("1. Using NPX (recommended - no installation needed):\n")
		instructions.WriteString("   npx pagefind --site site\n\n")
		instructions.WriteString("2. Using npm globally:\n")
		instructions.WriteString("   npm install -g pagefind\n\n")
		instructions.WriteString("3. Using Cargo (Rust):\n")
		instructions.WriteString("   cargo install pagefind --features extended\n\n")
		instructions.WriteString("4. Using Python:\n")
		instructions.WriteString("   pip3 install 'pagefind[extended]'\n\n")
		instructions.WriteString("5. Download standalone binary:\n")
		instructions.WriteString("   Visit: https://github.com/CloudCannon/pagefind/releases\n")
		instructions.WriteString("   Download: pagefind-linux-x64 or pagefind-linux-arm64\n")
		instructions.WriteString("   chmod +x pagefind-linux-*\n")
		instructions.WriteString("   sudo mv pagefind-linux-* /usr/local/bin/pagefind\n\n")

	default:
		instructions.WriteString("1. Using NPX (recommended - no installation needed):\n")
		instructions.WriteString("   npx pagefind --site site\n\n")
		instructions.WriteString("2. Using npm globally:\n")
		instructions.WriteString("   npm install -g pagefind\n\n")
		instructions.WriteString("3. Using Python:\n")
		instructions.WriteString("   pip install 'pagefind[extended]'\n\n")
		instructions.WriteString("4. Visit https://github.com/CloudCannon/pagefind/releases for standalone binaries\n\n")
	}

	instructions.WriteString("For more installation options, visit: https://pagefind.app/docs/installation/\n")
	instructions.WriteString("After installation, restart your terminal and try the build command again.")

	return instructions.String()
}

// ValidatePagefind checks if Pagefind is available and returns helpful error if not
func ValidatePagefind() error {
	info, err := DetectPagefind()
	if err != nil {
		return NewDependencyError("failed to check for Pagefind", err)
	}

	if !info.IsInstalled {
		instructions := GetPagefindInstallationInstructions()
		return NewDependencyError(fmt.Sprintf("Pagefind is required but not found.\n\n%s", instructions), nil)
	}

	return nil
}
