package internal

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// BuildOptions configures the build process
type BuildOptions struct {
	CSSOnly    bool
	SearchOnly bool
	Watch      bool
	Verbose    bool
}

// buildMutex prevents concurrent builds
var buildMutex sync.Mutex

// BuildResult contains information about a completed build
type BuildResult struct {
	Success     bool
	Duration    time.Duration
	CSSBuilt    bool
	SearchBuilt bool
	Errors      []string
}

// BuildCSS executes the CSS build process using the bin/build-css script
func BuildCSS(options BuildOptions) (*BuildResult, error) {
	result := &BuildResult{
		CSSBuilt: true,
	}
	start := time.Now()

	if options.Verbose {
		fmt.Println("🎨 Building CSS with Tailwind...")
	}

	// Validate Tailwind CLI first
	if err := ValidateTailwindCLI(); err != nil {
		result.Errors = append(result.Errors, err.Error())
		result.Success = false
		result.Duration = time.Since(start)
		return result, err
	}

	// Check if build script exists
	buildScript := "bin/build-css"
	if _, err := os.Stat(buildScript); os.IsNotExist(err) {
		errMsg := "CSS build script not found: bin/build-css"
		result.Errors = append(result.Errors, errMsg)
		result.Success = false
		result.Duration = time.Since(start)
		return result, NewFileSystemError(errMsg, err)
	}

	// Build command arguments
	var args []string
	if options.Watch {
		args = append(args, "--watch")
	}

	// Execute build script
	cmd := exec.Command(buildScript, args...)
	cmd.Dir = "."
	
	// Capture output for verbose mode
	if options.Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	// Execute the command
	var err error
	if options.Verbose || options.Watch {
		err = cmd.Run()
	} else {
		// Capture output for non-verbose mode
		output, cmdErr := cmd.CombinedOutput()
		err = cmdErr
		if err != nil && len(output) > 0 {
			// Include output in error message for debugging
			err = fmt.Errorf("%v\nOutput: %s", err, string(output))
		}
	}

	if !options.Watch {
		if err != nil {
			errMsg := fmt.Sprintf("CSS build failed: %v", err)
			result.Errors = append(result.Errors, errMsg)
			result.Success = false
			result.Duration = time.Since(start)
			return result, NewExternalError(errMsg, err)
		}

		// Check if output file was created
		outputFile := "site/style.css"
		if _, err := os.Stat(outputFile); os.IsNotExist(err) {
			errMsg := "CSS build completed but output file not found: " + outputFile
			result.Errors = append(result.Errors, errMsg)
			result.Success = false
			result.Duration = time.Since(start)
			return result, NewFileSystemError(errMsg, err)
		}

		if options.Verbose {
			fmt.Println("✅ CSS build completed successfully")
		}
	}

	result.Success = true
	result.Duration = time.Since(start)
	return result, nil
}

// BuildSearch executes the search index build process (placeholder)
func BuildSearch(options BuildOptions) (*BuildResult, error) {
	result := &BuildResult{
		SearchBuilt: true,
	}
	start := time.Now()

	if options.Verbose {
		fmt.Println("🔍 Building search index...")
	}

	// TODO: Implement Pagefind search index building
	// For now, just return success
	if options.Verbose {
		fmt.Println("⚠️  Search index building not yet implemented")
	}

	result.Success = true
	result.Duration = time.Since(start)
	return result, nil
}

// BuildAll executes the complete build process
func BuildAll(options BuildOptions) (*BuildResult, error) {
	// Prevent concurrent builds unless in watch mode
	if !options.Watch {
		buildMutex.Lock()
		defer buildMutex.Unlock()
	}

	result := &BuildResult{}
	start := time.Now()

	if options.Verbose {
		fmt.Println("🚀 Starting Garp build process...")
	}

	// Validate project structure
	if err := ValidateGarpProject(); err != nil {
		result.Errors = append(result.Errors, err.Error())
		result.Success = false
		result.Duration = time.Since(start)
		return result, err
	}

	var errors []string

	// Build CSS unless search-only is specified
	if !options.SearchOnly {
		cssResult, err := BuildCSS(options)
		result.CSSBuilt = cssResult.CSSBuilt
		if err != nil {
			errors = append(errors, err.Error())
			if cssResult != nil {
				errors = append(errors, cssResult.Errors...)
			}
		}
	}

	// Build search index unless css-only is specified
	if !options.CSSOnly {
		searchResult, err := BuildSearch(options)
		result.SearchBuilt = searchResult.SearchBuilt
		if err != nil {
			errors = append(errors, err.Error())
			if searchResult != nil {
				errors = append(errors, searchResult.Errors...)
			}
		}
	}

	result.Errors = errors
	result.Success = len(errors) == 0
	result.Duration = time.Since(start)

	if options.Verbose {
		if result.Success {
			fmt.Printf("✅ Build completed successfully in %v\n", result.Duration)
		} else {
			fmt.Printf("❌ Build failed after %v\n", result.Duration)
		}
	}

	return result, nil
}

// WatchFiles implements file watching for automatic rebuilds
func WatchFiles(options BuildOptions) error {
	if options.Verbose {
		fmt.Println("👀 Starting file watcher...")
		fmt.Println("Watching: input.css, site/")
		fmt.Println("Press Ctrl+C to stop watching")
	}

	// Use CSS watch mode which handles file watching internally
	watchOptions := options
	watchOptions.Watch = true
	watchOptions.Verbose = true // Always show output in watch mode
	
	_, err := BuildCSS(watchOptions)
	return err
}

// GetBuildInfo returns information about the current project's build setup
func GetBuildInfo() map[string]interface{} {
	info := make(map[string]interface{})
	
	// Check for required files
	files := map[string]bool{
		"input.css":           false,
		"tailwind.config.js":  false,
		"bin/build-css":       false,
		"site/":               false,
	}
	
	for file := range files {
		if _, err := os.Stat(file); err == nil {
			files[file] = true
		}
	}
	
	info["files"] = files
	
	// Check Tailwind CLI availability
	tailwindInfo, _ := DetectTailwindCLI()
	info["tailwind"] = map[string]interface{}{
		"installed": tailwindInfo.IsInstalled,
		"version":   tailwindInfo.Version,
		"path":      tailwindInfo.ExecutablePath,
	}
	
	// Check for output files
	outputs := map[string]bool{
		"site/style.css": false,
	}
	
	for file := range outputs {
		if stat, err := os.Stat(file); err == nil && !stat.IsDir() {
			outputs[file] = true
		}
	}
	
	info["outputs"] = outputs
	
	return info
}

// CleanBuildArtifacts removes generated build files
func CleanBuildArtifacts() error {
	filesToClean := []string{
		"site/style.css",
		"site/_pagefind/",
	}
	
	var errors []string
	
	for _, file := range filesToClean {
		if err := os.RemoveAll(file); err != nil {
			errors = append(errors, fmt.Sprintf("Failed to remove %s: %v", file, err))
		}
	}
	
	if len(errors) > 0 {
		return NewFileSystemError("Clean failed: "+strings.Join(errors, "; "), nil)
	}
	
	return nil
}