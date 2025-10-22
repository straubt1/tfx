// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

type RunDiscardView struct{ *BaseView }

func NewRunDiscardView() *RunDiscardView { return &RunDiscardView{NewBaseView()} }

func (v *RunDiscardView) Render(runID string) error {
	if v.IsJSON() {
		return v.renderer.RenderJSON(map[string]string{"discardedRunId": runID})
	}
	props := []PropertyPair{{Key: "Discarded run id", Value: runID}}
	return v.renderer.RenderProperties(props)
}
