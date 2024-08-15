package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/rest/server"
	"github.com/nmarsollier/authgo/user"
)

//	@Summary		Deshabilitar Usuario
//	@Description	Deshabilita un usuario en el sistema. El usuario logueado debe tener permisos "admin".
//	@Tags			Seguridad
//	@Accept			json
//	@Produce		json
//	@Param			userId			path	string	true	"ID del usuario a deshabilitar"
//	@Param			Authorization	header	string	true	"bearer {token}"
//	@Success		200				"No Content"
//	@Failure		400				{object}	errs.ValidationErr	"Bad Request"
//	@Failure		401				{object}	server.ErrorData	"Unauthorized"
//	@Failure		404				{object}	server.ErrorData	"Not Found"
//	@Failure		500				{object}	server.ErrorData	"Internal Server Error"
//	@Router			/v1/users/:userId/disable [post]
//
// Deshabilitar Usuario
func postUsersIdDisableRoute() {
	server.Router().POST(
		"/v1/users/:userID/disable",
		server.ValidateAdmin,
		disable,
	)
}

func disable(c *gin.Context) {
	userId := c.Param("userID")

	ctx := server.TestCtx(c)
	if err := user.Disable(userId, ctx...); err != nil {
		server.AbortWithError(c, err)
	}
}
