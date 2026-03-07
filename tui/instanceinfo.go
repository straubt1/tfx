// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"sort"
	"strings"

	"charm.land/lipgloss/v2"
	tea "charm.land/bubbletea/v2"
)

// renderInstanceInfoModal renders a centered modal popup showing connection,
// version, and health check information. It is designed to be composited on
// top of the current view via overlayInstanceInfoModal.
func (m Model) renderInstanceInfoModal() string {
	// Inner content width (before border + padding).
	innerW := 58
	if innerW > m.width-6 {
		innerW = m.width - 6
	}
	if innerW < 30 {
		innerW = 30
	}

	labelW := 16 // label column width inside the modal

	// Maximum scrollable content rows visible inside the modal.
	// Border (2) + padding rows (0) + hint row (1) are extra.
	maxRows := m.height - 10
	if maxRows < 5 {
		maxRows = 5
	}
	if maxRows > 18 {
		maxRows = 18
	}

	// ── Inline helpers ────────────────────────────────────────────────────────

	renderSec := func(title string) string {
		prefix := "  ── " + title + " "
		n := innerW - len([]rune(prefix))
		if n < 2 {
			n = 2
		}
		return contentTitleStyle.Render(prefix + strings.Repeat("─", n))
	}

	renderKV := func(label, value string) string {
		lp := detailLabelStyle.Width(labelW).Render("  " + label)
		var vs lipgloss.Style
		switch strings.ToUpper(strings.TrimSpace(value)) {
		case "UP", "OK", "HEALTHY", "TRUE":
			vs = contentStyle.Foreground(colorSuccess)
		case "DOWN", "ERROR", "UNHEALTHY", "DEGRADED", "FAIL", "FAILED", "FALSE":
			vs = contentStyle.Foreground(colorError)
		default:
			vs = contentStyle
		}
		maxVW := innerW - labelW
		if maxVW < 10 {
			maxVW = 10
		}
		return lp + vs.Render(truncateStr(value, maxVW))
	}

	// ── Build scrollable content lines ────────────────────────────────────────

	var all []string

	if m.c == nil {
		all = append(all, "  (no client)")
	} else {
		cl := m.c.Client

		all = append(all, renderSec("Instance"))
		all = append(all, renderKV("Application", cl.AppName()))
		all = append(all, renderKV("Hostname", m.hostname))
		all = append(all, renderKV("API Version", cl.RemoteAPIVersion()))
		if v := cl.RemoteTFEVersion(); v != "" {
			all = append(all, renderKV("TFE Monthly", v))
		}
		if v := cl.RemoteTFENumericVersion(); v != "" {
			all = append(all, renderKV("TFE Numeric", v))
		}
		all = append(all, "")

		all = append(all, renderSec("Health Check"))
		if m.healthCheckLoad {
			frame := spinnerFrames[m.spinnerIdx]
			all = append(all, renderKV("", frame+"  loading…"))
		} else if m.healthCheckErr != "" {
			all = append(all, renderKV("Error", m.healthCheckErr))
		} else if len(m.healthCheck) == 0 {
			all = append(all, renderKV("", "(no data)"))
		} else {
			keys := make([]string, 0, len(m.healthCheck))
			for k := range m.healthCheck {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				all = append(all, renderKV(k, m.healthCheck[k]))
			}
		}
	}

	// ── Scroll ────────────────────────────────────────────────────────────────

	maxScroll := len(all) - maxRows
	if maxScroll < 0 {
		maxScroll = 0
	}
	start := m.infoScroll
	if start > maxScroll {
		start = maxScroll
	}
	visible := all[start:]
	if len(visible) > maxRows {
		visible = visible[:maxRows]
	}

	// Pad remaining rows to keep modal height stable while scrolling.
	padded := make([]string, maxRows)
	copy(padded, visible)
	for i := len(visible); i < maxRows; i++ {
		padded[i] = ""
	}

	// Footer hint — always pinned at the bottom of the modal.
	padded = append(padded, contentPlaceholderStyle.Render(
		"  [↑↓] scroll  •  [r] refresh  •  [i/esc] close",
	))

	inner := strings.Join(padded, "\n")

	// Wrap in a rounded border with accent colour.
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorAccent).
		Background(colorBg).
		Foreground(colorFg).
		Padding(0, 1).
		Width(innerW).
		Render(inner)
}

// overlayInstanceInfoModal composites the instance info modal on top of the
// already-rendered full-screen base content using lipgloss Canvas + Layer.
func (m Model) overlayInstanceInfoModal(base string) string {
	modal := m.renderInstanceInfoModal()
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

// handleInstanceInfoModalKey processes keys while the instance info modal is visible.
func (m Model) handleInstanceInfoModalKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "i", "esc":
		m.showInstanceInfo = false
	case "up", "k":
		if m.infoScroll > 0 {
			m.infoScroll--
		}
	case "down", "j":
		m.infoScroll++
	case "r":
		// Re-fetch health check data.
		m.healthCheck = nil
		m.healthCheckLoad = true
		m.healthCheckErr = ""
		return m, tea.Batch(loadHealthCheck(m.c), tickSpinner())
	}
	return m, nil
}
