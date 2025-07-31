package cmd

import (
	"fmt"
	"garp/internal"
	"garp/internal/scaffold"

	"github.com/spf13/cobra"
)

var (
	enableForms  bool
	enableSearch bool
	disableSearch bool
)

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new Garp static site project",
	Long: `Create a new Garp project with the complete directory structure,
template files, and configuration needed for development.

This command generates:
- public/ directory for your website content
- Caddyfile for local development server
- Tailwind CSS v4 configuration and input.css
- Build scripts for CSS and search indexing
- Example content and starter template

Optional features:
- Form server (Ruby + Sinatra) for contact forms with --forms
- Search functionality (Pagefind) enabled by default`,
	Example: `  garp init my-blog
  garp init business-site --forms
  garp init portfolio --no-search
  garp init landing-page --forms --no-search`,
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
		if enableForms {
			fmt.Printf("✓ Forms enabled (Ruby form server)\n")
		}
		if enableSearch {
			fmt.Printf("✓ Search enabled (Pagefind)\n")
		}

		// Create project structure with options
		ps := scaffold.NewProjectStructure(projectName)
		ps.EnableForms = enableForms
		ps.EnableSearch = enableSearch

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

		// Create form server files if enabled
		if err := ps.CreateFormServerFiles(); err != nil {
			return err
		}

		fmt.Printf("\n✓ Project structure created successfully!\n")
		fmt.Printf("✓ Template files generated!\n")
		fmt.Printf("✓ Configuration files created!\n")
		fmt.Printf("✓ Build scripts are executable!\n")
		
		if enableForms {
			fmt.Printf("✓ Form server files created!\n")
		}
		
		fmt.Printf("\nNext steps:\n")
		fmt.Printf("  cd %s\n", projectName)
		fmt.Printf("  cp .env.example .env  # Configure environment variables\n")
		
		if enableForms {
			fmt.Printf("  # For forms: Set RESEND_API_KEY and email settings in .env\n")
			fmt.Printf("  bundle install        # Install Ruby dependencies\n")
		}
		
		fmt.Printf("  garp serve            # Start development server\n")
		
		if enableForms {
			fmt.Printf("  garp form-server      # Start form server (in another terminal)\n")
		}

		return nil
	},
}

func init() {
	initCmd.Flags().BoolVar(&enableForms, "forms", false, "Enable form server (Ruby + Sinatra)")
	initCmd.Flags().BoolVar(&enableSearch, "search", true, "Enable search functionality (Pagefind)")
	initCmd.Flags().BoolVar(&disableSearch, "no-search", false, "Disable search functionality")
	
	// Handle the --no-search flag properly
	initCmd.PreRun = func(cmd *cobra.Command, args []string) {
		if disableSearch {
			enableSearch = false
		}
	}
	
	rootCmd.AddCommand(initCmd)
}
