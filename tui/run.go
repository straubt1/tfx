// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/output"
)

// Run launches the TUI, reading connection settings from Viper (same source as CLI commands).
func Run() error {
	// Kill the CLI spinner before handing the terminal to Bubble Tea.
	// Without this, the spinner's goroutine writes "TFx is working..." to stdout
	// while Bubble Tea is rendering, corrupting the alt-screen display.
	output.Get().DisableSpinner()

	// Create an event bus for the API Inspector panel, then build a client that
	// always installs the logging transport so the TUI receives every HTTP call.
	bus := client.NewAPIEventBus()
	c, err := client.NewFromViperForTUI(bus)
	if err != nil {
		return errors.Wrap(err, "failed to create TFx client")
	}

	profileName := viper.GetString("profile")
	m := newModel(c, profileName)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		return errors.Wrap(err, "tui error")
	}
	return nil
}
