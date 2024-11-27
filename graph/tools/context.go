package tools

import (
	"context"
)

func GqlDeps(c context.Context) []interface{} {

	var deps []interface{}

	deps = append(deps, gqlLogger(c))

	return deps
}
