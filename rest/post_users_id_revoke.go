package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/rest/engine"
	"github.com/nmarsollier/authgo/user"
)

// Quita permisos al usuario indicado.
//
//	@Summary		Quitar permisos
//	@Description	Quita permisos al usuario indicado, el usuario logueado tiene que tener permiso "admin".
//	@Tags			Seguridad
//	@Accept			json
//	@Produce		json
//
//	@Param			userId			path	string				true	"ID del usuario a quitar permiso"
//	@Param			Authorization	header	string				true	"bearer {token}"
//	@Param			body			body	grantPermissionBody	true	"Permisos a Qutiar"
//	@Success		200				"No Content"
//	@Failure		400				{object}	apperr.ValidationErr	"Bad Request"
//	@Failure		401				{object}	engine.ErrorData		"Unauthorized"
//	@Failure		404				{object}	engine.ErrorData		"Not Found"
//	@Failure		500				{object}	engine.ErrorData		"Internal Server Error"
//	@Router			/v1/users/:userID/revoke [post]
func postUsersIdRevokeRoute() {
	engine.Router().POST(
		"/v1/users/:userID/revoke",
		engine.ValidateAdmin,
		revokePermission,
	)
}

type revokePermissionBody struct {
	Permissions []string `json:"permissions" binding:"required"`
}

func revokePermission(c *gin.Context) {
	body := revokePermissionBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		engine.AbortWithError(c, err)
		return
	}

	userId := c.Param("userID")

	ctx := engine.TestCtx(c)
	if err := user.Revoke(userId, body.Permissions, ctx...); err != nil {
		engine.AbortWithError(c, err)
		return
	}
}
