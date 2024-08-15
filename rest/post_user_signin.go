package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/rest/server"
	"github.com/nmarsollier/authgo/user"
)

//	@Summary		Login
//	@Description	Loguea un usuario en el sistema.
//	@Tags			Seguridad
//	@Accept			json
//	@Produce		json
//	@Param			body	body		user.SignInRequest	true	"Sign in information"
//	@Success		200		{object}	tokenResponse		"User Token"
//	@Failure		400		{object}	errs.ValidationErr	"Bad Request"
//	@Failure		401		{object}	server.ErrorData	"Unauthorized"
//	@Failure		404		{object}	server.ErrorData	"Not Found"
//	@Failure		500		{object}	server.ErrorData	"Internal Server Error"
//	@Router			/v1/user/signin [post]
//
// Loguea un usuario en el sistema.
func postUserSignInRoute() {
	server.Router().POST(
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
		server.AbortWithError(c, err)
		return
	}

	ctx := server.TestCtx(c)
	tokenString, err := user.SignIn(login, ctx...)
	if err != nil {
		server.AbortWithError(c, err)
		return
	}

	c.JSON(200, tokenResponse{
		Token: tokenString,
	})
}
