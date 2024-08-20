package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/rest/server"
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
//	@Failure		401		{object}	server.ErrorData	"Unauthorized"
//	@Failure		404		{object}	server.ErrorData	"Not Found"
//	@Failure		500		{object}	server.ErrorData	"Internal Server Error"
//	@Router			/v1/user [post]
//
// Registra un nuevo usuario en el sistema.
func postUsersRoute() {
	server.Router().POST(
		"/v1/user",
		signUp,
	)
}

func signUp(c *gin.Context) {
	body := user.SignUpRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		server.AbortWithError(c, err)
		return
	}

	ctx := server.GinCtx(c)
	token, err := user.SignUp(&body, ctx...)
	if err != nil {
		server.AbortWithError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"token": token,
	})
}
