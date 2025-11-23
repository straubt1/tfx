// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/logrusorgru/aurora"

	"github.com/straubt1/tfx/output"
)

// Renderer handles output formatting
type Renderer interface {
	// RenderError outputs an error
	RenderError(err error) error

	// Message outputs informational text (suppressed in JSON mode)
	Message(format string, args ...interface{})

	// MessageCommandHeader outputs command header information in green (suppressed in JSON mode)
	MessageCommandHeader(format string, args ...interface{})

	// MessageCommandFilter outputs filter information without separator line (suppressed in JSON mode)
	MessageCommandFilter(format string, args ...interface{})

	// RenderTable outputs a table
	RenderTable(headers []string, rows [][]interface{}) error

	// RenderFields outputs key-value pairs
	RenderFields(fields map[string]interface{}) error

	// RenderProperties outputs key-value pairs with formatting in order (bold keys, blue values, aligned)
	RenderProperties(properties []output.PropertyPair) error

	// RenderTags outputs tags with a header and formatted key-value pairs (bold keys, blue values, indented)
	RenderTags(label string, tags []output.PropertyPair) error

	// RenderJSON outputs raw JSON (used by JSON renderer for everything)
	RenderJSON(data interface{}) error
}

// TerminalRenderer renders output for human consumption
type TerminalRenderer struct {
	s       *spinner.Spinner
	running bool
	mu      sync.Mutex
	depth   int
}

func NewTerminalRenderer() *TerminalRenderer {
	// Initialize spinner with character set type [14] and start automatically
	sp := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	_ = sp.Color("cyan")
	sp.Suffix = "  TFx is working..."
	// Start the spinner immediately
	// sp.Start()

	return &TerminalRenderer{s: sp, running: true}
}

func (r *TerminalRenderer) RenderError(err error) error {
	r.stopSpinner()
	defer r.startSpinner()
	r.Message("%s %s", aurora.Red("✗"), aurora.Bold(aurora.Red("Error")))
	r.Message("")
	r.Message("%s", err.Error())
	r.Message("")
	// return err
	return nil // Error is rendered, do not return it
}

func (r *TerminalRenderer) Message(format string, args ...interface{}) {
	r.stopSpinner()
	defer r.startSpinner()
	// Allow colored output in messages
	fmt.Printf(format+"\n", args...)
}

func (r *TerminalRenderer) MessageCommandHeader(format string, args ...interface{}) {
	r.stopSpinner()
	defer r.startSpinner()
	// Apply green color to all arguments
	greenArgs := make([]interface{}, len(args))
	for i, arg := range args {
		greenArgs[i] = aurora.Green(arg)
	}
	// Format and print the message
	message := fmt.Sprintf(format, greenArgs...)
	r.Message("%s", message)

	// Calculate the visible length of the message (without ANSI color codes)
	// We need to strip color codes to get the actual display width
	visibleLength := len(fmt.Sprintf(format, args...))

	// Print separator line matching the header width
	separator := ""
	for i := 0; i < visibleLength; i++ {
		separator += "─"
	}
	r.Message("%s", separator)
}

func (r *TerminalRenderer) MessageCommandFilter(format string, args ...interface{}) {
	r.stopSpinner()
	defer r.startSpinner()
	// Simply print the message without separator line
	r.Message(format, args...)
}

func (r *TerminalRenderer) RenderTable(headers []string, rows [][]interface{}) error {
	r.stopSpinner()
	defer r.startSpinner()
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
	r.stopSpinner()
	defer r.startSpinner()
	// Calculate max key length for alignment
	maxLen := 0
	for key := range fields {
		if len(key) > maxLen {
			maxLen = len(key)
		}
	}

	// Print aligned fields
	for key, value := range fields {
		r.Message("%-*s %s", maxLen+1, aurora.Bold(key+":"), aurora.Blue(value))
	}

	return nil
}

func (r *TerminalRenderer) RenderProperties(properties []output.PropertyPair) error {
	r.stopSpinner()
	defer r.startSpinner()
	// Calculate max key length for alignment
	maxLen := 0
	for _, prop := range properties {
		if len(prop.Key) > maxLen {
			maxLen = len(prop.Key)
		}
	}

	// Print aligned properties with bold keys and blue values in order
	for _, prop := range properties {
		r.Message("%-*s  %s", maxLen+1, aurora.Bold(prop.Key+":"), aurora.Blue(fmt.Sprint(prop.Value)))
	}

	return nil
}

func (r *TerminalRenderer) RenderTags(label string, tags []output.PropertyPair) error {
	r.stopSpinner()
	defer r.startSpinner()
	// Print Tags header with blank line before
	r.Message("")
	r.Message("%s", aurora.Bold(label+":"))

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
		r.Message("  %-*s %s", maxLen+1, aurora.Bold(tag.Key+":"), aurora.Blue(fmt.Sprint(tag.Value)))
	}

	return nil
}

func (r *TerminalRenderer) RenderJSON(data interface{}) error {
	r.stopSpinner()
	defer r.startSpinner()
	// Terminal renderer doesn't typically use this, but implement for interface
	return json.NewEncoder(os.Stdout).Encode(data)
}

// spinner control helpers
func (r *TerminalRenderer) stopSpinner() {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Increase render depth; stop spinner only when entering outermost render
	r.depth++
	if r.depth == 1 && r.s != nil && r.running {
		r.s.Stop()
		r.running = false
	}
}

func (r *TerminalRenderer) startSpinner() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.depth > 0 {
		r.depth--
	}
	// Restart spinner only when exiting outermost render
	if r.depth == 0 && r.s != nil && !r.running {
		r.s.Start()
		r.running = true
	}
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

func (r *JSONRenderer) MessageCommandFilter(format string, args ...interface{}) {
	// Suppress filter messages in JSON mode
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

func (r *JSONRenderer) RenderProperties(properties []output.PropertyPair) error {
	// Convert ordered properties to map for JSON output
	result := make(map[string]interface{})
	for _, prop := range properties {
		result[prop.Key] = prop.Value
	}
	return json.NewEncoder(os.Stdout).Encode(result)
}

func (r *JSONRenderer) RenderTags(label string, tags []output.PropertyPair) error {
	return nil
}

func (r *JSONRenderer) RenderJSON(data interface{}) error {
	return json.NewEncoder(os.Stdout).Encode(data)
}
