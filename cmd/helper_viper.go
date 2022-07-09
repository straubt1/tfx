package cmd

import (
	"errors"
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/spf13/viper"
)

func viperString(flag string) *string {
	if viper.GetString(flag) == "" {
		// return nil
		value := ""
		return &value
	}
	value := viper.GetString(flag)
	return &value
}

func viperInt64(flag string) *int64 {
	value := viper.GetInt64(flag)
	return &value
}

func viperInt(flag string) *int {
	value := viper.GetInt(flag)
	return &value
}

func viperSemanticVersionString(flag string) (string, error) {
	v, err := semver.NewVersion(viper.GetString(flag))
	if err != nil {
		return "", errors.New("invalid semantic version")
	}
	return v.String(), nil
}

func viperBool(flag string) *bool {
	if !viper.GetBool(flag) {
		// return nil //TODO: figure out why defaults are not working
		value := false
		return &value
	}
	value := viper.GetBool(flag)
	return &value
}

func viperStringSlice(flag string) []string {
	value := viper.GetStringSlice(flag)
	if len(value) == 0 {
		return []string{}
	}
	return value
}

func viperStringSliceMap(flag string) (map[string]string, error) {
	m := make(map[string]string)
	values := viper.GetStringSlice(flag)

	for _, v := range values {
		// Expecting each value to be in "a=1" format
		s := strings.SplitN(v, "=", 2)
		if len(s) != 2 {
			return nil, errors.New("invalid env var")
		}
		m[s[0]] = s[1]
	}
	return m, nil
}
