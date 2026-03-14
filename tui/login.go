// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	tfe "github.com/hashicorp/go-tfe"

	"github.com/straubt1/tfx/pkg/browser"
	"github.com/straubt1/tfx/pkg/hclconfig"
)

// loginStep is the state machine for the login wizard.
type loginStep int

const (
	stepProfileList     loginStep = iota // existing-profile selector (skipped when no profiles)
	stepMenu                             // two options: browser / direct entry
	stepProfileName                      // two options: use hostname / enter custom name
	stepProfileNameEntry                 // text input for a custom profile name
	stepToken                            // masked token input
	stepValidating                       // spinner while fetching orgs
	stepTokenError                       // validation failed — re-enter or accept anyway
	stepOrgSelect                        // arrow-key org picker (2+ orgs)
	stepDone                             // success — profile written
	stepError                            // fatal write/config error
	stepCancelled                        // clean exit via q/esc/ctrl+c
)

// ── Message types ─────────────────────────────────────────────────────────────

type loginOrgsMsg    []*tfe.Organization
type loginErrMsg     struct{ err error }
type loginSpinnerMsg struct{}

// ── Commands ──────────────────────────────────────────────────────────────────

var loginSpinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

func loginTickSpinner() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(80 * time.Millisecond)
		return loginSpinnerMsg{}
	}
}

func loginFetchOrgs(hostname, token string) tea.Cmd {
	return func() tea.Msg {
		tfeClient, err := tfe.NewClient(&tfe.Config{
			Address: "https://" + hostname,
			Token:   token,
		})
		if err != nil {
			return loginErrMsg{err}
		}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		var all []*tfe.Organization
		for pageNum := 1; ; pageNum++ {
			result, err := tfeClient.Organizations.List(ctx, &tfe.OrganizationListOptions{
				ListOptions: tfe.ListOptions{PageNumber: pageNum, PageSize: 100},
			})
			if err != nil {
				return loginErrMsg{err}
			}
			all = append(all, result.Items...)
			if pageNum >= result.TotalPages {
				break
			}
		}
		if len(all) == 0 {
			return loginErrMsg{fmt.Errorf("token is valid but no organizations are accessible — check token permissions")}
		}
		return loginOrgsMsg(all)
	}
}

// ── Model ─────────────────────────────────────────────────────────────────────

// LoginModel is a self-contained Bubble Tea model for the inline login wizard.
type LoginModel struct {
	step                loginStep
	hostname            string              // target hostname for the create-new flow
	configPath          string
	profiles            []hclconfig.Profile // existing profiles loaded from config file
	profileCursor       int
	isUpdate            bool                // true when re-authing an existing profile
	selectedProfileName string              // final name written to the profile block
	profileNameCursor   int                 // stepProfileName: 0 = use hostname, 1 = enter custom
	nameRunes           []rune              // stepProfileNameEntry text buffer
	menuCursor          int                 // stepMenu: 0 = browser, 1 = direct
	useBrowser          bool
	tokenRunes          []rune
	tokenErr            error               // validation error shown in stepTokenError
	tokenErrCursor      int                 // 0 = re-enter, 1 = accept anyway
	orgs                []*tfe.Organization
	orgCursor           int
	selectedOrg         string
	resolvedToken       string
	spinnerIdx          int
	err                 error               // fatal write/config error
	width               int
}

func (m LoginModel) Init() tea.Cmd { return nil }

func (m LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil

	case tea.KeyPressMsg:
		return m.handleKey(msg)

	case tea.PasteMsg:
		switch m.step {
		case stepToken:
			for _, r := range msg.Content {
				if r != '\n' && r != '\r' && r != '\t' {
					m.tokenRunes = append(m.tokenRunes, r)
				}
			}
		case stepProfileNameEntry:
			for _, r := range msg.Content {
				if r != '\n' && r != '\r' {
					m.nameRunes = append(m.nameRunes, r)
				}
			}
		}
		return m, nil

	case loginSpinnerMsg:
		m.spinnerIdx = (m.spinnerIdx + 1) % len(loginSpinnerFrames)
		if m.step == stepValidating {
			return m, loginTickSpinner()
		}
		return m, nil

	case loginOrgsMsg:
		orgs := []*tfe.Organization(msg)
		if len(orgs) == 1 {
			m.selectedOrg = orgs[0].Name
			return m.finalize()
		}
		m.orgs = orgs
		m.orgCursor = 0
		m.step = stepOrgSelect
		return m, nil

	case loginErrMsg:
		// Validation failure — let the user retry or accept anyway rather than hard-exit.
		m.tokenErr = msg.err
		m.tokenErrCursor = 0
		m.step = stepTokenError
		return m, nil
	}
	return m, nil
}

