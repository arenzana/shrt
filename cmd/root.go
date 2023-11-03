// Copyright © 2023 Ismael Arenzana <isma@arenzana.org>

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ApplicationVersion = ""

var (
	cfgFile string
	quiet   bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "shrt",
	Version: ApplicationVersion,
	Short:   "An alternative Shlink service application to shorten URLs",
	Long:    `An alternative Shlink service application to interact with a Shlink instance without having to install the official client and server.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.SetVersionTemplate("Version: ")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.shrt.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", true, "config file (default is $HOME/.shrt.yaml)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".shrt" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".shrt")
		viper.SetConfigType("yaml")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		_ = fmt.Errorf("Error reading config file: %s", viper.ConfigFileUsed())
	}

	// Environment variables take precedence over config options
	viper.SetEnvPrefix("shrt")
	viper.AutomaticEnv() // read in environment variables that match
	_ = viper.BindEnv("shlink_url")
	_ = viper.BindEnv("api_key")
	_ = viper.BindEnv("timeout")

	if viper.Get("api_key") == "" {
		_ = fmt.Errorf("API Key not set")
		os.Exit(1)
	}
	// Set defaults
	viper.SetDefault("shlink_url", "https://shlink.io")
	viper.SetDefault("timeout", 10)

}
