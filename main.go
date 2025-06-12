package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/unklstewy/redbug/pkg/cli"
	"github.com/unklstewy/redbug_pulitzer/pkg/reporting"
)

var pulitzerCommand = cli.CommandInfo{
	Name:        "report",
	Description: "Report Generation System",
	Usage:       "report <command> [flags] [args]",
	Examples: []string{
		"report generate-index",
		"report view baofeng dm32uv",
		"report export baofeng dm32uv --format=pdf",
	},
	SubCommands: []cli.CommandInfo{
		{
			Name:        "generate-index",
			Description: "Generate or update the main index page",
			Usage:       "report generate-index",
		},
		{
			Name:        "view",
			Description: "View reports for a specific radio",
			Usage:       "report view <vendor> <model>",
		},
		{
			Name:        "export",
			Description: "Export reports to different formats",
			Usage:       "report export <vendor> <model> --format=<format>",
		},
	},
}

func main() {
	// Check for help flag
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		cli.PrintHelp(pulitzerCommand)
		return
	}

	if len(os.Args) < 2 {
		fmt.Println("Error: Command required")
		cli.PrintHelp(pulitzerCommand)
		os.Exit(1)
	}

	command := os.Args[1]

	// Find the matching subcommand
	var cmdInfo *cli.CommandInfo
	for _, cmd := range pulitzerCommand.SubCommands {
		if cmd.Name == command {
			cmdInfo = &cmd
			break
		}
	}

	if cmdInfo == nil {
		fmt.Printf("Unknown command: %s\n\n", command)
		cli.PrintHelp(pulitzerCommand)
		os.Exit(1)
	}

	// If help is requested for the subcommand
	if len(os.Args) > 2 && (os.Args[2] == "--help" || os.Args[2] == "-h") {
		cli.PrintHelp(*cmdInfo)
		os.Exit(0)
	}

	switch command {
	case "generate-index":
		reporting.UpdateIndexPage()
		fmt.Println("Index page updated successfully")

	case "view":
		if len(os.Args) < 4 {
			fmt.Println("Error: Vendor and model required")
			cli.PrintHelp(*cmdInfo)
			os.Exit(1)
		}

		vendor := os.Args[2]
		model := os.Args[3]

		// Check if reports exist
		reportPath := filepath.Join("reports", "protocol", "read", vendor, model)
		if _, err := os.Stat(reportPath); os.IsNotExist(err) {
			fmt.Printf("No reports found for %s %s\n", vendor, model)
			os.Exit(1)
		}

		fmt.Printf("Opening reports for %s %s...\n", vendor, model)
		// Here you would open the index.html in a browser

	case "export":
		if len(os.Args) < 4 {
			fmt.Println("Error: Vendor and model required")
			cli.PrintHelp(*cmdInfo)
			os.Exit(1)
		}

		vendor := os.Args[2]
		model := os.Args[3]

		// Default format
		format := "pdf"

		// Check for format flag
		for i := 4; i < len(os.Args); i++ {
			if strings.HasPrefix(os.Args[i], "--format=") {
				format = strings.TrimPrefix(os.Args[i], "--format=")
				break
			}
		}

		fmt.Printf("Exporting %s %s reports to %s format...\n", vendor, model, format)
		// Here you would implement the export functionality

	default:
		fmt.Printf("Unknown command: %s\n", command)
		cli.PrintHelp(pulitzerCommand)
		os.Exit(1)
	}
}