func (m LoginModel) handleKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	k := msg.String()

	if k == "ctrl+c" {
		m.step = stepCancelled
		return m, tea.Quit
	}

	switch m.step {
	case stepProfileList:
		maxCursor := len(m.profiles)
		switch k {
		case "j", "down":
			if m.profileCursor < maxCursor {
				m.profileCursor++
			}
		case "k", "up":
			if m.profileCursor > 0 {
				m.profileCursor--
			}
		case "enter":
			if m.profileCursor == 0 {
				m.isUpdate = false
				m.selectedProfileName = ""
				m.profileNameCursor = 0
				m.step = stepProfileName
			} else {
				p := m.profiles[m.profileCursor-1]
				m.hostname = p.Hostname
				m.selectedProfileName = p.Name
				m.isUpdate = true
				m.step = stepMenu
			}
		case "q", "esc":
			m.step = stepCancelled
			return m, tea.Quit
		}

	case stepMenu:
		switch k {
		case "j", "down":
			if m.menuCursor < 1 {
				m.menuCursor++
			}
		case "k", "up":
			if m.menuCursor > 0 {
				m.menuCursor--
			}
		case "enter":
			m.useBrowser = m.menuCursor == 0
			if m.useBrowser {
				tokenURL := fmt.Sprintf("https://%s/app/settings/tokens?source=tfx-login", m.hostname)
				_ = browser.Open(tokenURL)
			}
			m.step = stepToken
		case "q", "esc":
			if m.isUpdate && len(m.profiles) > 0 {
				m.step = stepProfileList
			} else {
				m.step = stepProfileName
			}
		}

	case stepProfileName:
		switch k {
		case "j", "down":
			if m.profileNameCursor < 1 {
				m.profileNameCursor++
			}
		case "k", "up":
			if m.profileNameCursor > 0 {
				m.profileNameCursor--
			}
		case "enter":
			if m.profileNameCursor == 0 {
				// Use "default" as the profile name
				m.selectedProfileName = "default"
				m.step = stepMenu
			} else {
				// Enter custom name
				m.nameRunes = nil
				m.step = stepProfileNameEntry
			}
		case "esc":
			if len(m.profiles) > 0 {
				m.step = stepProfileList
			} else {
				m.step = stepCancelled
				return m, tea.Quit
			}
		}

	case stepProfileNameEntry:
		switch k {
		case "enter":
			name := strings.TrimSpace(string(m.nameRunes))
			if name == "" {
				return m, nil
			}
			m.selectedProfileName = name
			m.step = stepMenu
		case "backspace":
			if len(m.nameRunes) > 0 {
				m.nameRunes = m.nameRunes[:len(m.nameRunes)-1]
			}
		case "esc":
			m.step = stepProfileName
		default:
			if isPrintable(k) {
				m.nameRunes = append(m.nameRunes, []rune(k)[0])
			}
		}

	case stepToken:
		switch k {
		case "enter":
			token := strings.TrimSpace(string(m.tokenRunes))
			if token == "" {
				return m, nil
			}
			m.resolvedToken = token
			m.step = stepValidating
			return m, tea.Batch(loginFetchOrgs(m.hostname, token), loginTickSpinner())
		case "backspace":
			if len(m.tokenRunes) > 0 {
				m.tokenRunes = m.tokenRunes[:len(m.tokenRunes)-1]
			}
		case "esc":
			m.step = stepMenu
		default:
			if isPrintable(k) {
				m.tokenRunes = append(m.tokenRunes, []rune(k)[0])
			}
		}

	case stepTokenError:
		switch k {
		case "j", "down":
			if m.tokenErrCursor < 1 {
				m.tokenErrCursor++
			}
		case "k", "up":
			if m.tokenErrCursor > 0 {
				m.tokenErrCursor--
			}
		case "enter":
			if m.tokenErrCursor == 0 {
				// Re-enter: clear token, go back to token input
				m.tokenRunes = nil
				m.resolvedToken = ""
				m.step = stepToken
			} else {
				// Accept anyway: save with no org (commented placeholder written by hclconfig)
				m.selectedOrg = ""
				return m.finalize()
			}
		case "esc":
			// Go back to token input keeping the current token for editing
			m.step = stepToken
		}

	case stepOrgSelect:
		switch k {
		case "j", "down":
			if m.orgCursor < len(m.orgs)-1 {
				m.orgCursor++
			}
		case "k", "up":
			if m.orgCursor > 0 {
				m.orgCursor--
			}
		case "enter":
			m.selectedOrg = m.orgs[m.orgCursor].Name
			return m.finalize()
		case "esc":
			m.tokenRunes = nil
			m.resolvedToken = ""
			m.step = stepToken
		}

	case stepValidating:
		// no keys while loading
	}

	return m, nil
}

