package rabbit

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/nmarsollier/authgo/tools/env"
	"github.com/streadway/amqp"
)

// ErrChannelNotInitialized Rabbit channel could not be initialized
var ErrChannelNotInitialized = errors.New("Channel not initialized")

var channel *amqp.Channel

type message struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func get() (*amqp.Channel, error) {
	if channel == nil {
		conn, err := amqp.Dial(env.Get().RabbitURL)
		if err != nil {
			return nil, err
		}

		ch, err := conn.Channel()
		if err != nil {
			return nil, err
		}
		channel = ch
	}
	if channel == nil {
		return nil, ErrChannelNotInitialized
	}
	return channel, nil
}

/**
 * @api {fanout} auth/fanout Invalidar Token
 * @apiGroup RabbitMQ POST
 *
 * @apiDescription AuthService enviá un broadcast a todos los usuarios cuando un token ha sido invalidado. Los clientes deben eliminar de sus caches las sesiones invalidadas.
 *
 * @apiSuccessExample {json} Mensaje
 *     {
 *        "type": "logout",
 *        "message": "{Token revocado}"
 *     }
 */
func SendLogout(token string) error {
	send := message{
		Type:    "logout",
		Message: token,
	}

	chanel, err := get()
	if err != nil {
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
		channel = nil
		return err
	}

	body, err := json.Marshal(send)
	if err != nil {
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
		channel = nil
		return err
	}

	log.Output(1, "Rabbit logout enviado")
	return nil
}
