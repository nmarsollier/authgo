package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/internal/rest/server"
	"github.com/nmarsollier/commongo/rst"
)

//	@Summary		Cambiar Password
//	@Description	Cambia la contraseña del usuario actual.
//	@Tags			Seguridad
//	@Accept			json
//	@Produce		json
//	@Param			body			body	changePasswordBody	true	"Passwords"
//	@Param			Authorization	header	string				true	"Bearer {token}"
//	@Success		200				"No Content"
//	@Failure		400				{object}	errs.ValidationErr	"Bad Request"
//	@Failure		401				{object}	rst.ErrorData		"Unauthorized"
//	@Failure		404				{object}	rst.ErrorData		"Not Found"
//	@Failure		500				{object}	rst.ErrorData		"Internal Server Error"
//	@Router			/users/password [post]
//
// Cambia la contraseña del usuario actual.
func getUserPasswordRoute(engine *gin.Engine) {
	engine.POST(
		"/users/password",
		server.IsAuthenticatedMiddleware,
		changePassword,
	)
}

type changePasswordBody struct {
	Current string `json:"currentPassword" binding:"required,min=1,max=100"`
	New     string `json:"newPassword" binding:"required,min=1,max=100"`
}

func changePassword(c *gin.Context) {
	body := changePasswordBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		rst.AbortWithError(c, err)
		return
	}
	token := server.GetCtxToken(c)

	di := server.GinDi(c)
	if err := di.UserService().ChangePassword(token.UserID.Hex(), body.Current, body.New); err != nil {
		rst.AbortWithError(c, err)
		return
	}
}
