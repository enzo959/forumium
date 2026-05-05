package routes

import (
	"github.com/enzo959/forumium/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/api/health", handlers.HealthCheck)
}

func SetupforgotPasswordRoutes(router *gin.Engine) {
	router.GET("/api/auth/forgot-password", handlers.forgotPassword)
}