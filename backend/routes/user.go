package routes

import (
	"database/sql"
	"errors"
	"net/http"
	"vita-track-ai/models"
	"vita-track-ai/repository"
	"vita-track-ai/service"
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
	user.Gender = signupRequest.Gender
	fileHeader := signupRequest.ProfilePic

	_, err = repository.GetUserModelByEmail(user.Email)

	if err == nil {
		context.JSON(http.StatusConflict, gin.H{
			"message": "User Already Exists",
			"error":   errors.New("User Already Exists").Error(),
		})
		return
	}

	if fileHeader != nil {
		storageKey, err := service.UploadProfilePicToS3(fileHeader, user.Email)

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to upload profile picture",
				"error":   err.Error(),
			})
			return
		}
		user.ProfilePic = &storageKey
	}

	err = repository.SaveUser(&user)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "There is some problem saving the user",
			"error":   err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{"users": user})
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

func getUserUsage(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)
	userUsage, err := repository.GetCurrentStorageUsed(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Some problem with fetching user storage usage",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User storage usage fetched successfully",
		"data":    userUsage,
	})
}

func updateProfile(context *gin.Context) {

	var updateUserReq models.UpdateUserRequest
	err := context.ShouldBind(&updateUserReq)

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
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "There is a problem in updating the user",
			"error":   err.Error(),
		})
		return
	}

	err = repository.UpdateUser(userModel)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "Some problem with database in updating the user",
			"error":   err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Update successful",
		"user":    userModel,
	})

}
