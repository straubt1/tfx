// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// buildProjectDetailSections assembles the sections shown in the project detail view.
func buildProjectDetailSections(proj *tfe.Project) []wsDetailSection {
	// ── General ──────────────────────────────────────────────────────────────
	general := wsDetailSection{title: "General"}
	general.rows = []wsDetailRow{
		{"Name", proj.Name},
		{"ID", proj.ID},
	}
	if proj.Description != "" {
		general.rows = append(general.rows, wsDetailRow{"Description", proj.Description})
	}
	if proj.DefaultExecutionMode != "" {
		general.rows = append(general.rows, wsDetailRow{"Default Execution Mode", proj.DefaultExecutionMode})
	}

	return []wsDetailSection{general}
}

// renderProjectDetailContent renders the full detail view for the selected project.
func (m Model) renderProjectDetailContent() string {
	h := m.contentHeight()
	if m.selectedProj == nil {
		lines := make([]string, h)
		for i := range lines {
			lines[i] = contentStyle.Width(m.width).Render("")
		}
		return strings.Join(lines, "\n")
	}

	sections := buildProjectDetailSections(m.selectedProj)

	var all []string
	all = append(all, contentStyle.Width(m.width).Render("")) // top padding

	for si, sec := range sections {
		all = append(all, m.renderDetailSectionHeader(sec.title))
		for _, row := range sec.rows {
			all = append(all, m.renderDetailKV(row.label, row.value))
		}
		if si < len(sections)-1 {
			all = append(all, contentStyle.Width(m.width).Render(""))
		}
	}
	all = append(all, contentStyle.Width(m.width).Render("")) // bottom padding

	// Clamp scroll and slice visible window.
	maxScroll := len(all) - h
	if maxScroll < 0 {
		maxScroll = 0
	}
	start := m.projDetScroll
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
		out[i] = contentStyle.Width(m.width).Render("")
	}
	return strings.Join(out, "\n")
}
