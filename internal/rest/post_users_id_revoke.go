package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/internal/rest/server"
)

//	@Summary		Quitar permisos
//	@Description	Quita permisos al usuario indicado, el usuario logueado tiene que tener permiso "admin".
//	@Tags			Seguridad
//	@Accept			json
//	@Produce		json
//	@Param			userId			path	string				true	"ID del usuario a quitar permiso"
//	@Param			Authorization	header	string				true	"Bearer {token}"
//	@Param			body			body	grantPermissionBody	true	"Permisos a Qutiar"
//	@Success		200				"No Content"
//	@Failure		400				{object}	errs.ValidationErr	"Bad Request"
//	@Failure		401				{object}	server.ErrorData	"Unauthorized"
//	@Failure		404				{object}	server.ErrorData	"Not Found"
//	@Failure		500				{object}	server.ErrorData	"Internal Server Error"
//	@Router			/users/:userID/revoke [post]
//
// Quita permisos al usuario indicado.
func postUsersIdRevokeRoute(engine *gin.Engine) {
	engine.POST(
		"/users/:userID/revoke",
		server.IsAdminMiddleware,
		revokePermission,
	)
}

type revokePermissionBody struct {
	Permissions []string `json:"permissions" binding:"required"`
}

func revokePermission(c *gin.Context) {
	body := revokePermissionBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		server.AbortWithError(c, err)
		return
	}

	userId := c.Param("userID")

	di := server.GinDi(c)
	if err := di.UserService().Revoke(userId, body.Permissions); err != nil {
		server.AbortWithError(c, err)
		return
	}
}
