// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package output

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

// PropertyPair represents a key-value pair for ordered rendering
type PropertyPair struct {
	Key   string
	Value interface{}
}
