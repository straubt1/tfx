// Copyright (c) Tom Straub (github.com/straubt1) 2025
// SPDX-License-Identifier: MIT

package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/straubt1/tfx/output"
	"github.com/straubt1/tfx/pkg/browser"
	"github.com/straubt1/tfx/pkg/hclconfig"
)

var loginCmd = &cobra.Command{
	Use:   "login [hostname]",
	Short: "Authenticate to HCP Terraform or Terraform Enterprise",
	Long: `Authenticate to HCP Terraform or Terraform Enterprise.

Opens your browser to the API token creation page, prompts you to paste the
token, then saves it (along with the selected organization) to ~/.tfx.hcl.

If no hostname is provided the default is app.terraform.io.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runLogin,
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func runLogin(cmd *cobra.Command, args []string) error {
	output.Get().DisableSpinner()

	hostname := "app.terraform.io"
	if len(args) == 1 {
		hostname = strings.TrimSpace(args[0])
	}

	configPath, err := hclconfig.DefaultConfigPath()
	if err != nil {
		return fmt.Errorf("finding config path: %w", err)
	}

	// ── Intro ─────────────────────────────────────────────────────────────────
	fmt.Printf("\nTFx will request an API token for %s using your browser.\n\n", hostname)
	fmt.Printf("If login is successful, TFx will store the token in plain text in\n")
	fmt.Printf("the following file for use by subsequent commands:\n    %s\n\n", configPath)

	// ── Confirmation ──────────────────────────────────────────────────────────
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Do you want to proceed?")
	fmt.Println("  Only 'yes' will be accepted to confirm.")
	fmt.Print("\nEnter a value: ")
	answer, _ := reader.ReadString('\n')
	if strings.TrimSpace(answer) != "yes" {
		fmt.Println("\nLogin cancelled.")
		return nil
	}

	// ── Open browser ──────────────────────────────────────────────────────────
	tokenURL := fmt.Sprintf("https://%s/app/settings/tokens?source=tfx-login", hostname)
	fmt.Printf("\nOpening browser to:\n  %s\n\n", tokenURL)
	if err := browser.Open(tokenURL); err != nil {
		fmt.Printf("  (Could not open browser automatically: %s)\n", err)
		fmt.Println("  Open the URL above manually to create a token.")
	}
	fmt.Println("Generate a token there, then paste it below.")

	// ── Read token (no echo) ──────────────────────────────────────────────────
	fmt.Print("\nEnter your API token (will not echo): ")
	var token string
	if term.IsTerminal(int(os.Stdin.Fd())) {
		tokenBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			return fmt.Errorf("reading token: %w", err)
		}
		token = strings.TrimSpace(string(tokenBytes))
	} else {
		// Piped input (e.g. in scripts) — read normally.
		line, _ := reader.ReadString('\n')
		token = strings.TrimSpace(line)
	}
	if token == "" {
		return fmt.Errorf("token cannot be empty")
	}

	// ── Validate token + fetch organizations ──────────────────────────────────
	fmt.Print("\nValidating token and fetching organizations... ")
	tfeClient, err := tfe.NewClient(&tfe.Config{
		Address: fmt.Sprintf("https://%s", hostname),
		Token:   token,
	})
	if err != nil {
		return fmt.Errorf("creating API client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var allOrgs []*tfe.Organization
	pageNum := 1
	for {
		result, err := tfeClient.Organizations.List(ctx, &tfe.OrganizationListOptions{
			ListOptions: tfe.ListOptions{PageNumber: pageNum, PageSize: 100},
		})
		if err != nil {
			return fmt.Errorf("\ntoken validation failed: %w", err)
		}
		allOrgs = append(allOrgs, result.Items...)
		if pageNum >= result.TotalPages {
			break
		}
		pageNum++
	}
	fmt.Println("✓")

	// ── Select organization ───────────────────────────────────────────────────
	var selectedOrg string
	switch len(allOrgs) {
	case 0:
		return fmt.Errorf("no organizations found for this token")
	case 1:
		selectedOrg = allOrgs[0].Name
		fmt.Printf("\n  [auto-selected] %s\n", selectedOrg)
	default:
		fmt.Printf("\nSelect an organization:\n")
		for i, org := range allOrgs {
			fmt.Printf("  %d. %s\n", i+1, org.Name)
		}
		fmt.Print("\nEnter number: ")
		numStr, _ := reader.ReadString('\n')
		n, err := strconv.Atoi(strings.TrimSpace(numStr))
		if err != nil || n < 1 || n > len(allOrgs) {
			return fmt.Errorf("invalid selection")
		}
		selectedOrg = allOrgs[n-1].Name
	}

	// ── Save to config ────────────────────────────────────────────────────────
	if err := hclconfig.WriteProfile(configPath, hostname, selectedOrg, token); err != nil {
		return fmt.Errorf("saving credentials: %w", err)
	}

	fmt.Printf("\nSaved profile %q to %s.\n", hostname, configPath)
	fmt.Println("\nTFx is ready. Try:")
	fmt.Printf("  tfx workspace list\n\n")
	return nil
}
