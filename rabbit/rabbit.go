package rabbit

import (
	"github.com/nmarsollier/authgo/engine/env"
	"github.com/nmarsollier/authgo/engine/log"
	"github.com/streadway/amqp"
)

func NewChannel(log log.LogRusEntry) (RabbitChannel, error) {
	conn, err := amqp.Dial(env.Get().RabbitURL)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return rabbitChannel{ch: channel}, nil
}
