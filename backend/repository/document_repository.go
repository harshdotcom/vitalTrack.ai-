package repository

import (
	"time"
	"vita-track-ai/database"
	"vita-track-ai/models"
)

func CreateDocument(doc *models.Document) error {
	return database.DB.Create(doc).Error
}

func GetDocumentByID(id string, userID int64) (*models.Document, error) {

	var doc models.Document

	err := database.DB.
		Where("id = ? AND user_id = ?", id, userID).
		First(&doc).Error

	return &doc, err
}

func DeleteDocument(id string, userID int64) error {

	return database.DB.
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.Document{}).Error
}

func GetDocumentsByMonth(userID int64, req models.CalendarRequest) ([]models.Document, error) {

	var docs []models.Document

	loc := time.Now().Location()

	start := time.Date(
		req.Year,
		time.Month(req.Month),
		1,
		0, 0, 0, 0,
		loc,
	)

	end := start.AddDate(0, 1, 0)

	query := database.DB.
		Where("user_id = ?", userID).
		Where("report_date >= ? AND report_date < ?", start, end)

	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}

	for _, tag := range req.Tags {
		query = query.Where("tags LIKE ?", "%"+tag+"%")
	}

	err := query.Find(&docs).Error
	return docs, err
}
