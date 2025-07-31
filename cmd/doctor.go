package cmd

import (
	"fmt"
	"garp/internal"

	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check system dependencies and project health",
	Long: `Run diagnostic checks to verify that all required dependencies are installed
and that the current project is properly configured.

This command checks:
  • Required system dependencies (caddy, ruby, etc.)
  • Project structure and configuration files
  • File permissions and accessibility
  • Common configuration issues`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🩺 Running Garp diagnostics...")
		fmt.Println()

		// Check dependencies
		fmt.Println("📦 Checking dependencies:")
		dependencies := internal.ValidateAllDependencies()

		allDepsOK := true
		for dep, err := range dependencies {
			if err != nil {
				fmt.Printf("  ❌ %s: Not available\n", dep)
				allDepsOK = false
			} else {
				fmt.Printf("  ✅ %s: Available\n", dep)
			}
		}

		if !allDepsOK {
			fmt.Println()
			fmt.Println("⚠️  Some optional dependencies are missing.")
			fmt.Println("   Use the specific commands to see installation instructions.")
		}

		fmt.Println()

		// Check comprehensive project configuration
		fmt.Println("📁 Checking project configuration:")
		configErrors := internal.ValidateProjectConfiguration()

		if len(configErrors) == 0 {
			fmt.Println("  ✅ Project configuration is valid")
		} else {
			fmt.Printf("  ❌ Found %d configuration issues:\n", len(configErrors))
			for i, err := range configErrors {
				fmt.Printf("     %d. %s\n", i+1, err.Error())
			}
		}

		// Additional project health checks
		fmt.Println()
		fmt.Println("🔧 Checking project health:")

		// Check writable directories
		if err := internal.ValidateWritableDirectory("public"); err != nil {
			fmt.Printf("  ❌ public/ directory: %s\n", err.Error())
		} else {
			fmt.Println("  ✅ public/ directory: Writable")
		}

		if err := internal.ValidateWritableDirectory("bin"); err != nil {
			fmt.Printf("  ⚠️  bin/ directory: %s\n", err.Error())
		} else {
			fmt.Println("  ✅ bin/ directory: Writable")
		}

		fmt.Println()
		fmt.Println("✨ Diagnostics complete!")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
