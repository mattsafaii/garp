package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build CSS and search index",
	Long: `Execute the build process which compiles Tailwind CSS 
and generates the search index with Pagefind.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Building project...")
		// TODO: Implement build orchestration
		return fmt.Errorf("build command not yet implemented")
	},
}

var (
	cssOnly    bool
	searchOnly bool
)

func init() {
	buildCmd.Flags().BoolVar(&cssOnly, "css-only", false, "Build only CSS files")
	buildCmd.Flags().BoolVar(&searchOnly, "search-only", false, "Build only search index")
	rootCmd.AddCommand(buildCmd)
}