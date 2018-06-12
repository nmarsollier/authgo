package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"github.com/nmarsollier/ms_auth_go/controller"
	"github.com/nmarsollier/ms_auth_go/tools/config"
)

func main() {
	r := gin.Default()

	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	r.Use(static.Serve("/", static.LocalFile("www", false)))

	r.POST("/auth/password", controller.ChangePassword)
	r.POST("/auth/signin", controller.SignIn)
	r.GET("/auth/signout", controller.SignOut)
	r.POST("/auth/signup", controller.SignUp)
	r.GET("/auth/currentUser", controller.CurrentUser)
	r.Run(fmt.Sprintf(":%d", config.Environment().Port))
}

func preflight(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
	c.JSON(http.StatusOK, struct{}{})
}