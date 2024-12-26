package rabbit

import (
	"encoding/json"

	"github.com/nmarsollier/authgo/internal/engine/log"
	"github.com/nmarsollier/authgo/internal/engine/rbt"
)

type SendLogoutService interface {
	SendLogout(token string) error
}

type sendLogoutService struct {
	channel rbt.RabbitChannel
	log     log.LogRusEntry
}

func NewSendLogoutService(
	log log.LogRusEntry,
	channel rbt.RabbitChannel,
) (SendLogoutService, error) {
	return &sendLogoutService{
		channel: channel,
		log:     log,
	}, nil
}

//	@Summary		Mensage Rabbit
//	@Description	SendLogout envía un broadcast a rabbit con logout. Esto no es Rest es RabbitMQ.
//	@Tags			Rabbit
//	@Accept			json
//	@Produce		json
//	@Param			body	body	message	true	"Token deshabilitado"
//	@Router			/rabbit/logout [put]
//
// SendLogout envía un broadcast a rabbit con logout
func (s *sendLogoutService) SendLogout(token string) error {
	logger := s.log.
		WithField(log.LOG_FIELD_CONTROLLER, "Rabbit").
		WithField(log.LOG_FIELD_RABBIT_ACTION, "Emit").
		WithField(log.LOG_FIELD_RABBIT_EXCHANGE, "auth").
		WithField(log.LOG_FIELD_RABBIT_QUEUE, "logout")

	corrId, _ := logger.Data()[log.LOG_FIELD_CORRELATION_ID].(string)
	send := message{
		CorrelationId: corrId,
		Message:       token,
	}

	err := s.channel.ExchangeDeclare(
		"auth",   // name
		"fanout", // type
	)
	if err != nil {
		logger.Error(err)
		return err
	}

	body, err := json.Marshal(send)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = s.channel.Publish(
		"auth", // exchange
		"",     // routing key
		body)
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Info("Rabbit logout enviado", string(body))
	return nil
}

type message struct {
	CorrelationId string `json:"correlation_id" example:"123123" `
	Message       string `json:"message" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbklEIjoiNjZiNjBlYzhlMGYzYzY4OTUzMzJlOWNmIiwidXNlcklEIjoiNjZhZmQ3ZWU4YTBhYjRjZjQ0YTQ3NDcyIn0.who7upBctOpmlVmTvOgH1qFKOHKXmuQCkEjMV3qeySg" `
}
