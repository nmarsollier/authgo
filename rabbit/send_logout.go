package rabbit

import (
	"encoding/json"

	"github.com/golang/glog"
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
	send := message{
		Type:    "logout",
		Message: token,
	}

	chanel, err := getChannel(ctx...)
	if err != nil {
		glog.Error(err)
		return err
	}

	err = chanel.ExchangeDeclare(
		"auth",   // name
		"fanout", // type
	)
	if err != nil {
		glog.Error(err)
		return err
	}

	body, err := json.Marshal(send)
	if err != nil {
		glog.Error(err)
		return err
	}

	err = chanel.Publish(
		"auth", // exchange
		"",     // routing key
		body)
	if err != nil {
		glog.Error(err)
		return err
	}

	glog.Info("Rabbit logout enviado", send)
	return nil
}

type message struct {
	Type    string `json:"type" example:"logout" `
	Message string `json:"message" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbklEIjoiNjZiNjBlYzhlMGYzYzY4OTUzMzJlOWNmIiwidXNlcklEIjoiNjZhZmQ3ZWU4YTBhYjRjZjQ0YTQ3NDcyIn0.who7upBctOpmlVmTvOgH1qFKOHKXmuQCkEjMV3qeySg" `
}