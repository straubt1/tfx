// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	"fmt"

	"github.com/fatih/color"
	tfe "github.com/hashicorp/go-tfe"
)

type PlanCreateView struct{ *BaseView }

type PlanCreateRenderOptions struct {
	RunID        string
	PlanID       string
	Hostname     string
	Organization string
	Workspace    string
}

func NewPlanCreateView() *PlanCreateView { return &PlanCreateView{NewBaseView()} }

type planCreateOutput struct {
	ID               string                      `json:"id"`
	Status           string                      `json:"status"`
	LogReadURL       string                      `json:"logReadUrl"`
	StatusTimestamps *planStatusTimestampsOutput `json:"statusTimestamps,omitempty"`
}

func (v *PlanCreateView) Render(plan *tfe.Plan, opts *PlanCreateRenderOptions) error {
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
		return v.Output().RenderJSON(planCreateOutput{
			ID:               plan.ID,
			Status:           string(plan.Status),
			LogReadURL:       plan.LogReadURL,
			StatusTimestamps: timestamps,
		})
	}

	props := []PropertyPair{
		{Key: "ID", Value: plan.ID},
		{Key: "Status", Value: string(plan.Status)},
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
			if err := v.Output().RenderTags("Statuses", statusTags); err != nil {
				return err
			}
		}
	}

	// Add the run and navigation info if options provided
	if opts != nil {
		v.Output().Message("Run ID: %s", color.BlueString(opts.RunID))
		v.Output().Message("Plan ID: %s", color.BlueString(opts.PlanID))
		v.Output().Message("Navigate: %s", color.BlueString(fmt.Sprintf("https://%s/app/%s/workspaces/%s/runs/%s", opts.Hostname, opts.Organization, opts.Workspace, opts.RunID)))
	}

	return nil
}
