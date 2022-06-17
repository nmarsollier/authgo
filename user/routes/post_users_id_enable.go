package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/rest"
	"github.com/nmarsollier/authgo/user"
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
	rest.Router().POST(
		"/v1/users/:userID/enable",
		rest.ValidateAdmin,
		enable,
	)
}

func enable(c *gin.Context) {
	if err := user.Enable(c.Param("userID")); err != nil {
		rest.AbortWithError(c, err)
	}
}
