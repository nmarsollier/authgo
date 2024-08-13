package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/rest/engine"
	"github.com/nmarsollier/authgo/user"
)

// Loguea un usuario en el sistema.
//
//	@Summary		Login
//	@Description	Loguea un usuario en el sistema.
//	@Tags			Seguridad
//	@Accept			json
//	@Produce		json
//	@Param			body	body		user.SignInRequest		true	"Sign in information"
//	@Success		200		{object}	tokenResponse			"User Token"
//	@Failure		400		{object}	apperr.ValidationErr	"Bad Request"
//	@Failure		401		{object}	engine.ErrorData		"Unauthorized"
//	@Failure		404		{object}	engine.ErrorData		"Not Found"
//	@Failure		500		{object}	engine.ErrorData		"Internal Server Error"
//	@Router			/v1/user/signin [post]
//
// Handler function
func postUserSignInRoute() {
	engine.Router().POST(
		"/v1/user/signin",
		signIn,
	)
}

type tokenResponse struct {
	Token string `json:"token"`
}

func signIn(c *gin.Context) {
	login := user.SignInRequest{}
	if err := c.ShouldBindJSON(&login); err != nil {
		engine.AbortWithError(c, err)
		return
	}

	ctx := engine.TestCtx(c)
	tokenString, err := user.SignIn(login, ctx...)
	if err != nil {
		engine.AbortWithError(c, err)
		return
	}

	c.JSON(200, tokenResponse{
		Token: tokenString,
	})
}
