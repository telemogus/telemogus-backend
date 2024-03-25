package db

import (
	"github.com/dgb35/telemogus_backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var err error
	var dsn = "host=localhost user=postgres password=12345dgb dbname=telemogus port=5432 sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}
	DB.AutoMigrate(&models.User{}, &models.Chat{}, &models.Message{})
}
