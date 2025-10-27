// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	tfe "github.com/hashicorp/go-tfe"
)

type ConfigVersionCreateView struct{ *BaseView }

func NewConfigVersionCreateView() *ConfigVersionCreateView {
	return &ConfigVersionCreateView{NewBaseView()}
}

func (v *ConfigVersionCreateView) Render(cv *tfe.ConfigurationVersion) error {
	if v.IsJSON() {
		return v.Output().RenderJSON(map[string]interface{}{
			"id":          cv.ID,
			"speculative": cv.Speculative,
			"status":      cv.Status,
		})
	}
	props := []PropertyPair{
		{Key: "Created", Value: "Success"},
		{Key: "ID", Value: cv.ID},
		{Key: "Speculative", Value: cv.Speculative},
	}
	return v.Output().RenderProperties(props)
}
