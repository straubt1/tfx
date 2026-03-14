// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// buildVariableDetailSections assembles the sections shown in the variable detail view.
func buildVariableDetailSections(v *tfe.Variable) []wsDetailSection {
	// ── General ──────────────────────────────────────────────────────────────
	general := wsDetailSection{title: "General"}
	general.rows = []wsDetailRow{
		{"Key", v.Key},
		{"ID", v.ID},
		{"Category", string(v.Category)},
		{"HCL", boolYesNo(v.HCL)},
		{"Sensitive", boolYesNo(v.Sensitive)},
	}

	// ── Value ─────────────────────────────────────────────────────────────────
	value := wsDetailSection{title: "Value"}
	if v.Sensitive {
		value.rows = []wsDetailRow{
			{"Value", "***** (sensitive)"},
		}
	} else if v.Value != "" {
		// Split multi-line values into one row per line so embedded newlines
		// don't break the row-counting / scroll math.
		valueLines := strings.Split(v.Value, "\n")
		rows := make([]wsDetailRow, len(valueLines))
		rows[0] = wsDetailRow{"Value", valueLines[0]}
		for i, line := range valueLines[1:] {
			rows[i+1] = wsDetailRow{"", line}
		}
		value.rows = rows
	} else {
		value.rows = []wsDetailRow{
			{"Value", "(empty)"},
		}
	}

	sections := []wsDetailSection{general, value}

	// ── Description ───────────────────────────────────────────────────────────
	if v.Description != "" {
		desc := wsDetailSection{title: "Description"}
		desc.rows = []wsDetailRow{
			{"", v.Description},
		}
		sections = append(sections, desc)
	}

	return sections
}

// renderVariableDetailContent renders the full detail view for the selected variable.
func (m Model) renderVariableDetailContent() string {
	h := m.contentHeight()
	if m.selectedVar == nil {
		lines := make([]string, h)
		for i := range lines {
			lines[i] = contentStyle.Width(m.innerWidth()).Render("")
		}
		return strings.Join(lines, "\n")
	}

	sections := buildVariableDetailSections(m.selectedVar)

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
	start := m.varDetScroll
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
