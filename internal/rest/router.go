package rest

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/internal/engine/env"
	"github.com/nmarsollier/authgo/internal/rest/server"
)

// Start this server
func Start() {
	InitRoutes(server.Router())
	server.Router().Run(fmt.Sprintf(":%d", env.Get().Port))
}

func InitRoutes(engine *gin.Engine) {
	getUserSignOutRoute(engine)
	getUsersCurrentRoute(engine)
	getUsersRoute(engine)
	getUserPasswordRoute(engine)
	postUserSignInRoute(engine)
	postUsersRoute(engine)
	postUsersIdDisableRoute(engine)
	postUsersIdEnableRoute(engine)
	postUsersIdGrantRoute(engine)
	postUsersIdRevokeRoute(engine)
}
