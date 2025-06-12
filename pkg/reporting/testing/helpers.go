package testing

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// GetTestDataPath returns the absolute path to the testdata directory
func GetTestDataPath() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "testdata")
}

// CreateTempReportDir creates a temporary directory for test reports
func CreateTempReportDir(t *testing.T) string {
	t.Helper()

	tempDir, err := ioutil.TempDir("", "redbug_pulitzer_test_")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Create reports subdirectories
	subDirs := []string{
		"protocol/read/baofeng/dm32uv",
		"protocol/write/baofeng/dm32uv",
		"api/baofeng/dm32uv",
	}

	for _, dir := range subDirs {
		err := os.MkdirAll(filepath.Join(tempDir, "reports", dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create subdirectory: %v", err)
		}
	}

	return tempDir
}

// CleanupTempDir removes a temporary test directory
func CleanupTempDir(t *testing.T, dir string) {
	t.Helper()

	if err := os.RemoveAll(dir); err != nil {
		t.Logf("Warning: Failed to clean up temp directory %s: %v", dir, err)
	}
}
