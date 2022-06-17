package routes

import (
	"github.com/gin-gonic/gin"
	rest "github.com/nmarsollier/authgo/rest"
	"github.com/nmarsollier/authgo/token"
)

/**
 * @api {get} /v1/user/signout Logout
 * @apiName Logout
 * @apiGroup Seguridad
 *
 * @apiDescription desloguea un usuario en el sistema, invalida el token.
 *
 * @apiSuccessExample {json} Respuesta
 *     HTTP/1.1 200 OK
 *
 * @apiUse AuthHeader
 * @apiUse OtherErrors
 */

func init() {
	rest.Router().GET(
		"/v1/user/signout",
		rest.ValidateLoggedIn,
		signOut,
	)
}

func signOut(c *gin.Context) {
	tokenString, _ := rest.RAWHeaderToken(c)

	if err := token.Invalidate(tokenString); err != nil {
		rest.AbortWithError(c, err)
		return
	}

	c.Done()
}
