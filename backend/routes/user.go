package routes

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
	"vita-track-ai/models"
	"vita-track-ai/repository"
	"vita-track-ai/service"
	"vita-track-ai/utility"

	"github.com/gin-gonic/gin"
)

func isEmailDisabled() bool {
	return os.Getenv("DISABLE_EMAIL_FLOW") == "true"
}

// @Summary User Signup
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Param email formData string true "Email"
// @Param password formData string true "Password"
// @Param name formData string true "Name"
// @Param dob formData string false "Date of Birth (YYYY-MM-DD)"
// @Param gender formData string false "Gender"
// @Param profile_pic formData file false "Profile Picture"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/signup [post]
func signup(context *gin.Context) {
	var signupRequest models.SignupRequest
	err := context.ShouldBind(&signupRequest) //not with JSON as it will be a form data :)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to pass the values into the user object",
			"error":   err.Error(),
		})
		return
	}

	var user models.User
	user.Email = signupRequest.Email
	user.Password = &signupRequest.Password
	user.Name = signupRequest.Name
	user.DOB = signupRequest.DOB
	if signupRequest.Gender != nil {
		user.Gender = *signupRequest.Gender
	}

	// fileHeader := signupRequest.ProfilePic

	existingUser, err := repository.GetUserModelByEmail(user.Email)

	if err == nil {
		if existingUser.IsVerified {
			context.JSON(http.StatusConflict, gin.H{
				"message": "User Already Exists",
				"error":   errors.New("User Already Exists").Error(),
			})
			return
		}
		// User exists but is not verified — delete and allow re-signup
		if delErr := repository.DeleteUserByEmail(user.Email); delErr != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to reset unverified user",
				"error":   delErr.Error(),
			})
			return
		}
	}

	user.IsVerified = false
	userId, err := repository.SaveUser(&user)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "There is some problem saving the user",
			"error":   err.Error(),
		})
		return
	}

	if signupRequest.ProfilePic != nil {
		if _, err := service.UploadUserProfileImage(userId, signupRequest.ProfilePic); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"message": "User created but profile picture upload failed",
				"error":   err.Error(),
			})
			return
		}
	}

	// When email flow is disabled, auto-verify the user and skip OTP entirely.
	if isEmailDisabled() {
		fmt.Println("I am here, means email is disabled", os.Getenv("DISABLE_EMAIL_FLOW"))
		if verifyErr := repository.MakeUserVerified(user.Email); verifyErr != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"message": "User created but could not auto-verify",
				"error":   verifyErr.Error(),
			})
			return
		}
		context.JSON(http.StatusOK, gin.H{
			"message":       "Signup successful. Email verification is disabled — you can log in directly.",
			"email_enabled": false,
		})
		return
	}

	// 🔹 Generate OTP
	otpModel := utility.GenerateOTP()
	otpModel.Id = userId
	otpModel.Email = user.Email

	repository.SaveOTP(&otpModel)

	err = utility.SendEmail(user.Email, *otpModel.OTP)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Cannot send the email",
			"error":   err.Error(),
		})

		return

	}

	context.JSON(http.StatusOK, gin.H{
		"message":       "Signup successful. Please verify OTP sent to your email.",
		"email_enabled": true,
	})
}

// @Summary User Login
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login payload"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Success 200 {object} map[string]interface{}
// @Router /users/login [post]
func login(context *gin.Context) {

	var loginRequest models.LoginRequest
	err := context.ShouldBindBodyWithJSON(&loginRequest)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to pass the values into the user object",
			"error":   err.Error(),
		})

		return
	}

	fmt.Printf("[LOGIN] Attempt — email: %s\n", loginRequest.Email)

	var user models.User
	user.Email = loginRequest.Email
	user.Password = &loginRequest.Password

	err = repository.ValidateCredential(&user) //user struct gets updated with db values here

	if err != nil {
		fmt.Printf("[LOGIN] ValidateCredential failed — email: %s | error: %s\n", loginRequest.Email, err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"error":   err.Error(),
		})

		return
	}

	fmt.Printf("[LOGIN] Credential valid — email: %s | is_verified: %v\n", user.Email, user.IsVerified)

	token, err := utility.GenerateToken(user.Email, user.UserId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Some problem in generating jwt token",
		})

		return
	}

	userResponse, err := service.BuildUserResponse(user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Some problem preparing the user response",
			"error":   err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Successfully login",
		"token":   token,
		"user":    userResponse,
	})

}

