// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package tui

import (
	"fmt"
	"image/color"
	"math"
	"os/exec"
	"runtime"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/version"
)

const (
	fixedLines = 10 // header(1) + profilebar(5) + box-top(1) + box-bottom(1) + statusbar(1) + clihint(1)
	minWidth   = 60
	minHeight  = 10
)

// spinnerFrames is the braille-sweep animation sequence for the loading indicator.
var spinnerFrames = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}

// wsTabs defines the ordered workspace sub-view tabs (left → right).
var wsTabs = []struct {
	label string
	view  viewType
}{
	{"Settings", viewWorkspaceSettings},
	{"Runs", viewRuns},
	{"Variables", viewVariables},
	{"Config Versions", viewConfigVersions},
	{"State Versions", viewStateVersions},
}

type viewType int

const (
	viewOrganizations viewType = iota // Phase 6: top-level org list (entry point)
	viewProjects
	viewWorkspaces
	viewWorkspaceSettings        // Settings tab (first workspace sub-view tab)
	viewRuns
	viewVariables                // Phase 5
	viewConfigVersions           // Phase 5
	viewStateVersions            // Phase 5
	viewWorkspaceDetail          // workspace detail (d key from workspace list or sub-views)
	viewOrgDetail                // organization detail (d key from org list)
	viewProjectDetail            // project detail (d key from project list)
	viewRunDetail                // run detail (enter from run list) — Phase 7
	viewVariableDetail           // variable detail (enter from variable list) — Phase 7
	viewStateVersionDetail       // state version detail (d from SV list) — Phase 7
	viewConfigVersionDetail      // config version detail (d from CV list) — Phase 7
	viewStateVersionViewer       // state version JSON viewer (enter from SV list, o from SV detail) — Phase 7b
	viewConfigVersionViewer      // CV file tree browser (enter from CV list, o from CV detail) — Phase 7c
	viewConfigVersionFileContent // CV file content viewer (enter from file browser) — Phase 7c
)

// Model is the root TUI model. All state lives here per the ELM architecture.
type Model struct {
	// Layout
	width  int
	height int
	ready  bool

	// Connection
	c            *client.TfxClient
	hostname     string
	org          string         // active org name (may change when user selects from org list)
	profileName  string         // active profile name from ~/.tfx.hcl
	configFile   string         // path to the config file (tilde-abbreviated for display)
	accountUser      *tfe.User      // currently authenticated user; nil until loaded
	accountToken     *tfe.UserToken // most-recently-used user token; nil until loaded
	accountTokenType accountResourceType // Unknown until loadAccountType resolves

	// View routing
	currentView viewType
	showHelp    bool

	// Loading / error / transient state
	loading      bool
	errMsg       string
	clipFeedback string // cleared on next keypress

	// Spinner state (Phase 5.5)
	spinnerIdx int

	// Organization list state (Phase 6)
	orgs         []*tfe.Organization
	orgCursor    int
	orgOffset    int
	orgFilter    string
	orgFiltering bool
	selectedOrg  *tfe.Organization

	// Project list state
	projects      []*tfe.Project
	projCursor    int
	projOffset    int
	projFilter    string
	projFiltering bool

	// Workspace list state
	workspaces   []*tfe.Workspace
	wsCursor     int
	wsOffset     int
	wsFilter     string
	wsFiltering  bool
	selectedProj *tfe.Project

	// Run list state
	runs         []*tfe.Run
	runCursor    int
	runOffset    int
	runFilter    string
	runFiltering bool
	selectedWS   *tfe.Workspace

	// Variable list state (Phase 5)
	variables    []*tfe.Variable
	varCursor    int
	varOffset    int
	varFilter    string
	varFiltering bool

	// Configuration version list state (Phase 5)
	configVersions []*tfe.ConfigurationVersion
	cvCursor       int
	cvOffset       int
	cvFilter       string
	cvFiltering    bool

	// State version list state (Phase 5)
	stateVersions []*tfe.StateVersion
	svCursor      int
	svOffset      int
	svFilter      string
	svFiltering   bool

	// Workspace detail state
	wsDetScroll        int
	wsDetPrevView      viewType   // view to return to when esc-ing from workspace detail
	wsSettingsScroll   int        // scroll offset for the Settings tab
	wsLatestChange     *time.Time // latest-change-at from API; nil until loaded

	// Organization detail state
	orgDetScroll int

	// Project detail state
	projDetScroll int

	// Run detail state (Phase 7)
	selectedRun  *tfe.Run
	runDetScroll int

	// Variable detail state (Phase 7)
	selectedVar  *tfe.Variable
	varDetScroll int

	// State version detail state (Phase 7)
	selectedSV  *tfe.StateVersion
	svDetScroll int

	// State version viewer state (Phase 7b)
	svJsonLines      []string
	svJsonScroll     int
	svJsonLoading    bool
	svJsonErr        string
	svViewerPrevView viewType // view to return to when esc-ing from the state version viewer

	// Config version detail state (Phase 7)
	selectedCV  *tfe.ConfigurationVersion
	cvDetScroll int

	// Config version viewer state (Phase 7c)
	cvFiles          []cvFile
	cvFileCursor     int
	cvFileOffset     int
	cvFileLoading    bool
	cvFileErr        string
	cvViewerPrevView viewType // view to return to when esc-ing from the config version viewer

	// Config version file content viewer state (Phase 7c)
	cvFileLines  []string
	cvFileScroll int
	cvFileName   string // base name of the currently viewed file

	// API Inspector panel state (Phase 8)
	showDebug       bool              // true = panel is visible (toggled with l)
	debugFocused    bool              // true = right panel has keyboard focus (Tab toggles)
	apiEvents       []client.APIEvent // ring buffer, max 100, newest at index 0
	debugCursor     int               // selected call index in filtered list
	debugDetailMode bool              // true = showing full request/response detail for selected call
	debugBodyScroll int               // scroll offset in the detail view
	debugFilter     string            // case-insensitive method+path filter
	debugFiltering  bool              // filter input is active

	// Instance info modal state (i key — composited on top of current view)
	showInstanceInfo bool // true = modal popup is visible
	infoScroll       int
	healthCheck      map[string]string // nil until loaded
	healthCheckLoad  bool
	healthCheckErr   string

	// CV archived modal state (shown when enter is pressed on an archived CV)
	showCVArchivedModal bool
	cvArchivedCV        *tfe.ConfigurationVersion // the archived CV being shown

	// Tape recorder — non-nil when --tape flag is set
	recorder *TapeRecorder
}

func newModel(c *client.TfxClient, profileName, configFile string, rec *TapeRecorder) Model {
	return Model{
		c:           c,
		hostname:    c.Hostname,
		org:         c.OrganizationName,
		profileName: profileName,
		configFile:  configFile,
		loading:     true,
		recorder:    rec,
	}
}

func (m Model) Init() tea.Cmd {
	// Start the org fetch, account fetch, spinner, and API event listener
	// simultaneously. The event listener runs unconditionally so the inspector
	// panel shows history even when opened after some calls have already completed.
	return tea.Batch(loadOrganizations(m.c), loadAccount(m.c), tickSpinner(), waitForAPIEvent(m.c))
}

// waitForAPIEvent returns a Cmd that blocks until the next API event arrives
// on the client's event bus, then delivers it as a tea.Msg.
// It is re-issued by Update() after each event to keep the subscription alive.
func waitForAPIEvent(c *client.TfxClient) tea.Cmd {
	if c.EventBus == nil {
		return nil
	}
	return func() tea.Msg {
		return <-c.EventBus.Receive()
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// ── Layout ────────────────────────────────────────────────────────────────
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		if m.recorder != nil {
			m.recorder.WriteHeader()
		}

	// ── Spinner ───────────────────────────────────────────────────────────────
	case spinnerTickMsg:
		m.spinnerIdx = (m.spinnerIdx + 1) % len(spinnerFrames)
		if m.loading || m.svJsonLoading || m.cvFileLoading || m.healthCheckLoad {
			return m, tickSpinner()
		}

	// ── API Inspector event bus ───────────────────────────────────────────────
	case client.APIEvent:
		// Prepend (newest first), keep ring buffer at most 100 entries.
		m.apiEvents = append([]client.APIEvent{msg}, m.apiEvents...)
		if len(m.apiEvents) > 100 {
			m.apiEvents = m.apiEvents[:100]
		}
		// Advance cursor so it stays on the same call when not at the top.
		if m.debugCursor > 0 {
			m.debugCursor++
		}
		// Re-subscribe unconditionally to keep the listener alive.
		return m, waitForAPIEvent(m.c)

	// ── Data loads ────────────────────────────────────────────────────────────
	case accountLoadedMsg:
		m.accountUser = msg.user
		m.accountTokenType = msg.resType
		if m.accountUser != nil {
			return m, loadAccountToken(m.c, m.accountUser.ID)
		}

	case accountTokenLoadedMsg:
		m.accountToken = (*tfe.UserToken)(msg)

	case wsLatestChangeMsg:
		m.wsLatestChange = (*time.Time)(msg)
		m.loading = false

	case orgsLoadedMsg:
		m.orgs = []*tfe.Organization(msg)
		m.errMsg = ""

		// Auto-select an org and jump straight to projects if:
		//   1. defaultOrganization is configured and found in the list, OR
		//   2. The token only has access to exactly one org.
		var autoOrg *tfe.Organization
		if m.org != "" {
			for _, o := range m.orgs {
				if o.Name == m.org {
					autoOrg = o
					break
				}
			}
		} else if len(m.orgs) == 1 {
			autoOrg = m.orgs[0]
		}
		if autoOrg != nil {
			m.selectedOrg = autoOrg
			m.org = autoOrg.Name
			m.projCursor = 0
			m.projOffset = 0
			m.projFilter = ""
			m.projFiltering = false
			return m, loadProjects(m.c, autoOrg.Name)
		}

		// Multiple orgs with no default — show the org list.
		m.loading = false
		m.currentView = viewOrganizations

	case projectsLoadedMsg:
		m.projects = []*tfe.Project(msg)
		m.loading = false
		m.currentView = viewProjects
		m.errMsg = ""

	case workspacesLoadedMsg:
		m.workspaces = []*tfe.Workspace(msg)
		m.loading = false
		m.currentView = viewWorkspaces
		m.wsOffset = 0
		m.errMsg = ""

	case runsLoadedMsg:
		m.runs = []*tfe.Run(msg)
		m.loading = false
		m.currentView = viewRuns
		m.runOffset = 0
		m.errMsg = ""

	case variablesLoadedMsg:
		m.variables = []*tfe.Variable(msg)
		m.loading = false
		m.currentView = viewVariables
		m.varOffset = 0
		m.errMsg = ""

	case configVersionsLoadedMsg:
		m.configVersions = []*tfe.ConfigurationVersion(msg)
		m.loading = false
		m.currentView = viewConfigVersions
		m.cvOffset = 0
		m.errMsg = ""

	case stateVersionsLoadedMsg:
		m.stateVersions = []*tfe.StateVersion(msg)
		m.loading = false
		m.currentView = viewStateVersions
		m.svOffset = 0
		m.errMsg = ""

	case runDetailLoadedMsg:
		// Silently update the selected run with full Plan/Apply data.
		// Does not change currentView or loading — detail view updates in-place.
		if msg != nil {
			m.selectedRun = (*tfe.Run)(msg)
		}

	case svJsonLoadedMsg:
		m.svJsonLines = msg.lines
		m.svJsonLoading = false

	case svJsonErrMsg:
		m.svJsonErr = msg.err.Error()
		m.svJsonLoading = false

	case healthCheckLoadedMsg:
		m.healthCheck = map[string]string(msg)
		m.healthCheckLoad = false

	case healthCheckErrMsg:
		m.healthCheckErr = msg.err.Error()
		m.healthCheckLoad = false

	case cvFilesLoadedMsg:
		m.cvFiles = msg.files
		m.cvFileLoading = false

	case cvFileContentLoadedMsg:
		m.cvFileLines = msg.lines
		m.cvFileName = msg.name
		m.cvFileScroll = 0

	case cvFileErrMsg:
		m.cvFileErr = msg.err.Error()
		m.cvFileLoading = false

	case fetchErrMsg:
		m.loading = false
		m.errMsg = msg.err.Error()

	// ── Key input ─────────────────────────────────────────────────────────────
	case tea.KeyPressMsg:
		// Record keypress to tape file if recording is active.
		if m.recorder != nil {
			m.recorder.Record(msg)
		}

		// Clear transient clipboard feedback on next key.
		m.clipFeedback = ""

		// Help overlay consumes all keys.
		if m.showHelp {
			if msg.String() == "esc" || msg.String() == "?" {
				m.showHelp = false
			}
			return m, nil
		}

		// ── Always-global keys — fire even during filter input ─────────────────
		// Only ctrl+c is truly unconditional; all other "global" shortcuts are
		// suppressed while filter input is active so typed characters don't
		// accidentally trigger navigation or quit.
		// shift+c in the inspector detail view copies a curl command.
		if msg.String() == "C" && m.debugFocused && m.debugDetailMode {
			events := m.filteredDebugEvents()
			if m.debugCursor < len(events) {
				curl := buildCurlCommand(events[m.debugCursor])
				if err := copyToClipboard(curl); err == nil {
					m.clipFeedback = "✓ curl command copied to clipboard"
				} else {
					m.clipFeedback = "clipboard unavailable"
				}
			}
			return m, nil
		}
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		// Instance info modal intercepts all remaining keys when open.
		if m.showInstanceInfo {
			return m.handleInstanceInfoModalKey(msg)
		}

		// CV archived modal intercepts all remaining keys when open.
		if m.showCVArchivedModal {
			return m.handleCVArchivedModalKey(msg)
		}

		// ── Filter input mode — must be checked before panel-toggle / quit ────
		// While the user is typing a filter query, only esc, backspace, enter
		// and printable characters are handled; everything else is silently
		// dropped so shortcuts like q, l, r, tab cannot fire mid-filter.
		if m.isFiltering() {
			return m.handleFilterKey(msg)
		}

		// ── Non-filter global keys (suppressed during filter input) ──────────
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "l":
			// Toggle the API Inspector panel.
			m.showDebug = !m.showDebug
			if !m.showDebug {
				m.debugFocused = false
				m.debugDetailMode = false
			}
			return m, nil
		case "tab":
			// Switch focus between left (main) and right (inspector) panels.
			if m.showDebug {
				m.debugFocused = !m.debugFocused
			}
			return m, nil
		}

		// ── Right panel gets all remaining keys when it has focus ────────────
		if m.debugFocused && m.showDebug {
			return m.handleDebugPanelKey(msg)
		}

		// Left-panel global keys.
		switch msg.String() {
		case "?":
			m.showHelp = true
			return m, nil
		case "esc":
			return m.navigateBack()
		case "r":
			return m.refresh()
		case "c":
			cmd := m.currentCliCmd()
			if err := copyToClipboard(cmd); err == nil {
				m.clipFeedback = "✓ copied tfx command to clipboard"
			} else {
				m.clipFeedback = "clipboard unavailable"
			}
			return m, nil
		case "i":
			// Open the instance info modal.
			m.showInstanceInfo = true
			m.infoScroll = 0
			m.healthCheck = nil
			m.healthCheckLoad = true
			m.healthCheckErr = ""
			return m, tea.Batch(loadHealthCheck(m.c), tickSpinner())
		}

		// View-specific navigation (only when not loading).
		if !m.loading {
			switch m.currentView {
			case viewOrganizations:
				return m.handleOrgsKey(msg)
			case viewProjects:
				return m.handleProjectsKey(msg)
			case viewWorkspaces:
				return m.handleWorkspacesKey(msg)
			case viewWorkspaceSettings:
				return m.handleWorkspaceSettingsKey(msg)
			case viewRuns:
				return m.handleRunsKey(msg)
			case viewVariables:
				return m.handleVariablesKey(msg)
			case viewConfigVersions:
				return m.handleConfigVersionsKey(msg)
			case viewStateVersions:
				return m.handleStateVersionsKey(msg)
			case viewWorkspaceDetail:
				return m.handleWorkspaceDetailKey(msg)
			case viewOrgDetail:
				return m.handleOrgDetailKey(msg)
			case viewProjectDetail:
				return m.handleProjectDetailKey(msg)
			case viewRunDetail:
				return m.handleRunDetailKey(msg)
			case viewVariableDetail:
				return m.handleVariableDetailKey(msg)
			case viewStateVersionDetail:
				return m.handleSVDetailKey(msg)
			case viewConfigVersionDetail:
				return m.handleCVDetailKey(msg)
			case viewStateVersionViewer:
				return m.handleSVJsonKey(msg)
			case viewConfigVersionViewer:
				return m.handleCVFilesKey(msg)
			case viewConfigVersionFileContent:
				return m.handleCVFileContentKey(msg)
			}
		}
	}

	return m, nil
}

