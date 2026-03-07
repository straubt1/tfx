// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// buildCVDetailSections assembles the sections shown in the config version detail view.
func buildCVDetailSections(cv *tfe.ConfigurationVersion) []wsDetailSection {
	// ── General ──────────────────────────────────────────────────────────────
	general := wsDetailSection{title: "General"}
	general.rows = []wsDetailRow{
		{"ID", cv.ID},
		{"Status", string(cv.Status)},
	}
	if cv.Source != "" {
		general.rows = append(general.rows, wsDetailRow{"Source", string(cv.Source)})
	}
	if cv.StatusTimestamps != nil && !cv.StatusTimestamps.FinishedAt.IsZero() {
		general.rows = append(general.rows, wsDetailRow{"Finished At", cv.StatusTimestamps.FinishedAt.UTC().Format("2006-01-02 15:04 UTC")})
	}

	// ── Settings ──────────────────────────────────────────────────────────────
	settings := wsDetailSection{title: "Settings"}
	settings.rows = []wsDetailRow{
		{"Speculative", boolYesNo(cv.Speculative)},
		{"Auto Queue Runs", boolYesNo(cv.AutoQueueRuns)},
		{"Provisional", boolYesNo(cv.Provisional)},
	}

	sections := []wsDetailSection{general, settings}

	// ── Timestamps ────────────────────────────────────────────────────────────
	if cv.StatusTimestamps != nil {
		ts := cv.StatusTimestamps
		timestamps := wsDetailSection{title: "Timestamps"}
		if !ts.QueuedAt.IsZero() {
			timestamps.rows = append(timestamps.rows, wsDetailRow{"Queued At", ts.QueuedAt.UTC().Format("2006-01-02 15:04 UTC")})
		}
		if !ts.StartedAt.IsZero() {
			timestamps.rows = append(timestamps.rows, wsDetailRow{"Started At", ts.StartedAt.UTC().Format("2006-01-02 15:04 UTC")})
		}
		if !ts.FinishedAt.IsZero() {
			timestamps.rows = append(timestamps.rows, wsDetailRow{"Finished At", ts.FinishedAt.UTC().Format("2006-01-02 15:04 UTC")})
		}
		if len(timestamps.rows) > 0 {
			sections = append(sections, timestamps)
		}
	}

	// ── VCS ───────────────────────────────────────────────────────────────────
	if ia := cv.IngressAttributes; ia != nil {
		vcs := wsDetailSection{title: "VCS"}
		if ia.Branch != "" {
			vcs.rows = append(vcs.rows, wsDetailRow{"Branch", ia.Branch})
		}
		if ia.CommitSHA != "" {
			short := ia.CommitSHA
			if len(short) > 12 {
				short = short[:12]
			}
			vcs.rows = append(vcs.rows, wsDetailRow{"Commit SHA", short})
		}
		if ia.CommitMessage != "" {
			vcs.rows = append(vcs.rows, wsDetailRow{"Commit Message", ia.CommitMessage})
		}
		if ia.SenderUsername != "" {
			vcs.rows = append(vcs.rows, wsDetailRow{"Sender", ia.SenderUsername})
		}
		if ia.IsPullRequest {
			vcs.rows = append(vcs.rows, wsDetailRow{"Pull Request", fmt.Sprintf("#%d", ia.PullRequestNumber)})
			if ia.PullRequestURL != "" {
				vcs.rows = append(vcs.rows, wsDetailRow{"PR URL", ia.PullRequestURL})
			}
		}
		if len(vcs.rows) > 0 {
			sections = append(sections, vcs)
		}
	}

	return sections
}

// renderConfigVersionDetailContent renders the full detail view for the selected config version.
func (m Model) renderConfigVersionDetailContent() string {
	h := m.contentHeight()
	if m.selectedCV == nil {
		lines := make([]string, h)
		for i := range lines {
			lines[i] = contentStyle.Width(m.mainWidth()).Render("")
		}
		return strings.Join(lines, "\n")
	}

	sections := buildCVDetailSections(m.selectedCV)

	var all []string
	all = append(all, contentStyle.Width(m.mainWidth()).Render("")) // top padding

	for si, sec := range sections {
		all = append(all, m.renderDetailSectionHeader(sec.title))
		for _, row := range sec.rows {
			all = append(all, m.renderDetailKV(row.label, row.value))
		}
		if si < len(sections)-1 {
			all = append(all, contentStyle.Width(m.mainWidth()).Render(""))
		}
	}
	all = append(all, contentStyle.Width(m.mainWidth()).Render("")) // bottom padding

	// Clamp scroll and slice visible window.
	maxScroll := len(all) - h
	if maxScroll < 0 {
		maxScroll = 0
	}
	start := m.cvDetScroll
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
		out[i] = contentStyle.Width(m.mainWidth()).Render("")
	}
	return strings.Join(out, "\n")
}
