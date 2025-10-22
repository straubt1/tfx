// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	tfe "github.com/hashicorp/go-tfe"
)

type RunCreateView struct{ *BaseView }

func NewRunCreateView() *RunCreateView { return &RunCreateView{NewBaseView()} }

type runCreateOutput struct {
	ID                   string `json:"id"`
	ConfigurationVersion string `json:"configurationVersion"`
	TerraformVersion     string `json:"terraformVersion"`
	Link                 string `json:"link"`
}

func (v *RunCreateView) Render(run *tfe.Run, link string) error {
	if v.IsJSON() {
		cv := ""
		if run.ConfigurationVersion != nil {
			cv = run.ConfigurationVersion.ID
		}
		return v.renderer.RenderJSON(runCreateOutput{
			ID:                   run.ID,
			ConfigurationVersion: cv,
			TerraformVersion:     run.TerraformVersion,
			Link:                 link,
		})
	}

	cv := ""
	if run.ConfigurationVersion != nil {
		cv = run.ConfigurationVersion.ID
	}
	props := []PropertyPair{
		{Key: "ID", Value: run.ID},
		{Key: "Configuration Version", Value: cv},
		{Key: "Terraform Version", Value: run.TerraformVersion},
		{Key: "Link", Value: link},
	}
	return v.renderer.RenderProperties(props)
}
