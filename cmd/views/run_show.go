// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	tfe "github.com/hashicorp/go-tfe"
)

type RunShowView struct{ *BaseView }

func NewRunShowView() *RunShowView { return &RunShowView{NewBaseView()} }

type runShowOutput struct {
	ID                   string `json:"id"`
	ConfigurationVersion string `json:"configurationVersion"`
	Status               string `json:"status"`
	Message              string `json:"message"`
	TerraformVersion     string `json:"terraformVersion"`
	Created              string `json:"created"`
}

func (v *RunShowView) Render(run *tfe.Run) error {
	if v.IsJSON() {
		cv := ""
		if run.ConfigurationVersion != nil {
			cv = run.ConfigurationVersion.ID
		}
		return v.renderer.RenderJSON(runShowOutput{
			ID:                   run.ID,
			ConfigurationVersion: cv,
			Status:               string(run.Status),
			Message:              run.Message,
			TerraformVersion:     run.TerraformVersion,
			Created:              FormatDateTime(run.CreatedAt),
		})
	}

	cv := ""
	if run.ConfigurationVersion != nil {
		cv = run.ConfigurationVersion.ID
	}
	props := []PropertyPair{
		{Key: "ID", Value: run.ID},
		{Key: "Configuration Version", Value: cv},
		{Key: "Status", Value: run.Status},
		{Key: "Message", Value: run.Message},
		{Key: "Terraform Version", Value: run.TerraformVersion},
		{Key: "Created", Value: FormatDateTime(run.CreatedAt)},
	}
	return v.renderer.RenderProperties(props)
}
