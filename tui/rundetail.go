// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// buildRunDetailSections assembles the sections shown in the run detail view.
func buildRunDetailSections(run *tfe.Run) []wsDetailSection {
	// ── General ──────────────────────────────────────────────────────────────
	general := wsDetailSection{title: "General"}
	general.rows = []wsDetailRow{
		{"ID", run.ID},
		{"Status", string(run.Status)},
	}
	if run.Message != "" {
		general.rows = append(general.rows, wsDetailRow{"Message", run.Message})
	}
	if run.Source != "" {
		general.rows = append(general.rows, wsDetailRow{"Source", string(run.Source)})
	}
	if run.TriggerReason != "" {
		general.rows = append(general.rows, wsDetailRow{"Trigger Reason", run.TriggerReason})
	}
	if run.TerraformVersion != "" {
		general.rows = append(general.rows, wsDetailRow{"Terraform Version", run.TerraformVersion})
	}
	if !run.CreatedAt.IsZero() {
		general.rows = append(general.rows, wsDetailRow{"Created", run.CreatedAt.UTC().Format("2006-01-02 15:04 UTC")})
	}

	// ── Flags ─────────────────────────────────────────────────────────────────
	flags := wsDetailSection{title: "Flags"}
	flags.rows = []wsDetailRow{
		{"Auto Apply", boolYesNo(run.AutoApply)},
		{"Is Destroy", boolYesNo(run.IsDestroy)},
		{"Plan Only", boolYesNo(run.PlanOnly)},
		{"Allow Empty Apply", boolYesNo(run.AllowEmptyApply)},
		{"Refresh Only", boolYesNo(run.RefreshOnly)},
		{"Has Changes", boolYesNo(run.HasChanges)},
	}

	sections := []wsDetailSection{general, flags}

	// ── Plan ──────────────────────────────────────────────────────────────────
	if run.Plan != nil {
		plan := wsDetailSection{title: "Plan"}
		plan.rows = []wsDetailRow{
			{"Plan ID", run.Plan.ID},
			{"Status", string(run.Plan.Status)},
			{"Has Changes", boolYesNo(run.Plan.HasChanges)},
			{"Additions", fmt.Sprintf("%d", run.Plan.ResourceAdditions)},
			{"Changes", fmt.Sprintf("%d", run.Plan.ResourceChanges)},
			{"Destructions", fmt.Sprintf("%d", run.Plan.ResourceDestructions)},
			{"Imports", fmt.Sprintf("%d", run.Plan.ResourceImports)},
		}
		sections = append(sections, plan)
	}

	// ── Apply ─────────────────────────────────────────────────────────────────
	if run.Apply != nil {
		apply := wsDetailSection{title: "Apply"}
		apply.rows = []wsDetailRow{
			{"Apply ID", run.Apply.ID},
			{"Status", string(run.Apply.Status)},
		}
		sections = append(sections, apply)
	}

	// ── VCS ───────────────────────────────────────────────────────────────────
	if run.ConfigurationVersion != nil && run.ConfigurationVersion.IngressAttributes != nil {
		ia := run.ConfigurationVersion.IngressAttributes
		vcs := wsDetailSection{title: "VCS"}
		if ia.CommitSHA != "" {
			short := ia.CommitSHA
			if len(short) > 12 {
				short = short[:12]
			}
			vcs.rows = append(vcs.rows, wsDetailRow{"Commit SHA", short})
		}
		if ia.Branch != "" {
			vcs.rows = append(vcs.rows, wsDetailRow{"Branch", ia.Branch})
		}
		if ia.CommitMessage != "" {
			vcs.rows = append(vcs.rows, wsDetailRow{"Commit Message", ia.CommitMessage})
		}
		if ia.SenderUsername != "" {
			vcs.rows = append(vcs.rows, wsDetailRow{"Sender", ia.SenderUsername})
		}
		if ia.CommitURL != "" {
			vcs.rows = append(vcs.rows, wsDetailRow{"Commit URL", ia.CommitURL})
		}
		if ia.IsPullRequest {
			vcs.rows = append(vcs.rows, wsDetailRow{"Pull Request", fmt.Sprintf("#%d", ia.PullRequestNumber)})
		}
		if len(vcs.rows) > 0 {
			sections = append(sections, vcs)
		}
	}

	return sections
}

// renderRunDetailContent renders the full detail view for the selected run.
func (m Model) renderRunDetailContent() string {
	h := m.contentHeight()
	if m.selectedRun == nil {
		lines := make([]string, h)
		for i := range lines {
			lines[i] = contentStyle.Width(m.innerWidth()).Render("")
		}
		return strings.Join(lines, "\n")
	}

	sections := buildRunDetailSections(m.selectedRun)

	var all []string
	all = append(all, contentStyle.Width(m.innerWidth()).Render("")) // top padding

	for si, sec := range sections {
		all = append(all, m.renderDetailSectionHeader(sec.title))
		for _, row := range sec.rows {
			all = append(all, m.renderDetailKV(row.label, row.value))
		}
		if si < len(sections)-1 {
			all = append(all, contentStyle.Width(m.innerWidth()).Render(""))
		}
	}
	all = append(all, contentStyle.Width(m.innerWidth()).Render("")) // bottom padding

	// Clamp scroll and slice visible window.
	maxScroll := len(all) - h
	if maxScroll < 0 {
		maxScroll = 0
	}
	start := m.runDetScroll
	if start > maxScroll {
		start = maxScroll
	}
	visible := all[start:]
	if len(visible) > h {
		visible = visible[:h]
	}
	out := make([]string, h)
	copy(out, visible)
	for i := len(visible); i < h; i++ {
		out[i] = contentStyle.Width(m.innerWidth()).Render("")
	}
	return strings.Join(out, "\n")
}
