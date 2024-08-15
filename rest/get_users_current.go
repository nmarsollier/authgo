package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/rest/server"
	"github.com/nmarsollier/authgo/user"
)

type UserResponse struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
	Login       string   `json:"login"`
}

//	@Summary		Usuario Actual
//	@Description	Obtiene información del usuario actual.
//	@Tags			Seguridad
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"bearer {token}"
//	@Success		200				{object}	UserResponse		"User data"
//	@Failure		400				{object}	errs.ValidationErr	"Bad Request"
//	@Failure		401				{object}	server.ErrorData	"Unauthorized"
//	@Failure		404				{object}	server.ErrorData	"Not Found"
//	@Failure		500				{object}	server.ErrorData	"Internal Server Error"
//	@Router			/v1/users/current [get]
//
// Obtiene información del usuario actual.
func getUsersCurrentRoute() {
	server.Router().GET(
		"/v1/users/current",
		server.ValidateLoggedIn,
		currentUser,
	)
}

func currentUser(c *gin.Context) {
	token := server.HeaderToken(c)

	ctx := server.TestCtx(c)
	user, err := user.Get(token.UserID.Hex(), ctx...)
	if err != nil {
		server.AbortWithError(c, err)
		return
	}

	c.JSON(200, UserResponse{
		Id:          user.ID.Hex(),
		Name:        user.Name,
		Permissions: user.Permissions,
		Login:       user.Login,
	})
}
