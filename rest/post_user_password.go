package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/rest/engine"
	"github.com/nmarsollier/authgo/user"
)

// @Summary		Cambiar Password
// @Description	Cambia la contraseña del usuario actual.
// @Tags			Seguridad
// @Accept			json
// @Produce		json
// @Param			body			body	changePasswordBody	true	"Passwords"
// @Param			Authorization	header	string				true	"bearer {token}"
// @Success		200				"No Content"
// @Failure		400				{object}	apperr.ValidationErr	"Bad Request"
// @Failure		401				{object}	engine.ErrorData		"Unauthorized"
// @Failure		404				{object}	engine.ErrorData		"Not Found"
// @Failure		500				{object}	engine.ErrorData		"Internal Server Error"
// @Router			/v1/user/password [post]
//
// Cambia la contraseña del usuario actual.
func getUserPasswordRoute() {
	engine.Router().POST(
		"/v1/user/password",
		engine.ValidateLoggedIn,
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
		engine.AbortWithError(c, err)
		return
	}
	token := engine.HeaderToken(c)

	ctx := engine.TestCtx(c)
	if err := user.ChangePassword(token.UserID.Hex(), body.Current, body.New, ctx...); err != nil {
		engine.AbortWithError(c, err)
		return
	}
}
