// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/straubt1/tfx/client"
)

// ── Filtering ─────────────────────────────────────────────────────────────────

// filteredDebugEvents returns the apiEvents slice with m.debugFilter applied.
// Returns newest-first (the natural order of m.apiEvents).
func (m Model) filteredDebugEvents() []client.APIEvent {
	if m.debugFilter == "" {
		return m.apiEvents
	}
	f := strings.ToLower(m.debugFilter)
	var out []client.APIEvent
	for _, e := range m.apiEvents {
		haystack := strings.ToLower(e.Method + " " + e.Path)
		if strings.Contains(haystack, f) {
			out = append(out, e)
		}
	}
	return out
}

// ── Key handler ───────────────────────────────────────────────────────────────

func (m Model) handleDebugPanelKey(msg tea.KeyPressMsg) (Model, tea.Cmd) {
	// Always handle Tab — toggles panel focus back to the left (main) view.
	if msg.String() == "tab" {
		m.debugFocused = false
		return m, nil
	}

	// ── Detail view keys ─────────────────────────────────────────────────
	if m.debugDetailMode {
		bodyH := m.contentHeight() - 2 // 2 = title + divider
		halfStep := bodyH / 2
		if halfStep < 1 {
			halfStep = 1
		}
		switch msg.String() {
		case "up", "k":
			m.debugBodyScroll--
			if m.debugBodyScroll < 0 {
				m.debugBodyScroll = 0
			}
		case "down", "j":
			m.debugBodyScroll++
		case "shift+up":
			m.debugBodyScroll -= bodyH
			if m.debugBodyScroll < 0 {
				m.debugBodyScroll = 0
			}
		case "shift+down":
			m.debugBodyScroll += bodyH
		case "ctrl+u":
			m.debugBodyScroll -= halfStep
			if m.debugBodyScroll < 0 {
				m.debugBodyScroll = 0
			}
		case "ctrl+d":
			m.debugBodyScroll += halfStep
		case "g":
			m.debugBodyScroll = 0
		case "G":
			m.debugBodyScroll = 9999 // clamped in renderer
		case "esc":
			m.debugDetailMode = false
			m.debugBodyScroll = 0
		}
		return m, nil
	}

	// ── List view keys ───────────────────────────────────────────────────
	events := m.filteredDebugEvents()
	n := len(events)

	switch msg.String() {
	case "up", "k":
		if m.debugCursor > 0 {
			m.debugCursor--
		}
	case "down", "j":
		if m.debugCursor < n-1 {
			m.debugCursor++
		}
	case "g":
		m.debugCursor = 0
	case "G":
		if n > 0 {
			m.debugCursor = n - 1
		}
	case "enter":
		if n > 0 {
			m.debugDetailMode = true
			m.debugBodyScroll = 0
		}

	// ── Filter ──────────────────────────────────────────────────────────
	case "/":
		m.debugFiltering = true

	// ── Escape: clear filter → unfocus ──────────────────────────────────
	case "esc":
		if m.debugFiltering {
			m.debugFiltering = false
			m.debugFilter = ""
			m.debugCursor = 0
		} else {
			m.debugFocused = false
		}

	// ── Filter input ────────────────────────────────────────────────────
	default:
		if m.debugFiltering {
			switch msg.String() {
			case "backspace":
				runes := []rune(m.debugFilter)
				if len(runes) > 0 {
					m.debugFilter = string(runes[:len(runes)-1])
				}
			case "enter":
				m.debugFiltering = false
				m.debugCursor = 0
			default:
				if isPrintable(msg.String()) {
					m.debugFilter += msg.String()
				}
			}
		}
	}
	return m, nil
}

// isPrintable returns true for single-rune printable key strings.
func isPrintable(s string) bool {
	if len([]rune(s)) != 1 {
		return false
	}
	r := []rune(s)[0]
	return r >= 32 && r != 127
}

// ── Panel-wide style set ──────────────────────────────────────────────────────

// debugPanelStyles holds every style used inside the debug panel so that the
// entire viewport background can be switched as a unit when focus changes.
// Focused state uses colorHeaderBg (slightly lighter) to make the active panel
// immediately visually distinct from the left content panel.
type debugPanelStyles struct {
	bg          lipgloss.Style // base background + default foreground
	punct       lipgloss.Style // dim decorative text (section headers, timestamps, hints)
	row         lipgloss.Style // non-selected list rows
	divider     lipgloss.Style // horizontal ─── divider line
	placeholder lipgloss.Style // empty-state italic placeholder text
	panelBg     color.Color // raw background colour (for method/status helpers)
}

