// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package output

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/viper"
)

var (
	once     sync.Once
	instance *Output
)

// OutputMode represents the output format
type OutputMode int

const (
	ModeTerminal OutputMode = iota
	ModeJSON
)

// PropertyPair represents a key-value pair for ordered rendering
type PropertyPair struct {
	Key   string
	Value interface{}
}

// Renderer defines the interface for rendering output
type Renderer interface {
	// RenderError outputs an error
	RenderError(err error) error

	// Message outputs informational text (suppressed in JSON mode)
	Message(format string, args ...interface{})

	// MessageCommandHeader outputs command header information (suppressed in JSON mode)
	MessageCommandHeader(format string, args ...interface{})

	// MessageCommandFilter outputs filter information (suppressed in JSON mode)
	MessageCommandFilter(format string, args ...interface{})

	// RenderTable outputs a table
	RenderTable(headers []string, rows [][]interface{}) error

	// RenderFields outputs key-value pairs
	RenderFields(fields map[string]interface{}) error

	// RenderProperties outputs key-value pairs with formatting in order
	RenderProperties(properties []PropertyPair) error

	// RenderTags outputs tags with a header and formatted key-value pairs
	RenderTags(label string, tags []PropertyPair) error

	// RenderJSON outputs raw JSON
	RenderJSON(data interface{}) error
}

// Output is the singleton that manages all output operations
type Output struct {
	mode     OutputMode
	renderer Renderer
	spinner  *Spinner
	logger   *Logger
	mu       sync.RWMutex
}

// initialize sets up the singleton Output instance
// Called automatically by Get() on first access
func initialize() {
	once.Do(func() {
		// Determine output mode from viper
		mode := ModeTerminal
		if viper.GetBool("json") {
			mode = ModeJSON
		}

		// Get logging configuration from environment
		logLevel := os.Getenv("TFX_LOG")
		logPath := os.Getenv("TFX_LOG_PATH")

		var renderer Renderer
		var spinner *Spinner

		if mode == ModeJSON {
			renderer = NewJSONRenderer()
			spinner = nil // No spinner in JSON mode
		} else {
			renderer = NewTerminalRenderer()
			spinner = NewSpinner()
		}

		logger := NewLogger(logLevel, logPath)

		instance = &Output{
			mode:     mode,
			renderer: renderer,
			spinner:  spinner,
			logger:   logger,
		}
	})
}

// Get returns the singleton Output instance
// Automatically initializes on first call using viper configuration
func Get() *Output {
	initialize()
	return instance
}

// Mode returns the current output mode
func (o *Output) Mode() OutputMode {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.mode
}

// IsJSON returns true if output mode is JSON
func (o *Output) IsJSON() bool {
	return o.Mode() == ModeJSON
}

// Renderer returns the underlying renderer
func (o *Output) Renderer() Renderer {
	return o.renderer
}

// Logger returns the underlying logger
func (o *Output) Logger() *Logger {
	return o.logger
}

// Spinner returns the spinner (nil in JSON mode)
func (o *Output) Spinner() *Spinner {
	return o.spinner
}

// RenderError renders an error
func (o *Output) RenderError(err error) error {
	if o.spinner != nil {
		o.spinner.Stop()
		defer o.spinner.Start()
	}
	return o.renderer.RenderError(err)
}

// Message outputs a message
func (o *Output) Message(format string, args ...interface{}) {
	if o.spinner != nil {
		o.spinner.Stop()
		defer o.spinner.Start()
	}
	o.renderer.Message(format, args...)
}

// MessageCommandHeader outputs a command header
func (o *Output) MessageCommandHeader(format string, args ...interface{}) {
	if o.spinner != nil {
		o.spinner.Stop()
		defer o.spinner.Start()
	}
	o.renderer.MessageCommandHeader(format, args...)
}

// MessageCommandFilter outputs filter information
func (o *Output) MessageCommandFilter(format string, args ...interface{}) {
	if o.spinner != nil {
		o.spinner.Stop()
		defer o.spinner.Start()
	}
	o.renderer.MessageCommandFilter(format, args...)
}

// RenderTable renders a table
func (o *Output) RenderTable(headers []string, rows [][]interface{}) error {
	if o.spinner != nil {
		o.spinner.Stop()
		defer o.spinner.Start()
	}
	return o.renderer.RenderTable(headers, rows)
}

// RenderFields renders key-value fields
func (o *Output) RenderFields(fields map[string]interface{}) error {
	if o.spinner != nil {
		o.spinner.Stop()
		defer o.spinner.Start()
	}
	return o.renderer.RenderFields(fields)
}

// RenderProperties renders ordered properties
func (o *Output) RenderProperties(properties []PropertyPair) error {
	if o.spinner != nil {
		o.spinner.Stop()
		defer o.spinner.Start()
	}
	return o.renderer.RenderProperties(properties)
}

// RenderTags renders tags with a label
func (o *Output) RenderTags(label string, tags []PropertyPair) error {
	if o.spinner != nil {
		o.spinner.Stop()
		defer o.spinner.Start()
	}
	return o.renderer.RenderTags(label, tags)
}

// RenderJSON renders data as JSON
func (o *Output) RenderJSON(data interface{}) error {
	if o.spinner != nil {
		o.spinner.Stop()
		defer o.spinner.Start()
	}
	return o.renderer.RenderJSON(data)
}

// Close stops the spinner if running and closes any open output streams
// This should be called before the application exits to ensure clean shutdown
func (o *Output) Close() error {
	// Stop spinner if it's running
	if o.spinner != nil {
		o.spinner.FinalStop()
		o.spinner.FinalStop()
		o.spinner.FinalStop()
	}

	// Future: Close any file handles or other resources
	// For now, this is primarily for stopping the spinner gracefully

	return nil
}

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
	// Suppress in JSON mode
}

func (r *JSONRenderer) MessageCommandFilter(format string, args ...interface{}) {
	// Suppress in JSON mode
}

func (r *JSONRenderer) RenderTable(headers []string, rows [][]interface{}) error {
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
