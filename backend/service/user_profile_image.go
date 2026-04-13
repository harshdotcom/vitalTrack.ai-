package service

import (
	"errors"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"vita-track-ai/models"
	"vita-track-ai/repository"

	"github.com/google/uuid"
)

var allowedProfileImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

const maxAllowedProfileImageSize int64 = 5 * 1024 * 1024

func UploadUserProfileImage(userID int64, fileHeader *multipart.FileHeader) (*models.UserProfileImage, error) {
	if fileHeader == nil {
		return nil, errors.New("profile picture is required")
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !allowedProfileImageExtensions[ext] {
		return nil, errors.New("profile picture must be a jpg, jpeg, png, or webp file")
	}

	if fileHeader.Size > maxAllowedProfileImageSize {
		return nil, errors.New("profile picture size must not exceed 5MB")
	}

	bucket := os.Getenv("AWS_BUCKET_NAME")
	if strings.TrimSpace(bucket) == "" {
		return nil, errors.New("aws bucket name is not configured")
	}

	existingImage, err := repository.GetUserProfileImageByUserID(userID)
	if err != nil {
		return nil, err
	}

	objectKey := buildProfileImageObjectKey(userID, ext)
	if err := UploadToS3(fileHeader, objectKey, bucket); err != nil {
		return nil, err
	}

	profileImage := &models.UserProfileImage{
		UserID:       userID,
		Bucket:       bucket,
		ObjectKey:    objectKey,
		OriginalName: fileHeader.Filename,
		MimeType:     fileHeader.Header.Get("Content-Type"),
		FileSize:     fileHeader.Size,
	}

	if err := repository.UpsertUserProfileImage(profileImage); err != nil {
		_ = DeleteFileFromS3(objectKey)
		return nil, err
	}

	if existingImage != nil && existingImage.ObjectKey != objectKey {
		_ = DeleteFileFromS3(existingImage.ObjectKey)
	}

	return profileImage, nil
}

func DeleteUserProfileImage(userID int64) error {
	existingImage, err := repository.GetUserProfileImageByUserID(userID)
	if err != nil {
		return err
	}

	if existingImage == nil {
		return nil
	}

	if err := DeleteFileFromS3(existingImage.ObjectKey); err != nil {
		return err
	}

	return repository.DeleteUserProfileImageByUserID(userID)
}

func buildProfileImageObjectKey(userID int64, ext string) string {
	return "profile-images/" + strconv.FormatInt(userID, 10) + "/" + uuid.NewString() + ext
}
