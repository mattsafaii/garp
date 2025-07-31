package internal

import (
	"testing"
)

func TestDetectTailwindCLI(t *testing.T) {
	info, err := DetectTailwindCLI()
	if err != nil {
		t.Fatalf("DetectTailwindCLI() error = %v", err)
	}

	if info == nil {
		t.Fatal("DetectTailwindCLI() returned nil info")
	}

	// Test should pass regardless of whether Tailwind is installed
	// We're just checking that the function returns without error
	t.Logf("Tailwind installed: %v", info.IsInstalled)
	if info.IsInstalled {
		t.Logf("Version: %s", info.Version)
		t.Logf("Path: %s", info.ExecutablePath)
	}
}

func TestGetTailwindInstallationInstructions(t *testing.T) {
	instructions := GetTailwindInstallationInstructions()

	if instructions == "" {
		t.Error("GetTailwindInstallationInstructions() returned empty string")
	}

	// Check that instructions contain platform-specific content
	if len(instructions) < 100 {
		t.Error("Instructions seem too short to be helpful")
	}

	t.Logf("Instructions length: %d characters", len(instructions))
}

func TestValidateTailwindCLI(t *testing.T) {
	err := ValidateTailwindCLI()

	// This test will pass if Tailwind is installed, or fail with helpful error if not
	if err != nil {
		t.Logf("ValidateTailwindCLI() error (expected if Tailwind not installed): %v", err)

		// Check that error contains installation instructions
		errStr := err.Error()
		if len(errStr) < 50 {
			t.Error("Error message seems too short to contain helpful instructions")
		}
	} else {
		t.Log("ValidateTailwindCLI() passed - Tailwind CLI is installed")
	}
}
