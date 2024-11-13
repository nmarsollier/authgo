package token

import (
	"github.com/nmarsollier/authgo/rabbit"
	"github.com/nmarsollier/authgo/tools/errs"
	"github.com/nmarsollier/authgo/tools/log"
)

// Invalidate invalida un token
func Invalidate(tokenString string, ctx ...interface{}) error {
	token, err := Validate(tokenString, ctx...)
	if err != nil {
		return errs.Unauthorized
	}

	if err = delete(token.ID, ctx...); err != nil {
		return err
	}

	cacheRemove(token)

	go func() {
		if err = rabbit.SendLogout("Bearer "+tokenString, ctx...); err != nil {
			log.Get(ctx...).Info("Rabbit logout no se pudo enviar")
		}
	}()

	return nil
}
