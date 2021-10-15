package main

import (
	"fmt"
	"kctlswitch/cmd"
	"os"
)

// TODO: implement proper input request

var srcPath string

func init() {
	homeDir, _ := os.UserHomeDir()
	srcPath = fmt.Sprintf("%s/.kctlswitch/bin/", homeDir)
}

func main() {
	// defer logger.Sync() // flushes buffer, if any
	cmd.Execute()

	//lib.AltFetchTags(slog)
	/* kctlVersions, err := lib.KctlVersionList(constraint, slog)
	if err != nil {
		slog.Error(err)
	}

	prompt := promptui.Select{
		Label: "Select kubectl version",
		Items: kctlVersions,
	}
	_, result, err := prompt.Run()
	if err != nil {
		slog.Error(err)
	}

	lib.DownloadKctl(fmt.Sprintf("v%s", result), srcPath, slog)
	lib.InstallKctlVersion("1.17.9", srcPath, "/usr/local/bin/", slog) */
}
