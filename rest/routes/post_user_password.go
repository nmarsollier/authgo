package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/model/user"
	"github.com/nmarsollier/authgo/rest/middlewares"
)

/**
 * @api {post} /v1/user/password Cambiar Password
 * @apiName Cambiar Password
 * @apiGroup Seguridad
 *
 * @apiDescription Cambia la contraseña del usuario actual.
 *
 * @apiExample {json} Body
 *    {
 *      "currentPassword" : "{Contraseña actual}",
 *      "newPassword" : "{Nueva Contraseña}",
 *    }
 *
 * @apiSuccessExample {json} Respuesta
 *     HTTP/1.1 200 OK
 *
 * @apiUse AuthHeader
 * @apiUse ParamValidationErrors
 * @apiUse OtherErrors
 */

func init() {
	router().POST(
		"/v1/user/password",
		middlewares.ValidateLoggedIn,
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

	payload := middlewares.HeaderToken(c)
	if err := user.ChangePassword(payload.UserID.Hex(), body.Current, body.New); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}
}

func validateChangePasswordBody(c *gin.Context) {
	body := changePasswordBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.Set("data", body)
	c.Next()
}
