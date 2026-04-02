package repository

import (
	"time"
	"vita-track-ai/database"
	"vita-track-ai/models"
)

func GetMonthlyAIAnalysisCount(userID int64, periodStart time.Time, periodEnd time.Time) (int64, error) {
	var count int64

	err := database.DB.
		Model(&models.MedicalReportDB{}).
		Joins("JOIN files ON files.id = medical_report_dbs.id").
		Where("files.uploaded_by = ? AND medical_report_dbs.created_at >= ? AND medical_report_dbs.created_at < ?", userID, periodStart, periodEnd).
		Count(&count).
		Error

	return count, err
}

func GetMonthlyAICreditTopUp(userID int64, periodStart time.Time) (int64, error) {
	var total int64

	err := database.DB.
		Model(&models.UserAICreditGrant{}).
		Where("user_id = ? AND effective_month = ?", userID, periodStart).
		Select("COALESCE(SUM(credits), 0)").
		Scan(&total).
		Error

	return total, err
}
