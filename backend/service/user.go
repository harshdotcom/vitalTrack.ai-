package service

import (
	"mime/multipart"
	"os"
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

func UploadProfilePicToS3(fileHeader *multipart.FileHeader, email string) (string, error) {
	storageKey := "profile-pics-" + email
	err := UploadToS3(fileHeader, storageKey, os.Getenv("AWS_BUCKET_NAME"))

	return storageKey, err

}

func ManageUserUpdateRequest(updateUserReq models.UpdateUserRequest, userId int64) (*models.User, error) {
	var shouldDeleteProfilePic bool = true
	userModel, err := repository.GetUserModelById(userId)

	if err != nil {
		return nil, err
	}

	if updateUserReq.DeleteProfilePic == nil || *updateUserReq.DeleteProfilePic == "false" {
		shouldDeleteProfilePic = false
	}

	fileHeader := updateUserReq.ProfilePic

	if fileHeader != nil && shouldDeleteProfilePic == false {

		if userModel.ProfilePic != nil {
			err = DeleteFileFromS3(*userModel.ProfilePic)

			if err != nil {
				return nil, err
			}
		}

		storageKey := "profile-pics-" + userModel.Email
		err = UploadToS3(fileHeader, storageKey, os.Getenv("AWS_BUCKET_NAME"))
		if err != nil {
			return nil, err
		}

		// fmt.Printf("The user profile pic is %v", userModel.ProfilePic)
		userModel.ProfilePic = &storageKey
	} else if shouldDeleteProfilePic == true {

		if userModel.ProfilePic != nil {
			err = DeleteFileFromS3(*userModel.ProfilePic)
			if err != nil {
				return nil, err
			}
		}

		userModel.ProfilePic = nil
	}

	if updateUserReq.Name != nil {
		userModel.Name = *updateUserReq.Name
	}

	if updateUserReq.DOB != nil {
		userModel.DOB = updateUserReq.DOB
	}

	return &userModel, nil

}
