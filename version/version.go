// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package version

var (
	Version    = "0.1.5"
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
		v += "\nBuild: " + Build
	}
	if Date != "" {
		v += "\nDate: " + Date
	}
	v += "\nBuilt By: " + BuiltBy

	return v
}
