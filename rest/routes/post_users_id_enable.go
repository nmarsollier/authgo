package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/model/user"
	"github.com/nmarsollier/authgo/rest/middlewares"
)

/**
 * @api {post} /v1/users/:userId/enable Habilitar Usuario
 * @apiName Habilitar Usuario
 * @apiGroup Seguridad
 *
 * @apiDescription Habilita un usuario en el sistema. El usuario logueado debe tener permisos "admin".
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
		"/v1/users/:userID/enable",
		middlewares.ValidateAdmin,
		enable,
	)
}

func enable(c *gin.Context) {
	if err := user.Enable(c.Param("userID")); err != nil {
		middlewares.AbortWithError(c, err)
	}
}
