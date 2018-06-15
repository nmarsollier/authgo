package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
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

	server.POST("/v1/user/password", ChangePassword)
	server.POST("/v1/user/signin", SignIn)
	server.GET("/v1/user/signout", SignOut)
	server.POST("/v1/user", SignUp)
	server.GET("/v1/users/current", CurrentUser)
	server.POST("/v1/users/:userID/grant", GrantPermission)
	server.POST("/v1/users/:userID/revoke", RevokePermission)
	server.POST("/v1/users/:userID/enable", Enable)
	server.POST("/v1/users/:userID/disable", Disable)

	server.Run(fmt.Sprintf(":%d", env.Get().Port))
}
