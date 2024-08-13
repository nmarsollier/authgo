package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/rest/engine"
	"github.com/nmarsollier/authgo/token"
)

// Desloguea un usuario en el sistema, invalida el token.
//
//	@Summary		Logout
//	@Description	Desloguea un usuario en el sistema, invalida el token.
//	@Tags			Seguridad
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header	string	true	"bearer {token}"
//
//	@Success		200				"No Content"
//
//	@Failure		500				{object}	engine.ErrorData	"Error response"
//	@Router			/v1/user/signout [get]
func getUserSignOutRoute() {
	engine.Router().GET(
		"/v1/user/signout",
		engine.ValidateLoggedIn,
		signOut,
	)
}

func signOut(c *gin.Context) {
	tokenString, _ := engine.HeaderAuthorization(c)

	ctx := engine.TestCtx(c)
	if err := token.Invalidate(tokenString, ctx...); err != nil {
		engine.AbortWithError(c, err)
		return
	}

	c.Done()
}
