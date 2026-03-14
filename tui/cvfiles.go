// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"fmt"
	"path/filepath"
	"strings"

	lipgloss "charm.land/lipgloss/v2"
	devicons "github.com/epilande/go-devicons"
)

// cvFile represents a single entry (file or directory) in an extracted config version.
type cvFile struct {
	relPath string // path relative to extraction root, e.g., "modules/main.tf"
	size    int64
	isDir   bool
}

// depth returns the nesting level (0 = top-level entry).
func (f cvFile) depth() int {
	if f.relPath == "" {
		return 0
	}
	return strings.Count(f.relPath, string(filepath.Separator))
}

// displayName returns the base name with a trailing "/" for directories.
func (f cvFile) displayName() string {
	name := filepath.Base(f.relPath)
	if f.isDir {
		return name + "/"
	}
	return name
}

// sizeStr returns a human-readable size for display.
func (f cvFile) sizeStr() string {
	if f.isDir {
		return "(dir)"
	}
	if f.size < 1024 {
		return fmt.Sprintf("%d B", f.size)
	}
	if f.size < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(f.size)/1024)
	}
	return fmt.Sprintf("%.1f MB", float64(f.size)/(1024*1024))
}

// ── Tree connector helpers ─────────────────────────────────────────────────────

// isLastChild reports whether the item at index i is the last sibling at depth `level`.
// When level == files[i].depth() it checks the file itself.
// When level < files[i].depth() it checks the ancestor at that depth.
func isLastChild(files []cvFile, i, level int) bool {
	if i >= len(files) {
		return true
	}
	fParts := strings.Split(filepath.ToSlash(files[i].relPath), "/")
	if level >= len(fParts) {
		return true
	}
	for j := i + 1; j < len(files); j++ {
		jParts := strings.Split(filepath.ToSlash(files[j].relPath), "/")
		if len(jParts) <= level {
			continue // shallower than this level — cannot be a sibling here
		}
		// Verify the first `level` path components match (same parent).
		sameParent := true
		for k := 0; k < level; k++ {
			if k >= len(jParts) || jParts[k] != fParts[k] {
				sameParent = false
				break
			}
		}
		if sameParent && jParts[level] != fParts[level] {
			return false // found a sibling after this item
		}
	}
	return true
}

// buildTreeConnector returns the box-drawing connector prefix for file i, e.g.:
//
//	""            (depth 0 — root items have no connector)
//	"├── "        (depth 1, non-last child)
//	"└── "        (depth 1, last child)
//	"    ├── "   (depth 2, non-last child, depth-1 parent was last)
//	"│   └── "   (depth 2, last child, depth-1 parent has siblings)
//
// Root-level (depth 0) items return an empty string — the icon alone signals
// the entry; connectors are only drawn for nested children.
// When computing continuation lines, depth-0 ancestors are treated as if they
// always have no vertical bar (since they carry no connector of their own).
//
// Inspired by the charm.land/lipgloss/v2/tree DefaultEnumerator style.
func buildTreeConnector(files []cvFile, i int) string {
	d := files[i].depth()
	if d == 0 {
		return "" // root items: no connector
	}
	var sb strings.Builder
	// Continuation lines for ancestor levels 1..d-1.
	// We skip level 0 (root) because root items carry no visual connector —
	// writing "    " for them would add an unnecessary 4-space indent to every
	// child, pushing connectors away from the left edge.
	for level := 1; level < d; level++ {
		if isLastChild(files, i, level) {
			sb.WriteString("    ") // last-child ancestor — no vertical bar continues
		} else {
			sb.WriteString("│   ") // ancestor has siblings below — draw continuation bar
		}
	}
	// Branch connector for this item.
	if isLastChild(files, i, d) {
		sb.WriteString("└── ")
	} else {
		sb.WriteString("├── ")
	}
	return sb.String()
}

// cvFileIcon returns the Nerd Font icon glyph and its hex colour for a file entry.
// For directories the standard Nerd Font folder icon is used (matches go-devicons DirStyle).
// For files, go-devicons performs an extension/name lookup; os.Lstat is called with just the
// display name so it always falls back to the name-based lookup path, which is what we want.
func cvFileIcon(f cvFile) (string, string) {
	if f.isDir {
		return "\uf07b", "#61AFEF" // Nerd Font folder icon (matches devicons.DirStyle)
	}
	s := devicons.IconForPath(f.displayName()) // Lstat fails → name/ext lookup
	return s.Icon, s.Color
}

// ── File browser rendering ─────────────────────────────────────────────────────

