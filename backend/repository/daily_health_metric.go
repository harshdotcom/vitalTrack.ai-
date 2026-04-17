package repository

import (
	"time"
	"vita-track-ai/database"
	"vita-track-ai/models"
)

func SaveHealthMetric(metric *models.DailyHealthMetric) error {
	return database.DB.Create(metric).Error
}

func GetHealthMetricsByMonth(userID int64, req models.CalendarRequest) ([]models.DailyHealthMetric, error) {
	var metrics []models.DailyHealthMetric

	loc := time.Now().Location()
	start := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, loc)
	end := start.AddDate(0, 1, 0)

	err := database.DB.
		Model(&models.DailyHealthMetric{}).
		Where("uploaded_by = ?", userID).
		Where("timestamp >= ? AND timestamp < ?", start, end).
		Order("timestamp ASC").
		Find(&metrics).Error

	return metrics, err
}

func DeleteHealthMetric(id string, userID int64) error {
	return database.DB.
		Where("id = ? AND uploaded_by = ?", id, userID).
		Delete(&models.DailyHealthMetric{}).Error
}
