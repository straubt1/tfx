package cmd

import (
	"time"
)

func timestamp(time time.Time) string {
	// TODO: normalize this to current users timezone
	return time.Format("2006-01-02 15:04:05")
}
