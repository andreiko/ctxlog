package ctxlog

import (
	"context"
	"github.com/apex/log"
)

const LogFieldsContextKey = "log_fields"

func GetContextualLogger(context context.Context, logger log.Interface) log.Interface {
	ctxValue := context.Value(LogFieldsContextKey)
	if ctxValue == nil {
		return logger
	}

	ctxFields, ok := ctxValue.(log.Fields)
	if !ok {
		return logger
	}

	return logger.WithFields(ctxFields)
}

func GetUpdatedLoggingContext(ctx context.Context, fields log.Fields) context.Context {
	if len(fields) < 1 {
		return ctx
	}

	newFields := make(log.Fields)

	ctxValue := ctx.Value(LogFieldsContextKey)
	if ctxValue != nil {
		if currentFields, ok := ctxValue.(log.Fields); ok {
			for name, value := range currentFields {
				newFields[name] = value
			}
		}
	}

	for name, value := range fields {
		newFields[name] = value
	}

	return context.WithValue(ctx, LogFieldsContextKey, newFields)
}
