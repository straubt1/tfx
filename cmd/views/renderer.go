// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/logrusorgru/aurora"
)

// Renderer handles output formatting
type Renderer interface {
	// RenderError outputs an error
	RenderError(err error) error

	// Message outputs informational text (suppressed in JSON mode)
	Message(format string, args ...interface{})

	// RenderTable outputs a table
	RenderTable(headers []string, rows [][]interface{}) error

	// RenderFields outputs key-value pairs
	RenderFields(fields map[string]interface{}) error

	// RenderJSON outputs raw JSON (used by JSON renderer for everything)
	RenderJSON(data interface{}) error
}

// TerminalRenderer renders output for human consumption
type TerminalRenderer struct{}

func NewTerminalRenderer() *TerminalRenderer {
	return &TerminalRenderer{}
}

func (r *TerminalRenderer) RenderError(err error) error {
	fmt.Println()
	fmt.Println(aurora.Red(fmt.Sprintf("Error: %s", err.Error())))
	return err
}

func (r *TerminalRenderer) Message(format string, args ...interface{}) {
	// Allow colored output in messages
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

func (r *JSONRenderer) RenderJSON(data interface{}) error {
	return json.NewEncoder(os.Stdout).Encode(data)
}
