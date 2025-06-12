package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestMainIntegration(t *testing.T) {
	// Skip in regular test runs
	if os.Getenv("RUN_INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration test. Set RUN_INTEGRATION_TESTS=1 to enable.")
	}

	// Create a temporary test directory
	tempDir, err := os.MkdirTemp("", "redbug_pulitzer_integration_")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set up test directory structure
	testDirs := []string{
		"reports/protocol/read/baofeng/dm32uv",
		"reports/protocol/write/baofeng/dm32uv",
		"reports/api/baofeng/dm32uv",
	}

	for _, dir := range testDirs {
		if err := os.MkdirAll(filepath.Join(tempDir, dir), 0755); err != nil {
			t.Fatalf("Failed to create directory structure: %v", err)
		}
	}

	// Create some sample HTML files
	sampleFiles := []struct {
		path    string
		content string
	}{
		{
			path: filepath.Join(tempDir, "reports/protocol/read/baofeng/dm32uv/test_read.html"),
			content: `<html><head><title>Test Read</title></head><body>
            <h1>Test Read Report</h1><p>This is a test file</p></body></html>`,
		},
		{
			path: filepath.Join(tempDir, "reports/protocol/write/baofeng/dm32uv/test_write.html"),
			content: `<html><head><title>Test Write</title></head><body>
            <h1>Test Write Report</h1><p>This is a test file</p></body></html>`,
		},
	}

	for _, file := range sampleFiles {
		if err := os.WriteFile(file.path, []byte(file.content), 0644); err != nil {
			t.Fatalf("Failed to write file %s: %v", file.path, err)
		}
	}

	// Save current directory
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Change to test directory
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to test directory: %v", err)
	}
	defer os.Chdir(currentDir)

	// Run the main program with generate-index command
	cmd := exec.Command("go", "run", filepath.Join(currentDir, "main.go"), "generate-index")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	// Check if index.html was created
	indexPath := filepath.Join(tempDir, "reports/index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Errorf("Expected index.html to be created, but it doesn't exist")
	}
}
