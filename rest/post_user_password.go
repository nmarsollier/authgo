package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/rest/engine"
	"github.com/nmarsollier/authgo/user"
)

// Cambia la contraseña del usuario actual.
//
//	@Summary		Cambiar Password
//	@Description	Cambia la contraseña del usuario actual.
//	@Tags			Seguridad
//	@Accept			json
//	@Produce		json
//
//	@Param			body			body	changePasswordBody	true	"Passwords"
//
//	@Param			Authorization	header	string				true	"bearer {token}"
//
//	@Success		200				"No Content"
//
//	@Failure		400				{object}	app_errors.ErrValidation	"Bad Request"
//
//	@Failure		401				{object}	app_errors.OtherErrors		"Unauthorized"
//	@Failure		404				{object}	app_errors.OtherErrors		"Not Found"
//	@Failure		500				{object}	app_errors.OtherErrors		"Internal Server Error"
//
//	@Router			/v1/user/password [post]
func getUserPasswordRoute() {
	engine.Router().POST(
		"/v1/user/password",
		engine.ValidateLoggedIn,
		validateChangePasswordBody,
		changePassword,
	)
}

type changePasswordBody struct {
	Current string `json:"currentPassword" binding:"required,min=1,max=100"`
	New     string `json:"newPassword" binding:"required,min=1,max=100"`
}

func changePassword(c *gin.Context) {
	body := c.MustGet("data").(changePasswordBody)

	var extraParams []interface{}
	if mocks, ok := c.Get("mocks"); ok {
		extraParams = mocks.([]interface{})
	}

	payload := engine.HeaderToken(c)
	if err := user.ChangePassword(payload.UserID.Hex(), body.Current, body.New, extraParams...); err != nil {
		engine.AbortWithError(c, err)
		return
	}
}

func validateChangePasswordBody(c *gin.Context) {
	body := changePasswordBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		engine.AbortWithError(c, err)
		return
	}

	c.Set("data", body)
	c.Next()
}
