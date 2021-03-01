package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/model/user"
	"github.com/nmarsollier/authgo/rest/middlewares"
)

/**
 * @api {post} /v1/user Registrar Usuario
 * @apiName Registrar Usuario
 * @apiGroup Seguridad
 *
 * @apiDescription Registra un nuevo usuario en el sistema.
 *
 * @apiExample {json} Body
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

func init() {
	router().POST(
		"/v1/user",
		validateSignUpBody,
		signUp,
	)
}

func signUp(c *gin.Context) {
	body := c.MustGet("data").(user.SignUpRequest)

	token, err := user.SignUp(&body)
	if err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"token": token,
	})
}

func validateSignUpBody(c *gin.Context) {
	body := user.SignUpRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		middlewares.AbortWithError(c, err)
		return
	}

	c.Set("data", body)
	c.Next()
}
