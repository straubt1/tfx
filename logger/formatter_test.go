// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package logger

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestNewColorHandler(t *testing.T) {
	t.Run("creates handler with options", func(t *testing.T) {
		var buf bytes.Buffer
		opts := &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
		handler := NewColorHandler(&buf, opts)

		if handler == nil {
			t.Fatal("NewColorHandler() returned nil")
		}
	})

	t.Run("creates handler with nil options", func(t *testing.T) {
		var buf bytes.Buffer
		handler := NewColorHandler(&buf, nil)

		if handler == nil {
			t.Fatal("NewColorHandler() with nil options returned nil")
		}
	})
}

func TestColorHandlerEnabled(t *testing.T) {
	tests := []struct {
		name      string
		minLevel  slog.Level
		testLevel slog.Level
		want      bool
	}{
		{"debug enabled when min is debug", slog.LevelDebug, slog.LevelDebug, true},
		{"info enabled when min is debug", slog.LevelDebug, slog.LevelInfo, true},
		{"debug disabled when min is info", slog.LevelInfo, slog.LevelDebug, false},
		{"info enabled when min is info", slog.LevelInfo, slog.LevelInfo, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			opts := &slog.HandlerOptions{
				Level: tt.minLevel,
			}
			handler := NewColorHandler(&buf, opts)

			got := handler.Enabled(context.Background(), tt.testLevel)
			if got != tt.want {
				t.Errorf("Enabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColorHandlerHandle(t *testing.T) {
	tests := []struct {
		name    string
		level   slog.Level
		message string
		attrs   []slog.Attr
	}{
		{
			name:    "trace level",
			level:   LevelTrace,
			message: "trace message",
			attrs:   []slog.Attr{slog.String("key", "value")},
		},
		{
			name:    "debug level",
			level:   slog.LevelDebug,
			message: "debug message",
			attrs:   []slog.Attr{slog.Int("count", 42)},
		},
		{
			name:    "info level",
			level:   slog.LevelInfo,
			message: "info message",
			attrs:   []slog.Attr{},
		},
		{
			name:    "warn level",
			level:   slog.LevelWarn,
			message: "warn message",
			attrs:   []slog.Attr{slog.Bool("flag", true)},
		},
		{
			name:    "error level",
			level:   slog.LevelError,
			message: "error message",
			attrs:   []slog.Attr{slog.String("error", "test error")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			opts := &slog.HandlerOptions{
				Level: LevelTrace,
			}
			handler := NewColorHandler(&buf, opts)

			record := slog.NewRecord(time.Now(), tt.level, tt.message, 0)
			for _, attr := range tt.attrs {
				record.AddAttrs(attr)
			}

			err := handler.Handle(context.Background(), record)
			if err != nil {
				t.Errorf("Handle() error = %v", err)
			}

			output := buf.String()
			if !strings.Contains(output, tt.message) {
				t.Errorf("Output doesn't contain message %q: %q", tt.message, output)
			}
		})
	}
}
