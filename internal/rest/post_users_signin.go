package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/internal/rest/server"
	"github.com/nmarsollier/authgo/internal/usecases"
	"github.com/nmarsollier/commongo/rst"
)

//	@Summary		Login
//	@Description	Loguea un usuario en el sistema.
//	@Tags			Seguridad
//	@Accept			json
//	@Produce		json
//	@Param			body	body		usecases.SignInRequest	true	"Sign in information"
//	@Success		200		{object}	usecases.TokenResponse	"User Token"
//	@Failure		400		{object}	errs.ValidationErr		"Bad Request"
//	@Failure		401		{object}	rst.ErrorData			"Unauthorized"
//	@Failure		404		{object}	rst.ErrorData			"Not Found"
//	@Failure		500		{object}	rst.ErrorData			"Internal Server Error"
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
		rst.AbortWithError(c, err)
		return
	}

	di := server.GinDi(c)
	token, err := di.SignInUseCase().SignIn(login)
	if err != nil {
		rst.AbortWithError(c, err)
		return
	}

	c.JSON(200, token)
}
