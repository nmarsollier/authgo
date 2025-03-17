package tools

import (
	"context"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/nmarsollier/authgo/internal/common/errs"
	"github.com/nmarsollier/authgo/internal/token"
	"github.com/nmarsollier/authgo/internal/user"
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

	log := GqlLogger(ctx)
	if !user.Granted(log, token.UserID.Hex(), "admin") {
		log.Warn("Unauthorized")
		return errs.Unauthorized
	}

	return nil
}

// HeaderToken Token data from Authorization header
func HeaderToken(ctx context.Context) (*token.Token, error) {
	log := GqlLogger(ctx)

	tokenString, err := TokenString(ctx)
	if err != nil {
		return nil, errs.Unauthorized
	}

	payload, err := token.Validate(log, tokenString)
	if err != nil {
		log.Error(err)
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
