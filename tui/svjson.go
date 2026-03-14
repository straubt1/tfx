// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"fmt"
	"strings"
)

// renderStateVersionJsonContent renders the scrollable, syntax-highlighted JSON
// viewer for the selected state version.
func (m Model) renderStateVersionJsonContent() string {
	h := m.contentHeight()

	// ── Loading ───────────────────────────────────────────────────────────────
	if m.svJsonLoading {
		lines := make([]string, h)
		mid := h / 2
		frame := spinnerFrames[m.spinnerIdx]
		for i := range lines {
			if i == mid {
				lines[i] = contentPlaceholderStyle.Width(m.innerWidth()).Render("  " + frame + "  Loading state JSON…")
			} else {
				lines[i] = contentStyle.Width(m.innerWidth()).Render("")
			}
		}
		return strings.Join(lines, "\n")
	}

	// ── Error ─────────────────────────────────────────────────────────────────
	if m.svJsonErr != "" {
		lines := make([]string, h)
		for i := range lines {
			if i == 0 {
				lines[i] = contentStyle.Width(m.innerWidth()).Render("  ✗  " + m.svJsonErr)
			} else {
				lines[i] = contentStyle.Width(m.innerWidth()).Render("")
			}
		}
		return strings.Join(lines, "\n")
	}

	// ── Empty ─────────────────────────────────────────────────────────────────
	if len(m.svJsonLines) == 0 {
		lines := make([]string, h)
		for i := range lines {
			lines[i] = contentStyle.Width(m.innerWidth()).Render("")
		}
		return strings.Join(lines, "\n")
	}

	// ── Syntax-highlighted line-numbered JSON ─────────────────────────────────
	numLines := len(m.svJsonLines)
	lineNumWidth := len(fmt.Sprintf("%d", numLines))
	// Layout: 2 margin + lineNumWidth + " │ " (3) + content
	contentWidth := m.innerWidth() - 2 - lineNumWidth - 3
	if contentWidth < 10 {
		contentWidth = 10
	}

	all := make([]string, 0, numLines)
	for i, line := range m.svJsonLines {
		lineNum := fmt.Sprintf("%*d", lineNumWidth, i+1)
		display := line
		if len(display) > contentWidth {
			display = display[:contentWidth-1] + "…"
		}
		row := "  " +
			detailLabelStyle.Render(lineNum) +
			contentDividerStyle.Render(" │ ") +
			colorizeJSONLine(display)
		all = append(all, contentStyle.Width(m.innerWidth()).Render(row))
	}

	// Clamp scroll and slice visible window.
	maxScroll := len(all) - h
	if maxScroll < 0 {
		maxScroll = 0
	}
	start := m.svJsonScroll
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

// ── JSON syntax colorizer ─────────────────────────────────────────────────────
//
// Works character-by-character on a single pretty-printed JSON line.
// No external dependencies — handles the regular structure produced by
// json.Indent("", "  ") which this viewer always uses.
//
// Color mapping (GitHub Dark palette):
//   keys          → colorAccent  (blue)
//   string values → colorSuccess (green)
//   numbers       → colorPurple  (purple)
//   true/false/null → colorLoading (amber)
//   punctuation   → colorDim     (dim gray)

// colorizeJSONLine applies syntax colors to one line of pretty-printed JSON.
func colorizeJSONLine(line string) string {
	// Preserve leading whitespace unstyled (indentation).
	trimmed := strings.TrimLeft(line, " \t")
	if trimmed == "" {
		return line
	}
	indent := line[:len(line)-len(trimmed)]
	return indent + tokenizeJSON(trimmed)
}

// tokenizeJSON colorizes the non-whitespace portion of a JSON line.
func tokenizeJSON(s string) string {
	if s == "" {
		return s
	}
	var out strings.Builder
	switch s[0] {
	case '"':
		// Either a key ("key": ...) or a standalone string value.
		end := jsonStringEnd(s, 0)
		if end < 0 {
			// Truncated string (line was cut) — leave unstyled.
			return s
		}
		str := s[:end]
		rest := s[end:]
		restTrim := strings.TrimLeft(rest, " ")
		if len(restTrim) > 0 && restTrim[0] == ':' {
			// It's a key: colour it blue.
			out.WriteString(jsonKeyStyle.Render(str))
			ws := rest[:len(rest)-len(restTrim)]
			out.WriteString(jsonPunctStyle.Render(ws + ":"))
			out.WriteString(jsonColorValue(restTrim[1:])) // everything after ':'
		} else {
			// Standalone string value (e.g., an array element).
			out.WriteString(jsonStringStyle.Render(str))
			if rest != "" {
				out.WriteString(jsonPunctStyle.Render(rest)) // trailing comma
			}
		}
	case '{', '}', '[', ']':
		// Structural bracket line (possibly with trailing comma).
		out.WriteString(jsonPunctStyle.Render(s))
	default:
		// Number, true, false, null — possibly with trailing comma.
		out.WriteString(jsonColorValue(s))
	}
	return out.String()
}

// jsonColorValue colorizes a JSON value fragment that may include an optional
// leading space and trailing comma, e.g.: ` "hello",` or ` 42,` or ` true`.
func jsonColorValue(s string) string {
	var out strings.Builder

	// Emit leading whitespace as-is (not styled).
	trimmed := strings.TrimLeft(s, " \t")
	if ws := s[:len(s)-len(trimmed)]; ws != "" {
		out.WriteString(ws)
	}
	s = trimmed
	if s == "" {
		return out.String()
	}

	// Peel off a trailing comma.
	trailer := ""
	if s[len(s)-1] == ',' {
		trailer = ","
		s = s[:len(s)-1]
	}
	if s == "" {
		out.WriteString(jsonPunctStyle.Render(trailer))
		return out.String()
	}

	switch s[0] {
	case '"':
		out.WriteString(jsonStringStyle.Render(s))
	case '{', '[', '}', ']':
		out.WriteString(jsonPunctStyle.Render(s))
	case 't', 'f', 'n': // true, false, null
		out.WriteString(jsonKeywordStyle.Render(s))
	default: // number (including negative numbers starting with '-')
		out.WriteString(jsonNumberStyle.Render(s))
	}
	if trailer != "" {
		out.WriteString(jsonPunctStyle.Render(trailer))
	}
	return out.String()
}

// jsonStringEnd returns the index immediately after the closing quote of a
// JSON string starting at pos. Returns -1 when the string is unclosed
// (i.e., the line was truncated mid-token).
func jsonStringEnd(s string, pos int) int {
	if pos >= len(s) || s[pos] != '"' {
		return -1
	}
	i := pos + 1
	for i < len(s) {
		switch s[i] {
		case '\\':
			i += 2 // skip the escaped character
		case '"':
			return i + 1
		default:
			i++
		}
	}
	return -1 // unclosed
}
