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
//
//	@Param			body	body		user.SignInRequest			true	"Sign in information"
//
//	@Success		200		{object}	tokenResponse				"User Token"
//
//	@Failure		400		{object}	app_errors.ErrValidation	"Bad Request"
//	@Failure		401		{object}	app_errors.OtherErrors		"Unauthorized"
//	@Failure		404		{object}	app_errors.OtherErrors		"Not Found"
//	@Failure		500		{object}	app_errors.OtherErrors		"Internal Server Error"
//
//	@Router			/v1/user/signin [post]
func postUserSignInRoute() {
	engine.Router().POST(
		"/v1/user/signin",
		validateSignInBody,
		signIn,
	)
}

type tokenResponse struct {
	Token string `json:"token"`
}

func signIn(c *gin.Context) {
	login := c.MustGet("data").(user.SignInRequest)

	var extraParams []interface{}
	if mocks, ok := c.Get("mocks"); ok {
		extraParams = mocks.([]interface{})
	}

	tokenString, err := user.SignIn(login, extraParams...)
	if err != nil {
		engine.AbortWithError(c, err)
		return
	}

	c.JSON(200, tokenResponse{
		Token: tokenString,
	})
}

func validateSignInBody(c *gin.Context) {
	login := user.SignInRequest{}
	if err := c.ShouldBindJSON(&login); err != nil {
		engine.AbortWithError(c, err)
		return
	}

	c.Set("data", login)
	c.Next()
}
