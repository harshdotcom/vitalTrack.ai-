package main

import (
	"fmt"
	"time"

	"vita-track-ai/database"
	"vita-track-ai/routes"
	"vita-track-ai/service"

	_ "vita-track-ai/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title Vita Track AI API
// @version 1.0
// @description API documentation for Vita Track AI
// @host localhost:8081
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	// Load env
	godotenv.Load()

	// Init DB & S3
	database.Init()
	service.InitS3()

	// Create server
	server := gin.New()
	server.Use(gin.Recovery())
	server.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		statusColor := param.StatusCodeColor()
		methodColor := param.MethodColor()
		resetColor := param.ResetColor()

		return fmt.Sprintf("[GIN] %s | %s%3d%s | %12v | %s | %s%-7s%s %s\n",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			statusColor, param.StatusCode, resetColor,
			param.Latency,
			param.ClientIP,
			methodColor, param.Method, resetColor,
			param.Path,
		)
	}))

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
