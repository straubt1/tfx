# Roadmap

As we look to the future, I want to place my thoughts on where the project is headed and some of the core work that needs to be done.

## Core Decisions

- Should this repository stay under my personal GitHub account or be moved to an organization?
- Should we enable GitHub Projects or other tooling for tracking work?
- Do we need to refresh the documentation site? and process?

## List of things to do

These are random things that have been on the back burner for a while.

- [x] Initialize a Unit Testing Framework
- [x] Clean up debugging tfe client saving to file
- [x] Error Handling and formatting
  - Get specific errors from the tfe client and format them nicely
- [x] List and show command output (what to display and how to format it for list/show)
  - List might only show id.name.description but when outputted to json, should it be the full object?
- [ ] Integration Testing (HCPT is easy, TFE versions is harder)
- [ ] Automated Releases and updates to the [brew tap repo](https://github.com/straubt1/homebrew-tap)
- [ ] Go version management is fragmented (actions, readme, and go.mod)
- [ ] Versioning is manual and hard coded in version.go
- [ ] GHA need renamed and cleaned up
- [ ] Each Command needs an example, like in Projects
- [ ] Update Org and Proj to make additional API call for agent pool settings (if set)

## Things that might be nice to have or a terrible idea

- [ ] Embedded json filtering (similar to azure cli with JMESPath)
- [ ] TUI for diving deep into things
- [ ] Having a diff across like entities (e.g. diff workspaces, project, etc...)

## Command Refactor

Items that came up while working on commands.

### General

- no c.Client calls in the cmd/ files
- no fmt.Println() or like in cmd/view files

### Workspace

- Refactor out old helper functions at the bottom of cmd/workspace.go
- Team Access
  - we list the team name, but should we add access?
- Remote Sharing
  - Do we like the view of listing project names?
```
Remote State Sharing Workspaces:
  - local-workspace
  - aws-drift-test
```

## Long Term Goals

- Path to 1.0
  - Production Ready
  - Consistent Automated Updates
  - Testing Framework and Coverage, including Integration Tests
- Improved Documentation
- Better error handling
- Solid framework for additional functionality
