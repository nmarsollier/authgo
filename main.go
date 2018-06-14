package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"github.com/nmarsollier/authgo/controller"
	"github.com/nmarsollier/authgo/tools/env"
)

func main() {
	if len(os.Args) > 1 {
		env.Load(os.Args[1])
	}

	r := gin.Default()

	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	r.Use(static.Serve("/", static.LocalFile(env.Get().WWWWPath, false)))

	r.POST("/auth/password", controller.ChangePassword)
	r.POST("/auth/signin", controller.SignIn)
	r.GET("/auth/signout", controller.SignOut)
	r.POST("/auth/signup", controller.SignUp)
	r.GET("/auth/currentUser", controller.CurrentUser)
	r.POST("/users/:userID/grant", controller.GrantPermission)
	r.POST("/users/:userID/revoke", controller.RevokePermission)
	r.POST("/users/:userID/enable", controller.Enable)
	r.POST("/users/:userID/disable", controller.Disable)
	r.Run(fmt.Sprintf(":%d", env.Get().Port))
}

func preflight(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
	c.JSON(http.StatusOK, struct{}{})
}
