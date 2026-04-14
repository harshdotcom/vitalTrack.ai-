package repository

import (
	"errors"
	"fmt"
	"os"
	"time"
	"vita-track-ai/database"
	"vita-track-ai/models"
	"vita-track-ai/utility"

	"gorm.io/gorm"
)

func GetUserModelByEmail(email string) (models.User, error) {
	var user models.User
	tx := database.DB.Preload("ProfileImage").Where("email = ?", email).First(&user)
	err := tx.Error

	return user, err
}

func GetUserModelById(id int64) (models.User, error) {
	var user models.User
	tx := database.DB.Preload("ProfileImage").Where("user_id = ?", id).First(&user)
	err := tx.Error

	return user, err
}

func SaveUser(u *models.User) (int64, error) {
	var err error

	if u.Password != nil {
		hashedPassword, err := utility.HashPassword(*u.Password)
		if err != nil {
			return -1, err
		}

		u.Password = hashedPassword
	}
	err = database.DB.Create(u).Error
	return u.UserId, err
}

func ValidateCredential(u *models.User) error {

	enteredPassword := u.Password
	fmt.Println(enteredPassword, "printing the added password")
	tx := database.DB.Preload("ProfileImage").Where("email = ?", u.Email).First(u)
	err := tx.Error
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("User with this email does not exists")
		}

		return err
	}

	isValid := utility.ValidateEnteredPassword(*enteredPassword, *u.Password)

	if !isValid {
		return errors.New("Entered Password is Incorrect")
	}

	if !u.IsVerified && os.Getenv("DISABLE_EMAIL_FLOW") != "true" {
		return errors.New("User is not verified")
	}

	return nil
}

func UpdateGoogleId(u *models.User) error {
	tx := database.DB.Model(u).Where("email = ?", u.Email)
	err := tx.Update("google_id", u.GoogleId).Error

	if err != nil {
		return err
	}

	// After the update, get the User ID (since GORM doesn't support RETURNING like raw SQL)
	err = database.DB.Preload("ProfileImage").Where("email = ?", u.Email).First(u).Error

	return err
}

// func UpdateUser(u *models.User) error {
// 	return database.DB.Model(&models.User{}).
// 		Where("user_id = ?", u.UserId).
// 		Updates(u).Error
// }

func UpdateUser(userModel *models.User) error {
	now := time.Now()
	updates := map[string]interface{}{
		"name":       userModel.Name,
		"dob":        userModel.DOB,
		"gender":     userModel.Gender,
		"updated_at": now,
	}

	if err := database.DB.Model(&models.User{}).
		Where("user_id = ?", userModel.UserId).
		Updates(updates).Error; err != nil {
		return err
	}

	userModel.UpdatedAt = now
	return nil
}

func DeleteUserByEmail(email string) error {
	return database.DB.Where("email = ?", email).Delete(&models.User{}).Error
}

func GetCurrentStorageUsed(userId int64) (*models.UserUsage, error) {
	var userUsage models.UserUsage
	userUsage.UserID = userId
	tx := database.DB.Table("user_usage").
		Select("total_storage_used").
		Where("user_id = ?", userId).
		Scan(&userUsage.TotalStorageUsed)

	return &userUsage, tx.Error
}

func MakeUserVerified(email string) error {
	return database.DB.Model(&models.User{}).
		Where("email = ?", email).
		Update("is_verified", true).Error
}

func UpdatePassword(email string, newHashedPassword string) error {
	return database.DB.Model(&models.User{}).
		Where("email = ?", email).
		Update("password", newHashedPassword).Error
}
