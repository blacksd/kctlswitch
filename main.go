package main

import (
	"kctlswitch/lib"

	"go.uber.org/zap"
)

var logger, _ = zap.NewProduction()
var slog = logger.Sugar()

func main() {
	defer logger.Sync() // flushes buffer, if any
	// lib.KctlVersionsList("<= 1.7")
	slog.Info("test")
	lib.DownloadKctl("v1.12.3", "./", slog)
}
