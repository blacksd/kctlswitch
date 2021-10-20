/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"kctlswitch/lib"
	"kctlswitch/logging"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a different kubectl version.",
	Args:  cobra.OnlyValidArgs,
	Long: `Download, verify and set a kubectl version.
	
	You can from the constraint`,
	Run: setKubectlVersion,
}

var constraint string
var srcPath string
var useLatestVersion bool
var noVerify bool

func init() {
	rootCmd.AddCommand(setCmd)

	homeDir, _ := os.UserHomeDir()
	srcPath = fmt.Sprintf("%s/.kctlswitch/bin/", homeDir)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	setCmd.Flags().StringVarP(&constraint, "constraint", "c", "", "The kubectl semver constraint to use.")
	setCmd.Flags().BoolVarP(&useLatestVersion, "use-latest", "l", false, "Use the latest version in constraint, if results are more than one.")
	setCmd.Flags().BoolVarP(&noVerify, "no-verify", "n", false, "Skip checking the version's hash.")

	setCmd.MarkFlagRequired("constraint")
}

func setKubectlVersion(cmd *cobra.Command, args []string) {

	myLoggerSet := logging.WithContext(cmd.Context())
	myLoggerSet.Info("set subcommand invoked")

	kctlVersions, err := lib.KctlVersionList(constraint, myLoggerSet)
	if err != nil {
		myLoggerSet.Fatal(err)
	}
	var result string
	if (len(kctlVersions) == 1) || useLatestVersion {
		result = kctlVersions[len(kctlVersions)-1]
	} else {
		prompt := promptui.Select{
			Label: "Select kubectl version",
			Items: kctlVersions,
		}
		_, result, err = prompt.Run()
		if err != nil {
			myLoggerSet.Fatal(err)
		}
	}

	lib.DownloadKctl(fmt.Sprintf("v%s", result), srcPath, noVerify, myLoggerSet)
	lib.InstallKctlVersion(result, srcPath, rootCmd.PersistentFlags().Lookup("path").Value.String(), myLoggerSet)
}
