// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"image/color"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

func variableColumns(width int) []column {
	idW := 30
	catW := 12  // "terraform" or "env"
	senW := 9   // "SENSITIVE" header; values are "yes" / "no"
	keyW := 30  // variable names are identifiers — cap at 30, rarely longer
	valW := width - idW - catW - senW - keyW - 12 // 2(cursor) + 5×2(col padding)
	if valW < 15 {
		valW = 15
	}
	return []column{
		{name: "KEY", width: keyW},
		{name: "VALUE", width: valW},
		{name: "SENSITIVE", width: senW},
		{name: "CATEGORY", width: catW},
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

// sensitiveStr returns a human-readable yes/no for v.Sensitive.
func sensitiveStr(v *tfe.Variable) string {
	if v.Sensitive {
		return "yes"
	}
	return "no"
}

// categoryFg returns a foreground color for the variable category.
func categoryFg(v *tfe.Variable) color.Color {
	switch v.Category {
	case tfe.CategoryTerraform:
		return colorAccent // blue — native terraform type
	case tfe.CategoryEnv:
		return colorSuccess // green — from environment
	default:
		return nil
	}
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
	lines = append(lines, m.renderWorkspaceTabStrip())
	if m.varFilter != "" || m.varFiltering {
		lines = append(lines, m.renderFilterBar(m.varFilter, m.varFiltering))
	}
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
			// CATEGORY is at index 3: KEY(0) VALUE(1) SENSITIVE(2) CATEGORY(3) ID(4)
			cellFgs := []color.Color{nil, nil, nil, categoryFg(v), nil}
			lines = append(lines, m.renderTableRowWithCellStyles(i == m.varCursor, []string{
				v.Key,
				variableValue(v),
				sensitiveStr(v),
				string(v.Category),
				v.ID,
			}, cols, cellFgs))
		}
	}

	for len(lines) < m.contentHeight() {
		lines = append(lines, contentStyle.Width(m.width).Render(""))
	}
	return strings.Join(lines[:m.contentHeight()], "\n")
}
