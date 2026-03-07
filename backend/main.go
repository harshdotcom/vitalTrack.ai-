package main

import (
	"time"

	"vita-track-ai/config"
	"vita-track-ai/database"
	"vita-track-ai/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// Load env
	godotenv.Load()

	// Init DB & S3
	database.Init()
	config.InitS3()

	// Create server
	server := gin.Default()

	// =========================
	// CORS CONFIGURATION
	// =========================
	server.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Register routes
	routes.RegisterRoutes(server)

	// Run server
	server.Run(":8081")
}
