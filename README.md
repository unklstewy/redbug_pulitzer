# REDBUG_PULITZER - Report Generation System

A tool for generating detailed HTML reports and API documentation for radio protocol analysis results.

## Features

- HTML report generation for protocol analysis
- API documentation for radio commands
- Interactive visualization of radio communication
- Index page for browsing multiple reports

## Usage

```
# API for report generation
import "github.com/unklstewy/redbug_pulitzer/pkg/reporting"

// Generate API documentation
reporting.GenerateAPIDocHTML(apiDocs, "output.html", reporting.ReadMode)

// Update index page
reporting.UpdateIndexPage()
```
