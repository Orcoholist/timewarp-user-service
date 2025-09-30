package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `gorm:"uniqueIndex" json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"` // admin, traveler
	RegisteredAt time.Time `json:"registered_at"`
	LastLogin    time.Time `json:"last_login"`
	TimeTravels  []TimeTravel `gorm:"foreignKey:UserID"`
}
