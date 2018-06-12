package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/errors"
	"github.com/nmarsollier/authgo/user"
)

/**
 * @api {post} /auth/password Cambiar Password
 * @apiName ChangePassword
 * @apiGroup Seguridad
 *
 * @apiDescription Cambia la contraseña del usuario actual.
 *
 * @apiParamExample {json} Body
 *    {
 *      "currentPassword" : "{Contraseña actual}",
 *      "verifyPassword" : "{Contraseña actual}"
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
// ChangePassword Change Password Controller
func ChangePassword(c *gin.Context) {
	payload, err := token.ValidateToken(c)

	if err != nil {
		errors.HandleError(c, err)
		return
	}

	body := changePasswordRequest{}

	if err := c.ShouldBindJSON(&body); err != nil {
		errors.HandleError(c, err)
		return
	}

	err = user.ChangePassword(payload.UserID, body.CurrentPassword, body.NewPassword)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.Done()
}

type changePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required"`
}
