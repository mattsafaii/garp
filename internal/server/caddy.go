package server

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"garp/internal"
)

// CaddyServer manages the Caddy server instance
type CaddyServer struct {
	Host       string
	Port       int
	ConfigFile string
	Process    *exec.Cmd
}

// NewCaddyServer creates a new CaddyServer instance
func NewCaddyServer(host string, port int) *CaddyServer {
	return &CaddyServer{
		Host:       host,
		Port:       port,
		ConfigFile: ".garp-caddyfile-temp",
	}
}

// Start starts the Caddy server
func (cs *CaddyServer) Start() error {
	// Check if Caddy is installed
	if err := cs.checkCaddyInstallation(); err != nil {
		return err
	}

	// Generate dynamic Caddyfile
	if err := cs.generateCaddyfile(); err != nil {
		return err
	}

	// Validate the generated Caddyfile
	if err := cs.ValidateConfiguration(); err != nil {
		return err
	}

	// Ensure cleanup of temp Caddyfile
	defer cs.cleanup()

	// Start Caddy server
	fmt.Printf("Starting Caddy server on %s:%d\n", cs.Host, cs.Port)
	fmt.Printf("Using configuration: %s\n", cs.ConfigFile)

	cs.Process = exec.Command("caddy", "run", "--adapter", "caddyfile", "--config", cs.ConfigFile)
	cs.Process.Stdout = os.Stdout
	cs.Process.Stderr = os.Stderr

	if err := cs.Process.Start(); err != nil {
		return internal.NewExternalError("failed to start Caddy server", err)
	}

	fmt.Printf("✓ Server started successfully!\n")
	fmt.Printf("📖 Visit: http://%s:%d\n", cs.Host, cs.Port)
	fmt.Printf("📁 Documentation: http://%s:%d/docs/\n", cs.Host, cs.Port)
	fmt.Printf("\nPress Ctrl+C to stop the server...\n\n")

	// Set up signal handling for graceful shutdown
	cs.setupSignalHandling()

	// Wait for the process to finish
	return cs.Process.Wait()
}

// Stop stops the Caddy server gracefully
func (cs *CaddyServer) Stop() error {
	if cs.Process == nil {
		return nil
	}

	fmt.Println("\n🛑 Stopping server...")

	// Try graceful shutdown first
	if err := cs.Process.Process.Signal(syscall.SIGTERM); err != nil {
		// If graceful shutdown fails, force kill
		if killErr := cs.Process.Process.Kill(); killErr != nil {
			return internal.NewExternalError("failed to stop server", killErr)
		}
	}

	// Wait for process to exit (with timeout)
	done := make(chan error, 1)
	go func() {
		done <- cs.Process.Wait()
	}()

	select {
	case <-done:
		fmt.Println("✓ Server stopped successfully")
		return nil
	case <-time.After(5 * time.Second):
		// Force kill if graceful shutdown takes too long
		cs.Process.Process.Kill()
		fmt.Println("✓ Server forcefully stopped")
		return nil
	}
}

// checkCaddyInstallation verifies that Caddy is installed and accessible
func (cs *CaddyServer) checkCaddyInstallation() error {
	// Check if caddy command is available
	_, err := exec.LookPath("caddy")
	if err != nil {
		return internal.NewExternalError(
			"Caddy not found - please install Caddy server",
			fmt.Errorf("installation instructions: https://caddyserver.com/docs/install"),
		)
	}

	// Check Caddy version (optional - for debugging)
	cmd := exec.Command("caddy", "version")
	output, err := cmd.Output()
	if err != nil {
		return internal.NewExternalError("failed to verify Caddy installation", err)
	}

	fmt.Printf("Using Caddy: %s", string(output))
	return nil
}

// setupSignalHandling sets up graceful shutdown on interrupt signals
func (cs *CaddyServer) setupSignalHandling() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		cs.Stop()
		os.Exit(0)
	}()
}

// ValidateConfiguration validates the Caddy configuration file
func (cs *CaddyServer) ValidateConfiguration() error {
	if err := internal.ValidateFile(cs.ConfigFile); err != nil {
		return err
	}

	// Run caddy validate to check configuration syntax
	cmd := exec.Command("caddy", "validate", "--adapter", "caddyfile", "--config", cs.ConfigFile)
	if err := cmd.Run(); err != nil {
		return internal.NewConfigurationError("Caddyfile configuration is invalid - please check syntax")
	}

	fmt.Println("✓ Caddyfile configuration is valid")
	return nil
}

// generateCaddyfile creates a dynamic Caddyfile for the current host/port
func (cs *CaddyServer) generateCaddyfile() error {
	caddyfileContent := fmt.Sprintf(`{
	auto_https off
	admin off
}

http://%s:%d {
	root * .
	
	redir / /docs/
	
	handle /docs/* {
		uri strip_prefix /docs
		rewrite * /site/docs{path}.html
		try_files {path} /site/docs/index.html
		
		# Apply templates for HTML files
		templates {
			mime text/html
		}
		file_server
	}
	
	handle /docs {
		rewrite * /site/docs/index.html
		templates {
			mime text/html
		}
		file_server
	}
	
	# Serve raw markdown files for now (temporary)
	handle /docs/markdown/* {
		file_server
	}
	
	handle /style.css {
		rewrite * /site/style.css
		file_server
	}
	
	handle /_pagefind/* {
		rewrite * /site{path}
		file_server
	}
	
	file_server
	
	encode gzip
	
	log {
		output stdout
		format console
		level INFO
	}
	
	header {
		Access-Control-Allow-Origin "*"
		Access-Control-Allow-Methods "GET, POST, OPTIONS"
		Access-Control-Allow-Headers "*"
	}
	
	handle_errors {
		@404 expression {http.error.status_code} == 404
		handle @404 {
			respond "Page not found. Try visiting /docs/ for documentation." 404
		}
	}
}`, cs.Host, cs.Port)

	// Write the dynamic Caddyfile
	file, err := os.Create(cs.ConfigFile)
	if err != nil {
		return internal.NewFileSystemError("failed to create dynamic Caddyfile", err)
	}
	defer file.Close()

	if _, err := file.WriteString(caddyfileContent); err != nil {
		return internal.NewFileSystemError("failed to write dynamic Caddyfile", err)
	}

	fmt.Printf("✓ Generated dynamic Caddyfile for %s:%d\n", cs.Host, cs.Port)
	return nil
}

// cleanup removes temporary files
func (cs *CaddyServer) cleanup() {
	if cs.ConfigFile != "" && cs.ConfigFile != "site/Caddyfile" {
		os.Remove(cs.ConfigFile)
	}
}
