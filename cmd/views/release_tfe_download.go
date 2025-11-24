// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	"fmt"
	"io"
	"strings"

	"code.cloudfoundry.org/bytefmt"
)

// ReleaseTfeDownloadView handles rendering for release tfe download command
type ReleaseTfeDownloadView struct {
	*BaseView
}

func NewReleaseTfeDownloadView() *ReleaseTfeDownloadView {
	return &ReleaseTfeDownloadView{
		BaseView: NewBaseView(),
	}
}

// ProgressWriter wraps an io.Writer to track progress and update a spinner
type ProgressWriter struct {
	writer       io.Writer
	totalBytes   int64
	writtenBytes int64
	view         *ReleaseTfeDownloadView
}

// NewProgressWriter creates a progress writer that updates the view's spinner
func (v *ReleaseTfeDownloadView) NewProgressWriter(w io.Writer, totalBytes int64) *ProgressWriter {
	pw := &ProgressWriter{
		writer:     w,
		totalBytes: totalBytes,
		view:       v,
	}

	// Initialize spinner message
	if v.Output().Spinner() != nil {
		v.Output().Spinner().UpdateMessage("Downloading...")
	}

	return pw
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	n, err := pw.writer.Write(p)
	if n > 0 {
		pw.writtenBytes += int64(n)
		spinner := pw.view.Output().Spinner()
		if spinner != nil {
			written := bytefmt.ByteSize(uint64(pw.writtenBytes))

			if pw.totalBytes > 0 {
				// Show progress bar with percentage when total size is known
				percentage := float64(pw.writtenBytes) / float64(pw.totalBytes) * 100
				total := bytefmt.ByteSize(uint64(pw.totalBytes))

				// Create visual progress bar with bounds checking
				barWidth := 30
				filled := int(float64(barWidth) * percentage / 100)

				// Ensure filled is within valid bounds
				if filled < 0 {
					filled = 0
				} else if filled > barWidth {
					filled = barWidth
				}

				bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
				spinner.UpdateMessage(fmt.Sprintf("Downloading [%s] %.1f%% (%s / %s)", bar, percentage, written, total))
			} else {
				// Show indeterminate progress (just bytes written, no percentage)
				spinner.UpdateMessage(fmt.Sprintf("Downloading... %s written", written))
			}
		}
	}
	return n, err
}

// WrittenBytes returns the total bytes written
func (pw *ProgressWriter) WrittenBytes() int64 {
	return pw.writtenBytes
}

// Render renders the download completion message
func (v *ReleaseTfeDownloadView) Render(outputPath string, bytesWritten int64) error {
	// // Reset spinner message
	// if v.Output().Spinner() != nil {
	// 	v.Output().Spinner().UpdateMessage("TFx is working...")
	// }

	if v.IsJSON() {
		return v.Output().RenderJSON(map[string]interface{}{
			"status": "Success",
			"file":   outputPath,
			"size":   bytesWritten,
		})
	}

	v.Output().Message("Download complete: %s (%s)", outputPath, bytefmt.ByteSize(uint64(bytesWritten)))
	return nil
}
