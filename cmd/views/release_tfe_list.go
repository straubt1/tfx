// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

// ReleaseTfeListView handles rendering for release tfe list command
type ReleaseTfeListView struct {
	*BaseView
}

func NewReleaseTfeListView() *ReleaseTfeListView {
	return &ReleaseTfeListView{
		BaseView: NewBaseView(),
	}
}

// releaseTfeListOutput is a JSON-safe representation of a TFE release for list views
type releaseTfeListOutput struct {
	Tag     string `json:"tag"`
	Created string `json:"created,omitempty"`
}

// Render renders TFE releases list
func (v *ReleaseTfeListView) Render(releases []map[string]interface{}) error {
	if v.IsJSON() {
		// Convert to JSON-safe representation
		output := make([]releaseTfeListOutput, len(releases))
		for i, rel := range releases {
			output[i] = releaseTfeListOutput{
				Tag:     rel["Tag"].(string),
				Created: rel["Created"].(string),
			}
		}
		return v.Output().RenderJSON(output)
	}

	// Terminal mode: render as table
	headers := []string{"Tag", "Date"}

	rows := make([][]interface{}, len(releases))
	for i, rel := range releases {
		rows[i] = []interface{}{
			rel["Tag"],
			rel["Created"],
		}
	}

	return v.Output().RenderTable(headers, rows)
}
