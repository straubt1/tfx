package data

import (
	"math"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/cmd/flags"
	"github.com/straubt1/tfx/logger"
)

// FetchWorkspaces fetches all workspaces for a given organization using pagination
func FetchWorkspaces(c *client.TfxClient, orgName string, options *flags.WorkspaceListFlags) ([]*tfe.Workspace, error) {
	// TODO: options to JSON
	logger.Debug("Fetching workspaces", "organization", orgName, "options", options)

	return client.FetchAll(c.Context, func(pageNumber int) ([]*tfe.Workspace, *client.Pagination, error) {
		logger.Trace("Fetching workspaces page", "organization", orgName, "page", pageNumber)

		opts := &tfe.WorkspaceListOptions{
			ListOptions: tfe.ListOptions{PageNumber: pageNumber, PageSize: 100},
			Include:     []tfe.WSIncludeOpt{"organization", "current_run"},
		}

		// Apply search and filter options if provided
		if options != nil {
			if options.Search != "" {
				opts.Search = options.Search
			}
			if options.WildcardName != "" {
				opts.WildcardName = options.WildcardName
			}
			if options.ProjectID != "" {
				opts.ProjectID = options.ProjectID
			}
			if options.Tags != "" {
				opts.Tags = options.Tags
			}
			if options.ExcludeTags != "" {
				opts.ExcludeTags = options.ExcludeTags
			}
			if options.RunStatus != "" {
				opts.CurrentRunStatus = options.RunStatus
			}
		}

		result, err := c.Client.Workspaces.List(c.Context, orgName, opts)
		if err != nil {
			logger.Error("Failed to fetch workspaces page", "organization", orgName, "page", pageNumber, "error", err)
			return nil, nil, err
		}

		logger.Trace("Workspaces page fetched", "organization", orgName, "page", pageNumber, "count", len(result.Items))
		return result.Items, client.NewPaginationFromTFE(result.Pagination), nil
	})
}

// FetchWorkspacesAcrossOrgs fetches workspaces across all organizations
func FetchWorkspacesAcrossOrgs(c *client.TfxClient, options *flags.WorkspaceListFlags) ([]*tfe.Workspace, error) {
	logger.Info("Fetching workspaces across all organizations", "options", options)

	orgs, err := FetchOrganizations(c, "")
	if err != nil {
		logger.Error("Failed to fetch organizations", "error", err)
		return nil, errors.Wrap(err, "failed to list organizations")
	}

	logger.Debug("Organizations fetched", "count", len(orgs))

	var allWorkspaces []*tfe.Workspace
	for _, org := range orgs {
		logger.Debug("Fetching workspaces for organization", "organization", org.Name)

		workspaces, err := FetchWorkspaces(c, org.Name, options)
		if err != nil {
			logger.Error("Failed to fetch workspaces", "organization", org.Name, "error", err)
			return nil, errors.Wrapf(err, "failed to list workspaces for organization %s", org.Name)
		}

		logger.Debug("Workspaces fetched for organization", "organization", org.Name, "count", len(workspaces))
		allWorkspaces = append(allWorkspaces, workspaces...)
	}

	logger.Info("All workspaces fetched successfully", "totalWorkspaces", len(allWorkspaces), "organizations", len(orgs))
	return allWorkspaces, nil
}

// FetchWorkspace fetches a single workspace by name in the specified organization
func FetchWorkspace(c *client.TfxClient, orgName string, workspaceName string) (*tfe.Workspace, error) {
	logger.Debug("Fetching workspace by name", "organization", orgName, "workspaceName", workspaceName)

	workspace, err := c.Client.Workspaces.Read(c.Context, orgName, workspaceName)
	if err != nil {
		logger.Error("Failed to fetch workspace", "organization", orgName, "workspaceName", workspaceName, "error", err)
		return nil, err
	}

	logger.Debug("Workspace fetched successfully", "organization", orgName, "workspaceName", workspaceName, "workspaceID", workspace.ID)
	return workspace, nil
}

// FetchWorkspaceRemoteStateConsumers fetches all remote state consumers for a workspace
func FetchWorkspaceRemoteStateConsumers(c *client.TfxClient, workspaceID string) ([]*tfe.Workspace, error) {
	logger.Debug("Fetching remote state consumers", "workspaceID", workspaceID)

	return client.FetchAll(c.Context, func(pageNumber int) ([]*tfe.Workspace, *client.Pagination, error) {
		logger.Trace("Fetching remote state consumers page", "workspaceID", workspaceID, "page", pageNumber)

		opts := &tfe.RemoteStateConsumersListOptions{
			ListOptions: tfe.ListOptions{PageNumber: pageNumber, PageSize: 100},
		}

		result, err := c.Client.Workspaces.ListRemoteStateConsumers(c.Context, workspaceID, opts)
		if err != nil {
			logger.Error("Failed to fetch remote state consumers page", "workspaceID", workspaceID, "page", pageNumber, "error", err)
			return nil, nil, err
		}

		logger.Trace("Remote state consumers page fetched", "workspaceID", workspaceID, "page", pageNumber, "count", len(result.Items))
		return result.Items, client.NewPaginationFromTFE(result.Pagination), nil
	})
}

