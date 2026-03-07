// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

func orgColumns(width int) []column {
	idW := 22
	emailW := 30
	nameW := width - idW - emailW - 10 // 2(cursor) + 3×2(col padding) + some slack
	if nameW < 20 {
		nameW = 20
	}
	return []column{
		{name: "NAME", width: nameW},
		{name: "EMAIL", width: emailW},
		{name: "ID", width: idW},
	}
}

func filteredOrgs(m Model) []*tfe.Organization {
	if m.orgFilter == "" {
		return m.orgs
	}
	f := strings.ToLower(m.orgFilter)
	var out []*tfe.Organization
	for _, o := range m.orgs {
		if strings.Contains(strings.ToLower(o.Name), f) ||
			strings.Contains(strings.ToLower(o.Email), f) ||
			strings.Contains(strings.ToLower(o.ExternalID), f) {
			out = append(out, o)
		}
	}
	return out
}

func (m Model) renderOrgsContent() string {
	cols := orgColumns(m.mainWidth())
	visible := m.orgVisibleRows()
	filtered := filteredOrgs(m)

	var lines []string
	if m.orgFilter != "" || m.orgFiltering {
		lines = append(lines, m.renderFilterBar(m.orgFilter, m.orgFiltering))
	}
	lines = append(lines, m.renderTableHeader(cols))
	lines = append(lines, m.renderTableDivider())

	if len(filtered) == 0 {
		lines = append(lines, contentPlaceholderStyle.Width(m.mainWidth()).Render("  No organizations found."))
	} else {
		end := m.orgOffset + visible
		if end > len(filtered) {
			end = len(filtered)
		}
		for i := m.orgOffset; i < end; i++ {
			o := filtered[i]
			lines = append(lines, m.renderTableRow(i == m.orgCursor, []string{o.Name, o.Email, o.ExternalID}, cols))
		}
	}

	for len(lines) < m.contentHeight() {
		lines = append(lines, contentStyle.Width(m.mainWidth()).Render(""))
	}
	return strings.Join(lines[:m.contentHeight()], "\n")
}
