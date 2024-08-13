package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/rest/engine"
	"github.com/nmarsollier/authgo/user"
)

// Deshabilitar Usuario
//
//	@Summary		Deshabilitar Usuario
//	@Description	Deshabilita un usuario en el sistema. El usuario logueado debe tener permisos "admin".
//	@Tags			Seguridad
//	@Accept			json
//	@Produce		json
//	@Param			userId			path	string	true	"ID del usuario a deshabilitar"
//	@Param			Authorization	header	string	true	"bearer {token}"
//	@Success		200				"No Content"
//	@Failure		400				{object}	apperr.ValidationErr	"Bad Request"
//	@Failure		401				{object}	engine.ErrorData		"Unauthorized"
//	@Failure		404				{object}	engine.ErrorData		"Not Found"
//	@Failure		500				{object}	engine.ErrorData		"Internal Server Error"
//	@Router			/v1/users/:userId/disable [post]
//
// Handler function
func postUsersIdDisableRoute() {
	engine.Router().POST(
		"/v1/users/:userID/disable",
		engine.ValidateAdmin,
		disable,
	)
}

func disable(c *gin.Context) {
	userId := c.Param("userID")

	ctx := engine.TestCtx(c)
	if err := user.Disable(userId, ctx...); err != nil {
		engine.AbortWithError(c, err)
	}
}
