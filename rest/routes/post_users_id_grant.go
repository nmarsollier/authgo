package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/model/user"
	"github.com/nmarsollier/authgo/rest/middlewares"
)

/**
 * @api {post} /v1/users/:userId/grant Otorga Permisos
 * @apiName Otorga Permisos
 * @apiGroup Seguridad
 *
 * @apiDescription Otorga permisos al usuario indicado, el usuario logueado tiene que tener permiso "admin".
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
		"/v1/users/:userID/grant",
		middlewares.ValidateAdmin,
		validateGrantBody,
		grantPermission,
	)
}

type grantPermissionBody struct {
	Permissions []string `json:"permissions" binding:"required"`
}

func grantPermission(c *gin.Context) {
	body := c.MustGet("data").(grantPermissionBody)

	if err := user.Grant(c.Param("userID"), body.Permissions); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}
}

func validateGrantBody(c *gin.Context) {
	body := grantPermissionBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.Set("data", body)
	c.Next()
}
