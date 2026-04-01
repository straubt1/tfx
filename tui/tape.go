// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
)

// TapeRecorder writes a VHS-compatible .tape file as the user interacts with the TUI.
// Every keypress is recorded with the elapsed time since the previous keypress as a Sleep command.
type TapeRecorder struct {
	path     string
	file     *os.File
	lastTime time.Time
	typeBuf  strings.Builder
	started  bool // true after the header has been written
}

// NewTapeRecorder creates a TapeRecorder that writes to path.
// The file is created (or truncated) immediately; the tape header is written
// on the first call to WriteHeader (triggered by the first tea.WindowSizeMsg).
func NewTapeRecorder(path string) (*TapeRecorder, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return &TapeRecorder{
		path: path,
		file: f,
	}, nil
}

// WriteHeader writes the VHS tape header to the file. It is idempotent — only the
// first call has effect. Call this once the terminal size is known (from WindowSizeMsg).
func (r *TapeRecorder) WriteHeader() {
	if r.started {
		return
	}
	r.started = true

	// Derive Output path: replace .tape extension with .gif, same directory.
	outputPath := strings.TrimSuffix(r.path, filepath.Ext(r.path)) + ".gif"

	fmt.Fprintf(r.file, "Output %s\n\n", outputPath)
	fmt.Fprintf(r.file, "Set Shell        \"zsh\"\n")
	fmt.Fprintf(r.file, "Set Width        1200\n")
	fmt.Fprintf(r.file, "Set Height       600\n")
	fmt.Fprintf(r.file, "Set FontSize     14\n")
	fmt.Fprintf(r.file, "Set FontFamily   \"RobotoMono Nerd Font\"\n")
	fmt.Fprintf(r.file, "Set Theme        \"GitHub Dark\"\n")
	fmt.Fprintf(r.file, "Set Padding      20\n")
	fmt.Fprintf(r.file, "Set WindowBar    \"Colorful\"\n")
	fmt.Fprintf(r.file, "Set BorderRadius 8\n")
	fmt.Fprintf(r.file, "Set TypingSpeed  100ms\n")
	fmt.Fprintf(r.file, "Set Framerate    60\n")
	fmt.Fprintf(r.file, "Set WaitTimeout  30s\n")
	fmt.Fprintf(r.file, "Sleep 1s\n")
	fmt.Fprintf(r.file, "Type \"./tfx tui\"\n")
	fmt.Fprintf(r.file, "Enter\n")
	fmt.Fprintf(r.file, "Sleep 2s\n")
	fmt.Fprintf(r.file, "\n# ============================================================\n")
	fmt.Fprintf(r.file, "# RECORDING STARTS HERE\n")
	fmt.Fprintf(r.file, "# ============================================================\n\n")

	r.lastTime = time.Now()
}

// Record writes the VHS command for a keypress, preceded by a Sleep that reflects
// how long has passed since the previous keypress.
func (r *TapeRecorder) Record(msg tea.KeyPressMsg) {
	if !r.started {
		return
	}

	now := time.Now()
	elapsed := now.Sub(r.lastTime)
	r.lastTime = now

	// Round to nearest 10ms and emit a Sleep if non-zero.
	ms := int(math.Round(float64(elapsed.Milliseconds())/10) * 10)
	if ms > 0 {
		fmt.Fprintf(r.file, "Sleep %dms\n", ms)
	}

	// Printable character — buffer consecutive chars so they coalesce into
	// a single Type "..." command (e.g. for filter input).
	if isPrintable(msg.String()) {
		r.typeBuf.WriteString(msg.String())
		return
	}

	// Non-printable key — flush any buffered printable chars first.
	r.flushType()

	fmt.Fprintf(r.file, "%s\n", keyToVHS(msg.String()))
}

// Flush writes any remaining buffered Type text and closes the file.
// Call this after the Bubble Tea program exits.
func (r *TapeRecorder) Flush() {
	if r.file == nil {
		return
	}
	r.flushType()
	r.file.Close()
	r.file = nil
}

// flushType writes the buffered printable characters as a single Type command.
func (r *TapeRecorder) flushType() {
	if r.typeBuf.Len() == 0 {
		return
	}
	fmt.Fprintf(r.file, "Type %q\n", r.typeBuf.String())
	r.typeBuf.Reset()
}

// keyToVHS maps a Bubble Tea v2 key string to the equivalent VHS tape command.
// Unknown keys are written as comments so the tape remains valid VHS syntax.
func keyToVHS(key string) string {
	switch key {
	case "enter":
		return "Enter"
	case "esc":
		return "Escape"
	case "backspace":
		return "Backspace"
	case "delete":
		return "Delete"
	case "tab":
		return "Tab"
	case "space":
		return "Space"
	case "up":
		return "Up"
	case "down":
		return "Down"
	case "left":
		return "Left"
	case "right":
		return "Right"
	case "pgup", "page_up":
		return "PageUp"
	case "pgdown", "page_down":
		return "PageDown"
	case "ctrl+c":
		return "Ctrl+C"
	case "ctrl+d":
		return "Ctrl+D"
	case "ctrl+l":
		return "Ctrl+L"
	case "ctrl+r":
		return "Ctrl+R"
	case "ctrl+z":
		return "Ctrl+Z"
	case "ctrl+a":
		return "Ctrl+A"
	case "ctrl+e":
		return "Ctrl+E"
	case "ctrl+u":
		return "Ctrl+U"
	case "ctrl+k":
		return "Ctrl+K"
	case "ctrl+w":
		return "Ctrl+W"
	case "shift+tab":
		return "Shift+Tab"
	case "shift+enter":
		return "Shift+Enter"
	}

	// alt+x → Alt+X
	if strings.HasPrefix(key, "alt+") {
		rest := key[4:]
		if len(rest) == 1 {
			return "Alt+" + strings.ToUpper(rest)
		}
		return "Alt+" + rest
	}

	// ctrl+x → Ctrl+X (catch-all for unlisted ctrl combos)
	if strings.HasPrefix(key, "ctrl+") {
		rest := key[5:]
		return "Ctrl+" + strings.ToUpper(rest)
	}

	// Unknown — write as a comment so the tape file stays valid VHS
	return fmt.Sprintf("# unknown: %s", key)
}
