/*
Copyright © 2021 Tom Straub <github.com/straubt1>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// `tfx workspace lock` command
	workspaceLockCmd = &cobra.Command{
		Use:   "lock",
		Short: "Lock a Workspace",
		Long:  "Lock a Workspace in a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return workspaceLock(
				getTfxClientContext(),
				*viperString("tfeOrganization"),
				*viperString("workspaceName"))
		},
	}

	workspaceLockAllCmd = &cobra.Command{
		Use:   "all",
		Short: "Lock All Workspaces",
		Long:  "Lock All Workspaces in a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return workspaceLockAll(
				getTfxClientContext(),
				*viperString("tfeOrganization"),
				*viperString("search"))
		},
	}

	// `tfx workspace unlock` command
	workspaceUnlockCmd = &cobra.Command{
		Use:   "unlock",
		Short: "Unlock a Workspace",
		Long:  "Unlock a Workspace in a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return workspaceUnlock(
				getTfxClientContext(),
				*viperString("tfeOrganization"),
				*viperString("workspaceName"))
		},
	}

	// `tfx workspace unlock all` command
	workspaceUnlockAllCmd = &cobra.Command{
		Use:   "all",
		Short: "Unlock All Workspaces",
		Long:  "Unlock All Workspaces in a TFx Organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return workspaceUnlockAll(
				getTfxClientContext(),
				*viperString("tfeOrganization"),
				*viperString("search"))
		},
	}
)

func init() {
	// `tfx workspace lock`
	workspaceLockCmd.Flags().StringP("workspaceName", "w", "", "Workspace name")
	workspaceLockCmd.MarkFlagRequired("workspaceName")

	// `tfx workspace lock all`
	workspaceLockAllCmd.Flags().StringP("search", "s", "", "Search string for Workspace Name (optional).")

	// `tfx workspace unlock`
	workspaceUnlockCmd.Flags().StringP("workspaceName", "w", "", "Workspace name")
	workspaceUnlockCmd.MarkFlagRequired("workspaceName")

	// `tfx workspace unlock all`
	workspaceUnlockAllCmd.Flags().StringP("search", "s", "", "Search string for Workspace Name (optional).")

	workspaceCmd.AddCommand(workspaceLockCmd)
	workspaceLockCmd.AddCommand(workspaceLockAllCmd)
	workspaceCmd.AddCommand(workspaceUnlockCmd)
	workspaceUnlockCmd.AddCommand(workspaceUnlockAllCmd)
}

func workspaceLock(c TfxClientContext, orgName string, workspaceName string) error {
	o.AddMessageUserProvided("Lock Workspace in Organization:", orgName)
	status, err := setWorkspaceLock(c, orgName, workspaceName, true)
	if err != nil {
		return errors.Wrap(err, "unable to lock workspace")
	}

	o.AddDeferredMessageRead(workspaceName, status)
	o.Close()

	return nil
}

func workspaceLockAll(c TfxClientContext, orgName string, searchString string) error {
	o.AddMessageUserProvided("Lock All Workspace in Organization:", orgName)
	workspaceList, err := workspaceListAllForOrganization(c, orgName, searchString)
	if err != nil {
		return errors.Wrap(err, "failed to list workspaces")
	}
	totalWorkspaces := len(workspaceList)

	o.AddFormattedMessageCalculated("Locking %d Workspaces, please wait...", totalWorkspaces)
	for _, ws := range workspaceList {
		status, err := setWorkspaceLock(c, orgName, ws.Name, true)
		if err != nil {
			o.AddDeferredMessageRead(ws.Name, err.Error())
		} else {
			o.AddDeferredMessageRead(ws.Name, status)
		}
	}

	o.Close()

	return nil
}

func workspaceUnlock(c TfxClientContext, orgName string, workspaceName string) error {
	o.AddMessageUserProvided("Unlock Workspace in Organization:", orgName)
	status, err := setWorkspaceLock(c, orgName, workspaceName, false)
	if err != nil {
		return errors.Wrap(err, "unable to unlock workspace")
	}

	o.AddDeferredMessageRead(workspaceName, status)
	o.Close()

	return nil
}

func workspaceUnlockAll(c TfxClientContext, orgName string, searchString string) error {
	o.AddMessageUserProvided("Unlock All Workspace in Organization:", orgName)
	workspaceList, err := workspaceListAllForOrganization(c, orgName, searchString)
	if err != nil {
		return errors.Wrap(err, "failed to list workspaces")
	}
	totalWorkspaces := len(workspaceList)

	o.AddFormattedMessageCalculated("Unlocking %d Workspaces, please wait...", totalWorkspaces)
	for _, ws := range workspaceList {
		status, err := setWorkspaceLock(c, orgName, ws.Name, false)
		if err != nil {
			o.AddDeferredMessageRead(ws.Name, err.Error())
		} else {
			o.AddDeferredMessageRead(ws.Name, status)
		}
	}

	o.Close()

	return nil
}

func setWorkspaceLock(c TfxClientContext, orgName string, workspaceName string, lockSet bool) (string, error) {
	w, err := c.Client.Workspaces.Read(c.Context, orgName, workspaceName)
	if err != nil {
		return "", errors.New("failed to read workspace id")
	}

	if lockSet {
		if w.Locked {
			return "Workspace already locked", nil
		}
		_, err := c.Client.Workspaces.Lock(c.Context, w.ID, tfe.WorkspaceLockOptions{
			Reason: tfe.String("Locked via TFx"),
		})
		if err != nil {
			return "", errors.Wrap(err, "unable to lock workspace")
		}
		return "Locked", nil
	} else {
		if !w.Locked {
			return "Workspace already unlocked", nil
		}
		// TODO: Force unlocking here to get around unlocking a WS that has an active run pending. Revisit impact.
		_, err := c.Client.Workspaces.ForceUnlock(c.Context, w.ID)
		// _, err := client.Workspaces.Unlock(ctx, w.ID)
		if err != nil {
			return "", errors.Wrap(err, "unable to unlock workspace")
		}
		return "Unlocked", nil
	}
}