func (m Model) View() tea.View {
	var content string

	if !m.ready {
		content = "\n  Initializing..."
	} else if m.width < minWidth || m.height < minHeight {
		content = fmt.Sprintf("\n  Terminal too small (%dx%d). Minimum: %dx%d.", m.width, m.height, minWidth, minHeight)
	} else if m.showHelp {
		content = m.renderHelpOverlay()
	} else if m.showDebug {
		// Split content box: unified ┌──┬──┐ / │  │  │ / └──┴──┘ border
		// with "API Inspector" title embedded in the top border.
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			m.renderHeader(),
			m.renderProfileBar(),
			m.renderSplitContentBox(),
			m.renderStatusBar(),
			m.renderCliHint(),
		)
	} else {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			m.renderHeader(),
			m.renderProfileBar(),
			m.renderContentBox(),
			m.renderStatusBar(),
			m.renderCliHint(),
		)
	}

	// Composite the instance info modal on top when visible.
	// (Not composited over the help overlay — help takes full priority.)
	if m.ready && m.showInstanceInfo && !m.showHelp && m.width >= minWidth && m.height >= minHeight {
		content = m.overlayInstanceInfoModal(content)
	}

	// Composite the CV archived modal on top when visible.
	if m.ready && m.showCVArchivedModal && !m.showHelp && m.width >= minWidth && m.height >= minHeight {
		content = m.overlayCVArchivedModal(content)
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

// ── Navigation ────────────────────────────────────────────────────────────────

func (m Model) navigateBack() (tea.Model, tea.Cmd) {
	switch m.currentView {
	case viewProjects:
		// Return to org list (don't re-fetch — orgs are already loaded).
		m.currentView = viewOrganizations
		m.projects = nil
		m.projCursor, m.projOffset = 0, 0
		m.projFilter, m.projFiltering = "", false
		m.selectedProj = nil
	case viewWorkspaces:
		m.currentView = viewProjects
		m.workspaces = nil
		m.wsCursor, m.wsOffset = 0, 0
		m.wsFilter, m.wsFiltering = "", false
		m.selectedProj = nil
	case viewWorkspaceSettings, viewRuns, viewVariables, viewConfigVersions, viewStateVersions:
		// All workspace tab views return to the workspace list and clear sub-view data.
		m.currentView = viewWorkspaces
		m.runs = nil
		m.runCursor, m.runOffset, m.runFilter, m.runFiltering = 0, 0, "", false
		m.variables = nil
		m.varCursor, m.varOffset, m.varFilter, m.varFiltering = 0, 0, "", false
		m.configVersions = nil
		m.cvCursor, m.cvOffset, m.cvFilter, m.cvFiltering = 0, 0, "", false
		m.stateVersions = nil
		m.svCursor, m.svOffset, m.svFilter, m.svFiltering = 0, 0, "", false
		m.selectedWS = nil
	case viewWorkspaceDetail:
		// Return to wherever d was pressed from (workspace list or a sub-view tab).
		m.currentView = m.wsDetPrevView
		m.wsDetScroll = 0
		// Only clear selectedWS when returning to the workspace list; sub-views still need it.
		if m.wsDetPrevView == viewWorkspaces {
			m.selectedWS = nil
		}
	case viewOrgDetail:
		// Return to org list; cursor/offset intentionally preserved.
		m.currentView = viewOrganizations
		m.orgDetScroll = 0
		m.selectedOrg = nil
	case viewProjectDetail:
		// Return to project list; cursor/offset intentionally preserved.
		m.currentView = viewProjects
		m.projDetScroll = 0
		m.selectedProj = nil
	case viewRunDetail:
		m.currentView = viewRuns
		m.runDetScroll = 0
		m.selectedRun = nil
	case viewVariableDetail:
		m.currentView = viewVariables
		m.varDetScroll = 0
		m.selectedVar = nil
	case viewStateVersionDetail:
		m.currentView = viewStateVersions
		m.svDetScroll = 0
		m.selectedSV = nil
	case viewStateVersionViewer:
		m.currentView = m.svViewerPrevView
		m.svJsonLines = nil
		m.svJsonScroll = 0
		m.svJsonLoading = false
		m.svJsonErr = ""
	case viewConfigVersionViewer:
		m.currentView = m.cvViewerPrevView
		m.cvFiles = nil
		m.cvFileCursor = 0
		m.cvFileOffset = 0
		m.cvFileLoading = false
		m.cvFileErr = ""
	case viewConfigVersionFileContent:
		m.currentView = viewConfigVersionViewer
		m.cvFileLines = nil
		m.cvFileScroll = 0
		m.cvFileName = ""
	case viewConfigVersionDetail:
		m.currentView = viewConfigVersions
		m.cvDetScroll = 0
		m.selectedCV = nil
	}
	return m, nil
}

func (m Model) refresh() (tea.Model, tea.Cmd) {
	m.loading = true
	m.errMsg = ""
	var cmd tea.Cmd
	switch m.currentView {
	case viewOrganizations:
		m.orgs = nil
		m.orgCursor, m.orgOffset = 0, 0
		cmd = loadOrganizations(m.c)
	case viewProjects:
		m.projects = nil
		m.projCursor, m.projOffset = 0, 0
		cmd = loadProjects(m.c, m.org)
	case viewWorkspaces:
		m.workspaces = nil
		m.wsCursor, m.wsOffset = 0, 0
		projectID := ""
		if m.selectedProj != nil {
			projectID = m.selectedProj.ID
		}
		cmd = loadWorkspaces(m.c, m.org, projectID)
	case viewWorkspaceSettings:
		m.wsLatestChange = nil
		if m.selectedWS != nil {
			cmd = loadWorkspaceLatestChange(m.c, m.selectedWS.ID)
		}
	case viewRuns:
		m.runs = nil
		m.runCursor, m.runOffset = 0, 0
		if m.selectedWS != nil {
			cmd = loadRuns(m.c, m.selectedWS.ID)
		}
	case viewVariables:
		m.variables = nil
		m.varCursor, m.varOffset = 0, 0
		if m.selectedWS != nil {
			cmd = loadVariables(m.c, m.selectedWS.ID)
		}
	case viewConfigVersions:
		m.configVersions = nil
		m.cvCursor, m.cvOffset = 0, 0
		if m.selectedWS != nil {
			cmd = loadConfigVersions(m.c, m.org, m.selectedWS.Name)
		}
	case viewStateVersions:
		m.stateVersions = nil
		m.svCursor, m.svOffset = 0, 0
		if m.selectedWS != nil {
			cmd = loadStateVersions(m.c, m.org, m.selectedWS.Name)
		}
	case viewStateVersionViewer:
		// Force re-download: bypass cache and restart the JSON load.
		m.loading = false
		if m.selectedSV != nil {
			m.svJsonLines = nil
			m.svJsonLoading = true
			m.svJsonErr = ""
			return m, tea.Batch(loadStateVersionJson(m.c, m.selectedSV.ID, true), tickSpinner())
		}
		return m, nil
	case viewConfigVersionViewer:
		// Force re-download: delete cached archive + extracted dir.
		m.loading = false
		if m.selectedCV != nil {
			m.cvFiles = nil
			m.cvFileLoading = true
			m.cvFileErr = ""
			return m, tea.Batch(loadCVFiles(m.c, m.selectedCV.ID, true), tickSpinner())
		}
		return m, nil
	case viewConfigVersionFileContent:
		// File is already on disk; no-op (re-read not needed for MVP).
		m.loading = false
		return m, nil
	case viewWorkspaceDetail, viewOrgDetail, viewProjectDetail,
		viewRunDetail, viewVariableDetail, viewStateVersionDetail, viewConfigVersionDetail:
		// Detail views show already-loaded data; nothing to refresh.
		m.loading = false
		return m, nil
	}
	return m, tea.Batch(cmd, tickSpinner())
}

// isWorkspaceSubView returns true when the current view is one of the
// workspace tab views (Settings, Runs, Variables, Config Versions, State Versions).
func (m Model) isWorkspaceSubView() bool {
	switch m.currentView {
	case viewWorkspaceSettings, viewRuns, viewVariables, viewConfigVersions, viewStateVersions:
		return true
	}
	return false
}

// switchWsTab transitions to the target workspace tab. If the data for that
// tab is already cached it switches instantly; otherwise it triggers a fetch.
func (m Model) switchWsTab(target viewType) (tea.Model, tea.Cmd) {
	m.currentView = target
	m.errMsg = ""

	switch target {
	case viewWorkspaceSettings:
		if m.wsLatestChange == nil && m.selectedWS != nil {
			return m, loadWorkspaceLatestChange(m.c, m.selectedWS.ID)
		}
		return m, nil
	case viewRuns:
		if m.runs != nil {
			return m, nil
		}
		m.loading = true
		return m, tea.Batch(loadRuns(m.c, m.selectedWS.ID), tickSpinner())
	case viewVariables:
		if m.variables != nil {
			return m, nil
		}
		m.loading = true
		return m, tea.Batch(loadVariables(m.c, m.selectedWS.ID), tickSpinner())
	case viewConfigVersions:
		if m.configVersions != nil {
			return m, nil
		}
		m.loading = true
		return m, tea.Batch(loadConfigVersions(m.c, m.org, m.selectedWS.Name), tickSpinner())
	case viewStateVersions:
		if m.stateVersions != nil {
			return m, nil
		}
		m.loading = true
		return m, tea.Batch(loadStateVersions(m.c, m.org, m.selectedWS.Name), tickSpinner())
	}
	return m, nil
}

// ── Key handlers ──────────────────────────────────────────────────────────────

func (m Model) isFiltering() bool {
	switch m.currentView {
	case viewOrganizations:
		return m.orgFiltering
	case viewProjects:
		return m.projFiltering
	case viewWorkspaces:
		return m.wsFiltering
	case viewRuns:
		return m.runFiltering
	case viewVariables:
		return m.varFiltering
	case viewConfigVersions:
		return m.cvFiltering
	case viewStateVersions:
		return m.svFiltering
	}
	return false
}

func (m Model) handleFilterKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		switch m.currentView {
		case viewOrganizations:
			m.orgFilter, m.orgFiltering = "", false
			m.orgCursor, m.orgOffset = 0, 0
		case viewProjects:
			m.projFilter, m.projFiltering = "", false
			m.projCursor, m.projOffset = 0, 0
		case viewWorkspaces:
			m.wsFilter, m.wsFiltering = "", false
			m.wsCursor, m.wsOffset = 0, 0
		case viewRuns:
			m.runFilter, m.runFiltering = "", false
			m.runCursor, m.runOffset = 0, 0
		case viewVariables:
			m.varFilter, m.varFiltering = "", false
			m.varCursor, m.varOffset = 0, 0
		case viewConfigVersions:
			m.cvFilter, m.cvFiltering = "", false
			m.cvCursor, m.cvOffset = 0, 0
		case viewStateVersions:
			m.svFilter, m.svFiltering = "", false
			m.svCursor, m.svOffset = 0, 0
		}
	case "enter":
		switch m.currentView {
		case viewOrganizations:
			m.orgFiltering = false
		case viewProjects:
			m.projFiltering = false
		case viewWorkspaces:
			m.wsFiltering = false
		case viewRuns:
			m.runFiltering = false
		case viewVariables:
			m.varFiltering = false
		case viewConfigVersions:
			m.cvFiltering = false
		case viewStateVersions:
			m.svFiltering = false
		}
	case "backspace":
		switch m.currentView {
		case viewOrganizations:
			if r := []rune(m.orgFilter); len(r) > 0 {
				m.orgFilter = string(r[:len(r)-1])
				m.orgCursor, m.orgOffset = 0, 0
			}
		case viewProjects:
			if r := []rune(m.projFilter); len(r) > 0 {
				m.projFilter = string(r[:len(r)-1])
				m.projCursor, m.projOffset = 0, 0
			}
		case viewWorkspaces:
			if r := []rune(m.wsFilter); len(r) > 0 {
				m.wsFilter = string(r[:len(r)-1])
				m.wsCursor, m.wsOffset = 0, 0
			}
		case viewRuns:
			if r := []rune(m.runFilter); len(r) > 0 {
				m.runFilter = string(r[:len(r)-1])
				m.runCursor, m.runOffset = 0, 0
			}
		case viewVariables:
			if r := []rune(m.varFilter); len(r) > 0 {
				m.varFilter = string(r[:len(r)-1])
				m.varCursor, m.varOffset = 0, 0
			}
		case viewConfigVersions:
			if r := []rune(m.cvFilter); len(r) > 0 {
				m.cvFilter = string(r[:len(r)-1])
				m.cvCursor, m.cvOffset = 0, 0
			}
		case viewStateVersions:
			if r := []rune(m.svFilter); len(r) > 0 {
				m.svFilter = string(r[:len(r)-1])
				m.svCursor, m.svOffset = 0, 0
			}
		}
	default:
		// Normalize "space" (Bubble Tea v2's key name for the space bar) to " ".
		s := msg.String()
		if s == "space" {
			s = " "
		}
		if r := []rune(s); len(r) == 1 && r[0] >= 32 {
			ch := string(r)
			switch m.currentView {
			case viewOrganizations:
				m.orgFilter += ch
				m.orgCursor, m.orgOffset = 0, 0
			case viewProjects:
				m.projFilter += ch
				m.projCursor, m.projOffset = 0, 0
			case viewWorkspaces:
				m.wsFilter += ch
				m.wsCursor, m.wsOffset = 0, 0
			case viewRuns:
				m.runFilter += ch
				m.runCursor, m.runOffset = 0, 0
			case viewVariables:
				m.varFilter += ch
				m.varCursor, m.varOffset = 0, 0
			case viewConfigVersions:
				m.cvFilter += ch
				m.cvCursor, m.cvOffset = 0, 0
			case viewStateVersions:
				m.svFilter += ch
				m.svCursor, m.svOffset = 0, 0
			}
		}
	}
	return m, nil
}

