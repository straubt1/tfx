// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	"encoding/json"
)

type PlanJSONOutputView struct{ *BaseView }

func NewPlanJSONOutputView() *PlanJSONOutputView { return &PlanJSONOutputView{NewBaseView()} }

func (v *PlanJSONOutputView) Render(jsonOutput []byte) error {
	// Parse the JSON to ensure it's valid, then pretty-print it
	var data interface{}
	err := json.Unmarshal(jsonOutput, &data)
	if err != nil {
		return v.RenderError(err)
	}

	// For both terminal and JSON output, render the parsed JSON
	// This ensures consistent formatting
	return v.Output().RenderJSON(data)
}
