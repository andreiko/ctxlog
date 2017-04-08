package ctxlog

import (
	"context"
	"github.com/apex/log"
)

const logFieldsContextKey = "log_fields"

// GetContextualLogger takes a context, a logger instance and tries to extend the logging
// instance with fields if they are stored in the context returning the extended logger.
// If context contains no fields, the logger instance is returned unmodified.
func GetContextualLogger(context context.Context, logger log.Interface) log.Interface {
    if logger == nil {
        logger = log.Log
    }

	ctxValue := context.Value(logFieldsContextKey)
	if ctxValue == nil {
		return logger
	}

	ctxFields, ok := ctxValue.(log.Fields)
	if !ok {
		return logger
	}

	return logger.WithFields(ctxFields)
}

// GetContextualLogger takes a context, a logger instance, a map of logging fields and returns
// the updated context that contains the fields and the updated logger. If the context already
// had some fields, they would be merged with the new fields map. The returned logger contains
// the merged map of fields.
func GetUpdatedLoggingContext(ctx context.Context, logger log.Interface, fields log.Fields) (context.Context, log.Interface) {
	newFields := make(log.Fields)

	ctxValue := ctx.Value(logFieldsContextKey)
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

    if logger == nil {
        logger = log.Log
    }

	return context.WithValue(ctx, logFieldsContextKey, newFields), logger.WithFields(newFields)
}
