package rabbit

import (
	"encoding/json"

	"github.com/nmarsollier/authgo/tools/log"
)

//	@Summary		Mensage Rabbit
//	@Description	SendLogout envía un broadcast a rabbit con logout. Esto no es Rest es RabbitMQ.
//	@Tags			Rabbit
//	@Accept			json
//	@Produce		json
//	@Param			body	body	message	true	"Token deshabilitado"
//	@Router			/rabbit/logout [put]
//
// SendLogout envía un broadcast a rabbit con logout
func SendLogout(token string, ctx ...interface{}) error {
	logger := log.Get(ctx...).
		WithField(log.LOG_FIELD_CONTROLLER, "Rabbit").
		WithField(log.LOG_FIELD_RABBIT_ACTION, "Emit").
		WithField(log.LOG_FIELD_RABBIT_EXCHANGE, "auth").
		WithField(log.LOG_FIELD_RABBIT_QUEUE, "logout")

	corrId, _ := logger.Data()[log.LOG_FIELD_CORRELATION_ID].(string)
	send := message{
		CorrelationId: corrId,
		Message:       token,
	}

	ctx = append(ctx, logger)
	chanel, err := getChannel(ctx...)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = chanel.ExchangeDeclare(
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

	err = chanel.Publish(
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
