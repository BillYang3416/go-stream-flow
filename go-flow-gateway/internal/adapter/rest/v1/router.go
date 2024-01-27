package v1

import (
	"net/http"
	"os"

	"github.com/bgg/go-flow-gateway/config"
	"github.com/bgg/go-flow-gateway/internal/usecase"
	"github.com/bgg/go-flow-gateway/pkg/logger"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Swagger docs.
	_ "github.com/bgg/go-flow-gateway/docs"
)

func NewRouter(cfg *config.Config, handler *gin.Engine, l logger.Logger, u usecase.UserProfile, uu usecase.UserUploadedFile, o usecase.OAuthDetail, c usecase.UserCredential) {

	// logging each http request
	handler.Use(gin.Logger())

	// Routers
	h := handler.Group("/api/v1")
	{
		NewUserProfileRoutes(h, u, l)
		NewAuthRoutes(cfg, h, u, l, o, c, os.Getenv("LINE_CHANNEL_ID"))
		NewUserUploadedFileRoutes(h, uu, l)
	}

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)
}

func CheckSessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		_, exists := session.Get("userID").(int)
		if !exists {
			sendErrorResponse(c, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}
