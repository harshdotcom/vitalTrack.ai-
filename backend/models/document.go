package models

import "time"

type Document struct {
	// ID         string    `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID            int64     `json:"user_id"`
	File              File      `gorm:"foreignKey:FileID;references:ID;constraint:OnDelete:CASCADE"`
	FileID            string    `json:"id" gorm:"column:id;primaryKey"`
	Category          string    `json:"category"`
	DocumentName      string    `json:"document_name"`
	Tags              string    `json:"tags"` // JSON string
	Status            string    `json:"status"`
	DocumentDate      time.Time `json:"document_date"`
	AnalysisGenerated bool      `json:"analysis_generated" gorm:"column:analysis_generated;->;-:migration"`
}

type CalendarRequest struct {
	Month    int      `json:"month"`
	Year     int      `json:"year"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
}

type UpdateDocumentRequest struct {
	Category     *string `form:"category"`
	DocumentName *string `form:"document_name"`
	Tags         *string `form:"tags"`
	DocumentDate *string `form:"document_date"`
}

type Cursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        string    `json:"id"`
}

type InfiniteScrollResponse struct {
	Data   []Document `json:"data"`
	Cursor string     `json:"cursor"`
}
