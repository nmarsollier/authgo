package usecases

import (
	"github.com/nmarsollier/authgo/internal/common/log"
	"github.com/nmarsollier/authgo/internal/token"
)

func InvalidateToken(
	log log.LogRusEntry,
	tokenString string,
) error {
	err := token.Invalidate(log, tokenString)

	if err != nil {
		return err
	}

	go func() {
		rabbit := sendLogoutPublisher(log)
		if rabbit == nil {
			log.Info("Rabbit logout no se pudo enviar")
			return
		}

		if err = rabbit.Publish("Bearer " + tokenString); err != nil {
			log.Info("Rabbit logout no se pudo enviar")
		}
	}()

	return nil
}
