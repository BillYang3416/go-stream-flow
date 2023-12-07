package v1

import (
	"net/http"

	"github.com/bgg/go-file-gate/config"
	"github.com/bgg/go-file-gate/internal/usecase"
	"github.com/bgg/go-file-gate/pkg/logger"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Swagger docs.
	_ "github.com/bgg/go-file-gate/docs"
)

func NewRouter(cfg *config.Config, handler *gin.Engine, l logger.Logger, u usecase.UserProfile, uu usecase.UserUploadedFile) {

	// logging each http request
	handler.Use(gin.Logger())

	// Routers
	h := handler.Group("/api/v1")
	{
		NewUserProfileRoutes(h, u, l)
		NewAuthRoutes(cfg, h, u, l)
		NewUserUploadedFileRoutes(h, uu, l)
	}

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)
}

func CheckSessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID, exists := session.Get("userID").(string)
		if !exists || userID == "" {
			sendErrorResponse(c, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}
