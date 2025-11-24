// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	tfe "github.com/hashicorp/go-tfe"
)

type PlanShowView struct{ *BaseView }

func NewPlanShowView() *PlanShowView { return &PlanShowView{NewBaseView()} }

type planStatusTimestampsOutput struct {
	QueuedAt        string `json:"queuedAt,omitempty"`
	StartedAt       string `json:"startedAt,omitempty"`
	FinishedAt      string `json:"finishedAt,omitempty"`
	CanceledAt      string `json:"canceledAt,omitempty"`
	ErroredAt       string `json:"erroredAt,omitempty"`
	ForceCanceledAt string `json:"forceCanceledAt,omitempty"`
}

type planShowOutput struct {
	ID                     string                      `json:"id"`
	Status                 string                      `json:"status"`
	HasChanges             bool                        `json:"hasChanges"`
	GeneratedConfiguration bool                        `json:"generatedConfiguration"`
	LogReadURL             string                      `json:"logReadUrl"`
	ResourceAdditions      int                         `json:"resourceAdditions"`
	ResourceChanges        int                         `json:"resourceChanges"`
	ResourceDestructions   int                         `json:"resourceDestructions"`
	ResourceImports        int                         `json:"resourceImports"`
	StatusTimestamps       *planStatusTimestampsOutput `json:"statusTimestamps,omitempty"`
}

func (v *PlanShowView) Render(plan *tfe.Plan) error {
	if v.IsJSON() {
		var timestamps *planStatusTimestampsOutput
		if plan.StatusTimestamps != nil {
			timestamps = &planStatusTimestampsOutput{
				QueuedAt:        FormatDateTime(plan.StatusTimestamps.QueuedAt),
				StartedAt:       FormatDateTime(plan.StatusTimestamps.StartedAt),
				FinishedAt:      FormatDateTime(plan.StatusTimestamps.FinishedAt),
				CanceledAt:      FormatDateTime(plan.StatusTimestamps.CanceledAt),
				ErroredAt:       FormatDateTime(plan.StatusTimestamps.ErroredAt),
				ForceCanceledAt: FormatDateTime(plan.StatusTimestamps.ForceCanceledAt),
			}
		}
		return v.Output().RenderJSON(planShowOutput{
			ID:                     plan.ID,
			Status:                 string(plan.Status),
			HasChanges:             plan.HasChanges,
			GeneratedConfiguration: plan.GeneratedConfiguration,
			LogReadURL:             plan.LogReadURL,
			ResourceAdditions:      plan.ResourceAdditions,
			ResourceChanges:        plan.ResourceChanges,
			ResourceDestructions:   plan.ResourceDestructions,
			ResourceImports:        plan.ResourceImports,
			StatusTimestamps:       timestamps,
		})
	}

	props := []PropertyPair{
		{Key: "ID", Value: plan.ID},
		{Key: "Status", Value: string(plan.Status)},
		{Key: "Has Changes", Value: plan.HasChanges},
		{Key: "Generated Configuration", Value: plan.GeneratedConfiguration},
		{Key: "Resource Additions", Value: plan.ResourceAdditions},
		{Key: "Resource Changes", Value: plan.ResourceChanges},
		{Key: "Resource Destructions", Value: plan.ResourceDestructions},
		{Key: "Resource Imports", Value: plan.ResourceImports},
		// {Key: "Log Read URL", Value: plan.LogReadURL}, #Lets only show this in JSON
	}

	if err := v.Output().RenderProperties(props); err != nil {
		return err
	}

	// Render status timestamps as indented tags
	if plan.StatusTimestamps != nil {
		statusTags := []PropertyPair{}
		if !plan.StatusTimestamps.QueuedAt.IsZero() {
			statusTags = append(statusTags, PropertyPair{Key: "Queued At", Value: FormatDateTime(plan.StatusTimestamps.QueuedAt)})
		}
		if !plan.StatusTimestamps.StartedAt.IsZero() {
			statusTags = append(statusTags, PropertyPair{Key: "Started At", Value: FormatDateTime(plan.StatusTimestamps.StartedAt)})
		}
		if !plan.StatusTimestamps.FinishedAt.IsZero() {
			statusTags = append(statusTags, PropertyPair{Key: "Finished At", Value: FormatDateTime(plan.StatusTimestamps.FinishedAt)})
		}
		if !plan.StatusTimestamps.CanceledAt.IsZero() {
			statusTags = append(statusTags, PropertyPair{Key: "Canceled At", Value: FormatDateTime(plan.StatusTimestamps.CanceledAt)})
		}
		if !plan.StatusTimestamps.ErroredAt.IsZero() {
			statusTags = append(statusTags, PropertyPair{Key: "Errored At", Value: FormatDateTime(plan.StatusTimestamps.ErroredAt)})
		}
		if !plan.StatusTimestamps.ForceCanceledAt.IsZero() {
			statusTags = append(statusTags, PropertyPair{Key: "Force Canceled At", Value: FormatDateTime(plan.StatusTimestamps.ForceCanceledAt)})
		}
		if len(statusTags) > 0 {
			return v.Output().RenderTags("Statuses", statusTags)
		}
	}

	return nil
}
