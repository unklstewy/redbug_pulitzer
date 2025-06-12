package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FormatFileSize returns a human-readable file size
func FormatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// FormatTime formats a time.Time as a human-readable string
func FormatTime(t time.Time) string {
	return t.Format("Jan 02, 2006 15:04:05")
}

// GetRelativePath returns a path relative to the current directory
func GetRelativePath(path string) string {
	pwd, err := os.Getwd()
	if err != nil {
		return path
	}
	rel, err := filepath.Rel(pwd, path)
	if err != nil {
		return path
	}
	return rel
}

// EnsureDirectoryExists creates a directory if it doesn't exist
func EnsureDirectoryExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

// GetFileExtension returns the extension of a file
func GetFileExtension(path string) string {
	return strings.ToLower(filepath.Ext(path))
}
