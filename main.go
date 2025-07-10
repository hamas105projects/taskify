package main

import (
	"log"
	"os"

	"taskify/config"
	"taskify/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	router := gin.Default()
	config.ConnectDatabase()

	api := router.Group("/api")
	{
		routes.AuthRoutes(api)
		routes.ProjectRoutes(api)
		routes.TaskRoutes(api)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
