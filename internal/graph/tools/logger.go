package tools

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/nmarsollier/authgo/internal/engine/env"
	"github.com/nmarsollier/authgo/internal/engine/log"
	uuid "github.com/satori/go.uuid"
)

func newLogger(ctx context.Context) log.LogRusEntry {
	operationContext := graphql.GetOperationContext(ctx)

	return log.Get(env.Get().FluentUrl).
		WithField(log.LOG_FIELD_CORRELATION_ID, getCorrelationId(ctx)).
		WithField(log.LOG_FIELD_CONTROLLER, "GraphQL").
		WithField(log.LOG_FIELD_HTTP_METHOD, operationContext.OperationName).
		WithField(log.LOG_FIELD_HTTP_PATH, operationContext.OperationName)
}

func getCorrelationId(ctx context.Context) string {
	operationContext := graphql.GetOperationContext(ctx)
	value := operationContext.Headers.Get("Authorization")

	if len(value) == 0 {
		value = uuid.NewV4().String()
	}

	return value
}
