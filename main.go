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

	r.Use(static.Serve("/", static.LocalFile(env.Get().WWWWPath, true)))

	r.POST("/v1/user/password", ChangePassword)
	r.POST("/v1/user/signin", SignIn)
	r.GET("/v1/user/signout", SignOut)
	r.POST("/v1/user", SignUp)
	r.GET("/v1/users/current", CurrentUser)
	r.POST("/v1/users/:userID/grant", GrantPermission)
	r.POST("/v1/users/:userID/revoke", RevokePermission)
	r.POST("/v1/users/:userID/enable", Enable)
	r.POST("/v1/users/:userID/disable", Disable)
	r.Run(fmt.Sprintf(":%d", env.Get().Port))
}
