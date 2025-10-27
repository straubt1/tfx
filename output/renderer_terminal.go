// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/logrusorgru/aurora"
)

// TerminalRenderer renders output for human consumption
type TerminalRenderer struct{}

func NewTerminalRenderer() *TerminalRenderer {
	return &TerminalRenderer{}
}

func (r *TerminalRenderer) RenderError(err error) error {
	fmt.Printf("%s %s\n", aurora.Red("✗"), aurora.Bold(aurora.Red("Error")))
	fmt.Println()
	fmt.Printf("%s\n", err.Error())
	fmt.Println()
	return nil // Error is rendered, do not return it
}

func (r *TerminalRenderer) Message(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

func (r *TerminalRenderer) MessageCommandHeader(format string, args ...interface{}) {
	// Apply green color to all arguments
	greenArgs := make([]interface{}, len(args))
	for i, arg := range args {
		greenArgs[i] = aurora.Green(arg)
	}
	message := fmt.Sprintf(format, greenArgs...)
	fmt.Printf("%s\n", message)

	// Calculate visible length and print separator
	visibleLength := len(fmt.Sprintf(format, args...))
	separator := ""
	for i := 0; i < visibleLength; i++ {
		separator += "─"
	}
	fmt.Printf("%s\n", separator)
}

func (r *TerminalRenderer) MessageCommandFilter(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

func (r *TerminalRenderer) RenderTable(headers []string, rows [][]interface{}) error {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// Convert headers to table.Row
	headerRow := make(table.Row, len(headers))
	for i, h := range headers {
		headerRow[i] = h
	}
	t.AppendHeader(headerRow)

	// Add rows
	for _, row := range rows {
		tableRow := make(table.Row, len(row))
		copy(tableRow, row)
		t.AppendRow(tableRow)
	}

	t.SetStyle(table.StyleRounded)
	t.SetAutoIndex(false)
	t.Style().Options.SeparateRows = false
	t.Render()

	return nil
}

func (r *TerminalRenderer) RenderFields(fields map[string]interface{}) error {
	// Calculate max key length for alignment
	maxLen := 0
	for key := range fields {
		if len(key) > maxLen {
			maxLen = len(key)
		}
	}

	// Print aligned fields
	for key, value := range fields {
		fmt.Printf("%-*s %s\n", maxLen+1, aurora.Bold(key+":"), aurora.Blue(value))
	}

	return nil
}

func (r *TerminalRenderer) RenderProperties(properties []PropertyPair) error {
	// Calculate max key length for alignment
	maxLen := 0
	for _, prop := range properties {
		if len(prop.Key) > maxLen {
			maxLen = len(prop.Key)
		}
	}

	// Print aligned properties
	for _, prop := range properties {
		fmt.Printf("%-*s  %s\n", maxLen+1, aurora.Bold(prop.Key+":"), aurora.Blue(fmt.Sprint(prop.Value)))
	}

	return nil
}

func (r *TerminalRenderer) RenderTags(label string, tags []PropertyPair) error {
	fmt.Println()
	fmt.Printf("%s\n", aurora.Bold(label+":"))

	if len(tags) == 0 {
		return nil
	}

	// Calculate max key length for alignment
	maxLen := 0
	for _, tag := range tags {
		if len(tag.Key) > maxLen {
			maxLen = len(tag.Key)
		}
	}

	// Print aligned tags, indented
	for _, tag := range tags {
		fmt.Printf("  %-*s %s\n", maxLen+1, aurora.Bold(tag.Key+":"), aurora.Blue(fmt.Sprint(tag.Value)))
	}

	return nil
}

func (r *TerminalRenderer) RenderJSON(data interface{}) error {
	return json.NewEncoder(os.Stdout).Encode(data)
}
