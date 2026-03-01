// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package output

import (
	"os"
	"sync"

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
			// Only create spinner if logging is not enabled
			// When TFX_LOG is set, disable spinner to avoid interference with log output
			if logLevel == "" || logLevel == "NONE" || logLevel == "none" {
				spinner = NewSpinner()
			} else {
				spinner = nil // No spinner when logging is enabled
			}
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
