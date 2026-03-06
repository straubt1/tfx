---
name: go-tfe-pr
description: Commit, push, and open a draft Pull Request from the go-tfe fork against hashicorp/go-tfe. Use after reviewing and approving the code changes made by go-tfe-branch. The go-tfe repo is at ../go-tfe/ relative to this project.
---

# go-tfe Pull Request Creation

You are creating a draft Pull Request from the `straubt1/go-tfe` fork against the upstream `hashicorp/go-tfe` repository.

## Your Task

$ARGUMENTS

If no arguments are provided, inspect the current state of the go-tfe repo and proceed accordingly.

## Step-by-Step Process

### 1. Verify repo state

```bash
cd /Users/tstraub/Projects/straubt1.github.com/go-tfe
git status
git branch --show-current
git log main..HEAD --oneline
```

- Confirm you are on a feature branch (not `main`)
- Confirm there are staged or unstaged changes, or commits ahead of main
- If already committed and pushed, skip to step 3

If anything looks wrong (on main, no changes, unexpected state), stop and report to the user.

### 2. Commit the changes

Stage all relevant files and commit. Use a concise, descriptive commit message following conventional commits style:

```
<type>: <short description>

<optional body>
```

Types: `feat`, `fix`, `refactor`, `test`, `docs`

Do not mention tooling or automation in the commit message.

**Before committing**: confirm `go vet ./...` and `go build ./...` pass. If they fail, report the errors — do not commit broken code.

```bash
cd /Users/tstraub/Projects/straubt1.github.com/go-tfe
go vet ./...
go build ./...
git add -p   # or add specific files by name
git commit -m "$(cat <<'EOF'
<type>: <description>

<body if needed>
EOF
)"
```

### 3. Push to fork

```bash
git push -u origin <branch-name>
```

### 4. Create draft PR against upstream

Use `gh pr create` targeting `hashicorp/go-tfe` as the base repo:

```bash
gh pr create \
  --repo hashicorp/go-tfe \
  --base main \
  --draft \
  --title "<concise title>" \
  --body "$(cat <<'EOF'
## Description

<1-3 sentences describing what this PR does and why>

## Changes

<bullet list of specific changes made>

## Testing

<brief note on what was tested — e.g., "go vet and go build pass", or integration test details if run>
EOF
)"
```

#### PR title guidelines
- Short (under 60 characters)
- Start with a verb: "Add", "Fix", "Update", "Remove"
- No period at the end
- Examples:
  - `Add SettingOverwrites to ProjectUpdateOptions`
  - `Fix workspace list pagination for large orgs`
  - `Add ReadWithOptions method to Projects`

#### PR body guidelines
- **Description**: What the change does and why it's needed. Mention the API field or endpoint if relevant.
- **Changes**: Bulleted list of the concrete modifications (struct fields added, methods added, files changed).
- **Testing**: Honest note on what was verified locally.
- Do not pad with filler. Keep it concise and factual.
- Do not mention tooling or automation.

### 5. Update CHANGELOG.md with PR number

After the PR is created, `gh pr create` returns the PR URL. Extract the PR number and update the CHANGELOG.md placeholder:

```bash
# Replace #XXXX with the actual PR number
sed -i '' 's/#XXXX/#<actual-number>/g' /Users/tstraub/Projects/straubt1.github.com/go-tfe/CHANGELOG.md
```

Then amend the commit to include the updated CHANGELOG:

```bash
cd /Users/tstraub/Projects/straubt1.github.com/go-tfe
git add CHANGELOG.md
git commit --amend --no-edit
git push --force-with-lease origin <branch-name>
```

### 6. Report back

Output:
- The PR URL
- The PR title and number
- A reminder that the PR is a draft — the user must manually mark it ready for review

## Important Rules

- Always create PRs as **draft** (`--draft` flag)
- Always target `hashicorp/go-tfe` as the base repo, not the fork
- Never force-push to main
- Do not commit broken code (vet/build must pass)
- Do not mention tooling or automation in PR title, body, or commit messages
