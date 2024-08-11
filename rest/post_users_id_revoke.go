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
//
//	@Failure		400				{object}	app_errors.ErrValidation	"Bad Request"
//	@Failure		401				{object}	app_errors.OtherErrors		"Unauthorized"
//	@Failure		404				{object}	app_errors.OtherErrors		"Not Found"
//	@Failure		500				{object}	app_errors.OtherErrors		"Internal Server Error"
//
//	@Router			/v1/users/:userID/revoke [post]
func postUsersIdRevokeRoute() {
	engine.Router().POST(
		"/v1/users/:userID/revoke",
		engine.ValidateAdmin,
		validateRevokeBody,
		revokePermission,
	)
}

type revokePermissionBody struct {
	UserId      string   `json:"user" binding:"required"`
	Permissions []string `json:"permissions" binding:"required"`
}

func revokePermission(c *gin.Context) {
	body := c.MustGet("data").(revokePermissionBody)

	var extraParams []interface{}
	if mocks, ok := c.Get("mocks"); ok {
		extraParams = mocks.([]interface{})
	}

	if err := user.Revoke(body.UserId, body.Permissions, extraParams...); err != nil {
		engine.AbortWithError(c, err)
		return
	}
}

func validateRevokeBody(c *gin.Context) {
	body := revokePermissionBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		engine.AbortWithError(c, err)
		return
	}

	c.Set("data", body)
	c.Next()
}
