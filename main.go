package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"github.com/nmarsollier/authgo/routes"
	"github.com/nmarsollier/authgo/tools/env"
)

func main() {
	if len(os.Args) > 1 {
		env.Load(os.Args[1])
	}

	server := gin.Default()

	server.Use(gzip.Gzip(gzip.DefaultCompression))

	server.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	server.Use(static.Serve("/", static.LocalFile(env.Get().WWWWPath, true)))

	server.POST("/v1/user/password", routes.ChangePassword)
	server.POST("/v1/user/signin", routes.SignIn)
	server.GET("/v1/user/signout", routes.SignOut)
	server.POST("/v1/user", routes.SignUp)
	server.GET("/v1/users/current", routes.CurrentUser)
	server.POST("/v1/users/:userID/grant", routes.GrantPermission)
	server.POST("/v1/users/:userID/revoke", routes.RevokePermission)
	server.POST("/v1/users/:userID/enable", routes.Enable)
	server.POST("/v1/users/:userID/disable", routes.Disable)
	server.GET("/v1/users", routes.Users)

	server.Run(fmt.Sprintf(":%d", env.Get().Port))
}