func (m Model) handleOrgsKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	filtered := filteredOrgs(m)
	n := len(filtered)
	vis := m.orgVisibleRows()

	switch msg.String() {
	case "up", "k":
		if m.orgCursor > 0 {
			m.orgCursor--
			if m.orgCursor < m.orgOffset {
				m.orgOffset = m.orgCursor
			}
		}
	case "down", "j":
		if m.orgCursor < n-1 {
			m.orgCursor++
			if m.orgCursor >= m.orgOffset+vis {
				m.orgOffset = m.orgCursor - vis + 1
			}
		}
	case "g":
		m.orgCursor, m.orgOffset = 0, 0
	case "G":
		if n > 0 {
			m.orgCursor = n - 1
			if n > vis {
				m.orgOffset = n - vis
			}
		}
	case "/":
		m.orgFiltering = true
	case "d":
		if n == 0 || m.orgCursor >= n {
			break
		}
		m.selectedOrg = filtered[m.orgCursor]
		m.orgDetScroll = 0
		m.currentView = viewOrgDetail
	case "enter":
		if n == 0 || m.orgCursor >= n {
			break
		}
		sel := filtered[m.orgCursor]
		m.selectedOrg = sel
		m.org = sel.Name
		m.loading = true
		m.errMsg = ""
		m.projCursor, m.projOffset, m.projFilter = 0, 0, ""
		return m, tea.Batch(loadProjects(m.c, sel.Name), tickSpinner())
	case "u":
		m = m.applyURL(m.hostnameURL(), "hostname", false)
	case "U":
		m = m.applyURL(m.hostnameURL(), "hostname", true)
	}
	return m, nil
}

func (m Model) handleOrgDetailKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.orgDetScroll > 0 {
			m.orgDetScroll--
		}
	case "down", "j":
		m.orgDetScroll++
	case "g":
		m.orgDetScroll = 0
	case "G":
		m.orgDetScroll = 9999
	case "u":
		m = m.applyURL(m.orgURL(), "org", false)
	case "U":
		m = m.applyURL(m.orgURL(), "org", true)
	}
	return m, nil
}

func (m Model) handleProjectsKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	filtered := filteredProjects(m)
	n := len(filtered)
	vis := m.projVisibleRows()

	switch msg.String() {
	case "up", "k":
		if m.projCursor > 0 {
			m.projCursor--
			if m.projCursor < m.projOffset {
				m.projOffset = m.projCursor
			}
		}
	case "down", "j":
		if m.projCursor < n-1 {
			m.projCursor++
			if m.projCursor >= m.projOffset+vis {
				m.projOffset = m.projCursor - vis + 1
			}
		}
	case "g":
		m.projCursor, m.projOffset = 0, 0
	case "G":
		if n > 0 {
			m.projCursor = n - 1
			if n > vis {
				m.projOffset = n - vis
			}
		}
	case "/":
		m.projFiltering = true
	case "d":
		if n == 0 || m.projCursor >= n {
			break
		}
		m.selectedProj = filtered[m.projCursor]
		m.projDetScroll = 0
		m.currentView = viewProjectDetail
	case "enter":
		if n == 0 || m.projCursor >= n {
			break
		}
		sel := filtered[m.projCursor]
		m.selectedProj = sel
		m.loading = true
		m.errMsg = ""
		m.wsCursor, m.wsOffset, m.wsFilter = 0, 0, ""
		return m, tea.Batch(loadWorkspaces(m.c, m.org, sel.ID), tickSpinner())
	case "u":
		m = m.applyURL(m.orgProjectsURL(), "projects", false)
	case "U":
		m = m.applyURL(m.orgProjectsURL(), "projects", true)
	}
	return m, nil
}

func (m Model) handleProjectDetailKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.projDetScroll > 0 {
			m.projDetScroll--
		}
	case "down", "j":
		m.projDetScroll++
	case "g":
		m.projDetScroll = 0
	case "G":
		m.projDetScroll = 9999
	case "u":
		m = m.applyURL(m.projURL(), "project", false)
	case "U":
		m = m.applyURL(m.projURL(), "project", true)
	}
	return m, nil
}

func (m Model) handleRunDetailKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.runDetScroll > 0 {
			m.runDetScroll--
		}
	case "down", "j":
		m.runDetScroll++
	case "g":
		m.runDetScroll = 0
	case "G":
		m.runDetScroll = 9999
	case "u":
		m = m.applyURL(m.runURL(), "run", false)
	case "U":
		m = m.applyURL(m.runURL(), "run", true)
	}
	return m, nil
}

func (m Model) handleVariableDetailKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.varDetScroll > 0 {
			m.varDetScroll--
		}
	case "down", "j":
		m.varDetScroll++
	case "g":
		m.varDetScroll = 0
	case "G":
		m.varDetScroll = 9999
	case "u":
		m = m.applyURL(m.wsVariablesListURL(), "variables", false)
	case "U":
		m = m.applyURL(m.wsVariablesListURL(), "variables", true)
	}
	return m, nil
}

func (m Model) handleSVDetailKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.svDetScroll > 0 {
			m.svDetScroll--
		}
	case "down", "j":
		m.svDetScroll++
	case "g":
		m.svDetScroll = 0
	case "G":
		m.svDetScroll = 9999
	case "o":
		if m.selectedSV != nil {
			m.svJsonLines = nil
			m.svJsonScroll = 0
			m.svJsonLoading = true
			m.svJsonErr = ""
			m.svViewerPrevView = viewStateVersionDetail
			m.currentView = viewStateVersionViewer
			return m, tea.Batch(loadStateVersionJson(m.c, m.selectedSV.ID, false), tickSpinner())
		}
	case "u":
		m = m.applyURL(m.svURL(), "state version", false)
	case "U":
		m = m.applyURL(m.svURL(), "state version", true)
	}
	return m, nil
}

func (m Model) handleSVJsonKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.svJsonScroll > 0 {
			m.svJsonScroll--
		}
	case "down", "j":
		m.svJsonScroll++
	case "shift+up":
		half := m.contentHeight() / 2
		if m.svJsonScroll > half {
			m.svJsonScroll -= half
		} else {
			m.svJsonScroll = 0
		}
	case "shift+down":
		m.svJsonScroll += m.contentHeight() / 2
	case "g":
		m.svJsonScroll = 0
	case "G":
		m.svJsonScroll = 9999
	}
	return m, nil
}

func (m Model) handleCVDetailKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cvDetScroll > 0 {
			m.cvDetScroll--
		}
	case "down", "j":
		m.cvDetScroll++
	case "g":
		m.cvDetScroll = 0
	case "G":
		m.cvDetScroll = 9999
	case "o":
		if m.selectedCV != nil {
			m.cvFiles = nil
			m.cvFileCursor = 0
			m.cvFileOffset = 0
			m.cvFileLoading = true
			m.cvFileErr = ""
			m.cvViewerPrevView = viewConfigVersionDetail
			m.currentView = viewConfigVersionViewer
			return m, tea.Batch(loadCVFiles(m.c, m.selectedCV.ID, false), tickSpinner())
		}
	}
	return m, nil
}

func (m Model) handleCVFilesKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	n := len(m.cvFiles)
	vis := m.cvFilesVisibleRows()
	switch msg.String() {
	case "up", "k":
		if m.cvFileCursor > 0 {
			m.cvFileCursor--
			if m.cvFileCursor < m.cvFileOffset {
				m.cvFileOffset = m.cvFileCursor
			}
		}
	case "down", "j":
		if m.cvFileCursor < n-1 {
			m.cvFileCursor++
			if m.cvFileCursor >= m.cvFileOffset+vis {
				m.cvFileOffset = m.cvFileCursor - vis + 1
			}
		}
	case "g":
		m.cvFileCursor, m.cvFileOffset = 0, 0
	case "G":
		m.cvFileCursor = n - 1
		if n > vis {
			m.cvFileOffset = n - vis
		}
	case "enter":
		if n == 0 || m.cvFileCursor >= n {
			break
		}
		sel := m.cvFiles[m.cvFileCursor]
		if sel.isDir {
			break // directories are not openable in MVP
		}
		m.cvFileLines = nil
		m.cvFileScroll = 0
		m.cvFileName = sel.displayName()
		m.currentView = viewConfigVersionFileContent
		return m, loadCVFileContent(m.selectedCV.ID, sel)
	case "p":
		// Copy the on-disk extraction path to the clipboard.
		if m.selectedCV != nil {
			if err := copyToClipboard(cvExtractDirPath(m.selectedCV.ID)); err == nil {
				m.clipFeedback = "✓ path copied to clipboard"
			} else {
				m.clipFeedback = "clipboard unavailable"
			}
		}
	}
	return m, nil
}

func (m Model) handleCVFileContentKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cvFileScroll > 0 {
			m.cvFileScroll--
		}
	case "down", "j":
		m.cvFileScroll++
	case "shift+up":
		half := m.contentHeight() / 2
		if m.cvFileScroll > half {
			m.cvFileScroll -= half
		} else {
			m.cvFileScroll = 0
		}
	case "shift+down":
		m.cvFileScroll += m.contentHeight() / 2
	case "g":
		m.cvFileScroll = 0
	case "G":
		m.cvFileScroll = 9999
	}
	return m, nil
}

