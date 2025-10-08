// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import "github.com/spf13/viper"

// BaseView provides common functionality for all command views
type BaseView struct {
	renderer Renderer
}

// NewBaseView creates a base view with the appropriate renderer
func NewBaseView(isJSON bool) *BaseView {
	var renderer Renderer
	if isJSON {
		renderer = NewJSONRenderer()
	} else {
		renderer = NewTerminalRenderer()
	}
	return &BaseView{renderer: renderer}
}

// NewBaseViewFromViper creates a base view using viper config
func NewBaseViewFromViper() *BaseView {
	return NewBaseView(viper.GetBool("json"))
}

// RenderError renders an error in the appropriate format
func (v *BaseView) RenderError(err error) error {
	return v.renderer.RenderError(err)
}

// IsJSON returns true if using JSON output mode
func (v *BaseView) IsJSON() bool {
	_, ok := v.renderer.(*JSONRenderer)
	return ok
}

// Renderer returns the underlying renderer
func (v *BaseView) Renderer() Renderer {
	return v.renderer
}
