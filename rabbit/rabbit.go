package rabbit

import (
	"encoding/json"

	"github.com/golang/glog"
	"github.com/nmarsollier/authgo/tools/app_errors"
	"github.com/nmarsollier/authgo/tools/env"
	"github.com/streadway/amqp"
)

// ErrChannelNotInitialized Rabbit channel could not be initialized
var ErrChannelNotInitialized = app_errors.NewCustom(400, "Channel not initialized")

var channel *amqp.Channel

type Rabbit interface {
	SendLogout(token string) error
}

type rabbitImpl struct {
}

var currentRabbit Rabbit

func Get(options ...interface{}) Rabbit {
	for _, o := range options {
		if ti, ok := o.(Rabbit); ok {
			return ti
		}
	}

	if currentRabbit == nil {
		currentRabbit = &rabbitImpl{}
	}
	return currentRabbit
}

type message struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func getChannel() (*amqp.Channel, error) {
	if channel == nil {
		conn, err := amqp.Dial(env.Get().RabbitURL)
		if err != nil {
			glog.Error(err)
			return nil, err
		}

		ch, err := conn.Channel()
		if err != nil {
			glog.Error(err)
			return nil, err
		}
		channel = ch
	}
	if channel == nil {
		return nil, ErrChannelNotInitialized
	}
	return channel, nil
}

// SendLogout envía un broadcast a rabbit con logout
//
//	@Summary		Mensage Rabbit
//	@Description	SendLogout envía un broadcast a rabbit con logout. Esto no es Rest es RabbitMQ.
//	@Tags			Rabbit
//	@Accept			json
//	@Produce		json
//	@Param			body	body	message	true	"Token deshabilitado"
//	@Router			/rabbit/logout [put]
func (r *rabbitImpl) SendLogout(token string) error {
	send := message{
		Type:    "logout",
		Message: token,
	}

	chanel, err := getChannel()
	if err != nil {
		glog.Error(err)
		channel = nil
		return err
	}

	err = chanel.ExchangeDeclare(
		"auth",   // name
		"fanout", // type
		false,    // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		glog.Error(err)
		channel = nil
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
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			Body: []byte(body),
		})
	if err != nil {
		glog.Error(err)
		channel = nil
		return err
	}

	glog.Info("Rabbit logout enviado", send)
	return nil
}
