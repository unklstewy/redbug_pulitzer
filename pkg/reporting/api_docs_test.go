package reporting

import (
	"os"
	"strings"
	"testing"

	reporttest "github.com/unklstewy/redbug_pulitzer/pkg/reporting/testing"
)

func TestGenerateAPIDocHTML(t *testing.T) {
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

	// Create sample API docs
	apiDocs := []CommandAPI{
		{
			Command:        "READ_CONFIG",
			HexValue:       "5200",
			ASCIIValue:     "R.",
			Description:    "Read radio configuration",
			ResponseType:   "ACK followed by data",
			ResponseHex:    "0644415441",
			ResponseASCII:  ".DATA",
			FrequencyCount: 12,
			TimingAverage:  "125ms",
		},
		{
			Command:        "WRITE_CONFIG",
			HexValue:       "5700",
			ASCIIValue:     "W.",
			Description:    "Write radio configuration",
			ResponseType:   "ACK or NAK",
			ResponseHex:    "06",
			ResponseASCII:  ".",
			FrequencyCount: 8,
			TimingAverage:  "230ms",
		},
	}

	// Test both read and write modes
	modes := []struct {
		mode     Mode
		filename string
	}{
		{ReadMode, "baofeng_dm32uv_read_api_docs.html"},
		{WriteMode, "baofeng_dm32uv_write_api_docs.html"},
	}

	for _, m := range modes {
		t.Run(string(m.mode), func(t *testing.T) {
			// Generate the API documentation
			GenerateAPIDocHTML(apiDocs, m.filename, m.mode, "baofeng", "dm32uv")

			// Determine the expected path
			var reportType string
			if m.mode == ReadMode {
				reportType = ReportTypeReadAPI
			} else {
				reportType = ReportTypeWriteAPI
			}

			expectedPath := GetReportPath("baofeng", "dm32uv", reportType, m.filename)

			// Check if the file was created
			if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
				t.Errorf("Expected file %s to be created, but it doesn't exist", expectedPath)
			}

			// Read the content and check if it contains expected elements
			content, err := os.ReadFile(expectedPath)
			if err != nil {
				t.Fatalf("Failed to read generated file: %v", err)
			}

			contentStr := string(content)
			expectedPhrases := []string{
				"READ_CONFIG", "WRITE_CONFIG",
				"Read radio configuration", "Write radio configuration",
				"5200", "5700",
				"baofeng", "dm32uv",
			}

			for _, phrase := range expectedPhrases {
				if !strings.Contains(contentStr, phrase) {
					t.Errorf("Expected generated HTML to contain '%s', but it doesn't", phrase)
				}
			}

			// Also check mode-specific content
			modeStr := "Read Mode"
			if m.mode == WriteMode {
				modeStr = "Write Mode"
			}

			if !strings.Contains(contentStr, modeStr) {
				t.Errorf("Expected generated HTML to contain '%s', but it doesn't", modeStr)
			}
		})
	}
}