// cvFilesPathBar renders the context bar shown above the column headers,
// inspired by superfile's top-bar + position counter design.
//
// Left side  — "▸  <parent-dir-of-selected-item>"  (left-truncated so the
//              deepest directory is always visible, à la superfile)
// Right side — "<cursor+1>/<total>" position indicator
func (m Model) cvFilesPathBar() string {
	// Determine the parent directory of the currently selected item.
	pathText := "/"
	if m.cvFileCursor >= 0 && m.cvFileCursor < len(m.cvFiles) {
		sel := m.cvFiles[m.cvFileCursor]
		dir := filepath.ToSlash(filepath.Dir(sel.relPath))
		if dir != "" && dir != "." {
			pathText = dir + "/"
		}
	}

	// Position indicator (right side).
	total := len(m.cvFiles)
	posText := ""
	if total > 0 {
		cur := m.cvFileCursor + 1
		if cur > total {
			cur = total
		}
		posText = fmt.Sprintf("%d/%d", cur, total)
	}

	// Render the position label (fixed width, right-aligned).
	posRendered := detailLabelStyle.Render("  " + posText + "  ")
	posW := lipgloss.Width(posRendered)

	// Available rune-width for the left path portion.
	leftAvailW := m.innerWidth() - posW - 2 // 2 for "  " left margin

	// Build the glyph + path text, then left-truncate so the tail is always visible.
	// Styled: ▸ in accent blue, path text in default fg.
	glyphStr := "▸  "
	fullPath := glyphStr + pathText
	truncated := truncateStrLeft(fullPath, leftAvailW)

	// Re-split at the glyph boundary for independent colouring.
	// If truncation ate into the glyph, just render the whole thing in accent.
	var pathRendered string
	glyphRunes := []rune(glyphStr)
	truncRunes := []rune(truncated)
	if len(truncRunes) >= len(glyphRunes) && string(truncRunes[:len(glyphRunes)]) == glyphStr {
		// Glyph survived truncation — colour separately.
		pathRendered = contentStyle.Render("  ") +
			breadcrumbActiveStyle.Render(glyphStr) +
			contentStyle.Render(string(truncRunes[len(glyphRunes):]))
	} else {
		// Glyph was truncated away (very narrow terminal) — just render in accent.
		pathRendered = contentStyle.Render("  ") +
			breadcrumbActiveStyle.Render(truncated)
	}

	// Fill the gap between the path and the position indicator.
	pathW := lipgloss.Width(pathRendered)
	gapW := m.innerWidth() - pathW - posW
	if gapW < 0 {
		gapW = 0
	}

	return pathRendered + contentStyle.Width(gapW).Render("") + posRendered
}

// cvFilesVisibleRows returns the number of file-browser rows that fit on screen.
func (m Model) cvFilesVisibleRows() int {
	h := m.contentHeight() - 3 // path bar + header + divider
	if h < 1 {
		return 1
	}
	return h
}

