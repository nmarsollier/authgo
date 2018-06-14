package controller

// Son controllers de usuario

import (
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/errors"
	"github.com/nmarsollier/authgo/user"

	"github.com/gin-gonic/gin"
)

// SignUp registra usuarios nuevos
/**
 * @api {post} /user Registrar Usuario
 * @apiName Registrar Usuario
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

// SignOut is the sign out controller
/**
 * @api {get} /user/signout Logout
 * @apiName Logout
 * @apiGroup Seguridad
 *
 * @apiDescription Desloguea un usuario en el sistema, invalida el token.
 *
 * @apiSuccessExample {json} Respuesta
 *     HTTP/1.1 200 OK
 *
 * @apiUse AuthHeader
 * @apiUse OtherErrors
 */
func SignOut(c *gin.Context) {
	err := token.Invalidate(c)

	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.Done()
}

// SignIn is the controller to sign in users
/**
 * @api {post} /user/signin Login
 * @apiName Login
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
	type signIn struct {
		Password string `json:"password" binding:"required"`
		Login    string `json:"login" binding:"required"`
	}

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

// CurrentUser is the controller to get the current logged in user
/**
 * @api {get} /users/current Usuario Actual
 * @apiName Usuario Actual
 * @apiGroup Seguridad
 *
 * @apiDescription Obtiene información del usuario actual.
 *
 * @apiSuccessExample {json} Respuesta
 *     HTTP/1.1 200 OK
 *     {
 *        "id": "{Id usuario}",
 *        "name": "{Nombre del usuario}",
 *        "login": "{Login de usuario}",
 *        "permissions": [
 *            "{Rol}"
 *        ]
 *     }
 *
 * @apiUse AuthHeader
 * @apiUse OtherErrors
 */
func CurrentUser(c *gin.Context) {
	payload, err := token.Validate(c)

	if err != nil {
		errors.Handle(c, err)
		return
	}

	user, err := user.GetUser(payload.UserID)

	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.JSON(200, gin.H{
		"id":          user.ID.Hex(),
		"name":        user.Name,
		"permissions": user.Permissions,
		"login":       user.Login,
	})
}

// ChangePassword Change Password Controller
/**
 * @api {post} /user/password Cambiar Password
 * @apiName Cambiar Password
 * @apiGroup Seguridad
 *
 * @apiDescription Cambia la contraseña del usuario actual.
 *
 * @apiParamExample {json} Body
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
func ChangePassword(c *gin.Context) {
	type changePassword struct {
		Current string `json:"currentPassword" binding:"required"`
		New     string `json:"newPassword" binding:"required"`
	}

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
