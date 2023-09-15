package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/rest/engine"
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

func getUsersRoute() {
	engine.Router().GET(
		"/v1/users",
		engine.ValidateAdmin,
		users,
	)
}

func users(c *gin.Context) {
	var extraParams []interface{}
	if mocks, ok := c.Get("mocks"); ok {
		extraParams = mocks.([]interface{})
	}

	user, err := user.Users(extraParams...)

	if err != nil {
		engine.AbortWithError(c, err)
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
