package tools

import (
	"context"
)

func GqlCtx(c context.Context) []interface{} {

	var ctx []interface{}

	ctx = append(ctx, gqlLogger(c))

	return ctx
}
