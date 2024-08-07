package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/authgo/rest/engine"
	"github.com/nmarsollier/authgo/user"
)

type UserResponse struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
	Login       string   `json:"login"`
}

// Obtiene información del usuario actual.
//
//	@Summary		Usuario Actual
//	@Description	Obtiene información del usuario actual.
//	@Tags			Seguridad
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string						true	"bearer {token}"
//	@Success		200				{object}	UserResponse				"User data"
//
//	@Failure		400				{object}	app_errors.ErrValidation	"Bad Request"
//	@Failure		401				{object}	app_errors.OtherErrors		"Unauthorized"
//	@Failure		404				{object}	app_errors.OtherErrors		"Not Found"
//	@Failure		500				{object}	app_errors.OtherErrors		"Internal Server Error"
//
//	@Router			/v1/users/current [get]
func getUsersCurrentRoute() {
	engine.Router().GET(
		"/v1/users/current",
		engine.ValidateLoggedIn,
		currentUser,
	)
}

func currentUser(c *gin.Context) {
	var extraParams []interface{}
	if mocks, ok := c.Get("mocks"); ok {
		extraParams = mocks.([]interface{})
	}

	token := engine.HeaderToken(c)

	user, err := user.Get(token.UserID.Hex(), extraParams...)
	if err != nil {
		engine.AbortWithError(c, err)
		return
	}

	c.JSON(200, UserResponse{
		Id:          user.ID.Hex(),
		Name:        user.Name,
		Permissions: user.Permissions,
		Login:       user.Login,
	})
}
