// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// buildSVDetailSections assembles the sections shown in the state version detail view.
func buildSVDetailSections(sv *tfe.StateVersion) []wsDetailSection {
	// ── General ──────────────────────────────────────────────────────────────
	general := wsDetailSection{title: "General"}
	general.rows = []wsDetailRow{
		{"ID", sv.ID},
		{"Serial", fmt.Sprintf("%d", sv.Serial)},
		{"Status", string(sv.Status)},
	}
	if sv.StateVersion > 0 {
		general.rows = append(general.rows, wsDetailRow{"State Version", fmt.Sprintf("%d", sv.StateVersion)})
	}
	if !sv.CreatedAt.IsZero() {
		general.rows = append(general.rows, wsDetailRow{"Created", sv.CreatedAt.UTC().Format("2006-01-02 15:04 UTC")})
	}

	// ── Terraform ─────────────────────────────────────────────────────────────
	terraform := wsDetailSection{title: "Terraform"}
	terraform.rows = []wsDetailRow{
		{"Resources Processed", boolYesNo(sv.ResourcesProcessed)},
	}
	if sv.TerraformVersion != "" {
		terraform.rows = append([]wsDetailRow{{"Terraform Version", sv.TerraformVersion}}, terraform.rows...)
	}

	sections := []wsDetailSection{general, terraform}

	// ── VCS ───────────────────────────────────────────────────────────────────
	if sv.VCSCommitSHA != "" {
		vcs := wsDetailSection{title: "VCS"}
		short := sv.VCSCommitSHA
		if len(short) > 12 {
			short = short[:12]
		}
		vcs.rows = []wsDetailRow{
			{"Commit SHA", short},
		}
		if sv.VCSCommitURL != "" {
			vcs.rows = append(vcs.rows, wsDetailRow{"Commit URL", sv.VCSCommitURL})
		}
		sections = append(sections, vcs)
	}

	return sections
}

// renderStateVersionDetailContent renders the full detail view for the selected state version.
func (m Model) renderStateVersionDetailContent() string {
	h := m.contentHeight()
	if m.selectedSV == nil {
		lines := make([]string, h)
		for i := range lines {
			lines[i] = contentStyle.Width(m.mainWidth()).Render("")
		}
		return strings.Join(lines, "\n")
	}

	sections := buildSVDetailSections(m.selectedSV)

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
	start := m.svDetScroll
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
