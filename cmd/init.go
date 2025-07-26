package cmd

import (
	"fmt"
	"garp-cli/internal"
	"garp-cli/internal/scaffold"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new Garp project",
	Long: `Create a new Garp project with the complete directory structure,
template files, and configuration needed for development.

This command generates:
- site/ directory with HTML templates and markdown content
- Caddyfile for local development server
- Tailwind CSS configuration and input.css
- Build scripts for CSS and search indexing
- Example content and documentation`,
	Example: `  garp init my-blog
  garp init documentation-site
  garp init`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := "my-site"
		if len(args) > 0 {
			projectName = args[0]
		}
		
		// Validate project name
		if err := internal.ValidateProjectName(projectName); err != nil {
			return err
		}
		
		fmt.Printf("Initializing new Garp project: %s\n", projectName)
		
		// Create project structure
		ps := scaffold.NewProjectStructure(projectName)
		
		// Validate project path
		if err := ps.ValidateProjectPath(); err != nil {
			return err
		}
		
		// Create directories
		if err := ps.CreateDirectories(); err != nil {
			return err
		}
		
		// Create template files
		if err := ps.CreateTemplateFiles(); err != nil {
			return err
		}
		
		// Create configuration files
		if err := ps.CreateConfigurationFiles(); err != nil {
			return err
		}
		
		fmt.Printf("\n✓ Project structure created successfully!\n")
		fmt.Printf("✓ Template files generated!\n")
		fmt.Printf("✓ Configuration files created!\n")
		fmt.Printf("✓ Build scripts are executable!\n")
		fmt.Printf("\nNext steps:\n")
		fmt.Printf("  cd %s\n", projectName)
		fmt.Printf("  cp .env.example .env  # Configure environment variables\n")
		fmt.Printf("  garp serve            # Start development server\n")
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}