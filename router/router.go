package router

import (
	"auth_service/handlers"

	"github.com/gorilla/mux"
)

// Функция для расстановки маршрутов, так как не считаю нужным,
// чтобы это было все в main, так как не оч хорошо для расширяемости
func InitRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/healthz", handlers.HealthCheck).Methods("GET")

	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/token", handlers.GenerateTokens).Methods("GET")
	authRouter.HandleFunc("/refresh", handlers.RefreshTokens).Methods("POST")

	return router
}
