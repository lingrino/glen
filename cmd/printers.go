package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
)

// print outputs a map of environment variables in the specified format.
func print(m map[string]string, format string) {
	switch format {
	case "export":
		printExport(m)
	case "json":
		printJSON(m)
	case "table":
		printTable(m)
	default:
		log.Fatalf("ERROR: Output type '%s' is not supported", format)
	}
}

// printExport outputs a map of environment variables in 'export' format,
// meaning the output can be immediately evaluated to export the variables.
func printExport(m map[string]string) {
	for k, v := range m {
		fmt.Printf("export %s=\"%s\"\n", k, v)
	}
}

// printJSON outputs a map of environment variables in JSON format.
func printJSON(m map[string]string) {
	json, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		log.Fatalf("ERROR: Failed to marshal the output into JSON")
	}
	fmt.Println(string(json))
}

// printTable outputs a map of environment variables in a table format.
func printTable(m map[string]string) {
	data := [][]string{}
	for k, v := range m {
		data = append(data, []string{k, v})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("|")
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetHeader([]string{"Key", "Value"})
	table.AppendBulk(data)
	table.Render()
}
