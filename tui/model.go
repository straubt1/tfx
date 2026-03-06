// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/version"
)

const (
	fixedLines = 4 // header + breadcrumb + statusbar + clihint
	minWidth   = 60
	minHeight  = 10
)

type viewType int

const (
	viewProjects   viewType = iota
	viewWorkspaces
)

// Model is the root TUI model. All state lives here per the ELM architecture.
type Model struct {
	// Layout
	width  int
	height int
	ready  bool

	// Connection
	c        *client.TfxClient
	hostname string
	org      string

	// View routing
	currentView viewType
	showHelp    bool

	// Loading / error state
	loading bool
	errMsg  string

	// Project list state
	projects      []*tfe.Project
	projCursor    int
	projOffset    int
	projFilter    string
	projFiltering bool

	// Workspace list state
	workspaces   []*tfe.Workspace
	wsCursor     int
	wsOffset     int
	wsFilter     string
	wsFiltering  bool
	selectedProj *tfe.Project
}

func newModel(c *client.TfxClient) Model {
	return Model{
		c:        c,
		hostname: c.Hostname,
		org:      c.OrganizationName,
		loading:  true,
	}
}

func (m Model) Init() tea.Cmd {
	return loadProjects(m.c, m.org)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

	case projectsLoadedMsg:
		m.projects = []*tfe.Project(msg)
		m.loading = false
		m.currentView = viewProjects
		m.errMsg = ""

	case workspacesLoadedMsg:
		m.workspaces = []*tfe.Workspace(msg)
		m.loading = false
		m.currentView = viewWorkspaces
		m.wsOffset = 0
		m.errMsg = ""

	case fetchErrMsg:
		m.loading = false
		m.errMsg = msg.err.Error()

	case tea.KeyPressMsg:
		// Help overlay consumes all keys.
		if m.showHelp {
			if msg.String() == "esc" || msg.String() == "?" {
				m.showHelp = false
			}
			return m, nil
		}

		// Filter input mode.
		if m.isFiltering() {
			return m.handleFilterKey(msg)
		}

		// Global keys.
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "?":
			m.showHelp = true
			return m, nil
		case "esc":
			return m.navigateBack()
		case "r":
			return m.refresh()
		}

		// View-specific navigation (only when not loading).
		if !m.loading {
			switch m.currentView {
			case viewProjects:
				return m.handleProjectsKey(msg)
			case viewWorkspaces:
				return m.handleWorkspacesKey(msg)
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

// ── Navigation ────────────────────────────────────────────────────────────────

func (m Model) navigateBack() (tea.Model, tea.Cmd) {
	switch m.currentView {
	case viewWorkspaces:
		m.currentView = viewProjects
		m.workspaces = nil
		m.wsCursor = 0
		m.wsOffset = 0
		m.wsFilter = ""
		m.wsFiltering = false
		m.selectedProj = nil
	}
	return m, nil
}

func (m Model) refresh() (tea.Model, tea.Cmd) {
	m.loading = true
	m.errMsg = ""
	switch m.currentView {
	case viewProjects:
		m.projects = nil
		m.projCursor = 0
		m.projOffset = 0
		return m, loadProjects(m.c, m.org)
	case viewWorkspaces:
		m.workspaces = nil
		m.wsCursor = 0
		m.wsOffset = 0
		projectID := ""
		if m.selectedProj != nil {
			projectID = m.selectedProj.ID
		}
		return m, loadWorkspaces(m.c, m.org, projectID)
	}
	return m, nil
}

// ── Key handlers ──────────────────────────────────────────────────────────────

func (m Model) isFiltering() bool {
	return (m.currentView == viewProjects && m.projFiltering) ||
		(m.currentView == viewWorkspaces && m.wsFiltering)
}

func (m Model) handleFilterKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		if m.currentView == viewProjects {
			m.projFilter = ""
			m.projFiltering = false
			m.projCursor = 0
			m.projOffset = 0
		} else {
			m.wsFilter = ""
			m.wsFiltering = false
			m.wsCursor = 0
			m.wsOffset = 0
		}
	case "enter":
		if m.currentView == viewProjects {
			m.projFiltering = false
		} else {
			m.wsFiltering = false
		}
	case "backspace":
		if m.currentView == viewProjects {
			r := []rune(m.projFilter)
			if len(r) > 0 {
				m.projFilter = string(r[:len(r)-1])
				m.projCursor = 0
				m.projOffset = 0
			}
		} else {
			r := []rune(m.wsFilter)
			if len(r) > 0 {
				m.wsFilter = string(r[:len(r)-1])
				m.wsCursor = 0
				m.wsOffset = 0
			}
		}
	default:
		r := []rune(msg.String())
		if len(r) == 1 && r[0] >= 32 {
			if m.currentView == viewProjects {
				m.projFilter += string(r)
				m.projCursor = 0
				m.projOffset = 0
			} else {
				m.wsFilter += string(r)
				m.wsCursor = 0
				m.wsOffset = 0
			}
		}
	}
	return m, nil
}

