package logging

import (
	"context"

	"go.uber.org/zap"
)

type loggerKeyType int

const loggerKey loggerKeyType = iota

var logger *zap.Logger
var sugaredLogger *zap.SugaredLogger

func init() {
	logger, _ = zap.NewProduction()
	sugaredLogger = logger.Sugar()
}

func NewContext(ctx context.Context, fields ...zap.Field) context.Context {
	return context.WithValue(ctx, loggerKey, WithContext(ctx).With(fields...))
}

func WithContext(ctx context.Context) *zap.SugaredLogger {
	if ctx == nil {
		return sugaredLogger
	}
	if ctxLogger, ok := ctx.Value(loggerKey).(zap.SugaredLogger); ok {
		return &ctxLogger
	} else {
		return sugaredLogger
	}
}

// USAGE:

/*
From a called func:

logging.WithContext(ctx).Info("This is a cane log!",zap.String("hey", variable))

From the caller building the context:


var myLocalContext = context.Background()
var myFunctionStuff string = "name"

reqCtx := logging.NewContext(myLocalContext, zap.String("context_name", myFunctionStuff))
ctxLogger := logging.WithContext(reqCtx)

ctxLogger.Info("Cane!")
*/
