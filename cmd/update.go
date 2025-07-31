package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/mattsafaii/garp/internal"
	"github.com/spf13/cobra"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Body    string `json:"body"`
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Garp to the latest version",
	Long: `Update Garp to the latest version available on GitHub.

This command will:
- Check the current version against the latest GitHub release
- Update via 'go install' if a newer version is available
- Show release notes for the new version`,
	Example: `  garp update
  garp update --check-only`,
	RunE: func(cmd *cobra.Command, args []string) error {
		checkOnly, _ := cmd.Flags().GetBool("check-only")

		internal.LogInfo("Checking for Garp updates...")

		// Get current version
		currentVersion := strings.TrimPrefix(version, "v")
		internal.LogDebug("Current version", "version", currentVersion)

		// Check latest version from GitHub
		latestRelease, err := getLatestRelease()
		if err != nil {
			if strings.Contains(err.Error(), "no releases found") {
				internal.LogInfo("ℹ️  No releases found on GitHub")
				fmt.Println("This appears to be a development version of Garp.")
				fmt.Println("💡 You can update by running: go install github.com/mattsafaii/garp@latest")
				return nil
			}
			internal.LogErrorWithError("Failed to check for updates", err)
			return fmt.Errorf("failed to check for updates: %w", err)
		}

		latestVersion := strings.TrimPrefix(latestRelease.TagName, "v")
		internal.LogDebug("Latest version", "version", latestVersion)

		// Compare versions
		if latestVersion == currentVersion {
			internal.LogInfo("✅ Garp is already up to date", "version", currentVersion)
			return nil
		}

		// Show update information
		fmt.Printf("🔄 Update available: %s → %s\n", currentVersion, latestVersion)
		if latestRelease.Name != "" {
			fmt.Printf("📝 Release: %s\n", latestRelease.Name)
		}

		if checkOnly {
			fmt.Println("\n💡 Run 'garp update' to install the latest version")
			return nil
		}

		// Show release notes if available
		if latestRelease.Body != "" {
			fmt.Printf("\n📋 Release notes:\n%s\n", formatReleaseNotes(latestRelease.Body))
		}

		// Confirm update
		if !confirmUpdate(latestVersion) {
			fmt.Println("❌ Update cancelled")
			return nil
		}

		// Perform update
		internal.LogInfo("Updating Garp...", "target_version", latestVersion)
		if err := performUpdate(); err != nil {
			internal.LogErrorWithError("Update failed", err)
			return fmt.Errorf("update failed: %w", err)
		}

		fmt.Printf("✅ Successfully updated to Garp %s\n", latestVersion)
		fmt.Println("💡 Run 'garp --version' to verify the update")

		return nil
	},
}

func getLatestRelease() (*GitHubRelease, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/mattsafaii/garp/releases/latest")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("no releases found on GitHub - this may be a development version")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, fmt.Errorf("failed to parse release info: %w", err)
	}

	return &release, nil
}

func formatReleaseNotes(body string) string {
	// Limit release notes to first 10 lines to keep output manageable
	lines := strings.Split(body, "\n")
	if len(lines) > 10 {
		lines = lines[:10]
		lines = append(lines, "... (truncated)")
	}
	return strings.Join(lines, "\n")
}

func confirmUpdate(version string) bool {
	fmt.Printf("\n🤔 Do you want to update to version %s? [y/N]: ", version)
	
	var response string
	fmt.Scanln(&response)
	
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

func performUpdate() error {
	internal.LogDebug("Performing update via go install")

	// Check if Go is available
	if _, err := exec.LookPath("go"); err != nil {
		return fmt.Errorf("Go is not installed or not in PATH - cannot perform automatic update")
	}

	// Run go install to update
	cmd := exec.Command("go", "install", "github.com/mattsafaii/garp@latest")
	cmd.Env = os.Environ()
	
	// Set GOOS and GOARCH to match current runtime
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOOS=%s", runtime.GOOS))
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOARCH=%s", runtime.GOARCH))

	internal.LogDebug("Running go install command")
	output, err := cmd.CombinedOutput()
	if err != nil {
		internal.LogDebug("go install failed", "output", string(output))
		return fmt.Errorf("go install failed: %w\nOutput: %s", err, string(output))
	}

	internal.LogDebug("go install completed successfully")
	return nil
}

func init() {
	updateCmd.Flags().Bool("check-only", false, "Only check for updates without installing")
	rootCmd.AddCommand(updateCmd)
}