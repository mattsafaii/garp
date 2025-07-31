package internal

import (
	"testing"
)

func TestDetectPagefind(t *testing.T) {
	info, err := DetectPagefind()
	if err != nil {
		t.Fatalf("DetectPagefind() error = %v", err)
	}

	if info == nil {
		t.Fatal("DetectPagefind() returned nil info")
	}

	// Test should pass regardless of whether Pagefind is installed
	// We're just checking that the function returns without error
	t.Logf("Pagefind installed: %v", info.IsInstalled)
	if info.IsInstalled {
		t.Logf("Version: %s", info.Version)
		t.Logf("Path: %s", info.ExecutablePath)
		t.Logf("Extended: %v", info.IsExtended)
	}
}

func TestGetPagefindInstallationInstructions(t *testing.T) {
	instructions := GetPagefindInstallationInstructions()

	if instructions == "" {
		t.Error("GetPagefindInstallationInstructions() returned empty string")
	}

	// Check that instructions contain platform-specific content
	if len(instructions) < 100 {
		t.Error("Instructions seem too short to be helpful")
	}

	// Should contain NPX method (recommended)
	if !containsString(instructions, "npx pagefind") {
		t.Error("Instructions should contain NPX installation method")
	}

	// Should contain installation URL
	if !containsString(instructions, "pagefind.app") {
		t.Error("Instructions should contain Pagefind documentation URL")
	}

	t.Logf("Instructions length: %d characters", len(instructions))
}

func TestValidatePagefind(t *testing.T) {
	err := ValidatePagefind()

	// This test will pass if Pagefind is installed, or fail with helpful error if not
	if err != nil {
		t.Logf("ValidatePagefind() error (expected if Pagefind not installed): %v", err)

		// Check that error contains installation instructions
		errStr := err.Error()
		if len(errStr) < 50 {
			t.Error("Error message seems too short to contain helpful instructions")
		}

		// Should contain NPX method in error
		if !containsString(errStr, "npx pagefind") {
			t.Error("Error message should contain NPX installation method")
		}
	} else {
		t.Log("ValidatePagefind() passed - Pagefind is installed")
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		(len(s) > len(substr) && containsString(s[1:], substr))
}

// More robust substring check
func containsStringRobust(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestPagefindDetectionMethods(t *testing.T) {
	// Test that different detection methods work independently
	info := &PagefindInfo{}

	// Test command detection (will fail if not installed, but shouldn't crash)
	t.Run("command_detection", func(t *testing.T) {
		result := checkPagefindCommand("pagefind", info)
		t.Logf("Direct command detection result: %v", result)
	})

	t.Run("npx_detection", func(t *testing.T) {
		npxInfo := &PagefindInfo{}
		result := checkNpxPagefind(npxInfo)
		t.Logf("NPX detection result: %v", result)
	})

	t.Run("python_detection", func(t *testing.T) {
		pyInfo := &PagefindInfo{}
		result := checkPythonPagefind(pyInfo)
		t.Logf("Python detection result: %v", result)
	})

	t.Run("path_detection", func(t *testing.T) {
		paths := getPagefindCommonPaths()
		t.Logf("Found %d common paths to check", len(paths))
		if len(paths) == 0 {
			t.Error("Should return at least some common paths")
		}
	})
}
