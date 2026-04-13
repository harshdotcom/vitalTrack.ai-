package models

import "time"

type UserProfileImage struct {
	ID           int64  `gorm:"primaryKey;autoIncrement"`
	UserID       int64  `gorm:"not null;uniqueIndex"`
	Bucket       string `gorm:"not null"`
	ObjectKey    string `gorm:"not null;uniqueIndex"`
	OriginalName string
	MimeType     string
	FileSize     int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
