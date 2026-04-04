package routes

import (
	"vita-track-ai/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {

	api := server.Group("/api/v1")

	// PUBLIC routes (no auth)
	registerUserRoutes(api)

	// PROTECTED routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())

	registerFileRoutes(protected)
	registerDocumentRoutes(protected)
	registerUserDetailRoutes(protected)
}

func registerUserRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	{
		users.POST("/signup", signup)
		users.POST("/verify-otp", verifyOTP)
		users.POST("/login", login)
		users.POST("/google", googleLogin)
		users.POST("/forgot-password", forgotPassword)
		users.POST("/reset-password", resetPassword)
	}
}

func registerFileRoutes(rg *gin.RouterGroup) {
	files := rg.Group("/files")
	{
		files.POST("/upload", uploadFile)
		files.GET("/:id", getFile)
		files.GET("/ocr/:id", getFileText)
		files.GET("/ai/:id", getFileAnalysis)
		files.DELETE("/:id", deleteFile)
	}
}

// ⭐ NEW
func registerDocumentRoutes(rg *gin.RouterGroup) {
	documents := rg.Group("/documents")
	{
		documents.POST("", createDocument)
		documents.GET("/:id", getDocument)
		documents.DELETE("/:id", deleteDocument)
		documents.POST("/calendar", getCalendarDocuments)
		documents.PATCH("/update/:id", updateDocument)
	}
}

func registerUserDetailRoutes(rg *gin.RouterGroup) {
	userDetails := rg.Group("/user-details")
	{
		userDetails.GET("/usage", getUserUsage)
		userDetails.GET("/ai-credits", getAICreditUsage)
		userDetails.PATCH("/update", updateProfile)
	}

}
