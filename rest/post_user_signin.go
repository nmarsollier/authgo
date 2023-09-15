package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/rest/engine"
	"github.com/nmarsollier/authgo/user"
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

func postUserSignInRoute() {
	engine.Router().POST(
		"/v1/user/signin",
		validateSignInBody,
		signIn,
	)
}

func signIn(c *gin.Context) {
	login := c.MustGet("data").(user.SignInRequest)

	var extraParams []interface{}
	if mocks, ok := c.Get("mocks"); ok {
		extraParams = mocks.([]interface{})
	}

	tokenString, err := user.SignIn(login, extraParams...)
	if err != nil {
		engine.AbortWithError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"token": tokenString,
	})
}

func validateSignInBody(c *gin.Context) {
	login := user.SignInRequest{}
	if err := c.ShouldBindJSON(&login); err != nil {
		engine.AbortWithError(c, err)
		return
	}

	c.Set("data", login)
	c.Next()
}
