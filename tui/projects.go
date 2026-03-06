// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

func projectColumns(width int) []column {
	idW := 30
	nameW := width - idW - 6 // 2 (cursor) + 2 (name pad) + 2 (id pad)
	if nameW < 20 {
		nameW = 20
	}
	return []column{
		{name: "NAME", width: nameW},
		{name: "ID", width: idW},
	}
}

func filteredProjects(m Model) []*tfe.Project {
	if m.projFilter == "" {
		return m.projects
	}
	f := strings.ToLower(m.projFilter)
	var out []*tfe.Project
	for _, p := range m.projects {
		if strings.Contains(strings.ToLower(p.Name), f) {
			out = append(out, p)
		}
	}
	return out
}

func (m Model) renderProjectsContent() string {
	cols := projectColumns(m.width)
	visible := m.projVisibleRows()
	filtered := filteredProjects(m)

	var lines []string
	lines = append(lines, m.renderTableHeader(cols))
	lines = append(lines, m.renderTableDivider())

	if len(filtered) == 0 {
		lines = append(lines, contentPlaceholderStyle.Width(m.width).Render("  No projects found."))
	} else {
		end := m.projOffset + visible
		if end > len(filtered) {
			end = len(filtered)
		}
		for i := m.projOffset; i < end; i++ {
			p := filtered[i]
			lines = append(lines, m.renderTableRow(i == m.projCursor, []string{p.Name, p.ID}, cols))
		}
	}

	if m.projFilter != "" || m.projFiltering {
		lines = append(lines, m.renderFilterBar(m.projFilter, m.projFiltering))
	}

	for len(lines) < m.contentHeight() {
		lines = append(lines, contentStyle.Width(m.width).Render(""))
	}
	return strings.Join(lines[:m.contentHeight()], "\n")
}
