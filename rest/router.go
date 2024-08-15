package rest

import (
	"fmt"

	"github.com/nmarsollier/authgo/rest/server"
	"github.com/nmarsollier/authgo/tools/env"
)

// Start this server
func Start() {
	InitRoutes()
	server.Router().Run(fmt.Sprintf(":%d", env.Get().Port))
}

func InitRoutes() {
	getUserSignOutRoute()
	getUsersCurrentRoute()
	getUsersRoute()
	getUserPasswordRoute()
	postUserSignInRoute()
	postUsersRoute()
	postUsersIdDisableRoute()
	postUsersIdEnableRoute()
	postUsersIdGrantRoute()
	postUsersIdRevokeRoute()
}
