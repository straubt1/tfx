// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

type PlanLogsView struct{ *BaseView }

func NewPlanLogsView() *PlanLogsView { return &PlanLogsView{NewBaseView()} }

type planLogsOutput struct {
	Logs []string `json:"logs"`
}

func (v *PlanLogsView) Render(logs []string) error {
	if v.IsJSON() {
		return v.Output().RenderJSON(planLogsOutput{
			Logs: logs,
		})
	}

	// Terminal mode: print logs directly
	for _, line := range logs {
		v.Output().Message(line)
	}
	return nil
}
