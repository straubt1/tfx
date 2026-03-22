// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package cmd

import (
	"fmt"
	"log"

	"github.com/go-viper/encoding/hcl"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/straubt1/tfx/output"
	"github.com/straubt1/tfx/pkg/hclconfig"
	"github.com/straubt1/tfx/version"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	// Required to leverage viper defaults for optional Flags
	bindPFlags = func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			log.Fatal(err.Error())
		}
	}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tfx",
	Short: "A CLI to easily interact with TFC/TFE.",
	Long: `Leveraging the API can become a burden for common tasks.
	TFx aims to ease that challenge for common and repeatable tasks. This application
	can be used to interact with either HCP Terraform or Terraform Enterprise.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	Version:       version.String(),
	// PersistentPreRunE binds flags to viper, resolves the active profile, then
	// validates that credentials are present for all commands except 'login'.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		bindPFlags(cmd, args)
		resolveActiveProfile()
		if cmd.Name() == "login" {
			return nil
		}
		if viper.GetString("tfeToken") == "" {
			return fmt.Errorf("no API token found — run 'tfx login' to authenticate")
		}
		if cmd.Name() != "tui" && viper.GetString("tfeOrganization") == "" {
			return fmt.Errorf("organization is required (--tfeOrganization, TFE_ORGANIZATION, or run 'tfx login')")
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()

	// Always close output system for clean shutdown
	output.Get().Close()

	if err != nil {
		log.Fatal(aurora.Red(err))
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file, can be used to store common flags, (default is ./.tfx.hcl).")
	rootCmd.PersistentFlags().String("tfeHostname", "app.terraform.io", "The hostname of TFE without the schema. Can also be set with the environment variable TFE_HOSTNAME.")
	rootCmd.PersistentFlags().String("tfeOrganization", "", "The name of the TFx Organization. Can also be set with the environment variable TFE_ORGANIZATION.")
	rootCmd.PersistentFlags().String("tfeToken", "", "The API token used to authenticate to TFx. Can also be set with the environment variable TFE_TOKEN.")
	rootCmd.PersistentFlags().String("profile", "", "Named profile to use from ~/.tfx.hcl.")

	// Add json output option
	rootCmd.PersistentFlags().BoolP("json", "j", false, "Will output command results as JSON.")

	// ENV aliases
	viper.BindEnv("tfeHostname", "TFE_HOSTNAME")
	viper.BindEnv("tfeOrganization", "TFE_ORGANIZATION")
	viper.BindEnv("tfeToken", "TFE_TOKEN")
	viper.BindEnv("profile", "TFX_PROFILE")

	// Turn off completion option
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in current & home directory with name ".tfx" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(".tfx")
	}

	// Load 3rd party extension for HCL
	codecRegistry := viper.NewCodecRegistry()
	codecRegistry.RegisterCodec("hcl", hcl.Codec{})
	viper.SetOptions(viper.WithCodecRegistry(codecRegistry))

	// read in environment variables that match
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	isConfigFile := false
	err := viper.ReadInConfig()
	if err == nil {
		isConfigFile = true // Capture information here to bring after all flags are loaded (namely which output type)
	} else {
		output.Get().Logger().Warn("Unable to parse config file, will continue without it.")
	}

	// Some hacking here to let viper use the cobra required flags, simplifies this checking
	// in one place rather than each command
	// More info: https://github.com/spf13/viper/issues/397
	postInitCommands(rootCmd.Commands())

	// Output system initializes automatically on first use via output.Get()
	// Show config file message if found
	if isConfigFile {
		out := output.Get()
		out.Message("Using config file: %s", viper.ConfigFileUsed())
	}
}

// resolveActiveProfile reads profile blocks from the loaded config file and
// promotes the active profile's values to the flat viper keys (tfeHostname,
// tfeToken, tfeOrganization) that the rest of the app reads.
//
// Profile selection priority:
//  1. --profile flag (or TFX_PROFILE env var) — use the named profile
//  2. No flag — use the first profile block in the file
//  3. Old flat format (no profile blocks) — nothing to do; flat keys already set
func resolveActiveProfile() {
	configPath := viper.ConfigFileUsed()
	if configPath == "" {
		return
	}

	profiles, err := hclconfig.ListProfiles(configPath)
	if err != nil || len(profiles) == 0 {
		return
	}

	profileFlag := viper.GetString("profile")
	var active *hclconfig.Profile
	if profileFlag != "" {
		for i := range profiles {
			if profiles[i].Name == profileFlag {
				active = &profiles[i]
				break
			}
		}
		// Unknown profile — let validation in PersistentPreRunE report the error.
		if active == nil {
			return
		}
	} else {
		// No --profile flag: prefer a profile named "default", fall back to first.
		for i := range profiles {
			if profiles[i].Name == "default" {
				active = &profiles[i]
				break
			}
		}
		if active == nil {
			active = &profiles[0]
		}
	}

	// Store the resolved profile name so callers (e.g. the TUI) can read it
	// back from Viper even when --profile was not explicitly provided.
	viper.Set("profile", active.Name)

	if active.Hostname != "" {
		viper.Set("tfeHostname", active.Hostname)
	}
	viper.Set("tfeToken", active.Token)
	viper.Set("tfeOrganization", active.Organization)
}

// copy.pasta function
func postInitCommands(commands []*cobra.Command) {
	for _, cmd := range commands {
		presetRequiredFlags(cmd)
		if cmd.HasSubCommands() {
			postInitCommands(cmd.Commands())
		}
	}
}

// copy.pasta function
func presetRequiredFlags(cmd *cobra.Command) {
	viper.BindPFlags(cmd.Flags())
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
			cmd.Flags().Set(f.Name, viper.GetString(f.Name))
		}
	})
}
