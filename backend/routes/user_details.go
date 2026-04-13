package routes

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"vita-track-ai/models"
	"vita-track-ai/repository"
	"vita-track-ai/service"

	"github.com/gin-gonic/gin"
)

// @Summary Get User Usage
// @Tags User-Details
// @Router /user-details/usage [get]
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

// @Summary Get AI Credit Usage
// @Tags User-Details
// @Router /user-details/ai-credits [get]
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

// @Summary Update User Profile
// @Tags User-Details
// @Accept json
// @Accept multipart/form-data
// @Produce json
// @Param name formData string false "Name"
// @Param delete_profile_pic formData string false "DeleteProfilePic"
// @Param dob formData string false "Date of Birth (YYYY-MM-DD)"
// @Param gender formData string false "Gender"
// @Param profile_pic formData file false "Profile Picture"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /user-details/update [patch]
func updateProfile(context *gin.Context) {

	var updateUserReq models.UpdateUserRequest
	contentType := context.ContentType()

	var err error
	if strings.HasPrefix(contentType, "application/json") {
		err = context.ShouldBindJSON(&updateUserReq)
	} else {
		err = context.ShouldBind(&updateUserReq)
	}

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to pass the values into the updateProfile",
			"error":   err.Error(),
		})

		return
	}

	userID := context.MustGet("user_id").(int64)

	userModel, err := service.ManageUserUpdateRequest(updateUserReq, userID)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "name cannot be empty") || strings.Contains(err.Error(), "dob must be in YYYY-MM-DD format") {
			statusCode = http.StatusBadRequest
		}

		context.JSON(statusCode, gin.H{
			"message": "There is a problem in updating the user",
			"error":   err.Error(),
		})
		return
	}

	userResponse, err := service.BuildUserResponse(*userModel)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "There is a problem in preparing the user response",
			"error":   err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Update successful",
		"user":    userResponse,
	})

}
