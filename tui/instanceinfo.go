// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"sort"
	"strings"

	"charm.land/lipgloss/v2"
	tea "charm.land/bubbletea/v2"
)

// renderInstanceInfoContent renders the instance info / health check view.
// Connection and version fields are populated immediately from ping response
// headers cached on client init; health check data arrives asynchronously.
func (m Model) renderInstanceInfoContent() string {
	h := m.contentHeight()

	if m.c == nil {
		out := make([]string, h)
		for i := range out {
			out[i] = contentStyle.Width(m.mainWidth()).Render("")
		}
		return strings.Join(out, "\n")
	}

	cl := m.c.Client

	// ── Build sections ───────────────────────────────────────────────────────

	var sections []wsDetailSection

	// Connection
	conn := wsDetailSection{title: "Connection"}
	conn.rows = []wsDetailRow{
		{"Hostname", m.hostname},
		{"Organization", m.org},
		{"Application", cl.AppName()},
	}
	sections = append(sections, conn)

	// Version
	ver := wsDetailSection{title: "Version"}
	ver.rows = []wsDetailRow{
		{"API Version", cl.RemoteAPIVersion()},
	}
	if v := cl.RemoteTFEVersion(); v != "" {
		ver.rows = append(ver.rows, wsDetailRow{"TFE Monthly", v})
	}
	if v := cl.RemoteTFENumericVersion(); v != "" {
		ver.rows = append(ver.rows, wsDetailRow{"TFE Numeric", v})
	}
	sections = append(sections, ver)

	// Health Check
	hc := wsDetailSection{title: "Health Check"}
	if m.healthCheckLoad {
		frame := spinnerFrames[m.spinnerIdx]
		hc.rows = []wsDetailRow{{"", frame + "  loading…"}}
	} else if m.healthCheckErr != "" {
		hc.rows = []wsDetailRow{{"Error", m.healthCheckErr}}
	} else if len(m.healthCheck) == 0 {
		hc.rows = []wsDetailRow{{"", "(no data)"}}
	} else {
		// Sort service names alphabetically for stable output.
		keys := make([]string, 0, len(m.healthCheck))
		for k := range m.healthCheck {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			hc.rows = append(hc.rows, wsDetailRow{k, m.healthCheck[k]})
		}
	}
	sections = append(sections, hc)

	// ── Assemble rendered rows ───────────────────────────────────────────────

	var all []string
	all = append(all, contentStyle.Width(m.mainWidth()).Render("")) // top padding
	for si, sec := range sections {
		all = append(all, m.renderDetailSectionHeader(sec.title))
		for _, row := range sec.rows {
			all = append(all, m.renderInstanceInfoKV(row.label, row.value))
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
	start := m.infoScroll
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

// renderInstanceInfoKV is like renderDetailKV but colors health check status
// values: UP/OK/HEALTHY → green, anything non-UP → red.
func (m Model) renderInstanceInfoKV(label, value string) string {
	maxValueWidth := m.mainWidth() - wsDetLabelWidth - 2
	if maxValueWidth < 10 {
		maxValueWidth = 10
	}
	labelPart := detailLabelStyle.Width(wsDetLabelWidth).Render("  " + label)

	// Choose value style based on status keyword.
	var valueStyle lipgloss.Style
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "UP", "OK", "HEALTHY", "TRUE":
		valueStyle = contentStyle.Foreground(colorSuccess)
	case "DOWN", "ERROR", "UNHEALTHY", "DEGRADED", "FAIL", "FAILED", "FALSE":
		valueStyle = contentStyle.Foreground(colorError)
	default:
		valueStyle = contentStyle
	}
	valuePart := valueStyle.Render(truncateStr(value, maxValueWidth))
	return m.padContent(labelPart+valuePart, contentStyle)
}

// handleInstanceInfoKey processes keys while the instance info view is active.
func (m Model) handleInstanceInfoKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.infoScroll > 0 {
			m.infoScroll--
		}
	case "down", "j":
		m.infoScroll++
	case "r":
		// Re-fetch health check.
		m.healthCheck = nil
		m.healthCheckLoad = true
		m.healthCheckErr = ""
		return m, tea.Batch(loadHealthCheck(m.c), tickSpinner())
	}
	return m, nil
}
