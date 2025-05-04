package utils

import (
	"auth_service/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Структура с описанием, того что я хочу поместить в JWT
type Claims struct {
	UserID string `json:"user_id"`
	IP     string `json:"ip"`
	jwt.RegisteredClaims
}

// Функция для генерации JWT
func GenerateJWT(userID, ip string) (string, error) {
	secret := []byte(config.AppConfig.JWTSecret)
	expMinutes := time.Minute * time.Duration(config.AppConfig.TokenExpiry)

	claims := Claims{
		UserID: userID,
		IP:     ip,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expMinutes)),
			ID:        uuid.NewString(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(secret)
}

// Функция для проверки валидности JWT с моей структурой
func ValidateJWT(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWTSecret), nil
	}, jwt.WithValidMethods([]string{"HS512"}))

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, err
	}
	return claims, nil
}

// Функция извлечения поля ID из JWT
// Неправильно, что не обрабатываю ошибки, но я еще не придумал,
// как на обработчике на это реагировать
// TODO: Сделать обработку ошибок!
func ExtractTokenID(tokenStr string) string {
	claims, _ := ValidateJWT(tokenStr)
	return claims.ID
}
