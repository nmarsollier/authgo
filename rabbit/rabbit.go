package rabbit

import (
	"github.com/golang/glog"
	"github.com/nmarsollier/authgo/tools/env"
	"github.com/streadway/amqp"
)

type message struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func getChannel(ctx ...interface{}) (RabbitChannel, error) {
	for _, o := range ctx {
		if ti, ok := o.(RabbitChannel); ok {
			return ti, nil
		}
	}

	conn, err := amqp.Dial(env.Get().RabbitURL)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	return rabbitChannel{ch: channel}, nil
}
