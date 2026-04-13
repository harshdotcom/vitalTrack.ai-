package service

import (
	"errors"
	"strings"
	"time"
	"vita-track-ai/models"
	"vita-track-ai/repository"
)

var MAX_ALLOWED_USER_STORAGE int64 = 104857600 //100MB in bytes

// var MAX_ALLOWED_USER_STORAGE int64 = 439870

func exceedStorageLimit(userId int64, fileSize int64) (bool, error) {
	userUsage, err := repository.GetCurrentStorageUsed(userId)
	if err != nil {
		return true, err
	}
	totalStorageUsed := userUsage.TotalStorageUsed

	if totalStorageUsed+fileSize > MAX_ALLOWED_USER_STORAGE {
		return true, nil
	}
	return false, nil
}

func ManageUserUpdateRequest(updateUserReq models.UpdateUserRequest, userId int64) (*models.User, error) {
	var shouldDeleteProfilePic bool
	userModel, err := repository.GetUserModelById(userId)

	if err != nil {
		return nil, err
	}

	if updateUserReq.DeleteProfilePic != nil {
		shouldDeleteProfilePic = *updateUserReq.DeleteProfilePic
	}

	fileHeader := updateUserReq.ProfilePic

	if fileHeader != nil && !shouldDeleteProfilePic {
		profileImage, err := UploadUserProfileImage(userId, fileHeader)
		if err != nil {
			return nil, err
		}
		userModel.ProfileImage = profileImage
		userModel.ProfilePic = nil
	} else if shouldDeleteProfilePic {
		if err := DeleteUserProfileImage(userId); err != nil {
			return nil, err
		}
		userModel.ProfilePic = nil
		userModel.ProfileImage = nil
		userModel.LegacyProfilePic = nil
	}

	if updateUserReq.Name != nil {
		trimmedName := strings.TrimSpace(*updateUserReq.Name)
		if trimmedName == "" {
			return nil, errors.New("name cannot be empty")
		}
		userModel.Name = trimmedName
	}

	if updateUserReq.DOB != nil {
		trimmedDOB := strings.TrimSpace(*updateUserReq.DOB)
		if trimmedDOB == "" {
			userModel.DOB = nil
		} else {
			parsedDOB, err := time.Parse("2006-01-02", trimmedDOB)
			if err != nil {
				return nil, errors.New("dob must be in YYYY-MM-DD format")
			}
			userModel.DOB = &parsedDOB
		}
	}

	if updateUserReq.Gender != nil {
		userModel.Gender = strings.TrimSpace(*updateUserReq.Gender)
	}

	err = repository.UpdateUser(&userModel)

	if err != nil {
		return nil, err
	}

	return &userModel, nil

}