func (m Model) handleWorkspacesKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	filtered := filteredWorkspaces(m)
	n := len(filtered)
	vis := m.wsVisibleRows()

	switch msg.String() {
	case "up", "k":
		if m.wsCursor > 0 {
			m.wsCursor--
			if m.wsCursor < m.wsOffset {
				m.wsOffset = m.wsCursor
			}
		}
	case "down", "j":
		if m.wsCursor < n-1 {
			m.wsCursor++
			if m.wsCursor >= m.wsOffset+vis {
				m.wsOffset = m.wsCursor - vis + 1
			}
		}
	case "g":
		m.wsCursor, m.wsOffset = 0, 0
	case "G":
		if n > 0 {
			m.wsCursor = n - 1
			if n > vis {
				m.wsOffset = n - vis
			}
		}
	case "/":
		m.wsFiltering = true
	case "enter":
		if n == 0 || m.wsCursor >= n {
			break
		}
		sel := filtered[m.wsCursor]
		m.selectedWS = sel
		m.errMsg = ""
		m.wsSettingsScroll = 0
		m.wsLatestChange = nil
		m.currentView = viewWorkspaceSettings
		return m, loadWorkspaceLatestChange(m.c, sel.ID)
	case "v":
		if n == 0 || m.wsCursor >= n {
			break
		}
		sel := filtered[m.wsCursor]
		m.selectedWS = sel
		m.loading = true
		m.errMsg = ""
		m.varCursor, m.varOffset, m.varFilter = 0, 0, ""
		return m, tea.Batch(loadVariables(m.c, sel.ID), tickSpinner())
	case "f":
		if n == 0 || m.wsCursor >= n {
			break
		}
		sel := filtered[m.wsCursor]
		m.selectedWS = sel
		m.loading = true
		m.errMsg = ""
		m.cvCursor, m.cvOffset, m.cvFilter = 0, 0, ""
		return m, tea.Batch(loadConfigVersions(m.c, m.org, sel.Name), tickSpinner())
	case "s":
		if n == 0 || m.wsCursor >= n {
			break
		}
		sel := filtered[m.wsCursor]
		m.selectedWS = sel
		m.loading = true
		m.errMsg = ""
		m.svCursor, m.svOffset, m.svFilter = 0, 0, ""
		return m, tea.Batch(loadStateVersions(m.c, m.org, sel.Name), tickSpinner())
	case "d":
		if n == 0 || m.wsCursor >= n {
			break
		}
		m.selectedWS = filtered[m.wsCursor]
		m.wsDetScroll = 0
		m.wsDetPrevView = viewWorkspaces
		m.currentView = viewWorkspaceDetail
	case "u":
		m = m.applyURL(m.projURL(), "workspace list", false)
	case "U":
		m = m.applyURL(m.projURL(), "workspace list", true)
	}
	return m, nil
}

func (m Model) handleWorkspaceDetailKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.wsDetScroll > 0 {
			m.wsDetScroll--
		}
	case "down", "j":
		m.wsDetScroll++
	case "g":
		m.wsDetScroll = 0
	case "G":
		m.wsDetScroll = 9999 // clamped to max by renderWorkspaceDetailContent
	case "u":
		if url := m.wsURL(); url != "" {
			if err := copyToClipboard(url); err == nil {
				m.clipFeedback = "✓ workspace URL copied"
			} else {
				m.clipFeedback = "clipboard unavailable"
			}
		}
	case "U":
		if url := m.wsURL(); url != "" {
			if err := openBrowser(url); err == nil {
				m.clipFeedback = "✓ opening in browser"
			} else {
				m.clipFeedback = "could not open browser"
			}
		}
	}
	return m, nil
}

func (m Model) handleWorkspaceSettingsKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.wsSettingsScroll > 0 {
			m.wsSettingsScroll--
		}
	case "down", "j":
		m.wsSettingsScroll++
	case "g":
		m.wsSettingsScroll = 0
	case "G":
		m.wsSettingsScroll = 9999 // clamped to max by renderWorkspaceSettingsContent
	case "left":
		// First tab — nothing to the left.
	case "right":
		return m.switchWsTab(viewRuns)
	case "u":
		if url := m.wsURL(); url != "" {
			if err := copyToClipboard(url); err == nil {
				m.clipFeedback = "✓ workspace URL copied"
			} else {
				m.clipFeedback = "clipboard unavailable"
			}
		}
		return m, nil
	case "U":
		if url := m.wsURL(); url != "" {
			if err := openBrowser(url); err == nil {
				m.clipFeedback = "✓ opening in browser"
			} else {
				m.clipFeedback = "could not open browser"
			}
		}
		return m, nil
	}
	return m, nil
}

func (m Model) handleRunsKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	filtered := filteredRuns(m)
	n := len(filtered)
	vis := m.runVisibleRows()

	switch msg.String() {
	case "up", "k":
		if m.runCursor > 0 {
			m.runCursor--
			if m.runCursor < m.runOffset {
				m.runOffset = m.runCursor
			}
		}
	case "down", "j":
		if m.runCursor < n-1 {
			m.runCursor++
			if m.runCursor >= m.runOffset+vis {
				m.runOffset = m.runCursor - vis + 1
			}
		}
	case "g":
		m.runCursor, m.runOffset = 0, 0
	case "G":
		if n > 0 {
			m.runCursor = n - 1
			if n > vis {
				m.runOffset = n - vis
			}
		}
	case "/":
		m.runFiltering = true
	case "left":
		return m.switchWsTab(viewWorkspaceSettings)
	case "right":
		return m.switchWsTab(viewVariables)
	case "enter":
		if n == 0 || m.runCursor >= n {
			break
		}
		sel := filtered[m.runCursor]
		m.selectedRun = sel
		m.runDetScroll = 0
		m.currentView = viewRunDetail
		// Trigger a background re-fetch to populate Plan/Apply/VCS fields.
		return m, loadRunDetail(m.c, sel.ID)
	case "u":
		m = m.applyURL(m.wsRunsListURL(), "runs", false)
		return m, nil
	case "U":
		m = m.applyURL(m.wsRunsListURL(), "runs", true)
		return m, nil
	}
	return m, nil
}

func (m Model) handleVariablesKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	filtered := filteredVariables(m)
	n := len(filtered)
	vis := m.varVisibleRows()

	switch msg.String() {
	case "up", "k":
		if m.varCursor > 0 {
			m.varCursor--
			if m.varCursor < m.varOffset {
				m.varOffset = m.varCursor
			}
		}
	case "down", "j":
		if m.varCursor < n-1 {
			m.varCursor++
			if m.varCursor >= m.varOffset+vis {
				m.varOffset = m.varCursor - vis + 1
			}
		}
	case "g":
		m.varCursor, m.varOffset = 0, 0
	case "G":
		if n > 0 {
			m.varCursor = n - 1
			if n > vis {
				m.varOffset = n - vis
			}
		}
	case "/":
		m.varFiltering = true
	case "enter":
		if n == 0 || m.varCursor >= n {
			break
		}
		m.selectedVar = filtered[m.varCursor]
		m.varDetScroll = 0
		m.currentView = viewVariableDetail
	case "left":
		return m.switchWsTab(viewRuns)
	case "right":
		return m.switchWsTab(viewConfigVersions)
	case "u":
		m = m.applyURL(m.wsVariablesListURL(), "variables", false)
		return m, nil
	case "U":
		m = m.applyURL(m.wsVariablesListURL(), "variables", true)
		return m, nil
	}
	return m, nil
}

func (m Model) handleConfigVersionsKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	filtered := filteredConfigVersions(m)
	n := len(filtered)
	vis := m.cvVisibleRows()

	switch msg.String() {
	case "up", "k":
		if m.cvCursor > 0 {
			m.cvCursor--
			if m.cvCursor < m.cvOffset {
				m.cvOffset = m.cvCursor
			}
		}
	case "down", "j":
		if m.cvCursor < n-1 {
			m.cvCursor++
			if m.cvCursor >= m.cvOffset+vis {
				m.cvOffset = m.cvCursor - vis + 1
			}
		}
	case "g":
		m.cvCursor, m.cvOffset = 0, 0
	case "G":
		if n > 0 {
			m.cvCursor = n - 1
			if n > vis {
				m.cvOffset = n - vis
			}
		}
	case "/":
		m.cvFiltering = true
	case "enter":
		if n == 0 || m.cvCursor >= n {
			break
		}
		cv := filtered[m.cvCursor]
		if cv.Status == tfe.ConfigurationArchived {
			m.cvArchivedCV = cv
			m.showCVArchivedModal = true
			return m, nil
		}
		m.selectedCV = cv
		m.cvFiles = nil
		m.cvFileCursor = 0
		m.cvFileOffset = 0
		m.cvFileLoading = true
		m.cvFileErr = ""
		m.cvViewerPrevView = viewConfigVersions
		m.currentView = viewConfigVersionViewer
		return m, tea.Batch(loadCVFiles(m.c, m.selectedCV.ID, false), tickSpinner())
	case "left":
		return m.switchWsTab(viewVariables)
	case "right":
		return m.switchWsTab(viewStateVersions)
	case "d":
		if n == 0 || m.cvCursor >= n {
			break
		}
		m.selectedCV = filtered[m.cvCursor]
		m.cvDetScroll = 0
		m.currentView = viewConfigVersionDetail
	}
	return m, nil
}

func (m Model) handleStateVersionsKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	filtered := filteredStateVersions(m)
	n := len(filtered)
	vis := m.svVisibleRows()

	switch msg.String() {
	case "up", "k":
		if m.svCursor > 0 {
			m.svCursor--
			if m.svCursor < m.svOffset {
				m.svOffset = m.svCursor
			}
		}
	case "down", "j":
		if m.svCursor < n-1 {
			m.svCursor++
			if m.svCursor >= m.svOffset+vis {
				m.svOffset = m.svCursor - vis + 1
			}
		}
	case "g":
		m.svCursor, m.svOffset = 0, 0
	case "G":
		if n > 0 {
			m.svCursor = n - 1
			if n > vis {
				m.svOffset = n - vis
			}
		}
	case "/":
		m.svFiltering = true
	case "enter":
		if n == 0 || m.svCursor >= n {
			break
		}
		m.selectedSV = filtered[m.svCursor]
		m.svJsonLines = nil
		m.svJsonScroll = 0
		m.svJsonLoading = true
		m.svJsonErr = ""
		m.svViewerPrevView = viewStateVersions
		m.currentView = viewStateVersionViewer
		return m, tea.Batch(loadStateVersionJson(m.c, m.selectedSV.ID, false), tickSpinner())
	case "left":
		return m.switchWsTab(viewConfigVersions)
	case "right":
		// Last tab — nothing to the right.
	case "d":
		if n == 0 || m.svCursor >= n {
			break
		}
		m.selectedSV = filtered[m.svCursor]
		m.svDetScroll = 0
		m.currentView = viewStateVersionDetail
	case "u":
		m = m.applyURL(m.wsSVsListURL(), "state versions", false)
		return m, nil
	case "U":
		m = m.applyURL(m.wsSVsListURL(), "state versions", true)
		return m, nil
	}
	return m, nil
}

// ── Layout helpers ────────────────────────────────────────────────────────────

func (m Model) contentHeight() int {
	h := m.height - fixedLines
	if h < 1 {
		return 1
	}
	return h
}

// debugPanelWidth returns the content width of the API Inspector panel.
// Takes ~67% of the terminal width (the remaining ~33% goes to the main view).
func (m Model) debugPanelWidth() int {
	return (m.width - 3) * 2 / 3 // -3 for left │, middle │, right │ border chars
}

// mainWidth returns the total width of the content box (outer, including side borders).
// When the API Inspector panel is open, the content area is narrowed by the
// panel width plus one separator column.
func (m Model) mainWidth() int {
	if m.showDebug {
		return m.width - m.debugPanelWidth() - 1
	}
	return m.width
}

// innerWidth returns the usable width for content inside the box borders (mainWidth - 2).
// All content renderers must use this so rows fit inside the │ side borders.
func (m Model) innerWidth() int {
	w := m.mainWidth() - 2
	if w < 1 {
		return 1
	}
	return w
}

// padContent fills a rendered content-area string to innerWidth() using the
// given style. Use this instead of pad() inside all content-area renderers
// so the row width automatically accounts for the debug panel when it is open.
func (m Model) padContent(rendered string, style lipgloss.Style) string {
	w := m.innerWidth() - lipgloss.Width(rendered)
	if w < 0 {
		w = 0
	}
	return rendered + style.Width(w).Render("")
}

func (m Model) orgVisibleRows() int {
	h := m.contentHeight() - 2 // table header + divider
	if m.orgFilter != "" || m.orgFiltering {
		h--
	}
	if h < 1 {
		return 1
	}
	return h
}

func (m Model) projVisibleRows() int {
	h := m.contentHeight() - 2
	if m.projFilter != "" || m.projFiltering {
		h--
	}
	if h < 1 {
		return 1
	}
	return h
}

func (m Model) wsVisibleRows() int {
	h := m.contentHeight() - 2
	if m.wsFilter != "" || m.wsFiltering {
		h--
	}
	if h < 1 {
		return 1
	}
	return h
}

func (m Model) runVisibleRows() int {
	h := m.contentHeight() - 3 // tab strip + table header + divider
	if m.runFilter != "" || m.runFiltering {
		h--
	}
	if h < 1 {
		return 1
	}
	return h
}

func (m Model) varVisibleRows() int {
	h := m.contentHeight() - 3 // tab strip + table header + divider
	if m.varFilter != "" || m.varFiltering {
		h--
	}
	if h < 1 {
		return 1
	}
	return h
}

func (m Model) cvVisibleRows() int {
	h := m.contentHeight() - 3 // tab strip + table header + divider
	if m.cvFilter != "" || m.cvFiltering {
		h--
	}
	if h < 1 {
		return 1
	}
	return h
}

func (m Model) svVisibleRows() int {
	h := m.contentHeight() - 3 // tab strip + table header + divider
	if m.svFilter != "" || m.svFiltering {
		h--
	}
	if h < 1 {
		return 1
	}
	return h
}

