package database

import (
	"os"
	"task5-pbi/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Config() {
	Connect()
	Migrate()
}

func Connect() {
	db, err := gorm.Open(postgres.Open(os.Getenv("DB")))
	if err != nil {
		panic(err.Error())
	}
	DB = db
}

func Migrate() {
	DB.AutoMigrate(&models.User{}, &models.Photo{})
}