// newDebugStyles returns the style set that matches the current focus state.
func (m Model) newDebugStyles() debugPanelStyles {
	if m.debugFocused {
		bg := colorHeaderBg
		return debugPanelStyles{
			bg:          lipgloss.NewStyle().Background(bg).Foreground(colorFg),
			punct:       lipgloss.NewStyle().Background(bg).Foreground(colorDim),
			row:         lipgloss.NewStyle().Background(bg).Foreground(colorFg),
			divider:     lipgloss.NewStyle().Background(bg).Foreground(colorBorder),
			placeholder: lipgloss.NewStyle().Background(bg).Foreground(colorDim).Italic(true),
			panelBg:     bg,
		}
	}
	return debugPanelStyles{
		bg:          contentStyle,
		punct:       jsonPunctStyle,
		row:         tableRowStyle,
		divider:     contentDividerStyle,
		placeholder: contentPlaceholderStyle,
		panelBg:     colorBg,
	}
}

// ── Panel dispatcher ──────────────────────────────────────────────────────────

// renderDebugPanel dispatches to the list or detail renderer.
func (m Model) renderDebugPanel() string {
	ds := m.newDebugStyles()
	if m.debugDetailMode {
		return m.renderDebugDetail(ds)
	}
	return m.renderDebugList(ds)
}

// ── List view ─────────────────────────────────────────────────────────────────

// renderDebugList renders the full-height scrollable call list.
func (m Model) renderDebugList(ds debugPanelStyles) string {
	pw := m.debugPanelWidth()
	h := m.contentHeight()
	events := m.filteredDebugEvents()
	n := len(events)

	// Clamp cursor.
	cursor := m.debugCursor
	if n == 0 {
		cursor = 0
	} else if cursor >= n {
		cursor = n - 1
	}

	var lines []string

	// ── Title bar — styled to show which panel is active ──────────────────
	// Focused: accent bg + blue title text acts as a visible "active" tab.
	// Unfocused: dim text on content bg recedes visually.
	titleStyle := debugTitleUnfocusedStyle
	hintStyle := jsonPunctStyle
	if m.debugFocused {
		titleStyle = debugTitleFocusedStyle
		hintStyle = filterBarStyle // colorHeaderBg bg, colorDim fg — matches title bg
	}
	focusedHint := ""
	if !m.debugFocused {
		focusedHint = hintStyle.Render("  [Tab] focus  ")
	} else {
		focusedHint = hintStyle.Render("  [Tab] main  [l] close  ")
	}
	titleLeft := titleStyle.Render("  API Inspector")
	gap := pw - lipgloss.Width(titleLeft) - lipgloss.Width(focusedHint)
	if gap < 0 {
		gap = 0
	}
	lines = append(lines, titleLeft+titleStyle.Width(gap).Render("")+focusedHint)

	// ── Filter bar (inline with title, shown when active) ─────────────────
	if m.debugFiltering || m.debugFilter != "" {
		curs := ""
		if m.debugFiltering {
			curs = "▌"
		}
		label := filterBarStyle.Render("  / ")
		text := filterBarActiveStyle.Render(m.debugFilter + curs)
		filterLine := label + text
		w := pw - lipgloss.Width(label) - lipgloss.Width(text)
		if w > 0 {
			filterLine += filterBarStyle.Width(w).Render("")
		}
		lines = append(lines, filterLine)
	}

	// ── Divider ───────────────────────────────────────────────────────────
	lines = append(lines, ds.divider.Width(pw).Render(strings.Repeat("─", pw)))

	// ── Call list rows ────────────────────────────────────────────────────
	listH := h - len(lines) // remaining rows for the list
	if listH < 1 {
		listH = 1
	}

	// Keep cursor visible.
	listOffset := 0
	if cursor >= listH {
		listOffset = cursor - listH + 1
	}

	for i := 0; i < listH; i++ {
		idx := listOffset + i
		if idx >= n {
			lines = append(lines, ds.row.Width(pw).Render(""))
			continue
		}
		e := events[idx]
		selected := idx == cursor
		lines = append(lines, m.renderDebugCallRow(e, pw, selected, ds))
	}

	// Pad to exactly h lines.
	for len(lines) < h {
		lines = append(lines, ds.bg.Width(pw).Render(""))
	}
	return strings.Join(lines[:h], "\n")
}

