package handlers

import (
	"auth_service/db"
	"auth_service/models"
	"auth_service/utils"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// Структура тела ответа с токенами
// Она же у меня является телом запроса для refresh
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Обработчик для генерации токенов по user_id
func GenerateTokens(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	clientIP := r.RemoteAddr

	accessToken, err := utils.GenerateJWT(userID, clientIP)
	if err != nil {
		http.Error(w, "Error creating JWT", http.StatusInternalServerError)
		return
	}

	refreshTokenRaw := utils.GenerateRandomToken()
	refreshTokenHash, err := utils.HashToken(refreshTokenRaw)
	if err != nil {
		http.Error(w, "Could not hash refresh token", http.StatusInternalServerError)
		return
	}

	rt := &models.RefreshToken{
		UserID:     userID,
		TokenHash:  refreshTokenHash,
		AccessUUID: utils.ExtractTokenID(accessToken),
		IP:         clientIP,
	}
	if err := db.DB.Create(rt).Error; err != nil {
		http.Error(w, "Database error on saving refresh token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenRaw,
	})
}

// Обработчик для обновления токенов
// Подумал что user_id тоже передавать было бы странно, а вот access токен
// было бы хорошо обрабатывать, но тяжелее тестировать когда он спрятан в
// заголовке запроса, поэтому он также как и refresh передается в теле
func RefreshTokens(w http.ResponseWriter, r *http.Request) {
	var request TokenResponse
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	clientIP := r.RemoteAddr

	claims, err := utils.ValidateJWT(request.AccessToken)
	if err != nil {
		http.Error(w, "Invalid access token", http.StatusUnauthorized)
		return
	}

	var storedRT models.RefreshToken
	result := db.DB.Where("access_uuid = ? AND used = ?", claims.ID, false).First(&storedRT)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Refresh token not found or already used", http.StatusUnauthorized)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	if !utils.CheckTokenHash(request.RefreshToken, storedRT.TokenHash) {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	if storedRT.IP != clientIP {
		utils.SendWarningEmail(storedRT.UserID, clientIP)
	}

	storedRT.Used = true
	if err := db.DB.Save(&storedRT).Error; err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	newAccessToken, err := utils.GenerateJWT(storedRT.UserID, clientIP)
	if err != nil {
		http.Error(w, "Could not generate new access token", http.StatusInternalServerError)
		return
	}

	newRefreshTokenRaw := utils.GenerateRandomToken()
	newRefreshTokenHash, err := utils.HashToken(newRefreshTokenRaw)
	if err != nil {
		http.Error(w, "Could not hash new refresh token", http.StatusInternalServerError)
		return
	}

	newAccessTokenUUID := utils.ExtractTokenID(newAccessToken)

	newRT := &models.RefreshToken{
		UserID:     storedRT.UserID,
		TokenHash:  newRefreshTokenHash,
		AccessUUID: newAccessTokenUUID,
		IP:         clientIP,
	}
	if err := db.DB.Create(newRT).Error; err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshTokenRaw,
	})
}
