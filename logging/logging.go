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
}

func NewContext(ctx context.Context, fields ...zap.Field) context.Context {
	return context.WithValue(ctx, loggerKey, WithContext(ctx).With(fields))
}

func WithContext(ctx context.Context) *zap.SugaredLogger {
	if ctx == nil {
		return rootSugaredLogger
	}
	if ctxLogger, ok := ctx.Value(loggerKey).(zap.SugaredLogger); ok {
		return &ctxLogger
	} else {
		return rootSugaredLogger
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
