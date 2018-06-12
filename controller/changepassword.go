package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/errors"
	"github.com/nmarsollier/authgo/user"
)

// ChangePassword Change Password Controller
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
func ChangePassword(c *gin.Context) {
	payload, err := token.Validate(c)

	if err != nil {
		errors.Handle(c, err)
		return
	}

	body := changePassword{}

	if err := c.ShouldBindJSON(&body); err != nil {
		errors.Handle(c, err)
		return
	}

	err = user.ChangePassword(payload.UserID, body.Current, body.New)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.Done()
}

type changePassword struct {
	Current string `json:"currentPassword" binding:"required"`
	New     string `json:"newPassword" binding:"required"`
}
