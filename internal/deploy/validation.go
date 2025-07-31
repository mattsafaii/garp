package deploy

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ValidationOptions configures deployment validation
type ValidationOptions struct {
	CheckLinks    bool
	CheckImages   bool
	CheckFileSize bool
	MaxFileSize   int64 // in bytes
	RequiredFiles []string
	Verbose       bool
}

// ValidationResult contains validation results
type ValidationResult struct {
	Success     bool
	Issues      []ValidationIssue
	FileCount   int
	TotalSize   int64
	LargestFile string
	LargestSize int64
}

// ValidationIssue represents a validation problem
type ValidationIssue struct {
	Type       string // "error", "warning"
	Category   string // "link", "image", "file", "size"
	Message    string
	File       string
	LineNumber int
}

// ValidateDeployment performs comprehensive pre-deployment validation
func ValidateDeployment(sourceDir string, options ValidationOptions) (*ValidationResult, error) {
	result := &ValidationResult{
		Success: true,
		Issues:  []ValidationIssue{},
	}

	if options.Verbose {
		fmt.Println("🔍 Running deployment validation...")
	}

	// Check if source directory exists
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("source directory '%s' does not exist", sourceDir)
	}

	// Walk through all files
	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		result.FileCount++
		result.TotalSize += info.Size()

		// Track largest file
		if info.Size() > result.LargestSize {
			result.LargestSize = info.Size()
			result.LargestFile = path
		}

		// Check individual file size
		if options.CheckFileSize && options.MaxFileSize > 0 && info.Size() > options.MaxFileSize {
			result.Issues = append(result.Issues, ValidationIssue{
				Type:     "warning",
				Category: "size",
				Message:  fmt.Sprintf("File size (%d bytes) exceeds limit (%d bytes)", info.Size(), options.MaxFileSize),
				File:     path,
			})
		}

		// Validate HTML files
		if strings.HasSuffix(strings.ToLower(path), ".html") || strings.HasSuffix(strings.ToLower(path), ".htm") {
			if err := validateHTMLFile(path, options, result); err != nil {
				if options.Verbose {
					fmt.Printf("Warning: Could not validate %s: %v\n", path, err)
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory: %v", err)
	}

	// Check for required files
	for _, requiredFile := range options.RequiredFiles {
		fullPath := filepath.Join(sourceDir, requiredFile)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			result.Issues = append(result.Issues, ValidationIssue{
				Type:     "error",
				Category: "file",
				Message:  fmt.Sprintf("Required file missing: %s", requiredFile),
				File:     fullPath,
			})
			result.Success = false
		}
	}

	// Check if there are any errors
	for _, issue := range result.Issues {
		if issue.Type == "error" {
			result.Success = false
			break
		}
	}

	if options.Verbose {
		fmt.Printf("📊 Validation completed: %d files, %d issues\n", result.FileCount, len(result.Issues))
	}

	return result, nil
}

// validateHTMLFile validates an individual HTML file
func validateHTMLFile(filePath string, options ValidationOptions, result *ValidationResult) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	htmlContent := string(content)

	// Check for broken internal links
	if options.CheckLinks {
		validateLinks(filePath, htmlContent, result)
	}

	// Check for missing images
	if options.CheckImages {
		validateImages(filePath, htmlContent, result)
	}

	return nil
}

// validateLinks checks for broken internal links in HTML content
func validateLinks(filePath, content string, result *ValidationResult) {
	// Regex to find href attributes
	linkRegex := regexp.MustCompile(`href=["']([^"']+)["']`)
	matches := linkRegex.FindAllStringSubmatch(content, -1)

	baseDir := filepath.Dir(filePath)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		href := match[1]

		// Skip external links, mailto, tel, and anchors
		if strings.HasPrefix(href, "http://") ||
			strings.HasPrefix(href, "https://") ||
			strings.HasPrefix(href, "mailto:") ||
			strings.HasPrefix(href, "tel:") ||
			strings.HasPrefix(href, "#") {
			continue
		}

		// Parse URL to handle query parameters and fragments
		parsedURL, err := url.Parse(href)
		if err != nil {
			continue
		}

		// Get the path part without query/fragment
		linkPath := parsedURL.Path
		if linkPath == "" {
			continue
		}

		// Resolve relative path
		var targetPath string
		if filepath.IsAbs(linkPath) {
			// Absolute path relative to site root
			siteRoot := filepath.Dir(filepath.Dir(filePath))   // Assuming filePath is like "site/page.html"
			targetPath = filepath.Join(siteRoot, linkPath[1:]) // Remove leading slash
		} else {
			// Relative path
			targetPath = filepath.Join(baseDir, linkPath)
		}

		// Check if target exists
		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			// Also try with .html extension
			if !strings.HasSuffix(targetPath, ".html") {
				if _, err := os.Stat(targetPath + ".html"); err == nil {
					continue // Found with .html extension
				}
			}

			result.Issues = append(result.Issues, ValidationIssue{
				Type:     "warning",
				Category: "link",
				Message:  fmt.Sprintf("Broken internal link: %s -> %s", href, targetPath),
				File:     filePath,
			})
		}
	}
}

// validateImages checks for missing images in HTML content
func validateImages(filePath, content string, result *ValidationResult) {
	// Regex to find src attributes in img tags
	imgRegex := regexp.MustCompile(`<img[^>]+src=["']([^"']+)["']`)
	matches := imgRegex.FindAllStringSubmatch(content, -1)

	baseDir := filepath.Dir(filePath)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		src := match[1]

		// Skip external images and data URLs
		if strings.HasPrefix(src, "http://") ||
			strings.HasPrefix(src, "https://") ||
			strings.HasPrefix(src, "data:") {
			continue
		}

		// Resolve relative path
		var targetPath string
		if filepath.IsAbs(src) {
			// Absolute path relative to site root
			siteRoot := filepath.Dir(filepath.Dir(filePath)) // Assuming filePath is like "site/page.html"
			targetPath = filepath.Join(siteRoot, src[1:])    // Remove leading slash
		} else {
			// Relative path
			targetPath = filepath.Join(baseDir, src)
		}

		// Check if image exists
		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			result.Issues = append(result.Issues, ValidationIssue{
				Type:     "warning",
				Category: "image",
				Message:  fmt.Sprintf("Missing image: %s -> %s", src, targetPath),
				File:     filePath,
			})
		}
	}
}

// GetDefaultValidationOptions returns recommended validation options
func GetDefaultValidationOptions() ValidationOptions {
	return ValidationOptions{
		CheckLinks:    true,
		CheckImages:   true,
		CheckFileSize: true,
		MaxFileSize:   10 * 1024 * 1024, // 10MB
		RequiredFiles: []string{
			"index.html",
			"style.css",
		},
		Verbose: false,
	}
}
