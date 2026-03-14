// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"image/color"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

func variableColumns(width int) []column {
	keyW := 24 // variable names are identifiers — cap at 24
	catW := 10 // "terraform" or "env"
	senW := 9  // "SENSITIVE" header; values are "yes" / "no"
	valW := width - keyW - catW - senW - 10 // 2(cursor) + 4×2(col padding)
	if valW < 5 {
		valW = 5
	}
	return []column{
		{name: "KEY", width: keyW},
		{name: "CATEGORY", width: catW},
		{name: "SENSITIVE", width: senW},
		{name: "VALUE", width: valW},
	}
}

// variableValue returns the display value for a variable, masking sensitive ones.
// Embedded newlines are collapsed to ↵ so the value fits on a single table row.
func variableValue(v *tfe.Variable) string {
	if v.Sensitive {
		return "••••••••"
	}
	return strings.ReplaceAll(v.Value, "\n", "↵")
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
	cols := variableColumns(m.innerWidth())
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
		lines = append(lines, contentPlaceholderStyle.Width(m.innerWidth()).Render("  No variables found."))
	} else {
		end := m.varOffset + visible
		if end > len(filtered) {
			end = len(filtered)
		}
		for i := m.varOffset; i < end; i++ {
			v := filtered[i]
			// KEY(0) CATEGORY(1) SENSITIVE(2) VALUE(3)
			cellFgs := []color.Color{nil, categoryFg(v), nil, nil}
			lines = append(lines, m.renderTableRowWithCellStyles(i == m.varCursor, []string{
				v.Key,
				string(v.Category),
				sensitiveStr(v),
				variableValue(v),
			}, cols, cellFgs))
		}
	}

	for len(lines) < m.contentHeight() {
		lines = append(lines, contentStyle.Width(m.innerWidth()).Render(""))
	}
	return strings.Join(lines[:m.contentHeight()], "\n")
}