func (m Model) handleProjectsKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	filtered := filteredProjects(m)
	n := len(filtered)
	vis := m.projVisibleRows()

	switch msg.String() {
	case "up", "k":
		if m.projCursor > 0 {
			m.projCursor--
			if m.projCursor < m.projOffset {
				m.projOffset = m.projCursor
			}
		}
	case "down", "j":
		if m.projCursor < n-1 {
			m.projCursor++
			if m.projCursor >= m.projOffset+vis {
				m.projOffset = m.projCursor - vis + 1
			}
		}
	case "g":
		m.projCursor = 0
		m.projOffset = 0
	case "G":
		if n > 0 {
			m.projCursor = n - 1
			if n > vis {
				m.projOffset = n - vis
			}
		}
	case "/":
		m.projFiltering = true
	case "enter":
		if n == 0 || m.projCursor >= n {
			break
		}
		sel := filtered[m.projCursor]
		m.selectedProj = sel
		m.loading = true
		m.errMsg = ""
		m.wsCursor = 0
		m.wsOffset = 0
		m.wsFilter = ""
		return m, loadWorkspaces(m.c, m.org, sel.ID)
	}
	return m, nil
}

func (m Model) handleWorkspacesKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	filtered := filteredWorkspaces(m)
	n := len(filtered)
	vis := m.wsVisibleRows()

	switch msg.String() {
	case "up", "k":
		if m.wsCursor > 0 {
			m.wsCursor--
			if m.wsCursor < m.wsOffset {
				m.wsOffset = m.wsCursor
			}
		}
	case "down", "j":
		if m.wsCursor < n-1 {
			m.wsCursor++
			if m.wsCursor >= m.wsOffset+vis {
				m.wsOffset = m.wsCursor - vis + 1
			}
		}
	case "g":
		m.wsCursor = 0
		m.wsOffset = 0
	case "G":
		if n > 0 {
			m.wsCursor = n - 1
			if n > vis {
				m.wsOffset = n - vis
			}
		}
	case "/":
		m.wsFiltering = true
	case "enter":
		// Phase 3: drill into runs.
	}
	return m, nil
}

// ── Layout helpers ────────────────────────────────────────────────────────────

func (m Model) contentHeight() int {
	h := m.height - fixedLines
	if h < 1 {
		return 1
	}
	return h
}

func (m Model) projVisibleRows() int {
	h := m.contentHeight() - 2 // header + divider
	if m.projFilter != "" || m.projFiltering {
		h--
	}
	if h < 1 {
		return 1
	}
	return h
}

func (m Model) wsVisibleRows() int {
	h := m.contentHeight() - 2
	if m.wsFilter != "" || m.wsFiltering {
		h--
	}
	if h < 1 {
		return 1
	}
	return h
}

// pad fills a rendered string to the full terminal width using the given style.
func (m Model) pad(rendered string, style lipgloss.Style) string {
	w := m.width - lipgloss.Width(rendered)
	if w < 0 {
		w = 0
	}
	return rendered + style.Width(w).Render("")
}

