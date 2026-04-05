package models

import "time"

type Document struct {
	// ID         string    `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID            int64     `json:"user_id"`
	File              File      `gorm:"foreignKey:FileID;references:ID;constraint:OnDelete:CASCADE"`
	FileID            string    `json:"id" gorm:"column:id;primaryKey"`
	Category          string    `json:"category"`
	ReportType        string    `json:"report_type"`
	FileType          string    `json:"file_type"`
	Tags              string    `json:"tags"` // JSON string
	Status            string    `json:"status"`
	ReportDate        time.Time `json:"report_date"`
	AnalysisGenerated bool      `json:"analysis_generated" gorm:"column:analysis_generated;->;-:migration"`
}

type CalendarRequest struct {
	Month    int      `json:"month"`
	Year     int      `json:"year"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
}

type UpdateDocumentRequest struct {
	Category   *string `form:"category"`
	ReportType *string `form:"report_type"`
	FileType   *string `form:"file_type"`
	Tags       *string `form:"tags"`
}
