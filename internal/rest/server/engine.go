package server

import (
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"github.com/nmarsollier/authgo/internal/token"
	"github.com/nmarsollier/commongo/rst"
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

		engine.Use(rst.ErrorHandler)

		engine.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return engine
}

// GetCtxToken Token data from Authorization header
func GetCtxToken(c *gin.Context) *token.Token {
	return c.MustGet("auth_header").(*token.Token)
}

func loadTokenFromHeader(c *gin.Context) (*token.Token, error) {
	di := GinDi(c)
	tokenString, err := rst.GetHeaderToken(c)
	if err != nil {
		di.Logger().Error(err)
		return nil, err
	}

	payload, err := di.TokenService().Validate(tokenString)
	if err != nil {
		di.Logger().Error(err)
		return nil, err
	}

	c.Set("auth_header", payload)

	return payload, nil
}
