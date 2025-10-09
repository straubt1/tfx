// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package cmd

import (
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/straubt1/tfx/client"
	"github.com/straubt1/tfx/data"
	"github.com/straubt1/tfx/flags"
)

// organizationCmd represents the organization command
var (
	// `tfx organization` commands
	organizationCmd = &cobra.Command{
		Use:     "organization",
		Aliases: []string{"org"},
		Short:   "Organization Commands",
		Long:    "Work with TFx Organizations",
		Example: `
List all Organizations:
tfx organization list

Show an Organization by name:
tfx organization show --name "my-org"`,
	}

	// `tfx organization list` command
	organizationListCmd = &cobra.Command{
		Use:     "list",
		Short:   "List Organizations",
		Long:    "List all Organizations available to the authenticated user.",
		Example: `tfx organization list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return organizationList()
		},
	}

	// `tfx organization show` command
	organizationShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show organization details",
		Long:  "Show Organization details.",
		Example: `
tfx organization show --name myorganization
tfx organization show -n myorganization`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdConfig, err := flags.ParseOrganizationShowFlags(cmd)
			if err != nil {
				return err
			}
			return organizationShow(cmdConfig)
		},
	}
)

func init() {
	// `tfx organization show`
	organizationShowCmd.Flags().StringP("name", "n", "", "Name of the organization.")
	organizationShowCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(organizationCmd)
	organizationCmd.AddCommand(organizationListCmd)
	organizationCmd.AddCommand(organizationShowCmd)
}

func organizationList() error {
	c, err := client.NewFromViper()
	if err != nil {
		return err
	}

	o.AddMessageUserProvided("List Organizations", "")
	orgs, err := data.FetchOrganizations(c, "")
	if err != nil {
		return errors.Wrap(err, "failed to list organizations")
	}

	o.AddTableHeader("Name", "Email", "External ID")
	for _, org := range orgs {
		o.AddTableRows(org.Name, org.Email, org.ExternalID)
	}

	return nil
}

func organizationShow(cmdConfig *flags.OrganizationShowFlags) error {
	c, err := client.NewFromViper()
	if err != nil {
		return err
	}

	var org *tfe.Organization
	readOptions := &tfe.OrganizationReadOptions{
		Include: []tfe.OrganizationIncludeOpt{},
	}

	o.AddMessageUserProvided("Organization Name:", cmdConfig.Name)
	org, err = data.FetchOrganization(c, cmdConfig.Name, readOptions)

	if err != nil {
		return errors.Wrap(err, "failed to read organization")
	}

	o.AddDeferredMessageRead("Name", org.Name)
	o.AddDeferredMessageRead("Email", org.Email)
	o.AddDeferredMessageRead("External ID", org.ExternalID)
	o.AddDeferredMessageRead("Created At", org.CreatedAt)
	o.AddDeferredMessageRead("Collaborator Auth Policy", org.CollaboratorAuthPolicy)
	o.AddDeferredMessageRead("Cost Estimation Enabled", org.CostEstimationEnabled)
	o.AddDeferredMessageRead("Owners Team SAML Role ID", org.OwnersTeamSAMLRoleID)
	o.AddDeferredMessageRead("SAML Enabled", org.SAMLEnabled)
	o.AddDeferredMessageRead("Session Remember Minutes", org.SessionRemember)
	o.AddDeferredMessageRead("Session Timeout Minutes", org.SessionTimeout)
	o.AddDeferredMessageRead("Two Factor Conformant", org.TwoFactorConformant)
	o.AddDeferredMessageRead("Trial Expires At", org.TrialExpiresAt)
	o.AddDeferredMessageRead("Default Execution Mode", org.DefaultExecutionMode)
	o.AddDeferredMessageRead("Is Unified", org.IsUnified)

	// Permissions
	if org.Permissions != nil {
		o.AddDeferredMessageRead("Can Create Team", org.Permissions.CanCreateTeam)
		o.AddDeferredMessageRead("Can Create Workspace", org.Permissions.CanCreateWorkspace)
		o.AddDeferredMessageRead("Can Create Workspace Migration", org.Permissions.CanCreateWorkspaceMigration)
		o.AddDeferredMessageRead("Can Destroy", org.Permissions.CanDestroy)
		o.AddDeferredMessageRead("Can Manage Run Tasks", org.Permissions.CanManageRunTasks)
		o.AddDeferredMessageRead("Can Traverse", org.Permissions.CanTraverse)
		o.AddDeferredMessageRead("Can Update", org.Permissions.CanUpdate)
	}

	return nil
}
