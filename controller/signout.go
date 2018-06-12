package controller

import (
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/errors"

	"github.com/gin-gonic/gin"
)

// SignOut is the sign out controller
/**
 * @api {get} /auth/signout Logout
 * @apiName SignOut
 * @apiGroup Seguridad
 *
 * @apiDescription Desloguea un usuario en el sistema, invalida el token.
 *
 * @apiSuccessExample {json} Respuesta
 *     HTTP/1.1 200 OK
 *
 * @apiUse AuthHeader
 * @apiUse OtherErrors
 */
func SignOut(c *gin.Context) {
	err := token.Invalidate(c)

	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.Done()
}
