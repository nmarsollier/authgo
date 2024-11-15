package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/rest/server"
	"github.com/nmarsollier/authgo/user"
)

//	@Summary		Usuario Actual
//	@Description	Obtiene información del usuario actual.
//	@Tags			Seguridad
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Bearer {token}"
//	@Success		200				{object}	user.UserResponse	"User data"
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

	ctx := server.GinCtx(c)
	user, err := user.Get(token.UserID.Hex(), ctx...)
	if err != nil {
		server.AbortWithError(c, err)
		return
	}

	c.JSON(200, user)
}