// ── Content routing ───────────────────────────────────────────────────────────

func (m Model) renderContent() string {
	if m.loading {
		return m.renderLoadingContent()
	}
	if m.errMsg != "" {
		return m.renderErrorContent()
	}
	switch m.currentView {
	case viewProjects:
		return m.renderProjectsContent()
	case viewWorkspaces:
		return m.renderWorkspacesContent()
	}
	return m.renderLoadingContent()
}

func (m Model) renderLoadingContent() string {
	h := m.contentHeight()
	lines := make([]string, h)
	mid := h / 2
	for i := range lines {
		if i == mid {
			lines[i] = contentPlaceholderStyle.Width(m.width).Render("  Loading…")
		} else {
			lines[i] = contentStyle.Width(m.width).Render("")
		}
	}
	return strings.Join(lines, "\n")
}

func (m Model) renderErrorContent() string {
	h := m.contentHeight()
	lines := make([]string, h)
	mid := h / 2
	for i := range lines {
		if i == mid {
			lines[i] = statusErrorStyle.Width(m.width).Render(fmt.Sprintf("  ✗  %s", m.errMsg))
		} else {
			lines[i] = contentStyle.Width(m.width).Render("")
		}
	}
	return strings.Join(lines, "\n")
}

// ── Fixed chrome ──────────────────────────────────────────────────────────────

func (m Model) renderHeader() string {
	app := headerAppStyle.Render(" TFx ")
	info := headerInfoStyle.Render(fmt.Sprintf(" %s  ⬥  %s ", m.hostname, m.org))
	ver := headerVersionStyle.Render(fmt.Sprintf(" v%s ", version.Version))

	used := lipgloss.Width(app) + lipgloss.Width(info) + lipgloss.Width(ver)
	gap := m.width - used
	if gap < 0 {
		gap = 0
	}
	return app + info + headerStyle.Width(gap).Render("") + ver
}

func (m Model) renderBreadcrumb() string {
	sep := breadcrumbSepStyle.Render("  /  ")
	orgPart := breadcrumbBarStyle.Render(fmt.Sprintf(" org: %s", m.org))

	var line string
	switch m.currentView {
	case viewProjects:
		line = orgPart + sep + breadcrumbActiveStyle.Render("projects ")
	case viewWorkspaces:
		projName := ""
		if m.selectedProj != nil {
			projName = m.selectedProj.Name
		}
		line = orgPart + sep +
			breadcrumbBarStyle.Render(fmt.Sprintf("project: %s", projName)) +
			sep + breadcrumbActiveStyle.Render("workspaces ")
	default:
		line = orgPart
	}
	return m.pad(line, breadcrumbBarStyle)
}

func (m Model) renderStatusBar() string {
	if m.loading {
		return m.pad(statusLoadingStyle.Render("  Loading…"), statusLoadingStyle)
	}
	if m.errMsg != "" {
		return m.pad(statusErrorStyle.Render(fmt.Sprintf("  ✗  %s", m.errMsg)), statusErrorStyle)
	}

	var msg string
	switch m.currentView {
	case viewProjects:
		fp := filteredProjects(m)
		if m.projFilter != "" {
			msg = fmt.Sprintf("  %d / %d projects  •  filter: %s", len(fp), len(m.projects), m.projFilter)
		} else {
			msg = fmt.Sprintf("  %d projects", len(m.projects))
		}
	case viewWorkspaces:
		fw := filteredWorkspaces(m)
		if m.wsFilter != "" {
			msg = fmt.Sprintf("  %d / %d workspaces  •  filter: %s", len(fw), len(m.workspaces), m.wsFilter)
		} else {
			msg = fmt.Sprintf("  %d workspaces", len(m.workspaces))
		}
	default:
		msg = "  Ready"
	}
	return m.pad(statusBarStyle.Render(msg), statusBarStyle)
}

