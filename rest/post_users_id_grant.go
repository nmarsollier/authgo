package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/rest/engine"
	"github.com/nmarsollier/authgo/user"
)

// Otorga permisos al usuario indicado.
//
//	@Summary		Haiblitar permisos
//	@Description	Otorga permisos al usuario indicado, el usuario logueado tiene que tener permiso "admin".
//	@Tags			Seguridad
//	@Accept			json
//	@Produce		json
//
//	@Param			userId			path	string				true	"ID del usuario a habilitar permiso"
//	@Param			Authorization	header	string				true	"bearer {token}"
//	@Param			body			body	grantPermissionBody	true	"Permisos a Habilitar"
//	@Success		200				"No Content"
//	@Failure		400				{object}	apperr.ValidationErr	"Bad Request"
//	@Failure		401				{object}	engine.ErrorData		"Unauthorized"
//	@Failure		404				{object}	engine.ErrorData		"Not Found"
//	@Failure		500				{object}	engine.ErrorData		"Internal Server Error"
//	@Router			/v1/users/:userID/grant [post]
func postUsersIdGrantRoute() {
	engine.Router().POST(
		"/v1/users/:userID/grant",
		engine.ValidateAdmin,
		grantPermission,
	)
}

type grantPermissionBody struct {
	Permissions []string `json:"permissions" binding:"required"`
}

func grantPermission(c *gin.Context) {
	body := grantPermissionBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		engine.AbortWithError(c, err)
		return
	}
	userId := c.Param("userID")

	ctx := engine.TestCtx(c)
	if err := user.Grant(userId, body.Permissions, ctx...); err != nil {
		engine.AbortWithError(c, err)
		return
	}
}