// ── Detail view ───────────────────────────────────────────────────────────────

// renderDebugDetail renders the full request/response detail for the selected call.
func (m Model) renderDebugDetail(ds debugPanelStyles) string {
	pw := m.debugPanelWidth()
	h := m.contentHeight()
	events := m.filteredDebugEvents()

	cursor := m.debugCursor
	if cursor >= len(events) {
		cursor = len(events) - 1
	}

	var lines []string

	// ── Title bar — always focused when in detail mode ────────────────────
	escHint := filterBarStyle.Render("  [esc] list  [Tab] main  ")
	titleLeft := debugTitleFocusedStyle.Render("  API Inspector › Detail")
	gap := pw - lipgloss.Width(titleLeft) - lipgloss.Width(escHint)
	if gap < 0 {
		gap = 0
	}
	lines = append(lines, titleLeft+debugTitleFocusedStyle.Width(gap).Render("")+escHint)

	// ── Divider ───────────────────────────────────────────────────────────
	lines = append(lines, ds.divider.Width(pw).Render(strings.Repeat("─", pw)))

	// ── Body ──────────────────────────────────────────────────────────────
	bodyH := h - len(lines)
	bodyLines := m.buildDebugBody(events, cursor, pw, ds)

	// Clamp scroll.
	scroll := m.debugBodyScroll
	if scroll >= len(bodyLines) && len(bodyLines) > 0 {
		scroll = len(bodyLines) - 1
	}
	if scroll > 0 {
		bodyLines = bodyLines[scroll:]
	}

	for i := 0; i < bodyH; i++ {
		if i < len(bodyLines) {
			lines = append(lines, bodyLines[i])
		} else {
			lines = append(lines, ds.bg.Width(pw).Render(""))
		}
	}

	// Pad to exactly h lines.
	for len(lines) < h {
		lines = append(lines, ds.bg.Width(pw).Render(""))
	}
	return strings.Join(lines[:h], "\n")
}

// ── Call list row ─────────────────────────────────────────────────────────────

// renderDebugCallRow renders one row in the call list.
// Layout: [cursor] [METHOD] [path…] [status] [duration]
func (m Model) renderDebugCallRow(e client.APIEvent, pw int, selected bool, ds debugPanelStyles) string {
	methodStr := debugMethodLabel(e.Method) // fixed 7 chars

	var statusStr string
	if e.Err != "" {
		statusStr = "ERR"
	} else {
		statusStr = fmt.Sprintf("%d", e.StatusCode)
	}
	// Pad duration to a fixed 6-char field (right-aligned) so the status and
	// duration columns are always the same width regardless of the value:
	//   "78ms"  → "  78ms"
	//   "512ms" → " 512ms"
	//   "2.0s"  → "  2.0s"
	durStr := fmt.Sprintf("%6s", debugDurLabel(e.Duration))

	// Fixed right columns: " STATUS  DUR  " — now always the same width.
	rightFixed := "  " + statusStr + "  " + durStr + "  "
	rightW := lipgloss.Width(rightFixed)

	// Available for path (-2 for cursor mark, -2 gap after method)
	pathW := pw - 2 - lipgloss.Width(methodStr) - 2 - rightW
	if pathW < 4 {
		pathW = 4
	}
	pathStr := truncateStr(e.Path, pathW)

	base := ds.row
	cursorMark := "  "
	if selected {
		base = tableRowSelectedStyle
		cursorMark = "▶ "
	}

	var row string
	if selected {
		inner := cursorMark + methodStr + "  " + pathStr
		row = base.Render(inner) +
			base.Width(pathW-lipgloss.Width(pathStr)).Render("") +
			base.Render(rightFixed)
	} else {
		statusStyled := debugStatusStyle(e.StatusCode, e.Err, ds.panelBg).Render(statusStr)
		durStyled := ds.punct.Render(durStr)
		rightStyled := base.Render("  ") + statusStyled + base.Render("  ") + durStyled + base.Render("  ")

		gap := pathW - lipgloss.Width(pathStr)
		if gap < 0 {
			gap = 0
		}
		row = base.Render(cursorMark) +
			debugMethodStyle(e.Method, ds.panelBg).Render(methodStr) +
			base.Render("  "+pathStr) +
			base.Width(gap).Render("") +
			rightStyled
	}

	// Pad row to full panel width.
	rowW := lipgloss.Width(row)
	if rowW < pw {
		row += base.Width(pw - rowW).Render("")
	}
	return row
}

