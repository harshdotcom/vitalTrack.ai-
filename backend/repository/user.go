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

func SaveUser(u *models.User) error {
	var err error

	if u.Password != nil {
		hashedPassword, err := utility.HashPassword(*u.Password)
		if err != nil {
			return err
		}

		u.Password = hashedPassword
	}
	err = database.DB.Create(u).Error
	return err
}

func ValidateCredential(u *models.User) error {

	enteredPassword := u.Password
	err := database.DB.Where("email = ?", u.Email).First(u).Error
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
