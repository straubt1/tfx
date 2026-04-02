// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/go-viper/encoding/hcl"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/straubt1/tfx/output"
	"github.com/straubt1/tfx/pkg/hclconfig"
	"github.com/straubt1/tfx/tui"
	"github.com/straubt1/tfx/version"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	// userChangedFlags records which persistent flags were explicitly set on the
	// command line. We snapshot this in initConfig() *before* postInitCommands
	// runs, because postInitCommands calls cmd.Flags().Set() which marks flags
	// as Changed even though the user didn't provide them.
	userChangedFlags map[string]bool

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
	RunE: func(cmd *cobra.Command, args []string) error {
		tapePath, _ := cmd.Flags().GetString("tape")
		return tui.Run(tapePath)
	},
	// PersistentPreRunE binds flags to viper, resolves the active profile, then
	// validates that credentials are present for all commands except 'login'.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		bindPFlags(cmd, args)

		if cmd.Name() == "login" {
			return nil
		}

		if err := resolveProfile(); err != nil {
			return err
		}
		if viper.GetString("token") == "" {
			return fmt.Errorf("no API token found — run 'tfx login' to authenticate")
		}
		if cmd.Name() != "tfx" && viper.GetString("organization") == "" {
			return fmt.Errorf("organization is required (--organization, TFE_ORGANIZATION, or run 'tfx login')")
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config-file", "", "Path to config file (env: TFX_CONFIG_FILE). Auto-discovered at ./.tfx.hcl (current dir) or ~/.tfx.hcl (home dir) when not set.")
	rootCmd.PersistentFlags().String("hostname", "app.terraform.io", "The hostname of TFE without the schema. Can also be set with the environment variable TFE_HOSTNAME.")
	rootCmd.PersistentFlags().String("organization", "", "The name of the TFx Organization. Can also be set with the environment variable TFE_ORGANIZATION.")
	rootCmd.PersistentFlags().String("token", "", "The API token used to authenticate to TFx. Can also be set with the environment variable TFE_TOKEN.")
	rootCmd.PersistentFlags().StringP("profile", "p", "", "Named profile to use from the config file. Defaults to \"default\". Can also be set with the environment variable TFX_PROFILE.")

	// Add json output option
	rootCmd.PersistentFlags().BoolP("json", "j", false, "Will output command results as JSON.")

	// ENV aliases
	viper.BindEnv("hostname", "TFE_HOSTNAME")
	viper.BindEnv("organization", "TFE_ORGANIZATION")
	viper.BindEnv("token", "TFE_TOKEN")
	viper.BindEnv("profile", "TFX_PROFILE")

	// Hidden flag for VHS tape recording
	rootCmd.Flags().String("tape", "", "Record TUI input to a .tape file for VHS (e.g. debug/demo.tape)")
	rootCmd.Flags().MarkHidden("tape")

	// Turn off completion option
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Resolve config file: --config-file flag > TFX_CONFIG_FILE env var > auto-search.
	if cfgFile == "" {
		cfgFile = os.Getenv("TFX_CONFIG_FILE")
	}
	if cfgFile != "" {
		// Resolve to an absolute path so the "Using config file:" message and
		// the TUI profile bar always show the fully-qualified location.
		if abs, err := filepath.Abs(cfgFile); err == nil {
			cfgFile = abs
		}
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search current directory first, then home directory.
		// This lets a local .tfx.hcl override the one in ~/.
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.SetConfigName(".tfx")
	}

	// Load 3rd party extension for HCL
	codecRegistry := viper.NewCodecRegistry()
	codecRegistry.RegisterCodec("hcl", hcl.Codec{})
	viper.SetOptions(viper.WithCodecRegistry(codecRegistry))

	// read in environment variables that match
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		output.Get().Logger().Warn("Unable to parse config file, will continue without it.")
	}

	// Snapshot which persistent flags were explicitly set by the user on the
	// command line. Must happen BEFORE postInitCommands, which calls
	// cmd.Flags().Set() and marks flags as Changed even for config-file values.
	captureUserFlags()

	// Some hacking here to let viper use the cobra required flags, simplifies this checking
	// in one place rather than each command
	// More info: https://github.com/spf13/viper/issues/397
	postInitCommands(rootCmd.Commands())

}

// captureUserFlags snapshots which persistent flags were explicitly set on the
// command line. Called in initConfig before postInitCommands corrupts Changed.
func captureUserFlags() {
	userChangedFlags = make(map[string]bool)
	rootCmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		if f.Changed {
			userChangedFlags[f.Name] = true
		}
	})
}

// resolveProfile loads profile blocks from the config file and merges the
// active profile's values into Viper with correct precedence:
//
//	CLI flags (--hostname, --organization, --token)  — highest
//	Environment variables (TFE_HOSTNAME, etc.)
//	Profile values from .tfx.hcl
//	Defaults                                         — lowest
func resolveProfile() error {
	configPath := viper.ConfigFileUsed()
	if configPath == "" {
		return nil
	}

	profiles, err := hclconfig.ListProfiles(configPath)
	if err != nil || len(profiles) == 0 {
		// Config file exists but can't be parsed or has no profiles.
		// Not fatal — flags or env vars may still provide credentials.
		return nil
	}

	// Select profile: --profile flag > "default" name
	profileName := viper.GetString("profile")
	var active *hclconfig.Profile
	if profileName != "" {
		// User explicitly asked for a profile — error if not found.
		for i := range profiles {
			if profiles[i].Name == profileName {
				active = &profiles[i]
				break
			}
		}
		if active == nil {
			return fmt.Errorf("profile %q not found in %s", profileName, configPath)
		}
	} else {
		// Look for "default" profile. If not found, silently skip —
		// flags or env vars may still provide credentials.
		for i := range profiles {
			if profiles[i].Name == "default" {
				active = &profiles[i]
				break
			}
		}
		if active == nil {
			return nil
		}
	}

	// Store the resolved profile name so callers (e.g. the TUI) can read it back.
	viper.Set("profile", active.Name)

	// Print config file and profile (unless JSON mode).
	if !viper.GetBool("json") {
		output.Get().Message("Using config file: %s (profile: %s)", aurora.Blue(configPath), aurora.Blue(active.Name))
	}

	// Merge profile values with correct precedence.
	// Only set when no higher-precedence source (flag or env var) exists.
	type mapping struct {
		viperKey string
		envVar   string
		value    string
	}
	for _, m := range []mapping{
		{"hostname", "TFE_HOSTNAME", active.Hostname},
		{"token", "TFE_TOKEN", active.Token},
		{"organization", "TFE_ORGANIZATION", active.Organization},
	} {
		if userChangedFlags[m.viperKey] {
			continue // CLI flag wins
		}
		if os.Getenv(m.envVar) != "" {
			continue // env var wins (already in Viper via BindEnv)
		}
		if m.value != "" {
			viper.Set(m.viperKey, m.value)
		}
	}

	return nil
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
