package repository

import (
	"vita-track-ai/database"
	"vita-track-ai/models"
)

func SaveHealthMetric(metric *models.DailyHealthMetric) error {
	return database.DB.Create(metric).Error
}