// finalize writes the profile and exits.
func (m LoginModel) finalize() (tea.Model, tea.Cmd) {
	name := m.selectedProfileName
	if name == "" {
		name = "default"
	}
	if err := hclconfig.WriteProfile(m.configPath, name, m.hostname, m.selectedOrg, m.resolvedToken); err != nil {
		m.err = err
		m.step = stepError
	} else {
		m.step = stepDone
	}
	return m, tea.Quit
}

// tokenLooksValid returns true when the token matches the HCP Terraform /
// TFE format: <prefix>.atlasv1.<suffix>
func tokenLooksValid(runes []rune) bool {
	return strings.Contains(string(runes), ".atlasv1.")
}

// ── Rendering ─────────────────────────────────────────────────────────────────

func (m LoginModel) View() tea.View {
	return tea.NewView(m.render())
}

// renderWidth caps to 88 columns; defaults to 72 before first WindowSizeMsg.
func (m LoginModel) renderWidth() int {
	if m.width > 0 && m.width <= 88 {
		return m.width
	}
	if m.width > 88 {
		return 88
	}
	return 72
}

func (m LoginModel) render() string {
	w := m.renderWidth()

	// ── Shared styles ──────────────────────────────────────────────────────────
	hdrBg    := lipgloss.NewStyle().Background(colorHeaderBg).Foreground(colorFg).Width(w)
	hdrTitle := lipgloss.NewStyle().Background(colorHeaderBg).Foreground(colorAccent).Bold(true).Padding(0, 1)
	hdrSub   := lipgloss.NewStyle().Background(colorHeaderBg).Foreground(colorDim).Padding(0, 1)
	divider  := lipgloss.NewStyle().Foreground(colorBorder).Render(strings.Repeat("─", w))

	dim      := lipgloss.NewStyle().Foreground(colorDim)
	accent   := lipgloss.NewStyle().Foreground(colorAccent)
	success  := lipgloss.NewStyle().Foreground(colorSuccess)
	warn     := lipgloss.NewStyle().Foreground(colorLoading)
	errStyle := lipgloss.NewStyle().Foreground(colorError)
	selected := lipgloss.NewStyle().Foreground(colorAccent).Bold(true)
	hint     := lipgloss.NewStyle().Foreground(colorDim).Italic(true)

	const pad = "  "

	var b strings.Builder

	// ── Header bar ────────────────────────────────────────────────────────────
	b.WriteString(hdrBg.Render(hdrTitle.Render("TFx Login")+hdrSub.Render(m.hostname)) + "\n")
	b.WriteString(divider + "\n\n")

	// ── Step content ──────────────────────────────────────────────────────────
	switch m.step {
	case stepProfileList:
		b.WriteString(pad + dim.Render("Select a profile:") + "\n\n")

		if m.profileCursor == 0 {
			b.WriteString(pad + selected.Render("▸ + Add new profile") + "\n")
		} else {
			b.WriteString(pad + dim.Render("  + Add new profile") + "\n")
		}

		for i, p := range m.profiles {
			cursor := i + 1
			label := p.Name
			if p.Hostname != "" && p.Hostname != p.Name {
				label += "  " + dim.Render("("+p.Hostname+")")
			}
			org := p.Organization
			if org == "" {
				org = "(no org)"
			}
			if m.profileCursor == cursor {
				b.WriteString(pad + selected.Render("▸ "+label) + "  " + dim.Render(org) + "\n")
			} else {
				b.WriteString(pad + dim.Render("  "+label) + "  " + dim.Render(org) + "\n")
			}
		}

		b.WriteString("\n" + pad + hint.Render("↑/↓ · Enter to select · q to quit") + "\n")

	case stepMenu:
		b.WriteString(pad + dim.Render("Authenticate to "+m.hostname) + "\n\n")

		menuItems := []string{
			"Open browser to create a token",
			"Enter token directly",
		}
		for i, item := range menuItems {
			if i == m.menuCursor {
				b.WriteString(pad + selected.Render("▸ "+item) + "\n")
			} else {
				b.WriteString(pad + dim.Render("  "+item) + "\n")
			}
		}

		b.WriteString("\n" + pad + hint.Render("↑/↓ · Enter to select · Esc to go back · q to quit") + "\n")

	case stepProfileName:
		b.WriteString(pad + dim.Render("Choose a name for this profile.") + "\n\n")

		options := []string{
			`Use "default"`,
			"Enter custom name",
		}
		for i, opt := range options {
			if i == m.profileNameCursor {
				b.WriteString(pad + selected.Render("▸ ") + accent.Render(opt) + "\n")
			} else {
				b.WriteString(pad + dim.Render("  "+opt) + "\n")
			}
		}

		b.WriteString("\n" + pad + hint.Render("↑/↓ · Enter to select · Esc to go back") + "\n")

	case stepProfileNameEntry:
		b.WriteString(pad + dim.Render("Enter a name for this profile.") + "\n\n")

		nameDisplay := string(m.nameRunes)
		b.WriteString(pad + "Name  " + accent.Render(nameDisplay) + dim.Render(" _") + "\n\n")
		b.WriteString(pad + hint.Render("Enter") + dim.Render(" to continue · ") +
			hint.Render("Backspace") + dim.Render(" to delete · ") +
			hint.Render("Esc") + dim.Render(" to go back") + "\n")

	case stepToken:
		if m.isUpdate {
			b.WriteString(pad + warn.Render("⚠  Re-authenticating "+m.hostname+" — this will replace the existing token.") + "\n\n")
		}

		if m.useBrowser {
			tokenURL := fmt.Sprintf("https://%s/app/settings/tokens?source=tfx-login", m.hostname)
			b.WriteString(pad + "Browser opened to:\n")
			b.WriteString(pad + dim.Render("  "+tokenURL) + "\n")
			b.WriteString(pad + dim.Render("Generate a token then paste or type it below.") + "\n\n")
		} else {
			b.WriteString(pad + dim.Render("Enter your API token below.") + "\n\n")
		}

		dots := strings.Repeat("●", len(m.tokenRunes))
		tokenLine := pad + "Token  " + accent.Render(dots) + dim.Render(" _")
		if len(m.tokenRunes) > 0 && tokenLooksValid(m.tokenRunes) {
			tokenLine += "  " + success.Render("✓ looks right")
		} else if len(m.tokenRunes) > 0 {
			tokenLine += "  " + dim.Render("paste your token from the browser")
		}
		b.WriteString(tokenLine + "\n\n")
		b.WriteString(pad + hint.Render("Enter") + dim.Render(" to continue · ") +
			hint.Render("Backspace") + dim.Render(" to delete · ") +
			hint.Render("Esc") + dim.Render(" to go back") + "\n")

	case stepValidating:
		spinner := loginSpinnerFrames[m.spinnerIdx]
		b.WriteString(pad + accent.Render(spinner) + "  Validating token and fetching organizations...\n")

	case stepTokenError:
		b.WriteString(pad + errStyle.Render("✗  "+m.tokenErr.Error()) + "\n\n")
		b.WriteString(pad + "How would you like to proceed?\n\n")

		options := []string{
			"Re-enter token",
			"Accept anyway",
		}
		for i, opt := range options {
			if i == m.tokenErrCursor {
				b.WriteString(pad + selected.Render("▸ "+opt) + "\n")
			} else {
				b.WriteString(pad + dim.Render("  "+opt) + "\n")
			}
		}
		b.WriteString("\n" + pad + hint.Render("↑/↓ · Enter to select · Esc to edit token") + "\n")

	case stepOrgSelect:
		b.WriteString(pad + "Select an organization:\n\n")
		for i, org := range m.orgs {
			if i == m.orgCursor {
				b.WriteString(pad + selected.Render("▸ "+org.Name) + "\n")
			} else {
				b.WriteString(pad + "  " + dim.Render(org.Name) + "\n")
			}
		}
		b.WriteString("\n" + pad + hint.Render("↑/↓ or j/k · Enter to select · Esc to go back") + "\n")

	case stepDone:
		action := "created"
		if m.isUpdate {
			action = "updated"
		}
		b.WriteString(pad + success.Render("✓  Profile for "+m.hostname+" has been "+action) + "\n")
		if m.selectedOrg != "" {
			b.WriteString(pad + success.Render("✓  Organization: "+m.selectedOrg) + "\n")
		} else {
			b.WriteString(pad + warn.Render("⚠  Organization not set — edit "+m.configPath+" to configure") + "\n")
		}
		b.WriteString(pad + success.Render("✓  Saved to "+m.configPath) + "\n\n")
		b.WriteString(pad + "Try:  " + accent.Render("tfx workspace list") + "\n")

	case stepError:
		b.WriteString(pad + errStyle.Render("✗  "+m.err.Error()) + "\n")
	}

	return b.String()
}

// ── Entry point ───────────────────────────────────────────────────────────────

// RunLogin launches the inline login wizard for the given hostname.
func RunLogin(hostname string) error {
	configPath, err := hclconfig.DefaultConfigPath()
	if err != nil {
		return err
	}

	profiles, _ := hclconfig.ListProfiles(configPath)

	initialStep := stepProfileName
	if len(profiles) > 0 {
		initialStep = stepProfileList
	}

	m := LoginModel{
		step:       initialStep,
		hostname:   hostname,
		configPath: configPath,
		profiles:   profiles,
	}
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return err
	}
	lm := finalModel.(LoginModel)
	if lm.step == stepCancelled {
		return nil
	}
	return lm.err
}
