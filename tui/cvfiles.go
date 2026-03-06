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

// cvFilesVisibleRows returns the number of file-browser rows that fit on screen.
func (m Model) cvFilesVisibleRows() int {
	h := m.contentHeight() - 2 // header + divider
	if h < 1 {
		return 1
	}
	return h
}

// renderConfigVersionFilesContent renders the config-version file tree browser.
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
	const sizeColW = 10 // right-aligned size column width
	nameColW := m.width - 2 - sizeColW - 2 // cursor(2) + name + gap(2) + size
	if nameColW < 10 {
		nameColW = 10
	}

	var lines []string

	// Header + divider
	hdr := m.pad(
		tableHeaderStyle.Render("  ")+
			tableHeaderStyle.Width(nameColW).Render("NAME")+
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

			indent := strings.Repeat("  ", f.depth())
			displayN := f.displayName()
			size := f.sizeStr()

			style := tableRowStyle
			cursor := "  "
			if selected {
				style = tableRowSelectedStyle
				cursor = "> "
			}

			var row string
			if !selected && tfFile {
				// Terraform/HCL files: diamond icon + purple name (unselected only).
				name := truncateStr(indent+"◆ "+displayN, nameColW)
				row = tableRowStyle.Render(cursor) +
					hclFileStyle.Width(nameColW).Render(name) +
					tableRowStyle.Render("  ") +
					tableRowStyle.Width(sizeColW).Render(size)
			} else {
				name := truncateStr(indent+displayN, nameColW)
				row = style.Render(cursor) +
					style.Width(nameColW).Render(name) +
					style.Render("  ") +
					style.Width(sizeColW).Render(size)
			}
			lines = append(lines, m.pad(row, style))
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
