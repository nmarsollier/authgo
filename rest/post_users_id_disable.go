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
//
//	@Param			userId			path	string	true	"ID del usuario a deshabilitar"
//	@Param			Authorization	header	string	true	"bearer {token}"
//	@Success		200				"No Content"
//
//	@Failure		400				{object}	app_errors.ErrValidation	"Bad Request"
//	@Failure		401				{object}	app_errors.OtherErrors		"Unauthorized"
//	@Failure		404				{object}	app_errors.OtherErrors		"Not Found"
//	@Failure		500				{object}	app_errors.OtherErrors		"Internal Server Error"
//
//	@Router			/v1/users/:userId/disable [post]
func postUsersIdDisableRoute() {
	engine.Router().POST(
		"/v1/users/:userID/disable",
		engine.ValidateAdmin,
		disable,
	)
}

func disable(c *gin.Context) {
	var extraParams []interface{}
	if mocks, ok := c.Get("mocks"); ok {
		extraParams = mocks.([]interface{})
	}

	if err := user.Disable(c.Param("userID"), extraParams...); err != nil {
		engine.AbortWithError(c, err)
	}
}
