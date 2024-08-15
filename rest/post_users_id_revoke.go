package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/rest/server"
	"github.com/nmarsollier/authgo/user"
)

//	@Summary		Quitar permisos
//	@Description	Quita permisos al usuario indicado, el usuario logueado tiene que tener permiso "admin".
//	@Tags			Seguridad
//	@Accept			json
//	@Produce		json
//	@Param			userId			path	string				true	"ID del usuario a quitar permiso"
//	@Param			Authorization	header	string				true	"bearer {token}"
//	@Param			body			body	grantPermissionBody	true	"Permisos a Qutiar"
//	@Success		200				"No Content"
//	@Failure		400				{object}	errs.ValidationErr	"Bad Request"
//	@Failure		401				{object}	server.ErrorData	"Unauthorized"
//	@Failure		404				{object}	server.ErrorData	"Not Found"
//	@Failure		500				{object}	server.ErrorData	"Internal Server Error"
//	@Router			/v1/users/:userID/revoke [post]
//
// Quita permisos al usuario indicado.
func postUsersIdRevokeRoute() {
	server.Router().POST(
		"/v1/users/:userID/revoke",
		server.ValidateAdmin,
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

	ctx := server.TestCtx(c)
	if err := user.Revoke(userId, body.Permissions, ctx...); err != nil {
		server.AbortWithError(c, err)
		return
	}
}
