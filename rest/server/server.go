package server

import (
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	_ "github.com/nmarsollier/authgo/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var engine *gin.Engine = nil

func Router(deps ...interface{}) *gin.Engine {
	if engine == nil {

		engine = gin.Default()

		engine.Use(gzip.Gzip(gzip.DefaultCompression))
		engine.Use(GinLoggerMiddleware(deps...))

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

		engine.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return engine
}

func AbortWithError(c *gin.Context, err error) {
	c.Error(err)
	c.Abort()
}

// Obtiene el contexto a serivcios externos
// En prod este contexto esta vacio.
func GinDeps(c *gin.Context) []interface{} {
	var deps []interface{}
	// mock_deps solo es para mocks en testing
	if mocks, ok := c.Get("mock_deps"); ok {
		deps = mocks.([]interface{})
	}

	deps = append(deps, ginLogger(c))

	return deps
}
