// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"strings"
	"unicode"
)

// isHCLFile reports whether a file name warrants HCL syntax highlighting.
func isHCLFile(name string) bool {
	lower := strings.ToLower(name)
	return strings.HasSuffix(lower, ".tf") ||
		strings.HasSuffix(lower, ".tfvars") ||
		strings.HasSuffix(lower, ".hcl")
}

// hclBlockKeywords is the set of HCL/Terraform top-level block type keywords
// that receive special purple highlighting.
var hclBlockKeywords = map[string]bool{
	"resource":   true,
	"data":       true,
	"variable":   true,
	"output":     true,
	"locals":     true,
	"module":     true,
	"provider":   true,
	"terraform":  true,
	"backend":    true,
	"required_providers": true,
}

// colorizeHCLLine applies lightweight syntax highlighting to a single line of
// HCL/Terraform source.  It reuses the JSON lipgloss styles where the semantic
// meaning overlaps:
//
//	hclBlockKeyStyle  (purple)  — block type keywords (resource, variable, …)
//	jsonKeyStyle      (blue)    — attribute keys  (key = …)
//	jsonStringStyle   (green)   — quoted strings
//	jsonNumberStyle   (purple)  — numbers
//	jsonKeywordStyle  (amber)   — true / false / null
//	jsonPunctStyle    (dim)     — = { } [ ]
//	contentStyle      (fg)      — everything else
func colorizeHCLLine(line string) string {
	if line == "" {
		return contentStyle.Render(line)
	}

	// Comments: # or // at the start (after optional whitespace)
	trimmed := strings.TrimLeft(line, " \t")
	if strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "//") {
		return jsonPunctStyle.Render(line)
	}
	// Block-style comment continuation
	if strings.HasPrefix(trimmed, "*") || strings.HasPrefix(trimmed, "/*") || strings.HasPrefix(trimmed, "*/") {
		return jsonPunctStyle.Render(line)
	}

	// Attempt token-level coloring.
	return colorizeHCLTokens(line)
}

// colorizeHCLTokens performs a simple left-to-right tokenization of an HCL line.
func colorizeHCLTokens(line string) string {
	var out strings.Builder
	i := 0
	n := len(line)

	for i < n {
		ch := line[i]

		// Leading / inline whitespace — pass through unstyled.
		if ch == ' ' || ch == '\t' {
			out.WriteString(contentStyle.Render(string(ch)))
			i++
			continue
		}

		// Quoted strings (double-quote only; HCL uses " not ')
		if ch == '"' {
			end := hclStringEnd(line, i+1)
			segment := line[i:end]
			out.WriteString(jsonStringStyle.Render(segment))
			i = end
			continue
		}

		// Heredoc indicator (<<) — treat rest of line as plain.
		if i+1 < n && ch == '<' && line[i+1] == '<' {
			out.WriteString(jsonPunctStyle.Render(line[i:]))
			break
		}

		// Punctuation: = { } [ ] ( ) , :
		if ch == '=' || ch == '{' || ch == '}' || ch == '[' || ch == ']' || ch == '(' || ch == ')' || ch == ',' || ch == ':' {
			out.WriteString(jsonPunctStyle.Render(string(ch)))
			i++
			continue
		}

		// Numbers (leading digit or minus-digit)
		if unicode.IsDigit(rune(ch)) || (ch == '-' && i+1 < n && unicode.IsDigit(rune(line[i+1]))) {
			j := i + 1
			for j < n && (unicode.IsDigit(rune(line[j])) || line[j] == '.' || line[j] == 'e' || line[j] == 'E' || line[j] == '+' || line[j] == '-') {
				j++
			}
			out.WriteString(jsonNumberStyle.Render(line[i:j]))
			i = j
			continue
		}

		// Identifier: keyword, attribute name, reference, etc.
		if isHCLIdentStart(ch) {
			j := i + 1
			for j < n && isHCLIdentChar(line[j]) {
				j++
			}
			word := line[i:j]

			// Peek ahead (skip spaces) to decide context.
			k := j
			for k < n && (line[k] == ' ' || line[k] == '\t') {
				k++
			}

			switch {
			case hclBlockKeywords[word]:
				// Block keyword at start of significant content (allow indentation).
				out.WriteString(hclBlockKeyStyle.Render(word))
			case word == "true" || word == "false" || word == "null":
				out.WriteString(jsonKeywordStyle.Render(word))
			case k < n && line[k] == '=':
				// attribute key (followed by '=')
				out.WriteString(jsonKeyStyle.Render(word))
			default:
				out.WriteString(contentStyle.Render(word))
			}
			i = j
			continue
		}

		// Anything else — plain foreground.
		out.WriteString(contentStyle.Render(string(ch)))
		i++
	}
	return out.String()
}

// hclStringEnd returns the index just past the closing '"' starting from pos.
// Handles simple backslash escapes and interpolation sequences ${…} (treated
// as opaque — the whole string is highlighted as a string).
func hclStringEnd(s string, pos int) int {
	for pos < len(s) {
		ch := s[pos]
		if ch == '\\' {
			pos += 2 // skip escaped char
			continue
		}
		if ch == '"' {
			return pos + 1
		}
		pos++
	}
	return len(s)
}

// isHCLIdentStart returns true for the first character of an HCL identifier.
func isHCLIdentStart(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

// isHCLIdentChar returns true for subsequent characters of an HCL identifier.
func isHCLIdentChar(ch byte) bool {
	return isHCLIdentStart(ch) || (ch >= '0' && ch <= '9') || ch == '-' || ch == '.'
}
