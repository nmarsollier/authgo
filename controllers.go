package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/errors"
	"github.com/nmarsollier/authgo/user"
)

type permission struct {
	Permissions []string `json:"permissions" binding:"required"`
}

// GrantPermission  Otorga Permisos
/**
 * @api {post} /users/:userId/grant Otorga Permisos
 * @apiName Otorga Permisos
 * @apiGroup Seguridad
 *
 * @apiDescription Otorga permisos al usuario indicado, el usuario logueado tiene que tener permiso "admin".
 *
 * @apiParamExample {json} Body
 *    {
 *      "permissions" : ["permiso", ...],
 *    }
 *
 * @apiSuccessExample {json} Respuesta
 *     HTTP/1.1 200 OK
 *
 * @apiUse AuthHeader
 * @apiUse ParamValidationErrors
 * @apiUse OtherErrors
 */
func GrantPermission(c *gin.Context) {
	tokenString, err := getTokenHeader(c)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	payload, err := token.Validate(tokenString)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	body := permission{}

	if err := c.ShouldBindJSON(&body); err != nil {
		errors.Handle(c, err)
		return
	}

	if !user.Granted(payload.UserID, "admin") {
		errors.Handle(c, errors.AccessLevel)
		return
	}

	userID := c.Param("userID")
	err = user.Grant(userID, body.Permissions)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.Done()
}

// RevokePermission  Revoca Permisos
/**
 * @api {post} /users/:userId/revoke Revoca Permisos
 * @apiName Revoca Permisos
 * @apiGroup Seguridad
 *
 * @apiDescription Quita permisos al usuario indicado, el usuario logueado tiene que tener permiso "admin".
 *
 * @apiParamExample {json} Body
 *    {
 *      "permissions" : ["permiso", ...],
 *    }
 *
 * @apiSuccessExample {json} Respuesta
 *     HTTP/1.1 200 OK
 *
 * @apiUse AuthHeader
 * @apiUse ParamValidationErrors
 * @apiUse OtherErrors
 */
func RevokePermission(c *gin.Context) {
	tokenString, err := getTokenHeader(c)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	payload, err := token.Validate(tokenString)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	body := permission{}

	if err := c.ShouldBindJSON(&body); err != nil {
		errors.Handle(c, err)
		return
	}

	if !user.Granted(payload.UserID, "admin") {
		errors.Handle(c, errors.AccessLevel)
		return
	}

	userID := c.Param("userID")
	err = user.Revoke(userID, body.Permissions)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.Done()
}

// Disable  Deshabilita un usuario
/**
 * @api {post} /users/:userId/disable Deshabilitar Usuario
 * @apiName Deshabilitar Usuario
 * @apiGroup Seguridad
 *
 * @apiDescription Deshabilita un usuario en el sistema.   El usuario logueado debe tener permisos "admin".
 *
 * @apiSuccessExample {json} Respuesta
 *     HTTP/1.1 200 OK
 *
 * @apiUse AuthHeader
 * @apiUse ParamValidationErrors
 * @apiUse OtherErrors
 */
func Disable(c *gin.Context) {
	tokenString, err := getTokenHeader(c)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	payload, err := token.Validate(tokenString)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	if !user.Granted(payload.UserID, "admin") {
		errors.Handle(c, errors.AccessLevel)
		return
	}

	userID := c.Param("userID")
	err = user.Disable(userID)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.Done()
}

// Enable  Habilita un usuario
/**
 * @api {post} /users/:userId/enable Habilitar Usuario
 * @apiName Habilitar Usuario
 * @apiGroup Seguridad
 *
 * @apiDescription Habilita un usuario en el sistema. El usuario logueado debe tener permisos "admin".
 *
 * @apiSuccessExample {json} Respuesta
 *     HTTP/1.1 200 OK
 *
 * @apiUse AuthHeader
 * @apiUse ParamValidationErrors
 * @apiUse OtherErrors
 */
func Enable(c *gin.Context) {
	tokenString, err := getTokenHeader(c)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	payload, err := token.Validate(tokenString)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	if !user.Granted(payload.UserID, "admin") {
		errors.Handle(c, errors.AccessLevel)
		return
	}

	userID := c.Param("userID")
	err = user.Enable(userID)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.Done()
}

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
	tokenString, err := getTokenHeader(c)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	err = token.Invalidate(tokenString)

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
	tokenString, err := getTokenHeader(c)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	payload, err := token.Validate(tokenString)

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
		Current string `json:"currentPassword" binding:"required,min=1,max=100"`
		New     string `json:"newPassword" binding:"required,min=1,max=100"`
	}

	tokenString, err := getTokenHeader(c)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	payload, err := token.Validate(tokenString)

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
