// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package cmd

import (
	"log"

	"github.com/go-viper/encoding/hcl"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/straubt1/tfx/output"
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
	SilenceUsage:     true,
	SilenceErrors:    true,
	Version:          version.String(),
	PersistentPreRun: bindPFlags, // Bind here to avoid having to call this in every subcommand
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

	// Add json output option
	rootCmd.PersistentFlags().BoolP("json", "j", false, "Will output command results as JSON.")

	// required
	rootCmd.MarkPersistentFlagRequired("tfeOrganization")
	rootCmd.MarkPersistentFlagRequired("tfeToken")

	// ENV aliases
	viper.BindEnv("tfeHostname", "TFE_HOSTNAME")
	viper.BindEnv("tfeOrganization", "TFE_ORGANIZATION")
	viper.BindEnv("tfeToken", "TFE_TOKEN")

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
		logWarning(err, "Unable to parse config file, will continue without it.")
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
