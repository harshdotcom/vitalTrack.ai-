package repository

import (
	"errors"
	"vita-track-ai/database"
	"vita-track-ai/models"
	"vita-track-ai/utility"

	"gorm.io/gorm"
)

func GetUserModelByEmail(email string) (models.User, error) {
	var user models.User
	tx := database.DB.Where("email = ?", email).First(&user)
	err := tx.Error

	return user, err
}

func GetUserModelById(id int64) (models.User, error) {
	var user models.User
	tx := database.DB.Where("user_id = ?", id).First(&user)
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
	tx := database.DB.Where("email = ?", u.Email).First(u)
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

	if u.IsVerified == false {
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
	err = database.DB.Where("email = ?", u.Email).First(u).Error

	return err
}

// func UpdateUser(u *models.User) error {
// 	return database.DB.Model(&models.User{}).
// 		Where("user_id = ?", u.UserId).
// 		Updates(u).Error
// }

func UpdateUser(userModel *models.User) error {
	query, err := database.ReadSQLFile("sql/UPDATE_USER.sql")
	if err != nil {
		return err
	}

	// if userModel.DOB != nil {
	// 	dobStr := userModel.DOB.Format("2006-01-02")
	// 	tempDOB, _ := time.Parse("2006-01-02", dobStr)
	// 	userModel.DOB = &tempDOB
	// }

	return database.DB.Exec(query, userModel.Name, userModel.DOB, userModel.Gender, userModel.ProfilePic, userModel.IsVerified, userModel.UserId).Error
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
