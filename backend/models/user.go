package models

import (
	"mime/multipart"
	"time"
)

type User struct {
	UserId           int64             `json:"user_id" gorm:"primaryKey;autoIncrement"`
	Email            string            `json:"email" binding:"required" gorm:"unique;not null"`
	Password         *string           `json:"-" binding:"required"`
	GoogleId         *string           `json:"google_id"`
	Name             string            `json:"name" binding:"required" gorm:"not null"`
	Age              *int              `json:"age"`
	Gender           string            `json:"gender"`
	LegacyProfilePic *string           `json:"-" gorm:"column:profile_pic"`
	ProfilePic       *string           `json:"profile_pic" gorm:"-"`
	ProfileImage     *UserProfileImage `json:"-" gorm:"foreignKey:UserID;references:UserId"`
	DOB              *time.Time        `json:"dob"`
	IsVerified       bool              `json:"is_verified" gorm:"default:false"`
	CreatedAt        time.Time         `json:"created_at"` //YYYY-MM-DD HH:MM:SS.microseconds stored in DB
	UpdatedAt        time.Time         `json:"updated_at"`
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
	Name             *string               `json:"name" form:"name"`
	DOB              *string               `json:"dob" form:"dob"`
	DeleteProfilePic *bool                 `json:"delete_profile_pic" form:"delete_profile_pic"`
	Gender           *string               `json:"gender" form:"gender"`
	ProfilePic       *multipart.FileHeader `json:"-" form:"profile_pic"`
}

// ===========================Forget Password Request=================================================
type ForgetPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	OTP         string `json:"otp" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type UserResponse struct {
	UserId     int64      `json:"user_id"`
	Email      string     `json:"email"`
	GoogleId   *string    `json:"google_id"`
	Name       string     `json:"name"`
	Age        *int       `json:"age"`
	Gender     string     `json:"gender"`
	ProfilePic *string    `json:"profile_pic"`
	DOB        *time.Time `json:"dob"`
	IsVerified bool       `json:"is_verified"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
