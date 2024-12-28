package tools

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/nmarsollier/authgo/internal/di"
)

func GqlDi(c context.Context) di.Injector {
	operationContext := graphql.GetOperationContext(c)
	context_deps, exist := operationContext.Variables["di"]
	if exist {
		return context_deps.(di.Injector)
	}

	deps := di.NewInjector(newLogger(c))
	operationContext.Variables["di"] = deps
	return deps
}