func (m Model) renderCliHint() string {
	var cliCmd string
	switch m.currentView {
	case viewProjects:
		cliCmd = "tfx project list"
	case viewWorkspaces:
		if m.selectedProj != nil {
			cliCmd = fmt.Sprintf("tfx workspace list --project-id %s", m.selectedProj.ID)
		} else {
			cliCmd = "tfx workspace list"
		}
	default:
		cliCmd = "tfx"
	}

	label := cliHintBarStyle.Render("  cmd: ")
	cmd := cliHintCmdStyle.Render(cliCmd)
	hints := cliHintBarStyle.Render("   •   ? help   •   q quit")
	return m.pad(label+cmd+hints, cliHintBarStyle)
}

// ── Help overlay ──────────────────────────────────────────────────────────────

func (m Model) renderHelpOverlay() string {
	type binding struct {
		key  string
		desc string
	}
	bindings := []binding{
		{"↑ / k", "move up"},
		{"↓ / j", "move down"},
		{"enter", "select / drill in"},
		{"esc", "go back / clear filter"},
		{"r", "refresh"},
		{"/", "filter"},
		{"g / G", "jump to top / bottom"},
		{"?", "toggle help"},
		{"q", "quit"},
	}

	lines := make([]string, 0, m.height)
	lines = append(lines, m.renderHeader())
	lines = append(lines, m.pad(helpTitleStyle.Render("  Keyboard Shortcuts"), helpTitleStyle))
	lines = append(lines, helpBarStyle.Width(m.width).Render(""))

	for _, b := range bindings {
		key := helpKeyStyle.Width(14).Render(b.key)
		desc := helpDescStyle.Render("  " + b.desc)
		line := helpBarStyle.Render("  ") + key + desc
		lines = append(lines, m.pad(line, helpBarStyle))
	}

	lines = append(lines, helpBarStyle.Width(m.width).Render(""))
	lines = append(lines, m.pad(helpBarStyle.Render("  Press ? or esc to close"), helpBarStyle))

	for len(lines) < m.height {
		lines = append(lines, helpBarStyle.Width(m.width).Render(""))
	}
	return strings.Join(lines[:m.height], "\n")
}

// ── Table rendering (shared by projects + workspaces) ────────────────────────

// column defines a table column's display name and character width.
type column struct {
	name  string
	width int
}

func (m Model) renderTableHeader(cols []column) string {
	parts := []string{tableHeaderStyle.Render("  ")} // cursor placeholder
	for _, col := range cols {
		parts = append(parts, tableHeaderStyle.Width(col.width).Render(col.name))
		parts = append(parts, tableHeaderStyle.Render("  "))
	}
	return m.pad(strings.Join(parts, ""), tableHeaderStyle)
}

func (m Model) renderTableDivider() string {
	return contentDividerStyle.Width(m.width).Render(strings.Repeat("─", m.width))
}

func (m Model) renderTableRow(selected bool, cells []string, cols []column) string {
	style := tableRowStyle
	cursor := "  "
	if selected {
		style = tableRowSelectedStyle
		cursor = "> "
	}

	parts := []string{style.Render(cursor)}
	for i, col := range cols {
		val := ""
		if i < len(cells) {
			val = truncateStr(cells[i], col.width)
		}
		parts = append(parts, style.Width(col.width).Render(val))
		parts = append(parts, style.Render("  "))
	}

	return m.pad(strings.Join(parts, ""), style)
}

func (m Model) renderFilterBar(filter string, active bool) string {
	prompt := filterBarStyle.Render("  / ")
	var text string
	if active {
		text = filterBarActiveStyle.Render(filter + "▌")
	} else {
		text = filterBarActiveStyle.Render(filter)
	}
	return m.pad(prompt+text, filterBarStyle)
}

// truncateStr truncates s to at most n runes, appending "…" if shortened.
func truncateStr(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	if n > 1 {
		return string(r[:n-1]) + "…"
	}
	return string(r[:n])
}
