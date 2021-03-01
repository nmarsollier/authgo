package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/model/user"
	"github.com/nmarsollier/authgo/rest/middlewares"
)

/**
 * @api {get} /v1/users/current Usuario Actual
 * @apiName Usuario Actual
 * @apiGroup Seguridad
 *
 * @apiDescription Obtiene información del usuario actual.
 *
 * @apiSuccessExample {json} Respuesta
 *     HTTP/1.1 200 OK
 *     {
 *        "id": "{Id usuario}",
 *        "name": "{Nombre del usuario}",
 *        "login": "{Login de usuario}",
 *        "permissions": [
 *            "{Permission}"
 *        ]
 *     }
 *
 * @apiUse AuthHeader
 * @apiUse OtherErrors
 */

func init() {
	router().GET(
		"/v1/users/current",
		middlewares.ValidateLoggedIn,
		currentUser,
	)
}

func currentUser(c *gin.Context) {
	token := middlewares.HeaderToken(c)

	user, err := user.Get(token.UserID.Hex())
	if err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"id":          user.ID.Hex(),
		"name":        user.Name,
		"permissions": user.Permissions,
		"login":       user.Login,
	})
}
