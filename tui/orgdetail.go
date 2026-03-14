// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// buildOrgDetailSections assembles the sections shown in the org detail view.
func buildOrgDetailSections(org *tfe.Organization) []wsDetailSection {
	// ── General ──────────────────────────────────────────────────────────────
	general := wsDetailSection{title: "General"}
	general.rows = []wsDetailRow{
		{"Name", org.Name},
	}
	if org.ExternalID != "" {
		general.rows = append(general.rows, wsDetailRow{"External ID", org.ExternalID})
	}
	if org.Email != "" {
		general.rows = append(general.rows, wsDetailRow{"Email", org.Email})
	}
	if !org.CreatedAt.IsZero() {
		general.rows = append(general.rows, wsDetailRow{"Created", org.CreatedAt.UTC().Format("2006-01-02 15:04 UTC")})
	}

	// ── Settings ──────────────────────────────────────────────────────────────
	settings := wsDetailSection{title: "Settings"}
	settings.rows = []wsDetailRow{
		{"Default Execution Mode", org.DefaultExecutionMode},
		{"Cost Estimation", boolYesNo(org.CostEstimationEnabled)},
		{"Two-Factor Conformant", boolYesNo(org.TwoFactorConformant)},
		{"SAML Enabled", boolYesNo(org.SAMLEnabled)},
		{"Assessments Enforced", boolYesNo(org.AssessmentsEnforced)},
		{"Allow Force Delete WS", boolYesNo(org.AllowForceDeleteWorkspaces)},
	}
	if authPolicy := string(org.CollaboratorAuthPolicy); authPolicy != "" {
		settings.rows = append(settings.rows, wsDetailRow{"Auth Policy", authPolicy})
	}

	// ── Session ───────────────────────────────────────────────────────────────
	session := wsDetailSection{title: "Session"}
	session.rows = []wsDetailRow{
		{"Timeout", fmt.Sprintf("%d minutes", org.SessionTimeout)},
		{"Remember", fmt.Sprintf("%d minutes", org.SessionRemember)},
	}

	return []wsDetailSection{general, settings, session}
}

// renderOrgDetailContent renders the full detail view for the selected organization.
func (m Model) renderOrgDetailContent() string {
	h := m.contentHeight()
	if m.selectedOrg == nil {
		lines := make([]string, h)
		for i := range lines {
			lines[i] = contentStyle.Width(m.innerWidth()).Render("")
		}
		return strings.Join(lines, "\n")
	}

	sections := buildOrgDetailSections(m.selectedOrg)

	var all []string
	all = append(all, contentStyle.Width(m.innerWidth()).Render("")) // top padding

	for si, sec := range sections {
		all = append(all, m.renderDetailSectionHeader(sec.title))
		for _, row := range sec.rows {
			all = append(all, m.renderDetailKV(row.label, row.value))
		}
		if si < len(sections)-1 {
			all = append(all, contentStyle.Width(m.innerWidth()).Render(""))
		}
	}
	all = append(all, contentStyle.Width(m.innerWidth()).Render("")) // bottom padding

	// Clamp scroll and slice visible window.
	maxScroll := len(all) - h
	if maxScroll < 0 {
		maxScroll = 0
	}
	start := m.orgDetScroll
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
		out[i] = contentStyle.Width(m.innerWidth()).Render("")
	}
	return strings.Join(out, "\n")
}
