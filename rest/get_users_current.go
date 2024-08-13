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

// @Summary		Usuario Actual
// @Description	Obtiene información del usuario actual.
// @Tags			Seguridad
// @Accept			json
// @Produce		json
// @Param			Authorization	header		string					true	"bearer {token}"
// @Success		200				{object}	UserResponse			"User data"
// @Failure		400				{object}	apperr.ValidationErr	"Bad Request"
// @Failure		401				{object}	engine.ErrorData		"Unauthorized"
// @Failure		404				{object}	engine.ErrorData		"Not Found"
// @Failure		500				{object}	engine.ErrorData		"Internal Server Error"
// @Router			/v1/users/current [get]
//
// Obtiene información del usuario actual.
func getUsersCurrentRoute() {
	engine.Router().GET(
		"/v1/users/current",
		engine.ValidateLoggedIn,
		currentUser,
	)
}

func currentUser(c *gin.Context) {
	token := engine.HeaderToken(c)

	ctx := engine.TestCtx(c)
	user, err := user.Get(token.UserID.Hex(), ctx...)
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
