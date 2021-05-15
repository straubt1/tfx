/*
Copyright © 2021 Tom Straub <github.com/straubt1>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/straubt1/tfx/version"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile string

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
	Long: `Leveraging the API for TFC/TFE can become a burden for common tasks.
	TFx aims to ease that challenge for common and repeatable tasks. This application
	can be used to interact with either Terraform Cloud or Terraform Enterprise.`,
	Version: version.String(),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file, can be used to store common flags, (default is ./.tfx.hcl).")
	rootCmd.PersistentFlags().String("tfeHostname", "app.terraform.io", ".")
	rootCmd.PersistentFlags().String("tfeOrganization", "", ".")
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

		// Search config in home directory with name ".tfx" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".tfx")
		// Search config in working directory with name ".tfx" (without extension).
		viper.AddConfigPath(".")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	// set the viper config file values
	// tfeHostname = viper.GetString("tfeHostname")
	// tfeToken = viper.GetString("tfeToken")
	// tfeOrganization = viper.GetString("tfeOrganization")

	// for _, e := range os.Environ() {
	// 	pair := strings.SplitN(e, "=", 2)
	// 	fmt.Println(pair[0])
	// }
}
