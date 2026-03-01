// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	"github.com/straubt1/tfx/output"
)

// BaseView provides common functionality for all command views
// It wraps the singleton output instance for convenience
type BaseView struct {
	out *output.Output
}

// NewBaseView creates a base view using the singleton output instance
func NewBaseView() *BaseView {
	return &BaseView{
		out: output.Get(),
	}
}

// Output returns the singleton output instance
func (v *BaseView) Output() *output.Output {
	return v.out
}

// RenderError renders an error in the appropriate format
func (v *BaseView) RenderError(err error) error {
	return v.out.RenderError(err)
}

// IsJSON returns true if using JSON output mode
func (v *BaseView) IsJSON() bool {
	return v.out.IsJSON()
}

// Renderer returns the underlying renderer
func (v *BaseView) Renderer() output.Renderer {
	return v.out.Renderer()
}

// PrintCommandHeader prints a command header message (suppressed in JSON mode)
func (v *BaseView) PrintCommandHeader(format string, args ...interface{}) {
	v.out.MessageCommandHeader(format, args...)
}

// PrintCommandFilter prints a filter message (suppressed in JSON mode)
func (v *BaseView) PrintCommandFilter(format string, args ...interface{}) {
	v.out.MessageCommandFilter(format, args...)
}

// PropertyPair is re-exported from output package for convenience
type PropertyPair = output.PropertyPair

// GetIfSpecified is re-exported from output package for convenience
var GetIfSpecified = output.GetIfSpecified
