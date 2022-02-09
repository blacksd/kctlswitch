package lib

import (
	"context"
	"kctlswitch/logging"

	"go.uber.org/zap"
)

var libCtx = logging.NewContext(context.Background(), "context", "lib")
var libLogger *zap.SugaredLogger

func init() {
	libLogger = logging.WithContext(libCtx)
	defer libLogger.Sync()
}
