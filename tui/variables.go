// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

func variableColumns(width int) []column {
	idW := 30
	catW := 12  // "terraform" or "env"
	valW := 20
	keyW := width - idW - catW - valW - 10 // 2(cursor) + 4×2(col padding)
	if keyW < 15 {
		keyW = 15
	}
	return []column{
		{name: "KEY", width: keyW},
		{name: "CATEGORY", width: catW},
		{name: "VALUE", width: valW},
		{name: "ID", width: idW},
	}
}

// variableValue returns the display value for a variable, masking sensitive ones.
func variableValue(v *tfe.Variable) string {
	if v.Sensitive {
		return "••••••••"
	}
	return v.Value
}

func filteredVariables(m Model) []*tfe.Variable {
	if m.varFilter == "" {
		return m.variables
	}
	f := strings.ToLower(m.varFilter)
	var out []*tfe.Variable
	for _, v := range m.variables {
		if strings.Contains(strings.ToLower(v.Key), f) ||
			strings.Contains(strings.ToLower(string(v.Category)), f) ||
			strings.Contains(strings.ToLower(v.ID), f) {
			out = append(out, v)
		}
	}
	return out
}

func (m Model) renderVariablesContent() string {
	cols := variableColumns(m.width)
	visible := m.varVisibleRows()
	filtered := filteredVariables(m)

	var lines []string
	lines = append(lines, m.renderTableHeader(cols))
	lines = append(lines, m.renderTableDivider())

	if len(filtered) == 0 {
		lines = append(lines, contentPlaceholderStyle.Width(m.width).Render("  No variables found."))
	} else {
		end := m.varOffset + visible
		if end > len(filtered) {
			end = len(filtered)
		}
		for i := m.varOffset; i < end; i++ {
			v := filtered[i]
			lines = append(lines, m.renderTableRow(i == m.varCursor, []string{
				v.Key,
				string(v.Category),
				variableValue(v),
				v.ID,
			}, cols))
		}
	}

	if m.varFilter != "" || m.varFiltering {
		lines = append(lines, m.renderFilterBar(m.varFilter, m.varFiltering))
	}

	for len(lines) < m.contentHeight() {
		lines = append(lines, contentStyle.Width(m.width).Render(""))
	}
	return strings.Join(lines[:m.contentHeight()], "\n")
}
