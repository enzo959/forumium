package routes

import (
	"github.com/enzo959/forumium/handlers"
	"github.com/enzo959/forumium/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	router.GET("/api/health", handlers.HealthCheck)

	auth := router.Group("/api/auth")

	auth.POST("/register", handlers.Register)
	auth.POST("/login", handlers.Login)
	auth.POST("/refresh", handlers.Refresh)
	auth.POST("/forgot-password", handlers.ForgotPassword)

	protected := router.Group("/api")
	protected.Use(middleware.AuthRequired())

	protected.POST("/auth/logout", handlers.Logout)

}