// pad fills a rendered string to the full terminal width using the given style.
func (m Model) pad(rendered string, style lipgloss.Style) string {
	w := m.width - lipgloss.Width(rendered)
	if w < 0 {
		w = 0
	}
	return rendered + style.Width(w).Render("")
}

// currentCliCmd returns the equivalent tfx CLI command for the active view.
func (m Model) currentCliCmd() string {
	switch m.currentView {
	case viewOrganizations:
		return "tfx organization list"
	case viewProjects:
		return "tfx project list"
	case viewWorkspaces:
		if m.selectedProj != nil {
			return fmt.Sprintf("tfx workspace list --project-id %s", m.selectedProj.ID)
		}
		return "tfx workspace list"
	case viewRuns:
		if m.selectedWS != nil {
			return fmt.Sprintf("tfx workspace run list -n %s", m.selectedWS.Name)
		}
		return "tfx workspace run list"
	case viewVariables:
		if m.selectedWS != nil {
			return fmt.Sprintf("tfx workspace variable list -n %s", m.selectedWS.Name)
		}
		return "tfx workspace variable list"
	case viewConfigVersions:
		if m.selectedWS != nil {
			return fmt.Sprintf("tfx workspace configuration-version list -n %s", m.selectedWS.Name)
		}
		return "tfx workspace configuration-version list"
	case viewStateVersions:
		if m.selectedWS != nil {
			return fmt.Sprintf("tfx workspace state-version list -n %s", m.selectedWS.Name)
		}
		return "tfx workspace state-version list"
	case viewWorkspaceSettings:
		if m.selectedWS != nil {
			return fmt.Sprintf("tfx workspace show -n %s", m.selectedWS.Name)
		}
		return "tfx workspace show"
	case viewWorkspaceDetail:
		if m.selectedWS != nil {
			return fmt.Sprintf("tfx workspace show -n %s", m.selectedWS.Name)
		}
		return "tfx workspace show"
	case viewOrgDetail:
		if m.selectedOrg != nil {
			return fmt.Sprintf("tfx organization show -n %s", m.selectedOrg.Name)
		}
		return "tfx organization show"
	case viewProjectDetail:
		if m.selectedProj != nil {
			return fmt.Sprintf("tfx project show --project-id %s", m.selectedProj.ID)
		}
		return "tfx project show"
	case viewRunDetail:
		if m.selectedRun != nil {
			return fmt.Sprintf("tfx workspace run show --id %s", m.selectedRun.ID)
		}
		return "tfx workspace run show"
	case viewVariableDetail:
		if m.selectedWS != nil && m.selectedVar != nil {
			return fmt.Sprintf("tfx workspace variable show -n %s --key %s", m.selectedWS.Name, m.selectedVar.Key)
		}
		return "tfx workspace variable show"
	case viewStateVersionDetail:
		if m.selectedSV != nil {
			return fmt.Sprintf("tfx workspace state-version show --state-id %s", m.selectedSV.ID)
		}
		return "tfx workspace state-version show"
	case viewStateVersionViewer:
		if m.selectedSV != nil {
			return fmt.Sprintf("tfx workspace state-version download --state-id %s", m.selectedSV.ID)
		}
		return "tfx workspace state-version download"
	case viewConfigVersionViewer, viewConfigVersionFileContent:
		if m.selectedCV != nil {
			return fmt.Sprintf("tfx workspace configuration-version download --id %s", m.selectedCV.ID)
		}
		return "tfx workspace configuration-version download"
	case viewConfigVersionDetail:
		if m.selectedCV != nil {
			return fmt.Sprintf("tfx workspace configuration-version show --id %s", m.selectedCV.ID)
		}
		return "tfx workspace configuration-version show"
	default:
		return "tfx"
	}
}

// copyToClipboard writes text to the system clipboard via a platform-native command.
func copyToClipboard(text string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		cmd = exec.Command("xclip", "-selection", "clipboard")
	case "windows":
		cmd = exec.Command("clip")
	default:
		return fmt.Errorf("clipboard not supported on %s", runtime.GOOS)
	}
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

// openBrowser opens url in the system default browser.
func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return fmt.Errorf("open browser not supported on %s", runtime.GOOS)
	}
	return cmd.Start() // Start (not Run) so we don't wait for the browser to exit.
}

// applyURL copies url to the clipboard (u) or opens it in the browser (U).
// Returns the updated model with clipFeedback set. label is shown in the
// "✓ <label> URL copied" message.
func (m Model) applyURL(url, label string, openInBrowser bool) Model {
	if url == "" {
		return m
	}
	if openInBrowser {
		if err := openBrowser(url); err == nil {
			m.clipFeedback = "✓ opening in browser"
		} else {
			m.clipFeedback = "could not open browser"
		}
	} else {
		if err := copyToClipboard(url); err == nil {
			m.clipFeedback = "✓ " + label + " URL copied"
		} else {
			m.clipFeedback = "clipboard unavailable"
		}
	}
	return m
}

// hostnameURL returns the base URL for the TFE/HCP Terraform instance.
func (m Model) hostnameURL() string {
	return fmt.Sprintf("https://%s", m.hostname)
}

// orgProjectsURL returns the URL for the projects list of the current org.
func (m Model) orgProjectsURL() string {
	return fmt.Sprintf("https://%s/app/%s/projects", m.hostname, m.org)
}

// wsURL returns the HCP Terraform / TFE web URL for the currently selected workspace.
func (m Model) wsURL() string {
	if m.selectedWS == nil {
		return ""
	}
	return fmt.Sprintf("https://%s/app/%s/workspaces/%s", m.hostname, m.org, m.selectedWS.Name)
}

// wsRunsListURL returns the URL for the runs list of the current workspace.
func (m Model) wsRunsListURL() string {
	if m.selectedWS == nil {
		return ""
	}
	return fmt.Sprintf("https://%s/app/%s/workspaces/%s/runs", m.hostname, m.org, m.selectedWS.Name)
}

// wsVariablesListURL returns the URL for the variables list of the current workspace.
func (m Model) wsVariablesListURL() string {
	if m.selectedWS == nil {
		return ""
	}
	return fmt.Sprintf("https://%s/app/%s/workspaces/%s/variables", m.hostname, m.org, m.selectedWS.Name)
}

// wsSVsListURL returns the URL for the state versions list of the current workspace.
func (m Model) wsSVsListURL() string {
	if m.selectedWS == nil {
		return ""
	}
	return fmt.Sprintf("https://%s/app/%s/workspaces/%s/states", m.hostname, m.org, m.selectedWS.Name)
}

// svURL returns the URL for the currently selected state version.
func (m Model) svURL() string {
	if m.selectedSV == nil || m.selectedWS == nil {
		return ""
	}
	return fmt.Sprintf("https://%s/app/%s/workspaces/%s/states/%s", m.hostname, m.org, m.selectedWS.Name, m.selectedSV.ID)
}

// orgURL returns the HCP Terraform / TFE web URL for the currently selected org's settings.
func (m Model) orgURL() string {
	if m.selectedOrg == nil {
		return ""
	}
	return fmt.Sprintf("https://%s/app/%s/settings/general", m.hostname, m.selectedOrg.Name)
}

// projURL returns the HCP Terraform / TFE web URL for the currently selected project's workspaces.
func (m Model) projURL() string {
	if m.selectedProj == nil {
		return ""
	}
	return fmt.Sprintf("https://%s/app/%s/projects/%s/workspaces", m.hostname, m.org, m.selectedProj.ID)
}

// runURL returns the HCP Terraform / TFE web URL for the currently selected run.
func (m Model) runURL() string {
	if m.selectedRun == nil || m.selectedWS == nil {
		return ""
	}
	return fmt.Sprintf("https://%s/app/%s/workspaces/%s/runs/%s", m.hostname, m.org, m.selectedWS.Name, m.selectedRun.ID)
}

// ── Content routing ───────────────────────────────────────────────────────────

func (m Model) renderContent() string {
	if m.loading {
		// Workspace sub-views show the tab strip even while loading so the
		// user can see which tab they switched to.
		if m.isWorkspaceSubView() {
			return m.renderWorkspaceDetailLoading()
		}
		return m.renderLoadingContent()
	}
	if m.errMsg != "" {
		return m.renderErrorContent()
	}
	switch m.currentView {
	case viewOrganizations:
		return m.renderOrgsContent()
	case viewProjects:
		return m.renderProjectsContent()
	case viewWorkspaces:
		return m.renderWorkspacesContent()
	case viewWorkspaceSettings:
		return m.renderWorkspaceSettingsContent()
	case viewRuns:
		return m.renderRunsContent()
	case viewVariables:
		return m.renderVariablesContent()
	case viewConfigVersions:
		return m.renderConfigVersionsContent()
	case viewStateVersions:
		return m.renderStateVersionsContent()
	case viewWorkspaceDetail:
		return m.renderWorkspaceDetailContent()
	case viewOrgDetail:
		return m.renderOrgDetailContent()
	case viewProjectDetail:
		return m.renderProjectDetailContent()
	case viewRunDetail:
		return m.renderRunDetailContent()
	case viewVariableDetail:
		return m.renderVariableDetailContent()
	case viewStateVersionDetail:
		return m.renderStateVersionDetailContent()
	case viewConfigVersionDetail:
		return m.renderConfigVersionDetailContent()
	case viewStateVersionViewer:
		return m.renderStateVersionJsonContent()
	case viewConfigVersionViewer:
		return m.renderConfigVersionFilesContent()
	case viewConfigVersionFileContent:
		return m.renderConfigVersionFileContent()
	}
	return m.renderLoadingContent()
}

// renderWorkspaceTabStrip renders the horizontal tab strip for workspace sub-views.
func (m Model) renderWorkspaceTabStrip() string {
	var sb strings.Builder
	sb.WriteString(tabBarStyle.Render(" "))
	for i, t := range wsTabs {
		if i > 0 {
			sb.WriteString(tabBarStyle.Render("  "))
		}
		if m.currentView == t.view {
			sb.WriteString(tabActiveStyle.Render(t.label))
		} else {
			sb.WriteString(tabInactiveStyle.Render(t.label))
		}
	}
	return m.padContent(sb.String(), tabBarStyle)
}

// renderWorkspaceDetailLoading renders the tab strip + spinner for workspace
// sub-views that are in a loading state (e.g., after switching tabs).
func (m Model) renderWorkspaceDetailLoading() string {
	h := m.contentHeight()
	lines := make([]string, 0, h)
	lines = append(lines, m.renderWorkspaceTabStrip())

	frame := spinnerFrames[m.spinnerIdx]
	mid := (h - 1) / 2
	for i := 0; i < h-1; i++ {
		if i == mid {
			lines = append(lines, contentPlaceholderStyle.Width(m.innerWidth()).Render("  "+frame+"  Loading…"))
		} else {
			lines = append(lines, contentStyle.Width(m.innerWidth()).Render(""))
		}
	}
	return strings.Join(lines, "\n")
}

func (m Model) renderLoadingContent() string {
	h := m.contentHeight()
	lines := make([]string, h)
	mid := h / 2
	frame := spinnerFrames[m.spinnerIdx]
	for i := range lines {
		if i == mid {
			lines[i] = contentPlaceholderStyle.Width(m.innerWidth()).Render("  " + frame + "  Loading…")
		} else {
			lines[i] = contentStyle.Width(m.innerWidth()).Render("")
		}
	}
	return strings.Join(lines, "\n")
}

func (m Model) renderErrorContent() string {
	h := m.contentHeight()
	lines := make([]string, h)
	mid := h / 2
	for i := range lines {
		if i == mid {
			lines[i] = statusErrorStyle.Width(m.innerWidth()).Render(fmt.Sprintf("  ✗  %s", m.errMsg))
		} else {
			lines[i] = contentStyle.Width(m.innerWidth()).Render("")
		}
	}
	return strings.Join(lines, "\n")
}

// ── Fixed chrome ──────────────────────────────────────────────────────────────

func (m Model) renderHeader() string {
	app := headerAppStyle.Render(" TFx ")
	info := headerInfoStyle.Render(" " + m.hostname)
	ver := headerVersionStyle.Render(fmt.Sprintf(" v%s ", version.Version))

	// Remote app name + TFE version — populated from ping response headers on
	// client init, so no extra API call is needed. Empty for HCP Terraform
	// (no version) or if the client isn't yet initialized.
	remoteInfo := ""
	if m.c != nil {
		if appName := m.c.Client.AppName(); appName != "" {
			remoteInfo = headerRemoteStyle.Render(" ⬥  " + appName + " ")
		}
	}

	used := lipgloss.Width(app) + lipgloss.Width(info) + lipgloss.Width(remoteInfo) + lipgloss.Width(ver)
	gap := m.width - used
	if gap < 0 {
		gap = 0
	}
	return app + info + remoteInfo + headerStyle.Width(gap).Render("") + ver
}

// expiresLabel formats a token expiry time as "YYYY-MM-DD (N days/hours/minutes)".
// Returns "never" for the zero value and "n/a" if already expired.
func expiresLabel(t time.Time) string {
	if t.IsZero() {
		return "never"
	}
	date := t.Format("2006-01-02")
	remaining := time.Until(t)
	if remaining <= 0 {
		return date + " (expired)"
	}
	var countdown string
	switch {
	case remaining >= 24*time.Hour:
		days := int(math.Round(remaining.Hours() / 24))
		countdown = fmt.Sprintf("%d days", days)
	case remaining >= time.Hour:
		hours := int(math.Round(remaining.Hours()))
		countdown = fmt.Sprintf("%d hours", hours)
	default:
		mins := int(math.Round(remaining.Minutes()))
		if mins < 1 {
			mins = 1
		}
		countdown = fmt.Sprintf("%d minutes", mins)
	}
	return fmt.Sprintf("%s (%s)", date, countdown)
}

