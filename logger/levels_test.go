// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package logger

import (
	"log/slog"
	"testing"
)

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  slog.Level
	}{
		{"trace", "trace", LevelTrace},
		{"trace uppercase", "TRACE", LevelTrace},
		{"trace with spaces", "  trace  ", LevelTrace},
		{"debug", "debug", slog.LevelDebug},
		{"debug uppercase", "DEBUG", slog.LevelDebug},
		{"info", "info", slog.LevelInfo},
		{"info uppercase", "INFO", slog.LevelInfo},
		{"warn", "warn", slog.LevelWarn},
		{"warning", "warning", slog.LevelWarn},
		{"error", "error", slog.LevelError},
		{"error uppercase", "ERROR", slog.LevelError},
		{"off", "off", LevelOff},
		{"off uppercase", "OFF", LevelOff},
		{"empty string", "", LevelOff},
		{"invalid", "invalid", LevelOff},
		{"unknown", "unknown123", LevelOff},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseLevel(tt.input)
			if got != tt.want {
				t.Errorf("ParseLevel(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestLevelString(t *testing.T) {
	tests := []struct {
		name  string
		level slog.Level
		want  string
	}{
		{"trace level", LevelTrace, "TRACE"},
		{"debug level", slog.LevelDebug, "DEBUG"},
		{"info level", slog.LevelInfo, "INFO"},
		{"warn level", slog.LevelWarn, "WARN"},
		{"error level", slog.LevelError, "ERROR"},
		{"off level", LevelOff, "OFF"},
		{"higher than off", slog.Level(200), "OFF"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LevelString(tt.level)
			if got != tt.want {
				t.Errorf("LevelString(%v) = %q, want %q", tt.level, got, tt.want)
			}
		})
	}
}