// ── Body builder ──────────────────────────────────────────────────────────────

// buildDebugBody returns request+response viewer lines for the selected event.
// Layout mirrors curl -v: method+path → timestamp → request headers → body →
// status+duration → response headers → response body.
func (m Model) buildDebugBody(events []client.APIEvent, cursor int, pw int, ds debugPanelStyles) []string {
	var lines []string

	if len(events) == 0 {
		lines = append(lines, ds.placeholder.Width(pw).Render("  No API calls yet."))
		return lines
	}
	if cursor < 0 || cursor >= len(events) {
		return lines
	}
	e := events[cursor]

	// addSection renders a dim "── LABEL ─────" divider line.
	addSection := func(label string) {
		lbl := ds.punct.Render("  ── " + label + " ")
		rem := pw - lipgloss.Width(lbl)
		if rem < 0 {
			rem = 0
		}
		lines = append(lines, lbl+ds.punct.Width(rem).Render(strings.Repeat("─", rem)))
	}

	// addHeaders renders each "Name: value" header line, dim-styled.
	addHeaders := func(headers []string) {
		for _, h := range headers {
			line := truncateStr(h, pw-2)
			lines = append(lines, ds.punct.Render("  "+line))
		}
	}

	// ── REQUEST ───────────────────────────────────────────────────────────
	addSection("REQUEST")

	// Method + path — wrapped across multiple lines when the URL is long.
	// Continuation lines are indented to align under the path start.
	//   GET  /api/v2/organizations/org/workspaces/ws-longid/runs?page_number=1
	//        &filter[...]=...
	methodStyled := debugMethodStyle(e.Method, ds.panelBg).Render(e.Method)
	chunkW := pw - len(e.Method) - 2 // width per line (same for first and continuations)
	if chunkW < 1 {
		chunkW = 1
	}
	indent := strings.Repeat(" ", len(e.Method)+2)
	remaining := e.Path
	first := remaining
	if len(first) > chunkW {
		first = remaining[:chunkW]
		remaining = remaining[chunkW:]
	} else {
		remaining = ""
	}
	lines = append(lines, methodStyled+ds.bg.Render("  "+first))
	for len(remaining) > 0 {
		chunk := remaining
		if len(chunk) > chunkW {
			chunk = remaining[:chunkW]
			remaining = remaining[chunkW:]
		} else {
			remaining = ""
		}
		lines = append(lines, ds.bg.Render(indent+chunk))
	}

	// Timestamp
	lines = append(lines, ds.punct.Render("  "+e.Timestamp.Format(time.RFC3339)))

	// Request headers
	if len(e.ReqHeaders) > 0 {
		lines = append(lines, ds.bg.Width(pw).Render(""))
		addHeaders(e.ReqHeaders)
	}

	// Request body (POST/PATCH/PUT)
	if e.ReqBody != "" {
		lines = append(lines, ds.bg.Width(pw).Render(""))
		for _, bl := range strings.Split(e.ReqBody, "\n") {
			lines = append(lines, ds.bg.Render("  ")+colorizeJSONLine(truncateStr(bl, pw-2)))
		}
	}

	lines = append(lines, ds.bg.Width(pw).Render(""))

	// ── RESPONSE ──────────────────────────────────────────────────────────
	addSection("RESPONSE")

	if e.Err != "" {
		lines = append(lines, statusErrorStyle.Render("  ✗  "+e.Err))
		return lines
	}

	// Status + duration
	statusStr := fmt.Sprintf("%d", e.StatusCode)
	statusStyled := debugStatusStyle(e.StatusCode, "", ds.panelBg).Render(statusStr)
	durStyled := ds.punct.Render("  " + debugDurLabel(e.Duration))
	lines = append(lines, statusStyled+durStyled)

	// Response headers
	if len(e.RespHeaders) > 0 {
		lines = append(lines, ds.bg.Width(pw).Render(""))
		addHeaders(e.RespHeaders)
	}

	// Response body
	lines = append(lines, ds.bg.Width(pw).Render(""))
	if e.RespBody != "" {
		for _, bl := range strings.Split(e.RespBody, "\n") {
			lines = append(lines, ds.bg.Render("  ")+colorizeJSONLine(truncateStr(bl, pw-2)))
		}
	} else {
		lines = append(lines, ds.punct.Render("  (empty body)"))
	}

	return lines
}

