// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	"fmt"

	tfe "github.com/hashicorp/go-tfe"
)

// ProjectShowView handles rendering for project show command
type ProjectShowView struct {
	*BaseView
}

func NewProjectShowView() *ProjectShowView {
	return &ProjectShowView{
		BaseView: NewBaseView(),
	}
}

// projectShowOutput is a JSON-safe representation of a project
type projectShowOutput struct {
	Organization                string            `json:"organization"`
	Name                        string            `json:"name"`
	ID                          string            `json:"id"`
	Description                 string            `json:"description"`
	DefaultExecutionMode        string            `json:"defaultExecutionMode"`
	AutoDestroyActivityDuration string            `json:"autoDestroyActivityDuration,omitempty"`
	Tags                        map[string]string `json:"tags,omitempty"`
}

// Render renders a single project's details
func (v *ProjectShowView) Render(orgName string, project *tfe.Project) error {
	if project == nil {
		return v.RenderError(fmt.Errorf("project not found"))
	}

	duration, err := GetIfSpecified(project.AutoDestroyActivityDuration)
	if err != nil {
		return v.RenderError(fmt.Errorf("failed to get auto-destroy duration: %w", err))
	}

	if v.IsJSON() {
		// JSON mode: convert to JSON-safe structure
		output := projectShowOutput{
			Organization:                orgName,
			Name:                        project.Name,
			ID:                          project.ID,
			Description:                 project.Description,
			DefaultExecutionMode:        project.DefaultExecutionMode,
			AutoDestroyActivityDuration: duration,
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

	// Terminal mode: render key fields in order
	properties := []PropertyPair{
		{Key: "Name", Value: project.Name},
		{Key: "ID", Value: project.ID},
		{Key: "Description", Value: project.Description},
		{Key: "Default Execution Mode", Value: project.DefaultExecutionMode},
		{Key: "Auto Destroy", Value: duration},
	}

	err = v.renderer.RenderProperties(properties)
	if err != nil {
		return err
	}

	// Extract and display tags
	tags := []PropertyPair{}
	for _, tag := range project.EffectiveTagBindings {
		tags = append(tags, PropertyPair{Key: tag.Key, Value: tag.Value})
	}
	err = v.renderer.RenderTags("Tags", tags)
	if err != nil {
		return err
	}

	return nil
}
