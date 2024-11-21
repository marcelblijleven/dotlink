/*
Copyright © 2024 Marcel Blijleven <marcelblijleven@gmail.com>

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

	"github.com/marcelblijleven/dotlink/pkg/utils"
	"github.com/spf13/cobra"
)

var dryRun bool

// linkCmd represents the link command
var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Create a symbolic link for all configured files in the current directory",
	Long: `Link will iterate over all (nested) files in the current directory and Create
  a symbolic link to the home directory.

  If a file matches an ignore pattern, it will not be linked.`,
	Run: func(cmd *cobra.Command, args []string) {
		// preflight := utils.Preflight{
		// 	SymlinkActions: []utils.SymlinkAction{},
		// 	DirActions:     []utils.DirAction{},
		// 	Errors:         []error{},
		// }
		// ignorePatterns := utils.CompilePatterns(configuration.IgnorePatterns())

		files, err := utils.FindFiles(source, configuration)
		cobra.CheckErr(err)

		for _, file := range files {
			fmt.Println(file)
		}

		// for _, sa := range preflight.SymlinkActions {
		// 	fmt.Println(sa.String(source))
		// }
	},
}

func init() {
	rootCmd.AddCommand(linkCmd)
	linkCmd.Flags().BoolVar(&dryRun, "dry-run", false, "")
}