// renderProfileBar renders five fixed rows beneath the main header:
//
//	Profile:            <name>
//	  type:             <token type>
//	  username:         <username>         API Version: <ver>
//	  email:            <email>            TFE Version: <ver>  (TFE only)
//	  token expiration: <expiry>
//
// All five rows are always emitted (with "…" placeholders while loading) so
// fixedLines stays constant and the layout does not shift after data arrives.
func (m Model) renderProfileBar() string {
	bg := lipgloss.NewStyle().Background(colorHeaderBg)
	lbl := lipgloss.NewStyle().Background(colorHeaderBg).Foreground(colorDim)
	value := lipgloss.NewStyle().Background(colorHeaderBg).Foreground(colorAccent)
	dim := lipgloss.NewStyle().Background(colorHeaderBg).Foreground(colorDim)
	na := lipgloss.NewStyle().Background(colorHeaderBg).Foreground(colorDim)

	// splitAt is the column where the right column begins.
	splitAt := m.width / 2
	if splitAt < 44 {
		splitAt = 44
	}

	// splitRow renders a row with indent=4, labelW=18, and an optional right
	// column that starts at splitAt. rk=="" suppresses the right column.
	splitRow := func(lk, lv string, lvs lipgloss.Style, rk, rv string, rvs lipgloss.Style) string {
		const labelW = 18
		leftLabel := fmt.Sprintf("    %-*s", labelW, lk+":")
		ll := lbl.Render(leftLabel)
		lval := lvs.Render(lv)
		leftUsed := lipgloss.Width(ll) + lipgloss.Width(lval)

		if rk == "" {
			gap := m.width - leftUsed
			if gap < 0 {
				gap = 0
			}
			return ll + lval + bg.Width(gap).Render("")
		}

		// Pad left side to splitAt, then render right label+value.
		leftPad := splitAt - leftUsed
		if leftPad < 1 {
			leftPad = 1
		}
		const rLabelW = 13
		rightLabel := fmt.Sprintf("%-*s", rLabelW, rk+":")
		rl := lbl.Render(rightLabel)
		rval := rvs.Render(rv)
		rightUsed := lipgloss.Width(rl) + lipgloss.Width(rval)
		rightPad := m.width - splitAt - rightUsed
		if rightPad < 0 {
			rightPad = 0
		}
		return ll + lval + bg.Width(leftPad).Render("") + rl + rval + bg.Width(rightPad).Render("")
	}

	// ── Left-column values ─────────────────────────────────────────────────

	profName := m.profileName
	if profName == "" {
		profName = "default"
	}
	profValue := value.Render(profName)
	if m.configFile != "" {
		profValue += dim.Render("  ("+m.configFile+")")
	}

	tokenType, ttStyle := "…", dim
	if m.accountTokenType != accountResourceTypeUnknown {
		switch m.accountTokenType {
		case accountResourceTypeUser:
			tokenType = "User Token"
		case accountResourceTypeTeam:
			tokenType = "Team Token"
		}
		ttStyle = value
	}

	uname, email, expires := "…", "…", "…"
	uStyle, eStyle, xStyle := dim, dim, dim
	if m.accountUser != nil {
		uname = m.accountUser.Username
		uStyle = value
		if m.accountUser.Email != "" {
			email = m.accountUser.Email
			eStyle = value
		} else {
			email = "—"
			eStyle = na
		}
	}
	if m.accountToken != nil {
		expires = expiresLabel(m.accountToken.ExpiredAt)
		xStyle = value
	} else if m.accountUser != nil {
		expires = "n/a"
		xStyle = na
	}

	// ── Right-column values ────────────────────────────────────────────────

	apiVer, tfeVersion := "", ""
	if m.c != nil {
		cl := m.c.Client
		apiVer = cl.RemoteAPIVersion()
		tfeVersion = cl.RemoteTFENumericVersion()
	}

	// Only show the TFE version label when running against TFE (empty on HCP Terraform).
	tfeVersionKey := ""
	if tfeVersion != "" {
		tfeVersionKey = "TFE Version"
	}

	// Profile row is inlined: name in accent, config path in dim.
	const profLabelW = 9
	profRow := func() string {
		padded := fmt.Sprintf("  %-*s", profLabelW, "Profile:")
		l := lbl.Render(padded)
		used := lipgloss.Width(l) + lipgloss.Width(profValue)
		gap := m.width - used
		if gap < 0 {
			gap = 0
		}
		return l + profValue + bg.Width(gap).Render("")
	}

	return strings.Join([]string{
		profRow(),
		splitRow("type", tokenType, ttStyle, "", "", dim),
		splitRow("username", uname, uStyle, "API Version", apiVer, value),
		splitRow("email", email, eStyle, tfeVersionKey, tfeVersion, value),
		splitRow("token expiration", expires, xStyle, "", "", dim),
	}, "\n")
}

// breadcrumbLine returns the styled navigation trail without any full-width
// padding. Used both by renderBreadcrumb (legacy, kept for the help overlay)
// and by renderContentTopBorder as the box title.
func (m Model) breadcrumbLine() string {
	sep := breadcrumbSepStyle.Render("  /  ")
	orgPart := breadcrumbBarStyle.Render(fmt.Sprintf(" org: %s", m.org))

	projName := ""
	if m.selectedProj != nil {
		projName = m.selectedProj.Name
	}
	wsName := ""
	if m.selectedWS != nil {
		wsName = m.selectedWS.Name
	}

	switch m.currentView {
	case viewOrganizations:
		return breadcrumbActiveStyle.Render(" organizations ")
	case viewProjects:
		return orgPart + sep + breadcrumbActiveStyle.Render("projects ")
	case viewWorkspaces:
		return orgPart + sep +
			breadcrumbBarStyle.Render(fmt.Sprintf("project: %s", projName)) +
			sep + breadcrumbActiveStyle.Render("workspaces ")
	case viewWorkspaceSettings:
		return orgPart + sep +
			breadcrumbBarStyle.Render(fmt.Sprintf("project: %s", projName)) +
			sep + breadcrumbBarStyle.Render(fmt.Sprintf("workspace: %s", wsName)) +
			sep + breadcrumbActiveStyle.Render("settings ")
	case viewRuns:
		return orgPart + sep +
			breadcrumbBarStyle.Render(fmt.Sprintf("project: %s", projName)) +
			sep + breadcrumbBarStyle.Render(fmt.Sprintf("workspace: %s", wsName)) +
			sep + breadcrumbActiveStyle.Render("runs ")
	case viewVariables:
		return orgPart + sep +
			breadcrumbBarStyle.Render(fmt.Sprintf("project: %s", projName)) +
			sep + breadcrumbBarStyle.Render(fmt.Sprintf("workspace: %s", wsName)) +
			sep + breadcrumbActiveStyle.Render("variables ")
	case viewConfigVersions:
		return orgPart + sep +
			breadcrumbBarStyle.Render(fmt.Sprintf("project: %s", projName)) +
			sep + breadcrumbBarStyle.Render(fmt.Sprintf("workspace: %s", wsName)) +
			sep + breadcrumbActiveStyle.Render("config versions ")
	case viewStateVersions:
		return orgPart + sep +
			breadcrumbBarStyle.Render(fmt.Sprintf("project: %s", projName)) +
			sep + breadcrumbBarStyle.Render(fmt.Sprintf("workspace: %s", wsName)) +
			sep + breadcrumbActiveStyle.Render("state versions ")
	case viewWorkspaceDetail:
		return orgPart + sep +
			breadcrumbBarStyle.Render(fmt.Sprintf("project: %s", projName)) +
			sep + breadcrumbBarStyle.Render(fmt.Sprintf("workspace: %s", wsName)) +
			sep + breadcrumbActiveStyle.Render("detail ")
	case viewOrgDetail:
		orgDetailName := ""
		if m.selectedOrg != nil {
			orgDetailName = m.selectedOrg.Name
		}
		return breadcrumbBarStyle.Render(fmt.Sprintf(" org: %s", orgDetailName)) +
			sep + breadcrumbActiveStyle.Render("detail ")
	case viewProjectDetail:
		return orgPart + sep +
			breadcrumbBarStyle.Render(fmt.Sprintf("project: %s", projName)) +
			sep + breadcrumbActiveStyle.Render("detail ")
	case viewRunDetail:
		runID := ""
		if m.selectedRun != nil {
			runID = m.selectedRun.ID
		}
		return orgPart + sep +
			breadcrumbBarStyle.Render(fmt.Sprintf("project: %s", projName)) +
			sep + breadcrumbBarStyle.Render(fmt.Sprintf("workspace: %s", wsName)) +
			sep + breadcrumbBarStyle.Render("runs") +
			sep + breadcrumbActiveStyle.Render(fmt.Sprintf("run: %s ", runID))
	case viewVariableDetail:
		varKey := ""
		if m.selectedVar != nil {
			varKey = m.selectedVar.Key
		}
		return orgPart + sep +
			breadcrumbBarStyle.Render(fmt.Sprintf("project: %s", projName)) +
			sep + breadcrumbBarStyle.Render(fmt.Sprintf("workspace: %s", wsName)) +
			sep + breadcrumbBarStyle.Render("variables") +
			sep + breadcrumbActiveStyle.Render(fmt.Sprintf("var: %s ", varKey))
	case viewStateVersionDetail:
		svSerial := ""
		if m.selectedSV != nil {
			svSerial = fmt.Sprintf("%d", m.selectedSV.Serial)
		}
		return orgPart + sep +
			breadcrumbBarStyle.Render(fmt.Sprintf("project: %s", projName)) +
			sep + breadcrumbBarStyle.Render(fmt.Sprintf("workspace: %s", wsName)) +
			sep + breadcrumbBarStyle.Render("state versions") +
			sep + breadcrumbActiveStyle.Render(fmt.Sprintf("sv: %s ", svSerial))
	case viewConfigVersionDetail:
		cvID := ""
		if m.selectedCV != nil {
			cvID = m.selectedCV.ID
		}
		return orgPart + sep +
			breadcrumbBarStyle.Render(fmt.Sprintf("project: %s", projName)) +
			sep + breadcrumbBarStyle.Render(fmt.Sprintf("workspace: %s", wsName)) +
			sep + breadcrumbBarStyle.Render("config versions") +
			sep + breadcrumbActiveStyle.Render(fmt.Sprintf("cv: %s ", cvID))
	case viewStateVersionViewer:
		svSerial := ""
		if m.selectedSV != nil {
			svSerial = fmt.Sprintf("%d", m.selectedSV.Serial)
		}
		return orgPart + sep +
			breadcrumbBarStyle.Render(fmt.Sprintf("project: %s", projName)) +
			sep + breadcrumbBarStyle.Render(fmt.Sprintf("workspace: %s", wsName)) +
			sep + breadcrumbBarStyle.Render("state versions") +
			sep + breadcrumbBarStyle.Render(fmt.Sprintf("sv: %s", svSerial)) +
			sep + breadcrumbActiveStyle.Render("json ")
	case viewConfigVersionViewer:
		cvID := ""
		if m.selectedCV != nil {
			cvID = m.selectedCV.ID
		}
		return orgPart + sep +
			breadcrumbBarStyle.Render(fmt.Sprintf("project: %s", projName)) +
			sep + breadcrumbBarStyle.Render(fmt.Sprintf("workspace: %s", wsName)) +
			sep + breadcrumbBarStyle.Render("config versions") +
			sep + breadcrumbBarStyle.Render(fmt.Sprintf("cv: %s", cvID)) +
			sep + breadcrumbActiveStyle.Render("files ")
	case viewConfigVersionFileContent:
		cvID := ""
		if m.selectedCV != nil {
			cvID = m.selectedCV.ID
		}
		return orgPart + sep +
			breadcrumbBarStyle.Render(fmt.Sprintf("project: %s", projName)) +
			sep + breadcrumbBarStyle.Render(fmt.Sprintf("workspace: %s", wsName)) +
			sep + breadcrumbBarStyle.Render("config versions") +
			sep + breadcrumbBarStyle.Render(fmt.Sprintf("cv: %s", cvID)) +
			sep + breadcrumbBarStyle.Render("files") +
			sep + breadcrumbActiveStyle.Render(m.cvFileName+" ")
	default:
		return orgPart
	}
}

// renderContentBox wraps the inner content with a full border (top, sides,
// bottom) giving the table area a k9s-style framed appearance.
// Content lines (from renderContent) are each flanked by │ side borders.
func (m Model) renderContentBox() string {
	b := contentBoxBorderStyle
	lb := b.Render("│")
	rb := b.Render("│")

	inner := m.renderContent()
	lines := strings.Split(inner, "\n")
	wrapped := make([]string, len(lines))
	for i, line := range lines {
		wrapped[i] = lb + line + rb
	}
	return strings.Join([]string{
		m.renderContentTopBorder(),
		strings.Join(wrapped, "\n"),
		m.renderContentBottomBorder(),
	}, "\n")
}

