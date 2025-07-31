package cmd

import (
	"fmt"
	"github.com/mattsafaii/garp/internal"
	"github.com/mattsafaii/garp/internal/server"

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
		// Log server start
		internal.LogInfo("Starting development server",
			"host", host,
			"port", fmt.Sprintf("%d", port))

		// Validate dependencies first
		if err := internal.ValidateExecutable("caddy"); err != nil {
			internal.LogErrorWithError("Caddy dependency check failed", err)
			return err
		}
		internal.LogDebug("Caddy dependency validated")

		// Validate port and host
		if err := internal.ValidatePort(port); err != nil {
			internal.LogErrorWithError("Invalid port", err, "port", fmt.Sprintf("%d", port))
			return err
		}
		if err := internal.ValidateHost(host); err != nil {
			internal.LogErrorWithError("Invalid host", err, "host", host)
			return err
		}
		internal.LogDebug("Host and port validated", "host", host, "port", fmt.Sprintf("%d", port))

		// Validate that we're in a Garp project
		if err := internal.ValidateGarpProject(); err != nil {
			internal.LogErrorWithError("Not a valid Garp project", err)
			return err
		}
		internal.LogDebug("Garp project structure validated")

		// Create and configure Caddy server
		caddyServer := server.NewCaddyServer(host, port)
		internal.LogDebug("Caddy server instance created")

		// Start the server (this will block until stopped)
		// Note: ValidateConfiguration is called inside Start() after Caddyfile generation
		internal.LogInfo("Starting Caddy server")
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
