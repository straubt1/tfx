// Copyright © 2021 Tom Straub <github.com/straubt1>

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"log"

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
	o       *output.Output

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
	can be used to interact with either Terraform Cloud or Terraform Enterprise.`,
	SilenceUsage:     true,
	SilenceErrors:    true,
	Version:          version.String(),
	PersistentPreRun: bindPFlags, // Bind here to avoid having to call this in every subcommand
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Close output stream always before exiting
	if err := rootCmd.Execute(); err != nil {
		o.Close()
		log.Fatal(aurora.Red(err))
	} else {
		o.Close()
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

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	isConfigFile := false
	if err := viper.ReadInConfig(); err == nil {
		isConfigFile = true // Capture information here to bring after all flags are loaded (namely which output type)
	}

	// Some hacking here to let viper use the cobra required flags, simplifies this checking
	// in one place rather than each command
	// More info: https://github.com/spf13/viper/issues/397
	postInitCommands(rootCmd.Commands())

	// Initialize output
	o = output.New(*viperBool("json"))
	// Print if config file was found
	if isConfigFile {
		o.AddMessageCalculated("Using config file:", viper.ConfigFileUsed())
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