// renderContentTopBorder draws the top edge of the content box.
// The navigation breadcrumb is embedded as the title:
//
//	┌─ org: my-org / workspaces ─────────────────────────────────┐
func (m Model) renderContentTopBorder() string {
	b := contentBoxBorderStyle

	title := m.breadcrumbLine()
	corner := b.Render("┌")
	prefix := b.Render("─ ")

	// Measure consumed width, compute fill for the right side.
	// mainWidth() is the full outer width including corners.
	used := lipgloss.Width(corner) + lipgloss.Width(prefix) + lipgloss.Width(title) + 2 // 2 = " ┐" at end
	fill := m.mainWidth() - used
	if fill < 0 {
		fill = 0
	}
	return corner + prefix + title + b.Render(" "+strings.Repeat("─", fill)+"┐")
}

// renderContentBottomBorder draws the bottom edge of the content box:
//
//	└────────────────────────────────────────────────────────────┘
func (m Model) renderContentBottomBorder() string {
	b := contentBoxBorderStyle
	inner := m.mainWidth() - 2 // leave room for └ and ┘
	if inner < 0 {
		inner = 0
	}
	return b.Render("└" + strings.Repeat("─", inner) + "┘")
}

// ── Split content box (API Inspector visible) ────────────────────────────────

// renderSplitContentBox renders the content area as two side-by-side panels
// inside a unified border with ┬/┴ connectors, replacing the old
// joinPanels(renderContentBox(), renderDebugPanel()) approach.
//
//	┌─ breadcrumb ─────────┬── API Inspector ──────┐
//	│ left content         │ right content          │
//	└──────────────────────┴───────────────────────┘
func (m Model) renderSplitContentBox() string {
	b := contentBoxBorderStyle
	leftW := m.innerWidth()
	rightW := m.debugPanelWidth()
	h := m.contentHeight()
	ds := m.newDebugStyles()

	leftContent := m.renderContent()
	rightContent := m.renderDebugPanel()

	leftLines := strings.Split(leftContent, "\n")
	rightLines := strings.Split(rightContent, "\n")

	lb := b.Render("│")
	mb := b.Render("│")
	rb := b.Render("│")

	rows := make([]string, h)
	for i := 0; i < h; i++ {
		var l, r string
		if i < len(leftLines) {
			l = leftLines[i]
		}
		if i < len(rightLines) {
			r = rightLines[i]
		}

		// Enforce exact widths — truncate or pad.
		if w := lipgloss.Width(l); w > leftW {
			l = ansi.Truncate(l, leftW, "")
		} else if w < leftW {
			l += contentStyle.Width(leftW - w).Render("")
		}
		if w := lipgloss.Width(r); w > rightW {
			r = ansi.Truncate(r, rightW, "")
		} else if w < rightW {
			r += ds.bg.Width(rightW - w).Render("")
		}

		rows[i] = lb + l + mb + r + rb
	}

	return strings.Join([]string{
		m.renderSplitTopBorder(),
		strings.Join(rows, "\n"),
		m.renderSplitBottomBorder(),
	}, "\n")
}

// renderSplitTopBorder draws the top edge of the split content box:
//
//	┌─ breadcrumb ─────────┬── API Inspector ──────┐
func (m Model) renderSplitTopBorder() string {
	b := contentBoxBorderStyle
	lw := m.innerWidth()
	rw := m.debugPanelWidth()

	// ── Left half: ┌─ breadcrumb ─────── (lw chars after ┌) ──
	// NOTE: use lipgloss.Width for display-width of box-drawing chars, not len().
	breadcrumb := m.breadcrumbLine()
	const prefixStr = "─ "
	const suffixStr = " "
	prefixW := lipgloss.Width(prefixStr)
	suffixW := lipgloss.Width(suffixStr)
	used := prefixW + lipgloss.Width(breadcrumb) + suffixW
	fill := lw - used
	if fill < 0 {
		// Breadcrumb too wide — truncate it.
		breadcrumb = ansi.Truncate(breadcrumb, lw-prefixW-suffixW-1, "")
		fill = lw - prefixW - lipgloss.Width(breadcrumb) - suffixW
		if fill < 0 {
			fill = 0
		}
	}
	leftPart := b.Render("┌"+prefixStr) + breadcrumb + b.Render(suffixStr+strings.Repeat("─", fill))

	// ── Right half: ┬── API Inspector ──────┐ (rw chars between ┬ and ┐) ──
	titleStyle := debugTitleUnfocusedStyle
	if m.debugFocused {
		titleStyle = debugTitleFocusedStyle
	}
	inspectorTitle := titleStyle.Render("API Inspector")
	const rPrefixStr = "── "
	const rSuffixStr = " "
	rPrefixW := lipgloss.Width(rPrefixStr)
	rSuffixW := lipgloss.Width(rSuffixStr)
	rUsed := rPrefixW + lipgloss.Width(inspectorTitle) + rSuffixW
	rFill := rw - rUsed
	if rFill < 0 {
		rFill = 0
	}
	rightPart := b.Render("┬"+rPrefixStr) + inspectorTitle + b.Render(rSuffixStr+strings.Repeat("─", rFill)+"┐")

	return leftPart + rightPart
}

// renderSplitBottomBorder draws the bottom edge of the split content box:
//
//	└──────────────────────┴───────────────────────┘
func (m Model) renderSplitBottomBorder() string {
	b := contentBoxBorderStyle
	return b.Render("└" + strings.Repeat("─", m.innerWidth()) + "┴" + strings.Repeat("─", m.debugPanelWidth()) + "┘")
}

// renderBreadcrumb is kept for the help overlay (which uses m.renderHeader()
// directly and builds its own layout). Normal views use renderContentBox.
func (m Model) renderBreadcrumb() string {
	return m.pad(m.breadcrumbLine(), breadcrumbBarStyle)
}

func (m Model) renderStatusBar() string {
	if m.loading {
		frame := spinnerFrames[m.spinnerIdx]
		return m.pad(statusLoadingStyle.Render("  "+frame+"  Loading…"), statusLoadingStyle)
	}
	if m.errMsg != "" {
		return m.pad(statusErrorStyle.Render(fmt.Sprintf("  ✗  %s", m.errMsg)), statusErrorStyle)
	}
	if m.clipFeedback != "" {
		return m.pad(statusSuccessStyle.Render("  "+m.clipFeedback), statusSuccessStyle)
	}

	var msg string
	switch m.currentView {
	case viewOrganizations:
		fo := filteredOrgs(m)
		if m.orgFilter != "" {
			msg = fmt.Sprintf("  %d / %d organizations  •  filter: %s", len(fo), len(m.orgs), m.orgFilter)
		} else {
			msg = fmt.Sprintf("  %d organizations", len(m.orgs))
		}
	case viewProjects:
		fp := filteredProjects(m)
		if m.projFilter != "" {
			msg = fmt.Sprintf("  %d / %d projects  •  filter: %s", len(fp), len(m.projects), m.projFilter)
		} else {
			msg = fmt.Sprintf("  %d projects", len(m.projects))
		}
	case viewWorkspaces:
		fw := filteredWorkspaces(m)
		if m.wsFilter != "" {
			msg = fmt.Sprintf("  %d / %d workspaces  •  filter: %s", len(fw), len(m.workspaces), m.wsFilter)
		} else {
			msg = fmt.Sprintf("  %d workspaces", len(m.workspaces))
		}
	case viewRuns:
		fr := filteredRuns(m)
		if m.runFilter != "" {
			msg = fmt.Sprintf("  %d / %d runs  •  filter: %s", len(fr), len(m.runs), m.runFilter)
		} else {
			msg = fmt.Sprintf("  %d runs", len(m.runs))
		}
	case viewVariables:
		fv := filteredVariables(m)
		if m.varFilter != "" {
			msg = fmt.Sprintf("  %d / %d variables  •  filter: %s", len(fv), len(m.variables), m.varFilter)
		} else {
			msg = fmt.Sprintf("  %d variables", len(m.variables))
		}
	case viewConfigVersions:
		fc := filteredConfigVersions(m)
		if m.cvFilter != "" {
			msg = fmt.Sprintf("  %d / %d config versions  •  filter: %s", len(fc), len(m.configVersions), m.cvFilter)
		} else {
			msg = fmt.Sprintf("  %d config versions", len(m.configVersions))
		}
	case viewStateVersions:
		fs := filteredStateVersions(m)
		if m.svFilter != "" {
			msg = fmt.Sprintf("  %d / %d state versions  •  filter: %s", len(fs), len(m.stateVersions), m.svFilter)
		} else {
			msg = fmt.Sprintf("  %d state versions", len(m.stateVersions))
		}
	case viewWorkspaceDetail:
		if m.selectedWS != nil {
			msg = fmt.Sprintf("  workspace: %s  •  ↑ ↓ to scroll", m.selectedWS.Name)
		}
	case viewOrgDetail:
		if m.selectedOrg != nil {
			msg = fmt.Sprintf("  org: %s  •  ↑ ↓ to scroll", m.selectedOrg.Name)
		}
	case viewProjectDetail:
		if m.selectedProj != nil {
			msg = fmt.Sprintf("  project: %s  •  ↑ ↓ to scroll", m.selectedProj.Name)
		}
	case viewRunDetail:
		if m.selectedRun != nil {
			msg = fmt.Sprintf("  run: %s  •  ↑ ↓ to scroll", m.selectedRun.ID)
		}
	case viewVariableDetail:
		if m.selectedVar != nil {
			msg = fmt.Sprintf("  variable: %s  •  ↑ ↓ to scroll", m.selectedVar.Key)
		}
	case viewStateVersionDetail:
		if m.selectedSV != nil {
			msg = fmt.Sprintf("  state version serial: %d  •  ↑ ↓ to scroll", m.selectedSV.Serial)
		}
	case viewConfigVersionDetail:
		if m.selectedCV != nil {
			msg = fmt.Sprintf("  config version: %s  •  ↑ ↓ to scroll", m.selectedCV.ID)
		}
	case viewStateVersionViewer:
		if m.svJsonLoading {
			frame := spinnerFrames[m.spinnerIdx]
			return m.pad(statusLoadingStyle.Render("  "+frame+"  Loading state JSON…"), statusLoadingStyle)
		}
		if m.svJsonErr != "" {
			return m.pad(statusErrorStyle.Render(fmt.Sprintf("  ✗  %s", m.svJsonErr)), statusErrorStyle)
		}
		numLines := len(m.svJsonLines)
		cur := m.svJsonScroll + 1
		if cur > numLines {
			cur = numLines
		}
		totalBytes := 0
		for _, l := range m.svJsonLines {
			totalBytes += len(l) + 1
		}
		var sizeStr string
		if totalBytes < 1024 {
			sizeStr = fmt.Sprintf("%d B", totalBytes)
		} else {
			sizeStr = fmt.Sprintf("%d KB", totalBytes/1024)
		}
		msg = fmt.Sprintf("  state JSON  •  line %d of %d  (%s)", cur, numLines, sizeStr)
	case viewConfigVersionViewer:
		if m.cvFileLoading {
			frame := spinnerFrames[m.spinnerIdx]
			return m.pad(statusLoadingStyle.Render("  "+frame+"  Downloading config version archive…"), statusLoadingStyle)
		}
		if m.cvFileErr != "" {
			return m.pad(statusErrorStyle.Render(fmt.Sprintf("  ✗  %s", m.cvFileErr)), statusErrorStyle)
		}
		// Build status bar with OSC 8 hyperlink for the disk path so the user can
		// Cmd+Click (macOS) / Ctrl+Click (Linux/Windows) to open the directory in
		// Finder or the native file manager.  Apply lipgloss styling first, then
		// wrap in the OSC 8 bytes — never pass raw OSC 8 through lipgloss Render.
		prefix := statusBarStyle.Render("  config version files  •  ")
		suffix := statusBarStyle.Render(fmt.Sprintf("  •  %d files", len(m.cvFiles)))
		var pathPart string
		if m.selectedCV != nil {
			absPath := cvExtractDirPath(m.selectedCV.ID)
			styledPath := statusBarStyle.Render(tildePath(absPath))
			pathPart = osc8FileLink(absPath, styledPath)
		}
		return m.pad(prefix+pathPart+suffix, statusBarStyle)
	case viewConfigVersionFileContent:
		numLines := len(m.cvFileLines)
		cur := m.cvFileScroll + 1
		if cur > numLines {
			cur = numLines
		}
		totalBytes := 0
		for _, l := range m.cvFileLines {
			totalBytes += len(l) + 1
		}
		var sizeStr string
		if totalBytes < 1024 {
			sizeStr = fmt.Sprintf("%d B", totalBytes)
		} else {
			sizeStr = fmt.Sprintf("%d KB", totalBytes/1024)
		}
		name := m.cvFileName
		if name == "" {
			name = "file"
		}
		msg = fmt.Sprintf("  %s  •  line %d of %d  (%s)", name, cur, numLines, sizeStr)
	default:
		msg = "  Ready"
	}

	// When the API inspector panel is focused, append a right-aligned badge.
	if m.debugFocused {
		label := "  [api inspector]  "
		if m.debugDetailMode {
			label = "  [api inspector › detail]  "
		}
		badge := statusInspectorStyle.Render(label)
		left := statusBarStyle.Render(msg)
		gap := m.width - lipgloss.Width(left) - lipgloss.Width(badge)
		if gap < 0 {
			gap = 0
		}
		return left + statusBarStyle.Width(gap).Render("") + badge
	}
	return m.pad(statusBarStyle.Render(msg), statusBarStyle)
}

