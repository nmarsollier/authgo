package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/rest"
	"github.com/nmarsollier/authgo/user"
)

/**
 * @api {post} /v1/users/:userId/disable Deshabilitar Usuario
 * @apiName Deshabilitar Usuario
 * @apiGroup Seguridad
 *
 * @apiDescription Deshabilita un usuario en el sistema.   El usuario logueado debe tener permisos "admin".
 *
 * @apiSuccessExample {json} Respuesta
 *     HTTP/1.1 200 OK
 *
 * @apiUse AuthHeader
 * @apiUse ParamValidationErrors
 * @apiUse OtherErrors
 */

func init() {
	rest.Router().POST(
		"/v1/users/:userID/disable",
		rest.ValidateAdmin,
		disable,
	)
}

func disable(c *gin.Context) {
	if err := user.Disable(c.Param("userID")); err != nil {
		rest.AbortWithError(c, err)
	}
}
