package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// В данной структуре хранятся переменные окружения
type Config struct {
	DatabaseDSN   string
	JWTSecret     string
	TokenExpiry   int64
	RefreshExpiry int64
}

// Единственный экземпляр конфига приложения
var AppConfig Config

// Внутрення функция, чтобы не проверять вручную
func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Missing %s environment variable", key)
	}

	return value
}

// Фуннция загружающая нужные переменные в конфиг
func InitConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	AppConfig.DatabaseDSN = getEnv("DATABASE_DSN")
	AppConfig.JWTSecret = getEnv("JWT_SECRET")
	tokenExpiry, err := strconv.ParseInt(getEnv("TOKEN_EXPIRY"), 10, 64)
	if err != nil {
		log.Fatalf("Invalid TOKEN_EXPIRY value: %v", err)
	}
	AppConfig.TokenExpiry = tokenExpiry
}