// @Summary Verify OTP
// @Tags User
// @Router /users/verify-otp [post]
func verifyOTP(context *gin.Context) {

	var req struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	otpModel, err := repository.GetOTPModelByEmail(req.Email)

	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "OTP not found"})
		return
	}

	if otpModel.OTP == nil || *otpModel.OTP != req.OTP {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid OTP"})
		return
	}

	if otpModel.OTPExpiresAt.Before(time.Now()) {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "OTP expired"})
		return
	}

	err = repository.MakeUserVerified(req.Email)

	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "There is some problem in verifying the user"})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Email verified successfully",
	})
}

// @Summary Forgot Password
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.ForgetPasswordRequest true "Forget Password payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/forgot-password [post]
func forgotPassword(context *gin.Context) {
	if isEmailDisabled() {
		context.JSON(http.StatusServiceUnavailable, gin.H{
			"message": "Password reset via email is currently disabled.",
		})
		return
	}

	var forgetPasswordRequest models.ForgetPasswordRequest
	err := context.ShouldBindJSON(&forgetPasswordRequest)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to pass the values into the ForgetPassword object",
			"error":   err.Error(),
		})
		return
	}

	user, err := repository.GetUserModelByEmail(forgetPasswordRequest.Email)

	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"message": "Unable to find the user with this email",
			"error":   err.Error(),
		})
		return
	}

	otpModel := utility.GenerateOTP()
	otpModel.Id = user.UserId
	otpModel.Email = user.Email

	err = repository.SaveOTP(&otpModel)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to generate OTP",
			"error":   err.Error(),
		})
		return
	}

	go utility.SendEmail(user.Email, *otpModel.OTP)

	context.JSON(http.StatusOK, gin.H{
		"message": "Please verify OTP sent to your email to change your password.",
	})
}

// @Summary Reset Password
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.ResetPasswordRequest true "Reset Password"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/reset-password [post]
func resetPassword(context *gin.Context) {
	var req models.ResetPasswordRequest

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	otpModel, err := repository.GetOTPModelByEmail(req.Email)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "OTP not found"})
		return
	}

	if otpModel.OTP == nil || *otpModel.OTP != req.OTP {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid OTP"})
		return
	}

	if otpModel.OTPExpiresAt.Before(time.Now()) {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "OTP expired"})
		return
	}

	hashedPassword, err := utility.HashPassword(req.NewPassword)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not hash password"})
		return
	}

	err = repository.UpdatePassword(req.Email, *hashedPassword)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update password"})
		return
	}

	repository.DeleteOTPByEmail(req.Email)

	context.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully",
	})
}

// @Summary Resend OTP
// @Tags User
// @Accept json
// @Produce json
// @Param request body object true "Email payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/resend-otp [post]
func resendOTP(context *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required"`
	}

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Email is required"})
		return
	}

	user, err := repository.GetUserModelByEmail(req.Email)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "No account found with this email"})
		return
	}

	if user.IsVerified {
		context.JSON(http.StatusBadRequest, gin.H{"message": "This account is already verified"})
		return
	}

	// Delete existing OTP (if any) and generate a fresh one
	repository.DeleteOTPByEmail(req.Email)

	otpModel := utility.GenerateOTP()
	otpModel.Id = user.UserId
	otpModel.Email = user.Email

	if err := repository.SaveOTP(&otpModel); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate OTP"})
		return
	}

	go utility.SendEmail(user.Email, *otpModel.OTP)

	context.JSON(http.StatusOK, gin.H{
		"message": "A new OTP has been sent to your email.",
	})
}

// @Summary Google Login
// @Tags User
// @Router /users/google [post]
func googleLogin(context *gin.Context) {
	var req models.GoogleLoginRequest
	err := context.ShouldBindJSON(&req)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to pass the values into the GoogleLoginRequest",
			"error":   err.Error(),
		})

		return
	}

	userResponse, token, err := service.AuthenticateGoogleUser(req.Token)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := "Google login failed"

		switch {
		case errors.Is(err, service.ErrInvalidGoogleToken), errors.Is(err, service.ErrGoogleEmailNotVerified), errors.Is(err, service.ErrGoogleEmailUnavailable):
			statusCode = http.StatusUnauthorized
			message = err.Error()
		case errors.Is(err, service.ErrGoogleTokenNotConfigured):
			statusCode = http.StatusServiceUnavailable
			message = err.Error()
		}

		context.JSON(statusCode, gin.H{
			"message": message,
			"error":   err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Google login successful",
		"user":    userResponse,
		"token":   token,
	})

}
