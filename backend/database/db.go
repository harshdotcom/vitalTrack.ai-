package database

import (
	"context"
	"os"
	"path/filepath"
	"vita-track-ai/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func generateDSN() string {
	// var dbtype string = os.Getenv("DB_TYPE")
	var username string = os.Getenv("DB_USER")
	var password string = os.Getenv("DB_PASSWORD")
	var dbhost string = os.Getenv("DB_HOST")
	var port string = os.Getenv("DB_PORT")
	var dbname string = os.Getenv("DB_NAME")
	var security string = os.Getenv("DB_SECURITY")
	var dsn string = "host=" + dbhost + " user=" + username + " password=" + password + " dbname=" + dbname + " port=" + port + " " + security
	return dsn
}

func Init() {
	dsn := generateDSN()
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	createTables()
}

func ReadSQLFile(path string) (string, error) {
	b, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return "", err
	}

	query := string(b)

	return query, nil
}

func RunSQLFile(path string) error {
	query, err := ReadSQLFile(path)

	if err != nil {
		return err
	}

	return DB.WithContext(context.Background()).Exec(query).Error
}

func createTables() {
	createUserTable()
	createUserProfileImageTable()
	createUserAICreditGrantTable()
	createFileTable()
	createMedicalRecordTable()
	createDocumentTable()
	createUserStorageMV()
	createOTPTable()
	createHealthMetricTable()
}

func createUserTable() {

	err := DB.AutoMigrate(&models.User{})
	if err != nil {
		panic("failed to migrate User table")
	}
}

func createUserAICreditGrantTable() {
	err := DB.AutoMigrate(&models.UserAICreditGrant{})
	if err != nil {
		panic("failed to migrate UserAICreditGrant table")
	}
}

func createUserProfileImageTable() {
	err := DB.AutoMigrate(&models.UserProfileImage{})
	if err != nil {
		panic("failed to migrate UserProfileImage table")
	}
}

func createFileTable() {
	err := DB.AutoMigrate(&models.File{})
	if err != nil {
		panic("failed to migrate File table")
	}
}

func createDocumentTable() {
	err := DB.AutoMigrate(&models.Document{})
	if err != nil {
		panic("failed to migrate Document table")
	}
}

func createMedicalRecordTable() {
	err := DB.AutoMigrate(&models.MedicalReportDB{})
	if err != nil {
		panic("failed to migrate Document table")
	}
}

func createHealthMetricTable() {
	err := DB.AutoMigrate(&models.DailyHealthMetric{})
	if err != nil {
		panic("failed to migrate Health Metric table")
	}
}

func createUserStorageMV() error {

	query, err := ReadSQLFile("sql/USER_STORAGE.sql")
	if err != nil {
		return err
	}

	return DB.WithContext(context.Background()).Exec(query).Error
}

func createOTPTable() {
	err := DB.AutoMigrate(&models.OneTimePassword{})
	if err != nil {
		panic("failed to migrate OTP table")
	}
}
