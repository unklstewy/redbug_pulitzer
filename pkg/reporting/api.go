// filepath: /home/sannis/REDBUG/redbug_pulitzer/pkg/reporting/api.go
package reporting

import (
    "fmt"
    "html/template"
    "log"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"

)

// Mode determines which color scheme and labels to use
type Mode string

const (
    ReadMode  Mode = "read"
    WriteMode Mode = "write"
)

// StyleConfig holds styling information for different modes
type StyleConfig struct {
    PrimaryColor   string
    SecondaryColor string
    Title          string
    Icon           string
    HeaderBgColor  string
    BorderColor    string
}

// GetStyleConfig returns appropriate styling based on mode
func GetStyleConfig(mode Mode) StyleConfig {
    if mode == ReadMode {
        return StyleConfig{
            PrimaryColor:   "#2980b9",
            SecondaryColor: "#3498db",
            Title:          "Read Protocol API Documentation",
            Icon:           "📥",
            HeaderBgColor:  "#eaf2f8",
            BorderColor:    "#3498db",
        }
    }
    return StyleConfig{
        PrimaryColor:   "#c0392b",
        SecondaryColor: "#e74c3c",
        Title:          "Write Protocol API Documentation",
        Icon:           "📤",
        HeaderBgColor:  "#f9ebea",
        BorderColor:    "#e74c3c",
    }
}
func GenerateAPIDocHTML(apiDocs []CommandAPI, filename string, mode Mode, vendor string, model string) {
	// Get styling configuration based on mode
	style := GetStyleConfig(mode)

	// Determine report type based on mode
	reportType := ReportTypeReadAPI
	if mode == WriteMode {
		reportType = ReportTypeWriteAPI
	}

	// Generate the proper path for the report file
	filepath := GetReportPath(vendor, model, reportType, filename)

	// Create template functions map
	funcMap := template.FuncMap{
		"successRateClass": func(rate string) string {
			// Parse the percentage value
			rate = strings.TrimSuffix(rate, "%")
			value, err := strconv.ParseFloat(rate, 64)
			if err != nil {
				return ""
			}

			if value >= 95 {
				return "success-high"
			} else if value >= 80 {
				return "success-medium"
			} else {
				return "success-low"
			}
		},
	}

	// HTML template with collapsible sections and search functionality
	htmlTemplate := `
<!DOCTYPE html>
<html>
<head>
    <title>DM-32UV {{.Style.Title}}</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 0; background-color: #f5f5f5; color: #333; }
        .container { max-width: 1200px; margin: 0 auto; background-color: white; padding: 20px; box-shadow: 0 0 10px rgba(0,0,0,0.1); }
        h1 { color: {{.Style.PrimaryColor}}; text-align: center; border-bottom: 3px solid {{.Style.SecondaryColor}}; padding-bottom: 10px; }
        h2 { color: {{.Style.PrimaryColor}}; margin-top: 30px; }
        
        /* Command Card Styling */
        .command-card { 
            border: 1px solid #ddd; 
            border-radius: 8px; 
            margin-bottom: 15px; 
            overflow: hidden;
            box-shadow: 0 2px 4px rgba(0,0,0,0.05);
        }
        .command-header { 
            background-color: #f8f8f8; 
            padding: 12px; 
            cursor: pointer;
            display: flex;
            justify-content: space-between;
            align-items: center;
            border-bottom: 1px solid #eee;
        }
        .command-name { font-size: 18px; font-weight: bold; color: {{.Style.PrimaryColor}}; }
        .command-meta { display: flex; gap: 15px; font-size: 13px; color: #666; }
        .command-body { padding: 0 15px 15px; }
        
        /* Details & Summary Styling */
        details { margin: 10px 0; }
        details summary { 
            cursor: pointer; 
            padding: 8px; 
            background-color: #f9f9f9; 
            border-radius: 4px;
            font-weight: bold;
        }
        details[open] summary { margin-bottom: 10px; }
        
        /* Response Styling */
        .response-container { 
            background-color: {{.Style.HeaderBgColor}}; 
            padding: 10px; 
            border-left: 4px solid {{.Style.BorderColor}}; 
            margin: 10px 0; 
            border-radius: 4px;
        }
        
        /* Data Formatting */
        .hex { 
            font-family: 'Courier New', monospace; 
            background-color: #f8f9fa; 
            padding: 4px 8px; 
            border-radius: 3px; 
            overflow-wrap: break-word;
            word-break: break-all;
        }
        .ascii { 
            font-family: 'Courier New', monospace; 
            color: #d35400; 
            overflow-wrap: break-word;
            word-break: break-all;
        }
        .description { color: #555; font-style: italic; }
        .field-label { font-weight: bold; color: #555; display: inline-block; width: 120px; }
        
        /* Stats and Metadata */
        .timing { color: #8e44ad; font-weight: bold; }
        .category { 
            display: inline-block; 
            background-color: #f8f9fa; 
            padding: 2px 6px; 
            border-radius: 3px; 
            font-size: 12px; 
            color: #7f8c8d; 
        }
        .success-rate { font-weight: bold; }
        .success-high { color: #27ae60; }
        .success-medium { color: #f39c12; }
        .success-low { color: #c0392b; }
        
        /* Search and Filter */
        .search-container {
            background-color: {{.Style.HeaderBgColor}};
            padding: 15px;
            border-radius: 5px;
            margin: 20px 0;
            position: sticky;
            top: 0;
            z-index: 100;
        }
        .search-input {
            padding: 10px;
            width: 70%;
            border: 1px solid #ddd;
            border-radius: 4px;
            font-size: 16px;
        }
        .search-type {
            padding: 10px;
            border: 1px solid #ddd;
            border-radius: 4px;
            margin-left: 10px;
        }
        .index-link {
            display: inline-block;
            margin-top: 10px;
            color: {{.Style.PrimaryColor}};
            text-decoration: none;
            font-weight: bold;
        }
        .index-link:hover {
            text-decoration: underline;
        }
        #command-count {
            margin-left: 10px;
            font-weight: bold;
        }
        .no-results {
            padding: 20px;
            text-align: center;
            font-style: italic;
            color: #555;
            display: none;
        }
        
        /* Quick Navigation */
        .quick-nav {
            position: fixed;
            top: 100px;
            right: 20px;
            background: white;
            border: 1px solid #ddd;
            border-radius: 8px;
            padding: 10px;
            max-width: 200px;
            max-height: 400px;
            overflow-y: auto;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .quick-nav h4 {
            margin-top: 0;
            border-bottom: 1px solid #eee;
            padding-bottom: 5px;
        }
        .quick-nav ul {
            list-style: none;
            padding: 0;
            margin: 0;
        }
        .quick-nav a {
            display: block;
            padding: 3px 0;
            text-decoration: none;
            color: {{.Style.PrimaryColor}};
            font-size: 13px;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
        }
        .quick-nav a:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>{{.Style.Icon}} {{.Style.Title}}</h1>
        
        <div class="search-container">
            <input type="text" id="search-input" class="search-input" placeholder="Search commands...">
            <select id="search-type" class="search-type">
                <option value="all">All Fields</option>
                <option value="command">Command Name</option>
                <option value="hex">HEX Value</option>
                <option value="ascii">ASCII Value</option>
                <option value="category">Category</option>
            </select>
            <span id="command-count">Showing {{len .ApiDocs}} commands</span>
            <br>
            <a href="/index.html" class="index-link">← Back to Index</a>
        </div>
        
        <div class="quick-nav" id="quick-nav">
            <h4>Quick Navigation</h4>
            <ul id="nav-list">
                {{range $index, $cmd := .ApiDocs}}
                <li><a href="#cmd-{{$index}}">{{$cmd.Command}}</a></li>
                {{end}}
            </ul>
        </div>

        <div id="command-list">
            {{range $index, $cmd := .ApiDocs}}
            <div class="command-card" id="cmd-{{$index}}" data-command="{{$cmd.Command}}" data-hex="{{$cmd.HexValue}}" data-ascii="{{$cmd.ASCIIValue}}" data-category="{{$cmd.DataCategory}}">
                <div class="command-header" onclick="toggleCommandDetails(this)">
                    <div class="command-name">{{$cmd.Command}}</div>
                    <div class="command-meta">
                        {{if $cmd.DataCategory}}
                        <span class="category">{{$cmd.DataCategory}}</span>
                        {{end}}
                        {{if $cmd.SuccessRate}}
                        <span class="success-rate {{successRateClass $cmd.SuccessRate}}">{{$cmd.SuccessRate}}</span>
                        {{end}}
                        <span class="timing">{{$cmd.TimingAverage}}</span>
                        <span class="frequency">{{$cmd.FrequencyCount}}×</span>
                    </div>
                </div>
                
                <div class="command-body" style="display:none;">
                    <details>
                        <summary>Command Details</summary>
                        <div>
                            <p><span class="field-label">Description:</span> <span class="description">{{$cmd.Description}}</span></p>
                            <p><span class="field-label">HEX:</span> <span class="hex">{{$cmd.HexValue}}</span></p>
                            <p><span class="field-label">ASCII:</span> <span class="ascii">{{$cmd.ASCIIValue}}</span></p>
                        </div>
                    </details>
                    
                    <details>
                        <summary>Response Details</summary>
                        <div class="response-container">
                            <p><strong>{{$cmd.ResponseType}}</strong></p>
                            <p><span class="field-label">HEX:</span> <span class="hex">{{$cmd.ResponseHex}}</span></p>
                            <p><span class="field-label">ASCII:</span> <span class="ascii">{{$cmd.ResponseASCII}}</span></p>
                        </div>
                    </details>
                    
                    <details>
                        <summary>Performance Metrics</summary>
                        <p><span class="field-label">Average Time:</span> <span class="timing">{{$cmd.TimingAverage}}</span></p>
                        <p><span class="field-label">Usage Count:</span> <span class="frequency">{{$cmd.FrequencyCount}} times</span></p>
                        {{if $cmd.SuccessRate}}
                        <p><span class="field-label">Success Rate:</span> <span class="success-rate {{successRateClass $cmd.SuccessRate}}">{{$cmd.SuccessRate}}</span></p>
                        {{end}}
                    </details>
                </div>
            </div>
            {{end}}
            
            <div class="no-results" id="no-results">
                No commands match your search.
            </div>
        </div>
        
        <div style="margin-top: 50px; text-align: center; color: #7f8c8d; border-top: 1px solid #bdc3c7; padding-top: 20px;">
            <p>Generated by DM-32UV Protocol Analyzer | {{.Style.Title}} | Generated: {{.GeneratedAt}}</p>
        </div>
    </div>

    <script>
        // Toggle command details when clicking on header
        function toggleCommandDetails(header) {
            const body = header.nextElementSibling;
            body.style.display = body.style.display === 'none' ? 'block' : 'none';
        }
        
        // Search functionality
        const searchInput = document.getElementById('search-input');
        const searchType = document.getElementById('search-type');
        const commandList = document.getElementById('command-list');
        const commandCount = document.getElementById('command-count');
        const noResults = document.getElementById('no-results');
        
        function performSearch() {
            const searchTerm = searchInput.value.toLowerCase();
            const searchField = searchType.value;
            let visibleCount = 0;
            
            const commandCards = document.querySelectorAll('.command-card');
            commandCards.forEach(card => {
                let match = false;
                
                if (searchTerm === '') {
                    match = true;
                } else {
                    switch(searchField) {
                        case 'command':
                            match = card.dataset.command.toLowerCase().includes(searchTerm);
                            break;
                        case 'hex':
                            match = card.dataset.hex.toLowerCase().includes(searchTerm);
                            break;
                        case 'ascii':
                            match = card.dataset.ascii.toLowerCase().includes(searchTerm);
                            break;
                        case 'category':
                            match = card.dataset.category && card.dataset.category.toLowerCase().includes(searchTerm);
                            break;
                        case 'all':
                        default:
                            match = card.dataset.command.toLowerCase().includes(searchTerm) ||
                                   card.dataset.hex.toLowerCase().includes(searchTerm) ||
                                   card.dataset.ascii.toLowerCase().includes(searchTerm) ||
                                   (card.dataset.category && card.dataset.category.toLowerCase().includes(searchTerm));
                            break;
                    }
                }
                
                card.style.display = match ? 'block' : 'none';
                if (match) visibleCount++;
            });
            
            commandCount.textContent = 'Showing ' + visibleCount + ' of {{len .ApiDocs}} commands';
            noResults.style.display = visibleCount === 0 ? 'block' : 'none';
            
            // Update quick nav
            updateQuickNav();
        }
        
        // Update quick navigation based on visible commands
        function updateQuickNav() {
            const navList = document.getElementById('nav-list');
            navList.innerHTML = '';
            
            const visibleCommands = document.querySelectorAll('.command-card[style="display: block"]');
            visibleCommands.forEach(cmd => {
                const id = cmd.id;
                const name = cmd.querySelector('.command-name').textContent;
                
                const li = document.createElement('li');
                const a = document.createElement('a');
                a.href = '#' + id;
                a.textContent = name;
                li.appendChild(a);
                navList.appendChild(li);
            });
        }
        
        searchInput.addEventListener('input', performSearch);
        searchType.addEventListener('change', performSearch);
        
        // Initialize on page load
        document.addEventListener('DOMContentLoaded', function() {
            // Open the first command by default
            const firstCommand = document.querySelector('.command-card');
            if (firstCommand) {
                const header = firstCommand.querySelector('.command-header');
                const body = header.nextElementSibling;
                body.style.display = 'block';
            }
        });
    </script>
</body>
</html>
`

	// Prepare template data
	data := struct {
		ApiDocs     []CommandAPI
		Style       StyleConfig
		GeneratedAt string
	}{
		ApiDocs:     apiDocs,
		Style:       style,
		GeneratedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	// Create and parse template
	tmpl, err := template.New("api").Funcs(funcMap).Parse(htmlTemplate)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	// Create output file
	file, err := os.Create(filepath)
	if err != nil {
		log.Fatalf("Error creating HTML file: %v", err)
	}
	defer file.Close()

	// Execute template
	err = tmpl.Execute(file, data)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	modeStr := "Read"
	if mode == WriteMode {
		modeStr = "Write"
	}
	fmt.Printf("%s API documentation saved to: %s\n", modeStr, filepath)

	// Update index page
	UpdateIndexPage()
}

// UpdateIndexPage updates the central index page with links to all reports

// Helper function to sort reports by timestamp (newest first)

// ReportInfo holds metadata about a report file
// type ReportInfo struct {
// 	Path    string
// 	Size    int64
// 	ModTime time.Time
// }

// findReports searches for all report files in common locations
func findReports() []ReportInfo {
	var reports []ReportInfo

	// New approach: recursively walk through the reports directory
	reportsDir := "reports"

	err := filepath.Walk(reportsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only include HTML files
		if strings.HasSuffix(path, ".html") {
			reports = append(reports, ReportInfo{
				Path:    path,
				Size:    info.Size(),
				ModTime: info.ModTime(),
			})
		}

		return nil
	})

	if err != nil {
		log.Printf("Error walking reports directory: %v", err)
	}

	return reports
}

