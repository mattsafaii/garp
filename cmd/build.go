package cmd

import (
	"fmt"
	"garp-cli/internal"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build CSS and search index",
	Long: `Execute the build process which compiles Tailwind CSS 
and generates the search index with Pagefind.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create build options from flags
		options := internal.BuildOptions{
			CSSOnly:    cssOnly,
			SearchOnly: searchOnly,
			Watch:      watch,
			Verbose:    verbose,
		}

		// Handle watch mode
		if watch {
			fmt.Println("🚀 Starting Garp in watch mode...")
			return internal.WatchFiles(options)
		}

		// Execute build
		result, err := internal.BuildAll(options)
		if err != nil {
			if result != nil && len(result.Errors) > 0 {
				for _, errMsg := range result.Errors {
					fmt.Printf("Error: %s\n", errMsg)
				}
			}
			return err
		}

		// Print summary
		if verbose || (!cssOnly && !searchOnly) {
			fmt.Printf("✅ Build completed successfully in %v\n", result.Duration)
			if result.CSSBuilt {
				fmt.Println("  📄 CSS compiled")
			}
			if result.SearchBuilt {
				fmt.Println("  🔍 Search index generated")
			}
		}

		return nil
	},
}

var (
	cssOnly    bool
	searchOnly bool
	watch      bool
	verbose    bool
)

func init() {
	buildCmd.Flags().BoolVar(&cssOnly, "css-only", false, "Build only CSS files")
	buildCmd.Flags().BoolVar(&searchOnly, "search-only", false, "Build only search index")
	buildCmd.Flags().BoolVar(&watch, "watch", false, "Watch for changes and rebuild automatically")
	buildCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show detailed build output")
	rootCmd.AddCommand(buildCmd)
}