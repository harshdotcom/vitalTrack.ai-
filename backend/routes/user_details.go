package routes

import (
	"net/http"
	"vita-track-ai/repository"

	"github.com/gin-gonic/gin"
)

func getUserUsage(context *gin.Context) {
	userId := context.MustGet("user_id").(int64)

	usage, err := repository.GetCurrentStorageUsed(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch storage usage",
			"error":   err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Storage usage fetched successfully",
		"usage":   usage,
	})
}

// func updateProfile(context *gin.Context) {}
