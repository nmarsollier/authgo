package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/internal/rest/server"
	"github.com/nmarsollier/authgo/internal/usecases"
)

// @Summary		Logout
// @Description	Desloguea un usuario en el sistema, invalida el token.
// @Tags			Seguridad
// @Accept			json
// @Produce		json
// @Param			Authorization	header	string	true	"Bearer {token}"
// @Success		200				"No Content"
// @Failure		500				{object}	server.ErrorData	"Error response"
// @Router			/users/signout [get]
//
// Desloguea un usuario en el sistema, invalida el token.
func getUserSignOutRoute(engine *gin.Engine) {
	engine.GET(
		"/users/signout",
		server.IsAuthenticatedMiddleware,
		signOut,
	)
}

func signOut(c *gin.Context) {
	tokenString, _ := server.GetHeaderToken(c)
	log := server.GinLogger(c)

	if err := usecases.InvalidateToken(log, tokenString); err != nil {
		server.AbortWithError(c, err)
		return
	}

	c.Done()
}
