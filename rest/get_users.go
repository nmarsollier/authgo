package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/rest/server"
	"github.com/nmarsollier/authgo/user"
)

//	@Summary		Listar Usuarios
//	@Description	Obtiene información de todos los usuarios. El usuario logueado debe tener permisos "admin".
//	@Tags			Seguridad
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Bearer {token}"
//	@Success		200				{array}		user.UserResponse	"Users"
//	@Failure		400				{object}	errs.ValidationErr	"Bad Request"
//	@Failure		401				{object}	server.ErrorData	"Unauthorized"
//	@Failure		404				{object}	server.ErrorData	"Not Found"
//	@Failure		500				{object}	server.ErrorData	"Internal Server Error"
//	@Router			/v1/users [get]
//
// Obtiene información de todos los usuarios.
func getUsersRoute() {
	server.Router().GET(
		"/v1/users",
		server.ValidateAdmin,
		users,
	)
}

func users(c *gin.Context) {
	ctx := server.GinCtx(c)
	result, err := user.FindAllUsers(ctx...)

	if err != nil {
		server.AbortWithError(c, err)
		return
	}

	c.JSON(200, result)
}
