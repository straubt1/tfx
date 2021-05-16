package version

var (
	Version    = "0.0.0"
	Prerelease = "dev"
	Build      = ""
	Date       = ""
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
		v += "\nDate:  " + Date
	}
	return v
}
