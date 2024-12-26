package graph

import (
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/nmarsollier/authgo/internal/engine/env"
	"github.com/nmarsollier/authgo/internal/engine/log"
	"github.com/nmarsollier/authgo/internal/graph/model"
	"github.com/nmarsollier/authgo/internal/graph/schema"
)

func Start() {
	logger := log.Get(env.Get().FluentUrl)
	port := env.Get().GqlPort
	srv := handler.NewDefaultServer(model.NewExecutableSchema(model.Config{Resolvers: &schema.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	logger.Info("GraphQL playground in port : ", port)
	logger.Error(http.ListenAndServe(fmt.Sprintf(":%d", env.Get().GqlPort), nil))
}
