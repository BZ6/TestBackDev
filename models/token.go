package models

import (
	"gorm.io/gorm"
)

// Модель хранения в бд для refresh токена
type RefreshToken struct {
	gorm.Model
	UserID     string `gorm:"index"`
	TokenHash  string
	AccessUUID string `gorm:"index"`
	Used       bool
	IP         string
}
