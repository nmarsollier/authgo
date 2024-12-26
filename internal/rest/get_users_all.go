package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/internal/rest/server"
)

//	@Summary		Listar Usuarios
//	@Description	Obtiene información de todos los usuarios. El usuario logueado debe tener permisos "admin".
//	@Tags			Seguridad
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Bearer {token}"
//	@Success		200				{array}		user.UserData		"Users"
//	@Failure		400				{object}	errs.ValidationErr	"Bad Request"
//	@Failure		401				{object}	server.ErrorData	"Unauthorized"
//	@Failure		404				{object}	server.ErrorData	"Not Found"
//	@Failure		500				{object}	server.ErrorData	"Internal Server Error"
//	@Router			/users/all [get]
//
// Obtiene información de todos los usuarios.
func getUsersRoute(engine *gin.Engine) {
	engine.GET(
		"/users/all",
		server.IsAdminMiddleware,
		users,
	)
}

func users(c *gin.Context) {
	di := server.GinDi(c)
	result, err := di.UserService().FindAllUsers()

	if err != nil {
		server.AbortWithError(c, err)
		return
	}

	c.JSON(200, result)
}
