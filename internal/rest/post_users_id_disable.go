package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/internal/rest/server"
)

//	@Summary		Deshabilitar Usuario
//	@Description	Deshabilita un usuario en el sistema. El usuario logueado debe tener permisos "admin".
//	@Tags			Seguridad
//	@Accept			json
//	@Produce		json
//	@Param			userId			path	string	true	"ID del usuario a deshabilitar"
//	@Param			Authorization	header	string	true	"Bearer {token}"
//	@Success		200				"No Content"
//	@Failure		400				{object}	errs.ValidationErr	"Bad Request"
//	@Failure		401				{object}	server.ErrorData	"Unauthorized"
//	@Failure		404				{object}	server.ErrorData	"Not Found"
//	@Failure		500				{object}	server.ErrorData	"Internal Server Error"
//	@Router			/users/:userId/disable [post]
//
// Deshabilitar Usuario
func postUsersIdDisableRoute(engine *gin.Engine) {
	engine.POST(
		"/users/:userID/disable",
		server.IsAdminMiddleware,
		disable,
	)
}

func disable(c *gin.Context) {
	userId := c.Param("userID")

	di := server.GinDi(c)
	if err := di.UserService().Disable(userId); err != nil {
		server.AbortWithError(c, err)
	}
}
