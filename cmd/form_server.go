package cmd

import (
	"fmt"
	"garp/internal"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"
)

var formServerCmd = &cobra.Command{
	Use:   "form-server",
	Short: "Start the contact form server",
	Long: `Start the Sinatra server for handling contact form submissions with Resend email integration.

The form server provides:
  • Contact form submission handling
  • Email delivery via Resend API  
  • Input validation and spam protection
  • Rate limiting and security measures
  • Structured JSON logging
  • CORS support for cross-origin requests

Environment variables:
  RESEND_API_KEY       - Required: Resend API key for email delivery
  RESEND_FROM_EMAIL    - Required: From email address (must be verified)
  RESEND_TO_EMAIL      - Required: Recipient email address
  GARP_FORM_HOST       - Optional: Host binding (default: 0.0.0.0)
  GARP_ENV             - Optional: Environment (development/production)

Examples:
  garp form-server                    Start server on default port 4567
  garp form-server --port 8080       Start server on custom port
  garp form-server --host localhost  Bind to localhost only`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return startFormServer()
	},
}

var (
	formPort int
	formHost string
)

func startFormServer() error {
	// Set environment variables for the Ruby server
	os.Setenv("GARP_FORM_PORT", strconv.Itoa(formPort))
	if formHost != "" {
		os.Setenv("GARP_FORM_HOST", formHost)
	}

	// Check if form-server.rb exists in current directory
	if _, err := os.Stat("form-server.rb"); os.IsNotExist(err) {
		return internal.NewFileSystemErrorWithContext(
			"form-server.rb not found in current directory",
			"Make sure you're running this command from a Garp project root directory that has forms enabled",
			err,
		)
	}

	// Check if Ruby is available
	if _, err := exec.LookPath("ruby"); err != nil {
		return internal.NewDependencyErrorWithSuggestions(
			"Ruby not found in PATH",
			err,
			[]string{
				"macOS: brew install ruby",
				"Ubuntu: sudo apt install ruby",
				"Windows: Download from https://rubyinstaller.org/",
				"Verify installation with: ruby --version",
			},
		)
	}

	fmt.Printf("🚀 Starting Garp Form Server...\n")
	fmt.Printf("📧 Port: %d\n", formPort)
	if formHost != "" {
		fmt.Printf("🌐 Host: %s\n", formHost)
	}
	fmt.Printf("📝 Logs: form-submissions.log\n")
	fmt.Printf("💡 Use Ctrl+C to stop the server\n\n")

	// Create and start the Ruby process
	rubyCmd := exec.Command("ruby", "form-server.rb")
	rubyCmd.Stdout = os.Stdout
	rubyCmd.Stderr = os.Stderr
	rubyCmd.Env = os.Environ()

	// Start the server
	if err := rubyCmd.Start(); err != nil {
		return internal.NewExternalError("Failed to start form server", err)
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	go func() {
		<-sigChan
		fmt.Printf("\n🛑 Shutting down form server...\n")

		// Kill the Ruby process
		if rubyCmd.Process != nil {
			rubyCmd.Process.Kill()
		}
		os.Exit(0)
	}()

	// Wait for the Ruby process to complete
	if err := rubyCmd.Wait(); err != nil {
		// Don't return error for normal shutdown (killed by signal)
		if err.Error() == "signal: killed" {
			return nil
		}
		return internal.NewExternalError("Form server encountered an error", err)
	}

	return nil
}

func init() {
	formServerCmd.Flags().IntVarP(&formPort, "port", "p", 4567, "Port for form server")
	formServerCmd.Flags().StringVarP(&formHost, "host", "H", "", "Host to bind to (default: 0.0.0.0)")
	rootCmd.AddCommand(formServerCmd)
}
