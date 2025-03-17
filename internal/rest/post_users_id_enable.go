package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/internal/rest/server"
	"github.com/nmarsollier/authgo/internal/user"
)

//	@Summary		Enable User
//	@Description	Habilita un usuario en el sistema. El usuario logueado debe tener permisos "admin".
//	@Tags			Seguridad
//	@Accept			json
//	@Produce		json
//	@Param			userId			path	string	true	"ID del usuario a habilitar"
//	@Param			Authorization	header	string	true	"Bearer {token}"
//	@Success		200				"No Content"
//	@Failure		400				{object}	errs.ValidationErr	"Bad Request"
//	@Failure		401				{object}	server.ErrorData	"Unauthorized"
//	@Failure		404				{object}	server.ErrorData	"Not Found"
//	@Failure		500				{object}	server.ErrorData	"Internal Server Error"
//	@Router			/users/:userId/enable [post]
//
// Habilita un usuario en el sistema.
func postUsersIdEnableRoute(engine *gin.Engine) {
	engine.POST(
		"/users/:userID/enable",
		server.IsAdminMiddleware,
		enable,
	)
}

func enable(c *gin.Context) {
	userId := c.Param("userID")
	log := server.GinLogger(c)

	if err := user.Enable(log, userId); err != nil {
		server.AbortWithError(c, err)
	}
}
