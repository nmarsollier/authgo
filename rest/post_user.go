package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/rest/engine"
	"github.com/nmarsollier/authgo/user"
)

// Registra un nuevo usuario en el sistema.
//
//	@Summary		Registrar Usuario
//	@Description	Registra un nuevo usuario en el sistema.
//	@Tags			Seguridad
//	@Accept			json
//	@Produce		json
//
//	@Param			body	body		user.SignUpRequest			true	"Informacion de ususario"
//
//	@Success		200		{object}	tokenResponse				"User Token"
//
//	@Failure		400		{object}	app_errors.ErrValidation	"Bad Request"
//	@Failure		401		{object}	app_errors.OtherErrors		"Unauthorized"
//	@Failure		404		{object}	app_errors.OtherErrors		"Not Found"
//	@Failure		500		{object}	app_errors.OtherErrors		"Internal Server Error"
//
//	@Router			/v1/user [post]
func postUsersRoute() {
	engine.Router().POST(
		"/v1/user",
		validateSignUpBody,
		signUp,
	)
}

func signUp(c *gin.Context) {
	body := c.MustGet("data").(user.SignUpRequest)

	var extraParams []interface{}
	if mocks, ok := c.Get("mocks"); ok {
		extraParams = mocks.([]interface{})
	}

	token, err := user.SignUp(&body, extraParams...)
	if err != nil {
		engine.AbortWithError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"token": token,
	})
}

func validateSignUpBody(c *gin.Context) {
	body := user.SignUpRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		engine.AbortWithError(c, err)
		return
	}

	c.Set("data", body)
	c.Next()
}
