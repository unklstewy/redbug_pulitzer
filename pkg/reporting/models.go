package reporting

import (
	"time"
)

// IndexReport represents the data structure for the index page
type IndexReport struct {
	// Fields used in api.go
	Path        string
	Size        int64
	Timestamp   string
	Title       string
	Type        string
	Description string

	// Fields used in index.go
	GeneratedAt      string
	LastUpdated      time.Time
	VendorReports    map[string]map[string][]ReportInfo
	TotalReportCount int
	VendorCount      int
	ModelCount       int
}

// ReportInfo holds metadata about a single report
type ReportInfo struct {
	Path    string
	Title   string
	Type    string
	Size    int64
	ModTime time.Time
}
