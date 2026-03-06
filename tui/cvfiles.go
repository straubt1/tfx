// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"fmt"
	"path/filepath"
	"strings"
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
//	"├── "        (top-level, non-last)
//	"│   └── "   (depth 1, last child, parent has siblings)
//	"    ├── "   (depth 1, non-last child, parent was last)
//
// Inspired by the charm.land/lipgloss/v2/tree DefaultEnumerator style.
func buildTreeConnector(files []cvFile, i int) string {
	d := files[i].depth()
	var sb strings.Builder
	// Continuation lines for each ancestor level.
	for level := 0; level < d; level++ {
		if isLastChild(files, i, level) {
			sb.WriteString("    ") // ancestor was last child — no vertical line needed
		} else {
			sb.WriteString("│   ") // ancestor has siblings below — draw continuation
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

// ── File browser rendering ─────────────────────────────────────────────────────

// cvFilesVisibleRows returns the number of file-browser rows that fit on screen.
func (m Model) cvFilesVisibleRows() int {
	h := m.contentHeight() - 2 // header + divider
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
				lines[i] = contentPlaceholderStyle.Width(m.width).Render("  " + frame + "  Downloading and extracting config version…")
			} else {
				lines[i] = contentStyle.Width(m.width).Render("")
			}
		}
		return strings.Join(lines, "\n")
	}

	// ── Error ─────────────────────────────────────────────────────────────────
	if m.cvFileErr != "" {
		lines := make([]string, h)
		for i := range lines {
			if i == 0 {
				lines[i] = contentStyle.Width(m.width).Render("  ✗  " + m.cvFileErr)
			} else {
				lines[i] = contentStyle.Width(m.width).Render("")
			}
		}
		return strings.Join(lines, "\n")
	}

	// ── File tree ─────────────────────────────────────────────────────────────
	// Layout: nameColW (connector + name) + gap(2) + sizeColW(10)
	// No separate cursor column — selection is shown via row highlight colour.
	const sizeColW = 10
	const gapW = 2
	nameColW := m.width - sizeColW - gapW
	if nameColW < 12 {
		nameColW = 12
	}

	var lines []string

	// Header + divider (indent header label to align with non-connector depth-0 prefix).
	hdr := m.pad(
		tableHeaderStyle.Render("    ")+ // 4-char indent to align with "├── " / "└── "
			tableHeaderStyle.Width(nameColW-4).Render("NAME")+
			tableHeaderStyle.Render("  ")+
			tableHeaderStyle.Width(sizeColW).Render("SIZE"),
		tableHeaderStyle,
	)
	lines = append(lines, hdr)
	lines = append(lines, m.renderTableDivider())

	if len(m.cvFiles) == 0 {
		lines = append(lines, contentPlaceholderStyle.Width(m.width).Render("  No files found."))
	} else {
		vis := m.cvFilesVisibleRows()
		end := m.cvFileOffset + vis
		if end > len(m.cvFiles) {
			end = len(m.cvFiles)
		}
		for i := m.cvFileOffset; i < end; i++ {
			f := m.cvFiles[i]
			selected := i == m.cvFileCursor
			tfFile := !f.isDir && isHCLFile(f.displayName())

			// Build connector string and measure its width (all ASCII → byte len).
			connStr := buildTreeConnector(m.cvFiles, i)
			connW := len(connStr)

			// Available width for the name portion after the connector.
			availW := nameColW - connW
			if availW < 0 {
				availW = 0
				connStr = connStr[:nameColW] // truncate connector on very narrow terminals
				connW = nameColW
			}

			displayN := f.displayName()
			size := f.sizeStr()

			// Compute name portion (with optional ◆ icon for HCL files).
			var namePart string
			switch {
			case !selected && tfFile:
				namePart = truncateStr("◆ "+displayN, availW)
			default:
				namePart = truncateStr(displayN, availW)
			}

			// Build the row from separately-styled segments.
			// Connectors are always dim; item text colour depends on type & selection.
			var row string
			switch {
			case selected:
				// Entire row in selection style.
				row = tableRowSelectedStyle.Render(connStr) +
					tableRowSelectedStyle.Width(availW).Render(namePart) +
					tableRowSelectedStyle.Render("  ") +
					tableRowSelectedStyle.Width(sizeColW).Render(size)
				lines = append(lines, m.pad(row, tableRowSelectedStyle))
			case tfFile:
				// Connectors dim, name/icon purple.
				row = jsonPunctStyle.Render(connStr) +
					hclFileStyle.Width(availW).Render(namePart) +
					tableRowStyle.Render("  ") +
					tableRowStyle.Width(sizeColW).Render(size)
				lines = append(lines, m.pad(row, tableRowStyle))
			case f.isDir:
				// Connectors dim, directory name bold/bright.
				row = jsonPunctStyle.Render(connStr) +
					contentTitleStyle.Width(availW).Render(namePart) +
					tableRowStyle.Render("  ") +
					tableRowStyle.Width(sizeColW).Render(size)
				lines = append(lines, m.pad(row, tableRowStyle))
			default:
				// Connectors dim, regular file name in default fg.
				row = jsonPunctStyle.Render(connStr) +
					tableRowStyle.Width(availW).Render(namePart) +
					tableRowStyle.Render("  ") +
					tableRowStyle.Width(sizeColW).Render(size)
				lines = append(lines, m.pad(row, tableRowStyle))
			}
		}
	}

	for len(lines) < h {
		lines = append(lines, contentStyle.Width(m.width).Render(""))
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
			lines[i] = contentStyle.Width(m.width).Render("")
		}
		return strings.Join(lines, "\n")
	}

	numLines := len(m.cvFileLines)
	lineNumWidth := len(fmt.Sprintf("%d", numLines))
	// Layout: 2 margin + lineNumWidth + " │ " (3) + content
	contentWidth := m.width - 2 - lineNumWidth - 3
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
		all = append(all, contentStyle.Width(m.width).Render(row))
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
		out[i] = contentStyle.Width(m.width).Render("")
	}
	return strings.Join(out, "\n")
}
