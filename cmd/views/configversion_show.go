// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	tfe "github.com/hashicorp/go-tfe"
)

type ConfigVersionShowView struct{ *BaseView }

func NewConfigVersionShowView() *ConfigVersionShowView { return &ConfigVersionShowView{NewBaseView()} }

func (v *ConfigVersionShowView) Render(cv *tfe.ConfigurationVersion) error {
	if v.IsJSON() {
		out := map[string]interface{}{
			"id":          cv.ID,
			"status":      cv.Status,
			"speculative": cv.Speculative,
		}
		if cv.ErrorMessage != "" {
			out["errorMessage"] = cv.ErrorMessage
		}
		if cv.IngressAttributes != nil {
			out["repo"] = cv.IngressAttributes.Identifier
			out["branch"] = cv.IngressAttributes.Branch
			out["commit"] = cv.IngressAttributes.CommitSHA
			out["message"] = cv.IngressAttributes.CommitMessage
			out["link"] = cv.IngressAttributes.CommitURL
		}
		return v.Output().RenderJSON(out)
	}
	props := []PropertyPair{
		{Key: "ID", Value: cv.ID},
		{Key: "Status", Value: cv.Status},
		{Key: "Speculative", Value: cv.Speculative},
	}
	if cv.ErrorMessage != "" {
		props = append(props, PropertyPair{Key: "Error Message", Value: cv.ErrorMessage})
	}
	if err := v.Output().RenderProperties(props); err != nil {
		return err
	}
	tags := []PropertyPair{}
	if cv.IngressAttributes != nil {
		tags = append(tags,
			PropertyPair{Key: "Repo", Value: cv.IngressAttributes.Identifier},
			PropertyPair{Key: "Branch", Value: cv.IngressAttributes.Branch},
			PropertyPair{Key: "Commit", Value: cv.IngressAttributes.CommitSHA},
			PropertyPair{Key: "Message", Value: cv.IngressAttributes.CommitMessage},
			PropertyPair{Key: "Link", Value: cv.IngressAttributes.CommitURL},
		)
	}
	return v.Output().RenderTags("VCS", tags)
}
