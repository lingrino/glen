package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/olekukonko/tablewriter"
)

// output outputs a map of environment variables in the specified format.
func output(m map[string]string, format string) {
	switch format {
	case "export":
		outputExport(m)
	case "json":
		outputJSON(m)
	case "table":
		outputTable(m)
	default:
		slog.Error("output type is not supported", "type", format)
		os.Exit(1)
	}
}

// outputExport outputs a map of environment variables in 'export' format,
// meaning the output can be immediately evaluated to export the variables.
func outputExport(m map[string]string) {
	for k, v := range m {
		fmt.Printf("export %s=\"%s\"\n", k, v)
	}
}

// outputJSON outputs a map of environment variables in JSON format.
func outputJSON(m map[string]string) {
	json, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		slog.Error("failed to marshal the output into JSON")
		os.Exit(1)
	}
	fmt.Println(string(json))
}

// outputTable outputs a map of environment variables in a table format.
func outputTable(m map[string]string) {
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
