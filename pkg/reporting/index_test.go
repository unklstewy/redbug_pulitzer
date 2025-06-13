package reporting

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	reporttest "github.com/unklstewy/redbug_pulitzer/pkg/reporting/testing"
)

func TestUpdateIndexPage(t *testing.T) {
	// Create a temporary test directory
	tempDir := reporttest.CreateTempReportDir(t)
	defer reporttest.CleanupTempDir(t, tempDir)

	// Save current working directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// Change to test directory for the test
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}
	// Restore original directory after test
	defer os.Chdir(originalWd)

	// Create some sample report files
	reportFiles := []struct {
		path    string
		content string
	}{
		{
			path:    filepath.Join("reports", "protocol", "read", "baofeng", "dm32uv", "read_analysis.html"),
			content: "<html><head><title>Read Analysis</title></head><body><h1>Read Analysis</h1></body></html>",
		},
		{
			path:    filepath.Join("reports", "protocol", "write", "baofeng", "dm32uv", "write_analysis.html"),
			content: "<html><head><title>Write Analysis</title></head><body><h1>Write Analysis</h1></body></html>",
		},
		{
			path:    filepath.Join("reports", "api", "baofeng", "dm32uv", "api_docs.html"),
			content: "<html><head><title>API Documentation</title></head><body><h1>API Documentation</h1></body></html>",
		},
	}

	// Create the sample report files
	for _, file := range reportFiles {
		// Create directory
		dir := filepath.Dir(file.path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}

		// Write file
		if err := os.WriteFile(file.path, []byte(file.content), 0644); err != nil {
			t.Fatalf("Failed to write file %s: %v", file.path, err)
		}
	}

	// Generate the index page
	UpdateIndexPage()

	// Check if index file was created
	indexPath := filepath.Join("reports", "index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Errorf("Expected index file %s to be created, but it doesn't exist", indexPath)
	}

	// Read the content and check if it contains expected elements
	content, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("Failed to read generated index file: %v", err)
	}

	contentStr := string(content)
	expectedPhrases := []string{
		"REDBUG Reports",
		"baofeng",
		"dm32uv",
		"read_analysis.html",
		"write_analysis.html",
		"api_docs.html",
	}

	for _, phrase := range expectedPhrases {
		if !strings.Contains(contentStr, phrase) {
			t.Errorf("Expected index HTML to contain '%s', but it doesn't", phrase)
		}
	}
}

func TestSortReportsByTime(t *testing.T) {
	now := time.Now()
	oneHourAgo := now.Add(-1 * time.Hour)
	twoDaysAgo := now.Add(-48 * time.Hour)

	reports := []ReportInfo{
		{
			Path:    "report1.html",
			Size:    1000,
			ModTime: oneHourAgo,
		},
		{
			Path:    "report2.html",
			Size:    2000,
			ModTime: now,
		},
		{
			Path:    "report3.html",
			Size:    3000,
			ModTime: twoDaysAgo,
		},
	}

	// Sort the reports
	sortReportsByTime(reports)

	// Check the order - newest first
	if reports[0].Path != "report2.html" || reports[1].Path != "report1.html" || reports[2].Path != "report3.html" {
		t.Errorf("Reports not sorted correctly by time. Got order: %v, %v, %v",
			reports[0].Path, reports[1].Path, reports[2].Path)
	}
}
