package repository

import (
	"fmt"
	"strings"
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
		Model(&models.Document{}).
		Preload("File").
		Select(`
			documents.*,
			CASE
				WHEN medical_report_dbs.id IS NOT NULL THEN true
				ELSE false
			END AS analysis_generated
		`).
		Joins("LEFT JOIN medical_report_dbs ON medical_report_dbs.id = documents.id").
		Where("documents.id = ? AND documents.user_id = ?", id, userID).
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
		Model(&models.Document{}).
		Select(`
			documents.*,
			CASE
				WHEN medical_report_dbs.id IS NOT NULL THEN true
				ELSE false
			END AS analysis_generated
		`).
		Joins("LEFT JOIN medical_report_dbs ON medical_report_dbs.id = documents.id").
		Where("documents.user_id = ?", userID).
		Where("documents.document_date >= ? AND documents.document_date < ?", start, end)

	if req.Category != "" {
		query = query.Where("documents.category = ?", req.Category)
	}

	for _, tag := range req.Tags {
		query = query.Where("tags LIKE ?", "%"+tag+"%")
	}

	err := query.Find(&docs).Error
	return docs, err
}

func UpdateDocument(userID int64, documentId string, updateDocReq *models.UpdateDocumentRequest) error {
	query := "UPDATE documents SET "
	args := []interface{}{}
	i := 1

	if updateDocReq.Category != nil {
		query += fmt.Sprintf("category = $%d, ", i)
		args = append(args, *updateDocReq.Category)
		i++
	}

	if updateDocReq.DocumentName != nil {
		query += fmt.Sprintf("document_name = $%d, ", i)
		args = append(args, *updateDocReq.DocumentName)
		i++
	}

	if updateDocReq.Tags != nil {
		query += fmt.Sprintf("tags = $%d, ", i)
		args = append(args, *updateDocReq.Tags)
		i++
	}

	if updateDocReq.DocumentDate != nil {
		parsedDate, _ := time.Parse("2006-01-02", *updateDocReq.DocumentDate)
		query += fmt.Sprintf("document_date = $%d, ", i)
		args = append(args, parsedDate)
		i++
	}

	// ❗ If no fields provided, avoid executing invalid query
	if len(args) == 0 {
		return nil // or return error like "nothing to update"
	}

	// Remove trailing comma
	query = strings.TrimSuffix(query, ", ")

	// Add WHERE clause
	query += fmt.Sprintf(" WHERE user_id = $%d AND id = $%d", i, i+1)
	args = append(args, userID, documentId)

	// Execute
	return database.DB.Exec(query, args...).Error
}

func GetDocumentsInfiniteScroll(cursor *models.Cursor, limit int64, userID int64) ([]models.Document, *models.Cursor, error) {
	var documents []models.Document

	// query := database.DB.
	// 	Where("user_id = ?", userID).
	// 	Order("document_date DESC, id DESC").
	// 	Limit(int(limit))
	query := database.DB.
		Model(&models.Document{}).
		Where("user_id = ?", userID).
		Order("document_date DESC, id DESC").
		Limit(int(limit))

	// Cursor condition
	if cursor != nil {
		query = query.Where(
			"(document_date < ?) OR (document_date = ? AND id < ?)",
			cursor.CreatedAt, cursor.CreatedAt, cursor.ID,
		)
	}

	err := query.Find(&documents).Error
	if err != nil {
		return nil, nil, err
	}
	if len(documents) == 0 {
		return documents, nil, nil
	}

	last := documents[len(documents)-1]

	nextCursor := &models.Cursor{
		CreatedAt: last.DocumentDate,
		ID:        last.FileID,
	}

	return documents, nextCursor, nil
}
