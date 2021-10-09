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
	"os"

	"github.com/manifoldco/promptui"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a different kubectl version.",
	Args:  cobra.OnlyValidArgs,
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: setKubectlVersion,
}

var Constraint string
var srcPath string

func init() {
	homeDir, _ := os.UserHomeDir()
	srcPath = fmt.Sprintf("%s/.kctlswitch/bin/", homeDir)
	rootCmd.AddCommand(setCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	setCmd.Flags().StringVarP(&Constraint, "constraint", "c", "", "The kubectl semver constraint to use.")
	setCmd.Flags().BoolP("use-latest", "l", false, "Use the latest version in constraint, if results are more than one.")

	setCmd.MarkFlagRequired("constraint")
}

func setKubectlVersion(cmd *cobra.Command, args []string) {
	kctlVersions, err := lib.KctlVersionList(Constraint, slog)
	if err != nil {
		slog.Error(err)
	}
	var result string
	if len(kctlVersions) > 1 {
		prompt := promptui.Select{
			Label: "Select kubectl version",
			Items: kctlVersions,
		}
		_, result, err = prompt.Run()
		if err != nil {
			slog.Error(err)
		}
	} else {
		result = kctlVersions[0]
	}
	lib.DownloadKctl(fmt.Sprintf("v%s", result), srcPath, slog)
	// TODO
	lib.InstallKctlVersion(result, srcPath, rootCmd.PersistentFlags().Lookup("path").Value.String(), slog)
}
