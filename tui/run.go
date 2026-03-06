// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Run launches the TUI, reading connection settings from Viper (same source as CLI commands).
func Run() error {
	m := newModel(
		viper.GetString("tfeHostname"),
		viper.GetString("tfeOrganization"),
	)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		return errors.Wrap(err, "tui error")
	}
	return nil
}
