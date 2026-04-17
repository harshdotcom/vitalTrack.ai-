package routes

import (
	"vita-track-ai/models"
	"vita-track-ai/service"

	"github.com/gin-gonic/gin"
)

// SaveHealthMetric godoc
// @Summary Save health metric
// @Description Create a new health metric entry for the logged-in user
// @Tags Health
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.SaveHealthMetricRequest true "Health metric payload"
// @Success 201 {object} map[string]interface{} "health metric saved successfully"
// @Failure 400 {object} map[string]interface{} "invalid request"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /health-metric/save [post]
func saveHealthMetric(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)
	var req models.SaveHealthMetricRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	metric, err := service.SaveHealtMetric(req, userID)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Success response
	c.JSON(201, gin.H{
		"message": "health metric saved successfully",
		"data":    metric,
	})
}

func deleteHealthMetric(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)
	metricID := c.Param("id")

	err := service.DeleteHealthMetric(metricID, userID)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "health metric deleted successfully",
	})
}
