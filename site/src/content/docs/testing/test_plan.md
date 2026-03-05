---
title: TFx Test Plan
---

This document is a structured walkthrough for early users to validate TFx against a live HCP Terraform or Terraform Enterprise instance. Work through each section in order — later sections depend on resources created earlier.

:::note[Feedback]
If you encounter any unexpected behavior, please [open an issue](https://github.com/straubt1/tfx/issues) with the command you ran, the output you received, and your TFx version (`tfx version`).
:::

## Setup

### 1. Installation

Verify TFx is installed and reporting the correct version.

```sh
tfx version
```

**Expected:** Version string is printed with no errors.

### 2. Configuration

Create a `.tfx.hcl` config file or export environment variables, then confirm connectivity:

```sh
# Option A: config file
cat <<EOF > .tfx.hcl
tfeHostname     = "app.terraform.io"
tfeOrganization = "my-org"
tfeToken        = "my-token"
EOF

# Option B: environment variables
export TFE_HOSTNAME=app.terraform.io
export TFE_ORGANIZATION=my-org
export TFE_TOKEN=my-token
```

Validate the configuration by listing organizations — this is the simplest authenticated call:

```sh
tfx organization list
```

**Expected:** Table of organizations is displayed with no authentication errors.

---

## Organization

### 3. List Organizations

```sh
tfx organization list
```

**Expected:** Table with at least one organization row.

### 4. List Organizations (alias)

```sh
tfx org list
```

**Expected:** Same output as above — confirms the `org` alias works.

### 5. List Organizations (JSON output)

```sh
tfx organization list --json
```

**Expected:** Valid JSON array. Pipe through `jq` to confirm structure:

```sh
tfx organization list --json | jq '.[0].Name'
```

### 6. Show Organization

Replace `<org-name>` with a name from the list above.

```sh
tfx organization show --name <org-name>
```

**Expected:** Detail view showing the organization name, email, and plan.

---

## Project

### 7. List Projects

```sh
tfx project list
```

**Expected:** Table of projects for the configured organization.

### 8. List Projects (alias + search)

```sh
tfx prj list --search default
```

**Expected:** Filtered list. If no results, try a different search term or omit `--search`.

### 9. List Projects (JSON)

```sh
tfx project list --json
```

**Expected:** Valid JSON array of project objects.

### 10. Show Project by Name

Replace `<project-name>` with a name from the list.

```sh
tfx project show --name <project-name>
```

**Expected:** Project ID, name, and description are displayed.

### 11. Show Project by ID

Use the ID returned in the previous step.

```sh
tfx project show --id <project-id>
```

**Expected:** Same output as the name lookup.

---

## Workspace

### 12. List Workspaces

```sh
tfx workspace list
```

**Expected:** Table of workspaces with name, ID, resource count, and status columns.

### 13. List Workspaces (alias)

```sh
tfx ws list
```

**Expected:** Same output — confirms the `ws` alias works.

### 14. List Workspaces (search filter)

```sh
tfx workspace list --search <partial-name>
```

**Expected:** Filtered list containing only workspaces whose names match the search string.

### 15. List Workspaces (wildcard filter)

```sh
tfx workspace list --wildcard-name "*dev*"
```

**Expected:** Workspaces whose names contain `dev`. If none, try a pattern that matches your environment.

### 16. List Workspaces (run status filter)

```sh
tfx workspace list --run-status errored
```

**Expected:** Only workspaces with a current run in `errored` state. An empty table is a valid result if none exist.

### 17. List Workspaces (project filter)

Use a project ID from [step 10](#10-show-project-by-name).

```sh
tfx workspace list --project-id <project-id>
```

**Expected:** Only workspaces that belong to the specified project.

### 18. List Workspaces (all organizations)

```sh
tfx workspace list --all
```

**Expected:** Table includes an `ORGANIZATION` column and workspaces from every org accessible to your token.

### 19. List Workspaces (JSON)

```sh
tfx workspace list --json
```

**Expected:** Valid JSON array.

### 20. Show Workspace

Replace `<workspace-name>` with a workspace name from the list. Record the workspace name — it is used throughout the rest of this plan.

```sh
tfx workspace show --name <workspace-name>
```

**Expected:** Detail view showing ID, Terraform version, execution mode, lock status, and current run information.

---

## Workspace Lock / Unlock

:::caution
These steps temporarily lock a workspace. Ensure the workspace is not in active use before proceeding.
:::

### 21. Lock a Workspace

```sh
tfx workspace lock --name <workspace-name>
```

**Expected:** Confirmation that the workspace is now locked.

### 22. Verify Lock

```sh
tfx workspace show --name <workspace-name>
```

**Expected:** `Locked: true` in the output.

### 23. Unlock a Workspace

```sh
tfx workspace unlock --name <workspace-name>
```

**Expected:** Confirmation that the workspace is now unlocked.

### 24. Verify Unlock

```sh
tfx workspace show --name <workspace-name>
```

**Expected:** `Locked: false` in the output.

### 25. Lock All (search scoped)

Use a search string that limits scope to avoid accidentally locking production workspaces.

```sh
tfx workspace lock all --search <safe-prefix>
```

**Expected:** All matching workspaces are locked and a summary is printed.

### 26. Unlock All (search scoped)

```sh
tfx workspace unlock all --search <safe-prefix>
```

**Expected:** All matching workspaces are unlocked.

---

## Workspace Variables

### 27. List Variables

```sh
tfx workspace variable list --name <workspace-name>
```

**Expected:** Table of workspace variables (may be empty for a new workspace).

### 28. List Variables (alias)

```sh
tfx ws var list --name <workspace-name>
```

**Expected:** Same output — confirms the `var` alias works.

### 29. Create a Terraform Variable

```sh
tfx workspace variable create \
  --name <workspace-name> \
  --key tfx_test_var \
  --value "hello-from-tfx"
```

**Expected:** Confirmation showing the new variable ID and key.

### 30. Create an Environment Variable

```sh
tfx workspace variable create \
  --name <workspace-name> \
  --key TFX_TEST_ENV \
  --value "env-value" \
  --env
```

**Expected:** Confirmation showing the new variable with `env` category.

### 31. Create a Sensitive Variable

```sh
tfx workspace variable create \
  --name <workspace-name> \
  --key tfx_secret \
  --value "s3cr3t" \
  --sensitive
```

**Expected:** Confirmation showing `Sensitive: true`. The value will not be readable afterward.

### 32. Create an HCL Variable

```sh
tfx workspace variable create \
  --name <workspace-name> \
  --key tfx_map \
  --value '{"key1"="val1","key2"="val2"}' \
  --hcl
```

**Expected:** Confirmation showing `HCL: true`.

### 33. Show a Variable

```sh
tfx workspace variable show \
  --name <workspace-name> \
  --key tfx_test_var
```

**Expected:** Detail view for `tfx_test_var`.

### 34. Update a Variable

```sh
tfx workspace variable update \
  --name <workspace-name> \
  --key tfx_test_var \
  --value "updated-value" \
  --description "Updated by TFx test plan"
```

**Expected:** Confirmation that the variable was updated.

### 35. List Variables (JSON, verify update)

```sh
tfx workspace variable list --name <workspace-name> --json | jq '.[] | select(.Key=="tfx_test_var")'
```

**Expected:** JSON object showing `Value: "updated-value"`.

### 36. Delete Variables (cleanup)

```sh
tfx workspace variable delete --name <workspace-name> --key tfx_test_var
tfx workspace variable delete --name <workspace-name> --key TFX_TEST_ENV
tfx workspace variable delete --name <workspace-name> --key tfx_secret
tfx workspace variable delete --name <workspace-name> --key tfx_map
```

**Expected:** Each deletion is confirmed with no errors.

---

## Workspace Configuration Versions

### 37. List Configuration Versions

```sh
tfx workspace configuration-version list --name <workspace-name>
```

**Expected:** Table of existing configuration versions (may be empty).

### 38. List Configuration Versions (alias + max-items)

```sh
tfx ws cv list --name <workspace-name> --max-items 5
```

**Expected:** At most 5 rows returned.

### 39. Create a Configuration Version

Run this from a directory that contains Terraform files (even a minimal `main.tf` is sufficient).

```sh
# Create a minimal Terraform config if needed
mkdir -p /tmp/tfx-test && echo 'terraform {}' > /tmp/tfx-test/main.tf

tfx workspace configuration-version create \
  --name <workspace-name> \
  --directory /tmp/tfx-test
```

**Expected:** Upload completes and a configuration version ID (`cv-...`) is printed. Record this ID.

### 40. Create a Speculative Configuration Version

```sh
tfx workspace configuration-version create \
  --name <workspace-name> \
  --directory /tmp/tfx-test \
  --speculative
```

**Expected:** Configuration version created with `Speculative: true`.

### 41. Show Configuration Version

Use the ID from step 39.

```sh
tfx workspace configuration-version show --id <cv-id>
```

**Expected:** Detail view showing status, upload URL, and speculative flag.

### 42. Download Configuration Version

```sh
tfx workspace configuration-version download \
  --id <cv-id> \
  --directory /tmp/tfx-download
```

**Expected:** Archive is saved to `/tmp/tfx-download`. Verify the file exists:

```sh
ls /tmp/tfx-download
```

---

## Workspace Runs

### 43. List Runs

```sh
tfx workspace run list --name <workspace-name>
```

**Expected:** Table of recent runs (up to 10 by default).

### 44. List Runs (max-items)

```sh
tfx workspace run list --name <workspace-name> --max-items 3
```

**Expected:** At most 3 rows.

### 45. Create a Run

```sh
tfx workspace run create \
  --name <workspace-name> \
  --message "TFx test plan run"
```

**Expected:** A new run ID (`run-...`) is returned. Record this ID.

### 46. Show a Run

```sh
tfx workspace run show --id <run-id>
```

**Expected:** Run details including status, message, and timestamps.

### 47. Discard a Run

If the run from step 45 is in `planned` or `planned_and_finished` status, discard it:

```sh
tfx workspace run discard --id <run-id>
```

**Expected:** Run is discarded. Confirm with `run show`.

### 48. Cancel Latest Run

If a run is queued or currently applying, cancel it:

```sh
tfx workspace run cancel --name <workspace-name>
```

**Expected:** The latest run on the workspace is cancelled.

---

## Workspace Plans

### 49. Create a Plan

```sh
tfx workspace plan create \
  --name <workspace-name> \
  --directory /tmp/tfx-test \
  --message "TFx test plan"
```

**Expected:** A run is queued and a plan ID is returned. Record the plan ID (`plan-...`) from the output.

### 50. Show Plan

```sh
tfx workspace plan show --id <plan-id>
```

**Expected:** Plan detail view showing status and resource counts.

### 51. Plan Logs

```sh
tfx workspace plan logs --id <plan-id>
```

**Expected:** Streaming or full plan log output (similar to `terraform plan` output).

### 52. Plan JSON Output

```sh
tfx workspace plan jsonoutput --id <plan-id>
```

**Expected:** Raw JSON plan output. Pipe through `jq` to verify structure:

```sh
tfx workspace plan jsonoutput --id <plan-id> | jq '.format_version'
```

### 53. Create a Speculative Plan

```sh
tfx workspace plan create \
  --name <workspace-name> \
  --directory /tmp/tfx-test \
  --speculative
```

**Expected:** Speculative plan is created and does not trigger an apply.

---

## Workspace State Versions

### 54. List State Versions

```sh
tfx workspace state-version list --name <workspace-name>
```

**Expected:** Table of state versions (may be empty for a new workspace).

### 55. List State Versions (alias + max-items)

```sh
tfx ws sv list --name <workspace-name> --max-items 5
```

**Expected:** At most 5 rows.

If the workspace has at least one state version, continue with steps 56–58. Otherwise skip to [Workspace Teams](#workspace-teams).

### 56. Show State Version

```sh
# Get a state version ID from the list above, then:
tfx workspace state-version show --state-id <sv-id>
```

**Expected:** Detail view with size, created timestamp, and download URL.

### 57. Download State Version

```sh
tfx workspace state-version download \
  --state-id <sv-id> \
  --directory /tmp/tfx-state
```

**Expected:** State file is saved. Verify:

```sh
ls /tmp/tfx-state && cat /tmp/tfx-state/*.json | jq '.version'
```

---

## Workspace Teams

### 58. List Workspace Teams

```sh
tfx workspace team list --name <workspace-name>
```

**Expected:** Table of teams with their access level for this workspace (may be empty).

---

## Registry — Modules

### 59. List Modules

```sh
tfx registry module list
```

**Expected:** Table of private registry modules (may be empty).

### 60. List Modules (JSON + all pages)

```sh
tfx registry module list --all --json
```

**Expected:** Full JSON array of all modules without pagination truncation.

### 61. Create a Module

```sh
tfx registry module create \
  --name tfx-test-module \
  --provider aws
```

**Expected:** Module is created and its ID is returned. The module has no versions yet.

### 62. Show Module

```sh
tfx registry module show \
  --name tfx-test-module \
  --provider aws
```

**Expected:** Module detail view showing name, provider, and empty versions.

### 63. Create a Module Version

Run from a directory that contains a valid Terraform module.

```sh
tfx registry module version create \
  --name tfx-test-module \
  --provider aws \
  --version 1.0.0 \
  --directory /tmp/tfx-test
```

**Expected:** Version is uploaded and `1.0.0` is confirmed.

### 64. List Module Versions

```sh
tfx registry module version list \
  --name tfx-test-module \
  --provider aws
```

**Expected:** Table showing `1.0.0`.

### 65. Download Module Version

```sh
tfx registry module version download \
  --name tfx-test-module \
  --provider aws \
  --version 1.0.0 \
  --directory /tmp/tfx-module-download
```

**Expected:** Module archive saved to `/tmp/tfx-module-download`.

### 66. Delete Module Version

```sh
tfx registry module version delete \
  --name tfx-test-module \
  --provider aws \
  --version 1.0.0
```

**Expected:** Version is deleted without error.

### 67. Delete Module (cleanup)

```sh
tfx registry module delete \
  --name tfx-test-module \
  --provider aws
```

**Expected:** Module is deleted without error.

---

## Registry — Providers

:::note
Provider registration requires a GPG key. Complete the [Admin — GPG Keys](#admin-gpg-keys) section first if you intend to test provider version creation.
:::

### 68. List Providers

```sh
tfx registry provider list
```

**Expected:** Table of private registry providers (may be empty).

### 69. List Providers (JSON + all pages)

```sh
tfx registry provider list --all --json
```

**Expected:** Full JSON array.

### 70. Create a Provider

```sh
tfx registry provider create --name tfx-test-provider
```

**Expected:** Provider is created and its ID is returned.

### 71. Show Provider

```sh
tfx registry provider show --name tfx-test-provider
```

**Expected:** Provider detail view.

### 72. Delete Provider (cleanup)

```sh
tfx registry provider delete --name tfx-test-provider
```

**Expected:** Provider is deleted without error.

---

## Admin — GPG Keys

:::note
These commands require admin-level API token permissions.
:::

### 73. List GPG Keys

```sh
tfx admin gpg list
```

**Expected:** Table of GPG keys for the organization (may be empty).

### 74. Create a GPG Key

Generate or use an existing GPG public key. Export the armored public key to a file:

```sh
# Example: export from local keyring
gpg --armor --export <key-fingerprint> > /tmp/tfx-test-public.asc

tfx admin gpg create \
  --namespace <org-name> \
  --public-key /tmp/tfx-test-public.asc
```

**Expected:** GPG key ID (`<key-id>`) is returned. Record it.

### 75. Show GPG Key

```sh
tfx admin gpg show \
  --namespace <org-name> \
  --id <key-id>
```

**Expected:** Key detail view showing fingerprint and namespace.

### 76. Delete GPG Key (cleanup)

```sh
tfx admin gpg delete \
  --namespace <org-name> \
  --id <key-id>
```

**Expected:** Key is deleted without error.

---

## Admin — Terraform Versions

:::note
These commands require admin-level API token permissions and are only applicable to Terraform Enterprise. HCP Terraform does not support custom Terraform version management.
:::

### 77. List Terraform Versions

```sh
tfx admin terraform-version list
```

**Expected:** Table of available Terraform versions.

### 78. List Terraform Versions (alias + search)

```sh
tfx admin tfv list --search 1.9
```

**Expected:** Filtered list of `1.9.x` versions.

### 79. Show Terraform Version

```sh
tfx admin terraform-version show --version 1.9.0
```

**Expected:** Detail view for version `1.9.0` including enabled/disabled status.

### 80. Create Official Terraform Version

```sh
tfx admin terraform-version create official \
  --version 1.9.9 \
  --disable
```

**Expected:** Version `1.9.9` is created in disabled state.

### 81. Enable a Terraform Version

```sh
tfx admin terraform-version enable --versions 1.9.9
```

**Expected:** Version is now enabled.

### 82. Disable a Terraform Version

```sh
tfx admin terraform-version disable --versions 1.9.9
```

**Expected:** Version is now disabled.

### 83. Delete Terraform Version (cleanup)

```sh
tfx admin terraform-version delete --version 1.9.9
```

**Expected:** Version is deleted without error.

---

## Global Flags

These tests verify cross-cutting behavior available on all commands.

### 84. JSON Output Flag (short form)

```sh
tfx workspace list -j
```

**Expected:** Same JSON output as `--json`.

### 85. Config File Flag

```sh
tfx organization list --config /path/to/custom/.tfx.hcl
```

**Expected:** Command runs using the specified config file (confirmation message includes the config path).

### 86. Hostname and Token Flags

```sh
tfx organization list \
  --tfeHostname app.terraform.io \
  --tfeOrganization <org-name> \
  --tfeToken <token>
```

**Expected:** Same output as using the config file — confirms inline flags override the config.

### 87. Missing Required Flag Error

```sh
tfx workspace show
```

**Expected:** Error message indicating `--name` is required. Exit code is non-zero.

### 88. Invalid Token Error

```sh
tfx organization list --tfeToken invalid-token-value
```

**Expected:** Authentication error is reported clearly. No panic or unhandled stack trace.

---

## Summary Checklist

After completing all sections, confirm:

- [ ] All `list` commands return results or a clear empty-state message
- [ ] All `show` commands display structured detail views
- [ ] All `create` → `show` → `delete` workflows complete without error
- [ ] `--json` flag produces valid JSON on every command tested
- [ ] Authentication and config errors produce clear, human-readable messages
- [ ] No panics or unhandled exceptions were encountered
