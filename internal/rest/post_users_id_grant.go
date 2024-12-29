package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/internal/rest/server"
	"github.com/nmarsollier/commongo/rst"
)

//	@Summary		Haiblitar permisos
//	@Description	Otorga permisos al usuario indicado, el usuario logueado tiene que tener permiso "admin".
//	@Tags			Seguridad
//	@Accept			json
//	@Produce		json
//	@Param			userId			path	string				true	"ID del usuario a habilitar permiso"
//	@Param			Authorization	header	string				true	"Bearer {token}"
//	@Param			body			body	grantPermissionBody	true	"Permisos a Habilitar"
//	@Success		200				"No Content"
//	@Failure		400				{object}	errs.ValidationErr	"Bad Request"
//	@Failure		401				{object}	rst.ErrorData		"Unauthorized"
//	@Failure		404				{object}	rst.ErrorData		"Not Found"
//	@Failure		500				{object}	rst.ErrorData		"Internal Server Error"
//	@Router			/users/:userID/grant [post]
//
// Otorga permisos al usuario indicado.
func postUsersIdGrantRoute(engine *gin.Engine) {
	engine.POST(
		"/users/:userID/grant",
		server.IsAdminMiddleware,
		grantPermission,
	)
}

type grantPermissionBody struct {
	Permissions []string `json:"permissions" binding:"required"`
}

func grantPermission(c *gin.Context) {
	body := grantPermissionBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		rst.AbortWithError(c, err)
		return
	}
	userId := c.Param("userID")

	di := server.GinDi(c)
	if err := di.UserService().Grant(userId, body.Permissions); err != nil {
		rst.AbortWithError(c, err)
		return
	}
}
