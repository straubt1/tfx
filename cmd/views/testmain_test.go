// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

// View tests run with JSON output mode so assertions can parse structured output.
// The output singleton (output.Get()) is initialized once per test binary via sync.Once;
// setting viper "json" = true before m.Run() ensures it initializes in JSON mode.

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	viper.Set("json", true)
	os.Exit(m.Run())
}

// captureOutput redirects os.Stdout, calls fn, then returns everything written.
func captureOutput(t *testing.T, fn func() error) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	old := os.Stdout
	os.Stdout = w

	fnErr := fn()

	w.Close()
	os.Stdout = old

	if fnErr != nil {
		t.Fatalf("view render returned unexpected error: %v", fnErr)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("failed to read captured output: %v", err)
	}
	return buf.String()
}
