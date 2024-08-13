package tests

import (
	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/rabbit"
)

func MockRabbitChannel(ctrl *gomock.Controller, times int) *rabbit.MockRabbitChannel {
	channel := rabbit.NewMockRabbitChannel(ctrl)
	channel.EXPECT().ExchangeDeclare(gomock.Any(), gomock.Any()).Return(nil).Times(times)
	return channel
}
