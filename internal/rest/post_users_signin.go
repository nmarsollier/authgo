package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/internal/rest/server"
	"github.com/nmarsollier/authgo/internal/usecases"
)

//	@Summary		Login
//	@Description	Loguea un usuario en el sistema.
//	@Tags			Seguridad
//	@Accept			json
//	@Produce		json
//	@Param			body	body		usecases.SignInRequest	true	"Sign in information"
//	@Success		200		{object}	usecases.TokenResponse	"User Token"
//	@Failure		400		{object}	errs.ValidationErr		"Bad Request"
//	@Failure		401		{object}	server.ErrorData		"Unauthorized"
//	@Failure		404		{object}	server.ErrorData		"Not Found"
//	@Failure		500		{object}	server.ErrorData		"Internal Server Error"
//	@Router			/users/signin [post]
//
// Loguea un usuario en el sistema.
func postUserSignInRoute(engine *gin.Engine) {
	engine.POST(
		"/users/signin",
		signIn,
	)
}

func signIn(c *gin.Context) {
	login := &usecases.SignInRequest{}
	if err := c.ShouldBindJSON(&login); err != nil {
		server.AbortWithError(c, err)
		return
	}

	log := server.GinLogger(c)
	token, err := usecases.SignIn(log, login)
	if err != nil {
		server.AbortWithError(c, err)
		return
	}

	c.JSON(200, token)
}
