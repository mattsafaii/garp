package cmd

import (
	"garp-cli/internal"
	"garp-cli/internal/server"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the local development server",
	Long: `Start the Caddy server for local development with 
live reloading and markdown rendering.

The development server provides:
- Automatic markdown rendering with Goldmark
- YAML frontmatter parsing
- Template variables for metadata
- Static file serving
- Live reloading during development`,
	Example: `  garp serve
  garp serve --port 3000
  garp serve --host 0.0.0.0 --port 8080`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate port and host
		if err := internal.ValidatePort(port); err != nil {
			return err
		}
		if err := internal.ValidateHost(host); err != nil {
			return err
		}
		
		// Validate that we're in a Garp project
		if err := internal.ValidateGarpProject(); err != nil {
			return err
		}
		
		// Create and configure Caddy server
		caddyServer := server.NewCaddyServer(host, port)
		
		// Start the server (this will block until stopped)
		// Note: ValidateConfiguration is called inside Start() after Caddyfile generation
		return caddyServer.Start()
	},
}

var (
	port int
	host string
)

func init() {
	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to serve on")
	serveCmd.Flags().StringVar(&host, "host", "localhost", "Host to bind to")
	rootCmd.AddCommand(serveCmd)
}