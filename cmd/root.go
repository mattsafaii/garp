package cmd

import (
	"garp-cli/internal"

	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
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
  garp deploy            Deploy to configured target`,
	Version: version,
	SuggestionsMinimumDistance: 2,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		internal.HandleError(err)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}