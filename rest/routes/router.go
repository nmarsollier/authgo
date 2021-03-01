package routes

import (
	"fmt"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"github.com/nmarsollier/authgo/rest/middlewares"
	"github.com/nmarsollier/authgo/tools/env"
)

// Start this server
func Start() {
	router().Run(fmt.Sprintf(":%d", env.Get().Port))
}

var engine *gin.Engine = nil

func router() *gin.Engine {
	if engine == nil {

		engine = gin.Default()

		engine.Use(gzip.Gzip(gzip.DefaultCompression))

		engine.Use(cors.Middleware(cors.Config{
			Origins:         "*",
			Methods:         "GET, PUT, POST, DELETE",
			RequestHeaders:  "Origin, Authorization, Content-Type",
			ExposedHeaders:  "",
			MaxAge:          50 * time.Second,
			Credentials:     true,
			ValidateHeaders: false,
		}))

		engine.Use(middlewares.ErrorHandler)

		engine.Use(static.Serve("/", static.LocalFile(env.Get().WWWWPath, true)))
	}

	return engine
}
