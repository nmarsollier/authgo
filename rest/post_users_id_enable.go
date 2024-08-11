package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/rest/engine"
	"github.com/nmarsollier/authgo/user"
)

// Habilita un usuario en el sistema.
//
//	@Summary		Enable User
//	@Description	Habilita un usuario en el sistema. El usuario logueado debe tener permisos "admin".
//	@Tags			Seguridad
//	@Accept			json
//	@Produce		json
//
//	@Param			userId			path	string	true	"ID del usuario a habilitar"
//	@Param			Authorization	header	string	true	"bearer {token}"
//	@Success		200				"No Content"
//
//	@Failure		400				{object}	app_errors.ErrValidation	"Bad Request"
//	@Failure		401				{object}	app_errors.OtherErrors		"Unauthorized"
//	@Failure		404				{object}	app_errors.OtherErrors		"Not Found"
//	@Failure		500				{object}	app_errors.OtherErrors		"Internal Server Error"
//
//	@Router			/v1/users/:userId/enable [post]
func postUsersIdEnableRoute() {
	engine.Router().POST(
		"/v1/users/:userID/enable",
		engine.ValidateAdmin,
		enable,
	)
}

func enable(c *gin.Context) {
	var extraParams []interface{}
	if mocks, ok := c.Get("mocks"); ok {
		extraParams = mocks.([]interface{})
	}

	if err := user.Enable(c.Param("userID"), extraParams...); err != nil {
		engine.AbortWithError(c, err)
	}
}
