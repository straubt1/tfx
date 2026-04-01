// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"strings"

	"charm.land/lipgloss/v2"
	tea "charm.land/bubbletea/v2"
	tfe "github.com/hashicorp/go-tfe"
)

// renderCVArchivedModal renders a small modal explaining that the selected
// configuration version has been archived and is no longer available.
func (m Model) renderCVArchivedModal() string {
	cv := m.cvArchivedCV

	innerW := 52
	if innerW > m.width-6 {
		innerW = m.width - 6
	}
	if innerW < 30 {
		innerW = 30
	}

	labelW := 10

	dimStyle := lipgloss.NewStyle().Foreground(colorDim)
	accentStyle := lipgloss.NewStyle().Foreground(colorAccent)
	errorStyle := lipgloss.NewStyle().Foreground(colorError)

	renderKV := func(label, value string) string {
		l := dimStyle.Render(label + ":")
		pad := labelW - len([]rune(label+":"))
		if pad < 1 {
			pad = 1
		}
		return "  " + l + strings.Repeat(" ", pad) + accentStyle.Render(value)
	}

	// ── Body rows ─────────────────────────────────────────────────────────────
	var rows []string

	rows = append(rows, "")
	rows = append(rows, "  "+errorStyle.Render("✗  This configuration version is archived"))
	rows = append(rows, "     and its contents are no longer available.")
	rows = append(rows, "")

	rows = append(rows, renderKV("ID", cv.ID))

	// Archived-at timestamp
	if cv.StatusTimestamps != nil && !cv.StatusTimestamps.ArchivedAt.IsZero() {
		rows = append(rows, renderKV("Archived", timestampWithRelative(cv.StatusTimestamps.ArchivedAt)))
	}

	// Source / means
	if cv.Source != "" {
		rows = append(rows, renderKV("Source", cvSourceLabel(cv.Source)))
	}

	rows = append(rows, "")
	rows = append(rows, "  "+dimStyle.Render("Press Esc or Enter to dismiss."))
	rows = append(rows, "")

	body := strings.Join(rows, "\n")

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorError).
		Width(innerW).
		Padding(0, 1)

	return borderStyle.Render(body)
}

// cvSourceLabel returns a human-readable label for a ConfigurationSource.
func cvSourceLabel(src tfe.ConfigurationSource) string {
	switch src {
	case tfe.ConfigurationSourceGithub:
		return "GitHub"
	case tfe.ConfigurationSourceGitlab:
		return "GitLab"
	case tfe.ConfigurationSourceBitbucket:
		return "Bitbucket"
	case tfe.ConfigurationSourceAdo:
		return "Azure DevOps"
	case tfe.ConfigurationSourceTerraform:
		return "Terraform CLI"
	default:
		if string(src) == "" {
			return "unknown"
		}
		return string(src)
	}
}

// overlayCVArchivedModal composites the CV archived modal centered over the
// already-rendered full-screen base content.
func (m Model) overlayCVArchivedModal(base string) string {
	modal := m.renderCVArchivedModal()
	mw := lipgloss.Width(modal)
	mh := lipgloss.Height(modal)

	x := (m.width - mw) / 2
	y := (m.height - mh) / 2
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	baseLayer := lipgloss.NewLayer(base).Z(0)
	modalLayer := lipgloss.NewLayer(modal).X(x).Y(y).Z(1)

	compositor := lipgloss.NewCompositor(baseLayer, modalLayer)
	canvas := lipgloss.NewCanvas(m.width, m.height)
	canvas.Compose(compositor)
	return canvas.Render()
}

// handleCVArchivedModalKey processes keys while the CV archived modal is open.
func (m Model) handleCVArchivedModalKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "enter", "q":
		m.showCVArchivedModal = false
		m.cvArchivedCV = nil
	}
	return m, nil
}
