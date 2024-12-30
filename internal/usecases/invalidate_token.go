package usecases

import (
	"github.com/nmarsollier/authgo/internal/token"
	"github.com/nmarsollier/commongo/log"
	"github.com/nmarsollier/commongo/rbt"
)

type InvalidateTokenUseCase interface {
	InvalidateToken(token string) error
}

func NewInvalidateTokenUseCase(
	log log.LogRusEntry,
	tokenService token.TokenService,
	sendLogout rbt.RabbitPublisher[string],
) InvalidateTokenUseCase {
	return &invalidateTokenUseCase{
		log:          log,
		tokenService: tokenService,
		rabbit:       sendLogout,
	}
}

type invalidateTokenUseCase struct {
	log          log.LogRusEntry
	tokenService token.TokenService
	rabbit       rbt.RabbitPublisher[string]
}

func (s *invalidateTokenUseCase) InvalidateToken(token string) error {
	err := s.tokenService.Invalidate(token)

	if err != nil {
		return err
	}

	go func() {
		if s.rabbit == nil {
			s.log.Info("Rabbit logout no se pudo enviar")
			return
		}

		if err = s.rabbit.Publish(
			"Bearer " + token); err != nil {
			s.log.Info("Rabbit logout no se pudo enviar")
		}
	}()

	return nil
}
