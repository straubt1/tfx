// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

type WorkspaceLockResult struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type WorkspaceLockView struct{ *BaseView }

func NewWorkspaceLockView() *WorkspaceLockView { return &WorkspaceLockView{NewBaseView()} }

func (v *WorkspaceLockView) RenderSingle(name, status string) error {
	if v.IsJSON() {
		return v.Output().RenderJSON(WorkspaceLockResult{Name: name, Status: status})
	}
	props := []PropertyPair{{Key: name, Value: status}}
	return v.Output().RenderProperties(props)
}

func (v *WorkspaceLockView) RenderBulk(results []WorkspaceLockResult) error {
	if v.IsJSON() {
		return v.Output().RenderJSON(results)
	}
	headers := []string{"Workspace", "Status"}
	rows := make([][]interface{}, len(results))
	for i, r := range results {
		rows[i] = []interface{}{r.Name, r.Status}
	}
	return v.Output().RenderTable(headers, rows)
}
