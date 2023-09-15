package engine

import (
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"github.com/nmarsollier/authgo/tools/env"
)

var engine *gin.Engine = nil

func TestRouter(props ...interface{}) *gin.Engine {
	engine = nil
	Router()
	if len(props) > 0 {
		engine.Use(func(c *gin.Context) {
			c.Set("mocks", props)
			c.Next()
		})
	}

	return engine
}

func Router() *gin.Engine {
	if engine == nil {
		engine = gin.Default()

		engine.Use(gzip.Gzip(gzip.DefaultCompression))

		engine.Use(cors.Middleware(cors.Config{
			Origins:         "*",
			Methods:         "GET, PUT, POST, DELETE",
			RequestHeaders:  "Origin, Authorization, Content-Type",
			ExposedHeaders:  "",
			MaxAge:          50 * time.Second,
			Credentials:     false,
			ValidateHeaders: false,
		}))

		engine.Use(ErrorHandler)

		engine.Use(static.Serve("/", static.LocalFile(env.Get().WWWWPath, true)))
	}

	return engine
}
