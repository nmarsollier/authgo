package rabbit

import (
	"github.com/nmarsollier/authgo/tools/env"
	"github.com/nmarsollier/authgo/tools/log"
	"github.com/streadway/amqp"
)

func getChannel(ctx ...interface{}) (RabbitChannel, error) {
	for _, o := range ctx {
		if ti, ok := o.(RabbitChannel); ok {
			return ti, nil
		}
	}

	conn, err := amqp.Dial(env.Get().RabbitURL)
	if err != nil {
		log.Get(ctx...).Error(err)
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Get(ctx...).Error(err)
		return nil, err
	}

	return rabbitChannel{ch: channel}, nil
}
