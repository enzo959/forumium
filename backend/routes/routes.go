package routes

import (
	"github.com/enzo959/forumium/handlers"
	"github.com/enzo959/forumium/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	router.GET("/api/health", handlers.HealthCheck)
	router.GET("/api/posts", handlers.GetPosts)
	router.GET("/api/categories", handlers.GetCategories)
	router.GET("/api/posts/:id", handlers.GetPostByID)
	router.GET("/api/posts/:id/comments", handlers.GetCommentsByPostID)
	router.GET("/api/users/:id/profile", handlers.GetUserProfile)

	auth := router.Group("/api/auth")

	auth.POST("/register", handlers.Register)
	auth.POST("/login", handlers.Login)
	auth.POST("/refresh", handlers.Refresh)
	auth.POST("/forgot-password", handlers.ForgotPassword)
	auth.POST("reset-password", handlers.ResetPassword)

	protected := router.Group("/api")
	protected.Use(middleware.AuthRequired())
	protected.POST("/auth/logout", handlers.Logout)
	protected.POST("/posts", handlers.CreatePost)

	protected.PUT("/posts/:id", handlers.UpdatePost)
	protected.DELETE("/posts/:id", handlers.DeletePost)
	protected.POST("/upload", handlers.UploadImage)
	protected.POST("/posts/:id/comments", handlers.CreateComment)
	protected.DELETE("/comments/:id", handlers.DeleteComment)
	protected.POST("/posts/:id/react", handlers.ReactToPost)
	

}
