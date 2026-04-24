package main

import (
	"TestProject/config"
	"TestProject/routes"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load() // load env variables from .env file
	if err != nil {
		fmt.Println("No .env file found, using environment variables")
	}
	config.MongoDB_Connection() // connected to manogoDB

	// CORS configuration
	corsConfig := cors.Config{ // cors configuration for allowing requests from any origin
		Origins:         "*",
		RequestHeaders:  "*",
		Methods:         "GET,POST,PUT,DELETE",
		Credentials:     false,
		ValidateHeaders: false,
		MaxAge:          1 * time.Minute,
	}

	// Create Gin router
	router := gin.Default()
	router.Use(cors.Middleware(corsConfig))

	// Setup routes
	routes.SetUpRoutes(router)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080" // Default port if not specified in .env
	}
	log.Printf("Server is running on port %s", PORT) // Start the server
	if err := router.Run(":" + PORT); err != nil {
		log.Println("No port found")
	}
}
