package controller

import (
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/errors"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/user"
)

/**
 * @api {get} /auth/currentUser Usuario Actual
 * @apiName CurrentUser
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
 *        "roles": [
 *            "{Rol}"
 *        ]
 *     }
 *
 * @apiUse AuthHeader
 * @apiUse OtherErrors
 */
// CurrentUser is the controller to get the current logged in user
func CurrentUser(c *gin.Context) {
	payload, err := token.ValidateToken(c)

	if err != nil {
		errors.HandleError(c, err)
		return
	}

	user, err := user.CurrentUser(payload.UserID)

	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"id":    user.ID(),
		"name":  user.Name,
		"roles": user.Roles,
		"login": user.Login,
	})
}
