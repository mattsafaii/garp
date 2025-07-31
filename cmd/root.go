package cmd

import (
	"fmt"
	"github.com/mattsafaii/garp/internal"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
	verbose bool
	debug   bool
)

var rootCmd = &cobra.Command{
	Use:   "garp",
	Short: "A legendary, no-nonsense static site engine",
	Long: `Garp is a lightweight static site framework that provides a simple, 
fast, production-ready way to ship content-driven websites.

Built for developers who value simplicity, performance, and maintainability 
over complex build processes.

Examples:
  garp init my-site      Create a new project called 'my-site'
  garp serve             Start development server on localhost:8080
  garp build             Build CSS and search index
  garp form-server       Start contact form server on port 4567
  garp deploy            Deploy to configured target`,
	Version:                    version,
	SuggestionsMinimumDistance: 2,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initializeLogging()
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		return internal.CloseGlobalLogger()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		internal.HandleError(err)
	}
}

// initializeLogging sets up the global logger based on command-line flags
func initializeLogging() error {
	config := internal.DefaultLoggerConfig()

	// Set log level based on flags
	if debug {
		config.Level = internal.LogLevelDebug
		config.Verbose = true
	} else if verbose {
		config.Level = internal.LogLevelInfo
		config.Verbose = true
	} else {
		config.Level = internal.LogLevelWarn
		config.Verbose = false
	}

	// Initialize the global logger
	if err := internal.InitializeGlobalLogger(config); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to initialize logging: %v\n", err)
		// Don't fail the command if logging initialization fails
	}

	return nil
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Add global flags for logging control
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug output (includes verbose)")
}
