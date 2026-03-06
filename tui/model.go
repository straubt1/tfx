// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"fmt"
	"image/color"
	"os/exec"
	"runtime"
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
	viewRuns
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

	// Loading / error / transient state
	loading      bool
	errMsg       string
	clipFeedback string // cleared on next keypress

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

	// Run list state
	runs         []*tfe.Run
	runCursor    int
	runOffset    int
	runFilter    string
	runFiltering bool
	selectedWS   *tfe.Workspace
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

	case runsLoadedMsg:
		m.runs = []*tfe.Run(msg)
		m.loading = false
		m.currentView = viewRuns
		m.runOffset = 0
		m.errMsg = ""

	case fetchErrMsg:
		m.loading = false
		m.errMsg = msg.err.Error()

	case tea.KeyPressMsg:
		// Clear transient clipboard feedback on next key.
		m.clipFeedback = ""

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
		case "c":
			cmd := m.currentCliCmd()
			if err := copyToClipboard(cmd); err == nil {
				m.clipFeedback = "✓ copied to clipboard"
			} else {
				m.clipFeedback = "clipboard unavailable"
			}
			return m, nil
		}

		// View-specific navigation (only when not loading).
		if !m.loading {
			switch m.currentView {
			case viewProjects:
				return m.handleProjectsKey(msg)
			case viewWorkspaces:
				return m.handleWorkspacesKey(msg)
			case viewRuns:
				return m.handleRunsKey(msg)
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
	case viewRuns:
		m.currentView = viewWorkspaces
		m.runs = nil
		m.runCursor = 0
		m.runOffset = 0
		m.runFilter = ""
		m.runFiltering = false
		m.selectedWS = nil
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
	case viewRuns:
		m.runs = nil
		m.runCursor = 0
		m.runOffset = 0
		if m.selectedWS != nil {
			return m, loadRuns(m.c, m.selectedWS.ID)
		}
	}
	return m, nil
}

// ── Key handlers ──────────────────────────────────────────────────────────────

func (m Model) isFiltering() bool {
	return (m.currentView == viewProjects && m.projFiltering) ||
		(m.currentView == viewWorkspaces && m.wsFiltering) ||
		(m.currentView == viewRuns && m.runFiltering)
}

func (m Model) handleFilterKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		switch m.currentView {
		case viewProjects:
			m.projFilter, m.projFiltering = "", false
			m.projCursor, m.projOffset = 0, 0
		case viewWorkspaces:
			m.wsFilter, m.wsFiltering = "", false
			m.wsCursor, m.wsOffset = 0, 0
		case viewRuns:
			m.runFilter, m.runFiltering = "", false
			m.runCursor, m.runOffset = 0, 0
		}
	case "enter":
		switch m.currentView {
		case viewProjects:
			m.projFiltering = false
		case viewWorkspaces:
			m.wsFiltering = false
		case viewRuns:
			m.runFiltering = false
		}
	case "backspace":
		switch m.currentView {
		case viewProjects:
			if r := []rune(m.projFilter); len(r) > 0 {
				m.projFilter = string(r[:len(r)-1])
				m.projCursor, m.projOffset = 0, 0
			}
		case viewWorkspaces:
			if r := []rune(m.wsFilter); len(r) > 0 {
				m.wsFilter = string(r[:len(r)-1])
				m.wsCursor, m.wsOffset = 0, 0
			}
		case viewRuns:
			if r := []rune(m.runFilter); len(r) > 0 {
				m.runFilter = string(r[:len(r)-1])
				m.runCursor, m.runOffset = 0, 0
			}
		}
	default:
		if r := []rune(msg.String()); len(r) == 1 && r[0] >= 32 {
			switch m.currentView {
			case viewProjects:
				m.projFilter += string(r)
				m.projCursor, m.projOffset = 0, 0
			case viewWorkspaces:
				m.wsFilter += string(r)
				m.wsCursor, m.wsOffset = 0, 0
			case viewRuns:
				m.runFilter += string(r)
				m.runCursor, m.runOffset = 0, 0
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
		m.projCursor, m.projOffset = 0, 0
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
		m.wsCursor, m.wsOffset, m.wsFilter = 0, 0, ""
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
		m.wsCursor, m.wsOffset = 0, 0
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
		if n == 0 || m.wsCursor >= n {
			break
		}
		sel := filtered[m.wsCursor]
		m.selectedWS = sel
		m.loading = true
		m.errMsg = ""
		m.runCursor, m.runOffset, m.runFilter = 0, 0, ""
		return m, loadRuns(m.c, sel.ID)
	}
	return m, nil
}

