package routes

import (
	"github.com/enzo959/forumium/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/api/health", handlers.HealthCheck)

	router.POST("/api/auth/register", handlers.Register)
	router.POST("/api/auth/login", handlers.Login)
	router.POST("/api/auth/refresh", handlers.Refresh)
}
