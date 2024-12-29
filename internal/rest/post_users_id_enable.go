package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/internal/rest/server"
	"github.com/nmarsollier/commongo/rst"
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
//	@Failure		401				{object}	rst.ErrorData		"Unauthorized"
//	@Failure		404				{object}	rst.ErrorData		"Not Found"
//	@Failure		500				{object}	rst.ErrorData		"Internal Server Error"
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
	di := server.GinDi(c)

	if err := di.UserService().Enable(userId); err != nil {
		rst.AbortWithError(c, err)
	}
}
