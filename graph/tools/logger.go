package tools

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/nmarsollier/authgo/tools/log"
	uuid "github.com/satori/go.uuid"
)

func newLogger(ctx context.Context, env ...interface{}) log.LogRusEntry {
	operationContext := graphql.GetOperationContext(ctx)

	return log.Get(env...).
		WithField(log.LOG_FIELD_CORRELATION_ID, getCorrelationId(ctx)).
		WithField(log.LOG_FIELD_CONTROLLER, "Rest").
		WithField(log.LOG_FIELD_HTTP_METHOD, operationContext.OperationName).
		WithField(log.LOG_FIELD_HTTP_PATH, operationContext.OperationName)
}

func gqlLogger(ctx context.Context) log.LogRusEntry {
	operationContext := graphql.GetOperationContext(ctx)

	logger, exist := operationContext.Variables["logger"]
	if !exist {
		return newLogger(ctx)
	}
	return logger.(log.LogRusEntry)
}

func getCorrelationId(ctx context.Context) string {
	operationContext := graphql.GetOperationContext(ctx)
	value := operationContext.Headers.Get("Authorization")

	if len(value) == 0 {
		value = uuid.NewV4().String()
	}

	return value
}
