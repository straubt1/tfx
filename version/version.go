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
