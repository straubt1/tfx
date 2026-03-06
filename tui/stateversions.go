// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

func svColumns(width int) []column {
	idW := 30
	statusW := 15
	serialW := 8
	tfVerW := 12
	createdW := 10
	_ = width
	return []column{
		{name: "STATUS", width: statusW},
		{name: "SERIAL", width: serialW},
		{name: "TF VERSION", width: tfVerW},
		{name: "CREATED", width: createdW},
		{name: "ID", width: idW},
	}
}

func filteredStateVersions(m Model) []*tfe.StateVersion {
	if m.svFilter == "" {
		return m.stateVersions
	}
	f := strings.ToLower(m.svFilter)
	var out []*tfe.StateVersion
	for _, sv := range m.stateVersions {
		if strings.Contains(strings.ToLower(string(sv.Status)), f) ||
			strings.Contains(strings.ToLower(sv.TerraformVersion), f) ||
			strings.Contains(strings.ToLower(sv.ID), f) {
			out = append(out, sv)
		}
	}
	return out
}

func (m Model) renderStateVersionsContent() string {
	cols := svColumns(m.width)
	visible := m.svVisibleRows()
	filtered := filteredStateVersions(m)

	var lines []string
	lines = append(lines, m.renderTableHeader(cols))
	lines = append(lines, m.renderTableDivider())

	if len(filtered) == 0 {
		lines = append(lines, contentPlaceholderStyle.Width(m.width).Render("  No state versions found."))
	} else {
		end := m.svOffset + visible
		if end > len(filtered) {
			end = len(filtered)
		}
		for i := m.svOffset; i < end; i++ {
			sv := filtered[i]
			lines = append(lines, m.renderTableRow(i == m.svCursor, []string{
				string(sv.Status),
				fmt.Sprintf("%d", sv.Serial),
				sv.TerraformVersion,
				relativeTime(sv.CreatedAt),
				sv.ID,
			}, cols))
		}
	}

	if m.svFilter != "" || m.svFiltering {
		lines = append(lines, m.renderFilterBar(m.svFilter, m.svFiltering))
	}

	for len(lines) < m.contentHeight() {
		lines = append(lines, contentStyle.Width(m.width).Render(""))
	}
	return strings.Join(lines[:m.contentHeight()], "\n")
}
