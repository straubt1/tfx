// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	"fmt"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/logrusorgru/aurora"
)

// ProjectShowView handles rendering for project show command
type ProjectShowView struct {
	*BaseView
}

func NewProjectShowView(isJSON bool) *ProjectShowView {
	return &ProjectShowView{
		BaseView: NewBaseView(isJSON),
	}
}

// projectShowOutput is a JSON-safe representation of a project
type projectShowOutput struct {
	Organization                string            `json:"organization"`
	Name                        string            `json:"name"`
	ID                          string            `json:"id"`
	Description                 string            `json:"description"`
	DefaultExecutionMode        string            `json:"defaultExecutionMode"`
	AutoDestroyActivityDuration *string           `json:"autoDestroyActivityDuration,omitempty"`
	Tags                        map[string]string `json:"tags,omitempty"`
}

// Render renders a single project's details
func (v *ProjectShowView) Render(orgName string, project *tfe.Project) error {
	// Print context messages (only in terminal)
	v.renderer.Message("Organization Name: %s", aurora.Green(orgName))

	if project == nil {
		return v.RenderError(fmt.Errorf("project not found"))
	}

	if v.IsJSON() {
		// JSON mode: convert to JSON-safe structure
		output := projectShowOutput{
			Organization:         orgName,
			Name:                 project.Name,
			ID:                   project.ID,
			Description:          project.Description,
			DefaultExecutionMode: project.DefaultExecutionMode,
		}

		// Handle nullable fields
		if project.AutoDestroyActivityDuration.IsSpecified() {
			if duration, err := project.AutoDestroyActivityDuration.Get(); err == nil {
				output.AutoDestroyActivityDuration = &duration
			}
		}

		// Extract tags
		if len(project.EffectiveTagBindings) > 0 {
			output.Tags = make(map[string]string)
			for _, tag := range project.EffectiveTagBindings {
				output.Tags[tag.Key] = tag.Value
			}
		}

		return v.renderer.RenderJSON(output)
	}

	// Terminal mode: render key fields
	fields := map[string]interface{}{
		"Name":                 project.Name,
		"ID":                   project.ID,
		"Description":          project.Description,
		"DefaultExecutionMode": project.DefaultExecutionMode,
	}

	// Render the fields first
	err := v.renderer.RenderFields(fields)
	if err != nil {
		return err
	}

	// Add optional fields
	if project.AutoDestroyActivityDuration.IsSpecified() {
		if duration, err := project.AutoDestroyActivityDuration.Get(); err == nil {
			v.renderer.Message("%-25s %s", aurora.Bold("Auto Destroy Activity Duration:"), aurora.Blue(duration))
		}
	}

	// Extract and display tags
	if len(project.EffectiveTagBindings) > 0 {
		v.renderer.Message("\n%s", aurora.Bold("Tags:"))
		for _, tag := range project.EffectiveTagBindings {
			v.renderer.Message("  %s %s", aurora.Bold(tag.Key+":"), aurora.Blue(tag.Value))
		}
	}

	return nil
}
