// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"fmt"
	"strings"
	"time"

	tfe "github.com/hashicorp/go-tfe"
)

const wsDetLabelWidth = 26

// wsDetailRow is a single label-value pair in the workspace detail view.
type wsDetailRow struct {
	label string
	value string
}

// wsDetailSection groups a set of detail rows under a named heading.
type wsDetailSection struct {
	title string
	rows  []wsDetailRow
}

// timestampWithRelative formats a time as "2006-01-02 15:04 UTC (Xd/h/m ago)".
func timestampWithRelative(t time.Time) string {
	return fmt.Sprintf("%s (%s)", t.UTC().Format("2006-01-02 15:04 UTC"), relativeTime(t))
}

// boolYesNo returns "yes" or "no" for a bool value.
func boolYesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

// buildWorkspaceDetailSections assembles the sections shown in the detail view.
func buildWorkspaceDetailSections(ws *tfe.Workspace) []wsDetailSection {
	// ── General ──────────────────────────────────────────────────────────────
	general := wsDetailSection{title: "General"}
	general.rows = []wsDetailRow{
		{"Name", ws.Name},
		{"ID", ws.ID},
	}
	if ws.Description != "" {
		general.rows = append(general.rows, wsDetailRow{"Description", ws.Description})
	}
	if !ws.CreatedAt.IsZero() {
		general.rows = append(general.rows, wsDetailRow{"Created", timestampWithRelative(ws.CreatedAt)})
	}

	// ── Configuration ─────────────────────────────────────────────────────────
	config := wsDetailSection{title: "Configuration"}
	config.rows = []wsDetailRow{
		{"Terraform Version", ws.TerraformVersion},
		{"Execution Mode", ws.ExecutionMode},
		{"Auto Apply", boolYesNo(ws.AutoApply)},
	}
	if ws.WorkingDirectory != "" {
		config.rows = append(config.rows, wsDetailRow{"Working Directory", ws.WorkingDirectory})
	}
	config.rows = append(config.rows,
		wsDetailRow{"Locked", boolYesNo(ws.Locked)},
		wsDetailRow{"Resource Count", fmt.Sprintf("%d", ws.ResourceCount)},
		wsDetailRow{"Allow Destroy Plan", boolYesNo(ws.AllowDestroyPlan)},
		wsDetailRow{"Speculative Enabled", boolYesNo(ws.SpeculativeEnabled)},
		wsDetailRow{"Queue All Runs", boolYesNo(ws.QueueAllRuns)},
		wsDetailRow{"File Triggers", boolYesNo(ws.FileTriggersEnabled)},
		wsDetailRow{"Global Remote State", boolYesNo(ws.GlobalRemoteState)},
	)

	sections := []wsDetailSection{general, config}

	// ── VCS ───────────────────────────────────────────────────────────────────
	if ws.VCSRepo != nil {
		vcs := wsDetailSection{title: "VCS"}
		vcs.rows = []wsDetailRow{
			{"Repository", ws.VCSRepo.Identifier},
		}
		if ws.VCSRepo.Branch != "" {
			vcs.rows = append(vcs.rows, wsDetailRow{"Branch", ws.VCSRepo.Branch})
		}
		if ws.VCSRepo.ServiceProvider != "" {
			vcs.rows = append(vcs.rows, wsDetailRow{"Provider", ws.VCSRepo.ServiceProvider})
		}
		if ws.VCSRepo.RepositoryHTTPURL != "" {
			vcs.rows = append(vcs.rows, wsDetailRow{"URL", ws.VCSRepo.RepositoryHTTPURL})
		}
		sections = append(sections, vcs)
	}

	// ── Stats ─────────────────────────────────────────────────────────────────
	stats := wsDetailSection{title: "Stats"}
	stats.rows = []wsDetailRow{
		{"Runs", fmt.Sprintf("%d", ws.RunsCount)},
		{"Run Failures", fmt.Sprintf("%d", ws.RunFailures)},
	}
	if len(ws.TagNames) > 0 {
		stats.rows = append(stats.rows, wsDetailRow{"Tags", strings.Join(ws.TagNames, ", ")})
	}
	sections = append(sections, stats)

	return sections
}

// renderWorkspaceDetailContent renders the full detail view for the selected workspace.
func (m Model) renderWorkspaceDetailContent() string {
	h := m.contentHeight()
	if m.selectedWS == nil {
		lines := make([]string, h)
		for i := range lines {
			lines[i] = contentStyle.Width(m.innerWidth()).Render("")
		}
		return strings.Join(lines, "\n")
	}

	sections := buildWorkspaceDetailSections(m.selectedWS)

	// Build a flat list of all rendered lines.
	var all []string
	all = append(all, contentStyle.Width(m.innerWidth()).Render("")) // top padding

	for si, sec := range sections {
		all = append(all, m.renderDetailSectionHeader(sec.title))
		for _, row := range sec.rows {
			all = append(all, m.renderDetailKV(row.label, row.value))
		}
		if si < len(sections)-1 {
			all = append(all, contentStyle.Width(m.innerWidth()).Render("")) // blank line between sections
		}
	}
	all = append(all, contentStyle.Width(m.innerWidth()).Render("")) // bottom padding

	// Clamp the scroll offset to valid range.
	maxScroll := len(all) - h
	if maxScroll < 0 {
		maxScroll = 0
	}
	start := m.wsDetScroll
	if start > maxScroll {
		start = maxScroll
	}

	// Slice the visible window, then pad to content height.
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

// renderDetailSectionHeader renders a full-width section divider line.
//
//	  ── Title ───────────────────────────────────
func (m Model) renderDetailSectionHeader(title string) string {
	prefix := "  ── " + title + " "
	n := m.innerWidth() - len([]rune(prefix))
	if n < 4 {
		n = 4
	}
	text := prefix + strings.Repeat("─", n)
	return contentTitleStyle.Width(m.innerWidth()).Render(text)
}

// renderDetailKV renders a single label-value row with a fixed-width label column.
func (m Model) renderDetailKV(label, value string) string {
	maxValueWidth := m.innerWidth() - wsDetLabelWidth - 2
	if maxValueWidth < 10 {
		maxValueWidth = 10
	}
	labelPart := detailLabelStyle.Width(wsDetLabelWidth).Render("  " + label)
	valuePart := contentStyle.Render(truncateStr(value, maxValueWidth))
	return m.padContent(labelPart+valuePart, contentStyle)
}
