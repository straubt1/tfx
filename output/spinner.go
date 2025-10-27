// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package output

import (
	"sync"
	"time"

	"github.com/briandowns/spinner"
)

// Spinner manages the loading spinner
type Spinner struct {
	s       *spinner.Spinner
	running bool
	mu      sync.Mutex
	depth   int // Track nested calls to prevent flicker
}

// NewSpinner creates a new spinner instance
func NewSpinner() *Spinner {
	sp := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	sp.Color("cyan")
	sp.Suffix = "  TFx is working..."
	sp.Start()

	return &Spinner{
		s:       sp,
		running: true,
	}
}

// Stop stops the spinner
// Uses depth tracking to handle nested stop/start calls
func (s *Spinner) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.depth++
	if s.depth == 1 && s.s != nil && s.running {
		s.s.Stop()
		s.running = false
	}
}

// FinalStop stops the spinner completely, ignoring depth tracking
// This should be called when you want to permanently stop the spinner
func (s *Spinner) FinalStop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.s != nil && s.running {
		s.s.Stop()
		s.running = false
		s.depth = 0
	}
}

// Start starts the spinner
// Uses depth tracking to handle nested stop/start calls
func (s *Spinner) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.depth > 0 {
		s.depth--
	}
	if s.depth == 0 && s.s != nil && !s.running {
		s.s.Start()
		s.running = true
	}
}

// UpdateMessage updates the spinner message
func (s *Spinner) UpdateMessage(msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.s != nil {
		s.s.Suffix = "  " + msg
	}
}

// IsRunning returns true if the spinner is currently running
func (s *Spinner) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}
