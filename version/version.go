package version

var (
	Version    = "0.0.1"
	Prerelease = "dev"
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
