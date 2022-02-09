package logging

import (
	"context"

	"go.uber.org/zap"
)

type loggerKeyType int

const loggerKey loggerKeyType = iota

var rootLogger *zap.Logger
var rootSugaredLogger *zap.SugaredLogger

func init() {
	rootLogger, _ = zap.NewProduction()
	rootSugaredLogger = rootLogger.Sugar()
	defer rootSugaredLogger.Sync()
}

func NewContext(ctx context.Context, fields ...interface{}) context.Context {
	return context.WithValue(ctx, loggerKey, WithContext(ctx).With(fields...))
}

func WithContext(ctx context.Context) *zap.SugaredLogger {
	if ctx == nil {
		return rootSugaredLogger
	}
	if ctxLogger, ok := ctx.Value(loggerKey).(*zap.SugaredLogger); ok {
		return ctxLogger
	} else {
		return rootSugaredLogger
	}
}
