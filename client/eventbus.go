// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package client

import "time"

// APIEvent captures a single TFE HTTP round-trip for the TUI API Inspector panel.
// It is published to an APIEventBus by the LoggingTransport after every request,
// regardless of whether TFX_LOG is set.
type APIEvent struct {
	Timestamp  time.Time
	Method     string        // HTTP method (GET, POST, PATCH, DELETE, …)
	URL        string        // Full URL — used for filter matching
	Path       string        // Path-only (scheme+host stripped) — used for display
	StatusCode int           // 0 if the round-trip errored before a response was received
	Duration   time.Duration // wall-clock time from request start to response end
	ReqHeaders []string      // sorted "Name: value" lines; sensitive values are [REDACTED]
	ReqBody    string        // Request body; empty for GET/HEAD/DELETE
	RespHeaders []string     // sorted "Name: value" lines from the HTTP response
	RespBody   string        // Response body, pretty-printed if valid JSON
	Err        string        // Non-empty when the round-trip returned an error
}

// APIEventBus is a goroutine-safe, non-blocking event sink.
// The LoggingTransport writes events; the TUI reads them via a blocking Bubble Tea Cmd.
type APIEventBus struct {
	ch chan APIEvent
}

// NewAPIEventBus returns a new bus with a 256-event buffer.
// A buffer this size prevents the HTTP transport from ever blocking even if the
// TUI is slow to drain (events are dropped when the buffer is full, not the
// HTTP request).
func NewAPIEventBus() *APIEventBus {
	return &APIEventBus{ch: make(chan APIEvent, 256)}
}

// Send publishes e to the bus. It never blocks: if the buffer is full the event
// is silently dropped rather than stalling the HTTP round-trip.
func (b *APIEventBus) Send(e APIEvent) {
	select {
	case b.ch <- e:
	default: // buffer full — drop rather than block the HTTP transport goroutine
	}
}

// Receive returns the read-only channel for use in a Bubble Tea Cmd:
//
//	func waitForAPIEvent(bus *client.APIEventBus) tea.Cmd {
//	    return func() tea.Msg { return <-bus.Receive() }
//	}
func (b *APIEventBus) Receive() <-chan APIEvent {
	return b.ch
}
