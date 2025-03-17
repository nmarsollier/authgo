package tools

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/nmarsollier/authgo/internal/common/log"
)

func GqlLogger(c context.Context) log.LogRusEntry {
	operationContext := graphql.GetOperationContext(c)
	context_deps, exist := operationContext.Variables["logger"]
	if exist {
		return context_deps.(log.LogRusEntry)
	}

	logger := newLogger(c)
	operationContext.Variables["logger"] = logger
	return logger
}
