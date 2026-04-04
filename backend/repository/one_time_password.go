package repository

import (
	"vita-track-ai/database"
	"vita-track-ai/models"

	"gorm.io/gorm/clause"
)

func SaveOTP(otpModel *models.OneTimePassword) error {
	return database.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		UpdateAll: true,
	}).Create(otpModel).Error
}

func GetOTPModelByEmail(email string) (models.OneTimePassword, error) {
	var otpModel models.OneTimePassword
	tx := database.DB.Where("email = ?", email).First(&otpModel)
	err := tx.Error

	return otpModel, err
}

func DeleteOTPByEmail(email string) error {
	return database.DB.Where("email = ?", email).Delete(&models.OneTimePassword{}).Error
}
