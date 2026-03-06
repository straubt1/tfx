// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import "charm.land/lipgloss/v2"

// Color palette — GitHub Dark inspired, works well on dark terminals.
var (
	colorBg        = lipgloss.Color("#0D1117")
	colorFg        = lipgloss.Color("#E6EDF3")
	colorAccent    = lipgloss.Color("#58A6FF")
	colorDim       = lipgloss.Color("#8B949E")
	colorPurple    = lipgloss.Color("#BC8CFF")
	colorBorder    = lipgloss.Color("#30363D")
	colorHeaderBg  = lipgloss.Color("#161B22")
	colorAppBg     = lipgloss.Color("#1F6FEB")
	colorSelected  = lipgloss.Color("#1C2128") // slightly lighter than bg for selected row
	colorError     = lipgloss.Color("#F85149")
	colorLoading   = lipgloss.Color("#D29922")
	colorSuccess   = lipgloss.Color("#3FB950") // GitHub dark green
)

var (
	// Header bar
	headerStyle = lipgloss.NewStyle().
			Background(colorHeaderBg).
			Foreground(colorFg)

	headerAppStyle = lipgloss.NewStyle().
			Background(colorAppBg).
			Foreground(colorFg).
			Bold(true).
			Padding(0, 1)

	headerInfoStyle = lipgloss.NewStyle().
			Background(colorHeaderBg).
			Foreground(colorDim).
			Padding(0, 1)

	headerVersionStyle = lipgloss.NewStyle().
				Background(colorHeaderBg).
				Foreground(colorPurple).
				Padding(0, 1)

	// Breadcrumb bar
	breadcrumbBarStyle = lipgloss.NewStyle().
				Background(colorBg).
				Foreground(colorDim)

	breadcrumbActiveStyle = lipgloss.NewStyle().
				Background(colorBg).
				Foreground(colorAccent).
				Bold(true)

	breadcrumbSepStyle = lipgloss.NewStyle().
				Background(colorBg).
				Foreground(colorBorder)

	// Content area
	contentStyle = lipgloss.NewStyle().
			Background(colorBg).
			Foreground(colorFg)

	contentTitleStyle = lipgloss.NewStyle().
				Background(colorBg).
				Foreground(colorFg).
				Bold(true)

	contentDividerStyle = lipgloss.NewStyle().
				Background(colorBg).
				Foreground(colorBorder)

	contentPlaceholderStyle = lipgloss.NewStyle().
				Background(colorBg).
				Foreground(colorDim).
				Italic(true)

	// Status bar
	statusBarStyle = lipgloss.NewStyle().
			Background(colorHeaderBg).
			Foreground(colorDim)

	// CLI hint bar
	cliHintBarStyle = lipgloss.NewStyle().
			Background(colorBg).
			Foreground(colorDim)

	cliHintCmdStyle = lipgloss.NewStyle().
			Background(colorBg).
			Foreground(colorAccent).
			Italic(true)

	// Help overlay
	helpTitleStyle = lipgloss.NewStyle().
			Background(colorHeaderBg).
			Foreground(colorFg).
			Bold(true)

	helpKeyStyle = lipgloss.NewStyle().
			Background(colorHeaderBg).
			Foreground(colorAccent).
			Bold(true)

	helpDescStyle = lipgloss.NewStyle().
			Background(colorHeaderBg).
			Foreground(colorFg)

	helpBarStyle = lipgloss.NewStyle().
			Background(colorHeaderBg).
			Foreground(colorDim)

	// Status bar variants
	statusLoadingStyle = lipgloss.NewStyle().
				Background(colorHeaderBg).
				Foreground(colorLoading)

	statusErrorStyle = lipgloss.NewStyle().
			Background(colorHeaderBg).
			Foreground(colorError)

	statusSuccessStyle = lipgloss.NewStyle().
				Background(colorHeaderBg).
				Foreground(colorSuccess)

	// Table
	tableHeaderStyle = lipgloss.NewStyle().
			Background(colorHeaderBg).
			Foreground(colorDim).
			Bold(true)

	tableRowStyle = lipgloss.NewStyle().
			Background(colorBg).
			Foreground(colorFg)

	tableRowSelectedStyle = lipgloss.NewStyle().
				Background(colorSelected).
				Foreground(colorAccent).
				Bold(true)

	// Filter bar
	filterBarStyle = lipgloss.NewStyle().
			Background(colorHeaderBg).
			Foreground(colorDim)

	filterBarActiveStyle = lipgloss.NewStyle().
				Background(colorHeaderBg).
				Foreground(colorFg)

	// Workspace detail view
	detailLabelStyle = lipgloss.NewStyle().
			Background(colorBg).
			Foreground(colorDim)

	// Workspace tab strip
	tabBarStyle = lipgloss.NewStyle().
			Background(colorHeaderBg).
			Foreground(colorDim)

	tabActiveStyle = lipgloss.NewStyle().
			Background(colorHeaderBg).
			Foreground(colorAccent).
			Bold(true).
			Underline(true).
			Padding(0, 1)

	tabInactiveStyle = lipgloss.NewStyle().
				Background(colorHeaderBg).
				Foreground(colorDim).
				Padding(0, 1)

	// JSON viewer syntax highlighting
	jsonKeyStyle     = lipgloss.NewStyle().Background(colorBg).Foreground(colorAccent)  // keys → blue
	jsonStringStyle  = lipgloss.NewStyle().Background(colorBg).Foreground(colorSuccess) // strings → green
	jsonNumberStyle  = lipgloss.NewStyle().Background(colorBg).Foreground(colorPurple)  // numbers → purple
	jsonKeywordStyle = lipgloss.NewStyle().Background(colorBg).Foreground(colorLoading) // true/false/null → amber
	jsonPunctStyle   = lipgloss.NewStyle().Background(colorBg).Foreground(colorDim)     // braces/colons/commas → dim

	// HCL/Terraform viewer syntax highlighting
	hclBlockKeyStyle = lipgloss.NewStyle().Background(colorBg).Foreground(colorPurple) // block type keywords → purple
)
