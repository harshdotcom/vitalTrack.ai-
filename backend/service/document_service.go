package service

import (
	"encoding/json"
	"net/http"
	"time"

	"vita-track-ai/models"
	"vita-track-ai/repository"

	"github.com/gin-gonic/gin"
)

type CreateDocumentRequest struct {
	FileID       string   `json:"file_id" binding:"required"`
	Category     string   `json:"category" binding:"required"`
	DocumentName string   `json:"document_name" binding:"required"`
	Tags         []string `json:"tags"`
	Date         string   `json:"document_date"` // ⭐ ADD THIS
}

func CreateDocument(c *gin.Context) {

	var req CreateDocumentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	// validate file exists
	_, err := repository.GetFileByID(req.FileID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid file_id",
		})
		return
	}

	tagsJSON, _ := json.Marshal(req.Tags)
	userID := c.MustGet("user_id").(int64)
	var parsedDate time.Time

	if req.Date != "" {
		var err error
		parsedDate, err = time.Parse("2006-01-02", req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid document_date format (use YYYY-MM-DD)",
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "document_date is required",
		})
		return
	}

	doc := models.Document{
		UserID:       userID,
		FileID:       req.FileID,
		Category:     req.Category,
		DocumentName: req.DocumentName,
		Tags:         string(tagsJSON),
		Status:       "uploaded",
		DocumentDate: parsedDate,
	}

	if err := repository.CreateDocument(&doc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create document",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"document_id": doc.FileID,
		"status":      doc.Status,
	})
}

func GetDocument(c *gin.Context) {

	id := c.Param("id")
	userID := c.MustGet("user_id").(int64)

	doc, err := repository.GetDocumentByID(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "document not found",
		})
		return
	}

	c.JSON(http.StatusOK, doc)
}

func DeleteDocument(c *gin.Context) {

	id := c.Param("id")
	userID := c.MustGet("user_id").(int64)

	err := repository.DeleteDocument(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "document not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "document deleted",
	})
}

func GetCalendarDocuments(c *gin.Context) {

	var req models.CalendarRequest
	userID := c.MustGet("user_id").(int64)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	docs, err := repository.GetDocumentsByMonth(userID, req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// ⭐ group by date (calendar structure)
	days := make(map[string]gin.H)

	for _, doc := range docs {

		date := doc.DocumentDate.Format("2006-01-02")

		if _, exists := days[date]; !exists {
			days[date] = gin.H{
				"count":     0,
				"documents": []models.Document{},
			}
		}

		entry := days[date]
		entry["count"] = entry["count"].(int) + 1
		entry["documents"] =
			append(entry["documents"].([]models.Document), doc)

		days[date] = entry
	}

	c.JSON(200, gin.H{
		"month": req.Month,
		"year":  req.Year,
		"days":  days,
	})
}

func UpdateDocument(userID int64, documentId string, updateDocReq *models.UpdateDocumentRequest) error {
	return repository.UpdateDocument(userID, documentId, updateDocReq)
}
