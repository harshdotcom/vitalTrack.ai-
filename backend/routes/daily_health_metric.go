package routes

import (
	"vita-track-ai/models"
	"vita-track-ai/service"

	"github.com/gin-gonic/gin"
)

func saveHealthMetric(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)
	var req models.SaveHealthMetricRequest
	req.UploadedBy = userID

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	metric, err := service.SaveHealtMetric(req)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}

	// Success response
	c.JSON(201, gin.H{
		"message": "health metric saved successfully",
		"data":    metric,
	})
}
