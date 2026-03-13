package service

import (
	"mime/multipart"
	"os"
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
