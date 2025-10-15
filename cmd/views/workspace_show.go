// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package view

import (
	"fmt"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/logrusorgru/aurora"
)

// WorkspaceShowView handles rendering for workspace show command
type WorkspaceShowView struct {
	*BaseView
}

func NewWorkspaceShowView() *WorkspaceShowView {
	return &WorkspaceShowView{
		BaseView: NewBaseView(),
	}
}

// workspaceShowOutput is a JSON-safe representation of a workspace
type workspaceShowOutput struct {
	Organization       string   `json:"organization"`
	Name               string   `json:"name"`
	ID                 string   `json:"id"`
	ResourceCount      int      `json:"resourceCount"`
	TerraformVersion   string   `json:"terraformVersion"`
	ExecutionMode      string   `json:"executionMode"`
	AutoApply          bool     `json:"autoApply"`
	WorkingDirectory   string   `json:"workingDirectory"`
	Locked             bool     `json:"locked"`
	GlobalRemoteState  bool     `json:"globalRemoteState"`
	CurrentRun         *runInfo `json:"currentRun,omitempty"`
	TeamAccess         []string `json:"teamAccess,omitempty"`
	RemoteStateSharing []string `json:"remoteStateSharing,omitempty"`
}

type runInfo struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Created string `json:"created"`
}

// Render renders a single workspace's details
func (v *WorkspaceShowView) Render(
	orgName string,
	workspace *tfe.Workspace,
	currentRun *tfe.Run,
	teamNames []interface{},
	remoteStateConsumers []*tfe.Workspace,
) error {
	if workspace == nil {
		return v.RenderError(fmt.Errorf("workspace not found"))
	}

	if v.IsJSON() {
		// JSON mode: convert to JSON-safe structure
		output := workspaceShowOutput{
			Organization:      orgName,
			Name:              workspace.Name,
			ID:                workspace.ID,
			ResourceCount:     workspace.ResourceCount,
			TerraformVersion:  workspace.TerraformVersion,
			ExecutionMode:     workspace.ExecutionMode,
			AutoApply:         workspace.AutoApply,
			WorkingDirectory:  workspace.WorkingDirectory,
			Locked:            workspace.Locked,
			GlobalRemoteState: workspace.GlobalRemoteState,
		}

		// Add current run if present
		if currentRun != nil {
			output.CurrentRun = &runInfo{
				ID:      currentRun.ID,
				Status:  string(currentRun.Status),
				Created: FormatDateTime(currentRun.CreatedAt),
			}
		}

		// Add team access
		if len(teamNames) > 0 {
			output.TeamAccess = make([]string, len(teamNames))
			for i, name := range teamNames {
				output.TeamAccess[i] = fmt.Sprint(name)
			}
		}

		// Add remote state sharing
		if !workspace.GlobalRemoteState && len(remoteStateConsumers) > 0 {
			output.RemoteStateSharing = make([]string, len(remoteStateConsumers))
			for i, ws := range remoteStateConsumers {
				output.RemoteStateSharing[i] = ws.Name
			}
		}

		return v.renderer.RenderJSON(output)
	}

	// Terminal mode: render key fields in order
	properties := []PropertyPair{
		{Key: "ID", Value: workspace.ID},
		{Key: "Resource Count", Value: workspace.ResourceCount},
		{Key: "Terraform Version", Value: workspace.TerraformVersion},
		{Key: "Execution Mode", Value: workspace.ExecutionMode},
		{Key: "Auto Apply", Value: workspace.AutoApply},
		{Key: "Working Directory", Value: workspace.WorkingDirectory},
		{Key: "Locked", Value: workspace.Locked},
		{Key: "Global State Sharing", Value: workspace.GlobalRemoteState},
	}

	// Add current run info
	if currentRun == nil {
		properties = append(properties, PropertyPair{Key: "Current Run", Value: "none"})
	} else {
		properties = append(properties,
			PropertyPair{Key: "Current Run Id", Value: currentRun.ID},
			PropertyPair{Key: "Current Run Status", Value: currentRun.Status},
			PropertyPair{Key: "Current Run Created", Value: FormatDateTime(currentRun.CreatedAt)},
		)
	}

	err := v.renderer.RenderProperties(properties)
	if err != nil {
		return err
	}

	// Render team access if present
	if len(teamNames) > 0 {
		fmt.Println()
		fmt.Printf("%s\n", aurora.Bold("Team Access:"))
		for _, name := range teamNames {
			fmt.Printf("  - %s\n", name)
		}
	}

	// Render remote state sharing if present
	if !workspace.GlobalRemoteState && len(remoteStateConsumers) > 0 {
		fmt.Println()
		fmt.Printf("%s\n", aurora.Bold("Remote State Sharing Workspaces:"))
		for _, ws := range remoteStateConsumers {
			fmt.Printf("  - %s\n", ws.Name)
		}
	}

	return nil
}
