/*
Copyright © 2021 Tom Straub <tstraub@hashicorp.com>

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
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile         string
	tfeHostname     string
	tfeToken        string
	tfeOrganization string
	// envs            string
	envs []string
)

// var client tfe.Client
// var ctx context.Context

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tfx",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tfx.yaml)")
	rootCmd.PersistentFlags().StringVar(&tfeHostname, "tfeHostname", "app.terraform.io", "The hostname of TFE. Defaults to TFC 'app.terraform.io'.")
	rootCmd.PersistentFlags().StringVar(&tfeToken, "tfeToken", "", "The API Token to interact with TFE or TFC.")
	// rootCmd.MarkPersistentFlagRequired("tfeToken")
	rootCmd.PersistentFlags().StringVar(&tfeOrganization, "tfeOrganization", "", "The TFE or TFC Organization.")
	planCmd.PersistentFlags().StringSliceVarP(&envs, "envs", "e", []string{}, "Array on ENV")

	// must bind to viper to pick up config file settings
	viper.BindPFlag("tfeHostname", rootCmd.PersistentFlags().Lookup("tfeHostname"))
	viper.BindPFlag("tfeToken", rootCmd.PersistentFlags().Lookup("tfeToken"))
	viper.BindPFlag("tfeOrganization", rootCmd.PersistentFlags().Lookup("tfeOrganization"))
	viper.BindPFlag("envs", rootCmd.PersistentFlags().Lookup("envs"))
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
		viper.AddConfigPath(".") // Look for config in the working directory
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	// set the viper config file values
	tfeHostname = viper.GetString("tfeHostname")
	tfeToken = viper.GetString("tfeToken")
	tfeOrganization = viper.GetString("tfeOrganization")

	// for _, e := range os.Environ() {
	// 	pair := strings.SplitN(e, "=", 2)
	// 	fmt.Println(pair[0])
	// }
}

// helper to get context and client
func getContext() (*tfe.Client, context.Context) {
	config := &tfe.Config{
		Address: "https://" + tfeHostname,
		Token:   tfeToken,
	}

	client, err := tfe.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// Create a context
	ctx := context.Background()

	return client, ctx
}
