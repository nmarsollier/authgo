package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/internal/rest/server"
	"github.com/nmarsollier/authgo/internal/usecases"
	"github.com/nmarsollier/commongo/rst"
)

//	@Summary		Registrar Usuario
//	@Description	Registra un nuevo usuario en el sistema.
//	@Tags			Seguridad
//	@Accept			json
//	@Produce		json
//	@Param			body	body		usecases.SignUpRequest	true	"Informacion de ususario"
//	@Success		200		{object}	usecases.TokenResponse	"User Token"
//	@Failure		400		{object}	errs.ValidationErr		"Bad Request"
//	@Failure		401		{object}	rst.ErrorData			"Unauthorized"
//	@Failure		404		{object}	rst.ErrorData			"Not Found"
//	@Failure		500		{object}	rst.ErrorData			"Internal Server Error"
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
		rst.AbortWithError(c, err)
		return
	}

	di := server.GinDi(c)
	token, err := di.SignUpUseCase().SignUp(&body)
	if err != nil {
		rst.AbortWithError(c, err)
		return
	}

	c.JSON(200, token)
}
