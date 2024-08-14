package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/rest/engine"
	"github.com/nmarsollier/authgo/user"
)

//	@Summary		Registrar Usuario
//	@Description	Registra un nuevo usuario en el sistema.
//	@Tags			Seguridad
//	@Accept			json
//	@Produce		json
//	@Param			body	body		user.SignUpRequest	true	"Informacion de ususario"
//	@Success		200		{object}	tokenResponse		"User Token"
//	@Failure		400		{object}	errs.ValidationErr	"Bad Request"
//	@Failure		401		{object}	engine.ErrorData	"Unauthorized"
//	@Failure		404		{object}	engine.ErrorData	"Not Found"
//	@Failure		500		{object}	engine.ErrorData	"Internal Server Error"
//	@Router			/v1/user [post]
//
// Registra un nuevo usuario en el sistema.
func postUsersRoute() {
	engine.Router().POST(
		"/v1/user",
		signUp,
	)
}

func signUp(c *gin.Context) {
	body := user.SignUpRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		engine.AbortWithError(c, err)
		return
	}

	ctx := engine.TestCtx(c)
	token, err := user.SignUp(&body, ctx...)
	if err != nil {
		engine.AbortWithError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"token": token,
	})
}
