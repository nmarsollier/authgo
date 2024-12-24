package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/rest/server"
)

//	@Summary		Usuario Actual
//	@Description	Obtiene información del usuario actual.
//	@Tags			Seguridad
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Bearer {token}"
//	@Success		200				{object}	user.UserData		"User data"
//	@Failure		400				{object}	errs.ValidationErr	"Bad Request"
//	@Failure		401				{object}	server.ErrorData	"Unauthorized"
//	@Failure		404				{object}	server.ErrorData	"Not Found"
//	@Failure		500				{object}	server.ErrorData	"Internal Server Error"
//	@Router			/users/current [get]
//
// Obtiene información del usuario actual.
func getUsersCurrentRoute(engine *gin.Engine) {
	engine.GET(
		"/users/current",
		server.IsAuthenticatedMiddleware,
		currentUser,
	)
}

func currentUser(c *gin.Context) {
	token := server.HeaderToken(c)

	di := server.GinDi(c)
	user, err := di.UserService().FindById(token.UserID.Hex())
	if err != nil {
		server.AbortWithError(c, err)
		return
	}

	c.JSON(200, user)
}
