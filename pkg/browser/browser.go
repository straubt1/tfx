// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

// Package browser opens URLs in the system's default web browser.
package browser

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Open launches url in the system's default browser.
// It uses Start (not Run) so execution returns immediately without waiting
// for the browser process to exit.
func Open(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return fmt.Errorf("opening browser not supported on %s", runtime.GOOS)
	}
	return cmd.Start()
}