// ── Style helpers ─────────────────────────────────────────────────────────────

func debugMethodLabel(method string) string {
	switch strings.ToUpper(method) {
	case "GET":
		return "GET    "
	case "POST":
		return "POST   "
	case "PATCH":
		return "PATCH  "
	case "PUT":
		return "PUT    "
	case "DELETE":
		return "DELETE "
	default:
		if len(method) < 7 {
			return method + strings.Repeat(" ", 7-len(method))
		}
		return method[:7]
	}
}

// debugMethodStyle returns a coloured style for the HTTP method verb.
// bg is the panel's current background colour so the glyph blends in correctly
// regardless of whether the panel is focused (colorHeaderBg) or not (colorBg).
func debugMethodStyle(method string, bg color.Color) lipgloss.Style {
	base := lipgloss.NewStyle().Background(bg)
	switch strings.ToUpper(method) {
	case "GET":
		return base.Foreground(colorAccent) // blue
	case "POST":
		return base.Foreground(colorSuccess) // green
	case "DELETE":
		return base.Foreground(colorError) // red
	case "PATCH", "PUT":
		return base.Foreground(colorLoading) // amber
	default:
		return base.Foreground(colorDim)
	}
}

// debugStatusStyle returns a coloured style for the HTTP status code.
// bg is the panel's current background colour (see debugMethodStyle).
func debugStatusStyle(code int, errStr string, bg color.Color) lipgloss.Style {
	base := lipgloss.NewStyle().Background(bg)
	if errStr != "" {
		return base.Foreground(colorError)
	}
	switch {
	case code >= 200 && code < 300:
		return base.Foreground(colorSuccess)
	case code >= 400 && code < 500:
		return base.Foreground(colorLoading)
	case code >= 500:
		return base.Foreground(colorError)
	default:
		return base.Foreground(colorDim)
	}
}

func debugDurLabel(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%dµs", d.Microseconds())
	}
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.1fs", d.Seconds())
}

// joinPanels zips the left (main content) and right (debug panel) strings
// line-by-line with a dim vertical-bar separator.
//
// It is a model method so it can access mainWidth() and debugPanelWidth() and
// ENFORCE exact column widths on every line from both sides. This makes the two
// panels mutually isolated: no renderer bug or off-by-one can bleed across the
// separator — short lines are padded, and each lane is always exactly its target
// width. This is the single source of truth for panel boundary enforcement.
func (m Model) joinPanels(left, right string) string {
	leftW := m.mainWidth()
	rightW := m.debugPanelWidth()

	// The right panel padding background matches the focused state so that
	// any short lines in the right panel are filled with the same bg colour
	// as the rest of the panel, not the left panel's colour.
	ds := m.newDebugStyles()

	leftLines := strings.Split(left, "\n")
	rightLines := strings.Split(right, "\n")
	maxLines := len(leftLines)
	if len(rightLines) > maxLines {
		maxLines = len(rightLines)
	}

	// Separator color tracks focus: accent when right panel is active,
	// dim when left panel is active — provides an extra focus cue.
	sepStyle := jsonPunctStyle
	if m.debugFocused {
		sepStyle = lipgloss.NewStyle().Background(colorBg).Foreground(colorAccent)
	}
	sep := sepStyle.Render("│")
	out := make([]string, maxLines)
	for i := 0; i < maxLines; i++ {
		l, r := "", ""
		if i < len(leftLines) {
			l = leftLines[i]
		}
		if i < len(rightLines) {
			r = rightLines[i]
		}

		// Enforce left panel width: truncate if too wide, pad if too narrow.
		// ansi.Truncate is ANSI-escape-code-aware and safely clips styled strings.
		lw := lipgloss.Width(l)
		switch {
		case lw > leftW:
			l = ansi.Truncate(l, leftW, "")
		case lw < leftW:
			l += contentStyle.Width(leftW - lw).Render("")
		}

		// Enforce right panel width: pad short lines with the panel's own
		// background (focused or not) so no left-panel colour bleeds through.
		rw := lipgloss.Width(r)
		switch {
		case rw > rightW:
			r = ansi.Truncate(r, rightW, "")
		case rw < rightW:
			r += ds.bg.Width(rightW - rw).Render("")
		}

		out[i] = l + sep + r
	}
	return strings.Join(out, "\n")
}
