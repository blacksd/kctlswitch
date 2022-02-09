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
	"errors"
	"fmt"
	"io"
	"kctlswitch/lib"
	"kctlswitch/logging"
	"net/http"
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
	RunE: setKubectlVersion,
}

var constraint string
var srcPath string
var useLatestVersion bool
var noVerify bool
var forceOverwrite bool

func init() {
	rootCmd.AddCommand(setCmd)

	homeDir, _ := os.UserHomeDir()
	srcPath = fmt.Sprintf("%s/.kctlswitch/bin/", homeDir)

	setCmd.Flags().StringVarP(&constraint, "constraint", "c", "", "The kubectl semver constraint to use.")
	setCmd.Flags().BoolVarP(&useLatestVersion, "use-latest", "l", false, "Use the latest version in constraint.")
	setCmd.Flags().BoolVarP(&noVerify, "no-verify", "n", false, "Skip checking the version's hash.")
	setCmd.Flags().BoolVarP(&forceOverwrite, "force", "f", false, "Overwrite a non-symlinked 'kubectl' existing in the destination.")
}

func setKubectlVersion(cmd *cobra.Command, args []string) error {

	setLogger := logging.WithContext(cmd.Context())
	setLogger.Debug("set subcommand invoked")

	if constraint == "" {
		if useLatestVersion {
			defaultConstraint, err := fetchKubernetesStableVersion()
			if err != nil {
				return errors.New("can't fetch the latest stable version")
			}
			constraint = defaultConstraint
		} else {
			return errors.New("no constraints and no latest, I have no clue what you want")
		}

	}

	kctlVersions, err := lib.KctlVersionList(constraint, setLogger)
	if err != nil {
		setLogger.Fatal(err)
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
			setLogger.Fatal(err)
		}
	}
	lib.DownloadKctl(fmt.Sprintf("v%s", result), srcPath, noVerify)
	lib.InstallKctlVersion(result, srcPath, rootCmd.PersistentFlags().Lookup("bin").Value.String(), forceOverwrite)
	return nil
}

func fetchKubernetesStableVersion() (string, error) {
	resp, err := http.DefaultClient.Get("https://dl.k8s.io/release/stable.txt")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		bodyString := string(bodyBytes)
		return bodyString, nil
	}
	return "", errors.New("unexpected error")
}
