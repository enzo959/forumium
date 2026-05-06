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
	router.POST("/api/auth/logout", handlers.Logout)
	router.POST("/api/auth/forgot-password", handlers.ForgotPassword)
	router.POST("/api/auth/reset-password", handlers.ResetPassword)
}
