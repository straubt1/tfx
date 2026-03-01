// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	tfe "github.com/hashicorp/go-tfe"
)

type ConfigVersionListView struct{ *BaseView }

func NewConfigVersionListView() *ConfigVersionListView { return &ConfigVersionListView{NewBaseView()} }

type configVersionRow struct {
	ID          string `json:"id"`
	Speculative bool   `json:"speculative"`
	Status      string `json:"status"`
	Repo        string `json:"repo"`
	Branch      string `json:"branch"`
	Commit      string `json:"commit"`
	Message     string `json:"message"`
}

func (v *ConfigVersionListView) Render(items []*tfe.ConfigurationVersion) error {
	if v.IsJSON() {
		out := make([]configVersionRow, len(items))
		for i, cv := range items {
			repo, branch, commit, message := "", "", "", ""
			if cv.IngressAttributes != nil {
				repo = cv.IngressAttributes.Identifier
				branch = cv.IngressAttributes.Branch
				commit = cv.IngressAttributes.CommitSHA
				if len(commit) > 7 {
					commit = commit[:7]
				}
				message = cv.IngressAttributes.CommitMessage
			}
			out[i] = configVersionRow{
				ID:          cv.ID,
				Speculative: cv.Speculative,
				Status:      string(cv.Status),
				Repo:        repo,
				Branch:      branch,
				Commit:      commit,
				Message:     message,
			}
		}
		return v.Output().RenderJSON(out)
	}
	headers := []string{"Id", "Speculative", "Status", "Repo", "Branch", "Commit", "Message"}
	rows := make([][]interface{}, len(items))
	for i, cv := range items {
		repo, branch, commit, message := "", "", "", ""
		if cv.IngressAttributes != nil {
			repo = cv.IngressAttributes.Identifier
			branch = cv.IngressAttributes.Branch
			commit = cv.IngressAttributes.CommitSHA
			if len(commit) > 7 {
				commit = commit[:7]
			}
			message = cv.IngressAttributes.CommitMessage
		}
		rows[i] = []interface{}{cv.ID, cv.Speculative, cv.Status, repo, branch, commit, message}
	}
	return v.Output().RenderTable(headers, rows)
}
