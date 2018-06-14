package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/token"
	"github.com/nmarsollier/authgo/tools/errors"
	"github.com/nmarsollier/authgo/user"
)

// GrantPermission  Otorga Permisos
/**
 * @api {post} /auth/:userId/grant Otorga Permisos
 * @apiName Grant
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
	payload, err := token.Validate(c)
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
 * @api {post} /auth/:userId/revoke Revoca Permisos
 * @apiName Revoke
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
	payload, err := token.Validate(c)
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
 * @api {post} /auth/:userId/disable Deshabilita un usuario.  El usuario logueado debe tener permisos "admin".
 * @apiName Disable
 * @apiGroup Seguridad
 *
 * @apiDescription Deshabilita un usuario en el sistema.
 *
 * @apiSuccessExample {json} Respuesta
 *     HTTP/1.1 200 OK
 *
 * @apiUse AuthHeader
 * @apiUse ParamValidationErrors
 * @apiUse OtherErrors
 */
func Disable(c *gin.Context) {
	payload, err := token.Validate(c)
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

// Enable  Deshabilita un usuario
/**
 * @api {post} /auth/:userId/enable Habilita un usuario
 * @apiName Enable
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
	payload, err := token.Validate(c)
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

type permission struct {
	Permissions []string `json:"permissions" binding:"required"`
}
