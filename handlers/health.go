package handlers

import "net/http"

// Обработчик для проверки состояния приложения
// Если приложение поднялось, то отправится StatusOK(200)
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
