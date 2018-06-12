package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/tools/errors"
	"github.com/nmarsollier/authgo/user"
)

// SignIn is the controller to sign in users
/**
 * @api {post} /auth/signin Login
 * @apiName Log in
 * @apiGroup Seguridad
 *
 * @apiDescription Loguea un usuario en el sistema.
 *
 * @apiParamExample {json} Body
 *    {
 *      "login": "{Login de usuario}",
 *      "password": "{Contraseña}"
 *    }
 *
 * @apiUse TokenResponse
 *
 * @apiUse ParamValidationErrors
 * @apiUse OtherErrors
 */
func SignIn(c *gin.Context) {
	login := signIn{}

	if err := c.ShouldBindJSON(&login); err != nil {
		errors.Handle(c, err)
		return
	}

	tokenString, err := user.SignIn(login.Login, login.Password)

	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.JSON(200, gin.H{
		"token": tokenString,
	})
}

type signIn struct {
	Password string `json:"password" binding:"required"`
	Login    string `json:"login" binding:"required"`
}
