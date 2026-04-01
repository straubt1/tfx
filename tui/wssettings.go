// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import "strings"

// renderWorkspaceSettingsContent renders the Settings tab for the selected
// workspace. It shows the same fields as the workspace detail view (d key)
// but embedded inside the tab strip rather than as a standalone view.
func (m Model) renderWorkspaceSettingsContent() string {
	tabStrip := m.renderWorkspaceTabStrip()

	if m.selectedWS == nil {
		h := m.contentHeight() - 1
		lines := make([]string, h)
		for i := range lines {
			lines[i] = contentStyle.Width(m.innerWidth()).Render("")
		}
		return tabStrip + "\n" + strings.Join(lines, "\n")
	}

	sections := buildWorkspaceDetailSections(m.selectedWS)

	// Inject "Last Updated" into the General section (index 0) using the
	// latest-change-at value fetched separately (not in tfe.Workspace).
	lastUpdated := "…"
	if m.wsLatestChange != nil {
		lastUpdated = timestampWithRelative(*m.wsLatestChange)
	}
	if len(sections) > 0 {
		sections[0].rows = append(sections[0].rows, wsDetailRow{"Last Updated", lastUpdated})
	}

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

	// Visible rows = content height minus tab strip row.
	h := m.contentHeight() - 1

	maxScroll := len(all) - h
	if maxScroll < 0 {
		maxScroll = 0
	}
	start := m.wsSettingsScroll
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

	return tabStrip + "\n" + strings.Join(out, "\n")
}
