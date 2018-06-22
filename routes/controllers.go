package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/security"
	"github.com/nmarsollier/authgo/tools/errors"
	"github.com/nmarsollier/authgo/user"
)

type permission struct {
	Permissions []string `json:"permissions" binding:"required"`
}

// GrantPermission  Otorga Permisos
/**
 * @api {post} /v1/users/:userId/grant Otorga Permisos
 * @apiName Otorga Permisos
 * @apiGroup Seguridad
 *
 * @apiDescription Otorga permisos al usuario indicado, el usuario logueado tiene que tener permiso "admin".
 *
 * @apiExample {json} Body
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
	token, err := validateAuthHeader(c)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	body := permission{}
	if err := c.ShouldBindJSON(&body); err != nil {
		errors.Handle(c, err)
		return
	}

	userService, err := user.NewService()
	if err != nil {
		errors.Handle(c, err)
		return
	}

	if !userService.Granted(token.UserID.Hex(), "admin") {
		errors.Handle(c, errors.AccessLevel)
		return
	}

	err = userService.Grant(c.Param("userID"), body.Permissions)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.Done()
}

// RevokePermission  Revoca Permisos
/**
 * @api {post} /v1/users/:userId/revoke Revoca Permisos
 * @apiName Revoca Permisos
 * @apiGroup Seguridad
 *
 * @apiDescription Quita permisos al usuario indicado, el usuario logueado tiene que tener permiso "admin".
 *
 * @apiExample {json} Body
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
	token, err := validateAuthHeader(c)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	body := permission{}
	if err := c.ShouldBindJSON(&body); err != nil {
		errors.Handle(c, err)
		return
	}

	userService, err := user.NewService()
	if err != nil {
		errors.Handle(c, err)
		return
	}

	if !userService.Granted(token.UserID.Hex(), "admin") {
		errors.Handle(c, errors.AccessLevel)
		return
	}

	err = userService.Revoke(c.Param("userID"), body.Permissions)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.Done()
}

// Disable  Deshabilita un usuario
/**
 * @api {post} /v1/users/:userId/disable Deshabilitar Usuario
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
	payload, err := validateAuthHeader(c)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	userService, err := user.NewService()
	if err != nil {
		errors.Handle(c, err)
		return
	}

	if !userService.Granted(payload.UserID.Hex(), "admin") {
		errors.Handle(c, errors.AccessLevel)
		return
	}

	err = userService.Disable(c.Param("userID"))
	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.Done()
}

// Enable  Habilita un usuario
/**
 * @api {post} /v1/users/:userId/enable Habilitar Usuario
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
	payload, err := validateAuthHeader(c)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	userService, err := user.NewService()
	if err != nil {
		errors.Handle(c, err)
		return
	}

	if !userService.Granted(payload.UserID.Hex(), "admin") {
		errors.Handle(c, errors.AccessLevel)
		return
	}

	err = userService.Enable(c.Param("userID"))
	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.Done()
}

// SignUp registra usuarios nuevos
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
func SignUp(c *gin.Context) {
	body := user.SignUpRequest{}
	if err := c.ShouldBindJSON(&body); err != nil {
		errors.Handle(c, err)
		return
	}

	userService, err := user.NewService()
	if err != nil {
		errors.Handle(c, err)
		return
	}

	token, err := userService.SignUp(&body)
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
 * @api {get} /v1/user/signout Logout
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
	tokenString, err := getAuthHeader(c)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	secService, err := security.NewService()
	if err != nil {
		errors.Handle(c, err)
		return
	}

	if err = secService.Invalidate(tokenString); err != nil {
		errors.Handle(c, err)
		return
	}

	c.Done()
}

// SignIn is the controller to sign in users
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

	userService, err := user.NewService()
	if err != nil {
		errors.Handle(c, err)
		return
	}

	tokenString, err := userService.SignIn(login.Login, login.Password)
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
 * @api {get} /v1/users/current Usuario Actual
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
 *            "{Permission}"
 *        ]
 *     }
 *
 * @apiUse AuthHeader
 * @apiUse OtherErrors
 */
func CurrentUser(c *gin.Context) {
	token, err := validateAuthHeader(c)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	userService, err := user.NewService()
	if err != nil {
		errors.Handle(c, err)
		return
	}

	user, err := userService.Get(token.UserID.Hex())

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
func ChangePassword(c *gin.Context) {
	payload, err := validateAuthHeader(c)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	type changePassword struct {
		Current string `json:"currentPassword" binding:"required,min=1,max=100"`
		New     string `json:"newPassword" binding:"required,min=1,max=100"`
	}
	body := changePassword{}
	if err := c.ShouldBindJSON(&body); err != nil {
		errors.Handle(c, err)
		return
	}

	userService, err := user.NewService()
	if err != nil {
		errors.Handle(c, err)
		return
	}

	err = userService.ChangePassword(payload.UserID.Hex(), body.Current, body.New)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	c.Done()
}

// Users Devuelve una lista de todos los usuarios
/**
 * @api {get} /v1/users Listar Usuarios
 * @apiName Listar Usuarios
 * @apiGroup Seguridad
 *
 * @apiDescription Obtiene información de todos los usuarios.
 *
 * @apiSuccessExample {json} Respuesta
 *     HTTP/1.1 200 OK
 *     [{
 *        "id": "{Id usuario}",
 *        "name": "{Nombre del usuario}",
 *        "login": "{Login de usuario}",
 *        "permissions": [
 *            "{Permission}"
 *        ],
 * 	      "enabled": true|false
 *     }, ...]
 *
 * @apiUse AuthHeader
 * @apiUse OtherErrors
 */
func Users(c *gin.Context) {
	payload, err := validateAuthHeader(c)
	if err != nil {
		errors.Handle(c, err)
		return
	}

	userService, err := user.NewService()
	if err != nil {
		errors.Handle(c, err)
		return
	}

	if !userService.Granted(payload.UserID.Hex(), "admin") {
		errors.Handle(c, errors.AccessLevel)
		return
	}

	user, err := userService.Users()

	if err != nil {
		errors.Handle(c, err)
		return
	}
	result := []gin.H{}
	for i := 0; i < len(user); i = i + 1 {
		result = append(result, gin.H{
			"id":          user[i].ID.Hex(),
			"name":        user[i].Name,
			"permissions": user[i].Permissions,
			"login":       user[i].Login,
			"enabled":     user[i].Enabled,
		})
	}

	c.JSON(200, result)
}
