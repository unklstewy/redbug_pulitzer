package reporting

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	reportutils "github.com/unklstewy/redbug_pulitzer/pkg/reporting/utils"
)

// UpdateIndexPage generates or updates the main index page
func UpdateIndexPage() error {
	fmt.Println("Generating index page...")

	// Create a new index report
	report := IndexReport{
		GeneratedAt:   time.Now().Format("January 2, 2006 15:04:05"),
		LastUpdated:   time.Now(),
		VendorReports: make(map[string]map[string][]ReportInfo),
	}

	// Find all reports
	err := filepath.Walk("reports", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and the index file itself
		if info.IsDir() || filepath.Base(path) == "index.html" {
			return nil
		}

		// Only include HTML files
		if !strings.HasSuffix(strings.ToLower(path), ".html") {
			return nil
		}

		// Parse the path to extract vendor and model
		parts := strings.Split(path, string(os.PathSeparator))
		if len(parts) < 5 {
			// Skip files that don't follow the expected directory structure
			return nil
		}

		// Extract report type, vendor, and model
		var reportType, vendor, model string

		if parts[1] == "protocol" {
			// For protocol reports: reports/protocol/read|write/vendor/model/...
			if len(parts) < 6 {
				return nil
			}
			reportType = parts[2] // read or write
			vendor = parts[3]
			model = parts[4]
		} else if parts[1] == "api" || parts[1] == "firmware" || parts[1] == "codeplug" || parts[1] == "cps" {
			// For other reports: reports/type/vendor/model/...
			reportType = parts[1]
			vendor = parts[2]
			model = parts[3]
		} else {
			// Unknown report type
			return nil
		}

		// Create vendor and model maps if they don't exist
		if _, ok := report.VendorReports[vendor]; !ok {
			report.VendorReports[vendor] = make(map[string][]ReportInfo)
		}
		if _, ok := report.VendorReports[vendor][model]; !ok {
			report.VendorReports[vendor][model] = []ReportInfo{}
		}

		// Add report to the list
		report.VendorReports[vendor][model] = append(report.VendorReports[vendor][model], ReportInfo{
			Path:    path,
			Title:   reportTypeToTitle(reportType, filepath.Base(path)),
			Type:    reportType,
			Size:    info.Size(),
			ModTime: info.ModTime(),
		})

		report.TotalReportCount++
		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking reports directory: %w", err)
	}

	// Count vendors and models
	report.VendorCount = len(report.VendorReports)
	for _, models := range report.VendorReports {
		report.ModelCount += len(models)
	}

	// Sort reports by time (newest first)
	for vendor, models := range report.VendorReports {
		for model, reports := range models {
			sortReportsByTime(reports)
			report.VendorReports[vendor][model] = reports
		}
	}

	// Generate HTML
	html := generateIndexHTML(report)

	// Create reports directory if it doesn't exist
	reportutils.EnsureDirectoryExists("reports")

	// Write index file
	indexPath := filepath.Join("reports", "index.html")
	err = os.WriteFile(indexPath, []byte(html), 0644)
	if err != nil {
		return fmt.Errorf("error writing index file: %w", err)
	}

	fmt.Printf("Index page generated: %s\n", indexPath)
	return nil
}

// sortReportsByTime sorts reports by modification time (newest first)
func sortReportsByTime(reports []ReportInfo) {
	sort.Slice(reports, func(i, j int) bool {
		return reports[i].ModTime.After(reports[j].ModTime)
	})
}

// reportTypeToTitle converts a report type to a human-readable title
func reportTypeToTitle(reportType, filename string) string {
	switch reportType {
	case "read":
		return "Read Protocol Analysis"
	case "write":
		return "Write Protocol Analysis"
	case "api":
		return "API Documentation"
	case "firmware":
		return "Firmware Analysis"
	case "codeplug":
		return "Codeplug Analysis"
	case "cps":
		return "CPS Analysis"
	default:
		// Try to extract a title from the filename
		return strings.ReplaceAll(strings.TrimSuffix(filename, filepath.Ext(filename)), "_", " ")
	}
}

// generateIndexHTML generates the HTML for the index page
func generateIndexHTML(report IndexReport) string {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>REDBUG Reports</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
        }
        h1, h2, h3 {
            color: #2c3e50;
        }
        .header {
            background-color: #3498db;
            color: white;
            padding: 20px;
            text-align: center;
            border-radius: 5px;
            margin-bottom: 20px;
        }
        .vendor-section {
            margin-bottom: 30px;
            border: 1px solid #ddd;
            border-radius: 5px;
            padding: 15px;
        }
        .vendor-header {
            background-color: #f8f9fa;
            padding: 10px;
            margin-bottom: 15px;
            border-radius: 3px;
        }
        .model-section {
            margin-bottom: 20px;
            padding-left: 20px;
        }
        .report-list {
            list-style-type: none;
            padding-left: 0;
        }
        .report-item {
            margin-bottom: 10px;
            padding: 10px;
            background-color: #f8f9fa;
            border-radius: 3px;
        }
        .report-link {
            text-decoration: none;
            color: #3498db;
            font-weight: bold;
        }
        .report-meta {
            color: #777;
            font-size: 0.9em;
        }
        .footer {
            margin-top: 30px;
            text-align: center;
            color: #777;
            font-size: 0.9em;
        }
        .stats {
            background-color: #f8f9fa;
            padding: 10px;
            border-radius: 3px;
            margin-bottom: 20px;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>REDBUG Radio Protocol Analysis Reports</h1>
        <p>Complete repository of analysis reports and documentation</p>
    </div>
    
    <div class="stats">
        <p><strong>Total Reports:</strong> ` + fmt.Sprintf("%d", report.TotalReportCount) + `</p>
        <p><strong>Vendors:</strong> ` + fmt.Sprintf("%d", report.VendorCount) + ` | <strong>Models:</strong> ` + fmt.Sprintf("%d", report.ModelCount) + `</p>
        <p><strong>Last Updated:</strong> ` + report.GeneratedAt + `</p>
    </div>
`

	// Sort vendors alphabetically
	var vendors []string
	for vendor := range report.VendorReports {
		vendors = append(vendors, vendor)
	}
	sort.Strings(vendors)

	// Add vendor sections
	for _, vendor := range vendors {
		html += `
    <div class="vendor-section">
        <div class="vendor-header">
            <h2>` + strings.Title(vendor) + `</h2>
        </div>
`

		// Sort models alphabetically
		var models []string
		for model := range report.VendorReports[vendor] {
			models = append(models, model)
		}
		sort.Strings(models)

		// Add model sections
		for _, model := range models {
			html += `
        <div class="model-section">
            <h3>` + strings.Title(model) + `</h3>
            <ul class="report-list">
`

			// Add reports
			reports := report.VendorReports[vendor][model]
			for _, r := range reports {
				html += `
                <li class="report-item">
                    <a class="report-link" href="` + r.Path + `">` + r.Title + `</a>
                    <div class="report-meta">
                        <span>Size: ` + reportutils.FormatFileSize(r.Size) + `</span> | 
                        <span>Modified: ` + reportutils.FormatTime(r.ModTime) + `</span>
                    </div>
                </li>
`
			}

			html += `
            </ul>
        </div>
`
		}

		html += `
    </div>
`
	}

	html += `
    <div class="footer">
        <p>Generated by REDBUG Pulitzer on ` + report.GeneratedAt + `</p>
    </div>
</body>
</html>
`

	return html
}
