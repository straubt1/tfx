// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package logger

import (
	"log/slog"
	"testing"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name      string
		logLevel  string
		logPath   string
		wantLevel slog.Level
	}{
		{
			name:      "init with debug level",
			logLevel:  "debug",
			logPath:   "",
			wantLevel: slog.LevelDebug,
		},
		{
			name:      "init with info level",
			logLevel:  "info",
			logPath:   "",
			wantLevel: slog.LevelInfo,
		},
		{
			name:      "init with trace level",
			logLevel:  "trace",
			logPath:   "",
			wantLevel: LevelTrace,
		},
		{
			name:      "init with log path",
			logLevel:  "info",
			logPath:   "/tmp/logs",
			wantLevel: slog.LevelInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init(tt.logLevel, tt.logPath)

			if GetLevel() != tt.wantLevel {
				t.Errorf("GetLevel() = %v, want %v", GetLevel(), tt.wantLevel)
			}

			if tt.logPath != "" && GetLogPath() != tt.logPath {
				t.Errorf("GetLogPath() = %q, want %q", GetLogPath(), tt.logPath)
			}
		})
	}
}

func TestIsEnabled(t *testing.T) {
	// Initialize logger with a known level
	Init("info", "")

	tests := []struct {
		name  string
		level slog.Level
		want  bool
	}{
		{"trace not enabled at info", LevelTrace, false},
		{"debug not enabled at info", slog.LevelDebug, false},
		{"info enabled at info", slog.LevelInfo, true},
		{"warn enabled at info", slog.LevelWarn, true},
		{"error enabled at info", slog.LevelError, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsEnabled(tt.level)
			if got != tt.want {
				t.Errorf("IsEnabled(%v) = %v, want %v", tt.level, got, tt.want)
			}
		})
	}
}

func TestLogFunctions(t *testing.T) {
	// Initialize logger
	Init("trace", "")

	// Test that logging functions don't panic
	t.Run("trace logging", func(t *testing.T) {
		Trace("test trace", "key", "value")
	})

	t.Run("debug logging", func(t *testing.T) {
		Debug("test debug", "key", "value")
	})

	t.Run("info logging", func(t *testing.T) {
		Info("test info", "key", "value")
	})

	t.Run("warn logging", func(t *testing.T) {
		Warn("test warn", "key", "value")
	})

	t.Run("error logging", func(t *testing.T) {
		Error("test error", "key", "value")
	})
}
