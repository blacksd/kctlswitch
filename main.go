package main

import (
	"fmt"
	"kctlswitch/lib"

	"github.com/manifoldco/promptui"
	"go.uber.org/zap"
)

// FEAT: improve logging with optional logfile and a better console experience
var logger, _ = zap.NewProduction()
var slog = logger.Sugar()

const (
	constraint string = "1.17.5 - 1.17.9"
)

func main() {
	defer logger.Sync() // flushes buffer, if any

	slog.Infof("Starting run with constraint %s", constraint)

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

	lib.DownloadKctl(fmt.Sprintf("v%s", result), "./", slog)
}
