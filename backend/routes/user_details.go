package routes

import (
	"net/http"
	"os"
	"strconv"
	"time"
	"vita-track-ai/models"
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

func getAICreditUsage(context *gin.Context) {
	userID := context.MustGet("user_id").(int64)
	periodStart, renewDate := currentMonthlyCreditWindow(time.Now())

	usedCredits, err := repository.GetMonthlyAIAnalysisCount(userID, periodStart, renewDate)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch AI credit usage",
			"error":   err.Error(),
		})
		return
	}

	topUpCredits, err := repository.GetMonthlyAICreditTopUp(userID, periodStart)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch AI credit top-up",
			"error":   err.Error(),
		})
		return
	}

	totalCredits := getBaseMonthlyAICredits() + topUpCredits
	leftCredits := totalCredits - usedCredits
	if leftCredits < 0 {
		leftCredits = 0
	}

	usage := models.AICreditUsage{
		UsedCredit:  usedCredits,
		LeftCredit:  leftCredits,
		TotalCredit: totalCredits,
		RenewDate:   renewDate,
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "AI credit usage fetched successfully",
		"usage":   usage,
	})
}

func currentMonthlyCreditWindow(now time.Time) (time.Time, time.Time) {
	location := now.Location()
	periodStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, location)
	renewDate := periodStart.AddDate(0, 1, 0)
	return periodStart, renewDate
}

func getBaseMonthlyAICredits() int64 {
	// AI Credit
	creditValue := os.Getenv("MONTHLY_AI_ANALYSIS_CREDITS")
	if creditValue == "" {
		return 2
	}

	credits, err := strconv.ParseInt(creditValue, 10, 64)
	if err != nil || credits < 0 {
		return 2
	}

	return credits
}

// func updateProfile(context *gin.Context) {}
