// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

func workspaceColumns(width int) []column {
	idW := 30
	statusW := 22
	nameW := width - idW - statusW - 8 // 2 (cursor) + 2+2+2 (padding between cols)
	if nameW < 20 {
		nameW = 20
	}
	return []column{
		{name: "NAME", width: nameW},
		{name: "STATUS", width: statusW},
		{name: "ID", width: idW},
	}
}

func workspaceStatus(ws *tfe.Workspace) string {
	if ws.CurrentRun == nil {
		return "—"
	}
	return string(ws.CurrentRun.Status)
}

func filteredWorkspaces(m Model) []*tfe.Workspace {
	if m.wsFilter == "" {
		return m.workspaces
	}
	f := strings.ToLower(m.wsFilter)
	var out []*tfe.Workspace
	for _, ws := range m.workspaces {
		if strings.Contains(strings.ToLower(ws.Name), f) {
			out = append(out, ws)
		}
	}
	return out
}

func (m Model) renderWorkspacesContent() string {
	cols := workspaceColumns(m.width)
	visible := m.wsVisibleRows()
	filtered := filteredWorkspaces(m)

	var lines []string
	if m.wsFilter != "" || m.wsFiltering {
		lines = append(lines, m.renderFilterBar(m.wsFilter, m.wsFiltering))
	}
	lines = append(lines, m.renderTableHeader(cols))
	lines = append(lines, m.renderTableDivider())

	if len(filtered) == 0 {
		lines = append(lines, contentPlaceholderStyle.Width(m.width).Render("  No workspaces found."))
	} else {
		end := m.wsOffset + visible
		if end > len(filtered) {
			end = len(filtered)
		}
		for i := m.wsOffset; i < end; i++ {
			ws := filtered[i]
			lines = append(lines, m.renderTableRow(i == m.wsCursor, []string{ws.Name, workspaceStatus(ws), ws.ID}, cols))
		}
	}

	for len(lines) < m.contentHeight() {
		lines = append(lines, contentStyle.Width(m.width).Render(""))
	}
	return strings.Join(lines[:m.contentHeight()], "\n")
}
