package main

import (
	"fmt"
	"kctlswitch/lib"
	"os"

	"github.com/manifoldco/promptui"
	"go.uber.org/zap"
)

// FEAT: improve logging with optional logfile and a better console experience
var logger, _ = zap.NewProduction()
var slog = logger.Sugar()

// TODO: implement proper input request
const (
	constraint string = "1.17.8 - 1.17.20"
)

var srcPath string

func init() {
	homeDir, _ := os.UserHomeDir()
	srcPath = fmt.Sprintf("%s/.kctlswitch/bin/", homeDir)
}

func main() {
	defer logger.Sync() // flushes buffer, if any

	slog.Infof("Starting run with constraint %s", constraint)

	//lib.AltFetchTags(slog)
	kctlVersions, err := lib.KctlVersionList(constraint, slog)
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
	lib.InstallKctlVersion("1.17.9", srcPath, "/usr/local/bin/", slog)
}
