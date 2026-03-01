// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package version

var (
	Version    = "dev" // overridden at build time via ldflags (goreleaser uses git tag)
	Prerelease = ""
	Build      = ""
	Date       = ""
	BuiltBy    = ""
)

func String() string {
	v := Version
	if Prerelease != "" {
		v += "-" + Prerelease
	}
	if Build != "" {
		v += "\nBuild:    " + Build
	}
	if Date != "" {
		v += "\nDate:     " + Date
	}
	builtBy := BuiltBy
	if builtBy == "" {
		builtBy = "local (no build metadata)"
	}
	v += "\nBuilt By: " + builtBy

	return v
}
