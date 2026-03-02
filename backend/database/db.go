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
	const dbtype string = "postgres"
	const username string = "postgres"
	const password string = "0000"
	const dbhost string = "localhost"
	var port string = os.Getenv("DB_PORT")
	const dbname string = "vitadb"
	const security string = "sslmode=disable"
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
	CreateDocumentTable()
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

func CreateDocumentTable() {
	err := DB.AutoMigrate(&models.Document{})
	if err != nil {
		panic("failed to migrate Document table")
	}
}
