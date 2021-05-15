package version

import (
	"fmt"
)

var Version = "0.0.1"
var Prerelease = "dev"

func String() string {
	if Prerelease != "" {
		return fmt.Sprintf("%s-%s", Version, Prerelease)
	}
	return Version
}
