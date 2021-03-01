package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/model/user"
	"github.com/nmarsollier/authgo/rest/middlewares"
)

/**
 * @api {post} /v1/users/:userId/revoke Revoca Permisos
 * @apiName Revoca Permisos
 * @apiGroup Seguridad
 *
 * @apiDescription Quita permisos al usuario indicado, el usuario logueado tiene que tener permiso "admin".
 *
 * @apiExample {json} Body
 *    {
 *      "permissions" : ["permiso", ...],
 *    }
 *
 * @apiSuccessExample {json} Respuesta
 *     HTTP/1.1 200 OK
 *
 * @apiUse AuthHeader
 * @apiUse ParamValidationErrors
 * @apiUse OtherErrors
 */

func init() {
	router().POST(
		"/v1/users/:userID/revoke",
		middlewares.ValidateAdmin,
		validateRevokeBody,
		revokePermission,
	)
}

type revokePermissionBody struct {
	Permissions []string `json:"permissions" binding:"required"`
}

func revokePermission(c *gin.Context) {
	body := c.MustGet("data").(revokePermissionBody)

	if err := user.Revoke(c.Param("userID"), body.Permissions); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}
}

func validateRevokeBody(c *gin.Context) {
	body := revokePermissionBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.Set("data", body)
	c.Next()
}
