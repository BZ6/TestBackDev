package handlers

import "net/http"

// Обработчик для проверки состояния приложения
// Если приложение поднялось, то отправится StatusNoContent(204)
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