// Add this function to determine the appropriate report path
func GetReportPath(vendor, model string, reportType string, fileName string) string {
	// Base reports directory
	baseDir := "reports"

	// Determine the subdirectory based on report type
	var subDir string
	if strings.Contains(reportType, "read") {
		subDir = filepath.Join("protocol", "read", vendor, model)
	} else if strings.Contains(reportType, "write") {
		subDir = filepath.Join("protocol", "write", vendor, model)
	} else if strings.Contains(reportType, "api") {
		subDir = filepath.Join("api", vendor, model)
	} else if strings.Contains(reportType, "codeplug") {
		subDir = filepath.Join("codeplug", vendor, model)
	} else if strings.Contains(reportType, "firmware") {
		subDir = filepath.Join("firmware", vendor, model)
	} else if strings.Contains(reportType, "cps") {
		subDir = filepath.Join("cps", vendor, model)
	} else {
		// Default to protocol for unrecognized types
		subDir = filepath.Join("protocol", "other")
	}

	// Create directory if it doesn't exist
	fullDir := filepath.Join(baseDir, subDir)
	os.MkdirAll(fullDir, 0755)

	return filepath.Join(fullDir, fileName)
}

// Add a constants for report types
const (
	ReportTypeReadAnalysis  = "read_analysis"
	ReportTypeWriteAnalysis = "write_analysis"
	ReportTypeReadAPI       = "read_api"
	ReportTypeWriteAPI      = "write_api"
	ReportTypeCodeplug      = "codeplug"
	ReportTypeFirmware      = "firmware"
	ReportTypeCPS           = "cps"
)
