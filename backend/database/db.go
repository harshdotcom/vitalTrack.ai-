package database

import (
	"os"
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

func createTables() {
	createUserTable()
	createFileTable()
	createMedicalRecordTable()
	createDocumentTable()
}

func createUserTable() {

	err := DB.AutoMigrate(&models.User{})
	if err != nil {
		panic("failed to migrate User table")
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
