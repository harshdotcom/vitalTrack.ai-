package routes

import (
	"database/sql"
	"errors"
	"net/http"
	"time"
	"vita-track-ai/models"
	"vita-track-ai/repository"
	"vita-track-ai/utility"

	"github.com/gin-gonic/gin"
)

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
	// user.Gender = signupRequest.Gender
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

	// 🔹 Generate OTP
	otp := utility.GenerateOTP()
	expiry := time.Now().Add(5 * time.Minute)

	user.IsVerified = false
	user.OTP = &otp
	user.OTPExpiresAt = &expiry

	err = repository.SaveUser(&user)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "There is some problem saving the user",
			"error":   err.Error(),
		})
		return
	}

	go utility.SendEmail(user.Email, otp)

	context.JSON(http.StatusOK, gin.H{
		"message": "Signup successful. Please verify OTP sent to your email.",
	})
}

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

	var user models.User
	user.Email = loginRequest.Email
	user.Password = &loginRequest.Password

	err = repository.ValidateCredential(&user) //user struct gets updated with db values here

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Password",
			"error":   err.Error(),
		})

		return
	}

	token, err := utility.GenerateToken(user.Email, user.UserId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Some problem in generating jwt token",
		})

		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Successfully login",
		"token":   token,
		"user":    user,
	})

}

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

	payload, err := utility.VerifyGoogkeIDTokenAndGetPayLoad(req.Token)

	if err != nil {
		context.JSON(401, gin.H{
			"message": "Invalid google token",
			"error":   err.Error(),
		})

		return
	}

	claims := payload.Claims

	email := utility.GetClaim("email", claims)
	name := utility.GetClaim("name", claims)
	picture := utility.GetClaim("picture", claims)
	googleId := payload.Subject

	var userModel models.User
	userModel, err = repository.GetUserModelByEmail(email)

	if err == sql.ErrNoRows {
		userModel.Email = email
		userModel.Name = name
		userModel.ProfilePic = &picture
		userModel.GoogleId = &googleId
		err = repository.SaveUser(&userModel)

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"message": "Problem in saving the user",
				"error":   err.Error(),
			})
			return
		}

	} else if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Some issue with the database",
			"error":   err.Error(),
		})

		return
	}

	if userModel.GoogleId == nil {
		userModel.GoogleId = &googleId
		err = repository.UpdateGoogleId(&userModel)

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"message": "Some problem in updating the google id",
				"error":   err.Error(),
			})

			return
		}
	}

	token, err := utility.GenerateToken(userModel.Email, userModel.UserId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Some problem with generating the token",
			"error":   err.Error(),
		})

		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Google login successful",
		"user":    userModel,
		"token":   token,
	})

}

func verifyOTP(context *gin.Context) {

	var req struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := repository.GetUserModelByEmail(req.Email)

	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	if user.OTP == nil || *user.OTP != req.OTP {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid OTP"})
		return
	}

	if user.OTPExpiresAt.Before(time.Now()) {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "OTP expired"})
		return
	}

	user.IsVerified = true
	user.OTP = nil
	user.OTPExpiresAt = nil

	repository.UpdateUser(&user)

	context.JSON(http.StatusOK, gin.H{
		"message": "Email verified successfully",
	})
}
