package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/model/token"
	"github.com/nmarsollier/authgo/rest/middlewares"
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
	router().GET(
		"/v1/user/signout",
		middlewares.ValidateLoggedIn,
		signOut,
	)
}

func signOut(c *gin.Context) {
	tokenString, _ := middlewares.RAWHeaderToken(c)

	if err := token.Invalidate(tokenString); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.Done()
}
