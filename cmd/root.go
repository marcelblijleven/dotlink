/*
Copyright © 2024 Marcel Blijleven

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
	"os"

	"github.com/marcelblijleven/dotlink/pkg/config"
	"github.com/spf13/cobra"
)

var (
	cfgFile       string
	target        string
	source        string
	ignoreExtra   []string
	configuration config.Config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dotlink",
	Short: "Create symbolic links to your home dir",
	Long: `Create symbolic links for every file
  and directory inside the directory.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .dotlink in the current directory)")
	rootCmd.PersistentFlags().StringVar(&target, "target", "", "target location for the symlinks. Defaults to the user's home dir")
	rootCmd.PersistentFlags().StringVar(&source, "source", "", "source location, defaults to the current directory")
	rootCmd.PersistentFlags().StringArrayVar(&ignoreExtra, "ignore", []string{}, "provide a pattern to ignore, you can provide this flag multiple times")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	config.InitConfig(cfgFile, &target, ignoreExtra, &configuration)
	fmt.Printf("Using target: %s\n\n", target)
}
