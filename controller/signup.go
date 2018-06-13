package controller

import (
	"github.com/nmarsollier/authgo/tools/errors"
	"github.com/nmarsollier/authgo/user"

	"github.com/gin-gonic/gin"
)

// SignUp registra usuarios nuevos
/**
 * @api {post} /auth/signup Registrar Usuario
 * @apiName signup
 * @apiGroup Seguridad
 *
 * @apiDescription Registra un nuevo usuario en el sistema.
 *
 * @apiParamExample {json} Body
 *    {
 *      "name": "{Nombre de Usuario}",
 *      "login": "{Login de usuario}",
 *      "password": "{Contraseña}"
 *    }
 *
 * @apiUse TokenResponse
 *
 * @apiUse ParamValidationErrors
 * @apiUse OtherErrors
 */
func SignUp(c *gin.Context) {
	body := user.SignUpRequest{}

	if err := c.ShouldBindJSON(&body); err != nil {
		errors.Handle(c, err)
		return
	}

	token, err := user.SignUp(&body)

	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.JSON(200, gin.H{
		"token": token,
	})

}
