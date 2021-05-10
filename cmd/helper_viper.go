package cmd

import "github.com/spf13/viper"

func viperString(flag string) *string {
	if viper.GetString(flag) == "" {
		// return nil
		value := ""
		return &value
	}
	value := viper.GetString(flag)
	return &value
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
