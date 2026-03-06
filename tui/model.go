// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/straubt1/tfx/version"
)

const (
	fixedLines = 4 // header + breadcrumb + statusbar + clihint
	minWidth   = 60
	minHeight  = 10
)

// Model is the root TUI model. All state lives here.
type Model struct {
	width    int
	height   int
	ready    bool
	showHelp bool
	hostname string
	org      string
}

func newModel(hostname, org string) Model {
	return Model{
		hostname: hostname,
		org:      org,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "?":
			m.showHelp = !m.showHelp
		case "esc":
			if m.showHelp {
				m.showHelp = false
			}
		}
	}
	return m, nil
}

func (m Model) View() tea.View {
	var content string

	if !m.ready {
		content = "\n  Initializing..."
	} else if m.width < minWidth || m.height < minHeight {
		content = fmt.Sprintf("\n  Terminal too small (%dx%d). Minimum: %dx%d.", m.width, m.height, minWidth, minHeight)
	} else if m.showHelp {
		content = m.renderHelpOverlay()
	} else {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			m.renderHeader(),
			m.renderBreadcrumb(),
			m.renderContent(),
			m.renderStatusBar(),
			m.renderCliHint(),
		)
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

// contentHeight returns the number of lines available for the main content area.
func (m Model) contentHeight() int {
	h := m.height - fixedLines
	if h < 1 {
		return 1
	}
	return h
}

// pad fills a line to the full terminal width using the given style.
func (m Model) pad(rendered string, style lipgloss.Style) string {
	w := m.width - lipgloss.Width(rendered)
	if w < 0 {
		w = 0
	}
	return rendered + style.Width(w).Render("")
}

func (m Model) renderHeader() string {
	app := headerAppStyle.Render(" TFx ")
	info := headerInfoStyle.Render(fmt.Sprintf(" %s  ⬥  %s ", m.hostname, m.org))
	ver := headerVersionStyle.Render(fmt.Sprintf(" v%s ", version.Version))

	// Fill the gap between info and version
	used := lipgloss.Width(app) + lipgloss.Width(info) + lipgloss.Width(ver)
	gap := m.width - used
	if gap < 0 {
		gap = 0
	}

	return app + info + headerStyle.Width(gap).Render("") + ver
}

func (m Model) renderBreadcrumb() string {
	org := breadcrumbActiveStyle.Render(fmt.Sprintf(" org: %s", m.org))
	sep := breadcrumbSepStyle.Render("  /  ")
	cur := breadcrumbBarStyle.Render("projects ")

	line := org + sep + cur
	return m.pad(line, breadcrumbBarStyle)
}

func (m Model) renderContent() string {
	h := m.contentHeight()
	lines := make([]string, h)

	title := contentTitleStyle.Width(m.width).Render("  Projects")
	divider := contentDividerStyle.Width(m.width).Render(strings.Repeat("─", m.width))
	placeholder := contentPlaceholderStyle.Width(m.width).Render("  TUI Phase 1 — data views coming in Phase 2. Press ? for keyboard shortcuts.")
	empty := contentStyle.Width(m.width).Render("")

	if h > 0 {
		lines[0] = title
	}
	if h > 1 {
		lines[1] = divider
	}
	if h > 2 {
		lines[2] = placeholder
	}
	for i := 3; i < h; i++ {
		lines[i] = empty
	}

	return strings.Join(lines, "\n")
}

func (m Model) renderStatusBar() string {
	status := statusBarStyle.Render("  Ready")
	return m.pad(status, statusBarStyle)
}

func (m Model) renderCliHint() string {
	label := cliHintBarStyle.Render("  cmd: ")
	cmd := cliHintCmdStyle.Render("tfx project list")
	hints := cliHintBarStyle.Render("   •   ? help   •   q quit")

	line := label + cmd + hints
	return m.pad(line, cliHintBarStyle)
}

func (m Model) renderHelpOverlay() string {
	type binding struct {
		key  string
		desc string
	}

	bindings := []binding{
		{"↑ / k", "move up"},
		{"↓ / j", "move down"},
		{"enter", "select / drill in"},
		{"esc", "go back"},
		{"r", "refresh"},
		{"/", "filter"},
		{"?", "toggle help"},
		{"q", "quit"},
	}

	lines := make([]string, 0, m.height)

	// Header (reuse the app header for context)
	lines = append(lines, m.renderHeader())

	// Title
	title := m.pad(helpTitleStyle.Render("  Keyboard Shortcuts"), helpTitleStyle)
	lines = append(lines, title)
	lines = append(lines, helpBarStyle.Width(m.width).Render(""))

	// Bindings
	for _, b := range bindings {
		key := helpKeyStyle.Width(12).Render(b.key)
		desc := helpDescStyle.Render("  " + b.desc)
		line := helpBarStyle.Render("  ") + key + desc
		lines = append(lines, m.pad(line, helpBarStyle))
	}

	lines = append(lines, helpBarStyle.Width(m.width).Render(""))
	close := m.pad(helpBarStyle.Render("  Press ? or esc to close"), helpBarStyle)
	lines = append(lines, close)

	// Pad remaining lines to fill the screen
	for len(lines) < m.height {
		lines = append(lines, helpBarStyle.Width(m.width).Render(""))
	}

	return strings.Join(lines[:m.height], "\n")
}
