package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/internal/rest/server"
	"github.com/nmarsollier/authgo/internal/usecases"
)

//	@Summary		Registrar Usuario
//	@Description	Registra un nuevo usuario en el sistema.
//	@Tags			Seguridad
//	@Accept			json
//	@Produce		json
//	@Param			body	body		usecases.SignUpRequest	true	"Informacion de ususario"
//	@Success		200		{object}	usecases.TokenResponse	"User Token"
//	@Failure		400		{object}	errs.ValidationErr		"Bad Request"
//	@Failure		401		{object}	server.ErrorData		"Unauthorized"
//	@Failure		404		{object}	server.ErrorData		"Not Found"
//	@Failure		500		{object}	server.ErrorData		"Internal Server Error"
//	@Router			/users/signup [post]
//
// Registra un nuevo usuario en el sistema.
func postUsersRoute(engine *gin.Engine) {
	engine.POST(
		"/users/signup",
		signUp,
	)
}

func signUp(c *gin.Context) {
	body := usecases.SignUpRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		server.AbortWithError(c, err)
		return
	}

	log := server.GinLogger(c)
	token, err := usecases.SignUp(log, &body)
	if err != nil {
		server.AbortWithError(c, err)
		return
	}

	c.JSON(200, token)
}
