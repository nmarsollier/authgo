package tools

import (
	"context"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/errs"
	"github.com/nmarsollier/authgo/tools/log"
	"github.com/nmarsollier/authgo/user"
)

func ValidateLoggedIn(ctx context.Context) error {
	_, err := HeaderToken(ctx)
	if err != nil {
		return err
	}
	return nil
}
func ValidateAdmin(ctx context.Context) error {
	token, err := HeaderToken(ctx)
	if err != nil {
		return err
	}

	env := GqlCtx(ctx)
	if !user.Granted(token.UserID.Hex(), "admin", env...) {
		log.Get(env...).Warn("Unauthorized")
		return errs.Unauthorized
	}

	return nil
}

// HeaderToken Token data from Authorization header
func HeaderToken(ctx context.Context) (*token.Token, error) {
	env := GqlCtx(ctx)

	tokenString, err := TokenString(ctx)
	if err != nil {
		return nil, errs.Unauthorized
	}

	payload, err := token.Validate(tokenString, env...)
	if err != nil {
		log.Get(env...).Error(err)
		return nil, err
	}

	return payload, nil
}

// HeaderToken Token data from Authorization header
func TokenString(ctx context.Context) (string, error) {
	operationContext := graphql.GetOperationContext(ctx)
	tokenString := operationContext.Headers.Get("Authorization")

	if strings.Index(strings.ToUpper(tokenString), "BEARER ") == 0 {
		tokenString = tokenString[7:]
	} else {
		return "", errs.Unauthorized
	}

	return tokenString, nil
}
