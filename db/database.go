package db

import (
	"auth_service/config"
	"auth_service/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Единственный экземпляр сессии подключения к базе данных
var DB *gorm.DB

// Фуннция подключения к базе данных
func InitDB() {
	var err error

	// При переменной окружения TEST=true, происходит подключение
	// к тестовой бд, но сейчас это просто пережиток прошлого и
	// логика уже другая, но может пригодиться позже еще
	if os.Getenv("TEST") == "true" {
		DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	} else {
		DB, err = gorm.Open(postgres.Open(config.AppConfig.DatabaseDSN), &gorm.Config{})
	}

	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	if err = DB.AutoMigrate(&models.RefreshToken{}); err != nil {
		log.Fatalf("Failed to migrate DB: %v", err)
	}

	log.Println("Database connected")
}
