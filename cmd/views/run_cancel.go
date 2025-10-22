// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

type RunCancelView struct{ *BaseView }

func NewRunCancelView() *RunCancelView { return &RunCancelView{NewBaseView()} }

func (v *RunCancelView) Render(runID string) error {
	if v.IsJSON() {
		return v.renderer.RenderJSON(map[string]string{"cancelledRunId": runID})
	}
	props := []PropertyPair{{Key: "Cancelled run id", Value: runID}}
	return v.renderer.RenderProperties(props)
}
