package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/rest/server"
	"github.com/nmarsollier/authgo/token"
)

// @Summary		Logout
// @Description	Desloguea un usuario en el sistema, invalida el token.
// @Tags			Seguridad
// @Accept			json
// @Produce		json
// @Param			Authorization	header	string	true	"bearer {token}"
// @Success		200				"No Content"
// @Failure		500				{object}	server.ErrorData	"Error response"
// @Router			/v1/user/signout [get]
//
// Desloguea un usuario en el sistema, invalida el token.
func getUserSignOutRoute() {
	server.Router().GET(
		"/v1/user/signout",
		server.ValidateLoggedIn,
		signOut,
	)
}

func signOut(c *gin.Context) {
	tokenString, _ := server.HeaderAuthorization(c)

	ctx := server.TestCtx(c)
	if err := token.Invalidate(tokenString, ctx...); err != nil {
		server.AbortWithError(c, err)
		return
	}

	c.Done()
}
