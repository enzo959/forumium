package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/enzo959/forumium/backend/config"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: no .env file found")
	}
	config.ConnectDatabase()

	router := gin.Default()

	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8080"
	}

	err = router.Run(":" + port)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
