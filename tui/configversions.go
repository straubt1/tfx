// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

func cvColumns(width int) []column {
	idW := 30
	statusW := 20
	sourceW := 15
	specW := 11 // "yes" / "no"
	// Use remaining width for nothing — columns are fixed for CVs.
	_ = width
	return []column{
		{name: "STATUS", width: statusW},
		{name: "SOURCE", width: sourceW},
		{name: "SPECULATIVE", width: specW},
		{name: "ID", width: idW},
	}
}

func filteredConfigVersions(m Model) []*tfe.ConfigurationVersion {
	if m.cvFilter == "" {
		return m.configVersions
	}
	f := strings.ToLower(m.cvFilter)
	var out []*tfe.ConfigurationVersion
	for _, cv := range m.configVersions {
		if strings.Contains(strings.ToLower(string(cv.Status)), f) ||
			strings.Contains(strings.ToLower(string(cv.Source)), f) ||
			strings.Contains(strings.ToLower(cv.ID), f) {
			out = append(out, cv)
		}
	}
	return out
}

func speculative(cv *tfe.ConfigurationVersion) string {
	if cv.Speculative {
		return "yes"
	}
	return "no"
}

func (m Model) renderConfigVersionsContent() string {
	cols := cvColumns(m.width)
	visible := m.cvVisibleRows()
	filtered := filteredConfigVersions(m)

	var lines []string
	lines = append(lines, m.renderWorkspaceTabStrip())
	if m.cvFilter != "" || m.cvFiltering {
		lines = append(lines, m.renderFilterBar(m.cvFilter, m.cvFiltering))
	}
	lines = append(lines, m.renderTableHeader(cols))
	lines = append(lines, m.renderTableDivider())

	if len(filtered) == 0 {
		lines = append(lines, contentPlaceholderStyle.Width(m.width).Render("  No configuration versions found."))
	} else {
		end := m.cvOffset + visible
		if end > len(filtered) {
			end = len(filtered)
		}
		for i := m.cvOffset; i < end; i++ {
			cv := filtered[i]
			lines = append(lines, m.renderTableRow(i == m.cvCursor, []string{
				string(cv.Status),
				string(cv.Source),
				speculative(cv),
				cv.ID,
			}, cols))
		}
	}

	for len(lines) < m.contentHeight() {
		lines = append(lines, contentStyle.Width(m.width).Render(""))
	}
	return strings.Join(lines[:m.contentHeight()], "\n")
}
