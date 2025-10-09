# Roadmap

As we look to the future, I want to place my thoughts on where the project is headed and some of the core work that needs to be done.

## Core Decisions

- Should this repository stay under my personal GitHub account or be moved to an organization?
- Should we enable GitHub Projects or other tooling for tracking work?
- Do we need to refresh the documentation site? and process?

## List of things to do

These are random things that have been on the back burner for a while.

- [x] Initialize a Unit Testing Framework
- [ ] Clean up debugging tfe client saving to file
- [ ] Integration Testing (HCPT is easy, TFE versions is harder)
- [ ] Automated Releases and updates to the [brew tap repo](https://github.com/straubt1/homebrew-tap)
- [ ] Go version management is fragmented (actions, readme, and go.mod)
- [ ] Versioning is manual and hard coded in version.go
- [ ] GHA need renamed and cleaned up
- [ ] Each Command needs an example, like in Projects
- [ ] List and show command output (what to display and how to format it for list/show)
  - List might only show id.name.description but when outputted to json, should it be the full object?

## Things that might be nice to have or a terrible idea

- [ ] Embedded json filtering (similar to azure cli with JMESPath)
- [ ] TUI for diving deep into things

## Long Term Goals

- Path to 1.0
  - Production Ready
  - Consistent Automated Updates
  - Testing Framework and Coverage, including Integration Tests
- Improved Documentation
- Better error handling
- Solid framework for additional functionality
