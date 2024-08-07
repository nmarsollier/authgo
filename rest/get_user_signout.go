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
//	@Failure		500				{object}	app_errors.OtherErrors	"Error response"
//	@Router			/v1/user/signout [get]
func getUserSignOutRoute() {
	engine.Router().GET(
		"/v1/user/signout",
		engine.ValidateLoggedIn,
		signOut,
	)
}

func signOut(c *gin.Context) {
	var extraParams []interface{}
	if mocks, ok := c.Get("mocks"); ok {
		extraParams = mocks.([]interface{})
	}

	tokenString, _ := engine.RAWHeaderToken(c)

	if err := token.Invalidate(tokenString, extraParams...); err != nil {
		engine.AbortWithError(c, err)
		return
	}

	c.Done()
}
