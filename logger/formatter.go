// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/logrusorgru/aurora"
)

// ColorHandler is a custom slog.Handler that formats logs with colors
type ColorHandler struct {
	opts slog.HandlerOptions
	out  io.Writer
}

// NewColorHandler creates a new ColorHandler
func NewColorHandler(out io.Writer, opts *slog.HandlerOptions) *ColorHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &ColorHandler{
		opts: *opts,
		out:  out,
	}
}

// Enabled reports whether the handler handles records at the given level
func (h *ColorHandler) Enabled(_ context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

// Handle formats and writes the log record
func (h *ColorHandler) Handle(_ context.Context, r slog.Record) error {
	timestamp := r.Time.Format("15:04:05")
	level := h.formatLevel(r.Level)

	// Build the message
	msg := fmt.Sprintf("%s %s %s",
		aurora.Faint(fmt.Sprintf("[%s]", timestamp)),
		level,
		r.Message,
	)

	// Add attributes (key-value pairs)
	r.Attrs(func(a slog.Attr) bool {
		msg += fmt.Sprintf(" %s=%s",
			aurora.Faint(a.Key),
			aurora.Faint(a.Value.String()),
		)
		return true
	})

	fmt.Fprintln(h.out, msg)
	return nil
}

// WithAttrs returns a new Handler with additional attributes
func (h *ColorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// For simplicity, we're not implementing attribute grouping
	// This would be needed for more advanced use cases
	return h
}

// WithGroup returns a new Handler with a group name
func (h *ColorHandler) WithGroup(name string) slog.Handler {
	// For simplicity, we're not implementing grouping
	// This would be needed for more advanced use cases
	return h
}

// formatLevel formats the log level with color
func (h *ColorHandler) formatLevel(level slog.Level) string {
	var levelStr string
	var colored string

	switch {
	case level < slog.LevelDebug:
		levelStr = "TRACE"
		colored = aurora.Faint(fmt.Sprintf("[%-5s]", levelStr)).String()
	case level < slog.LevelInfo:
		levelStr = "DEBUG"
		colored = aurora.Cyan(fmt.Sprintf("[%-5s]", levelStr)).String()
	case level < slog.LevelWarn:
		levelStr = "INFO"
		colored = aurora.Green(fmt.Sprintf("[%-5s]", levelStr)).String()
	case level < slog.LevelError:
		levelStr = "WARN"
		colored = aurora.Yellow(fmt.Sprintf("[%-5s]", levelStr)).String()
	default:
		levelStr = "ERROR"
		colored = aurora.Red(fmt.Sprintf("[%-5s]", levelStr)).String()
	}

	return colored
}
