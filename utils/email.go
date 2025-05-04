package utils

import (
	"log"
)

// Фиктивная функция для отправки писем с предупреждением
func SendWarningEmail(userID, newIP string) {
	log.Printf("[WARNING EMAIL] Suspicious IP detected for user %s: %s\n", userID, newIP)
}
