package repository

import (
	"errors"
	"vita-track-ai/database"
	"vita-track-ai/models"

	"gorm.io/gorm"
)

func GetUserProfileImageByUserID(userID int64) (*models.UserProfileImage, error) {
	var profileImage models.UserProfileImage
	err := database.DB.Where("user_id = ?", userID).First(&profileImage).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &profileImage, nil
}

func UpsertUserProfileImage(profileImage *models.UserProfileImage) error {
	return database.DB.
		Where("user_id = ?", profileImage.UserID).
		Assign(profileImage).
		FirstOrCreate(profileImage).Error
}

func DeleteUserProfileImageByUserID(userID int64) error {
	return database.DB.Where("user_id = ?", userID).Delete(&models.UserProfileImage{}).Error
}
