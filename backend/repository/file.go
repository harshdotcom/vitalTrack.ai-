package repository

import (
	"vita-track-ai/database"
	"vita-track-ai/models"
)

func CreateFile(file *models.File) error {
	return database.DB.Create(file).Error
}

func GetFileByID(id string) (*models.File, error) {
	var file models.File
	err := database.DB.First(&file, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func GetFilesByUser(userID string) ([]models.File, error) {
	var files []models.File
	err := database.DB.Where("uploaded_by = ?", userID).Find(&files).Error
	if err != nil {
		return nil, err
	}
	return files, nil
}

func GetS3Key(fileId string) (string, error) {
	var s3Key string

	err := database.DB.
		Model(&models.File{}).
		Select("s3_key").
		Where("id = ?", fileId).
		Scan(&s3Key).
		Error

	if err != nil {
		return "", err
	}
	return s3Key, nil
}

func DeleteFile(id string, userID int64) error {

	return database.DB.
		Where("id = ? AND uploaded_by = ?", id, userID).
		Delete(&models.File{}).Error
}
