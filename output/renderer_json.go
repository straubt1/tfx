// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package output

import (
	"encoding/json"
	"os"
)

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
