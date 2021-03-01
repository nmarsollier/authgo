package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/model/user"
	"github.com/nmarsollier/authgo/rest/middlewares"
)

/**
 * @api {post} /v1/user/signin Login
 * @apiName Login
 * @apiGroup Seguridad
 *
 * @apiDescription Loguea un usuario en el sistema.
 *
 * @apiExample {json} Body
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

func init() {
	router().POST(
		"/v1/user/signin",
		validateSignInBody,
		signIn,
	)
}

type signInBody struct {
	Password string `json:"password" binding:"required"`
	Login    string `json:"login" binding:"required"`
}

func signIn(c *gin.Context) {
	login := c.MustGet("data").(signInBody)

	tokenString, err := user.SignIn(login.Login, login.Password)
	if err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"token": tokenString,
	})
}

func validateSignInBody(c *gin.Context) {
	login := signInBody{}
	if err := c.ShouldBindJSON(&login); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.Set("data", login)
	c.Next()
}
