package reporting

import (
	"os"
	"path/filepath"
	"testing"

	reporttest "github.com/unklstewy/redbug_pulitzer/pkg/reporting/testing"
)

func TestGetReportPath(t *testing.T) {
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

	tests := []struct {
		name       string
		vendor     string
		model      string
		reportType string
		fileName   string
		expected   string
	}{
		{
			name:       "Read analysis report",
			vendor:     "baofeng",
			model:      "dm32uv",
			reportType: ReportTypeReadAnalysis,
			fileName:   "dm32uv_read_analysis.html",
			expected:   filepath.Join("reports", "protocol", "read", "baofeng", "dm32uv", "dm32uv_read_analysis.html"),
		},
		{
			name:       "Write analysis report",
			vendor:     "baofeng",
			model:      "dm32uv",
			reportType: ReportTypeWriteAnalysis,
			fileName:   "dm32uv_write_analysis.html",
			expected:   filepath.Join("reports", "protocol", "write", "baofeng", "dm32uv", "dm32uv_write_analysis.html"),
		},
		{
			name:       "API documentation",
			vendor:     "baofeng",
			model:      "dm32uv",
			reportType: ReportTypeReadAPI,
			fileName:   "dm32uv_read_api_docs.html",
			expected:   filepath.Join("reports", "api", "baofeng", "dm32uv", "dm32uv_read_api_docs.html"),
		},
		{
			name:       "Codeplug analysis",
			vendor:     "baofeng",
			model:      "dm32uv",
			reportType: ReportTypeCodeplug,
			fileName:   "dm32uv_codeplug_analysis.html",
			expected:   filepath.Join("reports", "codeplug", "baofeng", "dm32uv", "dm32uv_codeplug_analysis.html"),
		},
		{
			name:       "Firmware analysis",
			vendor:     "baofeng",
			model:      "dm32uv",
			reportType: ReportTypeFirmware,
			fileName:   "dm32uv_firmware_analysis.html",
			expected:   filepath.Join("reports", "firmware", "baofeng", "dm32uv", "dm32uv_firmware_analysis.html"),
		},
		{
			name:       "CPS analysis",
			vendor:     "baofeng",
			model:      "dm32uv",
			reportType: ReportTypeCPS,
			fileName:   "dm32uv_cps_analysis.html",
			expected:   filepath.Join("reports", "cps", "baofeng", "dm32uv", "dm32uv_cps_analysis.html"),
		},
		{
			name:       "Unknown report type",
			vendor:     "baofeng",
			model:      "dm32uv",
			reportType: "unknown",
			fileName:   "unknown_report.html",
			expected:   filepath.Join("reports", "protocol", "other", "unknown_report.html"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetReportPath(tt.vendor, tt.model, tt.reportType, tt.fileName)

			// Use filepath.Clean to normalize path separators for cross-platform testing
			if filepath.Clean(result) != filepath.Clean(tt.expected) {
				t.Errorf("GetReportPath() = %v, want %v", result, tt.expected)
			}

			// Verify the directory was created
			dir := filepath.Dir(result)
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				t.Errorf("Expected directory %s to be created, but it doesn't exist", dir)
			}
		})
	}
}
