package service

import (
	"encoding/json"
	"net/http"
	"sort"
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
	Date         string   `json:"document_date"`
}

func CreateDocument(c *gin.Context) {
	var req CreateDocumentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	docs, err := repository.GetDocumentsByMonth(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	metrics, err := repository.GetHealthMetricsByMonth(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	days := make(map[string]gin.H)

	for _, doc := range docs {
		date := doc.DocumentDate.Format("2006-01-02")
		day := ensureCalendarDay(days, date)
		day["count"] = day["count"].(int) + 1
		day["documents"] = append(day["documents"].([]gin.H), gin.H{
			"entry_type":         "document",
			"user_id":            doc.UserID,
			"File":               doc.File,
			"id":                 doc.FileID,
			"category":           doc.Category,
			"document_name":      doc.DocumentName,
			"tags":               doc.Tags,
			"status":             doc.Status,
			"document_date":      doc.DocumentDate,
			"analysis_generated": doc.AnalysisGenerated,
		})
		days[date] = day
	}

	for _, metric := range metrics {
		date := metric.Timestamp.Format("2006-01-02")
		day := ensureCalendarDay(days, date)
		day["count"] = day["count"].(int) + 1
		day["documents"] = append(day["documents"].([]gin.H), gin.H{
			"entry_type":         "direct_entry",
			"id":                 metric.ID,
			"user_id":            metric.UploadedBy,
			"category":           "Direct Entry",
			"document_name":      metric.MetricLabel(),
			"status":             "logged",
			"document_date":      metric.Timestamp,
			"analysis_generated": false,
			"timestamp":          metric.Timestamp,
			"metric_type":        metric.MetricType(),
			"metric_label":       metric.MetricLabel(),
			"metric_summary":     metric.MetricSummary(),
			"heart_rate":         metric.HeartRate,
			"weight":             metric.Weight,
			"blood_pressure":     metric.BloodPressure,
			"blood_sugar":        metric.BloodSugar,
			"notes":              metric.Notes,
			"sleep_hours":        metric.SleepHours,
			"steps":              metric.Steps,
			"calories":           metric.Calories,
			"oxygen_level":       metric.OxygenLevel,
		})
		days[date] = day
	}

	for date, rawDay := range days {
		entries := rawDay["documents"].([]gin.H)
		sort.Slice(entries, func(i, j int) bool {
			return getEntryTime(entries[i]["document_date"]).Before(getEntryTime(entries[j]["document_date"]))
		})
		rawDay["documents"] = entries
		days[date] = rawDay
	}

	c.JSON(http.StatusOK, gin.H{
		"month": req.Month,
		"year":  req.Year,
		"days":  days,
	})
}

func UpdateDocument(userID int64, documentId string, updateDocReq *models.UpdateDocumentRequest) error {
	return repository.UpdateDocument(userID, documentId, updateDocReq)
}

func ensureCalendarDay(days map[string]gin.H, date string) gin.H {
	if day, exists := days[date]; exists {
		return day
	}

	day := gin.H{
		"count":     0,
		"documents": []gin.H{},
	}
	days[date] = day
	return day
}

func getEntryTime(value interface{}) time.Time {
	if ts, ok := value.(time.Time); ok {
		return ts
	}

	return time.Time{}
}
