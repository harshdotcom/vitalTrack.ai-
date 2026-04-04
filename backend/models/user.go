package models

import (
	"mime/multipart"
	"time"
)

type User struct {
	UserId       int64   `json:"user_id" gorm:"primaryKey;autoIncrement"`
	Email        string  `json:"email" binding:"required" gorm:"unique;not null"`
	Password     *string `json:"password" binding:"required"`
	GoogleId     *string
	Name         string     `json:"name" binding:"required" gorm:"not null"`
	Gender       string     `json:"gender"`
	ProfilePic   *string    `json:"profile_pic"`
	DOB          *time.Time `json:"dob"`
	IsVerified   bool       `json:"is_verified" gorm:"default:false"`
	OTP          *string    `json:"-"` // hide from API
	OTPExpiresAt *time.Time

	CreatedAt time.Time //YYYY-MM-DD HH:MM:SS.microseconds stored in DB
	UpdatedAt time.Time
}

type UserUsage struct {
	UserID           int64 `gorm:"column:user_id;"`
	TotalStorageUsed int64 `gorm:"column:total_storage_used;"`
}

// =========================Signup Request==========================================================

type SignupRequest struct {
	Email      string                `form:"email" binding:"required,email"`
	Password   string                `form:"password" binding:"required,min=8"`
	Name       string                `form:"name" binding:"required"`
	DOB        *time.Time            `form:"dob" time_format:"2006-01-02"`
	Gender     *string               `form:"gender"`
	ProfilePic *multipart.FileHeader `form:"profile_pic"`
}

// =========================Login Requests==========================================================

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type GoogleLoginRequest struct {
	Token string `json:"token" binding:"required"`
}

// ===========================Update User Request=================================================
type UpdateUserRequest struct {
	Name             *string               `form:"name"`
	DOB              *time.Time            `form:"dob" time_format:"2006-01-02"`
	DeleteProfilePic *string               `form:"delete_profile_pic"`
	Gender           *string               `form:"gender"`
	ProfilePic       *multipart.FileHeader `form:"profile_pic"`
}