func (m Model) renderCliHint() string {
	label := cliHintBarStyle.Render("  cmd: ")
	cmd := cliHintCmdStyle.Render(m.currentCliCmd())

	var hints string
	switch {
	case m.debugFocused && m.showDebug && m.debugDetailMode:
		hints = cliHintBarStyle.Render("   •   ↑ ↓ scroll   shift+↑↓ page   •   c copy response   •   C copy curl   •   esc back to list   •   tab unfocus")
	case m.debugFocused && m.showDebug:
		hints = cliHintBarStyle.Render("   •   ↑ ↓ navigate   •   enter detail   /  filter   •   tab unfocus")
	case m.currentView == viewOrganizations:
		hints = cliHintBarStyle.Render("   •   enter projects   d detail   •   u url   U browser   •   c copy tfx cmd   •   ? help   •   q quit")
	case m.currentView == viewProjects:
		hints = cliHintBarStyle.Render("   •   enter workspaces   d detail   •   u url   U browser   •   c copy tfx cmd   •   ? help   •   q quit")
	case m.currentView == viewWorkspaces:
		hints = cliHintBarStyle.Render("   •   enter ws   v vars   f cvs   s svs   d detail   •   u url   U browser   •   c copy tfx cmd   •   ? help   •   q quit")
	case m.currentView == viewOrgDetail, m.currentView == viewProjectDetail, m.currentView == viewWorkspaceDetail:
		hints = cliHintBarStyle.Render("   •   ↑ ↓ scroll   •   u url   U browser   •   ? help   •   q quit")
	case m.currentView == viewRunDetail:
		hints = cliHintBarStyle.Render("   •   ↑ ↓ scroll   •   u url   U browser   •   ? help   •   q quit")
	case m.currentView == viewVariableDetail:
		hints = cliHintBarStyle.Render("   •   ↑ ↓ scroll   •   u url   U browser   •   ? help   •   q quit")
	case m.currentView == viewConfigVersions:
		hints = cliHintBarStyle.Render("   •   enter viewer   d detail   •   ← → switch tabs   •   c copy tfx cmd   •   ? help   •   q quit")
	case m.currentView == viewStateVersions:
		hints = cliHintBarStyle.Render("   •   enter viewer   d detail   •   ← → switch tabs   •   u url   U browser   •   c copy tfx cmd   •   ? help   •   q quit")
	case m.currentView == viewStateVersionDetail:
		hints = cliHintBarStyle.Render("   •   ↑ ↓ scroll   •   o state viewer   •   u url   U browser   •   ? help   •   q quit")
	case m.currentView == viewStateVersionViewer:
		hints = cliHintBarStyle.Render("   •   ↑ ↓ scroll   •   shift+↑↓ half page   •   r re-fetch   •   ? help   •   q quit")
	case m.currentView == viewConfigVersionDetail:
		hints = cliHintBarStyle.Render("   •   ↑ ↓ scroll   •   o config viewer   •   ? help   •   q quit")
	case m.currentView == viewConfigVersionViewer:
		hints = cliHintBarStyle.Render("   •   ↑ ↓ navigate   •   enter open   •   p copy path   •   r re-fetch   •   ? help   •   q quit")
	case m.currentView == viewConfigVersionFileContent:
		hints = cliHintBarStyle.Render("   •   ↑ ↓ scroll   •   shift+↑↓ half page   •   ? help   •   q quit")
	case m.isWorkspaceSubView():
		hints = cliHintBarStyle.Render("   •   enter detail   •   ← → switch tabs   •   u url   U browser   •   c copy tfx cmd   •   ? help   •   q quit")
	default:
		hints = cliHintBarStyle.Render("   •   c copy tfx cmd   •   ? help   •   q quit")
	}
	return m.pad(label+cmd+hints, cliHintBarStyle)
}

// ── Help overlay ──────────────────────────────────────────────────────────────

type helpBinding struct {
	key  string
	desc string
}

type helpSection struct {
	title    string
	bindings []helpBinding
}

// helpSections returns context-aware shortcut sections based on the current view.
func (m Model) helpSections() []helpSection {
	// Navigation — always shown.
	nav := helpSection{
		title: "Navigation",
		bindings: []helpBinding{
			{"↑ / k", "move up"},
			{"↓ / j", "move down"},
			{"enter", "select / drill in"},
			{"esc", "go back / clear filter"},
			{"g / shift+g", "jump to top / bottom"},
		},
	}

	// View-specific section.
	var view helpSection
	switch {
	case m.debugFocused && m.showDebug && m.debugDetailMode:
		view = helpSection{
			title: "API Inspector Detail",
			bindings: []helpBinding{
				{"↑ / ↓", "scroll one line"},
				{"⇧↑ / ⇧↓", "scroll full page"},
				{"^u / ^d", "scroll half-page"},
				{"c", "copy response body"},
				{"shift+c", "copy curl command"},
				{"esc", "back to call list"},
				{"tab", "switch to left panel"},
			},
		}
	case m.debugFocused && m.showDebug:
		view = helpSection{
			title: "API Inspector",
			bindings: []helpBinding{
				{"↑ / ↓", "navigate call list"},
				{"enter", "open request detail"},
				{"/", "filter calls"},
				{"esc", "clear filter / back"},
				{"tab", "switch to left panel"},
			},
		}
	case m.currentView == viewOrganizations:
		view = helpSection{
			title: "Organizations",
			bindings: []helpBinding{
				{"enter", "view projects"},
				{"d", "view org detail"},
				{"u", "copy URL"},
				{"shift+u", "open in browser"},
			},
		}
	case m.currentView == viewProjects:
		view = helpSection{
			title: "Projects",
			bindings: []helpBinding{
				{"enter", "view workspaces"},
				{"d", "view project detail"},
				{"u", "copy URL"},
				{"shift+u", "open in browser"},
			},
		}
	case m.currentView == viewWorkspaces:
		view = helpSection{
			title: "Workspaces",
			bindings: []helpBinding{
				{"enter", "view runs tab"},
				{"v", "view variables tab"},
				{"f", "view config versions tab"},
				{"s", "view state versions tab"},
				{"d", "view workspace detail"},
				{"u", "copy URL"},
				{"shift+u", "open in browser"},
			},
		}
	case m.currentView == viewConfigVersions:
		view = helpSection{
			title: "Config Versions",
			bindings: []helpBinding{
				{"enter", "open viewer"},
				{"d", "view detail"},
				{"← →", "switch tabs"},
			},
		}
	case m.currentView == viewStateVersions:
		view = helpSection{
			title: "State Versions",
			bindings: []helpBinding{
				{"enter", "open viewer"},
				{"d", "view detail"},
				{"← →", "switch tabs"},
				{"u", "copy URL"},
				{"shift+u", "open in browser"},
			},
		}
	case m.isWorkspaceSubView():
		view = helpSection{
			title: "Workspace Tabs",
			bindings: []helpBinding{
				{"enter", "view detail"},
				{"← →", "switch tabs"},
				{"u", "copy URL"},
				{"shift+u", "open in browser"},
			},
		}
	case m.currentView == viewStateVersionDetail:
		view = helpSection{
			title: "State Version Detail",
			bindings: []helpBinding{
				{"↑ ↓", "scroll"},
				{"o", "open state viewer"},
				{"u", "copy URL"},
				{"shift+u", "open in browser"},
			},
		}
	case m.currentView == viewConfigVersionDetail:
		view = helpSection{
			title: "Config Version Detail",
			bindings: []helpBinding{
				{"↑ ↓", "scroll"},
				{"o", "open config viewer"},
			},
		}
	case m.currentView == viewStateVersionViewer:
		view = helpSection{
			title: "State Viewer",
			bindings: []helpBinding{
				{"↑ ↓", "scroll one line"},
				{"⇧↑ / ⇧↓", "scroll half page"},
				{"r", "re-fetch state JSON"},
			},
		}
	case m.currentView == viewConfigVersionViewer:
		view = helpSection{
			title: "Config Viewer",
			bindings: []helpBinding{
				{"↑ ↓", "navigate file tree"},
				{"enter", "open file"},
				{"p", "copy cache path"},
				{"r", "re-fetch files"},
			},
		}
	case m.currentView == viewConfigVersionFileContent:
		view = helpSection{
			title: "File Viewer",
			bindings: []helpBinding{
				{"↑ ↓", "scroll one line"},
				{"⇧↑ / ⇧↓", "scroll half page"},
			},
		}
	default:
		// Generic detail views (org, project, workspace, run, variable detail).
		view = helpSection{
			title: "Detail View",
			bindings: []helpBinding{
				{"↑ ↓", "scroll"},
				{"u", "copy URL"},
				{"shift+u", "open in browser"},
			},
		}
	}

	// Tools — always shown.
	tools := helpSection{
		title: "Tools",
		bindings: []helpBinding{
			{"/", "filter"},
			{"r", "refresh"},
			{"c", "copy CLI command"},
			{"i", "toggle instance info"},
			{"l", "toggle API inspector"},
			{"?", "toggle help"},
			{"q", "quit"},
		},
	}

	sections := []helpSection{nav, view, tools}

	// API inspector section — shown when the panel is open but not focused.
	if m.showDebug && !m.debugFocused {
		sections = append(sections, helpSection{
			title: "API Inspector",
			bindings: []helpBinding{
				{"tab", "switch to inspector panel"},
			},
		})
	}

	return sections
}

func (m Model) renderHelpOverlay() string {
	sections := m.helpSections()

	lines := make([]string, 0, m.height)
	lines = append(lines, m.renderHeader())
	lines = append(lines, m.pad(helpTitleStyle.Render("  Keyboard Shortcuts"), helpTitleStyle))
	lines = append(lines, helpBarStyle.Width(m.width).Render(""))

	for si, sec := range sections {
		if si > 0 {
			lines = append(lines, helpBarStyle.Width(m.width).Render(""))
		}
		title := helpBarStyle.Render("  ") + helpSectionStyle.Render(sec.title)
		lines = append(lines, m.pad(title, helpBarStyle))
		for _, b := range sec.bindings {
			key := helpKeyStyle.Width(14).Render(b.key)
			desc := helpDescStyle.Render("  " + b.desc)
			line := helpBarStyle.Render("  ") + key + desc
			lines = append(lines, m.pad(line, helpBarStyle))
		}
	}

	lines = append(lines, helpBarStyle.Width(m.width).Render(""))
	lines = append(lines, m.pad(helpBarStyle.Render("  Press ? or esc to close"), helpBarStyle))

	for len(lines) < m.height {
		lines = append(lines, helpBarStyle.Width(m.width).Render(""))
	}
	return strings.Join(lines[:m.height], "\n")
}

// ── Table rendering (shared by all list views) ────────────────────────────────

// column defines a table column's display name and character width.
type column struct {
	name  string
	width int
}

func (m Model) renderTableHeader(cols []column) string {
	parts := []string{tableHeaderStyle.Render("  ")} // cursor placeholder
	for _, col := range cols {
		parts = append(parts, tableHeaderStyle.Width(col.width).Render(strings.ToUpper(col.name)))
		parts = append(parts, tableHeaderStyle.Render("  "))
	}
	return m.padContent(strings.Join(parts, ""), tableHeaderStyle)
}

func (m Model) renderTableDivider() string {
	w := m.innerWidth()
	return tableHeaderDividerStyle.Width(w).Render(strings.Repeat("─", w))
}

func (m Model) renderTableRow(selected bool, cells []string, cols []column) string {
	style := tableRowStyle
	cursor := "  "
	if selected {
		style = tableRowSelectedStyle
		cursor = "> "
	}

	parts := []string{style.Render(cursor)}
	for i, col := range cols {
		val := ""
		if i < len(cells) {
			val = truncateStr(cells[i], col.width)
		}
		parts = append(parts, style.Width(col.width).Render(val))
		parts = append(parts, style.Render("  "))
	}
	return m.padContent(strings.Join(parts, ""), style)
}

// renderTableRowWithCellStyles renders a row where individual cells can have
// a custom foreground color. cellFgs[i] overrides the fg for column i when
// the row is not selected (selection style always takes precedence).
func (m Model) renderTableRowWithCellStyles(selected bool, cells []string, cols []column, cellFgs []color.Color) string {
	base := tableRowStyle
	cursor := "  "
	if selected {
		base = tableRowSelectedStyle
		cursor = "> "
	}

	parts := []string{base.Render(cursor)}
	for i, col := range cols {
		val := ""
		if i < len(cells) {
			val = truncateStr(cells[i], col.width)
		}
		cellStyle := base.Width(col.width)
		if !selected && i < len(cellFgs) && cellFgs[i] != nil {
			cellStyle = cellStyle.Foreground(cellFgs[i])
		}
		parts = append(parts, cellStyle.Render(val))
		parts = append(parts, base.Render("  "))
	}
	return m.padContent(strings.Join(parts, ""), base)
}

func (m Model) renderFilterBar(filter string, active bool) string {
	prompt := filterBarStyle.Render("  / ")
	var text string
	if active {
		text = filterBarActiveStyle.Render(filter + "▌")
	} else {
		text = filterBarActiveStyle.Render(filter)
	}
	return m.padContent(prompt+text, filterBarStyle)
}

// truncateStr truncates s to at most n runes, appending "…" if shortened.
func truncateStr(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	if n > 1 {
		return string(r[:n-1]) + "…"
	}
	return string(r[:n])
}

// truncateStrLeft truncates s from the LEFT to at most n runes, prepending "…"
// if shortened.  This keeps the tail of the string visible — ideal for paths
// where the deepest directory is the most meaningful part (à la superfile).
func truncateStrLeft(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	if n > 1 {
		return "…" + string(r[len(r)-(n-1):])
	}
	return string(r[len(r)-n:])
}
