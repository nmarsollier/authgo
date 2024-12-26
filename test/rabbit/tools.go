package rabbit

import (
	"github.com/golang/mock/gomock"
	"github.com/nmarsollier/authgo/test/mock"
)

func DefaultMockRabbitChannel(ctrl *gomock.Controller, times int) *mock.MockRabbitChannel {
	channel := mock.NewMockRabbitChannel(ctrl)
	channel.EXPECT().ExchangeDeclare(gomock.Any(), gomock.Any()).Return(nil).Times(times)
	return channel
}
