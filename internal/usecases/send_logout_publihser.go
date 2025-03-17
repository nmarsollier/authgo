package usecases

import (
	"github.com/nmarsollier/authgo/internal/common/log"
	"github.com/nmarsollier/authgo/internal/common/rbt"
	"github.com/nmarsollier/authgo/internal/env"
)

var currSendLogout rbt.RabbitPublisher[string]

func sendLogoutPublisher(
	log log.LogRusEntry,
) rbt.RabbitPublisher[string] {
	currSendLogout, _ = rbt.NewRabbitPublisher[string](
		rbt.RbtLogger(env.Get().FluentURL, env.Get().ServerName, log.CorrelationId()),
		env.Get().RabbitURL,
		"auth",
		"fanout",
		"",
	)

	return currSendLogout
}
