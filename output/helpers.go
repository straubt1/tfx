// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package output

import "github.com/hashicorp/jsonapi"

// Helper functions for common operations

// GetIfSpecified retrieves the value from a NullableAttr if specified
// Returns empty string and nil error if not specified
func GetIfSpecified(property jsonapi.NullableAttr[string]) (string, error) {
	if property.IsSpecified() {
		if value, err := property.Get(); err != nil {
			return "", err
		} else {
			return value, nil
		}
	}
	return "", nil
}
