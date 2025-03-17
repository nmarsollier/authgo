package server

import (
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"github.com/nmarsollier/authgo/internal/token"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var engine *gin.Engine = nil

func Router() *gin.Engine {
	if engine == nil {
		engine = gin.Default()

		engine.Use(gzip.Gzip(gzip.DefaultCompression))
		engine.Use(DiInjectorMiddleware())

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

// GetCtxToken Token data from Authorization header
func GetCtxToken(c *gin.Context) *token.Token {
	return c.MustGet("auth_header").(*token.Token)
}

func loadTokenFromHeader(c *gin.Context) (*token.Token, error) {
	log := GinLogger(c)
	tokenString, err := GetHeaderToken(c)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	payload, err := token.Validate(log, tokenString)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	c.Set("auth_header", payload)

	return payload, nil
}
