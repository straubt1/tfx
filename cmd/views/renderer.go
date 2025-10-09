// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hashicorp/jsonapi"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/logrusorgru/aurora"
)

// Renderer handles output formatting
type Renderer interface {
	// RenderError outputs an error
	RenderError(err error) error

	// Message outputs informational text (suppressed in JSON mode)
	Message(format string, args ...interface{})

	// MessageCommandHeader outputs command header information in green (suppressed in JSON mode)
	MessageCommandHeader(format string, args ...interface{})

	// RenderTable outputs a table
	RenderTable(headers []string, rows [][]interface{}) error

	// RenderFields outputs key-value pairs
	RenderFields(fields map[string]interface{}) error

	// RenderProperties outputs key-value pairs with formatting in order (bold keys, blue values, aligned)
	RenderProperties(properties []PropertyPair) error

	// RenderTags outputs tags with a header and formatted key-value pairs (bold keys, blue values, indented)
	RenderTags(label string, tags []PropertyPair) error

	// RenderJSON outputs raw JSON (used by JSON renderer for everything)
	RenderJSON(data interface{}) error
}

// PropertyPair represents a key-value pair for ordered rendering
type PropertyPair struct {
	Key   string
	Value interface{}
}

// TerminalRenderer renders output for human consumption
type TerminalRenderer struct{}

func NewTerminalRenderer() *TerminalRenderer {
	return &TerminalRenderer{}
}

// Helper functions
// GetIfSpecified retrieves the value from a NullableAttr if specified
// Returns empty string and nil error if not specified
func GetIfSpecified(property jsonapi.NullableAttr[string]) (duration string, err error) {
	if property.IsSpecified() {
		if duration, err = property.Get(); err != nil {
			return "", err
		}
	}
	return duration, nil
}

func (r *TerminalRenderer) RenderError(err error) error {
	fmt.Printf("%s %s\n", aurora.Red("✗"), aurora.Bold(aurora.Red("Error")))
	fmt.Println()
	fmt.Printf("%s\n", err.Error())
	fmt.Println()
	// return err
	return nil // Error is rendered, do not return it
}

func (r *TerminalRenderer) Message(format string, args ...interface{}) {
	// Allow colored output in messages
	fmt.Printf(format+"\n", args...)
}

func (r *TerminalRenderer) MessageCommandHeader(format string, args ...interface{}) {
	// Apply green color to all arguments
	greenArgs := make([]interface{}, len(args))
	for i, arg := range args {
		greenArgs[i] = aurora.Green(arg)
	}
	// Format and print the message
	message := fmt.Sprintf(format, greenArgs...)
	fmt.Printf("%s\n", message)

	// Calculate the visible length of the message (without ANSI color codes)
	// We need to strip color codes to get the actual display width
	visibleLength := len(fmt.Sprintf(format, args...))

	// Print separator line matching the header width
	separator := ""
	for i := 0; i < visibleLength; i++ {
		separator += "─"
	}
	fmt.Printf("%s\n", separator)
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

	// Print aligned properties with bold keys and blue values in order
	for _, prop := range properties {
		fmt.Printf("%-*s  %s\n", maxLen+1, aurora.Bold(prop.Key+":"), aurora.Blue(prop.Value))
	}

	return nil
}

func (r *TerminalRenderer) RenderTags(label string, tags []PropertyPair) error {
	// Print Tags header with blank line before
	fmt.Printf("\n%s\n", aurora.Bold(label+":"))

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

	// Print aligned tags with bold keys and blue values, indented
	for _, tag := range tags {
		fmt.Printf("  %-*s %s\n", maxLen+1, aurora.Bold(tag.Key+":"), aurora.Blue(tag.Value))
	}

	return nil
}

func (r *TerminalRenderer) RenderJSON(data interface{}) error {
	// Terminal renderer doesn't typically use this, but implement for interface
	return json.NewEncoder(os.Stdout).Encode(data)
}

// JSONRenderer renders output as JSON
type JSONRenderer struct{}

func NewJSONRenderer() *JSONRenderer {
	return &JSONRenderer{}
}

func (r *JSONRenderer) RenderError(err error) error {
	errorOutput := map[string]interface{}{
		"error": err.Error(),
	}
	return json.NewEncoder(os.Stdout).Encode(errorOutput)
}

func (r *JSONRenderer) Message(format string, args ...interface{}) {
	// Suppress messages in JSON mode
}

func (r *JSONRenderer) MessageCommandHeader(format string, args ...interface{}) {
	// Suppress command header messages in JSON mode
}

func (r *JSONRenderer) RenderTable(headers []string, rows [][]interface{}) error {
	// Convert table to array of objects
	result := make([]map[string]interface{}, len(rows))
	for i, row := range rows {
		obj := make(map[string]interface{})
		for j, header := range headers {
			if j < len(row) {
				obj[header] = row[j]
			}
		}
		result[i] = obj
	}
	return json.NewEncoder(os.Stdout).Encode(result)
}

func (r *JSONRenderer) RenderFields(fields map[string]interface{}) error {
	return json.NewEncoder(os.Stdout).Encode(fields)
}

func (r *JSONRenderer) RenderProperties(properties []PropertyPair) error {
	// Convert ordered properties to map for JSON output
	result := make(map[string]interface{})
	for _, prop := range properties {
		result[prop.Key] = prop.Value
	}
	return json.NewEncoder(os.Stdout).Encode(result)
}

func (r *JSONRenderer) RenderTags(label string, tags []PropertyPair) error {
	return nil
}

func (r *JSONRenderer) RenderJSON(data interface{}) error {
	return json.NewEncoder(os.Stdout).Encode(data)
}
