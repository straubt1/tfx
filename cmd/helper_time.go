package cmd

import (
	"time"
)

func timestamp(time time.Time) string {
	// TODO: normalize this to current users timezone
	return time.Format("2006-01-02 15:04:05")
}

// Format date consistently
func FormatDateTime(t time.Time) string {
	return t.Format("Mon Jan _2 15:04 2006")
}