func (m Model) handleRunsKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	filtered := filteredRuns(m)
	n := len(filtered)
	vis := m.runVisibleRows()

	switch msg.String() {
	case "up", "k":
		if m.runCursor > 0 {
			m.runCursor--
			if m.runCursor < m.runOffset {
				m.runOffset = m.runCursor
			}
		}
	case "down", "j":
		if m.runCursor < n-1 {
			m.runCursor++
			if m.runCursor >= m.runOffset+vis {
				m.runOffset = m.runCursor - vis + 1
			}
		}
	case "g":
		m.runCursor, m.runOffset = 0, 0
	case "G":
		if n > 0 {
			m.runCursor = n - 1
			if n > vis {
				m.runOffset = n - vis
			}
		}
	case "/":
		m.runFiltering = true
	case "enter":
		// Phase 5: drill into run detail.
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
	h := m.contentHeight() - 2 // table header + divider
	if m.projFilter != "" || m.projFiltering {
		h-- // filter bar
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

func (m Model) runVisibleRows() int {
	h := m.contentHeight() - 2
	if m.runFilter != "" || m.runFiltering {
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

// currentCliCmd returns the equivalent tfx CLI command for the active view.
func (m Model) currentCliCmd() string {
	switch m.currentView {
	case viewProjects:
		return "tfx project list"
	case viewWorkspaces:
		if m.selectedProj != nil {
			return fmt.Sprintf("tfx workspace list --project-id %s", m.selectedProj.ID)
		}
		return "tfx workspace list"
	case viewRuns:
		if m.selectedWS != nil {
			return fmt.Sprintf("tfx workspace run list -n %s", m.selectedWS.Name)
		}
		return "tfx workspace run list"
	default:
		return "tfx"
	}
}

// copyToClipboard writes text to the system clipboard via a platform-native command.
func copyToClipboard(text string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		cmd = exec.Command("xclip", "-selection", "clipboard")
	case "windows":
		cmd = exec.Command("clip")
	default:
		return fmt.Errorf("clipboard not supported on %s", runtime.GOOS)
	}
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
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
	case viewRuns:
		return m.renderRunsContent()
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
	case viewRuns:
		projName := ""
		if m.selectedProj != nil {
			projName = m.selectedProj.Name
		}
		wsName := ""
		if m.selectedWS != nil {
			wsName = m.selectedWS.Name
		}
		line = orgPart + sep +
			breadcrumbBarStyle.Render(fmt.Sprintf("project: %s", projName)) +
			sep +
			breadcrumbBarStyle.Render(fmt.Sprintf("workspace: %s", wsName)) +
			sep + breadcrumbActiveStyle.Render("runs ")
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
	if m.clipFeedback != "" {
		return m.pad(statusSuccessStyle.Render("  "+m.clipFeedback), statusSuccessStyle)
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
	case viewRuns:
		fr := filteredRuns(m)
		if m.runFilter != "" {
			msg = fmt.Sprintf("  %d / %d runs  •  filter: %s", len(fr), len(m.runs), m.runFilter)
		} else {
			msg = fmt.Sprintf("  %d runs", len(m.runs))
		}
	default:
		msg = "  Ready"
	}
	return m.pad(statusBarStyle.Render(msg), statusBarStyle)
}

func (m Model) renderCliHint() string {
	label := cliHintBarStyle.Render("  cmd: ")
	cmd := cliHintCmdStyle.Render(m.currentCliCmd())
	hints := cliHintBarStyle.Render("   •   c copy   •   ? help   •   q quit")
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
		{"c", "copy CLI command"},
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

// ── Table rendering (shared by all list views) ────────────────────────────────

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

// renderTableRowWithCellStyles renders a row where individual cells can have
// a custom foreground color. cellFgs[i] overrides the fg for column i when
// the row is not selected (selection style always takes precedence).
func (m Model) renderTableRowWithCellStyles(selected bool, cells []string, cols []column, cellFgs []color.Color) string {
	base := tableRowStyle
	cursor := "  "
	if selected {
		base = tableRowSelectedStyle
		cursor = "> "
	}

	parts := []string{base.Render(cursor)}
	for i, col := range cols {
		val := ""
		if i < len(cells) {
			val = truncateStr(cells[i], col.width)
		}
		cellStyle := base.Width(col.width)
		if !selected && i < len(cellFgs) && cellFgs[i] != nil {
			cellStyle = cellStyle.Foreground(cellFgs[i])
		}
		parts = append(parts, cellStyle.Render(val))
		parts = append(parts, base.Render("  "))
	}
	return m.pad(strings.Join(parts, ""), base)
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