// FetchWorkspaceTeamAccess fetches all team access for a workspace
func FetchWorkspaceTeamAccess(c *client.TfxClient, workspaceID string, maxItems int) ([]*tfe.TeamAccess, error) {
	logger.Debug("Fetching team access", "workspaceID", workspaceID)

	if maxItems == 0 {
		maxItems = math.MaxInt
	}

	allItems := []*tfe.TeamAccess{}
	opts := tfe.TeamAccessListOptions{
		ListOptions: tfe.ListOptions{PageNumber: 1, PageSize: 100},
		WorkspaceID: workspaceID,
	}

	for {
		logger.Trace("Fetching team access page", "workspaceID", workspaceID, "page", opts.PageNumber)

		result, err := c.Client.TeamAccess.List(c.Context, &opts)
		if err != nil {
			logger.Error("Failed to fetch team access page", "workspaceID", workspaceID, "page", opts.PageNumber, "error", err)
			return nil, err
		}

		logger.Trace("Team access page fetched", "workspaceID", workspaceID, "page", opts.PageNumber, "count", len(result.Items))
		allItems = append(allItems, result.Items...)

		if result.CurrentPage >= result.TotalPages || len(allItems) >= maxItems {
			break
		}
		opts.PageNumber = result.NextPage
	}

	logger.Debug("Team access fetched successfully", "workspaceID", workspaceID, "count", len(allItems))
	return allItems, nil
}

// FilterWorkspaces filters workspaces based on run status and repository identifier
func FilterWorkspaces(workspaces []*tfe.Workspace, runStatus string, repoIdentifier string) []*tfe.Workspace {
	logger.Debug("Filtering workspaces", "runStatus", runStatus, "repoIdentifier", repoIdentifier, "totalWorkspaces", len(workspaces))

	var result []*tfe.Workspace

	for _, w := range workspaces {
		shouldInclude := hasRunStatus(*w, runStatus) && hasRepoIdentifier(*w, repoIdentifier)

		if shouldInclude {
			result = append(result, w)
		}
	}

	logger.Debug("Workspaces filtered", "filteredCount", len(result))
	return result
}

// hasRunStatus checks if a workspace matches the given run status filter
// If no run status given, return true
// Else return true only when a Current Run is available and matches
func hasRunStatus(w tfe.Workspace, runStatus string) bool {
	if runStatus == "" {
		return true // Empty means any run status should be included
	}
	if w.CurrentRun == nil {
		return false // Run status is not available, should not be included
	}
	return w.CurrentRun.Status == tfe.RunStatus(runStatus) // Status determines if it should be included
}

// hasRepoIdentifier checks if a workspace matches the given repository identifier filter
// If no repo given, return true
// Else return true only when a Repo identifier is available and matches
func hasRepoIdentifier(w tfe.Workspace, repoIdentifier string) bool {
	if repoIdentifier == "" {
		return true
	}
	if w.VCSRepo == nil {
		return false
	}
	return w.VCSRepo.Identifier == repoIdentifier
}

// ValidateRunStatus validates if the given run status is valid
func ValidateRunStatus(s string) bool {
	if s == "" {
		return true
	}
	validStatuses := []string{
		"pending", "plan_queued", "planning", "planned",
		"cost_estimating", "cost_estimated", "policy_checking",
		"policy_override", "policy_soft_failed", "policy_checked",
		"confirmed", "planned_and_finished", "apply_queued",
		"applying", "applied", "discarded", "errored",
		"canceled", "force_canceled",
	}
	for _, status := range validStatuses {
		if status == s {
			return true
		}
	}
	return false
}

// GetTeamAccessNames retrieves team names from team access objects
func GetTeamAccessNames(c *client.TfxClient, teamAccess []*tfe.TeamAccess) ([]interface{}, error) {
	logger.Debug("Fetching team names from team access", "count", len(teamAccess))

	var teamNames []interface{}
	for _, ta := range teamAccess {
		team, err := c.Client.Teams.Read(c.Context, ta.Team.ID)
		if err != nil {
			logger.Error("Failed to fetch team", "teamID", ta.Team.ID, "error", err)
			return nil, errors.Wrapf(err, "failed to read team %s", ta.Team.ID)
		}
		teamNames = append(teamNames, team.Name)
	}

	logger.Debug("Team names fetched successfully", "count", len(teamNames))
	return teamNames, nil
}
