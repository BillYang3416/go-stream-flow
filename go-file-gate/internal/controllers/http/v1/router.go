package v1

import (
	"net/http"

	"github.com/bgg/go-file-gate/internal/usecase"
	"github.com/bgg/go-file-gate/pkg/logger"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Swagger docs.
	_ "github.com/bgg/go-file-gate/docs"
)

func NewRouter(handler *gin.Engine, l logger.Logger, u usecase.UserProfile) {

	// logging each http request
	handler.Use(gin.Logger())

	// Routers
	h := handler.Group("/api/v1")
	{
		newUserProfileRoutes(h, u, l)
		NewAuthRoutes(h, u, l)
	}

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// Serve static files
	// handler.StaticFS("/static", gin.Dir("/internal/static", false)) for docker
	handler.StaticFS("/static", gin.Dir("../../internal/static", false))

	// Handle SPA client-side routing
	handler.NoRoute(func(c *gin.Context) {
		// c.File("/internal/static/index.html") for docker
		c.File("../../internal/static/index.html")
	})
}

func CheckSessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID, exists := session.Get("userID").(string)
		if !exists || userID == "" {
			sendErrorResponse(c, http.StatusUnauthorized, "unauthorized")
			return
		}
		c.Next()
	}
}
