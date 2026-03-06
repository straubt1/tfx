// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	tfe "github.com/hashicorp/go-tfe"
)

func runColumns(width int) []column {
	idW := 30
	statusW := 22
	createdW := 10
	msgW := width - idW - statusW - createdW - 10 // 2(cursor) + 4×2(col padding)
	if msgW < 15 {
		msgW = 15
	}
	return []column{
		{name: "STATUS", width: statusW},
		{name: "CREATED", width: createdW},
		{name: "MESSAGE", width: msgW},
		{name: "ID", width: idW},
	}
}

// runStatusFg returns a foreground color for a run status, or nil for the row default.
func runStatusFg(status tfe.RunStatus) color.Color {
	switch status {
	case tfe.RunApplied, tfe.RunPlannedAndFinished:
		return colorSuccess
	case tfe.RunErrored:
		return colorError
	case tfe.RunPlanning, tfe.RunApplying, tfe.RunFetching, tfe.RunQueuing,
		tfe.RunPlanQueued, tfe.RunApplyQueued:
		return colorLoading
	case tfe.RunPlanned, tfe.RunPolicyChecked, tfe.RunCostEstimated:
		return colorAccent
	case tfe.RunCanceled, tfe.RunDiscarded:
		return colorDim
	default:
		return nil // use row default foreground
	}
}

// relativeTime formats a time.Time as a human-readable relative duration.
func relativeTime(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	}
}

func filteredRuns(m Model) []*tfe.Run {
	if m.runFilter == "" {
		return m.runs
	}
	f := strings.ToLower(m.runFilter)
	var out []*tfe.Run
	for _, r := range m.runs {
		if strings.Contains(strings.ToLower(string(r.Status)), f) ||
			strings.Contains(strings.ToLower(r.Message), f) ||
			strings.Contains(strings.ToLower(r.ID), f) {
			out = append(out, r)
		}
	}
	return out
}

func (m Model) renderRunsContent() string {
	cols := runColumns(m.width)
	visible := m.runVisibleRows()
	filtered := filteredRuns(m)

	var lines []string
	lines = append(lines, m.renderTableHeader(cols))
	lines = append(lines, m.renderTableDivider())

	if len(filtered) == 0 {
		lines = append(lines, contentPlaceholderStyle.Width(m.width).Render("  No runs found."))
	} else {
		end := m.runOffset + visible
		if end > len(filtered) {
			end = len(filtered)
		}
		for i := m.runOffset; i < end; i++ {
			r := filtered[i]
			cells := []string{
				string(r.Status),
				relativeTime(r.CreatedAt),
				r.Message,
				r.ID,
			}
			// Status cell (index 0) gets a color based on run outcome.
			cellFgs := []color.Color{runStatusFg(r.Status), nil, nil, nil}
			lines = append(lines, m.renderTableRowWithCellStyles(i == m.runCursor, cells, cols, cellFgs))
		}
	}

	if m.runFilter != "" || m.runFiltering {
		lines = append(lines, m.renderFilterBar(m.runFilter, m.runFiltering))
	}

	for len(lines) < m.contentHeight() {
		lines = append(lines, contentStyle.Width(m.width).Render(""))
	}
	return strings.Join(lines[:m.contentHeight()], "\n")
}
