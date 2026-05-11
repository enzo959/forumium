package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/enzo959/forumium/config"
	"github.com/enzo959/forumium/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: no .env file found")
	}
	config.ConnectDatabase()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:5173"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    AllowCredentials: true,
	}))

	routes.SetupRoutes(router)

	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8080"
	}

	err = router.Run(":" + port)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}