package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/rest/server"
	"github.com/nmarsollier/authgo/user"
)

type UserDataResponse struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
	Login       string   `json:"login"`
	Enabled     bool     `json:"enabled"`
}

//	@Summary		Listar Usuarios
//	@Description	Obtiene información de todos los usuarios. El usuario logueado debe tener permisos "admin".
//	@Tags			Seguridad
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"bearer {token}"
//	@Success		200				{array}		UserDataResponse	"Users"
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
	user, err := user.Users(ctx...)

	if err != nil {
		server.AbortWithError(c, err)
		return
	}
	result := []UserDataResponse{}
	for i := 0; i < len(user); i = i + 1 {
		result = append(result, UserDataResponse{
			Id:          user[i].ID.Hex(),
			Name:        user[i].Name,
			Permissions: user[i].Permissions,
			Login:       user[i].Login,
			Enabled:     user[i].Enabled,
		})
	}

	c.JSON(200, result)
}
