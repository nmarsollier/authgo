package token

import (
	"github.com/nmarsollier/authgo/rabbit"
	"github.com/nmarsollier/authgo/tools/errs"
	"github.com/nmarsollier/authgo/tools/log"
)

// Invalidate invalida un token
func Invalidate(tokenString string, deps ...interface{}) error {
	token, err := Validate(tokenString, deps...)
	if err != nil {
		return errs.Unauthorized
	}

	if err = delete(token.ID, deps...); err != nil {
		return err
	}

	cacheRemove(token)

	go func() {
		if err = rabbit.SendLogout("Bearer "+tokenString, deps...); err != nil {
			log.Get(deps...).Info("Rabbit logout no se pudo enviar")
		}
	}()

	return nil
}
