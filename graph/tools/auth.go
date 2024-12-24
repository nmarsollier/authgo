package tools

import (
	"context"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/nmarsollier/authgo/engine/errs"
	"github.com/nmarsollier/authgo/token"
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

	di := GqlDi(ctx)
	if !di.UserService().Granted(token.UserID.Hex(), "admin") {
		di.Logger().Warn("Unauthorized")
		return errs.Unauthorized
	}

	return nil
}

// HeaderToken Token data from Authorization header
func HeaderToken(ctx context.Context) (*token.Token, error) {
	di := GqlDi(ctx)

	tokenString, err := TokenString(ctx)
	if err != nil {
		return nil, errs.Unauthorized
	}

	payload, err := di.TokenService().Validate(tokenString)
	if err != nil {
		di.Logger().Error(err)
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
