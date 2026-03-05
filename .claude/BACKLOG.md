# Skills Backlog

Future skills to build for Claude Code. Pick one up with: "look at my backlog and let's build the next skill."

## Completed

- [x] **new-data-function** — Conventions for the `data/` package (pagination, naming, SDK isolation)

## Planned

- [ ] **new-view** — Conventions for creating views in `cmd/views/` (embed `*BaseView`, `Render()`/`RenderError()`, terminal vs JSON output, `RenderProperties`/`RenderTags`/`RenderTable` primitives, dedicated JSON output structs with `omitempty`)

- [ ] **new-flags** — Conventions for creating flag structs in `cmd/flags/` (`*Flags` struct + `Parse*Flags(cmd)` function, `viper.GetString()`/`viper.GetBool()` for values, mutually exclusive / one-required patterns)

- [ ] **data-testing** — Conventions for testing the `data/` package (integration test setup with `//go:build integration`, `setupTest(t)` helper, `secrets/.env-int` config, co-located test files like `data/*_integration_test.go`)

- [ ] **go-tfe-pr** — Guide for contributing to the `hashicorp/go-tfe` SDK (fork workflow, struct field additions, jsonapi tags, test patterns, PR conventions)

- [ ] **prepare-release** — Pre-release checklist: verify docs are up to date for new/changed commands, run linting (`go vet`, `staticcheck`), run unit tests (`task test`), run build (`task go-build`), check CHANGELOG.md has an entry, run `task release-dry-run` to validate goreleaser
