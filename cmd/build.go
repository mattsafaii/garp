package cmd

import (
	"fmt"
	"github.com/mattsafaii/garp/internal"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build CSS and search index",
	Long: `Execute the build process which compiles Tailwind CSS 
and generates the search index with Pagefind.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Log build start
		internal.LogInfo("Starting build process",
			"css_only", fmt.Sprintf("%t", cssOnly),
			"search_only", fmt.Sprintf("%t", searchOnly),
			"watch", fmt.Sprintf("%t", watch))

		// Create build options from flags
		options := internal.BuildOptions{
			CSSOnly:    cssOnly,
			SearchOnly: searchOnly,
			Watch:      watch,
			Verbose:    verbose,
		}

		// Handle watch mode
		if watch {
			internal.LogInfo("Entering watch mode")
			fmt.Println("🚀 Starting Garp in watch mode...")
			return internal.WatchFiles(options)
		}

		// Execute build
		result, err := internal.BuildAll(options)
		if err != nil {
			internal.LogErrorWithError("Build failed", err,
				"css_built", fmt.Sprintf("%t", result != nil && result.CSSBuilt),
				"search_built", fmt.Sprintf("%t", result != nil && result.SearchBuilt))

			if result != nil && len(result.Errors) > 0 {
				for _, errMsg := range result.Errors {
					fmt.Printf("Error: %s\n", errMsg)
				}
			}
			return err
		}

		// Log successful build
		internal.LogInfo("Build completed successfully",
			"duration", result.Duration.String(),
			"css_built", fmt.Sprintf("%t", result.CSSBuilt),
			"search_built", fmt.Sprintf("%t", result.SearchBuilt))

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
)

func init() {
	buildCmd.Flags().BoolVar(&cssOnly, "css-only", false, "Build only CSS files")
	buildCmd.Flags().BoolVar(&searchOnly, "search-only", false, "Build only search index")
	buildCmd.Flags().BoolVar(&watch, "watch", false, "Watch for changes and rebuild automatically")
	rootCmd.AddCommand(buildCmd)
}
