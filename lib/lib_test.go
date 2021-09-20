package lib_test

import (
	"go.uber.org/zap"
)

var logger, _ = zap.NewProduction()
var slog = logger.Sugar()
