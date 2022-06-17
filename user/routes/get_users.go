package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/rest"
	"github.com/nmarsollier/authgo/user"
)

/**
 * @api {get} /v1/users Listar Usuarios
 * @apiName Listar Usuarios
 * @apiGroup Seguridad
 *
 * @apiDescription Obtiene información de todos los usuarios.
 *
 * @apiSuccessExample {json} Respuesta
 *     HTTP/1.1 200 OK
 *     [{
 *        "id": "{Id usuario}",
 *        "name": "{Nombre del usuario}",
 *        "login": "{Login de usuario}",
 *        "permissions": [
 *            "{Permission}"
 *        ],
 * 	      "enabled": true|false
 *     }, ...]
 *
 * @apiUse AuthHeader
 * @apiUse OtherErrors
 */

func init() {
	rest.Router().GET(
		"/v1/users",
		rest.ValidateAdmin,
		users,
	)
}

func users(c *gin.Context) {
	user, err := user.Users()

	if err != nil {
		rest.AbortWithError(c, err)
		return
	}
	result := []gin.H{}
	for i := 0; i < len(user); i = i + 1 {
		result = append(result, gin.H{
			"id":          user[i].ID.Hex(),
			"name":        user[i].Name,
			"permissions": user[i].Permissions,
			"login":       user[i].Login,
			"enabled":     user[i].Enabled,
		})
	}

	c.JSON(200, result)
}
