package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/rest/engine"
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

func getUserSignOutRoute() {
	engine.Router().GET(
		"/v1/user/signout",
		engine.ValidateLoggedIn,
		signOut,
	)
}

func signOut(c *gin.Context) {
	var extraParams []interface{}
	if mocks, ok := c.Get("mocks"); ok {
		extraParams = mocks.([]interface{})
	}

	tokenString, _ := engine.RAWHeaderToken(c)

	if err := token.Invalidate(tokenString, extraParams...); err != nil {
		engine.AbortWithError(c, err)
		return
	}

	c.Done()
}