// renderConfigVersionFilesContent renders the config-version file tree browser
// with Unicode box-drawing connectors (├── / └── / │) à la eza / lipgloss tree.
func (m Model) renderConfigVersionFilesContent() string {
	h := m.contentHeight()

	// ── Loading ───────────────────────────────────────────────────────────────
	if m.cvFileLoading {
		lines := make([]string, h)
		mid := h / 2
		frame := spinnerFrames[m.spinnerIdx]
		for i := range lines {
			if i == mid {
				lines[i] = contentPlaceholderStyle.Width(m.innerWidth()).Render("  " + frame + "  Downloading and extracting config version…")
			} else {
				lines[i] = contentStyle.Width(m.innerWidth()).Render("")
			}
		}
		return strings.Join(lines, "\n")
	}

	// ── Error ─────────────────────────────────────────────────────────────────
	if m.cvFileErr != "" {
		lines := make([]string, h)
		for i := range lines {
			if i == 0 {
				lines[i] = contentStyle.Width(m.innerWidth()).Render("  ✗  " + m.cvFileErr)
			} else {
				lines[i] = contentStyle.Width(m.innerWidth()).Render("")
			}
		}
		return strings.Join(lines, "\n")
	}

	// ── File tree ─────────────────────────────────────────────────────────────
	// Layout: nameColW (connector + name) + gap(2) + sizeColW(10)
	// No separate cursor column — selection is shown via row highlight colour.
	const sizeColW = 10
	const gapW = 2
	nameColW := m.innerWidth() - sizeColW - gapW
	if nameColW < 12 {
		nameColW = 12
	}

	var lines []string

	// Path context bar (superfile-inspired: current dir + position indicator).
	lines = append(lines, m.cvFilesPathBar())

	// Header + divider. NAME spans the full name column (root items have no connector prefix).
	hdr := m.padContent(
		tableHeaderStyle.Width(nameColW).Render("NAME")+
			tableHeaderStyle.Render("  ")+
			tableHeaderStyle.Width(sizeColW).Align(lipgloss.Right).Render("SIZE"),
		tableHeaderStyle,
	)
	lines = append(lines, hdr)
	lines = append(lines, m.renderTableDivider())

	if len(m.cvFiles) == 0 {
		lines = append(lines, contentPlaceholderStyle.Width(m.innerWidth()).Render("  No files found."))
	} else {
		vis := m.cvFilesVisibleRows()
		end := m.cvFileOffset + vis
		if end > len(m.cvFiles) {
			end = len(m.cvFiles)
		}
		for i := m.cvFileOffset; i < end; i++ {
			f := m.cvFiles[i]
			selected := i == m.cvFileCursor

			// Build connector string and measure its width (all ASCII → byte len).
			connStr := buildTreeConnector(m.cvFiles, i)
			connW := len(connStr)

			// Icon lookup (Nerd Font glyph + hex color for this file/dir type).
			iconGlyph, iconColor := cvFileIcon(f)
			iconStyle := lipgloss.NewStyle().Background(colorBg).Foreground(lipgloss.Color(iconColor))
			iconRendered := iconStyle.Render(iconGlyph + " ") // glyph + 1 trailing space
			iconW := lipgloss.Width(iconRendered)             // typically 2

			// Available width for the name portion after connector + icon.
			availW := nameColW - connW - iconW
			if availW < 0 {
				availW = 0
				if connW > nameColW {
					connStr = connStr[:nameColW] // truncate connector on very narrow terminals
					connW = nameColW
				}
			}

			displayN := f.displayName()
			size := f.sizeStr()
			namePart := truncateStr(displayN, availW)

			// Build the row from separately-styled segments.
			// Connectors are always dim; icon carries the file-type colour; name depends on
			// type (dir = bold, selected = selection highlight, file = default fg).
			var row string
			switch {
			case selected:
				// Entire row in selection style — icon loses its type colour on highlight.
				row = tableRowSelectedStyle.Render(connStr+iconGlyph+" ") +
					tableRowSelectedStyle.Width(availW).Render(namePart) +
					tableRowSelectedStyle.Render("  ") +
					tableRowSelectedStyle.Width(sizeColW).Align(lipgloss.Right).Render(size)
				lines = append(lines, m.padContent(row, tableRowSelectedStyle))
			case f.isDir:
				// Connector dim, icon its colour, directory name bold/bright.
				row = jsonPunctStyle.Render(connStr) +
					iconRendered +
					contentTitleStyle.Width(availW).Render(namePart) +
					tableRowStyle.Render("  ") +
					tableRowStyle.Width(sizeColW).Align(lipgloss.Right).Render(size)
				lines = append(lines, m.padContent(row, tableRowStyle))
			default:
				// Connector dim, icon its colour, file name in default fg.
				row = jsonPunctStyle.Render(connStr) +
					iconRendered +
					tableRowStyle.Width(availW).Render(namePart) +
					tableRowStyle.Render("  ") +
					tableRowStyle.Width(sizeColW).Align(lipgloss.Right).Render(size)
				lines = append(lines, m.padContent(row, tableRowStyle))
			}
		}
	}

	for len(lines) < h {
		lines = append(lines, contentStyle.Width(m.innerWidth()).Render(""))
	}
	return strings.Join(lines[:h], "\n")
}

// renderConfigVersionFileContent renders the scrollable line-numbered content
// viewer for a single file from the config version archive.
// .json files use JSON syntax highlighting; .tf/.tfvars/.hcl use HCL highlighting.
func (m Model) renderConfigVersionFileContent() string {
	h := m.contentHeight()

	if len(m.cvFileLines) == 0 {
		lines := make([]string, h)
		for i := range lines {
			lines[i] = contentStyle.Width(m.innerWidth()).Render("")
		}
		return strings.Join(lines, "\n")
	}

	numLines := len(m.cvFileLines)
	lineNumWidth := len(fmt.Sprintf("%d", numLines))
	// Layout: 2 margin + lineNumWidth + " │ " (3) + content
	contentWidth := m.innerWidth() - 2 - lineNumWidth - 3
	if contentWidth < 10 {
		contentWidth = 10
	}

	isJSON := strings.HasSuffix(strings.ToLower(m.cvFileName), ".json")
	isHCL := isHCLFile(m.cvFileName)

	all := make([]string, 0, numLines)
	for i, line := range m.cvFileLines {
		lineNum := fmt.Sprintf("%*d", lineNumWidth, i+1)
		display := line
		if len(display) > contentWidth {
			display = display[:contentWidth-1] + "…"
		}

		var colored string
		switch {
		case isJSON:
			colored = colorizeJSONLine(display)
		case isHCL:
			colored = colorizeHCLLine(display)
		default:
			colored = contentStyle.Render(display)
		}

		row := "  " +
			detailLabelStyle.Render(lineNum) +
			contentDividerStyle.Render(" │ ") +
			colored
		all = append(all, contentStyle.Width(m.innerWidth()).Render(row))
	}

	// Clamp scroll and slice visible window.
	maxScroll := len(all) - h
	if maxScroll < 0 {
		maxScroll = 0
	}
	start := m.cvFileScroll
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
