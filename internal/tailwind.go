package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// TailwindCLIInfo contains information about Tailwind CLI installation
type TailwindCLIInfo struct {
	IsInstalled    bool
	Version        string
	ExecutablePath string
}

// DetectTailwindCLI attempts to find and verify Tailwind CLI installation
func DetectTailwindCLI() (*TailwindCLIInfo, error) {
	info := &TailwindCLIInfo{}

	// Try different command names based on platform and installation method
	possibleCommands := []string{
		"tailwindcss",     // Standalone binary
		"tailwind",        // Alternative name
		"npx tailwindcss", // NPX execution
	}

	for _, cmd := range possibleCommands {
		if checkTailwindCommand(cmd, info) {
			return info, nil
		}
	}

	// Check common installation paths
	commonPaths := getTailwindCommonPaths()
	for _, path := range commonPaths {
		if checkTailwindPath(path, info) {
			return info, nil
		}
	}

	return info, nil
}

// checkTailwindCommand tests if a command works and gets version info
func checkTailwindCommand(cmd string, info *TailwindCLIInfo) bool {
	parts := strings.Fields(cmd)
	var execCmd *exec.Cmd

	if len(parts) > 1 {
		execCmd = exec.Command(parts[0], append(parts[1:], "--version")...)
	} else {
		execCmd = exec.Command(parts[0], "--version")
	}

	output, err := execCmd.Output()
	if err != nil {
		return false
	}

	version := strings.TrimSpace(string(output))
	if version != "" {
		info.IsInstalled = true
		info.Version = version
		info.ExecutablePath = cmd
		return true
	}

	return false
}

// checkTailwindPath tests if a specific path contains Tailwind CLI
func checkTailwindPath(path string, info *TailwindCLIInfo) bool {
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
		return true
	}

	return false
}

// getTailwindCommonPaths returns platform-specific common installation paths
func getTailwindCommonPaths() []string {
	var paths []string

	switch runtime.GOOS {
	case "windows":
		// Windows paths
		paths = append(paths,
			filepath.Join(os.Getenv("APPDATA"), "npm", "tailwindcss.exe"),
			filepath.Join(os.Getenv("LOCALAPPDATA"), "npm", "tailwindcss.exe"),
			"C:\\Program Files\\nodejs\\tailwindcss.exe",
		)
	case "darwin":
		// macOS paths
		homeDir, _ := os.UserHomeDir()
		paths = append(paths,
			"/usr/local/bin/tailwindcss",
			"/opt/homebrew/bin/tailwindcss",
			filepath.Join(homeDir, ".local", "bin", "tailwindcss"),
			filepath.Join(homeDir, ".npm-global", "bin", "tailwindcss"),
		)
	case "linux":
		// Linux paths
		homeDir, _ := os.UserHomeDir()
		paths = append(paths,
			"/usr/local/bin/tailwindcss",
			"/usr/bin/tailwindcss",
			filepath.Join(homeDir, ".local", "bin", "tailwindcss"),
			filepath.Join(homeDir, ".npm-global", "bin", "tailwindcss"),
		)
	}

	return paths
}

// GetTailwindInstallationInstructions returns platform-specific installation instructions
func GetTailwindInstallationInstructions() string {
	var instructions strings.Builder

	instructions.WriteString("Tailwind CSS CLI not found. Please install it using one of the following methods:\n\n")

	switch runtime.GOOS {
	case "windows":
		instructions.WriteString("For Windows:\n")
		instructions.WriteString("1. Using npm (recommended):\n")
		instructions.WriteString("   npm install -g tailwindcss\n\n")
		instructions.WriteString("2. Using yarn:\n")
		instructions.WriteString("   yarn global add tailwindcss\n\n")
		instructions.WriteString("3. Download standalone binary:\n")
		instructions.WriteString("   Visit: https://github.com/tailwindlabs/tailwindcss/releases\n")
		instructions.WriteString("   Download: tailwindcss-windows-x64.exe\n")
		instructions.WriteString("   Rename to: tailwindcss.exe\n")
		instructions.WriteString("   Add to your PATH\n\n")

	case "darwin":
		instructions.WriteString("For macOS:\n")
		instructions.WriteString("1. Using npm (recommended):\n")
		instructions.WriteString("   npm install -g tailwindcss\n\n")
		instructions.WriteString("2. Using Homebrew:\n")
		instructions.WriteString("   brew install tailwindcss\n\n")
		instructions.WriteString("3. Using yarn:\n")
		instructions.WriteString("   yarn global add tailwindcss\n\n")
		instructions.WriteString("4. Download standalone binary:\n")
		instructions.WriteString("   Visit: https://github.com/tailwindlabs/tailwindcss/releases\n")
		instructions.WriteString("   Download: tailwindcss-macos-x64 (Intel) or tailwindcss-macos-arm64 (Apple Silicon)\n")
		instructions.WriteString("   chmod +x tailwindcss-macos-*\n")
		instructions.WriteString("   mv tailwindcss-macos-* /usr/local/bin/tailwindcss\n\n")

	case "linux":
		instructions.WriteString("For Linux:\n")
		instructions.WriteString("1. Using npm (recommended):\n")
		instructions.WriteString("   npm install -g tailwindcss\n\n")
		instructions.WriteString("2. Using yarn:\n")
		instructions.WriteString("   yarn global add tailwindcss\n\n")
		instructions.WriteString("3. Download standalone binary:\n")
		instructions.WriteString("   Visit: https://github.com/tailwindlabs/tailwindcss/releases\n")
		instructions.WriteString("   Download: tailwindcss-linux-x64 or tailwindcss-linux-arm64\n")
		instructions.WriteString("   chmod +x tailwindcss-linux-*\n")
		instructions.WriteString("   sudo mv tailwindcss-linux-* /usr/local/bin/tailwindcss\n\n")

	default:
		instructions.WriteString("1. Using npm (recommended):\n")
		instructions.WriteString("   npm install -g tailwindcss\n\n")
		instructions.WriteString("2. Using yarn:\n")
		instructions.WriteString("   yarn global add tailwindcss\n\n")
		instructions.WriteString("3. Visit https://github.com/tailwindlabs/tailwindcss/releases for standalone binaries\n\n")
	}

	instructions.WriteString("After installation, restart your terminal and try the build command again.")

	return instructions.String()
}

// ValidateTailwindCLI checks if Tailwind CLI is available and returns helpful error if not
func ValidateTailwindCLI() error {
	info, err := DetectTailwindCLI()
	if err != nil {
		return NewDependencyError("failed to check for Tailwind CLI", err)
	}

	if !info.IsInstalled {
		instructions := GetTailwindInstallationInstructions()
		return NewDependencyError(fmt.Sprintf("Tailwind CSS CLI is required but not found.\n\n%s", instructions), nil)
	}

	return nil
}
