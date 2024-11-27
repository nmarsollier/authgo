package rabbit

import (
	"github.com/nmarsollier/authgo/tools/env"
	"github.com/nmarsollier/authgo/tools/log"
	"github.com/streadway/amqp"
)

func getChannel(deps ...interface{}) (RabbitChannel, error) {
	for _, o := range deps {
		if ti, ok := o.(RabbitChannel); ok {
			return ti, nil
		}
	}

	conn, err := amqp.Dial(env.Get().RabbitURL)
	if err != nil {
		log.Get(deps...).Error(err)
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Get(deps...).Error(err)
		return nil, err
	}

	return rabbitChannel{ch: channel}, nil
}
